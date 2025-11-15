package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/middleware"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Example 1: POST / — create a short URL (plain text)
func Example_postRoot() {
	mux := newMuxExample()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	fmt.Println(w.Code)
	fmt.Println(strings.TrimSpace(w.Body.String()))
	// Output:
	// 201
	// http://localhost:8088/0
}

// Example 2: GET /{id} — redirect by short URL
func Example_getRedirect() {
	mux := newMuxExample()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://foo.bar"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	r = httptest.NewRequest(http.MethodGet, "/0", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Location"))
	// Output:
	// 307
	// http://foo.bar
}

// Example 3: POST /api/shorten — shorten URL via JSON
func Example_postShortenJSON() {
	mux := newMuxExample()
	r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"http://foo.bar"}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	fmt.Println(w.Code)
	fmt.Println(strings.TrimSpace(w.Body.String()))
	// Output:
	// 201
	// {"result":"http://localhost:8088/0"}
}

// Example 4: POST /api/shorten/batch — shorten multiple URLs
func Example_postShortenBatch() {
	mux := newMuxExample()
	body := `[{"correlation_id":"a","original_url":"http://a.b"}]`
	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	fmt.Println(w.Code)
	fmt.Println(strings.TrimSpace(w.Body.String()))
	// Output:
	// 201
	// [{"correlation_id":"a","short_url":"http://localhost:8088/0"}]
}

// Example 5: GET /api/user/urls — list user URLs
func Example_getUserURLs() {
	mux := newMuxExample()

	// Create one short URL first
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	cookies := res.Cookies()
	defer res.Body.Close()

	// Fetch user's URLs
	r = httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	fmt.Println(w.Code)
	fmt.Println(strings.TrimSpace(w.Body.String()))
	// Output:
	// 200
	// [{"original_url":"http://example.com","short_url":"http://localhost:8088/0"}]
}

// Example 6: DELETE /api/user/urls — delete user URLs
func Example_deleteUserURLs() {
	mux := newMuxExample()

	// Create one short URL first
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	cookies := res.Cookies()
	defer res.Body.Close()

	// Delete it
	r = httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["0"]`))
	for _, c := range cookies {
		r.AddCookie(c)
	}
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	fmt.Println(w.Code)
	// Output:
	// 202
}

func newMuxExample() *chi.Mux {
	cfg := config.Config{
		ListenAddr:      ":8088",
		ShortenerPrefix: "http://localhost:8088",
	}
	logger, _ := zap.NewDevelopment()
	repo, _ := repository.NewMemoryRepo(cfg, logger)
	svc := service.NewService(repo, cfg, logger)
	h := NewHandler(svc, logger)
	auth := middleware.NewAuth(repo, cfg, logger)

	mux := chi.NewRouter()
	mux.Use(middleware.Compression, auth.Authentication)
	h.Register(mux)
	return mux
}
