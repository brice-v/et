package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"et/action"
	"et/defaults"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func defaultConfig() *config {
	return &config{
		Keybindings: defaults.Keybindings(),
		TabWidth:    defaults.TabWidth(),
	}
}

const version = "0.0.1"
const configFileName = "et.json"

type colorsConfig struct {
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	StatusBG   string `json:"status"`
}

type config struct {
	Colors      colorsConfig      `json:"colors"`
	Keybindings map[string]string `json:"keybindings"`
	TabWidth    int               `json:"tab_width"`
}

func parseColor(s string) (tcell.Color, error) {
	if s == "" {
		return color.Default, nil
	}
	c := color.GetColor(s)
	if c == color.Default {
		return c, fmt.Errorf("unknown color %q", s)
	}
	return c, nil
}

func findConfigPaths() (local, global string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	local = filepath.Join(wd, configFileName)

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", err
	}
	global = filepath.Join(configDir, "et", configFileName)
	return
}

func parseKeySpec(name string) (keySpec, error) {
	base := name
	modMask := tcell.ModNone

	type prefixEntry struct {
		s string
		m tcell.ModMask
	}
	prefixes := []prefixEntry{
		{"CtrlAltShift", tcell.ModCtrl | tcell.ModAlt | tcell.ModShift},
		{"CtrlShift", tcell.ModCtrl | tcell.ModShift},
		{"AltShift", tcell.ModAlt | tcell.ModShift},
		{"CtrlAlt", tcell.ModCtrl | tcell.ModAlt},
		{"Ctrl", tcell.ModCtrl},
		{"Alt", tcell.ModAlt},
		{"Shift", tcell.ModShift},
	}

	for _, p := range prefixes {
		if strings.HasPrefix(base, p.s) {
			modMask = p.m
			base = base[len(p.s):]
			break
		}
	}

	if modMask&tcell.ModCtrl != 0 && len(base) == 1 && base[0] >= 'A' && base[0] <= 'Z' {
		return keySpec{key: tcell.KeyCtrlA + tcell.Key(base[0]-'A')}, nil
	}

	namedKeys := map[string]tcell.Key{
		"Escape":    tcell.KeyEscape,
		"Enter":     tcell.KeyEnter,
		"Tab":       tcell.KeyTab,
		"Backspace": tcell.KeyBackspace,
		"Delete":    tcell.KeyDelete,
		"Insert":    tcell.KeyInsert,
		"Home":      tcell.KeyHome,
		"End":       tcell.KeyEnd,
		"Up":        tcell.KeyUp,
		"Down":      tcell.KeyDown,
		"Left":      tcell.KeyLeft,
		"Right":     tcell.KeyRight,
		"PgUp":      tcell.KeyPgUp,
		"PgDn":      tcell.KeyPgDn,
	}

	if k, ok := namedKeys[base]; ok {
		return keySpec{key: k, mods: modMask}, nil
	}

	if base == "Space" {
		return keySpec{key: tcell.KeyRune, str: " ", mods: modMask}, nil
	}

	if len(base) > 1 && base[0] == 'F' {
		var n int
		if _, err := fmt.Sscanf(base, "F%d", &n); err == nil && n >= 1 && n <= 64 {
			return keySpec{key: tcell.KeyF1 + tcell.Key(n-1), mods: modMask}, nil
		}
	}

	if len(base) == 1 {
		s := base
		if modMask&tcell.ModShift == 0 && s[0] >= 'A' && s[0] <= 'Z' {
			s = strings.ToLower(s)
		}
		return keySpec{key: tcell.KeyRune, str: s, mods: modMask}, nil
	}

	return keySpec{}, fmt.Errorf("unknown key name %q", name)
}

func loadConfig() *config {
	local, global, err := findConfigPaths()
	if err != nil {
		warn("could not determine config paths: %s", err)
		return defaultConfig()
	}

	for _, path := range []string{local, global} {
		data, err := os.ReadFile(path)
		if err != nil {
			if !os.IsNotExist(err) {
				warn("could not read %s: %s", path, err)
			}
			continue
		}
		var cfg config
		if err := json.Unmarshal(data, &cfg); err != nil {
			warn("could not parse %s: %s", path, err)
			continue
		}

		if cfg.Keybindings == nil {
			cfg.Keybindings = defaults.Keybindings()
		} else {
			valid := make(map[string]string, len(cfg.Keybindings))
			for name, act := range cfg.Keybindings {
				if _, err := parseKeySpec(name); err != nil {
					warn("invalid key %q: %s", name, err)
					continue
				}
				if _, err := action.Parse(act); err != nil {
					warn("unknown action %q for key %q: %s", act, name, err)
					continue
				}
				valid[name] = act
			}
			if len(valid) == 0 {
				valid = defaults.Keybindings()
			}
			cfg.Keybindings = valid
		}

		if cfg.Colors.Foreground != "" {
			if _, err := parseColor(cfg.Colors.Foreground); err != nil {
				warn("invalid foreground color %q: %s", cfg.Colors.Foreground, err)
				cfg.Colors.Foreground = ""
			}
		}
		if cfg.Colors.Background != "" {
			if _, err := parseColor(cfg.Colors.Background); err != nil {
				warn("invalid background color %q: %s", cfg.Colors.Background, err)
				cfg.Colors.Background = ""
			}
		}
		if cfg.Colors.StatusBG != "" {
			if _, err := parseColor(cfg.Colors.StatusBG); err != nil {
				warn("invalid status color %q: %s", cfg.Colors.StatusBG, err)
				cfg.Colors.StatusBG = ""
			}
		}

		if cfg.TabWidth <= 0 {
			warn("invalid tab_width %d, using default %d", cfg.TabWidth, defaults.TabWidth())
			cfg.TabWidth = defaults.TabWidth()
		}

		return &cfg
	}

	return defaultConfig()
}
