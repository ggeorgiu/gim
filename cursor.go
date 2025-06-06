package main

import "github.com/gdamore/tcell/v2"

type cursor struct {
	screen tcell.Screen
	x      int
	y      int
	prevX  int
	prevY  int
}

func (c *cursor) hold() {
	c.prevX = c.x
	c.prevY = c.y
}

func (c *cursor) rev() {
	c.x = c.prevX
	c.y = c.prevY
}

func (c *cursor) up() {
	c.y--
}

func (c *cursor) down() {
	c.y++
}

func (c *cursor) left() {
	c.x--
}

func (c *cursor) right() {
	c.x++
}

func (c *cursor) draw() {
	c.screen.ShowCursor(c.x, c.y)
	c.screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock, tcell.ColorBlue)
}
