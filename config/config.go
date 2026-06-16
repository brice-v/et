package config

import (
	"encoding/json"
	"et/defaults"
	"os"

	"github.com/gdamore/tcell/v3"
)

type Color struct {
	tcell.Color
}

type Key struct {
	tcell.Key
}

type Colors struct {
	Foreground Color `json:"foreground"`
	Background Color `json:"background"`
	StatusBar  Color `json:"status_bar"`
}

type KeyBindings struct {
	Quit []Key `json:"quit"`
}

type Config struct {
	Colors      Colors      `json:"colors"`
	KeyBindings KeyBindings `json:"keybindings"`
}

func NewDefault() *Config {
	return &Config{
		Colors: Colors{
			Foreground: Color{defaults.ColorForeground()},
			Background: Color{defaults.ColorBackground()},
			StatusBar:  Color{defaults.ColorStatusBar()},
		},
		KeyBindings: KeyBindings{
			Quit: makeKeysFromTcellKeys(defaults.KeyBindingsQuit()),
		},
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
