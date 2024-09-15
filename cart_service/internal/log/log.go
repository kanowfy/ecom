package log

import (
	"io"
	"log/slog"
)

func New(dest io.Writer, minLevel slog.Leveler, jsonFormat bool) *slog.Logger {
	var handler slog.Handler
	if jsonFormat {
		handler = slog.NewJSONHandler(dest, &slog.HandlerOptions{
			Level:     minLevel,
			AddSource: false,
		})
	} else {
		handler = slog.NewTextHandler(dest, &slog.HandlerOptions{
			Level:     minLevel,
			AddSource: false,
		})
	}

	return slog.New(handler)
}
