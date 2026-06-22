package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/demo/with-scaffold/internal/domain"
	"github.com/demo/with-scaffold/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubShortener struct {
	shortenFn func(ctx context.Context, url string) (domain.Link, error)
	resolveFn func(ctx context.Context, code string) (domain.Link, error)
}

func (s *stubShortener) Shorten(ctx context.Context, url string) (domain.Link, error) {
	return s.shortenFn(ctx, url)
}

func (s *stubShortener) Resolve(ctx context.Context, code string) (domain.Link, error) {
	return s.resolveFn(ctx, code)
}

func TestHandleShorten(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		svc        *stubShortener
		wantStatus int
		wantCode   string
	}{
		{
			name: "happy path",
			body: `{"url":"https://example.com"}`,
			svc: &stubShortener{
				shortenFn: func(_ context.Context, _ string) (domain.Link, error) {
					return domain.Link{Code: "abc123", OriginalURL: "https://example.com"}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantCode:   "abc123",
		},
		{
			name:       "missing url returns 400",
			body:       `{}`,
			svc:        &stubShortener{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json returns 400",
			body:       `not-json`,
			svc:        &stubShortener{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error returns 500",
			body: `{"url":"https://example.com"}`,
			svc: &stubShortener{
				shortenFn: func(_ context.Context, _ string) (domain.Link, error) {
					return domain.Link{}, assert.AnError
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			handler.NewShortenerHandler(tt.svc).RegisterRoutes(mux)

			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantCode != "" {
				var resp map[string]string
				require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
				assert.Equal(t, tt.wantCode, resp["short"])
			}
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		code         string
		svc          *stubShortener
		wantStatus   int
		wantLocation string
	}{
		{
			name: "known code redirects",
			code: "abc123",
			svc: &stubShortener{
				resolveFn: func(_ context.Context, _ string) (domain.Link, error) {
					return domain.Link{Code: "abc123", OriginalURL: "https://example.com"}, nil
				},
			},
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://example.com",
		},
		{
			name: "unknown code returns 404",
			code: "xxxxxx",
			svc: &stubShortener{
				resolveFn: func(_ context.Context, _ string) (domain.Link, error) {
					return domain.Link{}, domain.ErrNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			handler.NewShortenerHandler(tt.svc).RegisterRoutes(mux)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.code, nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantLocation != "" {
				assert.Equal(t, tt.wantLocation, rr.Header().Get("Location"))
			}
		})
	}
}
