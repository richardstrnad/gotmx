package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	http.Handler
	templates *template.Template
	store     DataStore
}

type Task struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	NextID int    `json:"next_id"`
	PrevID int    `json:"prev_id"`
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
	task, err := s.store.GetTask(1)
	if err != nil {
		log.Fatal(err)
	}
	data := Data{
		Title: "Hello World",
		Body:  "This is a test",
		Task:  task,
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
	Title string
	Body  string
	Task  Task
}

func (s *Server) dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is data"))
}

func NewServer(store DataStore) *Server {
	s := new(Server)
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(s.indexHandler))
	router.Handle("/about", http.HandlerFunc(s.aboutHandler))
	router.Handle("/data", http.HandlerFunc(s.dataHandler))
	router.Handle("/task/", http.HandlerFunc(s.getTaskHandler))
	router.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	s.Handler = router

	templates, err := template.ParseGlob(templateDir + "*.html")
	if err != nil {
		log.Fatal(err)
	}
	s.templates = templates

	s.store = store

	return s
}
