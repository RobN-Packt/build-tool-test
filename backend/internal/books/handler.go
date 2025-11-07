package books

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	service Service
}

type bookResource struct {
	ID            int64   `json:"id" example:"1"`
	Title         string  `json:"title" example:"Clean Code"`
	Author        string  `json:"author" example:"Robert C. Martin"`
	ISBN          string  `json:"isbn" example:"9780132350884"`
	Price         float64 `json:"price" example:"39.5"`
	Stock         int     `json:"stock" example:"12"`
	Description   *string `json:"description,omitempty" example:"Agile software craftsmanship techniques"`
	PublishedDate string  `json:"publishedDate" format:"date" example:"2008-08-01"`
	CreatedAt     string  `json:"createdAt" format:"date-time"`
	UpdatedAt     string  `json:"updatedAt" format:"date-time"`
}

type listBooksResponse struct {
	Books []bookResource `json:"books"`
}

type bookOKResponse struct {
	Body bookResource
}

type listBooksOKResponse struct {
	Body listBooksResponse
}

type deleteBookResponse struct {
	Status int
}

// RegisterRoutes binds the books routes to the provided Huma API.
func RegisterRoutes(api huma.API, service Service) {
	h := handler{service: service}

	huma.Register(api, huma.Operation{
		OperationID:   "list-books",
		Method:        http.MethodGet,
		Path:          "/books",
		Summary:       "List books",
		Description:   "Retrieve the catalog of books that are currently in stock.",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusOK,
	}, h.list)

	huma.Register(api, huma.Operation{
		OperationID: "get-book",
		Method:      http.MethodGet,
		Path:        "/books/{id}",
		Summary:     "Get book",
		Tags:        []string{"Books"},
	}, h.get)

	huma.Register(api, huma.Operation{
		OperationID:   "create-book",
		Method:        http.MethodPost,
		Path:          "/books",
		Summary:       "Create book",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusCreated,
	}, h.create)

	huma.Register(api, huma.Operation{
		OperationID: "update-book",
		Method:      http.MethodPut,
		Path:        "/books/{id}",
		Summary:     "Update book",
		Tags:        []string{"Books"},
	}, h.update)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-book",
		Method:        http.MethodDelete,
		Path:          "/books/{id}",
		Summary:       "Delete book",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusNoContent,
	}, h.delete)
}

func (h handler) list(ctx context.Context, _ *struct{}) (*listBooksOKResponse, error) {
	items, err := h.service.List(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to list books", err)
	}

	return &listBooksOKResponse{
		Body: listBooksResponse{Books: mapBooks(items)},
	}, nil
}

func (h handler) get(ctx context.Context, input *struct {
	ID int64 `path:"id" example:"1"`
}) (*bookOKResponse, error) {
	book, err := h.service.Get(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}

	return &bookOKResponse{Body: mapBook(book)}, nil
}

func (h handler) create(ctx context.Context, input *struct {
	Body CreateBookInput
}) (*bookOKResponse, error) {
	book, err := h.service.Create(ctx, input.Body)
	if err != nil {
		return nil, mapError(err)
	}

	return &bookOKResponse{
		Body: mapBook(book),
	}, nil
}

func (h handler) update(ctx context.Context, input *struct {
	ID   int64 `path:"id" example:"1"`
	Body UpdateBookInput
}) (*bookOKResponse, error) {
	book, err := h.service.Update(ctx, input.ID, input.Body)
	if err != nil {
		return nil, mapError(err)
	}

	return &bookOKResponse{Body: mapBook(book)}, nil
}

func (h handler) delete(ctx context.Context, input *struct {
	ID int64 `path:"id" example:"1"`
}) (*deleteBookResponse, error) {
	if err := h.service.Delete(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}

	return &deleteBookResponse{Status: http.StatusNoContent}, nil
}

func mapBooks(items []Book) []bookResource {
	books := make([]bookResource, len(items))
	for i, item := range items {
		books[i] = mapBook(item)
	}
	return books
}

func mapBook(book Book) bookResource {
	var description *string
	trimmed := strings.TrimSpace(book.Description)
	if trimmed != "" {
		description = &trimmed
	}

	return bookResource{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		ISBN:          book.ISBN,
		Price:         book.Price,
		Stock:         book.Stock,
		Description:   description,
		PublishedDate: book.PublishedDate.Format("2006-01-02"),
		CreatedAt:     book.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     book.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, ErrBookNotFound):
		return huma.Error404NotFound("book not found", err)
	case errors.Is(err, ErrDuplicateBook):
		return huma.Error409Conflict("book already exists", err)
	case errors.Is(err, ErrInvalidInput):
		return huma.Error400BadRequest("invalid book payload", err)
	case errors.Is(err, ErrInvalidPublishedDate):
		return huma.Error400BadRequest("invalid published date", err)
	default:
		return huma.Error500InternalServerError("unexpected error", err)
	}
}
