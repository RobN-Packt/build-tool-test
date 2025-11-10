package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/bookshop/apps/api/internal/domain"
	"github.com/example/bookshop/apps/api/internal/service"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	created domain.Book
	saveErr error
}

func (m *mockRepo) List(ctx context.Context) ([]domain.Book, error) {
	return nil, nil
}

func (m *mockRepo) Get(ctx context.Context, id string) (domain.Book, error) {
	return domain.Book{}, domain.ErrNotFound
}

func (m *mockRepo) Create(ctx context.Context, book domain.Book) (domain.Book, error) {
	if m.saveErr != nil {
		return domain.Book{}, m.saveErr
	}
	m.created = book
	return book, nil
}

func (m *mockRepo) Update(ctx context.Context, book domain.Book) (domain.Book, error) {
	return book, m.saveErr
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	return m.saveErr
}

func TestCreateBookSuccess(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewBookService(repo)

	book, err := svc.CreateBook(context.Background(), domain.CreateBookParams{
		Title:    "  The Go Programming Language  ",
		Author:   "Alan Donovan",
		Price:    42.50,
		Currency: "usd",
		Stock:    5,
	})

	require.NoError(t, err)
	require.NotEmpty(t, book.ID)
	require.Equal(t, "The Go Programming Language", book.Title)
	require.Equal(t, "Alan Donovan", book.Author)
	require.Equal(t, "USD", book.Currency)
	require.Equal(t, 5, book.Stock)
	require.WithinDuration(t, time.Now(), book.CreatedAt, time.Second)
	require.WithinDuration(t, time.Now(), book.UpdatedAt, time.Second)
}

func TestCreateBookValidationError(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewBookService(repo)

	_, err := svc.CreateBook(context.Background(), domain.CreateBookParams{
		Title:    "",
		Author:   "",
		Price:    -1,
		Currency: "BAD",
		Stock:    -1,
	})

	require.Error(t, err)
	var vErrs domain.ValidationErrors
	require.True(t, errors.As(err, &vErrs))
	require.GreaterOrEqual(t, len(vErrs), 4)

	fields := make(map[string]struct{})
	for _, v := range vErrs {
		fields[v.Field] = struct{}{}
	}

	for _, name := range []string{"title", "author", "price", "stock"} {
		_, ok := fields[name]
		require.Truef(t, ok, "expected validation error for %s", name)
	}
}

func TestUpdateBookValidation(t *testing.T) {
	repo := &mockRepo{}
	svc := service.NewBookService(repo)

	_, err := svc.UpdateBook(context.Background(), "missing", domain.UpdateBookParams{
		Price: func() *float32 { v := float32(-5); return &v }(),
	})

	require.Error(t, err)
}
