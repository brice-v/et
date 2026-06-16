package config

import (
	"encoding/json"
	"et/defaults"
	"os"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Color struct {
	tcell.Color
}

func (c *Color) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	c.Color = color.GetColor(s)
	return nil
}

func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

type Colors struct {
	Foreground Color `json:"foreground"`
	Background Color `json:"background"`
	StatusBar  Color `json:"status_bar"`
}

type Config struct {
	Colors Colors `json:"colors"`
}

func NewDefault() *Config {
	return &Config{
		Colors: Colors{
			Foreground: Color{defaults.ColorForeground()},
			Background: Color{defaults.ColorBackground()},
			StatusBar:  Color{defaults.ColorStatusBar()},
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
