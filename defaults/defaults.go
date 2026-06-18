package defaults

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type KeyBinding struct {
	Key       tcell.Key
	Modifiers tcell.ModMask
}

func ColorBackground() tcell.Color {
	return color.Black
}

func ColorForeground() tcell.Color {
	return color.White
}

func ColorStatusBar() tcell.Color {
	return color.DarkCyan
}

func KeyBindingsQuit() []KeyBinding {
	return []KeyBinding{
		{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl},
		{Key: tcell.KeyQ, Modifiers: tcell.ModNone},
		{Key: tcell.KeyEscape, Modifiers: tcell.ModNone},
	}
}

func TabWidth() int {
	return 4
}
