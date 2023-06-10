package main

import (
	"net/http"
)

var templateDir = "html/"

func main() {
	server := NewServer()
	http.ListenAndServe(":80", server)
}
