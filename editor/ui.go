package editor

import (
	"et/consts"
	"et/lexer"
	"fmt"
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
	if e.cfg.DisableHighlighting || e.hl == nil {
		return
	}
	offset := e.lPad - 1
	l := e.hl.newLexer(line)
	for tok := l.NextToken(); tok.Type != lexer.TTEof; tok = l.NextToken() {
		if tok.Type == lexer.TTIllegal {
			continue
		}
		hlStyle := e.hl.getStyle(tok.HlStyleType)
		if hlStyle == nil {
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

func (e *Editor) drawContent() {
	numLines := e.buffer.NumLines()
	lastLine := e.vScrollOffset + (e.sh - e.sbh) - 1
	for fileLine := e.vScrollOffset; fileLine <= lastLine && fileLine < numLines; fileLine++ {
		screenLine := fileLine - e.vScrollOffset
		line := e.buffer.Line(fileLine)
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
	modStr := ""
	if e.buffer.IsDirty() {
		modStr = " [*]"
	}
	statusMsg := fmt.Sprintf(" et - %s%s | %s to quit", fnameStr, modStr, quitKeyBindsString)
	for i, ch := range statusMsg {
		if i >= e.sw {
			break
		}
		e.s.SetContent(i, statusBarH, ch, nil, statusStyle)
	}

	ft := ""
	if e.buffer.IsOpen() {
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

func (e *Editor) drawPrompt() {
	if e.promptLabel == nil {
		return
	}
	h := e.sh - 1
	statusStyle := e.baseStyle.Background(e.cfg.Colors.StatusBar.Color)
	for x := range e.sw {
		e.s.SetContent(x, h, ' ', nil, statusStyle)
	}
	for i, ch := range e.promptLabel {
		if i >= e.sw {
			break
		}
		e.s.SetContent(i, h, ch, nil, statusStyle)
	}
	for i, ch := range e.promptInput {
		pos := len(e.promptLabel) + i
		if pos >= e.sw {
			break
		}
		e.s.SetContent(pos, h, ch, nil, statusStyle)
	}
}

func (e *Editor) drawLineNumbersOrTilde() {
	tilde := []rune(e.cfg.LeftPadString)
	useLineNums := e.cfg.ShowLineNumbers && e.buffer.IsOpen()
	numLines := e.buffer.NumLines()
	if useLineNums {
		maxLinesAsStr := fmt.Sprintf("%d", numLines)
		e.lPad = len([]rune(maxLinesAsStr)) + 1
	} else {
		e.lPad = len(tilde) + 1
	}
	for y := range e.sh - e.sbh {
		var ch []rune
		fileLine := e.vScrollOffset + y
		if useLineNums && fileLine < numLines {
			ch = []rune(fmt.Sprintf("%*d ", e.lPad-1, fileLine+1))
		} else {
			ch = tilde
		}
		for i := range ch {
			e.s.SetContent(i, y, ch[i], nil, e.baseStyle)
		}
	}
}

func (e *Editor) drawWelcomeMessage() {
	if e.buffer.IsOpen() {
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
	e.drawPrompt()
	e.drawWelcomeMessage()
	e.s.ShowCursor(e.cx, e.cy)
	e.s.Show()
}
