package main

import (
    "context"
    "log/slog"
    "os"

    "gofr.dev/pkg/gofr"

    "github.com/example/book-poc/apps/api/handlers"
    "github.com/example/book-poc/apps/api/internal/repository"
    "github.com/example/book-poc/apps/api/internal/service"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

    app := gofr.New()
    app.Logger = logger

    repo := repository.NewInMemoryBookRepository()
    publisher, err := service.NewSQSPublisher(context.Background())
    if err != nil {
        logger.Error("failed to configure SQS publisher", "error", err)
    }
    var purchasePublisher service.PurchasePublisher = publisher
    if purchasePublisher == nil {
        purchasePublisher = service.NewLoggerPublisher(logger)
    }
    bookService := service.NewBookService(repo, purchasePublisher)

    handler := handlers.NewBookHandler(bookService)
    handler.Register(app)

    app.Run()
}
