package main

// main is the composition root. Wire all layers here — nowhere else.
// Replace the stubs below with real dependencies as the service grows.
//
// Import order: stdlib → internal → external (enforced by goimports)
func main() {
	// TODO: initialise logger (e.g. slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	// TODO: open DB connection (e.g. sql.Open("pgx", os.Getenv("DATABASE_URL")))
	// TODO: wire repository  (e.g. repository.NewUserStore(db))
	// TODO: wire service     (e.g. service.NewUserService(userStore, logger))
	// TODO: wire handler     (e.g. handler.NewRouter(userSvc))
	// TODO: start HTTP server
}
