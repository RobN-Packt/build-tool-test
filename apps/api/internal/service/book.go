package service

import (
	"context"
	"strings"
	"time"

	"github.com/example/bookshop/apps/api/internal/domain"
	"github.com/google/uuid"
)

// BookRepository describes storage requirements for the service.
type BookRepository interface {
	List(ctx context.Context) ([]domain.Book, error)
	Get(ctx context.Context, id string) (domain.Book, error)
	Create(ctx context.Context, book domain.Book) (domain.Book, error)
	Update(ctx context.Context, book domain.Book) (domain.Book, error)
	Delete(ctx context.Context, id string) error
}

// BookService contains business logic for books.
type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) ListBooks(ctx context.Context) ([]domain.Book, error) {
	return s.repo.List(ctx)
}

func (s *BookService) GetBook(ctx context.Context, id string) (domain.Book, error) {
	return s.repo.Get(ctx, id)
}

func (s *BookService) CreateBook(ctx context.Context, params domain.CreateBookParams) (domain.Book, error) {
	if params.Currency == "" {
		params.Currency = "USD"
	}
	params.Currency = strings.ToUpper(strings.TrimSpace(params.Currency))
	params.Title = strings.TrimSpace(params.Title)
	params.Author = strings.TrimSpace(params.Author)

	if err := domain.ValidateCreate(params); err != nil {
		return domain.Book{}, err
	}

	now := time.Now().UTC()
	book := domain.Book{
		ID:        uuid.NewString(),
		Title:     params.Title,
		Author:    params.Author,
		Price:     params.Price,
		Currency:  params.Currency,
		Stock:     params.Stock,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(ctx, book)
}

func (s *BookService) UpdateBook(ctx context.Context, id string, update domain.UpdateBookParams) (domain.Book, error) {
	if update.Currency != nil {
		trimmed := strings.ToUpper(strings.TrimSpace(*update.Currency))
		update.Currency = &trimmed
	}
	if update.Title != nil {
		trimmed := strings.TrimSpace(*update.Title)
		update.Title = &trimmed
	}
	if update.Author != nil {
		trimmed := strings.TrimSpace(*update.Author)
		update.Author = &trimmed
	}

	if err := domain.ValidateUpdate(update); err != nil {
		return domain.Book{}, err
	}

	book, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}

	domain.ApplyUpdate(&book, update)
	book.UpdatedAt = time.Now().UTC()

	return s.repo.Update(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
