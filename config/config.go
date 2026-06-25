package config

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Key struct {
	tcell.Key
	Modifiers tcell.ModMask
}

type Color struct {
	color.Color
}

type ColorMap struct {
	Keywords1     []string `json:"keywords1"`
	Color1        Color    `json:"color1"`
	Keywords2     []string `json:"keywords2"`
	Color2        Color    `json:"color2"`
	Keywords3     []string `json:"keywords3"`
	Color3        Color    `json:"color3"`
	StringTokens  []string `json:"string_tokens"`
	ColorString   Color    `json:"color_string"`
	Operators     string   `json:"operators"`
	SpecialTokens []string `json:"special_tokens"`
	SpecialColor  Color    `json:"special_color"`
	CommentToken  string   `json:"comment_token"`
	CommentColor  Color    `json:"comment_color"`
}

type Colors struct {
	Foreground Color               `json:"foreground"`
	Background Color               `json:"background"`
	StatusBar  Color               `json:"status_bar"`
	Languages  map[string]ColorMap `json:"languages"`
}

type KeyBindings struct {
	Quit []Key `json:"quit"`
	Find Key   `json:"find"`
}

func (c *Config) GetQuitKeyBindingsAsStr() string {
	var s strings.Builder
	s.WriteByte('[')
	for i, k := range c.KeyBindings.Quit {
		s.WriteString(k.String())
		if i != len(c.KeyBindings.Quit)-1 {
			s.WriteByte(',')
		}
	}
	s.WriteByte(']')
	return s.String()
}

type Config struct {
	Colors              Colors            `json:"colors"`
	KeyBindings         KeyBindings       `json:"keybindings"`
	TabWidth            int               `json:"tab_width"`
	LeftPadString       string            `json:"left_pad_string"`
	ShowLineNumbers     bool              `json:"show_line_numbers"`
	FileExtensions      map[string]string `json:"file_extensions"`
	DisableHighlighting bool              `json:"disable_highlighting"`
	CursorStyle         string            `json:"cursor_style"`
	CursorColor         Color             `json:"cursor_color"`
}

func NewDefault() *Config {
	return &Config{
		Colors: Colors{
			Foreground: DefaultColorForeground(),
			Background: DefaultColorBackground(),
			StatusBar:  DefaultColorStatusBar(),
			Languages:  DefaultLanguagesColorMap(),
		},
		KeyBindings: KeyBindings{
			Quit: DefaultKeyBindingsQuit(),
			Find: DefaultKeyBindingFind(),
		},
		TabWidth:            DefaultTabWidth(),
		LeftPadString:       DefaultLeftPadString(),
		ShowLineNumbers:     DefaultShowLineNumbers(),
		FileExtensions:      DefaultFileExtensions(),
		DisableHighlighting: DefaultDisableHighlighting(),
		CursorStyle:         DefaultCursorStyle(),
		CursorColor:         DefaultCursorColor(),
	}
}

func Parse(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := NewDefault()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
