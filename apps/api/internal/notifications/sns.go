package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/example/bookapi/internal/domain"
)

const bookCreatedEventType = "BOOK_CREATED"

// SNSBookEventPublisher publishes book domain events to Amazon SNS.
type SNSBookEventPublisher struct {
	client   *sns.Client
	topicARN string
	logger   *slog.Logger
}

// NewSNSBookEventPublisher constructs a publisher backed by SNS.
func NewSNSBookEventPublisher(client *sns.Client, topicARN string, logger *slog.Logger) *SNSBookEventPublisher {
	if logger == nil {
		logger = slog.Default()
	}
	return &SNSBookEventPublisher{
		client:   client,
		topicARN: topicARN,
		logger:   logger,
	}
}

// PublishBookCreated sends a BOOK_CREATED event for the provided book.
func (p *SNSBookEventPublisher) PublishBookCreated(ctx context.Context, book domain.Book) error {
	if p == nil || p.client == nil || p.topicARN == "" {
		return fmt.Errorf("sns book event publisher is not fully configured")
	}

	payload := map[string]any{
		"type":   bookCreatedEventType,
		"bookId": book.ID.String(),
		"title":  book.Title,
		"price":  book.Price,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal book created payload: %w", err)
	}

	p.logger.Info("attempting to publish book created message",
		"bookId", book.ID,
		"topicArn", p.topicARN,
	)

	resp, err := p.client.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(p.topicARN),
		Message:  aws.String(string(jsonBytes)),
	})
	if err != nil {
		p.logger.Error("failed to publish book created message",
			"error", err,
			"bookId", book.ID,
			"topicArn", p.topicARN,
		)
		return fmt.Errorf("publish SNS message: %w", err)
	}

	p.logger.Info("published book created message",
		"bookId", book.ID,
		"topicArn", p.topicARN,
		"messageId", aws.ToString(resp.MessageId),
	)

	return nil
}
