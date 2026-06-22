package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("link not found")

type Link struct {
	Code        string
	OriginalURL string
	CreatedAt   time.Time
}
