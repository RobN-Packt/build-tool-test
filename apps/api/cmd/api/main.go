package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/example/bookapi/internal/http/handlers"
	"github.com/example/bookapi/internal/http/middleware"
	"github.com/example/bookapi/internal/repo"
	"github.com/example/bookapi/internal/repo/migrations"
	"github.com/example/bookapi/internal/service"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("api", flag.ExitOnError)
	migrateOnly := fs.Bool("migrate", false, "apply database migrations and exit")
	if err := fs.Parse(args); err != nil {
		return err
	}

	_ = godotenv.Load()

	port := envOrDefault("PORT", "8080")
	dsn := os.Getenv("DB_DSN")
	if strings.TrimSpace(dsn) == "" {
		return errors.New("DB_DSN is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("parse DB_DSN: %w", err)
	}
	logDatabaseConfig(cfg)

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	if err := applyMigrations(ctx, pool); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	if *migrateOnly {
		slog.Info("migrations applied, exiting")
		return nil
	}

	httpHandler := buildHTTPHandler(pool, slog.Default())

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func envOrDefault(key, defaultValue string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return defaultValue
	}
	return v
}

func applyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := migrations.Files.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)

	for _, name := range files {
		sqlBytes, err := migrations.Files.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("exec migration %s: %w", name, err)
		}
	}
	return nil
}

type healthOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

type booksHealthOutput struct {
	Body struct {
		Status    string    `json:"status"`
		CheckedAt time.Time `json:"checkedAt"`
	}
}

func buildHTTPHandler(pool *pgxpool.Pool, logger *slog.Logger) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	bookRepo := repo.NewBookRepository(pool)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService, logger)

	router := mux.NewRouter()
	config := huma.DefaultConfig("Book API", "1.0.0")
	api := humamux.New(router, config)

	registerHealthRoutes(api, bookRepo, logger)
	handlers.RegisterBookRoutes(api, bookHandler)

	return middleware.CORS(middleware.Logger(router))
}

func registerHealthRoutes(api huma.API, bookRepo *repo.BookRepository, logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}
	huma.Register(api, huma.Operation{
		OperationID:   "healthz",
		Method:        http.MethodGet,
		Path:          "/healthz",
		Summary:       "Service health check",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, _ *struct{}) (*healthOutput, error) {
		out := &healthOutput{}
		out.Body.Status = "ok"
		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "books-health",
		Method:        http.MethodGet,
		Path:          "/health/books",
		Summary:       "Check database access via books table",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, _ *struct{}) (*booksHealthOutput, error) {
		checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := bookRepo.CheckBooks(checkCtx); err != nil {
			logger.Error("books health check failed", "error", err)
			return nil, huma.NewError(http.StatusInternalServerError, fmt.Sprintf("books health check failed: %v", err))
		}

		out := &booksHealthOutput{}
		out.Body.Status = "ok"
		out.Body.CheckedAt = time.Now().UTC()
		return out, nil
	})
}

func logDatabaseConfig(cfg *pgxpool.Config) {
	if cfg == nil || cfg.ConnConfig == nil {
		return
	}
	logger := slog.Default()
	logger.Info("database configuration",
		"host", cfg.ConnConfig.Host,
		"port", cfg.ConnConfig.Port,
		"database", cfg.ConnConfig.Database,
		"user", cfg.ConnConfig.User,
		"sslmode", cfg.ConnConfig.Config.RuntimeParams["sslmode"],
	)
}
