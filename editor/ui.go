package editor

import (
	"fmt"
	"strings"

	"github.com/brice-v/et/consts"
	"github.com/brice-v/et/lexer"
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

func (e *Editor) updateLineHighlight(bufL, lineNumberOnScreen int, line []rune) {
	if e.cfg.DisableHighlighting || e.hl == nil {
		return
	}
	l := e.hl.newLexer(line)
	for tok := l.NextToken(); tok.Type != lexer.TTEof; tok = l.NextToken() {
		if tok.Type == lexer.TTIllegal {
			continue
		}
		hlStyle := e.hl.getStyle(tok.HlStyleType)
		if hlStyle == nil {
			continue
		}
		x := e.vx(bufL, tok.Position)
		for i, ch := range tok.Literal {
			sx := x + i
			if !e.inContent(sx) {
				break
			}
			e.s.SetContent(sx, lineNumberOnScreen, ch, nil, *hlStyle)
		}
	}
}

func (e *Editor) drawContent() {
	for bufL, last := e.vLines(); bufL <= last; bufL++ {
		screenLine := e.vy(bufL)
		line := e.buffer.Line(bufL)
		e.drawLine(screenLine, line)
		e.updateLineHighlight(bufL, screenLine, line)
	}
}

func (e *Editor) drawMatches() {
	if len(e.Find.Matches) == 0 || len(e.promptInput) == 0 {
		return
	}
	for idx, m := range e.Find.Matches {
		line := e.vy(m.line)
		if !e.inView(line) {
			continue
		}
		x := e.vx(m.line, m.col)
		bgColor := e.cfg.Colors.MatchHighlight.Color
		if idx == e.Find.CurrentMatchIdx {
			bgColor = e.cfg.Colors.CurrentMatchHighlight.Color
		}
		for i, ch := range e.promptInput {
			sx := x + i
			if !e.inContent(sx) {
				break
			}
			_, style, _ := e.s.Get(sx, line)
			e.s.SetContent(sx, line, ch, nil, style.Background(bgColor))
		}
	}
}

func (e *Editor) drawStatusBar() {
	statusStyle := e.baseStyle.Background(e.cfg.Colors.StatusBar.Color)
	statusBarH := e.vh()
	for x := range e.sw {
		e.s.SetContent(x, statusBarH, ' ', nil, statusStyle)
	}
	quitKeyBindsString := e.cfg.KeyBindings.Quit.String()
	fnameStr := e.fileName
	if e.fileName == "" {
		fnameStr = "<new file>"
	}
	modStr := ""
	if e.buffer.IsDirty() {
		modStr = " [*]"
	}
	statusMsg := fmt.Sprintf(" et - %s%s | %s to quit", fnameStr, modStr, quitKeyBindsString)
	if e.awaitingChord {
		statusMsg += " [awaiting chord...]"
	}
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
	leStr := strings.ToUpper(e.buffer.LineEnding())
	tsStr := "Tabs"
	if e.expandTabs {
		tsStr = "Spaces"
	}
	cursorLine := e.vScrollOffset + e.cy + 1
	cursorCol := e.stickyCol + 1
	posStr := fmt.Sprintf("%s [%s %s] Ln %d, Col %d", ft, leStr, tsStr, cursorLine, cursorCol)
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
	for y := range e.vh() {
		var ch []rune
		bufL := e.bufY(y)
		if useLineNums && bufL < numLines {
			ch = []rune(fmt.Sprintf("%*d ", e.lPad-1, bufL+1))
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

func (e *Editor) drawTerminal() {
	if !e.termOpen || e.term == nil {
		return
	}
	e.term.Draw()
	th := e.terminalHeight()
	sepY := e.sh - th - 1
	if sepY >= 0 {
		sepStyle := e.baseStyle.Background(e.cfg.Colors.StatusBar.Color)
		for x := range e.sw {
			e.s.SetContent(x, sepY, ' ', nil, sepStyle)
		}
		title := " TERMINAL (" + e.termShell + ")"
		for i, ch := range title {
			if i >= e.sw {
				break
			}
			e.s.SetContent(i, sepY, ch, nil, sepStyle)
		}
	}
}

func (e *Editor) drawTerminalCursor(th int) {
	if !e.termOpen || e.term == nil {
		e.s.ShowCursor(e.cx, e.cy)
		return
	}
	row, col, style, vis := e.term.Cursor()
	if vis {
		e.s.SetCursorStyle(style)
		e.s.ShowCursor(col, row+th)
	} else {
		e.s.HideCursor()
	}
}

func (e *Editor) drawEditorArea(th int) int {
	oldSh := e.sh
	e.sh -= th + 1
	if e.sh < 0 {
		e.sh = 0
	}
	e.updateViewport()
	e.drawLineNumbersOrTilde()
	e.clampCursor()
	e.drawContent()
	e.drawMatches()
	e.drawStatusBar()
	e.drawPrompt()
	e.drawWelcomeMessage()
	e.sh = oldSh
	return oldSh - th
}

func (e *Editor) Draw() {
	e.s.Clear()
	e.sw, e.sh = e.s.Size()
	th := e.terminalHeight()
	termOff := e.drawEditorArea(th)
	e.drawTerminal()
	e.drawTerminalCursor(termOff)
	e.s.Show()
}
