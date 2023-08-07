package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"nhooyr.io/websocket"
)

func (s *Server) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r,
		&websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		log.Print(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = s.subscribe(r.Context(), c)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}

func (s *Server) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	stringID := strings.TrimPrefix(r.URL.Path, "/task/")
	id, err := strconv.Atoi(stringID)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	task, err := s.Service.Store.GetTask(id)
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
	task, err := s.Service.Store.GetTask(id)
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
