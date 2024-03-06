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

    authService, err := services.NewAuthService(database)

	if err != nil {
		log.Fatal(err)
	}
    authHandler := handlers.NewAuthHandler(authService)

    mux.HandleFunc("POST /signup", authHandler.SignUp);
    mux.HandleFunc("POST /login", authHandler.LogIn);
    mux.HandleFunc("/logout", authHandler.LogOut);
    mux.HandleFunc("POST /refresh", authHandler.Refresh);

	urlService, err := services.NewUrlService(database)

	if err != nil {
        log.Fatal(err)
	}

	urlHandler := handlers.NewUrlHandler(urlService)
	mux.HandleFunc("POST /url", authHandler.Middleware(urlHandler.Post, true))
	mux.HandleFunc("GET /url/{key}", authHandler.Middleware(urlHandler.Get, true))
	mux.HandleFunc("GET /u/{key}", urlHandler.Redirect)

	pagesHandler := handlers.NewPagesHandler(authService, urlService);
    mux.HandleFunc("/", authHandler.Middleware(pagesHandler.Home, false));
    mux.HandleFunc("/signup", authHandler.Middleware(pagesHandler.SignUp, false));
    mux.HandleFunc("/login", authHandler.Middleware(pagesHandler.Login, false));
    mux.HandleFunc("/dashboard", authHandler.Middleware(pagesHandler.Dashboard, false));

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe("localhost:3000", mux)
}
