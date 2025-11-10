package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"github.com/example/bookshop/apps/api/internal/http/handlers"
	appmw "github.com/example/bookshop/apps/api/internal/http/middleware"
	"github.com/example/bookshop/apps/api/internal/repo"
	"github.com/example/bookshop/apps/api/internal/service"
)

var migrateOnly = flag.Bool("migrate-only", false, "run migrations and exit")

func main() {
	flag.Parse()

	dsn := getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/books?sslmode=disable")
	port := getEnv("PORT", "8080")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := repo.OpenPool(ctx, dsn)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	if err := repo.Migrate(ctx, pool); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	if *migrateOnly {
		log.Println("migrations complete")
		return
	}

	repository := repo.NewBookRepository(pool)
	bookService := service.NewBookService(repository)

	router := http.NewServeMux()
	cfg := huma.DefaultConfig("Book Shop API", "1.0.0")
	cfg.Servers = []*huma.Server{{URL: "/"}}
	api := humago.New(router, cfg)

	handlers.Register(api, &handlers.BookHandler{Service: bookService})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      appmw.Logger(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("API listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
