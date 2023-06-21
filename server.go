package main

import (
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
	Title   string
	Body    string
	Task    Task
	Version string
}

type Config struct {
	Data          Data
	PartialUpdate bool
	Path          string
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
	task, err := s.store.GetTask(1)
	if err != nil {
		log.Fatal(err)
	}
	data := Data{
		Title:   "Hello World",
		Body:    "This is a test",
		Task:    task,
		Version: version,
	}
	config := Config{
		Data: data,
	}
	err = s.routeHandler("index", config, w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Title: "About",
		Body:  "This is about",
	}
	config := Config{
		Data: data,
	}
	err := s.routeHandler("about", config, w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) routeHandler(name string, config Config, w http.ResponseWriter, r *http.Request) error {
	var err error
	config.Path = r.URL.Path
	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Add("hx-push", r.URL.Path)
		config.PartialUpdate = true
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
	log.Print(cookie)
}

func NewServer(store DataStore) *Server {
	s := new(Server)
	router := http.NewServeMux()
	rootMux := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(s.indexHandler))
	router.Handle("/about", http.HandlerFunc(s.aboutHandler))
	router.Handle("/data", http.HandlerFunc(s.dataHandler))
	router.Handle("/task/", http.HandlerFunc(s.getTaskHandler))
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
	tmpl["index"] = template.Must(index.ParseFiles(templateDir + "views/index.html"))
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
