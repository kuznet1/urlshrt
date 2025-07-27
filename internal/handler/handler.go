package handler

import (
	"errors"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/model"
	"github.com/kuznet1/urlshrt/internal/repository"
	"github.com/kuznet1/urlshrt/internal/service"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h Handler) Root(w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(r.URL.Path, "/")
	if id != "" {
		h.lengthen(w, r, id)
	} else {
		h.shorten(w, r)
	}
}

func (h Handler) shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("method %s not allowed, allowed is %s", r.Method, http.MethodPost), http.StatusBadRequest)
		return
	}

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

	resp := fmt.Sprintf("http://localhost:8080/%s", urlID)

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(resp))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write response: %s", err), http.StatusInternalServerError)
	}
}

func (h Handler) lengthen(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("method %s not allowed, allowed is %s", r.Method, http.MethodGet), http.StatusBadRequest)
		return
	}

	url, err := h.svc.Lengthen(id)
	if errors.Is(err, model.ErrIDParsing) {
		http.Error(w, fmt.Sprintf("unable to parse %q: it must be alphanumeric", id), http.StatusBadRequest)
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, fmt.Sprintf("url for shortening %q doesn't exist", id), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to lengthen url: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
