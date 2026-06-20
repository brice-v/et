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
	// #2c2e39
	return color.NewRGBColor(44, 46, 57)
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

type ColorMap struct {
	Keywords1    []string
	Color1       color.Color
	Keywords2    []string
	Color2       color.Color
	Keywords3    []string
	Color3       color.Color
	StringTokens []string
	ColorString  color.Color
}

func LanguagesColorMap() map[string]ColorMap {
	return map[string]ColorMap{
		"go": {
			Keywords1: []string{"break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var"},
			// #ff78c5
			Color1:    color.NewRGBColor(255, 120, 197),
			Keywords2: []string{"uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128", "byte", "rune", "uint", "int", "uintptr", "error"},
			// #87e1f6
			Color2:    color.NewRGBColor(135, 225, 246),
			Keywords3: []string{"print", "println", "make", "append", "len", "copy", "panic", "recover", "min", "max", "clear", "delete", "real", "imag", "new", "init"},
			// #4ad36d
			Color3:       color.NewRGBColor(74, 211, 109),
			StringTokens: []string{"`", `"`, "'"},
			// #d9e180
			ColorString: color.NewRGBColor(217, 225, 128),
		},
	}
}

func FileExtensions() map[string]string {
	return map[string]string{
		"go": "go",
	}
}
