package repository

import (
    "errors"
    "sort"
    "sync"

    "github.com/example/book-poc/apps/api/domain"
)

var (
    ErrBookNotFound     = errors.New("book not found")
    ErrDuplicateTitle   = errors.New("book title already exists")
    ErrInsufficientStock = errors.New("insufficient stock")
)

type BookRepository interface {
    List(limit int, cursor string) ([]domain.Book, *string, error)
    Get(id string) (*domain.Book, error)
    Create(book domain.Book) (*domain.Book, error)
    Update(book domain.Book) (*domain.Book, error)
    Delete(id string) error
}

type InMemoryBookRepository struct {
    mu    sync.RWMutex
    books map[string]domain.Book
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
    return &InMemoryBookRepository{
        books: make(map[string]domain.Book),
    }
}

func (r *InMemoryBookRepository) List(limit int, cursor string) ([]domain.Book, *string, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    if limit <= 0 {
        limit = 50
    }

    books := make([]domain.Book, 0, len(r.books))
    for _, b := range r.books {
        books = append(books, b)
    }
    sort.Slice(books, func(i, j int) bool {
        return books[i].CreatedAt.Before(books[j].CreatedAt)
    })

    // Simple pagination using offset encoded in cursor
    start := 0
    if cursor != "" {
        for i, b := range books {
            if b.ID == cursor {
                start = i + 1
                break
            }
        }
    }

    end := start + limit
    if end > len(books) {
        end = len(books)
    }

    slice := books[start:end]
    var next *string
    if end < len(books) {
        n := slice[len(slice)-1].ID
        next = &n
    }
    return slice, next, nil
}

func (r *InMemoryBookRepository) Get(id string) (*domain.Book, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    book, ok := r.books[id]
    if !ok {
        return nil, ErrBookNotFound
    }
    copy := book
    return &copy, nil
}

func (r *InMemoryBookRepository) Create(book domain.Book) (*domain.Book, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    for _, existing := range r.books {
        if existing.Title == book.Title {
            return nil, ErrDuplicateTitle
        }
    }

    r.books[book.ID] = book
    copy := book
    return &copy, nil
}

func (r *InMemoryBookRepository) Update(book domain.Book) (*domain.Book, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.books[book.ID]; !ok {
        return nil, ErrBookNotFound
    }

    for id, existing := range r.books {
        if id != book.ID && existing.Title == book.Title {
            return nil, ErrDuplicateTitle
        }
    }

    r.books[book.ID] = book
    copy := book
    return &copy, nil
}

func (r *InMemoryBookRepository) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.books[id]; !ok {
        return ErrBookNotFound
    }

    delete(r.books, id)
    return nil
}
