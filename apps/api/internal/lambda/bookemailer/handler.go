package bookemailer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

const (
	expectedType = "BOOK_CREATED"
	emailBodyTpl = "A new book has been added:\nTitle: %s\nPrice: £%.2f\n\nBook ID: %s\n\nThis is an automated message."
)

// Handler processes BOOK_CREATED SNS messages and emails end users.
type Handler struct {
	sender        EmailSender
	fallbackEmail string
	logger        *slog.Logger
}

// New wires dependencies from environment variables and AWS configuration.
func New(ctx context.Context) (*Handler, error) {
	sesRegion := os.Getenv("SES_REGION")
	if sesRegion == "" {
		return nil, errors.New("SES_REGION env var must be set")
	}

	fromEmail := os.Getenv("FROM_EMAIL")
	if fromEmail == "" {
		return nil, errors.New("FROM_EMAIL env var must be set")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(sesRegion))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Handler{
		sender: &sesSender{
			client: ses.NewFromConfig(cfg),
			source: fromEmail,
		},
		fallbackEmail: os.Getenv("TO_EMAIL"),
		logger:        logger,
	}, nil
}

// Handle satisfies aws-lambda-go expectations for an SNS handler.
func (h *Handler) Handle(ctx context.Context, event events.SNSEvent) error {
	for _, record := range event.Records {
		if err := h.processRecord(ctx, record); err != nil {
			h.logger.Error("failed to process SNS record",
				"error", err,
				"messageId", record.SNS.MessageID,
			)
		}
	}

	return nil
}

func (h *Handler) processRecord(ctx context.Context, record events.SNSEventRecord) error {
	var msg BookCreatedMessage
	if err := json.Unmarshal([]byte(record.SNS.Message), &msg); err != nil {
		return fmt.Errorf("parse sns message: %w", err)
	}

	if err := validateMessage(msg); err != nil {
		return err
	}

	recipient := msg.UserEmail
	if recipient == "" {
		recipient = h.fallbackEmail
	}

	if recipient == "" {
		return errors.New("no recipient email provided")
	}

	bookID := msg.BookID.String()
	price := *msg.Price

	subject := fmt.Sprintf("New book added: %s", msg.Title)
	body := fmt.Sprintf(emailBodyTpl, msg.Title, price, bookID)

	if err := h.sender.Send(ctx, recipient, subject, body); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	h.logger.Info("book created email sent",
		"bookId", bookID,
		"title", msg.Title,
		"recipient", recipient,
		"messageId", record.SNS.MessageID,
	)

	return nil
}

// BookCreatedMessage mirrors the SNS message schema.
type BookCreatedMessage struct {
	Type      string         `json:"type"`
	BookID    flexibleString `json:"bookId"`
	Title     string         `json:"title"`
	Price     *float64       `json:"price"`
	UserEmail string         `json:"userEmail"`
}

func validateMessage(msg BookCreatedMessage) error {
	switch {
	case msg.Type != expectedType:
		return fmt.Errorf("unexpected message type: %s", msg.Type)
	case strings.TrimSpace(msg.BookID.String()) == "":
		return errors.New("bookId must be provided")
	case msg.Title == "":
		return errors.New("title must be provided")
	case msg.Price == nil:
		return errors.New("price must be provided")
	default:
		return nil
	}
}

// EmailSender abstracts SES interactions for easier testing.
type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

type sesSender struct {
	client *ses.Client
	source string
}

func (s *sesSender) Send(ctx context.Context, to, subject, body string) error {
	if to == "" {
		return errors.New("missing destination email")
	}

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Subject: &types.Content{Data: aws.String(subject)},
			Body: &types.Body{
				Text: &types.Content{Data: aws.String(body)},
			},
		},
		Source: aws.String(s.source),
	}

	_, err := s.client.SendEmail(ctx, input)
	return err
}

type flexibleString string

func (fs flexibleString) String() string {
	return string(fs)
}

func (fs *flexibleString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*fs = flexibleString(s)
		return nil
	}

	var num json.Number
	if err := json.Unmarshal(data, &num); err == nil {
		*fs = flexibleString(num.String())
		return nil
	}

	return errors.New("bookId must be a string or number")
}
