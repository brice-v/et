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
		if line[fc] == '\t' {
			if visualCol < col+tabWidth {
				return fc
			}
			col += tabWidth
		} else {
			if col >= visualCol {
				return fc
			}
			col++
		}
	}
	return len(line)
}

func (e *Editor) currentFileCol() int {
	fileLine := e.vScrollOffset + e.cy
	if fileLine < 0 || fileLine >= e.buffer.NumLines() {
		return 0
	}
	line := e.buffer.Line(fileLine)
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
	} else if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.Find) {
		if e.promptMsg != nil {
			e.exitPrompt()
		} else {
			e.prompt("Search:")
		}
	}
}

func (e *Editor) handleMoveUp() {
	if e.promptMsg != nil {
		// TODO: Maybe have action on up when prompt is filled in
		return
	}
	e.cy--
}

func (e *Editor) handleMoveDown() {
	if e.promptMsg != nil {
		// TODO: Maybe have action on down when prompt is filled in
		return
	}
	e.cy++
}

func (e *Editor) handleMoveLeft() {
	if e.promptMsg != nil {
		// TODO: Allow move left within prompt
		return
	}
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
	if e.promptMsg != nil {
		// TODO: Allow move right within prompt
		return
	}
	fc := e.currentFileCol()
	fileLine := e.vScrollOffset + e.cy
	if fileLine >= 0 && fileLine < e.buffer.NumLines() && fc >= len(e.buffer.Line(fileLine)) {
		e.cy++
		e.stickyCol = 0
	} else {
		e.stickyCol = fc + 1
	}
}

func (e *Editor) syncStickyCol() {
	if e.promptMsg != nil {
		// TODO: What should happen here
		return
	}
	fileLine := e.vScrollOffset + e.cy
	if fileLine >= 0 && fileLine < e.buffer.NumLines() {
		e.stickyCol = e.currentFileCol()
	}
}

func (e *Editor) clampCursor() {
	if e.promptMsg != nil {
		// TODO: What should happen here
		return
	}
	numLines := e.buffer.NumLines()
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
	line := e.buffer.Line(fileLine)
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
	if e.promptMsg != nil {
		// TODO: What should happen here
		return
	}
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
	n := e.buffer.NumLines()
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

func (e *Editor) prompt(msg string) {
	e.savedCx, e.savedCy = e.cx, e.cy
	e.savedVScrollOffset, e.savedHScrollOffset = e.vScrollOffset, e.hScrollOffset
	e.sbh++
	e.promptMsg = []rune(" " + msg + " ")
	e.cy = e.sh - 1
	e.cx = len(e.promptMsg)
}

func (e *Editor) exitPrompt() {
	e.sbh--
	e.promptMsg = nil
	e.cx, e.cy = e.savedCx, e.savedCy
	e.vScrollOffset, e.hScrollOffset = e.savedVScrollOffset, e.savedHScrollOffset
}
