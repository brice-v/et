package editor

import (
	"fmt"
	"log/slog"
	"strings"
)

type findMode int

const (
	findModeExact findMode = iota
	findModeIgnoreCase
	findModeRegex
)

func (fm findMode) String() string {
	switch fm {
	case findModeExact:
		return "exact"
	case findModeIgnoreCase:
		return "ignore case"
	case findModeRegex:
		return "regex"
	default:
		return fmt.Sprintf("unknown findMode: %d", fm)
	}
}

type matchPos struct {
	line int
	col  int
}

type FindState struct {
	Mode             findMode
	Matches          []matchPos
	FoundCx, FoundCy int
	CurrentMatchIdx  int
	LastSearchTerm   string
}

func (e *Editor) findMatches(input string) {
	e.Find.Matches = e.Find.Matches[:0]
	for lineNo, line := range e.buffer.lines {
		lineText := string(line)
		n := e.findIndex(lineText, input)
		if n == -1 {
			continue
		}
		e.Find.Matches = append(e.Find.Matches, matchPos{line: lineNo, col: n})
	}
	e.Find.CurrentMatchIdx = 0
	if len(e.Find.Matches) > 0 {
		e.displayCurrentMatch()
	} else {
		e.vScrollOffset = e.savedVScrollOffset
		e.hScrollOffset = e.savedHScrollOffset
		e.Find.FoundCx = -1
	}
}

func (e *Editor) findIndex(haystack, needle string) int {
	switch e.Find.Mode {
	case findModeExact:
		return strings.Index(haystack, needle)
	case findModeIgnoreCase:
		return strings.Index(strings.ToLower(haystack), strings.ToLower(needle))
	case findModeRegex:
		slog.Warn("regex find mode not yet supported")
		return -1
	default:
		slog.Warn("incorrect find mode being used for findIndex", "findMode", e.Find.Mode.String())
		return -1
	}
}

func (e *Editor) displayCurrentMatch() {
	if len(e.Find.Matches) == 0 {
		return
	}
	m := e.Find.Matches[e.Find.CurrentMatchIdx]
	vh := e.vh()
	savedCy, savedCx := e.cy, e.cx
	e.vScrollOffset = max(0, m.line-vh/2)
	e.cy = m.line - e.vScrollOffset
	e.stickyCol = m.col
	e.adjustViewport()
	e.clampCursorPos()
	e.Find.FoundCx, e.Find.FoundCy = e.cx, e.cy
	e.cy, e.cx = savedCy, savedCx
}

func (e *Editor) findNextMatch() {
	if len(e.Find.Matches) == 0 {
		return
	}
	e.Find.CurrentMatchIdx = (e.Find.CurrentMatchIdx + 1) % len(e.Find.Matches)
	e.displayCurrentMatch()
}

func (e *Editor) findPreviousMatch() {
	if len(e.Find.Matches) == 0 {
		return
	}
	e.Find.CurrentMatchIdx = (e.Find.CurrentMatchIdx - 1 + len(e.Find.Matches)) % len(e.Find.Matches)
	e.displayCurrentMatch()
}

func (e *Editor) getPromptFindLabel() string {
	return fmt.Sprintf("Search [%s] ([%s,%s,%s] modes, [%s,%s] matches):",
		e.Find.Mode.String(),
		e.cfg.KeyBindings.Find.String(),
		e.cfg.KeyBindings.FindSecondary1Chord.String(),
		e.cfg.KeyBindings.FindSecondary2Chord.String(),
		e.cfg.KeyBindings.FindNext.String(),
		e.cfg.KeyBindings.FindPrevious.String())
}
