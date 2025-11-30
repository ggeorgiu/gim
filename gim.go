package main

import (
	"github.com/gdamore/tcell/v2"
)

const (
	numberLineWidth  = 4
	editorColumStart = 5
)

type gim struct {
	screen     tcell.Screen
	mode       mode
	cursor     *cursor
	editor     *editor
	statusLine *statusLine
	cmdLine    *cmdLine
	numberLine *numberLine
	running    bool
}

func newGim(
	s tcell.Screen,
	c *cursor,
	e *editor,
	nl *numberLine,
	sl *statusLine,
	cl *cmdLine,
) *gim {
	g := gim{
		running:    true,
		screen:     s,
		cursor:     c,
		numberLine: nl,
		editor:     e,
		statusLine: sl,
		cmdLine:    cl,
	}

	return &g
}

func (g *gim) Run() {
	defer g.screen.Fini()

	for g.running {
		g.Draw()
		ev := g.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			g.Refresh()
			g.screen.Sync()
		case *tcell.EventKey:
			g.HandleKey(ev)
		}
	}
}

func (g *gim) Refresh() {
	w, h := g.screen.Size()

	g.editor.refresh(bounds{editorColumStart, w, 0, h - 3})
	g.cmdLine.refresh(lineBounds{0, h - 1})
	g.statusLine.refresh(lineBounds{x: w, y: h - 2})
	g.numberLine.refresh(columnBounds{0, h, numberLineWidth})
}

func (g *gim) Draw() {
	g.screen.Clear()

	g.cursor.draw()
	g.editor.draw()
	g.statusLine.draw()
	g.numberLine.draw()

	if g.mode == command {
		g.cmdLine.draw()
	}

	g.screen.Show()
}

func (g *gim) HandleKey(ev *tcell.EventKey) {
	switch g.mode {
	case normal:
		g.handleNormalMode(ev)
	case insert:
		g.handleInsertMode(ev)
	case command:
		g.handleCommandMode(ev)
	default:
		return
	}
}

func (g *gim) handleCommandMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyESC:
		g.editor.recvCursor()
		g.cmdLine.reset()
		g.setMode(normal)

		return
	case tcell.KeyEnter:
		g.execCmd()

		g.editor.recvCursor()
		g.cmdLine.reset()
		g.setMode(normal)
		return
	default:
		g.cmdLine.handleKey(ev)
	}
}

func (g *gim) execCmd() {
	switch g.cmdLine.cmd {
	case "q":
		g.running = false
		return
	case "w":
		if err := g.saveContent(); err != nil {
			g.statusLine.error(err.Error())
		}
		return
	default:
	}
}

func (g *gim) handleNormalMode(ev *tcell.EventKey) {
	r := ev.Rune()
	switch r {
	case ':':
		g.cmdLine.recvCursor()
		g.setMode(command)
		return
	case 'i':
		g.setMode(insert)
		return
	default:
		g.editor.handleKeyInNormalMode(ev)
	}
}

func (g *gim) handleInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyESC:
		g.setMode(normal)
		return
	default:
		g.editor.handleKeyInInsertMode(ev)
	}
}

func (g *gim) setMode(m mode) {
	g.editor.setMode(m)
	g.mode = m
}

func (g *gim) saveContent() error {
	return g.editor.saveContent()
}
