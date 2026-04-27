package main

import (
	"crypto/tls"
	"fmt"
	"gallery-server/db"
	"gallery-server/handlers"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	slog.Info("Started server!")

	db.InitDB("images.db")
	slog.Info("database initialized")

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

	acm := &autocert.Manager{
		Cache:      autocert.DirCache("/var/www/.cache"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("gallery.wmmp.xyz"),
	}

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: &tls.Config{GetCertificate: acm.GetCertificate},
		Handler:   handlers.LoggingMiddleware(http.DefaultServeMux),
	}

	go http.ListenAndServe(":80", acm.HTTPHandler(nil))

	if err := server.ListenAndServeTLS("", ""); err != nil {
		slog.Error("Server Failed", "err", err)
		os.Exit(1)
	}

}
