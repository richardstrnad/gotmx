package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

var version = ""
var signer = NewSigner("secret")

type Server struct {
	http.Handler
	templates   *template.Template
	store       DataStore
	signer      *Signer
	websocket   *WebSocket
	templateMap map[string]*template.Template
}

type Task struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	NextID int    `json:"next_id"`
	PrevID int    `json:"prev_id"`
	Target string `json:"target"`
}

type DataStore interface {
	GetTask(id int) (Task, error)
}

type Data struct {
	Title string
	Body  string
	Task  Task
}

type SEO struct {
	Description string
}

type Config struct {
	Data          Data
	SEO           SEO
	PartialUpdate bool
	Path          string
	Version       string
	Dark          bool
}

func (s *Server) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	stringID := strings.TrimPrefix(r.URL.Path, "/task/")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task, err := s.store.GetTask(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.templates.ExecuteTemplate(w, "task", task)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // Check path here
		http.NotFound(w, r)
		return
	}
	data := Data{
		Title: "Index",
		Body:  "This is a test",
	}
	seo := SEO{
		Description: "This is the index page",
	}
	config := Config{
		Data: data,
		SEO:  seo,
	}
	err := s.routeHandler("index", config, w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) taskHandler(w http.ResponseWriter, r *http.Request) {
	stringID := strings.TrimPrefix(r.URL.Path, "/task/")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task, err := s.store.GetTask(id)
	if err != nil {
		log.Fatal(err)
	}
	data := Data{
		Title: task.Title,
		Task:  task,
	}
	seo := SEO{
		Description: fmt.Sprintf("Page about Task %s", task.Title),
	}
	config := Config{
		Data: data,
		SEO:  seo,
	}
	err = s.routeHandler("task", config, w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Title: "About",
		Body:  "This is about",
	}
	seo := SEO{
		Description: "This is the about page",
	}
	config := Config{
		Data: data,
		SEO:  seo,
	}
	err := s.routeHandler("about", config, w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func darkMode(w http.ResponseWriter, config *Config) {
	darkMode := w.Header().Get("Dark-Mode")
	if darkMode == "enabled" {
		config.Dark = true
	}
}

func versionHandler(config *Config) {
	config.Version = version
}

func (s *Server) routeHandler(name string, config Config, w http.ResponseWriter, r *http.Request) error {
	var err error
	config.Path = r.URL.Path
	darkMode(w, &config)
	versionHandler(&config)
	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Add("hx-push", r.URL.Path)
		config.PartialUpdate = true
		// the order here matters, first the <head> parts should be sent as this
		// allows the DOM to be properly built in the frontend.
		err = s.templateMap[name].ExecuteTemplate(w, "seo", config)
		err = s.templateMap[name].ExecuteTemplate(w, "content", config)
		err = s.templateMap[name].ExecuteTemplate(w, "header", config)
	} else {
		err = s.templateMap[name].ExecuteTemplate(w, "template.html", config)
	}
	return err
}

func (s *Server) dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Print(r.FormValue("email"))
	}
	w.Write([]byte("This is data"))
}

func cookieMiddleware(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		sessionID := uuid.New().String()
		signedSessionID := signer.Sign(sessionID)
		cookieValue := sessionID + "|" + signedSessionID
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: cookieValue,
		})
	}
	darkMode, err := r.Cookie("dark-mode")
	if err != nil {
		log.Print(err)
	} else {
		w.Header().Add("dark-mode", darkMode.Value)
	}
	log.Print(cookie)
}

func NewServer(store DataStore) *Server {
	s := new(Server)
	router := http.NewServeMux()
	rootMux := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(s.indexHandler))
	router.Handle("/about", http.HandlerFunc(s.aboutHandler))
	router.Handle("/data", http.HandlerFunc(s.dataHandler))
	router.Handle("/task/", http.HandlerFunc(s.taskHandler))
	// router.Handle("/task/", http.HandlerFunc(s.getTaskHandler))
	router.Handle("/ws/subscribe", http.HandlerFunc(s.subscribeHandler))

	handler := MiddleWare{router}

	// We exclude some paths from the middleware
	rootMux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	rootMux.Handle("/", handler)

	s.Handler = rootMux

	templates, err := template.New("dummy").Funcs(getFunctions()).ParseGlob(templateDir + "/*.html")
	template.Must(templates.ParseGlob(templateDir + "/components/*.html"))
	template.Must(templates.ParseGlob(templateDir + "/icons/*.html"))

	tmpl := make(map[string]*template.Template)
	about := template.Must(templates.Clone())
	index := template.Must(templates.Clone())
	task := template.Must(templates.Clone())
	tmpl["index"] = template.Must(index.ParseFiles(templateDir + "views/index.html"))
	tmpl["task"] = template.Must(task.ParseFiles(templateDir + "views/task.html"))
	tmpl["about"] = template.Must(about.ParseFiles(templateDir + "views/about.html"))
	s.templateMap = tmpl

	if err != nil {
		log.Fatal(err)
	}
	s.templates = templates

	s.store = store

	signer := NewSigner("secret")
	if err != nil {
		log.Fatal(err)
	}
	s.signer = signer

	envVersion := os.Getenv("VERSION")
	if envVersion != "" {
		version = envVersion
	}

	s.websocket = new(WebSocket)
	s.websocket.subscribers = make(map[*subscriber]struct{})
	s.websocket.subscriberMessageBuffer = 16
	s.websocket.publishLimiter = rate.NewLimiter(rate.Every(time.Millisecond*100), 8)

	return s
}
