package utils

import (
	"log/slog"
	"os"
	"strings"
)

func NewSlog(level string) {
	var l slog.Leveler

	switch strings.ToLower(level) {
	case "error":
		l = slog.LevelError
	case "warn":
		l = slog.LevelWarn
	case "info":
		l = slog.LevelInfo
	case "debug":
		l = slog.LevelDebug
	default:
		l = slog.LevelInfo
	}
	// h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l, AddSource: true})
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: l})
	slog.SetDefault(slog.New(h))
}
