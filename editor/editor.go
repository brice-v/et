package editor

import (
	"et/config"
	"log/slog"
	"strings"

	"github.com/gdamore/tcell/v3"
)

type promptMode int

const (
	promptModeNormal promptMode = iota
	promptModeFind
)

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

	// Exit is a flag to trigger exit
	Exit bool
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

func (et *Editor) HandlePromptMode() {
	if et.promptMode == promptModeNormal {
		return
	}

	switch et.promptMode {
	case promptModeFind:
		input := string(et.promptInput)
		if input == "" {
			et.vScrollOffset = et.savedVScrollOffset
			et.hScrollOffset = et.savedHScrollOffset
			et.foundCx = -1
			return
		}
		et.findMatches(input)
	default:
		slog.Warn("unknown promptMode", "promptMode", et.promptMode)
	}
}

func (et *Editor) findMatches(input string) {
	for lineNo, line := range et.buffer.lines {
		lineText := string(line)
		// TODO: Update to support ignore case and regex
		n := strings.Index(lineText, input)
		if n == -1 {
			continue
		}
		et.displayFound(lineNo, n)
	}
	// TODO: Highlight matches on screen
}

func (et *Editor) displayFound(lineNo, col int) {
	vh := et.sh - et.sbh
	if vh <= 0 {
		return
	}

	savedCy, savedCx := et.cy, et.cx

	et.vScrollOffset = max(0, lineNo-vh/2)
	et.cy = lineNo - et.vScrollOffset
	et.stickyCol = col
	et.adjustViewport()
	et.clampCursorPos()

	et.foundCx, et.foundCy = et.cx, et.cy

	et.cy, et.cx = savedCy, savedCx
}
