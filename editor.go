package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log/slog"
)

type Mode int

const (
	Normal Mode = iota
	Insert
	Command
)

func (m Mode) String() string {
	switch m {
	case Normal:
		return "NRM"
	case Insert:
		return "INS"
	case Command:
		return "CMD"
	default:
		return "???"
	}
}

type Editor struct {
	screen     tcell.Screen
	lines      []string
	cursorX    int
	cursorY    int
	mode       Mode
	statusLine string
}

func NewEditor() (*Editor, error) {
	screen, err := initScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to init scren, err: %w", err)
	}

	defaultMode := Normal
	e := Editor{
		screen:  screen,
		lines:   []string{""},
		cursorX: 0,
		cursorY: 0,
		mode:    defaultMode,
	}

	return &e, nil
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

func (e *Editor) Run() {
	for {
		e.Draw()
		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.screen.Sync()
		case *tcell.EventKey:
			if !e.HandleKey(ev) {
				return
			}
		}
	}
}

func (e *Editor) Draw() {
	e.screen.Clear()

	for y, line := range e.lines {
		for x, ch := range line {
			e.screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
		}
	}

	// set content line
	_, h := e.screen.Size()
	e.statusLine = ">" + e.mode.String()
	for i, ch := range e.statusLine {
		e.screen.SetContent(i, h-1, ch, nil, tcell.StyleDefault)
	}

	e.screen.ShowCursor(e.cursorX, e.cursorY)
	e.screen.Show()
}

func (e *Editor) HandleKey(ev *tcell.EventKey) bool {
	switch e.mode {
	case Normal:
		return e.handleNormalMode(ev)
	case Insert:
		return e.handleInsertMode(ev)
	case Command:
		return e.handleCommandMode(ev)
	default:
		return false
	}
}

func (e *Editor) handleCommandMode(ev *tcell.EventKey) bool {
	r := ev.Rune()
	switch r {
	case 'i':
		e.mode = Insert
	case 'n':
		e.mode = Normal
	}

	return true
}
func (e *Editor) handleNormalMode(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyESC:
		e.mode = Command
		return true
	}

	r := ev.Rune()
	switch r {
	case 'l':
		e.increaseX()
	case 'h':
		e.decreaseX()
	case 'j':
		e.increaseY()
	case 'k':
		e.decreaseY()
	}

	return true
}

func (e *Editor) handleInsertMode(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyDEL:
		if e.cursorX == 0 && e.cursorY == 0 {
			return true
		}
		if e.cursorX == 0 {
			e.lines = e.lines[:len(e.lines)]
		}

		line := e.lines[e.cursorY]
		e.lines[e.cursorY] = line[:e.cursorX-1] + line[e.cursorX:]
		e.decreaseX()
		return true
	case tcell.KeyEnter:
		e.lines = append(e.lines, "")
		e.cursorX = 0
		e.increaseY()
		return true
	case tcell.KeyESC:
		e.mode = Command
		return true
	}

	r := ev.Rune()
	line := e.lines[e.cursorY]
	e.lines[e.cursorY] = line[:e.cursorX] + string(r) + line[e.cursorX:]
	slog.Info("increase x", "x", e.cursorX)
	e.increaseX()
	slog.Info("increased x", "x", e.cursorX)

	return true
}

func (e *Editor) decreaseX() {
	if e.cursorX-1 < 0 {
		return
	}

	e.cursorX--
}

func (e *Editor) increaseX() {
	if e.cursorX+1 >= len(e.lines[e.cursorY]) {
		return
	}

	e.cursorX++
}

func (e *Editor) decreaseY() {
	if e.cursorY-1 < 0 {
		return
	}

	e.cursorY--
}

func (e *Editor) increaseY() {
	if e.cursorY+1 >= len(e.lines) {
		return
	}

	e.cursorY++
}
