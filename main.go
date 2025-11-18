package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crud/db"
	"crud/middleware"
	"crud/routes"

	"github.com/joho/godotenv"
)

func init() {
	// Initialize structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables", "error", err)
	}

	connStr := os.Getenv("DATABASE_STRING")
	if connStr == "" {
		slog.Error("Connection string not provided")
		os.Exit(1)
	}

	if err := db.Init(context.Background(), connStr); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	// db.Close() will be called during shutdown
	slog.Info("Application initialized successfully")
}

func main() {
	slog.Info("Starting Library Management Server...")

	router := routes.NewRouter()

	router.Use(middleware.RecoverMiddleware)
	router.Use(middleware.LoggingMiddleware)

	api := router.Group("/api/v1")

	apiRoutes := routes.NewRoutes()
	routes.AttachRoutes(api, apiRoutes)

	// Server config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Server running", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	} else {
		slog.Info("server stopped gracefully")
	}

	// close DB
	db.Close()
	slog.Info("Server shutdown complete")
}
