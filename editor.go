package main

import (
	"github.com/gdamore/tcell/v2"
)

type editor struct {
	screen  tcell.Screen
	cursor  *cursor
	bounds  bounds
	mode    mode
	content []string
}

func newEditor(screen tcell.Screen, c *cursor) *editor {
	e := editor{
		screen:  screen,
		content: []string{""},
		cursor:  c,
	}

	return &e
}

func (e *editor) refresh(b bounds) {
	e.bounds = b
}

func (e *editor) handleKeyInInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyDEL:
		if e.cursor.x == e.bounds.x1 && e.cursor.y == 0 {
			return
		}
		if e.cursor.x == e.bounds.x1 {
			e.content = e.content[:len(e.content)-1]
			e.cursor.up()
			e.cursor.x = len(e.content[e.cursor.y]) + e.bounds.x1
			return
		}

		line := e.content[e.cursor.y]
		e.content[e.cursor.y] = line[:e.cursor.x-e.bounds.x1-1] + line[e.cursor.x-e.bounds.x1:]
		e.cursor.left()
		return
	case tcell.KeyEnter:
		e.content = append(e.content, "")
		e.cursor.down()
		e.cursor.x = e.bounds.x1
		return
	}

	r := ev.Rune()
	line := e.content[e.cursor.y]
	e.content[e.cursor.y] = line[:e.cursor.x-e.bounds.x1] + string(r) + line[e.cursor.x-e.bounds.x1:]
	e.cursor.right()
}

func (e *editor) handleKeyInNormalMode(ev *tcell.EventKey) {
	switch ev.Rune() {
	case 'l':
		if e.cursor.x == len(e.currentLine()) {
			return
		}

		e.cursor.right()
		return
	case 'h':
		if e.cursor.x == e.bounds.x1 {
			return
		}

		e.cursor.left()
		return
	case 'j':
		if e.cursor.y == len(e.content)-1 {
			return
		}

		e.cursor.down()
		return
	case 'k':
		if e.cursor.y == e.bounds.y1 {
			return
		}

		e.cursor.up()
		return
	}
}

func (e *editor) draw() {
	for y, line := range e.content {
		for x, ch := range line {
			e.screen.SetContent(x+e.bounds.x1, y+e.bounds.y1, ch, nil, tcell.StyleDefault)
		}
	}
}

func (e *editor) recvCursor() {
	e.cursor.rev()
}

func (e *editor) currentLine() string {
	return e.content[e.cursor.y]
}

func (e *editor) cursorX() int {
	if e.mode == command {
		return e.cursor.prevX - e.bounds.x1
	}

	return e.cursor.x - e.bounds.x1
}
func (e *editor) cursorY() int {
	if e.mode == command {
		return e.cursor.prevY
	}

	return e.cursor.y
}

func (e *editor) setMode(m mode) {
	e.mode = m
}
