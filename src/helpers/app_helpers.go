package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

func ResolvePath(path string) (string, error) {
	path = os.ExpandEnv(path)

	if len(path) >= 2 && path[0] == '~' && (path[1] == '/' || path[1] == '\\') {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory - %w", err)
		}
		path = filepath.Join(home, path[2:])
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path %s - %w", path, err)
	}

	return filepath.Clean(abs), nil
}

func FormatDateString(raw string) string {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.UTC().Format("02 Jan 2006 15:04 UTC")
}

func InitLogger(level string) zerolog.Logger {
	zlevel, err := zerolog.ParseLevel(level)
	if err != nil {
		zlevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(zlevel)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
	}

	return zerolog.New(writer).With().Timestamp().Logger()
}
