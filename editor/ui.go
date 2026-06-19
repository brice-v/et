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
	for x := e.lPad; x < e.sw; x++ {
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
	for i := range e.sh - e.sbh {
		var l []rune
		if i < numLines {
			l = e.fileContentLines[i]
		}
		e.drawLine(i, l)
	}
}

func (e *Editor) drawStatusBar() {
	statusStyle := e.baseStyle.Background(e.cfg.Colors.StatusBar.Color)
	statusBarH := e.sh - e.sbh
	for x := range e.sw {
		e.s.SetContent(x+e.lPad, statusBarH, ' ', nil, statusStyle)
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

func (e *Editor) drawLineNumbersOrTilde() {
	ch := []rune(e.cfg.LeftPadString)
	useLineNums := e.cfg.ShowLineNumbers && e.fileContentLines != nil
	if useLineNums {
		maxLinesDisplayed := len(e.fileContentLines)
		maxLinesAsStr := fmt.Sprintf("%d", maxLinesDisplayed)
		// +2 so that other things using this allow for extra padding to the right
		e.lPad = len([]rune(maxLinesAsStr)) + 2
	} else {
		e.lPad = len(ch)
	}
	for y := range e.sh - e.sbh {
		if useLineNums {
			// TODO: Eventually needs to be based off of fileContentLines
			// e.lPad-1 so that the space at the end is printed properly
			// y+1 to be 1 based indexed
			ch = []rune(fmt.Sprintf("%*d ", e.lPad-1, y+1))
		}
		for i := range ch {
			e.s.SetContent(i, y, ch[i], nil, e.baseStyle)
		}
	}
}

func (e *Editor) Draw() {
	e.s.Clear()
	e.sw, e.sh = e.s.Size()
	e.drawLineNumbersOrTilde()
	e.drawContent()
	e.drawStatusBar()
	e.s.Show()
}
