package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/example/bookshop/apps/api/internal/domain"
	"github.com/example/bookshop/apps/api/internal/service"
	"github.com/example/bookshop/apps/api/openapi"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// BookHandler wires Huma operations to the book service.
type BookHandler struct {
	Service *service.BookService
}

type ListBooksOutput struct {
	Body []openapi.Book
}

type GetBookInput struct {
	ID string `path:"id" doc:"Book ID"`
}

type GetBookOutput struct {
	Body openapi.Book
}

type CreateBookInput struct {
	Body openapi.BookCreate
}

type CreateBookOutput struct {
	Body openapi.Book
}

type UpdateBookInput struct {
	ID   string `path:"id" doc:"Book ID"`
	Body openapi.BookUpdate
}

type UpdateBookOutput struct {
	Body openapi.Book
}

type DeleteBookInput struct {
	ID string `path:"id" doc:"Book ID"`
}

type DeleteBookOutput struct{}

// Register attaches all book routes to the given API.
func Register(api huma.API, handler *BookHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "listBooks",
		Summary:       "List books",
		Method:        http.MethodGet,
		Path:          "/books",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusOK,
	}, handler.list)

	huma.Register(api, huma.Operation{
		OperationID:   "createBook",
		Summary:       "Create book",
		Method:        http.MethodPost,
		Path:          "/books",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusCreated,
	}, handler.create)

	huma.Register(api, huma.Operation{
		OperationID:   "getBook",
		Summary:       "Get book",
		Method:        http.MethodGet,
		Path:          "/books/{id}",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusOK,
	}, handler.get)

	huma.Register(api, huma.Operation{
		OperationID:   "updateBook",
		Summary:       "Update book",
		Method:        http.MethodPut,
		Path:          "/books/{id}",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusOK,
	}, handler.update)

	huma.Register(api, huma.Operation{
		OperationID:   "deleteBook",
		Summary:       "Delete book",
		Method:        http.MethodDelete,
		Path:          "/books/{id}",
		Tags:          []string{"Books"},
		DefaultStatus: http.StatusNoContent,
	}, handler.delete)
}

func (h *BookHandler) list(ctx context.Context, _ *struct{}) (*ListBooksOutput, error) {
	books, err := h.Service.ListBooks(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &ListBooksOutput{Body: make([]openapi.Book, 0, len(books))}
	for _, book := range books {
		resp.Body = append(resp.Body, toOpenAPI(book))
	}
	return resp, nil
}

func (h *BookHandler) get(ctx context.Context, input *GetBookInput) (*GetBookOutput, error) {
	book, err := h.Service.GetBook(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &GetBookOutput{Body: toOpenAPI(book)}, nil
}

func (h *BookHandler) create(ctx context.Context, input *CreateBookInput) (*CreateBookOutput, error) {
	params := domain.CreateBookParams{
		Title:    input.Body.Title,
		Author:   input.Body.Author,
		Price:    input.Body.Price,
		Currency: valueOrEmpty(input.Body.Currency),
		Stock:    valueOrZero(input.Body.Stock),
	}
	book, err := h.Service.CreateBook(ctx, params)
	if err != nil {
		return nil, mapError(err)
	}
	return &CreateBookOutput{Body: toOpenAPI(book)}, nil
}

func (h *BookHandler) update(ctx context.Context, input *UpdateBookInput) (*UpdateBookOutput, error) {
	update := domain.UpdateBookParams{
		Title:    input.Body.Title,
		Author:   input.Body.Author,
		Price:    input.Body.Price,
		Currency: input.Body.Currency,
		Stock:    input.Body.Stock,
	}
	book, err := h.Service.UpdateBook(ctx, input.ID, update)
	if err != nil {
		return nil, mapError(err)
	}
	return &UpdateBookOutput{Body: toOpenAPI(book)}, nil
}

func (h *BookHandler) delete(ctx context.Context, input *DeleteBookInput) (*DeleteBookOutput, error) {
	if err := h.Service.DeleteBook(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteBookOutput{}, nil
}

func toOpenAPI(book domain.Book) openapi.Book {
	return openapi.Book{
		Id:        mustUUID(book.ID),
		Title:     book.Title,
		Author:    book.Author,
		Price:     book.Price,
		Currency:  book.Currency,
		Stock:     book.Stock,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
}

func valueOrEmpty(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func mustUUID(id string) openapi_types.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return openapi_types.UUID(uuid.Nil)
	}
	return openapi_types.UUID(parsed)
}

func valueOrZero(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	var vErrs domain.ValidationErrors
	if errors.As(err, &vErrs) {
		details := make([]error, len(vErrs))
		for i, v := range vErrs {
			details[i] = fmt.Errorf("%s: %s", v.Field, v.Reason)
		}
		return huma.Error400BadRequest("validation failed", details...)
	}

	if errors.Is(err, domain.ErrNotFound) {
		return huma.Error404NotFound("book not found")
	}

	return huma.Error500InternalServerError("internal server error")
}
