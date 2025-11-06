package service

import (
    "context"
    "encoding/json"
    "os"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/sqs"

    "github.com/example/book-poc/apps/api/domain"
)

type SQSPublisher struct {
    client   *sqs.Client
    queueURL string
}

func NewSQSPublisher(ctx context.Context) (*SQSPublisher, error) {
    queueURL := os.Getenv("SQS_QUEUE_URL")
    if queueURL == "" {
        return nil, nil
    }
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, err
    }
    client := sqs.NewFromConfig(cfg)
    return &SQSPublisher{client: client, queueURL: queueURL}, nil
}

func (p *SQSPublisher) PublishPurchase(ctx context.Context, input domain.PurchaseInput) (string, error) {
    if p == nil {
        return "", nil
    }
    payload := map[string]any{
        "bookId":     input.BookID,
        "quantity":   input.Quantity,
        "customerId": input.CustomerID,
    }
    bodyBytes, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }
    body := string(bodyBytes)
    out, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
        QueueUrl:    &p.queueURL,
        MessageBody: &body,
    })
    if err != nil {
        return "", err
    }
    if out.MessageId == nil {
        return "", fmt.Errorf("message id missing")
    }
    return *out.MessageId, nil
}
