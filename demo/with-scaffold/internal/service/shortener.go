package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/demo/with-scaffold/internal/domain"
)

// linkStore is defined here, by the consumer (service), not in repository.
// Swap Postgres for any other store without touching this package.
type linkStore interface {
	FindByURL(ctx context.Context, url string) (domain.Link, error)
	FindByCode(ctx context.Context, code string) (domain.Link, error)
	Save(ctx context.Context, link domain.Link) error
}

type Shortener struct {
	store linkStore
}

func NewShortener(store linkStore) *Shortener {
	return &Shortener{store: store}
}

func (s *Shortener) Shorten(ctx context.Context, url string) (domain.Link, error) {
	existing, err := s.store.FindByURL(ctx, url)
	if err == nil {
		return existing, nil
	}
	if err != domain.ErrNotFound {
		return domain.Link{}, fmt.Errorf("shortener.Shorten: lookup: %w", err)
	}

	code, err := generateCode()
	if err != nil {
		return domain.Link{}, fmt.Errorf("shortener.Shorten: generate code: %w", err)
	}

	link := domain.Link{Code: code, OriginalURL: url}
	if err := s.store.Save(ctx, link); err != nil {
		return domain.Link{}, fmt.Errorf("shortener.Shorten: save: %w", err)
	}
	return link, nil
}

func (s *Shortener) Resolve(ctx context.Context, code string) (domain.Link, error) {
	link, err := s.store.FindByCode(ctx, code)
	if err != nil {
		return domain.Link{}, fmt.Errorf("shortener.Resolve: %w", err)
	}
	return link, nil
}

func generateCode() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// base32 → uppercase; lowercase + trim padding → 6 or 7 chars, slice to 6
	code := strings.ToLower(base32.StdEncoding.EncodeToString(b))
	code = strings.TrimRight(code, "=")
	if len(code) > 6 {
		code = code[:6]
	}
	return code, nil
}
