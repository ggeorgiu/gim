package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type cmdLine struct {
	screen tcell.Screen
	cursor *cursor
	bound  lineBounds
	cmd    string
}

func newCmdLine(s tcell.Screen, c *cursor) *cmdLine {
	cl := cmdLine{
		screen: s,
		cursor: c,
	}

	return &cl
}

func (c *cmdLine) handleKey(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyDEL {
		if len(c.cmd) == 0 {
			return
		}
		c.cmd = c.cmd[:len(c.cmd)-1]
		c.cursor.left()
		return
	}

	r := ev.Rune()
	c.cmd += string(r)
	c.cursor.right()
}

func (c *cmdLine) refresh(b lineBounds) {
	c.bound = b
}

func (c *cmdLine) draw() {
	line := fmt.Sprintf(">> :%s", c.cmd)
	for i, ch := range line {
		c.screen.SetContent(i, c.bound.y, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}
}

func (c *cmdLine) recvCursor() {
	c.cursor.hold()
	c.cursor.x = c.bound.x + 4
	c.cursor.y = c.bound.y
}

func (c *cmdLine) reset() {
	c.cmd = ""
}
