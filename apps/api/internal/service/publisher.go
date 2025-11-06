package service

import (
    "context"
    "log/slog"

    "github.com/google/uuid"

    "github.com/example/book-poc/apps/api/domain"
)

type LoggerPublisher struct {
    logger *slog.Logger
}

func NewLoggerPublisher(logger *slog.Logger) *LoggerPublisher {
    return &LoggerPublisher{logger: logger}
}

func (p *LoggerPublisher) PublishPurchase(ctx context.Context, input domain.PurchaseInput) (string, error) {
    messageID := uuid.NewString()
    p.logger.InfoContext(ctx, "purchase enqueued", "message_id", messageID, "book_id", input.BookID, "quantity", input.Quantity, "customer_id", input.CustomerID)
    return messageID, nil
}
