package editor

import (
	"et/keys"

	"github.com/gdamore/tcell/v3"
)

func (e *Editor) HandleKeyPress(k *tcell.EventKey) {
	keyAsRune := ""
	if k.Key() == tcell.KeyRune {
		keyAsRune = k.Str()
	} else if k.Key() == tcell.KeyUp {
		e.cy--
	} else if k.Key() == tcell.KeyDown {
		e.cy++
	} else if k.Key() == tcell.KeyLeft {
		e.cx--
	} else if k.Key() == tcell.KeyRight {
		e.cx++
	} else if k.Key() == tcell.KeyEnter && k.Modifiers() == tcell.ModNone {
		e.cy++
	} else if k.Key() == tcell.KeyEnter && k.Modifiers() == tcell.ModShift {
		e.cy--
	}
	e.updateViewport()
	e.clampCursor()
	if keys.IsKeyAny(k.Key(), keyAsRune, k.Modifiers(), e.cfg.KeyBindings.Quit) {
		e.Exit = true
	}
}

func (e *Editor) clampCursor() {
	numLines := len(e.fileContentLines)
	if numLines == 0 {
		e.cy = 0
		e.cx = e.lPad
		return
	}
	if e.cx < e.lPad {
		e.cx = e.lPad
	}
	fileLine := e.scrollOffset + e.cy
	if fileLine < 0 {
		fileLine = 0
	}
	if fileLine >= numLines {
		fileLine = numLines - 1
	}
	maxCx := e.lPad + len(e.fileContentLines[fileLine])
	if e.cx > maxCx {
		e.cx = maxCx
	}
}

func (e *Editor) updateViewport() {
	vh := e.sh - e.sbh
	if vh <= 0 {
		e.scrollOffset = 0
		e.cy = 0
		return
	}

	// Keep cy on screen
	if e.cy >= vh {
		e.scrollOffset += e.cy - vh + 1
		e.cy = vh - 1
	} else if e.cy < 0 {
		e.scrollOffset += e.cy
		e.cy = 0
	}
	if e.scrollOffset < 0 {
		e.scrollOffset = 0
	}

	// Keep cursor file line in bounds
	n := len(e.fileContentLines)
	if n > 0 {
		if fileLine := e.scrollOffset + e.cy; fileLine >= n {
			e.scrollOffset -= fileLine - n + 1
			if e.scrollOffset < 0 {
				e.cy += e.scrollOffset
				e.scrollOffset = 0
			}
		}
	} else {
		e.scrollOffset = 0
		e.cy = 0
	}
}
