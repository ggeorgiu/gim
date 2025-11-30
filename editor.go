package main

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type slidingView struct {
	height int
	from   int
	to     int
}

func (sv *slidingView) down(maxLen int) {
	if sv.to == maxLen {
		return
	}

	sv.from++
	sv.to++
}
func (sv *slidingView) up() {
	if sv.from == 0 {
		return
	}

	sv.from--
	sv.to--
}

type editor struct {
	screen      tcell.Screen
	file        *os.File
	cursor      *cursor
	bounds      bounds
	mode        mode
	content     []string
	slidingView slidingView
}

func newEditor(screen tcell.Screen, c *cursor, file *os.File) (*editor, error) {
	var all []byte
	if file != nil {
		c, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		all = c
	}

	e := editor{
		screen:      screen,
		cursor:      c,
		content:     toSlice(all),
		file:        file,
		slidingView: slidingView{},
	}

	return &e, nil
}

func toSlice(all []byte) []string {
	val := string(all)
	return strings.Split(val, "\n")
}

func (e *editor) refresh(b bounds) {
	e.bounds = b
	e.slidingView.height = b.y2
	if e.slidingView.to == 0 {
		e.slidingView.to = e.slidingView.height
	} else {
		e.slidingView.to = e.slidingView.to - (e.slidingView.to - e.slidingView.height)
	}
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
		if e.cursor.y == e.slidingView.height-1 {
			e.slidingView.down(len(e.content) - 1)
			return
		}

		e.cursor.down()
		return
	case 'k':
		if e.cursor.y == e.bounds.y1 {
			e.slidingView.up()
			return
		}

		e.cursor.up()
		return
	}
}

func (e *editor) draw() {
	var idx int
	slog.Info("from and to", "from", e.slidingView.from, "to", e.slidingView.to)
	for i := e.slidingView.from; i < e.slidingView.to; i++ {
		line := e.content[i]
		for x, ch := range line {
			e.screen.SetContent(x+e.bounds.x1, idx+e.bounds.y1, ch, nil, tcell.StyleDefault)
		}
		idx++
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

func (e *editor) saveContent() error {
	var content []byte
	for i := 0; i < len(e.content)-1; i++ {
		content = append(content, []byte(e.content[i])...)
		content = append(content, []byte("\n")...)
	}
	content = append(content, []byte(e.content[len(e.content)-1])...)
	err := os.WriteFile(e.file.Name(), content, 0644)
	if err != nil {
		return err
	}

	return nil
}
