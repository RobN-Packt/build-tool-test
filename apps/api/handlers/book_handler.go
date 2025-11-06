package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"

    "gofr.dev/pkg/gofr"

    "github.com/example/book-poc/apps/api/domain"
    httputil "github.com/example/book-poc/apps/api/internal/http"
    "github.com/example/book-poc/apps/api/internal/service"
)

type BookHandler struct {
    service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
    return &BookHandler{service: service}
}

func (h *BookHandler) Register(app *gofr.Gofr) {
    app.GET("/healthz", h.Health)
    app.GET("/books", h.List)
    app.POST("/books", h.Create)
    app.GET("/books/{id}", h.Get)
    app.PUT("/books/{id}", h.Update)
    app.DELETE("/books/{id}", h.Delete)
    app.POST("/books/{id}/purchase", h.Purchase)
}

func (h *BookHandler) Health(ctx *gofr.Context) (interface{}, error) {
    return map[string]string{"status": "ok"}, nil
}

func (h *BookHandler) List(ctx *gofr.Context) (interface{}, error) {
    limit := intFromQuery(ctx, "limit", 50)
    cursor := ctx.Request().URL.Query().Get("cursor")
    books, next, err := h.service.List(ctx.Request().Context(), limit, cursor)
    if err != nil {
        return errorResponse(ctx, err)
    }
    resp := map[string]any{"data": books}
    if next != nil {
        resp["next_cursor"] = *next
    }
    return resp, nil
}

func (h *BookHandler) Get(ctx *gofr.Context) (interface{}, error) {
    id := ctx.Param("id")
    book, err := h.service.Get(ctx.Request().Context(), id)
    if err != nil {
        return errorResponse(ctx, err)
    }
    return map[string]any{"data": book}, nil
}

func (h *BookHandler) Create(ctx *gofr.Context) (interface{}, error) {
    var input domain.CreateBookInput
    if err := json.NewDecoder(ctx.Request().Body).Decode(&input); err != nil {
        return badRequest(ctx, err)
    }
    book, err := h.service.Create(ctx.Request().Context(), input)
    if err != nil {
        return errorResponse(ctx, err)
    }
    ctx.Response().Status(http.StatusCreated)
    return map[string]any{"data": book}, nil
}

func (h *BookHandler) Update(ctx *gofr.Context) (interface{}, error) {
    var input domain.UpdateBookInput
    if err := json.NewDecoder(ctx.Request().Body).Decode(&input); err != nil {
        return badRequest(ctx, err)
    }
    id := ctx.Param("id")
    book, err := h.service.Update(ctx.Request().Context(), id, input)
    if err != nil {
        return errorResponse(ctx, err)
    }
    return map[string]any{"data": book}, nil
}

func (h *BookHandler) Delete(ctx *gofr.Context) (interface{}, error) {
    id := ctx.Param("id")
    if err := h.service.Delete(ctx.Request().Context(), id); err != nil {
        return errorResponse(ctx, err)
    }
    ctx.Response().Status(http.StatusNoContent)
    return nil, nil
}

func (h *BookHandler) Purchase(ctx *gofr.Context) (interface{}, error) {
    id := ctx.Param("id")
    var body struct {
        Quantity   int    `json:"quantity"`
        CustomerID string `json:"customer_id"`
    }
    if err := json.NewDecoder(ctx.Request().Body).Decode(&body); err != nil {
        return badRequest(ctx, err)
    }
    messageID, err := h.service.Purchase(ctx.Request().Context(), domain.PurchaseInput{BookID: id, Quantity: body.Quantity, CustomerID: body.CustomerID})
    if err != nil {
        return errorResponse(ctx, err)
    }
    ctx.Response().Status(http.StatusAccepted)
    return map[string]any{"message_id": messageID, "status": "accepted"}, nil
}

func intFromQuery(ctx *gofr.Context, key string, def int) int {
    raw := ctx.Request().URL.Query().Get(key)
    if raw == "" {
        return def
    }
    v, err := strconv.Atoi(raw)
    if err != nil {
        return def
    }
    return v
}

func errorResponse(ctx *gofr.Context, err error) (interface{}, error) {
    ctx.Response().Status(httputil.StatusCode(err))
    return httputil.ErrorResponse{Error: err.Error(), Code: httputil.Code(err)}, nil
}

func badRequest(ctx *gofr.Context, err error) (interface{}, error) {
    ctx.Response().Status(http.StatusBadRequest)
    return httputil.ErrorResponse{Error: err.Error(), Code: "bad_request"}, nil
}

func (h *BookHandler) WithContext(ctx context.Context) context.Context {
    return ctx
}
