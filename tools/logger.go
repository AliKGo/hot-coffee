package tools

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger() {
	file, err := os.OpenFile(*Dir+"/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		return
	}

	Logger = slog.New(slog.NewJSONHandler(file, nil))
	Logger.Info("Logger initialized")
}
