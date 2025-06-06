package main

import (
	"github.com/gdamore/tcell/v2"
)

type mode int

const (
	normal mode = iota
	insert
	command
)

func (m mode) String() string {
	switch m {
	case normal:
		return "NRM"
	case insert:
		return "INS"
	case command:
		return "CMD"
	default:
		return "???"
	}
}

type gim struct {
	running    bool
	mode       mode
	screen     tcell.Screen
	cursor     *cursor
	editor     *editor
	statusLine *statusLine
	cmdLine    *cmdLine
	numberLine *numberLine
}

func newGim(s tcell.Screen) *gim {
	c := &cursor{screen: s, x: 5}
	e := newEditor(s, c)
	nl := newNumberLine(s, e, c)

	g := gim{
		running:    true,
		screen:     s,
		cursor:     c,
		numberLine: nl,
		editor:     e,
		statusLine: newStatusLine(s, e),
		cmdLine:    newCmdLine(s, c),
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

	g.editor.refresh(bounds{5, 0, 0, h - 3})
	g.numberLine.refresh(bounds{0, 4, 0, h})
	g.statusLine.refresh(bounds{w, 0, h - 2, 0})
	g.cmdLine.refresh(bounds{0, 0, h - 1, 0})
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
