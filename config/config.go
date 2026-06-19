package config

import (
	"encoding/json"
	"et/defaults"
	"os"
	"strings"

	"github.com/gdamore/tcell/v3"
)

type Color struct {
	tcell.Color
}

type Key struct {
	tcell.Key
	Modifiers tcell.ModMask
}

type Colors struct {
	Foreground Color `json:"foreground"`
	Background Color `json:"background"`
	StatusBar  Color `json:"status_bar"`
}

type KeyBindings struct {
	Quit []Key `json:"quit"`
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
	Colors          Colors      `json:"colors"`
	KeyBindings     KeyBindings `json:"keybindings"`
	TabWidth        int         `json:"tab_width"`
	LeftPadString   string      `json:"left_pad_string"`
	ShowLineNumbers bool        `json:"show_line_numbers"`
}

func NewDefault() *Config {
	return &Config{
		Colors: Colors{
			Foreground: Color{defaults.ColorForeground()},
			Background: Color{defaults.ColorBackground()},
			StatusBar:  Color{defaults.ColorStatusBar()},
		},
		KeyBindings: KeyBindings{
			Quit: makeKeysFromKeyBinding(defaults.KeyBindingsQuit()),
		},
		TabWidth:        defaults.TabWidth(),
		LeftPadString:   "~",
		ShowLineNumbers: true,
	}
}

func Parse(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
