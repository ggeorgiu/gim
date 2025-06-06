package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type statusLine struct {
	screen tcell.Screen
	bounds bounds
	editor *editor
}

func newStatusLine(s tcell.Screen, e *editor) *statusLine {
	l := statusLine{
		screen: s,
		editor: e,
	}

	return &l
}

func (sl *statusLine) refresh(b bounds) {
	sl.bounds = b
}

func (sl *statusLine) draw() {
	line := fmt.Sprintf("> %s <", sl.editor.mode.String())
	for i, ch := range line {
		sl.screen.SetContent(i, sl.bounds.y1, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}

	cpos := fmt.Sprintf("[ line: %2d | col: %2d ]", sl.editor.cursorY(), sl.editor.cursorX())
	for i, ch := range cpos {
		pos := sl.bounds.x1 - len(cpos) + i
		sl.screen.SetContent(pos, sl.bounds.y1, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}
}
