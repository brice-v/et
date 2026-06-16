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

func ColorStatusBar() tcell.Color {
	return color.DarkCyan
}
