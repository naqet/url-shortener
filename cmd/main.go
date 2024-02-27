package main

import (
	"log"
	"net/http"

	"github.com/naqet/url-shortener/db"
	"github.com/naqet/url-shortener/handlers"
	"github.com/naqet/url-shortener/services"
)

func main() {
	database, err := db.NewDB("./db/sqlite.db")

	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()
	mux := http.NewServeMux()

	urlService, err := services.NewUrlService(database)

	if err != nil {
		log.Fatal(err)
	}

	urlHandler := handlers.NewUrlHandler(&urlService)
	mux.HandleFunc("POST /short", urlHandler.Post)
	mux.HandleFunc("GET /short/{key}", urlHandler.Get)
	mux.HandleFunc("GET /{key}", urlHandler.Redirect)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe("localhost:3000", mux)
}
