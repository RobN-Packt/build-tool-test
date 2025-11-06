package books

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	listFn   func(ctx context.Context) ([]Book, error)
	getFn    func(ctx context.Context, id int64) (Book, error)
	createFn func(ctx context.Context, book Book) (Book, error)
	updateFn func(ctx context.Context, book Book) (Book, error)
	deleteFn func(ctx context.Context, id int64) error
}

func (m mockRepository) List(ctx context.Context) ([]Book, error) {
	return m.listFn(ctx)
}

func (m mockRepository) GetByID(ctx context.Context, id int64) (Book, error) {
	return m.getFn(ctx, id)
}

func (m mockRepository) Create(ctx context.Context, book Book) (Book, error) {
	return m.createFn(ctx, book)
}

func (m mockRepository) Update(ctx context.Context, book Book) (Book, error) {
	return m.updateFn(ctx, book)
}

func (m mockRepository) Delete(ctx context.Context, id int64) error {
	return m.deleteFn(ctx, id)
}

func TestService_Create_ValidInput(t *testing.T) {
	repo := mockRepository{
		createFn: func(_ context.Context, book Book) (Book, error) {
			book.ID = 1
			book.CreatedAt = time.Now()
			book.UpdatedAt = book.CreatedAt
			return book, nil
		},
	}

	service := NewService(repo)

	result, err := service.Create(context.Background(), CreateBookInput{
		Title:         "Test",
		Author:        "Author",
		ISBN:          "1234567890123",
		Price:         10.5,
		Stock:         5,
		Description:   "Desc",
		PublishedDate: "2020-01-02",
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Test", result.Title)
}

func TestService_Create_InvalidInput(t *testing.T) {
	service := NewService(mockRepository{})

	_, err := service.Create(context.Background(), CreateBookInput{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}

func TestService_Update_NotFound(t *testing.T) {
	repo := mockRepository{
		updateFn: func(_ context.Context, _ Book) (Book, error) {
			return Book{}, ErrBookNotFound
		},
	}

	service := NewService(repo)

	_, err := service.Update(context.Background(), 1, UpdateBookInput{
		Title:         "Test",
		Author:        "Author",
		ISBN:          "1234567890123",
		Price:         10,
		Stock:         2,
		Description:   "",
		PublishedDate: "2020-01-02",
	})

	assert.ErrorIs(t, err, ErrBookNotFound)
}

func TestService_Get_InvalidID(t *testing.T) {
	service := NewService(mockRepository{})

	_, err := service.Get(context.Background(), 0)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidInput))
}
