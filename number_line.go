package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"math"
)

type numberLine struct {
	screen tcell.Screen
	bounds bounds
	editor *editor
}

func newNumberLine(s tcell.Screen, e *editor) *numberLine {
	return &numberLine{
		screen: s,
		editor: e,
	}
}

func (nl *numberLine) draw() {
	for y := 0; y < len(nl.editor.content); y++ {
		val := fmt.Sprintf("%2d │", int(math.Abs(float64(nl.editor.cursorY()-y))))
		if y == nl.editor.cursorY() {
			val = fmt.Sprintf(" %2d│", y)
		}

		for x, r := range val {
			nl.screen.SetContent(nl.bounds.x1+x, y, r, nil, tcell.StyleDefault)
		}
	}
}

func (nl *numberLine) refresh(b bounds) {
	nl.bounds = b
}
