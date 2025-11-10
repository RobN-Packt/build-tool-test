package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/example/bookapi/internal/domain"
)

type ValidationError struct {
	Fields map[string]string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", v.Fields)
}

type BookRepository interface {
	Create(ctx context.Context, book domain.Book) error
	Get(ctx context.Context, id uuid.UUID) (domain.Book, error)
	List(ctx context.Context) ([]domain.Book, error)
	Update(ctx context.Context, book domain.Book) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BookService struct {
	repo BookRepository
	now  func() time.Time
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{
		repo: repo,
		now:  time.Now,
	}
}

type BookInput struct {
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Stock    int     `json:"stock"`
}

func (s *BookService) CreateBook(ctx context.Context, input BookInput) (domain.Book, error) {
	if err := validateBookInput(input); err != nil {
		return domain.Book{}, err
	}

	now := s.now().UTC()
	book := domain.Book{
		ID:        uuid.New(),
		Title:     strings.TrimSpace(input.Title),
		Author:    strings.TrimSpace(input.Author),
		Price:     input.Price,
		Currency:  normalizeCurrency(input.Currency),
		Stock:     input.Stock,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, book); err != nil {
		return domain.Book{}, err
	}
	return book, nil
}

func (s *BookService) GetBook(ctx context.Context, id uuid.UUID) (domain.Book, error) {
	return s.repo.Get(ctx, id)
}

func (s *BookService) ListBooks(ctx context.Context) ([]domain.Book, error) {
	return s.repo.List(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, id uuid.UUID, input BookInput) (domain.Book, error) {
	if err := validateBookInput(input); err != nil {
		return domain.Book{}, err
	}

	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}

	existing.Title = strings.TrimSpace(input.Title)
	existing.Author = strings.TrimSpace(input.Author)
	existing.Price = input.Price
	existing.Currency = normalizeCurrency(input.Currency)
	existing.Stock = input.Stock
	existing.UpdatedAt = s.now().UTC()

	if err := s.repo.Update(ctx, existing); err != nil {
		return domain.Book{}, err
	}
	return existing, nil
}

func (s *BookService) DeleteBook(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func normalizeCurrency(curr string) string {
	curr = strings.TrimSpace(curr)
	if curr == "" {
		return "USD"
	}
	return strings.ToUpper(curr)
}

func validateBookInput(input BookInput) error {
	errors := make(map[string]string)

	title := strings.TrimSpace(input.Title)
	if title == "" {
		errors["title"] = "required"
	} else if !withinLength(title, 1, 200) {
		errors["title"] = "must be 1-200 characters"
	}

	author := strings.TrimSpace(input.Author)
	if author == "" {
		errors["author"] = "required"
	} else if !withinLength(author, 1, 200) {
		errors["author"] = "must be 1-200 characters"
	}

	if input.Price < 0 {
		errors["price"] = "must be >= 0"
	}

	currency := normalizeCurrency(input.Currency)
	if len(currency) != 3 || strings.ToUpper(currency) != currency {
		errors["currency"] = "must be ISO 4217 code"
	}

	if input.Stock < 0 {
		errors["stock"] = "must be >= 0"
	}

	if len(errors) > 0 {
		return ValidationError{Fields: errors}
	}
	return nil
}

func withinLength(value string, min, max int) bool {
	length := utf8.RuneCountInString(value)
	return length >= min && length <= max
}
