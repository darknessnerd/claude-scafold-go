package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/demo/with-scaffold/internal/domain"
)

type PostgresLinkStore struct {
	db *sql.DB
}

func NewPostgresLinkStore(db *sql.DB) *PostgresLinkStore {
	return &PostgresLinkStore{db: db}
}

func (s *PostgresLinkStore) FindByURL(ctx context.Context, url string) (domain.Link, error) {
	var link domain.Link
	err := s.db.QueryRowContext(ctx,
		// Parameterized — no injection possible
		`SELECT code, original_url, created_at FROM links WHERE original_url = $1`,
		url,
	).Scan(&link.Code, &link.OriginalURL, &link.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Link{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Link{}, fmt.Errorf("postgres.FindByURL: %w", err)
	}
	return link, nil
}

func (s *PostgresLinkStore) FindByCode(ctx context.Context, code string) (domain.Link, error) {
	var link domain.Link
	err := s.db.QueryRowContext(ctx,
		`SELECT code, original_url, created_at FROM links WHERE code = $1`,
		code,
	).Scan(&link.Code, &link.OriginalURL, &link.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Link{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Link{}, fmt.Errorf("postgres.FindByCode: %w", err)
	}
	return link, nil
}

func (s *PostgresLinkStore) Save(ctx context.Context, link domain.Link) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO links (code, original_url) VALUES ($1, $2)`,
		link.Code, link.OriginalURL,
	)
	if err != nil {
		return fmt.Errorf("postgres.Save: %w", err)
	}
	return nil
}
