package books

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidInput         = errors.New("invalid book payload")
	ErrInvalidPublishedDate = errors.New("invalid published date format, expected YYYY-MM-DD")
)

// Service exposes domain operations for books.
type Service interface {
	List(ctx context.Context) ([]Book, error)
	Get(ctx context.Context, id int64) (Book, error)
	Create(ctx context.Context, input CreateBookInput) (Book, error)
	Update(ctx context.Context, id int64, input UpdateBookInput) (Book, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

// NewService constructs a Service instance.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context) ([]Book, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, id int64) (Book, error) {
	if id <= 0 {
		return Book{}, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *service) Create(ctx context.Context, input CreateBookInput) (Book, error) {
	if err := validateInput(input.Title, input.Author, input.ISBN, input.Price, input.Stock, input.PublishedDate); err != nil {
		return Book{}, err
	}

	publishedDate, err := parseDate(input.PublishedDate)
	if err != nil {
		return Book{}, err
	}

	book := Book{
		Title:         strings.TrimSpace(input.Title),
		Author:        strings.TrimSpace(input.Author),
		ISBN:          strings.TrimSpace(input.ISBN),
		Price:         input.Price,
		Stock:         input.Stock,
		Description:   strings.TrimSpace(input.Description),
		PublishedDate: publishedDate,
	}

	return s.repo.Create(ctx, book)
}

func (s *service) Update(ctx context.Context, id int64, input UpdateBookInput) (Book, error) {
	if id <= 0 {
		return Book{}, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	if err := validateInput(input.Title, input.Author, input.ISBN, input.Price, input.Stock, input.PublishedDate); err != nil {
		return Book{}, err
	}

	publishedDate, err := parseDate(input.PublishedDate)
	if err != nil {
		return Book{}, err
	}

	book := Book{
		ID:            id,
		Title:         strings.TrimSpace(input.Title),
		Author:        strings.TrimSpace(input.Author),
		ISBN:          strings.TrimSpace(input.ISBN),
		Price:         input.Price,
		Stock:         input.Stock,
		Description:   strings.TrimSpace(input.Description),
		PublishedDate: publishedDate,
	}

	return s.repo.Update(ctx, book)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func validateInput(title, author, isbn string, price float64, stock int, published string) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(author) == "" || strings.TrimSpace(isbn) == "" {
		return fmt.Errorf("%w: title, author and isbn are required", ErrInvalidInput)
	}

	if price < 0 {
		return fmt.Errorf("%w: price cannot be negative", ErrInvalidInput)
	}

	if stock < 0 {
		return fmt.Errorf("%w: stock cannot be negative", ErrInvalidInput)
	}

	if strings.TrimSpace(published) == "" {
		return fmt.Errorf("%w: published date is required", ErrInvalidInput)
	}

	return nil
}

func parseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, ErrInvalidPublishedDate
	}

	return parsed, nil
}
