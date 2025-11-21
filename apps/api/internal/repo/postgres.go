package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/example/bookapi/internal/domain"
)

var ErrNotFound = errors.New("book not found")

type BookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{pool: pool}
}

func (r *BookRepository) Create(ctx context.Context, book domain.Book) error {
	const query = `
		INSERT INTO books (id, title, author, price, currency, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		book.ID,
		book.Title,
		book.Author,
		book.Price,
		book.Currency,
		book.Stock,
		book.CreatedAt,
		book.UpdatedAt,
	)
	return err
}

func (r *BookRepository) Get(ctx context.Context, id uuid.UUID) (domain.Book, error) {
	const query = `
		SELECT id, title, author, price, currency, stock, created_at, updated_at
		FROM books
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)

	book, err := scanBook(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, ErrNotFound
		}
		return domain.Book{}, err
	}
	return book, nil
}

func (r *BookRepository) List(ctx context.Context) ([]domain.Book, error) {
	const query = `
		SELECT id, title, author, price, currency, stock, created_at, updated_at
		FROM books
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		book, err := scanBook(rows)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return books, nil
}

func (r *BookRepository) Update(ctx context.Context, book domain.Book) error {
	const query = `
		UPDATE books
		SET title = $2,
			author = $3,
			price = $4,
			currency = $5,
			stock = $6,
			updated_at = $7
		WHERE id = $1
	`
	tag, err := r.pool.Exec(ctx, query,
		book.ID,
		book.Title,
		book.Author,
		book.Price,
		book.Currency,
		book.Stock,
		book.UpdatedAt,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *BookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM books WHERE id = $1`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *BookRepository) CheckBooks(ctx context.Context) error {
	const query = `SELECT 1 FROM books LIMIT 1`
	var sentinel int
	err := r.pool.QueryRow(ctx, query).Scan(&sentinel)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("books health query: %w", err)
	}
	return nil
}

func scanBook(row pgx.Row) (domain.Book, error) {
	var book domain.Book
	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Price,
		&book.Currency,
		&book.Stock,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		return domain.Book{}, fmt.Errorf("scan book: %w", err)
	}
	return book, nil
}
