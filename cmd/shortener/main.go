package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stasenkoin/URL-ShortenerAI/internal/handler"
)

func main() {
	h := handler.New()

	r := chi.NewRouter()
	r.Post("/", h.ShortenURL)
	r.Get("/{id}", h.GetURL)
	r.NotFound(h.BadRequest)
	r.MethodNotAllowed(h.BadRequest)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
