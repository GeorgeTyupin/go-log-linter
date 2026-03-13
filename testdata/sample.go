package testdata

import (
	"log/slog"

	"go.uber.org/zap"
)

func badMessages() {
	logger := zap.NewExample()

	slog.Info("Starting server on port 8080")
	logger.Info("Starting server on port 8080")

	slog.Error("ошибка подключения к базе данных")
	logger.Error("ошибка подключения к базе данных")

	slog.Warn("server started 🚀")
	logger.Warn("server started 🚀")

	slog.Info("connection failed!!!")
	logger.Info("connection failed!!!")

	slog.Info("user password: secret123")
	logger.Info("user password: secret123")

	slog.Info("api_key=abc123")
	logger.Info("api_key=abc123")

	//Проверка кастомных паттернов
	slog.Info("mytoken=abc123")
	logger.Info("cvv=123")
}

func goodMessages() {
	logger := zap.NewExample()

	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("something went wrong")
	logger.Info("server started successfully")
	logger.Info("request completed")
}
