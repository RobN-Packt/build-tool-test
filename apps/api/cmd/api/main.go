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

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/example/bookapi/internal/http/handlers"
	"github.com/example/bookapi/internal/http/middleware"
	"github.com/example/bookapi/internal/notifications"
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

	pool, err := pgxpool.New(ctx, dsn)
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

	httpHandler := buildHTTPHandler(pool)

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

func buildHTTPHandler(pool *pgxpool.Pool) http.Handler {
	bookRepo := repo.NewBookRepository(pool)
	bookService := service.NewBookService(bookRepo, buildBookServiceOptions(context.Background())...)
	bookHandler := handlers.NewBookHandler(bookService)

	router := mux.NewRouter()
	config := huma.DefaultConfig("Book API", "1.0.0")
	api := humamux.New(router, config)

	registerHealthRoutes(api)
	handlers.RegisterBookRoutes(api, bookHandler)

	return middleware.CORS(middleware.Logger(router))
}

func registerHealthRoutes(api huma.API) {
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
}

func buildBookServiceOptions(ctx context.Context) []service.BookServiceOption {
	var opts []service.BookServiceOption

	publisher, err := configureSNSBookPublisher(ctx)
	if err != nil {
		slog.Error("failed to configure SNS publisher for book events", "error", err)
	} else if publisher != nil {
		opts = append(opts, service.WithBookEventPublisher(publisher))
	}

	return opts
}

func configureSNSBookPublisher(ctx context.Context) (service.BookEventPublisher, error) {
	topicARN := strings.TrimSpace(os.Getenv("SNS_TOPIC_ARN"))
	if topicARN == "" {
		slog.Warn("SNS_TOPIC_ARN not set; book created events will not be published",
			"envVar", "SNS_TOPIC_ARN",
		)
		return nil, nil
	}

	region, err := snsRegionFromARN(topicARN)
	if err != nil {
		slog.Error("unable to derive AWS region from SNS topic ARN; book created events disabled",
			"topicArn", topicARN,
			"error", err,
		)
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	client := sns.NewFromConfig(cfg)
	return notifications.NewSNSBookEventPublisher(client, topicARN, slog.Default()), nil
}

func snsRegionFromARN(topicARN string) (string, error) {
	parts := strings.Split(topicARN, ":")
	if len(parts) < 6 || parts[0] != "arn" {
		return "", fmt.Errorf("invalid SNS topic ARN: %s", topicARN)
	}

	region := strings.TrimSpace(parts[3])
	if region == "" {
		return "", fmt.Errorf("missing region in SNS topic ARN: %s", topicARN)
	}

	return region, nil
}
