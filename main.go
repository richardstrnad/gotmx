package main

import (
	"net/http"
)

var templateDir = "html/"

func main() {
	store := NewInMemoryDataStore()
	server := NewServer(store)
	http.ListenAndServe(":80", server)
}
