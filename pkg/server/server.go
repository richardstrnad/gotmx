package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/richardstrnad/gotmx/pkg/gotmx"
	"github.com/richardstrnad/gotmx/pkg/service"
	"github.com/richardstrnad/gotmx/pkg/signer"
	"golang.org/x/time/rate"
)

var version = ""
var templateDir = "html/"
var sign = signer.NewSigner("secret")

type Server struct {
	http.Handler
	templates   *template.Template
	signer      *signer.Signer
	websocket   *WebSocket
	templateMap map[string]*template.Template
	Service     *service.Service
}

type Data struct {
	Title string
	Body  string
	Task  gotmx.Task
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
	if s.templateMap[name] == nil {
		return fmt.Errorf("template %s not found", name)
	}

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
		signedSessionID := sign.Sign(sessionID)
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

func parseViews(tmpl map[string]*template.Template, templates *template.Template) (map[string]*template.Template, error) {
	views, err := filepath.Glob(
		path.Join(templateDir, "./views/*.html"),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, view := range views {
		name := filepath.Base(view)
		if filepath.Ext(name) == ".html" {
			name = strings.TrimSuffix(name, ".html")
			c := template.Must(templates.Clone())
			tmpl[name] = template.Must(c.ParseFiles(view))
		}
	}
	return tmpl, nil
}

func NewServer(svc *service.Service) *Server {
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
	tmpl, err = parseViews(tmpl, templates)
	if err != nil {
		log.Fatal(err)
	}
	s.templateMap = tmpl
	s.templates = templates

	sign := signer.NewSigner("secret")
	if err != nil {
		log.Fatal(err)
	}
	s.signer = sign

	envVersion := os.Getenv("VERSION")
	if envVersion != "" {
		version = envVersion
	}

	s.websocket = new(WebSocket)
	s.websocket.subscribers = make(map[*subscriber]struct{})
	s.websocket.subscriberMessageBuffer = 16
	s.websocket.publishLimiter = rate.NewLimiter(rate.Every(time.Millisecond*100), 8)

	s.Service = svc

	return s
}
