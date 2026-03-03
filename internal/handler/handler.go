package handler

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

const (
	baseURL  = "http://localhost:8080"
	idLength = 8
	charset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Handler struct {
	mu      sync.RWMutex
	storage map[string]string
}

func New() *Handler {
	return &Handler{
		storage: make(map[string]string),
	}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(strings.TrimSpace(string(body))) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	id := h.generateID()

	h.mu.Lock()
	h.storage[id] = originalURL
	h.mu.Unlock()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(baseURL + "/" + id))
}

func (h *Handler) GetURL(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	h.mu.RLock()
	originalURL, ok := h.storage[id]
	h.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func (h *Handler) BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (h *Handler) generateID() string {
	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
