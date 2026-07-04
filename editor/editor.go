package editor

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/brice-v/et/config"

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
	Find             FindState
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

	// Exit is a flag to trigger exit
	Exit bool

	// awaitingChord is true after the chord prefix key has been pressed
	awaitingChord bool
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
			e.Find.FoundCx = -1
			return
		}
		if input != e.Find.LastSearchTerm {
			e.Find.LastSearchTerm = input
			e.findMatches(input)
		}
	default:
		slog.Warn("promptMode being used for HandlePromptMode not supported", "promptMode", e.promptMode.String())
	}
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
