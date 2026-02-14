package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"crud/db"
	"crud/routes"

	"github.com/joho/godotenv"
)

var router http.Handler

func init() {

	// logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// load env
	_ = godotenv.Load()

	connStr := os.Getenv("DATABASE_STRING")
	if connStr == "" {
		slog.Error("DATABASE_STRING not set")
		return
	}

	// init db
	err := db.Init(context.Background(), connStr)
	if err != nil {
		slog.Error("DB init failed", slog.Any("error", err))
		return
	}

	// init router
	router = routes.SetupRouter()

	slog.Info("Vercel server initialized")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
