package editor

import (
	"et/consts"
	"et/keys"
	"log/slog"

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
	case tcell.KeyEnter:
		e.handleEnter()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		e.handleBackspace()
	case tcell.KeyDelete:
		e.handleDelete()
	case tcell.KeyTab:
		e.handleInsertRune("\t")
	}
	e.updateViewport()
	e.clampCursor()
	if key == tcell.KeyLeft || key == tcell.KeyRight {
		e.syncStickyCol()
	}
	if keys.IsKeyAny(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.Quit) {
		e.Exit = true
	} else if key == tcell.KeyRune && !e.Exit {
		e.handleInsertRune(keyAsRune)
	} else if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.Find) {
		if e.promptLabel != nil {
			e.exitPrompt()
		} else {
			e.promptMode = promptModeFind
			e.prompt("Search:")
		}
	}
}

func (e *Editor) handleMoveUp() {
	if e.promptLabel != nil {
		// TODO: Maybe have action on up when prompt is filled in
		return
	}
	e.cy--
}

func (e *Editor) handleMoveDown() {
	if e.promptLabel != nil {
		// TODO: Maybe have action on down when prompt is filled in
		return
	}
	e.cy++
}

func (e *Editor) handleMoveLeft() {
	if e.promptLabel != nil {
		if e.cx > len(e.promptLabel) {
			e.cx--
		}
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
	if e.promptLabel != nil {
		maxCx := len(e.promptLabel) + len(e.promptInput)
		if e.cx < maxCx {
			e.cx++
		}
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
	if e.promptLabel != nil {
		// TODO: What should happen here
		return
	}
	fileLine := e.vScrollOffset + e.cy
	if fileLine >= 0 && fileLine < e.buffer.NumLines() {
		e.stickyCol = e.currentFileCol()
	}
}

func (e *Editor) clampCursorPos() {
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

func (e *Editor) clampCursor() {
	if e.promptLabel != nil {
		return
	}
	e.clampCursorPos()
}

func (e *Editor) adjustViewport() {
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

func (e *Editor) updateViewport() {
	if e.promptLabel != nil {
		return
	}
	e.adjustViewport()
}

func (e *Editor) handleInsertRune(r string) {
	if len(r) == 0 {
		return
	}
	if e.promptLabel != nil {
		ci := e.cx - len(e.promptLabel)
		runes := []rune(r)
		newInput := make([]rune, 0, len(e.promptInput)+len(runes))
		newInput = append(newInput, e.promptInput[:ci]...)
		newInput = append(newInput, runes...)
		newInput = append(newInput, e.promptInput[ci:]...)
		e.promptInput = newInput
		e.cx += len(runes)
		return
	}
	fileLine := e.vScrollOffset + e.cy
	fc := e.currentFileCol()
	for _, ch := range r {
		e.buffer.InsertRune(fileLine, fc, ch)
		fc++
	}
	e.stickyCol = fc
}

func (e *Editor) handleEnter() {
	if e.promptLabel != nil {
		// TODO: Maybe submit prompt on enter
		slog.Info("prompt entered", "prompt", string(e.promptInput))
		e.exitPrompt()
		return
	}
	fileLine := e.vScrollOffset + e.cy
	fc := e.currentFileCol()
	e.buffer.SplitLine(fileLine, fc)
	e.cy++
	e.stickyCol = 0
}

func (e *Editor) handleBackspace() {
	if e.promptLabel != nil {
		ci := e.cx - len(e.promptLabel)
		if ci > 0 {
			e.promptInput = append(e.promptInput[:ci-1], e.promptInput[ci:]...)
			e.cx--
		}
		return
	}
	fileLine := e.vScrollOffset + e.cy
	fc := e.currentFileCol()
	if fc > 0 {
		e.buffer.DeleteRune(fileLine, fc-1)
		e.stickyCol = fc - 1
	} else if fileLine > 0 {
		prevLineLen := len(e.buffer.Line(fileLine - 1))
		e.buffer.JoinLine(fileLine - 1)
		e.cy--
		e.stickyCol = prevLineLen
	}
}

func (e *Editor) handleDelete() {
	if e.promptLabel != nil {
		ci := e.cx - len(e.promptLabel)
		if ci < len(e.promptInput) {
			e.promptInput = append(e.promptInput[:ci], e.promptInput[ci+1:]...)
		}
		return
	}
	fileLine := e.vScrollOffset + e.cy
	fc := e.currentFileCol()
	if fileLine < e.buffer.NumLines() && fc < len(e.buffer.Line(fileLine)) {
		e.buffer.DeleteRune(fileLine, fc)
	} else if fileLine < e.buffer.NumLines()-1 {
		e.buffer.JoinLine(fileLine)
	}
}

func (e *Editor) prompt(label string) {
	e.savedCx, e.savedCy = e.cx, e.cy
	e.savedVScrollOffset, e.savedHScrollOffset = e.vScrollOffset, e.hScrollOffset
	e.foundCx = -1
	e.sbh++
	e.promptLabel = []rune(" " + label + " ")
	e.promptInput = []rune{}
	e.cy = e.sh - 1
	e.cx = len(e.promptLabel)
}

func (e *Editor) exitPrompt() {
	e.sbh--
	e.promptLabel = nil
	e.promptInput = nil
	if e.foundCx >= 0 {
		e.cx, e.cy = e.foundCx, e.foundCy
	} else {
		e.cx, e.cy = e.savedCx, e.savedCy
	}
	e.promptMode = promptModeNormal
}
