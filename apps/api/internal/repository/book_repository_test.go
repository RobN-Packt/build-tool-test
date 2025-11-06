package repository

import (
    "testing"
    "time"

    "github.com/example/book-poc/apps/api/domain"
)

func TestInMemoryBookRepository_CRUD(t *testing.T) {
    repo := NewInMemoryBookRepository()
    now := time.Now()
    book := domain.Book{ID: "1", Title: "Test", Author: "Author", Price: 10, Currency: "USD", Stock: 5, CreatedAt: now, UpdatedAt: now}

    if _, err := repo.Create(book); err != nil {
        t.Fatalf("create failed: %v", err)
    }

    fetched, err := repo.Get("1")
    if err != nil {
        t.Fatalf("get failed: %v", err)
    }
    if fetched.Title != "Test" {
        t.Errorf("expected title Test, got %s", fetched.Title)
    }

    book.Title = "Updated"
    if _, err := repo.Update(book); err != nil {
        t.Fatalf("update failed: %v", err)
    }

    if err := repo.Delete("1"); err != nil {
        t.Fatalf("delete failed: %v", err)
    }

    if _, err := repo.Get("1"); err == nil {
        t.Fatalf("expected error after delete")
    }
}

func TestInMemoryBookRepository_Duplicate(t *testing.T) {
    repo := NewInMemoryBookRepository()
    now := time.Now()
    book := domain.Book{ID: "1", Title: "Same", Author: "A", Price: 10, Currency: "USD", Stock: 1, CreatedAt: now, UpdatedAt: now}
    _, _ = repo.Create(book)

    book2 := book
    book2.ID = "2"
    if _, err := repo.Create(book2); err != ErrDuplicateTitle {
        t.Fatalf("expected duplicate title error, got %v", err)
    }
}
