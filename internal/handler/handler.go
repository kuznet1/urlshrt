package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/model"
	"github.com/kuznet1/urlshrt/internal/service"
	"io"
	"net/http"
)

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) Handler {
	return Handler{svc: svc}
}

func (h Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read body: %s", err), http.StatusBadRequest)
		return
	}

	body := string(bytes)
	url, err := h.svc.Shorten(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to shorten url: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url))
}

func (h Handler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	var req model.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode body: %s", err), http.StatusBadRequest)
	}

	url, err := h.svc.Shorten(req.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to shorten url: %s", err), http.StatusInternalServerError)
		return
	}

	resp := model.ShortenResponse{
		Result: url,
	}

	respJSON(w, resp, http.StatusCreated)
}

func respJSON(w http.ResponseWriter, resp any, code int) {
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
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
