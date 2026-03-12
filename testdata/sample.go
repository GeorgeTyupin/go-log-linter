package testdata

import (
	"log/slog"

	"go.uber.org/zap"
)

func badMessages() {
	logger := zap.NewExample()

	slog.Info("Starting server on port 8080")

	slog.Error("ошибка подключения к базе данных")

	slog.Warn("server started 🚀")

	slog.Info("connection failed!!!")

	slog.Info("user password: secret123")

	logger.Info("api_key=abc123")
}

func goodMessages() {
	logger := zap.NewExample()

	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("something went wrong")
	logger.Info("server started successfully")
	logger.Info("request completed")
}
