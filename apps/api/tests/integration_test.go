package tests

import (
    "context"
    "encoding/json"
    "net/http"
    "os/exec"
    "testing"
    "time"
)

func TestIntegration_ListBooks(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()

    cmd := exec.CommandContext(ctx, "go", "run", "./cmd/api")
    if err := cmd.Start(); err != nil {
        t.Fatalf("failed to start api: %v", err)
    }
    defer func() {
        _ = cmd.Process.Kill()
        _ = cmd.Wait()
    }()

    waitCtx, waitCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer waitCancel()
    if err := waitForHealthy(waitCtx); err != nil {
        t.Fatalf("api did not become healthy: %v", err)
    }

    resp, err := http.Get("http://localhost:8080/books")
    if err != nil {
        t.Fatalf("failed to call books endpoint: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp.StatusCode)
    }

    var payload map[string]any
    if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }
    if payload["data"] == nil {
        t.Fatalf("expected data key in response")
    }
}

func waitForHealthy(ctx context.Context) error {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            resp, err := http.Get("http://localhost:8080/healthz")
            if err != nil {
                continue
            }
            resp.Body.Close()
            if resp.StatusCode == http.StatusOK {
                return nil
            }
        }
    }
}
