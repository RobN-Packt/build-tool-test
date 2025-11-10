package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	ErrNotFound     = errors.New("book not found")
	errCurrency     = regexp.MustCompile(`^[A-Z]{3}$`)
	ErrInvalidInput = errors.New("validation failed")
)

// Book represents the persisted book entity.
type Book struct {
	ID        string
	Title     string
	Author    string
	Price     float32
	Currency  string
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateBookParams describe the payload required to create a book.
type CreateBookParams struct {
	Title    string
	Author   string
	Price    float32
	Currency string
	Stock    int
}

// UpdateBookParams holds optional fields for book updates.
type UpdateBookParams struct {
	Title    *string
	Author   *string
	Price    *float32
	Currency *string
	Stock    *int
}

// ValidationError represents a single invalid field.
type ValidationError struct {
	Field  string
	Reason string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Reason)
}

// ValidationErrors aggregates multiple validation errors.
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	parts := make([]string, len(v))
	for i, err := range v {
		parts[i] = err.Error()
	}
	return fmt.Sprintf("%d validation error(s): %s", len(v), strings.Join(parts, "; "))
}

func (v ValidationErrors) Is(target error) bool {
	return target == ErrInvalidInput
}

// ValidateCreate validates book creation payload.
func ValidateCreate(params CreateBookParams) error {
	errs := make(ValidationErrors, 0)
	if strings.TrimSpace(params.Title) == "" {
		errs = append(errs, ValidationError{Field: "title", Reason: "cannot be empty"})
	} else if len([]rune(params.Title)) > 200 {
		errs = append(errs, ValidationError{Field: "title", Reason: "must be <= 200 characters"})
	}

	if strings.TrimSpace(params.Author) == "" {
		errs = append(errs, ValidationError{Field: "author", Reason: "cannot be empty"})
	} else if len([]rune(params.Author)) > 200 {
		errs = append(errs, ValidationError{Field: "author", Reason: "must be <= 200 characters"})
	}

	if params.Price < 0 {
		errs = append(errs, ValidationError{Field: "price", Reason: "must be >= 0"})
	}

	if params.Currency != "" && !errCurrency.MatchString(params.Currency) {
		errs = append(errs, ValidationError{Field: "currency", Reason: "must be a 3-letter ISO code"})
	}

	if params.Stock < 0 {
		errs = append(errs, ValidationError{Field: "stock", Reason: "must be >= 0"})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// ValidateUpdate validates update payload. Nil values are skipped.
func ValidateUpdate(update UpdateBookParams) error {
	errs := make(ValidationErrors, 0)

	if update.Title != nil {
		t := strings.TrimSpace(*update.Title)
		if t == "" {
			errs = append(errs, ValidationError{Field: "title", Reason: "cannot be empty"})
		} else if len([]rune(t)) > 200 {
			errs = append(errs, ValidationError{Field: "title", Reason: "must be <= 200 characters"})
		}
	}

	if update.Author != nil {
		a := strings.TrimSpace(*update.Author)
		if a == "" {
			errs = append(errs, ValidationError{Field: "author", Reason: "cannot be empty"})
		} else if len([]rune(a)) > 200 {
			errs = append(errs, ValidationError{Field: "author", Reason: "must be <= 200 characters"})
		}
	}

	if update.Price != nil && *update.Price < 0 {
		errs = append(errs, ValidationError{Field: "price", Reason: "must be >= 0"})
	}

	if update.Currency != nil && !errCurrency.MatchString(strings.ToUpper(*update.Currency)) {
		errs = append(errs, ValidationError{Field: "currency", Reason: "must be a 3-letter ISO code"})
	}

	if update.Stock != nil && *update.Stock < 0 {
		errs = append(errs, ValidationError{Field: "stock", Reason: "must be >= 0"})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// ApplyUpdate mutates the given book with values from update, assuming validation has passed.
func ApplyUpdate(book *Book, update UpdateBookParams) {
	if update.Title != nil {
		book.Title = strings.TrimSpace(*update.Title)
	}
	if update.Author != nil {
		book.Author = strings.TrimSpace(*update.Author)
	}
	if update.Price != nil {
		book.Price = *update.Price
	}
	if update.Currency != nil {
		book.Currency = strings.ToUpper(strings.TrimSpace(*update.Currency))
	}
	if update.Stock != nil {
		book.Stock = *update.Stock
	}
}
