package main

import (
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application error", "err", err)
	}
}

func run() error {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	handler := slog.NewTextHandler(logFile, nil)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	e, err := NewEditor()
	if err != nil {
		return err
	}

	e.Run()
	return nil
}
