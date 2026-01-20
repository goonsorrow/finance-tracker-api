package logger

import (
	"log/slog"
	"os"
	"time"
)

func InitSlogger() *slog.Logger {
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false, // добавляет информацию о файле и строке
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("15:04:05"))
				}
			}
			if a.Key == slog.LevelKey {
				if level, ok := a.Value.Any().(slog.Level); ok {
					switch level {
					case slog.LevelInfo:
						a.Value = slog.StringValue("INF")
					case slog.LevelDebug:
						a.Value = slog.StringValue("DBG")
					case slog.LevelWarn:
						a.Value = slog.StringValue("WRN")
					case slog.LevelError:
						a.Value = slog.StringValue("ERR")
					}
				}

			}
			return a
		},
	})

	slogger := slog.New(textHandler)
	slog.SetDefault(slogger)
	return slogger
}
