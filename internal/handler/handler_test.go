package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
			want{"/1", http.MethodGet, http.StatusTemporaryRedirect, "", "http://example.com"},
		}, {
			"non-existing url",
			want{"/2", http.MethodGet, http.StatusNotFound, "url for shortening \"2\" doesn't exist\n", ""},
		}, {
			"bad url",
			want{"/-1", http.MethodGet, http.StatusBadRequest, "unable to parse \"-1\": it must be alphanumeric\n", ""},
		}, {
			"bad method",
			want{"/1", http.MethodPost, http.StatusMethodNotAllowed, "", ""},
		},
	}

	cfg := config.Config{
		ListenAddr:      ":8088",
		ShortenerPrefix: "http://localhost:8088",
	}
	svc := service.NewService(&repository.MemoryRepo{})
	h := NewHandler(svc, cfg)
	mux := chi.NewRouter()
	mux.Post("/", h.Shorten)
	mux.Get("/{id}", h.Lengthen)

	t.Run("shorten url", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		res := w.Result()
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "http://localhost:8088/1", string(resBody))
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
