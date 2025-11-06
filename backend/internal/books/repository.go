package books

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrDuplicateBook = errors.New("book already exists")
	ErrBookNotFound  = errors.New("book not found")
)

// Repository exposes persistence operations for books.
type Repository interface {
	List(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id int64) (Book, error)
	Create(ctx context.Context, book Book) (Book, error)
	Update(ctx context.Context, book Book) (Book, error)
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *sql.DB
}

// NewRepository returns a new Repository instance backed by the provided DB.
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context) ([]Book, error) {
	query := `
        SELECT id, title, author, isbn, price, stock, description, published_date, created_at, updated_at
        FROM books
        ORDER BY id;
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		book, err := scanRow(rows)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (Book, error) {
	query := `
        SELECT id, title, author, isbn, price, stock, description, published_date, created_at, updated_at
        FROM books
        WHERE id = $1;
    `

	row := r.db.QueryRowContext(ctx, query, id)
	book, err := scanRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Book{}, ErrBookNotFound
		}

		return Book{}, err
	}

	return book, nil
}

func (r *repository) Create(ctx context.Context, book Book) (Book, error) {
	query := `
        INSERT INTO books (title, author, isbn, price, stock, description, published_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at;
    `

	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query,
		book.Title,
		book.Author,
		book.ISBN,
		book.Price,
		book.Stock,
		nullIfEmpty(book.Description),
		book.PublishedDate,
	).Scan(&book.ID, &createdAt, &updatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return Book{}, ErrDuplicateBook
		}

		return Book{}, err
	}

	book.CreatedAt = createdAt
	book.UpdatedAt = updatedAt

	return book, nil
}

func (r *repository) Update(ctx context.Context, book Book) (Book, error) {
	query := `
        UPDATE books
        SET title = $1,
            author = $2,
            isbn = $3,
            price = $4,
            stock = $5,
            description = $6,
            published_date = $7,
            updated_at = NOW()
        WHERE id = $8
        RETURNING created_at, updated_at;
    `

	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query,
		book.Title,
		book.Author,
		book.ISBN,
		book.Price,
		book.Stock,
		nullIfEmpty(book.Description),
		book.PublishedDate,
		book.ID,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Book{}, ErrBookNotFound
		}

		if isUniqueViolation(err) {
			return Book{}, ErrDuplicateBook
		}

		return Book{}, err
	}

	book.CreatedAt = createdAt
	book.UpdatedAt = updatedAt

	return book, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM books WHERE id = $1;`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrBookNotFound
	}

	return nil
}

func scanRow(scanner interface{ Scan(dest ...any) error }) (Book, error) {
	var (
		book        Book
		description sql.NullString
	)

	err := scanner.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Price,
		&book.Stock,
		&description,
		&book.PublishedDate,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return Book{}, err
	}

	if description.Valid {
		book.Description = description.String
	}

	return book, nil
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return sql.NullString{}
	}

	return value
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	// postgres unique violation error code: 23505
	// we perform string contains check to avoid driver specific types.
	return strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
		strings.Contains(err.Error(), "23505")
}
