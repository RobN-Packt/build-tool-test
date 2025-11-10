//go:build integration

package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/example/bookshop/apps/api/internal/http/handlers"
	appmw "github.com/example/bookshop/apps/api/internal/http/middleware"
	"github.com/example/bookshop/apps/api/internal/repo"
	"github.com/example/bookshop/apps/api/internal/service"
	"github.com/example/bookshop/apps/api/openapi"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBookCRUDIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set")
	}

	ctx := context.Background()
	pool, err := repo.OpenPool(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })

	require.NoError(t, repo.Migrate(ctx, pool))
	require.NoError(t, truncateBooks(ctx, pool))

	repository := repo.NewBookRepository(pool)
	svc := service.NewBookService(repository)

	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("Book Shop API", "1.0.0"))
	handlers.Register(api, &handlers.BookHandler{Service: svc})

	ts := httptest.NewServer(appmw.Logger(router))
	t.Cleanup(ts.Close)
	client := ts.Client()

	createBody := map[string]any{
		"title":    "Integration Testing in Go",
		"author":   "Jane Developer",
		"price":    39.99,
		"currency": "USD",
		"stock":    10,
	}
	created := postBook(t, client, ts.URL+"/books", createBody)
	bookID := created.Id.String()
	require.NotEmpty(t, bookID)

	fetched := getBook(t, client, ts.URL+"/books/"+bookID)
	require.Equal(t, bookID, fetched.Id.String())
	require.Equal(t, "Jane Developer", fetched.Author)

	listed := listBooks(t, client, ts.URL+"/books")
	require.NotEmpty(t, listed)

	updateBody := map[string]any{
		"price": 29.99,
		"stock": 5,
	}
	updated := putBook(t, client, ts.URL+"/books/"+bookID, updateBody)
	require.Equal(t, float32(29.99), updated.Price)
	require.Equal(t, 5, updated.Stock)
	require.True(t, updated.UpdatedAt.After(updated.CreatedAt) || updated.UpdatedAt.Equal(updated.CreatedAt))

	deleteBook(t, client, ts.URL+"/books/"+bookID)
	req, err := http.NewRequest(http.MethodGet, ts.URL+"/books/"+bookID, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func postBook(t *testing.T, client *http.Client, url string, body map[string]any) openapi.Book {
	t.Helper()
	data, err := json.Marshal(body)
	require.NoError(t, err)
	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	var book openapi.Book
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&book))
	return book
}

func getBook(t *testing.T, client *http.Client, url string) openapi.Book {
	t.Helper()
	resp, err := client.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var book openapi.Book
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&book))
	return book
}

func listBooks(t *testing.T, client *http.Client, url string) []openapi.Book {
	t.Helper()
	resp, err := client.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var books []openapi.Book
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&books))
	return books
}

func putBook(t *testing.T, client *http.Client, url string, body map[string]any) openapi.Book {
	t.Helper()
	data, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var book openapi.Book
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&book))
	return book
}

func deleteBook(t *testing.T, client *http.Client, url string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func truncateBooks(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, "TRUNCATE TABLE books")
	return err
}
