package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/example/bookapi/internal/domain"
	"github.com/example/bookapi/internal/repo"
	"github.com/example/bookapi/internal/service"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

type BookPayload struct {
	Title    string  `json:"title" example:"The Go Programming Language"`
	Author   string  `json:"author" example:"Alan A. A. Donovan"`
	Price    float64 `json:"price" example:"49.99"`
	Currency string  `json:"currency" example:"USD"`
	Stock    int     `json:"stock" example:"10"`
}

type BookResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BookIDInput struct {
	ID uuid.UUID `path:"id"`
}

type CreateBookInput struct {
	Body BookPayload `body:""`
}

type CreateBookOutput struct {
	Body BookResponse
}

type GetBookOutput struct {
	Body BookResponse
}

type ListBooksOutput struct {
	Body struct {
		Books []BookResponse `json:"books"`
	}
}

type UpdateBookInput struct {
	ID   uuid.UUID   `path:"id"`
	Body BookPayload `body:""`
}

type UpdateBookOutput struct {
	Body BookResponse
}

func RegisterBookRoutes(api huma.API, handler *BookHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "list-books",
		Method:        http.MethodGet,
		Path:          "/books",
		Summary:       "List books",
		DefaultStatus: http.StatusOK,
	}, handler.listBooks)

	huma.Register(api, huma.Operation{
		OperationID:   "get-book",
		Method:        http.MethodGet,
		Path:          "/books/{id}",
		Summary:       "Get book by ID",
		DefaultStatus: http.StatusOK,
	}, handler.getBook)

	huma.Register(api, huma.Operation{
		OperationID:   "create-book",
		Method:        http.MethodPost,
		Path:          "/books",
		Summary:       "Create book",
		DefaultStatus: http.StatusCreated,
	}, handler.createBook)

	huma.Register(api, huma.Operation{
		OperationID:   "update-book",
		Method:        http.MethodPut,
		Path:          "/books/{id}",
		Summary:       "Update book",
		DefaultStatus: http.StatusOK,
	}, handler.updateBook)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-book",
		Method:        http.MethodDelete,
		Path:          "/books/{id}",
		Summary:       "Delete book",
		DefaultStatus: http.StatusNoContent,
	}, handler.deleteBook)
}

func (h *BookHandler) listBooks(ctx context.Context, _ *struct{}) (*ListBooksOutput, error) {
	books, err := h.service.ListBooks(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}

	result := make([]BookResponse, 0, len(books))
	for _, b := range books {
		result = append(result, toBookResponse(b))
	}

	output := &ListBooksOutput{}
	output.Body.Books = result
	return output, nil
}

func (h *BookHandler) getBook(ctx context.Context, input *BookIDInput) (*GetBookOutput, error) {
	book, err := h.service.GetBook(ctx, input.ID)
	if err != nil {
		if err == repo.ErrNotFound {
			return nil, huma.NewError(http.StatusNotFound, "book not found")
		}
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}

	return &GetBookOutput{Body: toBookResponse(book)}, nil
}

func (h *BookHandler) createBook(ctx context.Context, input *CreateBookInput) (*CreateBookOutput, error) {
	book, err := h.service.CreateBook(ctx, service.BookInput{
		Title:    input.Body.Title,
		Author:   input.Body.Author,
		Price:    input.Body.Price,
		Currency: input.Body.Currency,
		Stock:    input.Body.Stock,
	})
	if err != nil {
		switch e := err.(type) {
		case service.ValidationError:
			return nil, huma.NewError(http.StatusBadRequest, "validation error", fmt.Errorf("fields: %v", e.Fields))
		default:
			return nil, huma.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	output := &CreateBookOutput{Body: toBookResponse(book)}
	return output, nil
}

func (h *BookHandler) updateBook(ctx context.Context, input *UpdateBookInput) (*UpdateBookOutput, error) {
	book, err := h.service.UpdateBook(ctx, input.ID, service.BookInput{
		Title:    input.Body.Title,
		Author:   input.Body.Author,
		Price:    input.Body.Price,
		Currency: input.Body.Currency,
		Stock:    input.Body.Stock,
	})
	if err != nil {
		switch e := err.(type) {
		case service.ValidationError:
			return nil, huma.NewError(http.StatusBadRequest, "validation error", fmt.Errorf("fields: %v", e.Fields))
		default:
			if err == repo.ErrNotFound {
				return nil, huma.NewError(http.StatusNotFound, "book not found")
			}
			return nil, huma.NewError(http.StatusInternalServerError, err.Error())
		}
	}
	return &UpdateBookOutput{Body: toBookResponse(book)}, nil
}

func (h *BookHandler) deleteBook(ctx context.Context, input *BookIDInput) (*struct{}, error) {
	err := h.service.DeleteBook(ctx, input.ID)
	if err != nil {
		if err == repo.ErrNotFound {
			return nil, huma.NewError(http.StatusNotFound, "book not found")
		}
		return nil, huma.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil, nil
}

func toBookResponse(book domain.Book) BookResponse {
	return BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		Price:     book.Price,
		Currency:  book.Currency,
		Stock:     book.Stock,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
}
