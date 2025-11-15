package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Benchmark 1: POST / — create a short URL (plain text)
func Benchmark_postRoot(b *testing.B) {
	mux, err := newMux(b)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
	}
}

// Benchmark 2: GET /{id} — redirect by short URL
func Benchmark_getRedirect(b *testing.B) {
	mux, err := newMux(b)
	if err != nil {
		b.Fatal(err)
	}
	// Create one URL first
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://foo.bar"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = httptest.NewRequest(http.MethodGet, "/0", nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, r)
	}
}

// Benchmark 3: POST /api/shorten — shorten URL via JSON
func Benchmark_postShortenJSON(b *testing.B) {
	mux, err := newMux(b)
	if err != nil {
		b.Fatal(err)
	}
	body := `{"url":"http://foo.bar"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
	}
}

// Benchmark 4: POST /api/shorten/batch — shorten multiple URLs
func Benchmark_postShortenBatch(b *testing.B) {
	mux, err := newMux(b)
	if err != nil {
		b.Fatal(err)
	}
	body := `[{"correlation_id":"a","original_url":"http://a.b"}]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
	}
}

// Benchmark 5: GET /api/user/urls — list user URLs
func Benchmark_getUserURLs(b *testing.B) {
	mux, err := newMux(b)
	if err != nil {
		b.Fatal(err)
	}
	// Create one URL to have data
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	cookies := res.Cookies()
	defer res.Body.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		for _, c := range cookies {
			r.AddCookie(c)
		}
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, r)
	}
}

// Benchmark 6: DELETE /api/user/urls — delete user URLs
func Benchmark_deleteUserURLs(b *testing.B) {
	mux, err := newMux(b)
	if err != nil {
		b.Fatal(err)
	}
	// Create one URL to delete
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	cookies := res.Cookies()
	defer res.Body.Close()

	body := `["0"]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(body))
		for _, c := range cookies {
			r.AddCookie(c)
		}
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, r)
	}
}
