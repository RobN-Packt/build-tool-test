package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/example/bookapi/internal/domain"
	"github.com/example/bookapi/internal/repo"
	"github.com/example/bookapi/internal/service"
	"github.com/example/bookapi/openapi"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

type BookIDInput struct {
	ID uuid.UUID `path:"id"`
}

type CreateBookInput struct {
	Body openapi.BookCreate `body:""`
}

type CreateBookOutput struct {
	Body openapi.Book
}

type GetBookOutput struct {
	Body openapi.Book
}

type ListBooksOutput struct {
	Body struct {
		Books []openapi.Book `json:"books"`
	}
}

type UpdateBookInput struct {
	ID   uuid.UUID          `path:"id"`
	Body openapi.BookUpdate `body:""`
}

type UpdateBookOutput struct {
	Body openapi.Book
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

	result := make([]openapi.Book, 0, len(books))
	for _, b := range books {
		result = append(result, toOpenAPIBook(b))
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

	return &GetBookOutput{Body: toOpenAPIBook(book)}, nil
}

func (h *BookHandler) createBook(ctx context.Context, input *CreateBookInput) (*CreateBookOutput, error) {
	book, err := h.service.CreateBook(ctx, toServiceCreateInput(input.Body))
	if err != nil {
		switch e := err.(type) {
		case service.ValidationError:
			return nil, huma.NewError(http.StatusBadRequest, "validation error", fmt.Errorf("fields: %v", e.Fields))
		default:
			return nil, huma.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	output := &CreateBookOutput{Body: toOpenAPIBook(book)}
	return output, nil
}

func (h *BookHandler) updateBook(ctx context.Context, input *UpdateBookInput) (*UpdateBookOutput, error) {
	book, err := h.service.UpdateBook(ctx, input.ID, toServiceUpdateInput(input.Body))
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
	return &UpdateBookOutput{Body: toOpenAPIBook(book)}, nil
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

func toOpenAPIBook(book domain.Book) openapi.Book {
	return openapi.Book{
		Id:        openapi_types.UUID(book.ID),
		Title:     book.Title,
		Author:    book.Author,
		Price:     float32(book.Price),
		Currency:  book.Currency,
		Stock:     book.Stock,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	}
}

func toServiceCreateInput(body openapi.BookCreate) service.BookCreateInput {
	return service.BookCreateInput{
		Title:    body.Title,
		Author:   body.Author,
		Price:    float64(body.Price),
		Currency: body.Currency,
		Stock:    body.Stock,
	}
}

func toServiceUpdateInput(body openapi.BookUpdate) service.BookUpdateInput {
	var result service.BookUpdateInput
	if body.Title != nil {
		value := *body.Title
		result.Title = &value
	}
	if body.Author != nil {
		value := *body.Author
		result.Author = &value
	}
	if body.Price != nil {
		value := float64(*body.Price)
		result.Price = &value
	}
	if body.Currency != nil {
		value := *body.Currency
		result.Currency = &value
	}
	if body.Stock != nil {
		value := *body.Stock
		result.Stock = &value
	}
	return result
}
