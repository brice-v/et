package editor

import (
	"et/config"
	"log/slog"
	"os"
	"strings"

	"github.com/gdamore/tcell/v3"
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
	// vScrollOffset is the first visible line in the viewport
	vScrollOffset int
	// hScrollOffset is the first visible column in the viewport
	hScrollOffset      int
	savedVScrollOffset int
	savedHScrollOffset int
	// stickyCol is the file column for vertical movement that gets "stuck"
	stickyCol int

	cfg *config.Config

	fileName         string
	fileContentLines [][]rune
	fileExtension    string

	// promptMsg is the message presented to the user at the bottom of the screen
	promptMsg []rune

	// Exit is a flag to trigger exit
	Exit bool
}

func New(s tcell.Screen, cfg *config.Config, fileName string) *Editor {
	fcl := getFileContent(fileName)
	baseStyle := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	splitFilename := strings.Split(fileName, ".")
	fileExtension := ""
	if len(splitFilename) > 0 {
		fileExtension = splitFilename[len(splitFilename)-1]
	}
	return &Editor{
		s:                s,
		sbh:              1,
		baseStyle:        baseStyle,
		cfg:              cfg,
		fileName:         fileName,
		fileContentLines: fcl,
		fileExtension:    fileExtension,
		hl:               NewHighlightState(cfg, fileExtension),
	}
}

func getFileContent(fileName string) [][]rune {
	var fcl [][]rune = nil
	if fileName != "" {
		data, err := os.ReadFile(fileName)
		if err != nil {
			slog.Warn("could not read file", "err", err)
			return fcl
		}
		lines := strings.Split(string(data), "\n")
		fcl = make([][]rune, len(lines))
		for i, line := range lines {
			fcl[i] = []rune(line)
		}
	}
	return fcl
}
