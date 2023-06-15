package main

import (
	"net/http"
)

type MiddleWare struct {
	handler http.Handler
}

func (c MiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookieMiddleware(w, r)
	c.handler.ServeHTTP(w, r)
}
