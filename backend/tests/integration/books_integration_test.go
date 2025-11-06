package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"gofr.dev/pkg/gofr"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/cursor/bookshop/backend/internal/books"
	"github.com/cursor/bookshop/backend/internal/database"
	"github.com/cursor/bookshop/backend/internal/database/migrations"
)

func TestBooksCRUD(t *testing.T) {
	ctx := context.Background()

	ensureDockerAvailable(t, ctx)

	pgContainer := startPostgresContainer(ctx, t)
	defer func() {
		require.NoError(t, pgContainer.Terminate(ctx))
	}()

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	mappedPort, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dbCfg := database.Config{
		Host:     host,
		Port:     mappedPort.Port(),
		User:     "bookshop",
		Password: "bookshop",
		Name:     "bookshop",
		SSLMode:  "disable",
	}

	httpPort := freePort(t)
	previousHTTPPort := os.Getenv("HTTP_PORT")
	require.NoError(t, os.Setenv("HTTP_PORT", fmt.Sprintf("%d", httpPort)))
	defer func() {
		if previousHTTPPort == "" {
			require.NoError(t, os.Unsetenv("HTTP_PORT"))
		} else {
			require.NoError(t, os.Setenv("HTTP_PORT", previousHTTPPort))
		}
	}()

	app := gofr.New()

	db, err := database.Connect(ctx, dbCfg)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, migrations.Apply(ctx, db))

	repo := books.NewRepository(db)
	service := books.NewService(repo)
	books.RegisterRoutes(app, service)

	go app.Run()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		require.NoError(t, app.Shutdown(shutdownCtx))
	}()

	baseURL := fmt.Sprintf("http://localhost:%d", httpPort)
	waitForServer(t, baseURL+"/.well-known/health")

	// List seeded books
	resp := doRequest(t, http.MethodGet, baseURL+"/books", nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var list struct {
		Data []books.Book `json:"data"`
	}
	decodeResponse(t, resp, &list)
	require.GreaterOrEqual(t, len(list.Data), 5)

	// Create a new book
	createBody := books.CreateBookInput{
		Title:         "Integration Testing in Go",
		Author:        "Jane Doe",
		ISBN:          "9991112223334",
		Price:         29.99,
		Stock:         7,
		Description:   "Hands-on guide to integration testing",
		PublishedDate: "2024-01-15",
	}

	resp = doRequest(t, http.MethodPost, baseURL+"/books", createBody)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var created books.Book
	decodeResponse(t, resp, &created)
	require.NotZero(t, created.ID)
	require.Equal(t, createBody.Title, created.Title)

	// Update the book
	updateBody := books.UpdateBookInput{
		Title:         "Integration Testing in Go - Updated",
		Author:        createBody.Author,
		ISBN:          createBody.ISBN,
		Price:         31.99,
		Stock:         9,
		Description:   createBody.Description,
		PublishedDate: createBody.PublishedDate,
	}

	resp = doRequest(t, http.MethodPut, fmt.Sprintf("%s/books/%d", baseURL, created.ID), updateBody)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var updated books.Book
	decodeResponse(t, resp, &updated)
	require.Equal(t, updateBody.Title, updated.Title)
	require.Equal(t, updateBody.Price, updated.Price)

	// Delete the book
	resp = doRequest(t, http.MethodDelete, fmt.Sprintf("%s/books/%d", baseURL, created.ID), nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func ensureDockerAvailable(t *testing.T, ctx context.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("docker not available: %v", err)
	}
	defer cli.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if _, err := cli.Ping(pingCtx); err != nil {
		t.Skipf("docker not available: %v", err)
	}
}

func startPostgresContainer(ctx context.Context, t *testing.T) tc.Container {
	req := tc.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "bookshop",
			"POSTGRES_USER":     "bookshop",
			"POSTGRES_DB":       "bookshop",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	return container
}

func waitForServer(t *testing.T, url string) {
	client := &http.Client{Timeout: 1 * time.Second}
	deadline := time.Now().Add(30 * time.Second)

	for {
		if time.Now().After(deadline) {
			t.Fatalf("server not ready after timeout waiting for %s", url)
		}

		resp, err := client.Get(url)
		if err == nil && resp.StatusCode < 500 {
			_ = resp.Body.Close()
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func doRequest(t *testing.T, method, url string, body interface{}) *http.Response {
	var payload []byte
	var err error

	if body != nil {
		payload, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func decodeResponse(t *testing.T, resp *http.Response, target interface{}) {
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	require.NoError(t, decoder.Decode(target))
}

func freePort(t *testing.T) int {
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
