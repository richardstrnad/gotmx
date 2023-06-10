package main

import (
	"html/template"
	"log"
	"net/http"
)

type Server struct {
	http.Handler
	templates *template.Template
}

type Task struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type DataStore interface {
	GetTask(id int) (Task, error)
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Title: "Hello World",
		Body:  "This is a test",
	}
	err := s.templates.ExecuteTemplate(w, "template.html", data)
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
}

func (s *Server) dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is data"))
}

func NewServer() *Server {
	s := new(Server)
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(s.indexHandler))
	router.Handle("/about", http.HandlerFunc(s.aboutHandler))
	router.Handle("/data", http.HandlerFunc(s.dataHandler))
	router.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	s.Handler = router

	templates, err := template.ParseGlob(templateDir + "*.html")
	if err != nil {
		log.Fatal(err)
	}
	s.templates = templates

	return s
}
