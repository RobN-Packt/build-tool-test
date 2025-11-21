package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestBookCRUDIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if strings.TrimSpace(dsn) == "" {
		t.Skip("TEST_DB_DSN not set; skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Skipf("unable to create pool for TEST_DB_DSN: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("unable to connect to TEST_DB_DSN: %v", err)
	}
	defer pool.Close()

	require.NoError(t, applyMigrations(ctx, pool))

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "TRUNCATE TABLE books")
	})

	server := httptest.NewServer(buildHTTPHandler(pool, nil))
	defer server.Close()

	client := server.Client()

	createPayload := map[string]any{
		"title":    "Integration Testing with Go",
		"author":   "John Doe",
		"price":    24.99,
		"currency": "usd",
		"stock":    3,
	}
	createBody, err := json.Marshal(createPayload)
	require.NoError(t, err)

	createResp, err := client.Post(server.URL+"/books", "application/json", bytes.NewReader(createBody))
	require.NoError(t, err)
	defer createResp.Body.Close()
	require.Equal(t, http.StatusCreated, createResp.StatusCode)

	var created bookResponse
	require.NoError(t, json.NewDecoder(createResp.Body).Decode(&created))
	_, err = uuid.Parse(created.ID)
	require.NoError(t, err)
	require.Equal(t, "Integration Testing with Go", created.Title)
	require.Equal(t, "USD", created.Currency)
	require.True(t, created.CreatedAt.After(time.Time{}))

	getResp, err := client.Get(server.URL + "/books/" + created.ID)
	require.NoError(t, err)
	defer getResp.Body.Close()
	require.Equal(t, http.StatusOK, getResp.StatusCode)

	var fetched bookResponse
	require.NoError(t, json.NewDecoder(getResp.Body).Decode(&fetched))
	require.Equal(t, created.ID, fetched.ID)

	listResp, err := client.Get(server.URL + "/books")
	require.NoError(t, err)
	defer listResp.Body.Close()
	require.Equal(t, http.StatusOK, listResp.StatusCode)

	var listBody struct {
		Books []bookResponse `json:"books"`
	}
	require.NoError(t, json.NewDecoder(listResp.Body).Decode(&listBody))
	require.Len(t, listBody.Books, 1)

	updatePayload := map[string]any{
		"title":    "Integration Testing with Go - Second Edition",
		"author":   "John Doe",
		"price":    29.99,
		"currency": "EUR",
		"stock":    7,
	}
	updateBody, err := json.Marshal(updatePayload)
	require.NoError(t, err)

	updateReq, err := http.NewRequest(http.MethodPut, server.URL+"/books/"+created.ID, bytes.NewReader(updateBody))
	require.NoError(t, err)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := client.Do(updateReq)
	require.NoError(t, err)
	defer updateResp.Body.Close()
	require.Equal(t, http.StatusOK, updateResp.StatusCode)

	var updated bookResponse
	require.NoError(t, json.NewDecoder(updateResp.Body).Decode(&updated))
	require.Equal(t, "Integration Testing with Go - Second Edition", updated.Title)
	require.Equal(t, "EUR", updated.Currency)
	require.Greater(t, updated.Stock, created.Stock)

	deleteReq, err := http.NewRequest(http.MethodDelete, server.URL+"/books/"+created.ID, nil)
	require.NoError(t, err)
	deleteResp, err := client.Do(deleteReq)
	require.NoError(t, err)
	defer deleteResp.Body.Close()
	require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	finalGetResp, err := client.Get(server.URL + "/books/" + created.ID)
	require.NoError(t, err)
	defer finalGetResp.Body.Close()
	require.Equal(t, http.StatusNotFound, finalGetResp.StatusCode)
}

func TestHealthEndpoint(t *testing.T) {
	handler := buildHTTPHandler(nil, nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOpenAPIExposed(t *testing.T) {
	handler := buildHTTPHandler(nil, nil)
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/openapi.json")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Contains(t, resp.Header.Get("Content-Type"), "openapi+json")
}

type bookResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
