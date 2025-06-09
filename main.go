package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func main() {
	if err := run(os.Args); err != nil {
		slog.Error("app error", "err", err)
	}
}

func run(args []string) error {
	logFile, err := os.OpenFile(".debug/debug.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}

	var content io.Reader
	content = strings.NewReader("")
	if len(args) == 2 {
		file, err := os.Open(args[1])
		if err != nil {
			return err
		}

		content = file
	}

	handler := slog.NewTextHandler(logFile, nil)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	screen, err := initScreen()
	if err != nil {
		return fmt.Errorf("failed to init scren, err: %w", err)
	}

	c := &cursor{screen: screen, x: editorColumStart}
	e, err := newEditor(screen, c, content)
	if err != nil {
		return err
	}
	nl := newNumberLine(screen, e)

	sl := newStatusLine(screen, e)
	cl := newCmdLine(screen, c)
	g := newGim(screen, c, e, nl, sl, cl)
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
