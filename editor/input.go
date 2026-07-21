package editor

import (
	"log/slog"
	"strings"

	"github.com/brice-v/et/config"
	"github.com/brice-v/et/consts"
	"github.com/brice-v/et/keys"

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
	keyAsRune := k.Str()
	key := k.Key()
	e.chordInvalidSuffix = ""
	if e.handleAwaitingChord(key, keyAsRune, k) {
		return
	}
	if e.startChord(key, keyAsRune, k.Modifiers()) {
		return
	}

	if e.termOpen && e.term != nil {
		e.term.HandleEvent(k)
		return
	}
	switch key {
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
	case tcell.KeyTab, tcell.KeyBacktab:
		if e.promptMode == promptModeNormal {
			if e.expandTabs {
				e.handleInsertRune(strings.Repeat(" ", e.cfg.TabWidth))
			} else {
				e.handleInsertRune("\t")
			}
		}
	}
	e.updateViewport()
	e.clampCursor()
	if key == tcell.KeyLeft || key == tcell.KeyRight {
		e.syncStickyCol()
	}

	if key == tcell.KeyRune {
		e.handleInsertRune(keyAsRune)
	} else if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.ExitPrompt) {
		if e.promptLabel != nil {
			e.exitPrompt()
		}
	}
	if e.promptMode == promptModeFind && key != tcell.KeyRune {
		if keys.IsKeyAny(key, keyAsRune, k.Modifiers(), []config.Key{e.cfg.KeyBindings.FindNext, e.cfg.KeyBindings.FindPrevious}) {
			if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.FindPrevious) {
				e.findPreviousMatch()
			} else {
				e.findNextMatch()
			}
		}
	}
}

// handleAwaitingChord processes a key press while awaiting a chord suffix.
// It returns true if the key was consumed by chord handling.
func (e *Editor) handleAwaitingChord(key tcell.Key, keyAsRune string, k *tcell.EventKey) bool {
	if !e.awaitingChord {
		return false
	}
	e.awaitingChord = false
	if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.ExitPrompt) {
		return true
	}
	// Global chord actions
	if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.Quit.Suffix) {
		e.Exit = true
		return true
	}
	if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.Find.Suffix) {
		if e.promptMode == promptModeFind {
			e.Find.Mode = findModeExact
			e.Find.LastSearchTerm = ""
			e.updatePromptLabel(e.getPromptFindLabel())
		} else {
			e.promptMode = promptModeFind
			e.prompt(e.getPromptFindLabel())
		}
		return true
	}
	if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.ToggleLineEnding.Suffix) {
		e.buffer.ToggleLineEnding()
		return true
	}
	if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.ToggleExpandTabs.Suffix) {
		e.ToggleExpandTabs()
		return true
	}
	if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.ToggleTerminal.Suffix) {
		e.ToggleTerminal()
		return true
	}
	if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.TerminalIncreaseChord.Suffix) {
		e.IncreaseTerminalHeight()
		return true
	}
	if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.TerminalDecreaseChord.Suffix) {
		e.DecreaseTerminalHeight()
		return true
	}
	// Find mode specific suffixes
	if e.promptMode == promptModeFind {
		if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.FindSecondary1Chord.Suffix) {
			e.Find.Mode = findModeIgnoreCase
			e.Find.LastSearchTerm = ""
			e.updatePromptLabel(e.getPromptFindLabel())
			return true
		}
		if keys.IsKey(key, keyAsRune, k.Modifiers(), e.cfg.KeyBindings.FindSecondary2Chord.Suffix) {
			e.Find.Mode = findModeRegex
			e.Find.LastSearchTerm = ""
			e.updatePromptLabel(e.getPromptFindLabel())
			return true
		}
	}
	e.chordInvalidSuffix = "invalid suffix " + k.Name()
	return true
}

// startChord checks if the key matches the configured chord prefix and enters chord awaiting mode.
// It returns true if a chord was started.
func (e *Editor) startChord(key tcell.Key, keyAsRune string, mods tcell.ModMask) bool {
	if keys.IsKey(key, keyAsRune, mods, e.cfg.ChordPrefix) {
		e.awaitingChord = true
		return true
	}
	return false
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

func (e *Editor) updatePromptLabel(label string) {
	oldLen := len(e.promptLabel)
	e.promptLabel = []rune(" " + label + " ")
	e.cx += len(e.promptLabel) - oldLen
}

func (e *Editor) prompt(label string) {
	e.savedCx, e.savedCy = e.cx, e.cy
	e.savedVScrollOffset, e.savedHScrollOffset = e.vScrollOffset, e.hScrollOffset
	e.Find.FoundCx = -1
	e.sbh++
	e.updatePromptLabel(label)
	e.promptInput = []rune{}
	// The prompt is drawn at e.sh - 1 during drawEditorArea, where e.sh is
	// reduced by terminalHeight() + 1, so the absolute screen line is
	// e.sh - terminalHeight() - 2.
	e.cy = e.sh - e.terminalHeight() - 2
	if e.cy < 0 {
		e.cy = 0
	}
	e.cx = len(e.promptLabel)
}

func (e *Editor) exitPrompt() {
	e.sbh--
	e.promptLabel = nil
	e.promptInput = nil
	if e.Find.FoundCx >= 0 {
		e.cx, e.cy = e.Find.FoundCx, e.Find.FoundCy
	} else {
		e.cx, e.cy = e.savedCx, e.savedCy
	}
	e.promptMode = promptModeNormal
}
