package editor

import (
	"et/consts"
	"et/lexer"
	"fmt"
	"log/slog"
)

func (e *Editor) drawLine(lineNumberOnScreen int, line []rune) {
	lineLen := 0
	if line != nil {
		lineLen = len(line)
	}
	idx := e.hScrollOffset
	for x := e.lPad; x < e.sw; x++ {
		ch := ' '
		if idx >= 0 && idx < lineLen {
			ch = line[idx]
		}
		if ch == '\t' {
			for twOffset := range e.cfg.TabWidth {
				if x+twOffset < e.sw {
					e.s.SetContent(x+twOffset, lineNumberOnScreen, ' ', nil, e.baseStyle)
				}
			}
			x += e.cfg.TabWidth - 1
		} else {
			e.s.SetContent(x, lineNumberOnScreen, ch, nil, e.baseStyle)
		}
		idx++
	}
}

func (e *Editor) updateLineHighlight(lineNumberOnScreen int, line []rune) {
	if e.cfg.DisableHighlighting {
		return
	}
	offset := e.lPad - 1
	l := lexer.New(line, e.hldb, e.hlOperators, e.commentToken, e.stringTokens)
	for tok := l.NextToken(); tok.Type != lexer.TTEof; tok = l.NextToken() {
		switch tok.Type {
		case lexer.TTIdent, lexer.TTString, lexer.TTNumber, lexer.TTComment:
			hlStyle, ok := e.getHighlightStyle(tok.HlStyleType)
			if !ok {
				continue
			}
			if hlStyle == nil {
				slog.Warn("hlStyle is nil for token", "tok", tok)
				continue
			}
			twOffset := 0
			if l.TabCount != 0 {
				twOffset = l.TabCount*e.cfg.TabWidth - l.TabCount
			}
			for i, ch := range tok.Literal {
				e.s.SetContent(tok.Position+offset+twOffset+i, lineNumberOnScreen, ch, nil, *hlStyle)
			}
		}
	}
}

func (e *Editor) drawContent() {
	numLines := len(e.fileContentLines)
	lastLine := e.vScrollOffset + (e.sh - e.sbh) - 1
	for fileLine := e.vScrollOffset; fileLine <= lastLine && fileLine < numLines; fileLine++ {
		screenLine := fileLine - e.vScrollOffset
		line := e.fileContentLines[fileLine]
		e.drawLine(screenLine, line)
		e.updateLineHighlight(screenLine, line)
	}
}

func (e *Editor) drawStatusBar() {
	statusStyle := e.baseStyle.Background(e.cfg.Colors.StatusBar.Color)
	statusBarH := e.sh - e.sbh
	for x := range e.sw {
		e.s.SetContent(x, statusBarH, ' ', nil, statusStyle)
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

	ft := ""
	if e.fileContentLines != nil {
		if e.fileExtension != "" {
			if fileType, ok := e.cfg.FileExtensions[e.fileExtension]; ok {
				ft = "[" + fileType + "]"
			}
		} else {
			ft = "[unknown]"
		}
	}
	cursorLine := e.vScrollOffset + e.cy + 1
	cursorCol := e.stickyCol + 1
	posStr := fmt.Sprintf("%s Ln %d, Col %d", ft, cursorLine, cursorCol)
	posX := max(e.sw-len(posStr), 0)
	for i, ch := range posStr {
		x := posX + i
		if x >= e.sw {
			break
		}
		e.s.SetContent(x, statusBarH, ch, nil, statusStyle)
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
		e.lPad = len(ch) + 1
	}
	for y := range e.sh - e.sbh {
		if useLineNums {
			ch = []rune(fmt.Sprintf("%*d ", e.lPad-1, e.vScrollOffset+y+1))
		}
		for i := range ch {
			e.s.SetContent(i, y, ch[i], nil, e.baseStyle)
		}
	}
}

func (e *Editor) drawWelcomeMessage() {
	if e.fileContentLines != nil {
		return
	}
	y := (e.sh / 2) - 2
	for wi, message := range consts.WelcomeMessages {
		x := (e.sw / 2) - (len(message) / 2)
		for i, ch := range message {
			e.s.SetContent(x+i, y, ch, nil, e.baseStyle)
		}
		if wi == 1 {
			y += 2
		} else {
			y++
		}
	}
}

func (e *Editor) Draw() {
	e.s.Clear()
	e.sw, e.sh = e.s.Size()
	e.updateViewport()
	e.drawLineNumbersOrTilde()
	e.clampCursor()
	e.drawContent()
	e.drawStatusBar()
	e.drawWelcomeMessage()
	e.s.ShowCursor(e.cx, e.cy)
	e.s.Show()
}
