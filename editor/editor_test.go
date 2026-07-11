package editor

import (
	"testing"

	"github.com/brice-v/et/config"
	"github.com/gdamore/tcell/v3"
)

type testScreen struct {
	w, h int
}

func (ts *testScreen) Init() error                                            { return nil }
func (ts *testScreen) Fini()                                                   {}
func (ts *testScreen) Clear()                                                  {}
func (ts *testScreen) Fill(rune, tcell.Style)                                  {}
func (ts *testScreen) SetContent(int, int, rune, []rune, tcell.Style)           {}
func (ts *testScreen) Get(int, int) (string, tcell.Style, int)                  { return "", tcell.StyleDefault, 1 }
func (ts *testScreen) Put(int, int, string, tcell.Style) (string, int)          { return "", 0 }
func (ts *testScreen) PutStr(int, int, string)                                   {}
func (ts *testScreen) PutStrStyled(int, int, string, tcell.Style)                {}
func (ts *testScreen) SetStyle(tcell.Style)                                      {}
func (ts *testScreen) ShowCursor(int, int)                                       {}
func (ts *testScreen) HideCursor()                                               {}
func (ts *testScreen) SetCursorStyle(tcell.CursorStyle, ...tcell.Color)          {}
func (ts *testScreen) Show()                                                     {}
func (ts *testScreen) Sync()                                                     {}
func (ts *testScreen) CharacterSet() string                                      { return "" }
func (ts *testScreen) RegisterRuneFallback(rune, string)                         {}
func (ts *testScreen) UnregisterRuneFallback(rune)                               {}
func (ts *testScreen) Resize(int, int, int, int)                                 {}
func (ts *testScreen) Suspend() error                                            { return nil }
func (ts *testScreen) Resume() error                                             { return nil }
func (ts *testScreen) Beep() error                                               { return nil }
func (ts *testScreen) SetSize(int, int)                                          {}
func (ts *testScreen) Colors() int                                               { return 256 }
func (ts *testScreen) EventQ() chan tcell.Event                                  { return nil }
func (ts *testScreen) EnableMouse(...tcell.MouseFlags)                           {}
func (ts *testScreen) DisableMouse()                                             {}
func (ts *testScreen) EnablePaste()                                              {}
func (ts *testScreen) DisablePaste()                                             {}
func (ts *testScreen) EnableFocus()                                              {}
func (ts *testScreen) DisableFocus()                                             {}
func (ts *testScreen) LockRegion(int, int, int, int, bool)                       {}
func (ts *testScreen) Tty() (tcell.Tty, bool)                                    { return nil, false }
func (ts *testScreen) SetTitle(string)                                           {}
func (ts *testScreen) SetClipboard([]byte)                                       {}
func (ts *testScreen) GetClipboard()                                             {}
func (ts *testScreen) HasClipboard() bool                                        { return false }
func (ts *testScreen) ShowNotification(string, string)                           {}
func (ts *testScreen) KeyboardProtocol() tcell.KeyProtocol                       { return 0 }
func (ts *testScreen) Terminal() (string, string)                                { return "", "" }

func (ts *testScreen) Size() (int, int) {
	return ts.w, ts.h
}

func newTestEditor(cfg *config.Config) *Editor {
	s := &testScreen{w: 80, h: 24}
	return New(s, cfg, "")
}

func TestTerminalHeightClosed(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	if got := e.terminalHeight(); got != 0 {
		t.Errorf("terminalHeight() when closed = %d, want 0", got)
	}
}

func TestTerminalHeightPercentage(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	e.termOpen = true
	e.sh = 40
	e.termHeight = 0
	if got := e.terminalHeight(); got != 10 {
		t.Errorf("terminalHeight() = %d, want 10", got)
	}
}

func TestTerminalHeightPercentageMin(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	e.termOpen = true
	e.sh = 4
	e.termHeight = 0
	if got := e.terminalHeight(); got != 3 {
		t.Errorf("terminalHeight() = %d, want 3", got)
	}
}

func TestTerminalHeightCustomPercentage(t *testing.T) {
	cfg := config.NewDefault()
	cfg.TerminalHeightPercentage = 0.5
	e := newTestEditor(cfg)
	e.termOpen = true
	e.sh = 40
	e.termHeight = 0
	if got := e.terminalHeight(); got != 20 {
		t.Errorf("terminalHeight() with 0.5 = %d, want 20", got)
	}
}

func TestTerminalHeightOverride(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	e.termOpen = true
	e.termHeight = 15
	e.sh = 40
	if got := e.terminalHeight(); got != 15 {
		t.Errorf("terminalHeight() with override = %d, want 15", got)
	}
}

func TestTerminalHeightOverrideMin(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	e.termOpen = true
	e.termHeight = 1
	e.sh = 40
	if got := e.terminalHeight(); got != 3 {
		t.Errorf("terminalHeight() with override 1 = %d, want 3 (min)", got)
	}
}

func TestIncreaseTerminalHeight(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	e.termOpen = true
	e.sh = 40
	e.IncreaseTerminalHeight()
	if e.termHeight == 0 {
		t.Errorf("termHeight = 0 after IncreaseTerminalHeight, expected > 0")
	}
	want := 10 + 1
	if e.termHeight != want {
		t.Errorf("termHeight after increase = %d, want %d", e.termHeight, want)
	}
}

func TestDecreaseTerminalHeight(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	e.termOpen = true
	e.sh = 40
	e.termHeight = 10
	e.DecreaseTerminalHeight()
	if e.termHeight != 9 {
		t.Errorf("termHeight after decrease = %d, want 9", e.termHeight)
	}
}

func TestDecreaseTerminalHeightMin(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	e.termOpen = true
	e.termHeight = 3
	e.DecreaseTerminalHeight()
	if e.termHeight != 3 {
		t.Errorf("termHeight after decrease from 3 = %d, want 3 (min)", e.termHeight)
	}
}

func TestIncreaseDecreaseClosed(t *testing.T) {
	e := newTestEditor(config.NewDefault())
	e.termOpen = false
	e.termHeight = 0
	e.IncreaseTerminalHeight()
	if e.termHeight != 0 {
		t.Errorf("termHeight = %d after increase when closed, want 0", e.termHeight)
	}
	e.DecreaseTerminalHeight()
	if e.termHeight != 0 {
		t.Errorf("termHeight = %d after decrease when closed, want 0", e.termHeight)
	}
}
