package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
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
	screen        tcell.Screen
	lines         []string
	cursorX       int
	cursorY       int
	prevX         int
	prevY         int
	mode          Mode
	cmd           string
	statusLine    string
	running       bool
	contentOffset int
}

func NewEditor() (*Editor, error) {
	screen, err := initScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to init scren, err: %w", err)
	}

	defaultMode := Normal
	e := Editor{
		running: true,
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
	defer e.screen.Fini()

	for e.running {
		e.Draw()
		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.screen.Sync()
		case *tcell.EventKey:
			e.HandleKey(ev)
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

	// set cmd line
	if e.mode == Command {
		_, h := e.screen.Size()
		cmdLine := fmt.Sprintf(">> :%s", e.cmd)
		for i, ch := range cmdLine {
			e.screen.SetContent(i, h-1, ch, nil, tcell.StyleDefault)
		}
	}

	// set content line
	w, h := e.screen.Size()
	e.statusLine = fmt.Sprintf("> %s <", e.mode.String())
	for i, ch := range e.statusLine {
		e.screen.SetContent(i, h-2, ch, nil, tcell.StyleDefault)
	}

	// set cursor position
	x := e.cursorX
	y := e.cursorY
	if e.mode == Command {
		x = e.prevX
		y = e.prevY
	}

	cpos := fmt.Sprintf("[ line: %d | col:  %d ]", y, x)
	for i, ch := range cpos {
		pos := w - len(cpos) + i
		e.screen.SetContent(pos, h-2, ch, nil, tcell.StyleDefault)
	}

	e.screen.ShowCursor(e.cursorX, e.cursorY)
	e.screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
	e.screen.Show()
}

func (e *Editor) HandleKey(ev *tcell.EventKey) {
	switch e.mode {
	case Normal:
		e.handleNormalMode(ev)
	case Insert:
		e.handleInsertMode(ev)
	case Command:
		e.handleCommandMode(ev)
	default:
		return
	}
}

func (e *Editor) handleCommandMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyESC:
		e.mode = Normal
		e.cursorX = e.prevX
		e.cursorY = e.prevY

		e.cmd = ""
		return
	case tcell.KeyEnter:
		e.mode = Normal
		e.execCmd()
		return
	case tcell.KeyDEL:
		if len(e.cmd) == 0 {
			return
		}
		e.cmd = e.cmd[:len(e.cmd)-1]
		e.cursorX--

		return
	default:
		r := ev.Rune()
		e.cmd = e.cmd + string(r)
		e.cursorX++
	}
}

func (e *Editor) execCmd() {
	switch e.cmd {
	case "q":
		e.running = false
		return
	default:
	}
}
func (e *Editor) handleNormalMode(ev *tcell.EventKey) {
	r := ev.Rune()
	switch r {
	case ':':
		e.prevX = e.cursorX
		e.prevY = e.cursorY

		_, h := e.screen.Size()
		e.cursorX = 4 + len(e.cmd)
		e.cursorY = h - 1

		e.mode = Command
		return
	case 'i':
		e.mode = Insert
		return
	case 'l':
		e.increaseX()
		return
	case 'h':
		e.decreaseX()
		return
	case 'j':
		e.increaseY()
		return
	case 'k':
		e.decreaseY()
		return
	}
}

func (e *Editor) handleInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyDEL:
		if e.cursorX == 0 && e.cursorY == 0 {
			return
		}
		if e.cursorX == 0 {
			e.lines = e.lines[:len(e.lines)]
			e.decreaseY()
			e.cursorX = len(e.lines[e.cursorY])
			return
		}

		line := e.lines[e.cursorY]
		e.lines[e.cursorY] = line[:e.cursorX-1] + line[e.cursorX:]
		e.decreaseX()
		return
	case tcell.KeyEnter:
		e.lines = append(e.lines, "")
		e.cursorX = 0
		e.increaseY()
		return
	case tcell.KeyESC:
		e.mode = Normal
		return
	}

	r := ev.Rune()
	line := e.lines[e.cursorY]
	e.lines[e.cursorY] = line[:e.cursorX] + string(r) + line[e.cursorX:]
	e.increaseX()
}

func (e *Editor) decreaseX() {
	if e.cursorX-1 < 0 {
		return
	}

	e.cursorX--
}

func (e *Editor) increaseX() {
	if e.cursorX+1 > len(e.lines[e.cursorY]) {
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
