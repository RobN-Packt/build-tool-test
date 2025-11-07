package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	humachi "github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/cursor/bookshop/backend/internal/books"
	"github.com/cursor/bookshop/backend/internal/database"
	"github.com/cursor/bookshop/backend/internal/database/migrations"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbCfg := database.FromEnv()
	db, err := database.Connect(ctx, dbCfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := migrations.Apply(ctx, db); err != nil {
		log.Fatalf("failed to apply database migrations: %v", err)
	}

	repo := books.NewRepository(db)
	service := books.NewService(repo)

	port := getenvDefault("HTTP_PORT", "8080")
	baseURL := getenvDefault("PUBLIC_BASE_URL", fmt.Sprintf("http://localhost:%s", port))

	router := chi.NewRouter()
	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	api := humachi.New(router, huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   "Bookshop API",
				Version: "1.0.0",
				Description: "A RESTful API for managing the bookshop catalog, " +
					"including CRUD operations and seeded data.",
			},
			Servers: []*huma.Server{
				{URL: baseURL},
			},
		},
		OpenAPIPath: "/openapi",
		DocsPath:    "/docs",
	})

	books.RegisterRoutes(api, service)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Bookshop API listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown server gracefully: %v", err)
	}
}

func getenvDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
