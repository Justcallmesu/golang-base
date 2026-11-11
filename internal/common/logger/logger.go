package logger

import (
	"log/slog"
	"os"
	"time"

	charmlog "github.com/charmbracelet/log"
)

func New(appEnv string) *slog.Logger {
	var handler slog.Handler

	if appEnv == "development" {
		handler = charmlog.NewWithOptions(os.Stderr, charmlog.Options{
			ReportTimestamp: true,
			TimeFormat:      time.StampMilli,
			Level:           charmlog.DebugLevel,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return slog.New(handler)
}
