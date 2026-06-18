package editor

import (
	"fmt"
)

func (e *Editor) drawLine(lineNumberOnScreen int, line []rune) {
	lineLen := 0
	if line != nil {
		lineLen = len(line)
	}
	runeIndex := 0
	for x := 1; x < e.sw; x++ {
		ch := ' '
		if runeIndex < lineLen {
			ch = line[runeIndex]
		}
		if ch == '\t' {
			for twOffset := range e.cfg.TabWidth {
				e.s.SetContent(x+twOffset, lineNumberOnScreen, ' ', nil, e.baseStyle)
			}
			x += e.cfg.TabWidth - 1
		} else {
			e.s.SetContent(x, lineNumberOnScreen, ch, nil, e.baseStyle)
		}
		runeIndex++
	}
}

func (e *Editor) drawContent() {
	numLines := len(e.fileContentLines)
	for i := range e.sh {
		var l []rune
		if i < numLines {
			l = e.fileContentLines[i]
		}
		e.drawLine(i, l)
	}
}

func (e *Editor) drawStatusBar() {
	statusStyle := e.baseStyle.Background(e.cfg.Colors.StatusBar.Color)
	statusBarH := e.sh - 1
	for x := range e.sw {
		e.s.SetContent(x+1, statusBarH, ' ', nil, statusStyle)
	}
	quitKeyBindsString := e.cfg.GetQuitKeyBindingsAsStr()
	fnameStr := e.fileName
	if e.fileName == "" {
		fnameStr = "<new file>"
	}
	statusMsg := fmt.Sprintf(" et - %s | %s to quit", fnameStr, quitKeyBindsString)
	for i, ch := range statusMsg {
		if i >= e.sw {
			break
		}
		e.s.SetContent(i, statusBarH, ch, nil, statusStyle)
	}
}

func (e *Editor) Draw() {
	e.s.Clear()
	e.sw, e.sh = e.s.Size()
	e.drawContent()
	e.drawStatusBar()
	e.s.Show()
}
