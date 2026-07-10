package editor

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/brice-v/et/config"
	"github.com/brice-v/et/terminal"

	"github.com/gdamore/tcell/v3"
)

type promptMode int

const (
	promptModeNormal promptMode = iota
	promptModeFind
)

func (pm promptMode) String() string {
	switch pm {
	case promptModeNormal:
		return "normal"
	case promptModeFind:
		return "find"
	default:
		return fmt.Sprintf("unknown promptMode: %d", pm)
	}
}

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
	Find             FindState
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
	absPath       string
	buffer        *Buffer
	fileExtension string

	// promptLabel is the message presented to the user at the bottom of the screen
	promptLabel []rune
	// promptInput is the users actual input into the prompt
	promptInput []rune
	promptMode  promptMode

	// Exit is a flag to trigger exit
	Exit bool

	// awaitingChord is true after the chord prefix key has been pressed
	awaitingChord bool

	// expandTabs controls whether tabs are converted to spaces on insert
	expandTabs bool

	// Terminal integration
	term        *terminal.VT
	termOpen    bool
	termStarted bool
	termShell   string
}

func New(s tcell.Screen, cfg *config.Config, fileName string) *Editor {
	baseStyle := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	splitFilename := strings.Split(fileName, ".")
	fileExtension := ""
	if len(splitFilename) > 0 {
		fileExtension = splitFilename[len(splitFilename)-1]
	}
	s.SetCursorStyle(config.CursorStyleFromString(cfg.CursorStyle), cfg.CursorColor.Color)

	absPath := ""
	if fileName != "" {
		abs, err := filepath.Abs(fileName)
		if err == nil {
			absPath = abs
		}
	} else {
		cwd, err := os.Getwd()
		if err == nil {
			absPath = cwd
		}
	}
	slog.Info("editor opened", "path", absPath)

	return &Editor{
		s:             s,
		sbh:           1,
		baseStyle:     baseStyle,
		cfg:           cfg,
		fileName:      fileName,
		absPath:       absPath,
		buffer:        NewBuffer(fileName),
		fileExtension: fileExtension,
		hl:            NewHighlightState(cfg, fileExtension),
		promptMode:    promptModeNormal,
		expandTabs:    cfg.ExpandTabs,
	}
}

func (e *Editor) HandlePromptMode() {
	if e.promptMode == promptModeNormal {
		return
	}

	switch e.promptMode {
	case promptModeFind:
		input := string(e.promptInput)
		if input == "" {
			e.vScrollOffset = e.savedVScrollOffset
			e.hScrollOffset = e.savedHScrollOffset
			e.Find.FoundCx = -1
			return
		}
		if input != e.Find.LastSearchTerm {
			e.Find.LastSearchTerm = input
			e.findMatches(input)
		}
	default:
		slog.Warn("promptMode being used for HandlePromptMode not supported", "promptMode", e.promptMode.String())
	}
}

// Generic Helpers

func (e *Editor) vy(bufY int) int {
	return bufY - e.vScrollOffset
}

func (e *Editor) vx(bufLine, bufX int) int {
	if bufLine < 0 || bufLine >= e.buffer.NumLines() {
		return e.lPad
	}
	line := e.buffer.Line(bufLine)
	return e.lPad + e.visualCol(line, bufX) - e.visualCol(line, e.hScrollOffset)
}

func (e *Editor) vh() int {
	return e.sh - e.sbh
}

func (e *Editor) vLines() (first, last int) {
	first = e.vScrollOffset
	last = e.vScrollOffset + e.vh() - 1
	if n := e.buffer.NumLines(); last >= n {
		last = n - 1
	}
	return first, last
}

func (e *Editor) inContent(x int) bool {
	return x >= e.lPad && x < e.sw
}

func (e *Editor) inView(y int) bool {
	return y >= 0 && y < e.vh()
}

func (e *Editor) bufX(bufLine, vx int) int {
	if bufLine < 0 || bufLine >= e.buffer.NumLines() {
		return 0
	}
	line := e.buffer.Line(bufLine)
	visualCol := vx - e.lPad + e.visualCol(line, e.hScrollOffset)
	if visualCol < 0 {
		return 0
	}
	return e.fileCol(line, visualCol)
}

func (e *Editor) bufY(vy int) int {
	return e.vScrollOffset + vy
}

// termSurface wraps a tcell.Screen to offset drawing to a specific region
type termSurface struct {
	screen tcell.Screen
	offY   int
	width  int
	height int
}

func (ts *termSurface) SetContent(x, y int, ch rune, comb []rune, style tcell.Style) {
	ts.screen.SetContent(x, y+ts.offY, ch, comb, style)
}

func (ts *termSurface) Size() (int, int) {
	return ts.width, ts.height
}

// shellCommand returns the default shell command for the current platform.
func shellCommand() string {
	if runtime.GOOS == "windows" {
		shell := os.Getenv("COMSPEC")
		if shell == "" {
			return "cmd.exe"
		}
		return shell
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "bash"
	}
	return shell
}

// ToggleExpandTabs switches between tab and space insertion.
func (e *Editor) ToggleExpandTabs() {
	e.expandTabs = !e.expandTabs
}

// terminalHeight returns the height of the terminal panel (minimum 3 rows)
func (e *Editor) terminalHeight() int {
	if !e.termOpen {
		return 0
	}
	return max(e.sh/4, 3)
}

// ToggleTerminal opens or closes the integrated terminal
func (e *Editor) ToggleTerminal() {
	if e.termOpen {
		e.termOpen = false
		return
	}
	if e.term == nil {
		e.term = terminal.New()
	}
	e.termOpen = true
	th := e.terminalHeight()
	vt := e.term
	e.updateTermSurface(th)
	vt.Resize(e.sw, th)
	vt.Attach(func(ev tcell.Event) {
		select {
		case e.s.EventQ() <- ev:
		default:
		}
	})
	if !e.termStarted {
		e.termStarted = true
		e.termShell = shellCommand()
		cmd := exec.Command(e.termShell)
		go func() {
			if err := vt.Start(cmd); err != nil {
				slog.Error("terminal start", "err", err)
			}
		}()
	}
}

func (e *Editor) updateTermSurface(th int) {
	if e.term == nil {
		return
	}
	e.term.SetSurface(&termSurface{
		screen: e.s,
		offY:   e.sh - th,
		width:  e.sw,
		height: th,
	})
}

func (e *Editor) ResizeTerminal() {
	if e.term == nil || !e.termOpen {
		return
	}
	th := e.terminalHeight()
	e.updateTermSurface(th)
	e.term.Resize(e.sw, th)
}
