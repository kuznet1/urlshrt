package handler

import (
	"bytes"
	"compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/middleware"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const repoFile = "test-repo.json"

type want struct {
	url            string
	verb           string
	code           int
	response       string
	locationHeader string
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name string
		want want
	}{
		{
			"existing url",
			want{"/0", http.MethodGet, http.StatusTemporaryRedirect, "", "http://example.com"},
		}, {
			"non-existing url",
			want{"/1", http.MethodGet, http.StatusNotFound, "url for shortening \"1\" doesn't exist\n", ""},
		}, {
			"bad url",
			want{"/-1", http.MethodGet, http.StatusBadRequest, "unable to parse \"-1\": it must be alphanumeric\n", ""},
		}, {
			"bad method",
			want{"/0", http.MethodPost, http.StatusMethodNotAllowed, "", ""},
		},
	}

	mux, err := newMux(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("shorten url", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "http://localhost:8088/0", string(resBody))
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := test.want
			r := httptest.NewRequest(want.verb, want.url, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.locationHeader, res.Header.Get("Location"))
		})
	}
}

func TestShortenJSON(t *testing.T) {
	mux, err := newMux(t)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"http://foo.bar"}`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"result":"http://localhost:8088/0"}`, string(resBody))
}

func TestBacthShorten(t *testing.T) {
	mux, err := newMux(t)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(`[{"correlation_id":"foo","original_url":"http://foo.bar"}]`))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `[{"correlation_id":"foo","short_url":"http://localhost:8088/0"}]`, string(resBody))
}

func TestShortenJSONGzip(t *testing.T) {
	mux, err := newMux(t)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("shorten with gzip", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"http://foo.bar"}`))
		r.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, "gzip", res.Header.Get("Content-Encoding"))

		gz, err := gzip.NewReader(res.Body)
		require.NoError(t, err)
		defer gz.Close()

		resBody, err := io.ReadAll(gz)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, `{"result":"http://localhost:8088/0"}`, string(resBody))
	})

	t.Run("empty body without gzip", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/0", nil)
		r.Header.Set("Accept-Encoding", "gzip")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "", res.Header.Get("Content-Encoding"))
		assert.Equal(t, "", string(resBody))
		require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
	})
}

func TestPostBodyGzip(t *testing.T) {
	mux, err := newMux(t)
	if err != nil {
		t.Fatal(err)
	}

	url := "http://foo.bar"
	t.Run("shorten with gzip body", func(t *testing.T) {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		_, err := gz.Write([]byte(url))
		require.NoError(t, err)
		gz.Close()
		r := httptest.NewRequest(http.MethodPost, "/", &buf)
		r.Header.Set("Content-Encoding", "gzip")

		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "http://localhost:8088/0", string(resBody))
	})

	t.Run("check redirect", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/0", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, url, res.Header.Get("Location"))
		require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
	})
}

func newMux(t *testing.T) (*chi.Mux, error) {
	cfg := config.Config{
		ListenAddr:      ":8088",
		ShortenerPrefix: "http://localhost:8088",
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	repo, err := repository.NewMemoryRepo(repoFile, logger)
	if err != nil {
		return nil, err
	}

	svc := service.NewService(repo, cfg)
	h := NewHandler(svc, logger)
	mux := chi.NewRouter()
	mux.Use(middleware.Compression)
	mux.Post("/", h.Shorten)
	mux.Get("/{id}", h.Lengthen)
	mux.Post("/api/shorten", h.ShortenJSON)
	mux.Post("/api/shorten/batch", h.ShortenBatch)

	t.Cleanup(func() {
		os.Remove(repoFile)
	})

	return mux, nil
}
