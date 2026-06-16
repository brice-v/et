package defaults

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func ColorBackground() tcell.Color {
	return color.Black
}

func ColorForeground() tcell.Color {
	return color.White
}

func ColorStatus() tcell.Color {
	return color.DarkCyan
}

func TabWidth() int {
	return 4
}

func Keybindings() map[string]string {
	return map[string]string{
		"CtrlQ":  "quit",
		"Escape": "quit",
	}
}
