package handler

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/service"
	"io"
	"net/http"
)

type Handler struct {
	svc service.Service
	cfg config.Config
}

func NewHandler(svc service.Service, cfg config.Config) Handler {
	return Handler{svc: svc, cfg: cfg}
}

func (h Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read body: %s", err), http.StatusBadRequest)
		return
	}

	body := string(bytes)
	urlID, err := h.svc.Shorten(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to shorten url: %s", err), http.StatusInternalServerError)
		return
	}

	resp := h.cfg.ShortenerPrefix + "/" + urlID

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(resp))
}

func (h Handler) Lengthen(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, err := h.svc.Lengthen(id)
	var httpErr *errs.HTTPError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Error(), httpErr.Code())
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to lengthen url: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
