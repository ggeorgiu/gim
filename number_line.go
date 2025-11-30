package main

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
)

type numberLine struct {
	screen tcell.Screen
	bounds columnBounds
	editor *editor
}

func newNumberLine(s tcell.Screen, e *editor) *numberLine {
	return &numberLine{
		screen: s,
		editor: e,
	}
}

func (nl *numberLine) refresh(b columnBounds) {
	nl.bounds = b
}

func (nl *numberLine) draw() {
	for y := range nl.bounds.y {
		val := fmt.Sprintf("%2d │", int(math.Abs(float64(nl.editor.cursorLineAt-y))))

		if y == nl.editor.cursorY() {
			val = fmt.Sprintf(" %2d│", nl.editor.lineIdx)
		}

		for x, r := range val {
			nl.screen.SetContent(nl.bounds.x+x, y, r, nil, tcell.StyleDefault)
		}
	}
}
