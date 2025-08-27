package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kuznet1/urlshrt/internal/errs"
	"github.com/kuznet1/urlshrt/internal/model"
	"github.com/kuznet1/urlshrt/internal/service"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Handler struct {
	svc    service.Service
	logger *zap.Logger
}

func NewHandler(svc service.Service, logger *zap.Logger) Handler {
	return Handler{svc: svc, logger: logger}
}

func (h Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read body: %s", err), http.StatusBadRequest)
		return
	}

	body := string(bytes)
	url, err := h.svc.Shorten(body)
	var duplicatedError *errs.DuplicatedURLError
	isDuplicatedError := errors.As(err, &duplicatedError)
	if err != nil && !isDuplicatedError {
		internalError("failed to shorten url", err, h.logger, w)
		return
	}

	if isDuplicatedError {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	w.Write([]byte(url))
}

func (h Handler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	var req model.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode body: %s", err), http.StatusBadRequest)
	}

	url, err := h.svc.Shorten(req.URL)
	var duplicatedError *errs.DuplicatedURLError
	isDuplicatedError := errors.As(err, &duplicatedError)
	if err != nil && !isDuplicatedError {
		internalError("failed to shorten url", err, h.logger, w)
		return
	}

	resp := model.ShortenResponse{
		Result: url,
	}

	if isDuplicatedError {
		respJSON(w, resp, http.StatusConflict, h.logger)
	} else {
		respJSON(w, resp, http.StatusCreated, h.logger)
	}
}

func (h Handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	var req []model.BatchShortenRequestItem
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode body: %s", err), http.StatusBadRequest)
	}

	var urls []string
	for _, reqItem := range req {
		urls = append(urls, reqItem.OriginalURL)
	}

	shortenLinks, err := h.svc.BatchShorten(urls)
	var duplicatedError *errs.DuplicatedURLError
	if errors.As(err, &duplicatedError) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		internalError("failed to shorten urls", err, h.logger, w)
		return
	}

	resp := []model.BatchShortenResponseItem{}
	for i := range urls {
		resp = append(resp, model.BatchShortenResponseItem{
			CorrelationID: req[i].CorrelationID,
			ShortURL:      shortenLinks[i],
		})
	}

	respJSON(w, resp, http.StatusCreated, h.logger)
}

func respJSON(w http.ResponseWriter, resp any, code int, logger *zap.Logger) {
	data, err := json.Marshal(resp)
	if err != nil {
		internalError("failed to encode response", err, logger, w)
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
		internalError("failed to lengthen url", err, h.logger, w)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func internalError(msg string, err error, logger *zap.Logger, w http.ResponseWriter) {
	logger.Error(msg, zap.Error(err))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
