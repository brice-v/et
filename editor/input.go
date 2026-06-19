package editor

import (
	"et/consts"
	"et/keys"

	"github.com/gdamore/tcell/v3"
)

func fileToVisualCol(line []rune, fileCol int, tabWidth int) int {
	col := 0
	max := min(fileCol, len(line))
	for i := range max {
		if line[i] == '\t' {
			col += tabWidth
		} else {
			col++
		}
	}
	return col
}

func visualToFileCol(line []rune, visualCol int, tabWidth int) int {
	col := 0
	for fc := range line {
		if col > visualCol {
			return fc - 1
		}
		if line[fc] == '\t' {
			col += tabWidth
		} else {
			col++
		}
	}
	return len(line)
}

func (e *Editor) currentFileCol() int {
	fileLine := e.vScrollOffset + e.cy
	if fileLine < 0 || fileLine >= len(e.fileContentLines) {
		return 0
	}
	line := e.fileContentLines[fileLine]
	visualCol := fileToVisualCol(line, e.hScrollOffset, e.cfg.TabWidth) + (e.cx - e.lPad)
	return visualToFileCol(line, visualCol, e.cfg.TabWidth)
}

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
	fc := e.currentFileCol()
	fileLine := e.vScrollOffset + e.cy
	if fc <= 0 && fileLine > 0 {
		e.cy--
		e.stickyCol = consts.StickyColMax
	} else if fc > 0 {
		e.stickyCol = fc - 1
	}
}

func (e *Editor) handleMoveRight() {
	fc := e.currentFileCol()
	fileLine := e.vScrollOffset + e.cy
	if fileLine >= 0 && fileLine < len(e.fileContentLines) && fc >= len(e.fileContentLines[fileLine]) {
		e.cy++
		e.stickyCol = 0
	} else {
		e.stickyCol = fc + 1
	}
}

func (e *Editor) syncStickyCol() {
	fileLine := e.vScrollOffset + e.cy
	if fileLine >= 0 && fileLine < len(e.fileContentLines) {
		e.stickyCol = e.currentFileCol()
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
	line := e.fileContentLines[fileLine]
	lineLen := len(line)

	fc := max(min(e.stickyCol, lineLen), 0)

	// Convert file column to visual column, accounting for tab expansion
	vc := fileToVisualCol(line, fc, e.cfg.TabWidth)
	scrollVisual := fileToVisualCol(line, e.hScrollOffset, e.cfg.TabWidth)

	// Adjust horizontal scroll to keep cursor visible on screen
	textAreaWidth := e.sw - e.lPad
	if textAreaWidth > 0 {
		if vc >= scrollVisual+textAreaWidth {
			targetVisual := vc - (textAreaWidth - 1)
			e.hScrollOffset = visualToFileCol(line, targetVisual, e.cfg.TabWidth)
			scrollVisual = fileToVisualCol(line, e.hScrollOffset, e.cfg.TabWidth)
		}
	}
	if vc < scrollVisual {
		e.hScrollOffset = fc
		scrollVisual = fileToVisualCol(line, e.hScrollOffset, e.cfg.TabWidth)
	}
	if e.hScrollOffset < 0 {
		e.hScrollOffset = 0
		scrollVisual = fileToVisualCol(line, 0, e.cfg.TabWidth)
	}

	e.cx = vc - scrollVisual + e.lPad
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
