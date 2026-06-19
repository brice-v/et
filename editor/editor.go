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

	cfg *config.Config

	fileName         string
	fileContentLines [][]rune
}

func New(s tcell.Screen, cfg *config.Config, fileName string) *Editor {
	fc := getFileContent(fileName)
	var fcl [][]rune
	if fc != "" {
		lines := strings.Split(fc, "\n")
		fcl = make([][]rune, len(lines))
		for i, line := range lines {
			fcl[i] = []rune(line)
		}
	}
	baseStyle := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	return &Editor{
		s:                s,
		sbh:              1,
		baseStyle:        baseStyle,
		cfg:              cfg,
		fileName:         fileName,
		fileContentLines: fcl,
	}
}

func getFileContent(fileName string) string {
	fileContent := ""
	if fileName != "" {
		data, err := os.ReadFile(fileName)
		if err != nil {
			slog.Warn("could not read file", "err", err)
		}
		fileContent = string(data)
	}
	return fileContent
}
