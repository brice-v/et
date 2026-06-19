package editor

import (
	"et/consts"
	"et/keys"

	"github.com/gdamore/tcell/v3"
)

func (e *Editor) HandleKeyPress(k *tcell.EventKey) {
	keyAsRune := ""
	key := k.Key()
	switch key {
	case tcell.KeyRune:
		keyAsRune = k.Str()
	case tcell.KeyUp:
		e.handleMoveUp()
	case tcell.KeyDown:
		e.handleMoveDown()
	case tcell.KeyLeft:
		e.handleMoveLeft()
	case tcell.KeyRight:
		e.handleMoveRight()
	}
	e.updateViewport()
	e.clampCursor()
	if key == tcell.KeyLeft || key == tcell.KeyRight {
		e.syncStickyCol()
	}
	if keys.IsKeyAny(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.Quit) {
		e.Exit = true
	}
}

func (e *Editor) handleMoveUp() {
	e.cy--
}

func (e *Editor) handleMoveDown() {
	e.cy++
}

func (e *Editor) handleMoveLeft() {
	fc := (e.cx - e.lPad) + e.hScrollOffset
	fileLine := e.vScrollOffset + e.cy
	if fc <= 0 && fileLine > 0 {
		e.cy--
		e.stickyCol = consts.StickyColMax
	} else if fc > 0 {
		e.cx--
		e.stickyCol = (e.cx - e.lPad) + e.hScrollOffset
	}
}

func (e *Editor) handleMoveRight() {
	fc := (e.cx - e.lPad) + e.hScrollOffset
	fileLine := e.vScrollOffset + e.cy
	if fileLine >= 0 && fileLine < len(e.fileContentLines) && fc >= len(e.fileContentLines[fileLine]) {
		e.cy++
		e.stickyCol = 0
	} else {
		e.cx++
		e.stickyCol = (e.cx - e.lPad) + e.hScrollOffset
	}
}

func (e *Editor) syncStickyCol() {
	fileLine := e.vScrollOffset + e.cy
	if fileLine >= 0 && fileLine < len(e.fileContentLines) {
		e.stickyCol = (e.cx - e.lPad) + e.hScrollOffset
	}
}

func (e *Editor) clampCursor() {
	numLines := len(e.fileContentLines)
	if numLines == 0 {
		e.cy = 0
		e.cx = e.lPad
		e.hScrollOffset = 0
		return
	}
	if e.cx < e.lPad {
		e.cx = e.lPad
	}
	fileLine := max(e.vScrollOffset+e.cy, 0)
	if fileLine >= numLines {
		fileLine = numLines - 1
	}
	lineLen := len(e.fileContentLines[fileLine])

	fc := max(min(e.stickyCol, lineLen), 0)

	// Adjust horizontal scroll to keep cursor visible on screen
	textAreaWidth := e.sw - e.lPad
	if textAreaWidth > 0 {
		if fc >= textAreaWidth {
			e.hScrollOffset = fc - (textAreaWidth - 1)
		}
	}
	if fc < e.hScrollOffset {
		e.hScrollOffset = fc
	}
	if e.hScrollOffset < 0 {
		e.hScrollOffset = 0
	}

	e.cx = fc - e.hScrollOffset + e.lPad
}

func (e *Editor) updateViewport() {
	vh := e.sh - e.sbh
	if vh <= 0 {
		e.vScrollOffset = 0
		e.cy = 0
		return
	}

	// Keep cy on screen
	if e.cy >= vh {
		e.vScrollOffset += e.cy - vh + 1
		e.cy = vh - 1
	} else if e.cy < 0 {
		e.vScrollOffset += e.cy
		e.cy = 0
	}
	if e.vScrollOffset < 0 {
		e.vScrollOffset = 0
	}

	// Keep cursor file line in bounds
	n := len(e.fileContentLines)
	if n > 0 {
		if fileLine := e.vScrollOffset + e.cy; fileLine >= n {
			e.vScrollOffset -= fileLine - n + 1
			if e.vScrollOffset < 0 {
				e.cy += e.vScrollOffset
				e.vScrollOffset = 0
			}
		}
	} else {
		e.vScrollOffset = 0
		e.cy = 0
	}
}
