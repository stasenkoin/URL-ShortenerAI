package main

import (
	"log"
	"net/http"

	"github.com/stasenkoin/URL-ShortenerAI/internal/handler"
)

func main() {
	h := handler.New()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", h.ShortenURL)
	mux.HandleFunc("GET /{id}", h.GetURL)
	mux.HandleFunc("/", h.BadRequest)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
