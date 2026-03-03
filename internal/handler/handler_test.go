package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestShortenURL(t *testing.T) {
	t.Run("valid URL returns 201 and shortened URL", func(t *testing.T) {
		h := New("http://localhost:8080")
		body := strings.NewReader("https://practicum.yandex.ru/")
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()

		h.ShortenURL(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", res.StatusCode)
		}

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		shortURL := string(responseBody)
		if !strings.HasPrefix(shortURL, "http://localhost:8080/") {
			t.Errorf("expected shortened URL to start with http://localhost:8080/, got %s", shortURL)
		}
	})

	t.Run("empty body returns 400", func(t *testing.T) {
		h := New("http://localhost:8080")
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		w := httptest.NewRecorder()

		h.ShortenURL(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestGetURL(t *testing.T) {
	h := New("http://localhost:8080")
	body := strings.NewReader("https://practicum.yandex.ru/")
	req := httptest.NewRequest(http.MethodPost, "/", body)
	w := httptest.NewRecorder()
	h.ShortenURL(w, req)

	shortURL := w.Body.String()
	parts := strings.Split(shortURL, "/")
	id := parts[len(parts)-1]

	t.Run("existing ID returns 307 with Location", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		h.GetURL(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusTemporaryRedirect {
			t.Errorf("expected status 307, got %d", res.StatusCode)
		}

		location := res.Header.Get("Location")
		if location != "https://practicum.yandex.ru/" {
			t.Errorf("expected Location https://practicum.yandex.ru/, got %s", location)
		}
	})

	t.Run("non-existing ID returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		h.GetURL(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}
