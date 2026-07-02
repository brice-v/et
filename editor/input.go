package editor

import (
	"et/consts"
	"et/keys"
	"log/slog"

	"github.com/gdamore/tcell/v3"
)

func (e *Editor) visualCol(line []rune, fileCol int) int {
	col := 0
	max := min(fileCol, len(line))
	for i := range max {
		if line[i] == '\t' {
			col += e.cfg.TabWidth
		} else {
			col++
		}
	}
	return col
}

func (e *Editor) fileCol(line []rune, visualCol int) int {
	col := 0
	for fc := range line {
		if line[fc] == '\t' {
			if visualCol < col+e.cfg.TabWidth {
				return fc
			}
			col += e.cfg.TabWidth
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
	return e.bufX(e.bufY(e.cy), e.cx)
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
	bufL := e.vScrollOffset + e.cy
	if fc <= 0 && bufL > 0 {
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
	bufL := e.vScrollOffset + e.cy
	if bufL >= 0 && bufL < e.buffer.NumLines() && fc >= len(e.buffer.Line(bufL)) {
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
	bufL := e.vScrollOffset + e.cy
	if bufL >= 0 && bufL < e.buffer.NumLines() {
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
	bufL := max(e.vScrollOffset+e.cy, 0)
	if bufL >= numLines {
		bufL = numLines - 1
	}
	line := e.buffer.Line(bufL)
	fc := max(min(e.stickyCol, len(line)), 0)
	vc := e.visualCol(line, fc)

	if textAreaWidth := e.sw - e.lPad; textAreaWidth > 0 {
		scrollVisual := e.visualCol(line, e.hScrollOffset)
		if vc < scrollVisual {
			scrollVisual = vc
		} else if vc >= scrollVisual+textAreaWidth {
			scrollVisual = vc - (textAreaWidth - 1)
		}
		e.hScrollOffset = e.fileCol(line, max(scrollVisual, 0))
	} else {
		e.hScrollOffset = 0
	}

	e.cx = e.vx(bufL, fc)
}

func (e *Editor) clampCursor() {
	if e.promptLabel != nil {
		return
	}
	e.clampCursorPos()
}

func (e *Editor) adjustViewport() {
	vh := e.vh()
	n := e.buffer.NumLines()
	if vh <= 0 || n == 0 {
		e.vScrollOffset = 0
		e.cy = 0
		return
	}

	// Desired file line, clamped to buffer
	fl := max(0, min(e.vScrollOffset+e.cy, n-1))
	e.cy = max(0, min(e.cy, vh-1))
	e.vScrollOffset = fl - e.cy

	if e.vScrollOffset < 0 {
		e.cy += e.vScrollOffset
		e.vScrollOffset = 0
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
	bufL := e.vScrollOffset + e.cy
	fc := e.currentFileCol()
	for _, ch := range r {
		e.buffer.InsertRune(bufL, fc, ch)
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
	bufL := e.vScrollOffset + e.cy
	fc := e.currentFileCol()
	e.buffer.SplitLine(bufL, fc)
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
	bufL := e.vScrollOffset + e.cy
	fc := e.currentFileCol()
	if fc > 0 {
		e.buffer.DeleteRune(bufL, fc-1)
		e.stickyCol = fc - 1
	} else if bufL > 0 {
		prevLineLen := len(e.buffer.Line(bufL - 1))
		e.buffer.JoinLine(bufL - 1)
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
	bufL := e.vScrollOffset + e.cy
	fc := e.currentFileCol()
	if bufL < e.buffer.NumLines() && fc < len(e.buffer.Line(bufL)) {
		e.buffer.DeleteRune(bufL, fc)
	} else if bufL < e.buffer.NumLines()-1 {
		e.buffer.JoinLine(bufL)
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
