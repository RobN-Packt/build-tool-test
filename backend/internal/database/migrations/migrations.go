package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

// Apply executes the embedded SQL migrations against the provided database connection.
func Apply(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("nil database connection")
	}

	files := []string{
		"sql/001_create_books.sql",
		"sql/002_seed_books.sql",
	}

	for _, file := range files {
		content, err := sqlFiles.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}

		statement := strings.TrimSpace(string(content))
		if statement == "" {
			continue
		}

		if _, err := db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("execute migration %s: %w", file, err)
		}
	}

	return nil
}
