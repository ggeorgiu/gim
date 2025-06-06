package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type statusLine struct {
	screen tcell.Screen
	bounds bounds
	cursor *cursor
	mode   mode
}

func newStatusLine(s tcell.Screen, c *cursor) *statusLine {
	l := statusLine{
		screen: s,
		cursor: c,
	}

	return &l
}

func (sl *statusLine) setMode(m mode) {
	sl.mode = m
}

func (sl *statusLine) refresh(b bounds) {
	sl.bounds = b
}

func (sl *statusLine) draw() {
	line := fmt.Sprintf("> %s <", sl.mode.String())
	for i, ch := range line {
		sl.screen.SetContent(i, sl.bounds.y1, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}

	var x = sl.cursor.x
	var y = sl.cursor.y
	if sl.mode == command {
		x = sl.cursor.prevX
		y = sl.cursor.prevY
	}

	cpos := fmt.Sprintf("[ line: %d | col:  %d ]", y, x)
	for i, ch := range cpos {
		pos := sl.bounds.x1 - len(cpos) + i
		sl.screen.SetContent(pos, sl.bounds.y1, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}
}
