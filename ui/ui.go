package ui

import (
	"et/config"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v3"
)

func drawLine(s tcell.Screen, baseStyle tcell.Style, w, lineNumberOnScreen int, line string) {
	lineRunes := []rune(line)
	lineLen := len(lineRunes)
	for x := range w {
		ch := ' '
		if x < lineLen && lineLen != 0 {
			ch = lineRunes[x]
		}
		s.SetContent(x, lineNumberOnScreen, ch, nil, baseStyle)
	}
}

func drawContent(s tcell.Screen, baseStyle tcell.Style, w, h int, fileContent string) {
	lines := strings.Split(fileContent, "\n")
	numLines := len(lines)
	for i := range h {
		l := ""
		if i < numLines {
			l = lines[i]
		}
		drawLine(s, baseStyle, w, i, l)
	}
}

func drawStatusBar(s tcell.Screen, cfg *config.Config, baseStyle tcell.Style, w, h int, fileName string) {
	statusStyle := baseStyle.Background(cfg.Colors.StatusBar.Color)
	statusBarH := h - 1
	for x := range w {
		s.SetContent(x, statusBarH, ' ', nil, statusStyle)
	}
	quitKeyBindsString := cfg.GetQuitKeyBindingsAsStr()
	fnameStr := fileName
	if fileName == "" {
		fnameStr = "<new file>"
	}
	statusMsg := fmt.Sprintf(" et - %s | %s to quit", fnameStr, quitKeyBindsString)
	for i, ch := range statusMsg {
		if i >= w {
			break
		}
		s.SetContent(i, statusBarH, ch, nil, statusStyle)
	}
}

func Draw(s tcell.Screen, cfg *config.Config, fileName, fileContent string) {
	style := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	s.Clear()
	w, h := s.Size()
	drawContent(s, style, w, h, fileContent)
	drawStatusBar(s, cfg, style, w, h, fileName)
	s.Show()
}
