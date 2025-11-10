package repo

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/example/bookshop/apps/api/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func OpenPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return pool, nil
}

// BookRepository provides persistence for books backed by Postgres.
type BookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{pool: pool}
}

// Migrate executes embedded SQL migrations in lexical order.
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, name := range files {
		content, err := migrationFiles.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		statements := splitSQL(string(content))
		for _, stmt := range statements {
			if strings.TrimSpace(stmt) == "" {
				continue
			}
			if _, err := pool.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("exec migration %s: %w", name, err)
			}
		}
	}

	return nil
}

func splitSQL(sql string) []string {
	replaced := strings.ReplaceAll(sql, "\r", "")
	parts := strings.Split(replaced, ";")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func scanBook(row pgx.Row) (domain.Book, error) {
	var book domain.Book
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Currency, &book.Stock, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Book{}, domain.ErrNotFound
		}
		return domain.Book{}, err
	}
	book.CreatedAt = createdAt.UTC()
	book.UpdatedAt = updatedAt.UTC()
	return book, nil
}

// List returns all books ordered by created_at descending.
func (r *BookRepository) List(ctx context.Context) ([]domain.Book, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, author, price, currency, stock, created_at, updated_at
		FROM books
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := make([]domain.Book, 0)
	for rows.Next() {
		book, err := scanBook(rows)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, rows.Err()
}

// Get retrieves a single book by ID.
func (r *BookRepository) Get(ctx context.Context, id string) (domain.Book, error) {
	return scanBook(r.pool.QueryRow(ctx, `
		SELECT id, title, author, price, currency, stock, created_at, updated_at
		FROM books
		WHERE id = $1`, id))
}

// Create inserts a new book record.
func (r *BookRepository) Create(ctx context.Context, book domain.Book) (domain.Book, error) {
	if book.ID == "" {
		book.ID = uuid.NewString()
	}
	return scanBook(r.pool.QueryRow(ctx, `
		INSERT INTO books (id, title, author, price, currency, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, author, price, currency, stock, created_at, updated_at`,
		book.ID, book.Title, book.Author, book.Price, book.Currency, book.Stock, book.CreatedAt, book.UpdatedAt))
}

// Update replaces all book fields.
func (r *BookRepository) Update(ctx context.Context, book domain.Book) (domain.Book, error) {
	return scanBook(r.pool.QueryRow(ctx, `
		UPDATE books
		SET title = $2,
		    author = $3,
		    price = $4,
		    currency = $5,
		    stock = $6,
		    created_at = $7,
		    updated_at = $8
		WHERE id = $1
		RETURNING id, title, author, price, currency, stock, created_at, updated_at`,
		book.ID, book.Title, book.Author, book.Price, book.Currency, book.Stock, book.CreatedAt, book.UpdatedAt))
}

// Delete removes a book by ID.
func (r *BookRepository) Delete(ctx context.Context, id string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
