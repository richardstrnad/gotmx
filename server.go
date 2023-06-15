package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var VERSION = ""
var signer = NewSigner("secret")

type Server struct {
	http.Handler
	templates *template.Template
	store     DataStore
	signer    *Signer
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
		Version: VERSION,
	}
	err = s.templates.ExecuteTemplate(w, "template.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Title: "About",
		Body:  "This is about",
	}
	err := s.templates.ExecuteTemplate(w, "template.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

type Data struct {
	Title   string
	Body    string
	Task    Task
	Version string
}

func (s *Server) dataHandler(w http.ResponseWriter, r *http.Request) {
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

	handler := MiddleWare{router}

	// We exclude some paths from the middleware
	rootMux.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	rootMux.Handle("/", handler)

	s.Handler = rootMux

	templates, err := template.New("dummy").Funcs(getFunctions()).ParseGlob(templateDir + "/*.html")
	template.Must(templates.ParseGlob(templateDir + "/components/*.html"))
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

	VERSION = os.Getenv("VERSION")

	return s
}
