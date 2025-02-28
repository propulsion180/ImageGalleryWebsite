package main

import (
	"gallery-server/handlers"
	"io/fs"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", handlers.RootHandler)

	fs := http.FileServer(http.Dir("/public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
