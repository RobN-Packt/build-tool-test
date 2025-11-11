package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/example/bookapi/internal/domain"
	"github.com/example/bookapi/internal/repo"
)

func TestBookServiceCreate_Success(t *testing.T) {
	mockRepo := newMockBookRepo()
	svc := NewBookService(mockRepo)
	now := time.Date(2025, 11, 10, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	input := BookCreateInput{
		Title:    "  The Go Programming Language ",
		Author:   "Alan Donovan",
		Price:    49.99,
		Currency: "",
		Stock:    5,
	}

	book, err := svc.CreateBook(context.Background(), input)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, book.ID)
	require.Equal(t, "The Go Programming Language", book.Title)
	require.Equal(t, "Alan Donovan", book.Author)
	require.Equal(t, "USD", book.Currency)
	require.Equal(t, now, book.CreatedAt)
	require.Equal(t, now, book.UpdatedAt)

	stored, ok := mockRepo.store[book.ID]
	require.True(t, ok)
	require.Equal(t, stored, book)
}

func TestBookServiceCreate_ValidationError(t *testing.T) {
	mockRepo := newMockBookRepo()
	svc := NewBookService(mockRepo)

	_, err := svc.CreateBook(context.Background(), BookCreateInput{
		Price:    -1,
		Currency: "US",
		Stock:    -5,
	})
	require.Error(t, err)

	validationErr, ok := err.(ValidationError)
	require.True(t, ok)
	require.Contains(t, validationErr.Fields, "title")
	require.Contains(t, validationErr.Fields, "author")
	require.Contains(t, validationErr.Fields, "price")
	require.Contains(t, validationErr.Fields, "currency")
	require.Contains(t, validationErr.Fields, "stock")
}

func TestBookServiceUpdate_PartialSuccess(t *testing.T) {
	mockRepo := newMockBookRepo()
	svc := NewBookService(mockRepo)
	now := time.Date(2025, 11, 10, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return now }

	existing := domain.Book{
		ID:        uuid.New(),
		Title:     "Original",
		Author:    "Author",
		Price:     10,
		Currency:  "USD",
		Stock:     5,
		CreatedAt: now.Add(-time.Hour),
		UpdatedAt: now.Add(-time.Hour),
	}
	mockRepo.store[existing.ID] = existing

	newTitle := "Updated Title"
	newPrice := 12.5
	input := BookUpdateInput{
		Title: &newTitle,
		Price: &newPrice,
	}

	updated, err := svc.UpdateBook(context.Background(), existing.ID, input)
	require.NoError(t, err)
	require.Equal(t, newTitle, updated.Title)
	require.Equal(t, newPrice, updated.Price)
	require.Equal(t, existing.Author, updated.Author)
	require.True(t, updated.UpdatedAt.After(existing.UpdatedAt))
}

func TestBookServiceUpdate_MissingBody(t *testing.T) {
	mockRepo := newMockBookRepo()
	svc := NewBookService(mockRepo)
	bookID := uuid.New()
	mockRepo.store[bookID] = domain.Book{ID: bookID}

	_, err := svc.UpdateBook(context.Background(), bookID, BookUpdateInput{})
	require.Error(t, err)
	_, ok := err.(ValidationError)
	require.True(t, ok)
}

type mockBookRepo struct {
	store map[uuid.UUID]domain.Book
}

func newMockBookRepo() *mockBookRepo {
	return &mockBookRepo{store: make(map[uuid.UUID]domain.Book)}
}

func (m *mockBookRepo) Create(_ context.Context, book domain.Book) error {
	m.store[book.ID] = book
	return nil
}

func (m *mockBookRepo) Get(_ context.Context, id uuid.UUID) (domain.Book, error) {
	book, ok := m.store[id]
	if !ok {
		return domain.Book{}, repo.ErrNotFound
	}
	return book, nil
}

func (m *mockBookRepo) List(_ context.Context) ([]domain.Book, error) {
	var result []domain.Book
	for _, book := range m.store {
		result = append(result, book)
	}
	return result, nil
}

func (m *mockBookRepo) Update(_ context.Context, book domain.Book) error {
	if _, ok := m.store[book.ID]; !ok {
		return repo.ErrNotFound
	}
	m.store[book.ID] = book
	return nil
}

func (m *mockBookRepo) Delete(_ context.Context, id uuid.UUID) error {
	if _, ok := m.store[id]; !ok {
		return repo.ErrNotFound
	}
	delete(m.store, id)
	return nil
}
