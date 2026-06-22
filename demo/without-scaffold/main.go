// WITHOUT SCAFFOLD
// Claude asked: "build a URL shortener in Go with Postgres"
// No project rules, no architecture constraints, no security guidance.
//
// This code works. It will pass a quick smoke test.
// Every problem is annotated — these are the gaps the scaffold prevents.
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	_ "github.com/lib/pq"
)

// [PROBLEM 1] Global mutable state.
// No constructor injection — db wired as package-level var.
// Impossible to swap for tests without modifying this file.
var db *sql.DB

func main() {
	var err error

	// [PROBLEM 2] Hardcoded credentials.
	// Password in source = leaks in git history, logs, CI artifacts.
	// Scaffold rule 04-security: "Never hardcode credentials."
	db, err = sql.Open("postgres", "host=localhost user=postgres password=secret dbname=urls sslmode=disable")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	http.HandleFunc("/", redirectHandler)

	fmt.Println("listening on :8080")
	// [PROBLEM 3] ListenAndServe error silently swallowed.
	http.ListenAndServe(":8080", nil)
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Short string `json:"short"`
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 405)
		return
	}

	var req ShortenRequest
	// [PROBLEM 4] Decode error silently swallowed.
	// Invalid JSON = empty struct, no 400 returned to caller.
	json.NewDecoder(r.Body).Decode(&req)

	// [PROBLEM 5] SQL INJECTION.
	// User-supplied URL interpolated directly into query string.
	// Input: ' OR '1'='1 — dumps entire table.
	// Scaffold rule 04-security: "Parameterized queries always."
	row := db.QueryRow(fmt.Sprintf("SELECT code FROM links WHERE original_url = '%s'", req.URL))
	var existingCode string
	if err := row.Scan(&existingCode); err == nil {
		json.NewEncoder(w).Encode(ShortenResponse{Short: existingCode})
		return
	}

	// [PROBLEM 6] math/rand is not cryptographically safe.
	// Codes are predictable — attacker can enumerate short URLs.
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	shortCode := string(code)

	// [PROBLEM 7] SQL INJECTION on INSERT.
	// Same class as problem 5 — both params injectable.
	_, err = db.Exec(fmt.Sprintf("INSERT INTO links (code, original_url) VALUES ('%s', '%s')", shortCode, req.URL))
	if err != nil {
		// [PROBLEM 8] Error not wrapped with context.
		// "db error" in logs — impossible to trace which query failed.
		http.Error(w, "db error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShortenResponse{Short: shortCode})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]
	if code == "" {
		http.NotFound(w, r)
		return
	}

	// [PROBLEM 9] SQL INJECTION on lookup. Third occurrence.
	row := db.QueryRow(fmt.Sprintf("SELECT original_url FROM links WHERE code = '%s'", code))
	var url string
	if err := row.Scan(&url); err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

// [PROBLEM 10] Zero tests.
// All logic lives in handlers — impossible to unit test without HTTP + live DB.
// Scaffold rule 03-testing: "No business logic in handler; service layer 100% tested."
