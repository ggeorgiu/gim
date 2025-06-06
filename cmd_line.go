package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

type cmdLine struct {
	screen tcell.Screen
	cursor *cursor
	bound  bounds
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
	switch ev.Key() {
	case tcell.KeyDEL:
		if len(c.cmd) == 0 {
			return
		}
		c.cmd = c.cmd[:len(c.cmd)-1]
		c.cursor.left()
		return
	}

	r := ev.Rune()
	c.cmd = c.cmd + string(r)
	c.cursor.right()
}

func (c *cmdLine) refresh(b bounds) {
	c.bound = b
}

func (c *cmdLine) draw() {
	line := fmt.Sprintf(">> :%s", c.cmd)
	for i, ch := range line {
		c.screen.SetContent(i, c.bound.y1, ch, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	}
}

func (c *cmdLine) recvCursor() {
	c.cursor.hold()
	c.cursor.x = c.bound.x1 + 4
	c.cursor.y = c.bound.y1
}

func (c *cmdLine) reset() {
	c.cmd = ""
}
