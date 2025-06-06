package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log/slog"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		slog.Error("app error", "err", err)
	}
}

func run(_ []string) error {
	logFile, err := os.OpenFile(".debug/debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	handler := slog.NewTextHandler(logFile, nil)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	screen, err := initScreen()
	if err != nil {
		return fmt.Errorf("failed to init scren, err: %w", err)
	}

	g := newGim(screen)
	g.Run()

	return nil
}

func initScreen() (tcell.Screen, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.SetStyle(tcell.StyleDefault)
	return screen, nil
}
