package main

import (
	"fmt"
	"gallery-server/db"
	"gallery-server/handlers"
	"log"
	"net/http"
)

func main() {
	db.InitDB("images.db")

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/all", handlers.AllImageHandler)
	http.HandleFunc("/image", handlers.ImageHandler)
	http.HandleFunc("/addimage", handlers.NewImageHandler)
	http.HandleFunc("/delimage", handlers.DeleteImageHandler)
	http.HandleFunc("/setimage", handlers.UpdateImageHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/tknlgn", handlers.CookieLoginHandler)
	http.HandleFunc("/signup", handlers.SignUpHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	fs2 := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", fs2))
	fmt.Println("starting servers")
	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
	}))
	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}
