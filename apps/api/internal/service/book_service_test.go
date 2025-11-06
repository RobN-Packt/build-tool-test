package service

import (
    "context"
    "testing"
    "time"

    "github.com/example/book-poc/apps/api/domain"
    "github.com/example/book-poc/apps/api/internal/repository"
)

type stubPublisher struct {
    lastInput domain.PurchaseInput
}

func (s *stubPublisher) PublishPurchase(ctx context.Context, input domain.PurchaseInput) (string, error) {
    s.lastInput = input
    return "message-1", nil
}

func TestBookService_CreateAndPurchase(t *testing.T) {
    repo := repository.NewInMemoryBookRepository()
    publisher := &stubPublisher{}
    svc := NewBookService(repo, publisher)
    svc.clock = func() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) }

    book, err := svc.Create(context.Background(), domain.CreateBookInput{Title: "Test", Author: "Author", Price: 10, Currency: "usd", Stock: 2})
    if err != nil {
        t.Fatalf("create failed: %v", err)
    }

    _, err = svc.Purchase(context.Background(), domain.PurchaseInput{BookID: book.ID, Quantity: 1, CustomerID: "c1"})
    if err != nil {
        t.Fatalf("purchase failed: %v", err)
    }

    stored, _ := repo.Get(book.ID)
    if stored.Stock != 1 {
        t.Fatalf("expected stock 1, got %d", stored.Stock)
    }

    if publisher.lastInput.Quantity != 1 {
        t.Fatalf("expected publisher to receive quantity=1")
    }
}

func TestBookService_Validation(t *testing.T) {
    repo := repository.NewInMemoryBookRepository()
    publisher := &stubPublisher{}
    svc := NewBookService(repo, publisher)

    if _, err := svc.Create(context.Background(), domain.CreateBookInput{Title: "", Author: "A", Price: 10, Currency: "USD", Stock: 1}); err == nil {
        t.Fatalf("expected validation error")
    }
}
