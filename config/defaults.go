package config

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func DefaultColorBackground() Color {
	return Color{color.NewRGBColor(40, 42, 53)}
}

func DefaultColorForeground() Color {
	return Color{color.White}
}

func DefaultColorStatusBar() Color {
	return Color{color.DarkCyan}
}

func DefaultTabWidth() int {
	return 4
}

func DefaultLeftPadString() string {
	return "~"
}

func DefaultShowLineNumbers() bool {
	return true
}

func DefaultKeyBindingsQuit() []Key {
	return []Key{
		{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl},
		{Key: tcell.KeyQ, Modifiers: tcell.ModNone},
		{Key: tcell.KeyEscape, Modifiers: tcell.ModNone},
	}
}

func DefaultKeyBindingFind() Key {
	return Key{Key: tcell.KeyF, Modifiers: tcell.ModCtrl}
}

func DefaultLanguagesColorMap() map[string]ColorMap {
	return map[string]ColorMap{
		"go": {
			Keywords1:    []string{"break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var"},
			Color1:       Color{Color: color.NewRGBColor(255, 120, 197)},
			Keywords2:    []string{"uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128", "byte", "rune", "uint", "int", "uintptr", "error"},
			Color2:       Color{Color: color.NewRGBColor(161, 231, 250)},
			Keywords3:    []string{"print", "println", "make", "append", "len", "copy", "panic", "recover", "min", "max", "clear", "delete", "real", "imag", "new", "init"},
			Color3:       Color{Color: color.NewRGBColor(134, 247, 137)},
			StringTokens: []string{"`", `"`, "'"},
			ColorString:  Color{Color: color.NewRGBColor(243, 250, 154)},
			Operators:    "+-*/!|^&%=~{}[]:()",
			// Percentage formatting should also go here, but probably only if its in a string
			SpecialTokens: []string{"nil", "(", ")", "true", "false", "iota"},
			// Numbers use Special Color as well
			SpecialColor: Color{Color: color.NewRGBColor(183, 149, 243)},
			CommentToken: "//",
			CommentColor: Color{Color: color.NewRGBColor(101, 114, 160)},
		},
	}
}

func DefaultFileExtensions() map[string]string {
	return map[string]string{
		"go": "go",
	}
}

func DefaultDisableHighlighting() bool {
	return false
}
