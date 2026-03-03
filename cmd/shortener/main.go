package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stasenkoin/URL-ShortenerAI/internal/config"
	"github.com/stasenkoin/URL-ShortenerAI/internal/handler"
)

func main() {
	cfg := config.ParseFlags()

	h := handler.New(cfg.BaseURL)

	r := chi.NewRouter()
	r.Post("/", h.ShortenURL)
	r.Get("/{id}", h.GetURL)
	r.NotFound(h.BadRequest)
	r.MethodNotAllowed(h.BadRequest)

	log.Printf("Starting server on %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
