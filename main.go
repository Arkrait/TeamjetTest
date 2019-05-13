package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
)

func main() {
	go InitServer()

	log.Print("initialized server")
	select {}
}

func InitServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	SetupRoutes(r)

	http.ListenAndServe(":8080", r)
}
