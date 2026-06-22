package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/demo/with-scaffold/internal/domain"
)

// shortenerService defined here, by the consumer (handler), not in service/.
type shortenerService interface {
	Shorten(ctx context.Context, url string) (domain.Link, error)
	Resolve(ctx context.Context, code string) (domain.Link, error)
}

type ShortenerHandler struct {
	svc shortenerService
}

func NewShortenerHandler(svc shortenerService) *ShortenerHandler {
	return &ShortenerHandler{svc: svc}
}

func (h *ShortenerHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /shorten", h.handleShorten)
	mux.HandleFunc("GET /{code}", h.handleRedirect)
	mux.HandleFunc("GET /healthz", h.handleHealth)
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Short string `json:"short"`
}

func (h *ShortenerHandler) handleShorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	link, err := h.svc.Shorten(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenResponse{Short: link.Code})
}

func (h *ShortenerHandler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	link, err := h.svc.Resolve(r.Context(), code)
	if errors.Is(err, domain.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, link.OriginalURL, http.StatusMovedPermanently)
}

func (h *ShortenerHandler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
