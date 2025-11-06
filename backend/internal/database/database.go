package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"gofr.dev/pkg/gofr/config"
)

// Config holds database connection configuration.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// FromGoFr loads the configuration from a GoFr config provider.
func FromGoFr(cfg config.Config) Config {
	return Config{
		Host:     cfg.Get("DB_HOST"),
		Port:     cfg.GetOrDefault("DB_PORT", "5432"),
		User:     cfg.Get("DB_USER"),
		Password: cfg.Get("DB_PASSWORD"),
		Name:     cfg.Get("DB_NAME"),
		SSLMode:  cfg.GetOrDefault("DB_SSL_MODE", "disable"),
	}
}

// Connect initialises a PostgreSQL connection using the provided configuration.
func Connect(ctx context.Context, cfg Config) (*sql.DB, error) {
	if cfg.Host == "" || cfg.User == "" || cfg.Name == "" {
		return nil, fmt.Errorf("database configuration incomplete")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(15)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
