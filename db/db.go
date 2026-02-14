package db

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	pingTimeout = 5 * time.Second
	maxOpen     = 25
	maxIdle     = 25
	idleTime    = 5 * time.Minute
	maxLifetime = 30 * time.Minute
)

// Init initializes Database and sets the global DB variable.
func Init(ctx context.Context, connStr string) error {
	var err error

	slog.Info("Connecting to Database...")

	// Initialize connection using pq
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to open Database connection", slog.Any("error", err))
		return err
	}

	// Connection pool tuning
	DB.SetMaxOpenConns(maxOpen)
	DB.SetMaxIdleConns(maxIdle)
	DB.SetConnMaxIdleTime(idleTime)
	DB.SetConnMaxLifetime(maxLifetime)

	// Ping with timeout
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := DB.PingContext(pingCtx); err != nil {
		slog.Error("Database ping failed", slog.Any("error", err))
		return err
	}

	slog.Info("Database connected successfully")
	return nil
}

// Close closes the global DB instance
func Close() {
	if DB == nil {
		return
	}

	slog.Info("Closing Database connection pool...")

	if err := DB.Close(); err != nil {
		slog.Error("Failed to close DB connection", slog.Any("error", err))
	}
}
