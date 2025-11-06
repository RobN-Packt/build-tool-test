package books

import (
	"errors"
	"strconv"

	"gofr.dev/pkg/gofr"
	gofrHTTP "gofr.dev/pkg/gofr/http"
)

// Handler wires HTTP requests to the books service.
type Handler struct {
	service Service
}

// RegisterRoutes binds the books routes to the provided GoFr app.
func RegisterRoutes(app *gofr.App, service Service) {
	h := Handler{service: service}

	app.GET("/books", h.list)
	app.GET("/books/{id}", h.get)
	app.POST("/books", h.create)
	app.PUT("/books/{id}", h.update)
	app.DELETE("/books/{id}", h.delete)
}

func (h Handler) list(c *gofr.Context) (interface{}, error) {
	books, err := h.service.List(c)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"data": books}, nil
}

func (h Handler) get(c *gofr.Context) (interface{}, error) {
	id, err := parseID(c.PathParam("id"))
	if err != nil {
		return nil, err
	}

	value := c.PathParam("id")

	book, err := h.service.Get(c, id)
	if err != nil {
		return nil, translateError("id", value, err)
	}

	return book, nil
}

func (h Handler) create(c *gofr.Context) (interface{}, error) {
	var payload CreateBookInput

	if err := c.Bind(&payload); err != nil {
		return nil, gofrHTTP.ErrorInvalidParam{Params: []string{"body"}}
	}

	book, err := h.service.Create(c, payload)
	if err != nil {
		return nil, translateError("body", "", err)
	}

	return book, nil
}

func (h Handler) update(c *gofr.Context) (interface{}, error) {
	id, err := parseID(c.PathParam("id"))
	if err != nil {
		return nil, err
	}

	var payload UpdateBookInput
	if err := c.Bind(&payload); err != nil {
		return nil, gofrHTTP.ErrorInvalidParam{Params: []string{"body"}}
	}

	value := c.PathParam("id")

	book, err := h.service.Update(c, id, payload)
	if err != nil {
		return nil, translateError("id", value, err)
	}

	return book, nil
}

func (h Handler) delete(c *gofr.Context) (interface{}, error) {
	id, err := parseID(c.PathParam("id"))
	if err != nil {
		return nil, err
	}

	value := c.PathParam("id")

	if err := h.service.Delete(c, id); err != nil {
		return nil, translateError("id", value, err)
	}

	return map[string]string{"message": "book deleted"}, nil
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, gofrHTTP.ErrorInvalidParam{Params: []string{"id"}}
	}

	return id, nil
}

func translateError(field, value string, err error) error {
	switch err {
	case ErrBookNotFound:
		return gofrHTTP.ErrorEntityNotFound{Name: field, Value: value}
	case ErrDuplicateBook:
		return gofrHTTP.ErrorEntityAlreadyExist{}
	}

	if errors.Is(err, ErrInvalidInput) || errors.Is(err, ErrInvalidPublishedDate) {
		return gofrHTTP.ErrorInvalidParam{Params: []string{field}}
	}

	return err
}
