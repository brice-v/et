package editor

import (
	"et/config"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gdamore/tcell/v3"
)

type promptMode int

const (
	promptModeNormal promptMode = iota
	promptModeFind
)

func (pm promptMode) String() string {
	switch pm {
	case promptModeNormal:
		return "normal"
	case promptModeFind:
		return "find"
	default:
		return fmt.Sprintf("unknown promptMode: %d", pm)
	}
}

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

type Editor struct {
	s tcell.Screen
	// sw, sh screen width and height, calculated every Draw
	sw, sh int
	// lPad is the padding needed for line numbers or tilde on the left
	lPad int
	// sbh is the status bar height (always drawn at bottom but including
	//  this so it can be more than 1 high)
	// (defaults to 1)
	sbh       int
	baseStyle tcell.Style
	hl        *HighlightState

	// cx, cy cursor x and y position
	cx, cy int
	// savedCx, savedCy are the saved cursor positions when prompted
	savedCx, savedCy int
	// foundCx, foundCy are the last search match screen position (-1 = no active match)
	foundCx, foundCy int
	// vScrollOffset is the first visible line in the viewport
	vScrollOffset int
	// hScrollOffset is the first visible column in the viewport
	hScrollOffset      int
	savedVScrollOffset int
	savedHScrollOffset int
	// stickyCol is the file column for vertical movement that gets "stuck"
	stickyCol int

	cfg *config.Config

	fileName      string
	buffer        *Buffer
	fileExtension string

	// promptLabel is the message presented to the user at the bottom of the screen
	promptLabel []rune
	// promptInput is the users actual input into the prompt
	promptInput []rune
	promptMode  promptMode

	findMode  findMode
	hlMatches []matchPos

	// Exit is a flag to trigger exit
	Exit bool
}

type matchPos struct {
	line int
	col  int
}

func New(s tcell.Screen, cfg *config.Config, fileName string) *Editor {
	baseStyle := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	splitFilename := strings.Split(fileName, ".")
	fileExtension := ""
	if len(splitFilename) > 0 {
		fileExtension = splitFilename[len(splitFilename)-1]
	}
	s.SetCursorStyle(config.CursorStyleFromString(cfg.CursorStyle), cfg.CursorColor.Color)
	return &Editor{
		s:             s,
		sbh:           1,
		baseStyle:     baseStyle,
		cfg:           cfg,
		fileName:      fileName,
		buffer:        NewBuffer(fileName),
		fileExtension: fileExtension,
		hl:            NewHighlightState(cfg, fileExtension),
		promptMode:    promptModeNormal,
	}
}

func (e *Editor) HandlePromptMode() {
	if e.promptMode == promptModeNormal {
		return
	}

	switch e.promptMode {
	case promptModeFind:
		input := string(e.promptInput)
		if input == "" {
			e.vScrollOffset = e.savedVScrollOffset
			e.hScrollOffset = e.savedHScrollOffset
			e.foundCx = -1
			return
		}
		e.findMatches(input)
	default:
		slog.Warn("promptMode being used for HandlePromptMode not supported", "promptMode", e.promptMode.String())
	}
}

func (e *Editor) findMatches(input string) {
	if e.hlMatches == nil || len(e.hlMatches) != 0 {
		e.hlMatches = []matchPos{}
	}
	for lineNo, line := range e.buffer.lines {
		lineText := string(line)
		n := e.findIndex(lineText, input)
		if n == -1 {
			continue
		}
		e.displayFound(lineNo, n)
	}
	first, last := e.vLines()
	for i := first; i <= last; i++ {
		line := e.buffer.Line(i)
		lineText := string(line)
		n := e.findIndex(lineText, input)
		if n != -1 {
			e.hlMatches = append(e.hlMatches, matchPos{line: i, col: n})
		}
	}
}

func (e *Editor) findIndex(haystack, needle string) int {
	switch e.findMode {
	case findModeExact:
		return strings.Index(haystack, needle)
	case findModeIgnoreCase:
		return strings.Index(strings.ToLower(haystack), strings.ToLower(needle))
	case findModeRegex:
		// TODO: Support regex
		slog.Warn("regex find mode not yet supported")
		return -1
	default:
		slog.Warn("incorrect find mode being used for findIndex", "findMode", e.findMode.String())
		return -1
	}
}

func (e *Editor) displayFound(lineNo, col int) {
	vh := e.vh()

	savedCy, savedCx := e.cy, e.cx

	e.vScrollOffset = max(0, lineNo-vh/2)
	e.cy = lineNo - e.vScrollOffset
	e.stickyCol = col
	e.adjustViewport()
	e.clampCursorPos()

	e.foundCx, e.foundCy = e.cx, e.cy

	e.cy, e.cx = savedCy, savedCx
}

// Generic Helpers

func (e *Editor) vy(bufY int) int {
	return bufY - e.vScrollOffset
}

func (e *Editor) vx(bufLine, bufX int) int {
	if bufLine < 0 || bufLine >= e.buffer.NumLines() {
		return e.lPad
	}
	line := e.buffer.Line(bufLine)
	return e.lPad + e.visualCol(line, bufX) - e.visualCol(line, e.hScrollOffset)
}

func (e *Editor) vh() int {
	return e.sh - e.sbh
}

func (e *Editor) vLines() (first, last int) {
	first = e.vScrollOffset
	last = e.vScrollOffset + e.vh() - 1
	if n := e.buffer.NumLines(); last >= n {
		last = n - 1
	}
	return first, last
}

func (e *Editor) inContent(x int) bool {
	return x >= e.lPad && x < e.sw
}

func (e *Editor) inView(y int) bool {
	return y >= 0 && y < e.vh()
}

func (e *Editor) bufX(bufLine, vx int) int {
	if bufLine < 0 || bufLine >= e.buffer.NumLines() {
		return 0
	}
	line := e.buffer.Line(bufLine)
	visualCol := vx - e.lPad + e.visualCol(line, e.hScrollOffset)
	if visualCol < 0 {
		return 0
	}
	return e.fileCol(line, visualCol)
}

func (e *Editor) bufY(vy int) int {
	return e.vScrollOffset + vy
}

func (e *Editor) getPromptFindLabel() string {
	return fmt.Sprintf("Search [%s] ([%s,%s,%s] to change modes):",
		e.findMode.String(),
		e.cfg.KeyBindings.Find.String(),
		e.cfg.KeyBindings.FindSecondary1.String(),
		e.cfg.KeyBindings.FindSecondary2.String())
}
