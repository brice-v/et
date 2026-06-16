package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v3"
)

func expandTabs(line string, tabWidth int) string {
	var b strings.Builder
	col := 0
	for _, ch := range line {
		if ch == '\t' {
			spaces := tabWidth - col%tabWidth
			b.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			b.WriteRune(ch)
			col++
		}
	}
	return b.String()
}

func drawScreen(screen tcell.Screen, filename string, hasFile bool, fileLines []string, style, statusStyle tcell.Style, tabWidth int) {
	screen.Clear()
	width, height := screen.Size()

	for y := range height {
		for x := range width {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}

	bodyHeight := height - 1

	if !hasFile {
		for y := range bodyHeight {
			screen.SetContent(0, y, '~', nil, style)
		}
		title := fmt.Sprintf("et - %s", version)
		titleRunes := []rune(title)
		tx := (width - len(titleRunes)) / 2
		ty := bodyHeight / 2
		for i, ch := range titleRunes {
			if tx+i < width {
				screen.SetContent(tx+i, ty, ch, nil, style)
			}
		}
	} else {
		lineNumWidth := 1
		for n := len(fileLines); n >= 10; n /= 10 {
			lineNumWidth++
		}
		if lineNumWidth < 2 {
			lineNumWidth = 2
		}

		maxY := min(len(fileLines), bodyHeight)
		for y := range maxY {
			numStr := fmt.Sprintf("%*d", lineNumWidth, y+1)
			for i, ch := range numStr {
				screen.SetContent(i, y, ch, nil, style)
			}
			screen.SetContent(lineNumWidth, y, '|', nil, style)
			line := expandTabs(fileLines[y], tabWidth)
			for i, ch := range line {
				screen.SetContent(lineNumWidth+1+i, y, ch, nil, style)
			}
		}
	}

	for x := range width {
		screen.SetContent(x, height-1, ' ', nil, statusStyle)
	}
	statusMsg := fmt.Sprintf(" et - %s", filename)
	if !hasFile {
		statusMsg = " et - <new file>"
	}
	for x, ch := range []rune(statusMsg) {
		if x >= width {
			break
		}
		screen.SetContent(x, height-1, ch, nil, statusStyle)
	}
	screen.Show()
}
