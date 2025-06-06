package main

import (
	"github.com/gdamore/tcell/v2"
)

type mode int

const (
	Normal mode = iota
	Insert
	Command
)

func (m mode) String() string {
	switch m {
	case Normal:
		return "NRM"
	case Insert:
		return "INS"
	case Command:
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
}

func newGim(s tcell.Screen) *gim {
	c := &cursor{screen: s}
	g := gim{
		running:    true,
		screen:     s,
		editor:     newEditor(s, c),
		statusLine: newStatusLine(s),
		cmdLine:    newCmdLine(s, c),
		cursor:     c,
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

	g.editor.refresh(bounds{0, 0, 0, h - 4})
	g.statusLine.refresh(bounds{w, 0, h - 3, 0})
	g.cmdLine.refresh(bounds{0, 0, h - 2, 0})
}

func (g *gim) Draw() {
	g.screen.Clear()

	g.cursor.draw()
	g.editor.draw()

	g.statusLine.setMode(g.mode)
	g.statusLine.draw()

	if g.mode == Command {
		g.cmdLine.draw()
	}

	g.screen.Show()
}

func (g *gim) HandleKey(ev *tcell.EventKey) {
	switch g.mode {
	case Normal:
		g.handleNormalMode(ev)
	case Insert:
		g.handleInsertMode(ev)
	case Command:
		g.handleCommandMode(ev)
	default:
		return
	}
}

func (g *gim) handleCommandMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyESC:
		g.editor.recvCursor()
		g.mode = Normal

		return
	case tcell.KeyEnter:
		g.editor.recvCursor()
		g.mode = Normal

		g.execCmd()
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
		g.mode = Command
		return
	case 'i':
		g.mode = Insert
		return
	case 'l':
		g.cursor.right()
		return
	case 'h':
		g.cursor.left()
		return
	case 'j':
		g.cursor.down()
		return
	case 'k':
		g.cursor.up()
		return
	}
}

func (g *gim) handleInsertMode(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyESC:
		g.mode = Normal
		return
	default:
		g.editor.handleKeyInInsertMode(ev)
	}
}
