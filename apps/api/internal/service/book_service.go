package service

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/google/uuid"

    "github.com/example/book-poc/apps/api/domain"
    "github.com/example/book-poc/apps/api/internal/repository"
)

type PurchasePublisher interface {
    PublishPurchase(ctx context.Context, input domain.PurchaseInput) (string, error)
}

type BookService struct {
    repo      repository.BookRepository
    publisher PurchasePublisher
    clock     func() time.Time
}

func NewBookService(repo repository.BookRepository, publisher PurchasePublisher) *BookService {
    return &BookService{
        repo:      repo,
        publisher: publisher,
        clock:     time.Now,
    }
}

func (s *BookService) List(ctx context.Context, limit int, cursor string) ([]domain.Book, *string, error) {
    return s.repo.List(limit, cursor)
}

func (s *BookService) Get(ctx context.Context, id string) (*domain.Book, error) {
    return s.repo.Get(id)
}

func (s *BookService) Create(ctx context.Context, input domain.CreateBookInput) (*domain.Book, error) {
    if err := validateCreate(input); err != nil {
        return nil, err
    }
    now := s.clock()
    book := domain.Book{
        ID:        uuid.NewString(),
        Title:     strings.TrimSpace(input.Title),
        Author:    strings.TrimSpace(input.Author),
        Price:     input.Price,
        Currency:  strings.ToUpper(input.Currency),
        Stock:     input.Stock,
        CreatedAt: now,
        UpdatedAt: now,
    }
    return s.repo.Create(book)
}

func (s *BookService) Update(ctx context.Context, id string, input domain.UpdateBookInput) (*domain.Book, error) {
    if err := validateUpdate(input); err != nil {
        return nil, err
    }
    existing, err := s.repo.Get(id)
    if err != nil {
        return nil, err
    }
    existing.Title = strings.TrimSpace(input.Title)
    existing.Author = strings.TrimSpace(input.Author)
    existing.Price = input.Price
    existing.Currency = strings.ToUpper(input.Currency)
    existing.Stock = input.Stock
    existing.UpdatedAt = s.clock()
    return s.repo.Update(*existing)
}

func (s *BookService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(id)
}

func (s *BookService) Purchase(ctx context.Context, input domain.PurchaseInput) (string, error) {
    if input.Quantity <= 0 {
        return "", fmt.Errorf("quantity must be positive")
    }

    book, err := s.repo.Get(input.BookID)
    if err != nil {
        return "", err
    }
    if book.Stock < input.Quantity {
        return "", repository.ErrInsufficientStock
    }

    book.Stock -= input.Quantity
    book.UpdatedAt = s.clock()
    if _, err := s.repo.Update(*book); err != nil {
        return "", err
    }

    if s.publisher == nil {
        return "", errors.New("purchase publisher is not configured")
    }

    return s.publisher.PublishPurchase(ctx, input)
}

func validateCreate(input domain.CreateBookInput) error {
    if strings.TrimSpace(input.Title) == "" {
        return fmt.Errorf("title is required")
    }
    if strings.TrimSpace(input.Author) == "" {
        return fmt.Errorf("author is required")
    }
    if input.Price < 0 {
        return fmt.Errorf("price must be non-negative")
    }
    if len(strings.TrimSpace(input.Currency)) != 3 {
        return fmt.Errorf("currency must be 3 letters")
    }
    if input.Stock < 0 {
        return fmt.Errorf("stock must be non-negative")
    }
    return nil
}

func validateUpdate(input domain.UpdateBookInput) error {
    return validateCreate(domain.CreateBookInput(input))
}
