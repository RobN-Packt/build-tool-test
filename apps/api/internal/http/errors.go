package httputil

import (
    "errors"
    "net/http"

    apierrors "github.com/example/book-poc/apps/api/internal/repository"
)

type ErrorResponse struct {
    Error string `json:"error"`
    Code  string `json:"code"`
    Details any `json:"details,omitempty"`
}

func StatusCode(err error) int {
    switch {
    case err == nil:
        return http.StatusOK
    case errors.Is(err, apierrors.ErrBookNotFound):
        return http.StatusNotFound
    case errors.Is(err, apierrors.ErrDuplicateTitle):
        return http.StatusConflict
    case errors.Is(err, apierrors.ErrInsufficientStock):
        return http.StatusBadRequest
    default:
        return http.StatusBadRequest
    }
}

func Code(err error) string {
    switch {
    case errors.Is(err, apierrors.ErrBookNotFound):
        return "not_found"
    case errors.Is(err, apierrors.ErrDuplicateTitle):
        return "duplicate"
    case errors.Is(err, apierrors.ErrInsufficientStock):
        return "insufficient_stock"
    default:
        return "bad_request"
    }
}
