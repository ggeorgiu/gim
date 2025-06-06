package main

import (
	"github.com/gdamore/tcell/v2"
)

type editor struct {
	screen  tcell.Screen
	cursor  *cursor
	bounds  bounds
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
