package main

import (
	handlers "indexer/Handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Use(cors.Default().Handler)

	router.Post("/api/indexer/", handlers.Indexar)

	router.Get("/api/indexer/{term}-{max}", handlers.GetCorreos)
	
	http.ListenAndServe(":5000", router)

}
