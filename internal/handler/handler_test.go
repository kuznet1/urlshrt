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
	"time"
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
		code := put(t, mux)
		require.Equal(t, http.StatusCreated, code)
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

func TestDuplicated(t *testing.T) {
	mux, err := newMux(t)
	if err != nil {
		t.Fatal(err)
	}

	code := put(t, mux)
	require.Equal(t, http.StatusCreated, code)

	code = put(t, mux)
	require.Equal(t, http.StatusConflict, code)
}

func put(t *testing.T, mux *chi.Mux) int {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8088/0", string(resBody))
	return res.StatusCode
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

func TestBatchDelete(t *testing.T) {
	mux, err := newMux(t)
	if err != nil {
		t.Fatal(err)
	}

	var cookies []*http.Cookie
	t.Run("shorten urls", func(t *testing.T) {
		cookies = putWithCookie(t, mux, "http://example.com") // user 1
		putWithCookie(t, mux, "http://foo.bar")               // user 2
	})

	t.Run("delete", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(`["0","1"]`))
		for _, c := range cookies {
			r.AddCookie(c)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		require.NoError(t, err)
		require.Equal(t, http.StatusAccepted, w.Code)
	})

	t.Run("check redirect 1", func(t *testing.T) {
		deadline := time.Now().Add(time.Second * 5)
		for {
			r := httptest.NewRequest(http.MethodGet, "/0", nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode == http.StatusGone {
				break
			}

			if time.Now().After(deadline) {
				t.Fatalf("deadline exceeded")
			}
		}
	})

	t.Run("check redirect 2", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/1", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, "http://foo.bar", res.Header.Get("Location"))
		require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
	})
}

func TestUrlsByUser(t *testing.T) {
	mux, err := newMux(t)
	if err != nil {
		t.Fatal(err)
	}

	var cookies []*http.Cookie
	t.Run("shorten url", func(t *testing.T) {
		cookies = putWithCookie(t, mux, "http://example.com")
	})

	t.Run("user urls", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		for _, c := range cookies {
			r.AddCookie(c)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `[{"original_url":"http://example.com","short_url":"http://localhost:8088/0"}]`, string(resBody))
	})
}

func putWithCookie(t *testing.T, mux *chi.Mux, url string) []*http.Cookie {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(url))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return res.Cookies()
}

func newMux(t testing.TB) (*chi.Mux, error) {
	cfg := config.Config{
		ListenAddr:      ":8088",
		ShortenerPrefix: "http://localhost:8088",
		FileStoragePath: repoFile,
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	repo, err := repository.NewMemoryRepo(cfg, logger)
	if err != nil {
		return nil, err
	}

	svc := service.NewService(repo, cfg, logger)
	h := NewHandler(svc, logger)
	auth := middleware.NewAuth(repo, cfg, logger)
	mux := chi.NewRouter()
	mux.Use(middleware.Compression, auth.Authentication)
	h.Register(mux)
	t.Cleanup(func() {
		os.Remove(repoFile)
	})

	return mux, nil
}
