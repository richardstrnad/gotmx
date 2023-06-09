package main

import (
	"html/template"
	"log"
	"net/http"
)

var templateDir = "html/"
var templateFiles = []string{
	templateDir + "template.html",
	templateDir + "header.html",
	templateDir + "footer.html",
}

type Server struct {
	http.Handler
	templates *template.Template
}

type Data struct {
	Title string
	Body  string
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

	templates := template.Must(template.ParseFiles(templateFiles...))
	s.templates = templates

	return s
}

func main() {
	server := NewServer()
	http.ListenAndServe(":80", server)
}
