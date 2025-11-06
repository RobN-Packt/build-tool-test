package main

import (
	"context"

	"gofr.dev/pkg/gofr"

	"github.com/cursor/bookshop/backend/internal/books"
	"github.com/cursor/bookshop/backend/internal/database"
	"github.com/cursor/bookshop/backend/internal/database/migrations"
)

func main() {
	app := gofr.New()

	ctx := context.Background()

	dbCfg := database.FromGoFr(app.Config)

	db, err := database.Connect(ctx, dbCfg)
	if err != nil {
		app.Logger().Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := migrations.Apply(ctx, db); err != nil {
		app.Logger().Fatalf("failed to apply migrations: %v", err)
	}

	repo := books.NewRepository(db)
	service := books.NewService(repo)
	books.RegisterRoutes(app, service)

	app.Run()
}
