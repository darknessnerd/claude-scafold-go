package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/demo/with-scaffold/internal/handler"
	"github.com/demo/with-scaffold/internal/repository"
	"github.com/demo/with-scaffold/internal/service"
)

func main() {
	// Secrets from environment — never hardcoded.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Wire dependencies: cmd/main.go is the only file that imports all layers.
	store := repository.NewPostgresLinkStore(db)
	svc := service.NewShortener(store)
	h := handler.NewShortenerHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}
