package migrations

import "embed"

// Files contains the SQL migration files.
//
//go:embed *.sql
var Files embed.FS
