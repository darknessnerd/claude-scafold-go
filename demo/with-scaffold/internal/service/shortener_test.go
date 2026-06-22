package service

import (
	"context"
	"testing"

	"github.com/demo/with-scaffold/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubStore is an in-memory linkStore — no DB, no mocking framework needed
// for a struct this small.
type stubStore struct {
	byURL  map[string]domain.Link
	byCode map[string]domain.Link
}

func newStubStore() *stubStore {
	return &stubStore{
		byURL:  make(map[string]domain.Link),
		byCode: make(map[string]domain.Link),
	}
}

func (s *stubStore) FindByURL(_ context.Context, url string) (domain.Link, error) {
	l, ok := s.byURL[url]
	if !ok {
		return domain.Link{}, domain.ErrNotFound
	}
	return l, nil
}

func (s *stubStore) FindByCode(_ context.Context, code string) (domain.Link, error) {
	l, ok := s.byCode[code]
	if !ok {
		return domain.Link{}, domain.ErrNotFound
	}
	return l, nil
}

func (s *stubStore) Save(_ context.Context, link domain.Link) error {
	s.byURL[link.OriginalURL] = link
	s.byCode[link.Code] = link
	return nil
}

func TestShortener_Shorten(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(*stubStore)
		url         string
		wantNewCode bool
		wantErr     bool
	}{
		{
			name:        "new URL generates a code",
			url:         "https://example.com",
			wantNewCode: true,
		},
		{
			name: "existing URL returns same code",
			setup: func(s *stubStore) {
				s.byURL["https://example.com"] = domain.Link{Code: "abc123", OriginalURL: "https://example.com"}
				s.byCode["abc123"] = s.byURL["https://example.com"]
			},
			url:         "https://example.com",
			wantNewCode: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := newStubStore()
			if tt.setup != nil {
				tt.setup(store)
			}
			svc := NewShortener(store)

			link, err := svc.Shorten(context.Background(), tt.url)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.url, link.OriginalURL)
			assert.NotEmpty(t, link.Code)

			if !tt.wantNewCode {
				assert.Equal(t, "abc123", link.Code)
			}
		})
	}
}

func TestShortener_Resolve(t *testing.T) {
	t.Parallel()

	t.Run("known code returns link", func(t *testing.T) {
		t.Parallel()
		store := newStubStore()
		store.byCode["abc123"] = domain.Link{Code: "abc123", OriginalURL: "https://example.com"}

		svc := NewShortener(store)
		link, err := svc.Resolve(context.Background(), "abc123")

		require.NoError(t, err)
		assert.Equal(t, "https://example.com", link.OriginalURL)
	})

	t.Run("unknown code returns wrapped ErrNotFound", func(t *testing.T) {
		t.Parallel()
		svc := NewShortener(newStubStore())

		_, err := svc.Resolve(context.Background(), "xxxxxx")

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	codes := make(map[string]struct{})
	for i := 0; i < 100; i++ {
		code, err := generateCode()
		require.NoError(t, err)
		assert.Len(t, code, 6, "code must be exactly 6 chars")
		codes[code] = struct{}{}
	}
	// 100 codes, crypto/rand — collisions astronomically unlikely
	assert.Greater(t, len(codes), 95, "expected high uniqueness across 100 codes")
}
