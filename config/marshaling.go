package config

import (
	"encoding/json"
	"et/defaults"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func (c *Color) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	c.Color = color.GetColor(s)
	return nil
}

func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (k *Key) String() string {
	var parts []string
	if k.Modifiers&tcell.ModCtrl != 0 {
		parts = append(parts, "ctrl")
	}
	if k.Modifiers&tcell.ModShift != 0 {
		parts = append(parts, "shift")
	}
	if k.Modifiers&tcell.ModAlt != 0 {
		parts = append(parts, "alt")
	}
	base := ""
	switch k.Key {
	case tcell.KeyQ:
		base = "q"
	case tcell.KeyEscape:
		base = "esc"
	default:
		return ""
	}
	if len(parts) == 0 {
		return base
	}
	return strings.Join(parts, "+") + "+" + base
}

func (k *Key) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	kk, mod, err := parseKeyBinding(s)
	if err != nil {
		return err
	}
	k.Key = kk
	k.Modifiers = mod
	return nil
}

func parseKeyBinding(s string) (tcell.Key, tcell.ModMask, error) {
	parts := strings.Split(s, "+")
	if len(parts) == 0 {
		return 0, 0, fmt.Errorf("empty key binding")
	}

	var mod tcell.ModMask
	var keyParts []string
	for _, p := range parts {
		switch p {
		case "ctrl":
			mod |= tcell.ModCtrl
		case "shift":
			mod |= tcell.ModShift
		case "alt":
			mod |= tcell.ModAlt
		default:
			keyParts = append(keyParts, p)
		}
	}

	if len(keyParts) != 1 {
		return 0, 0, fmt.Errorf("failed to parse %s as key binding", s)
	}

	keyStr := keyParts[0]
	var key tcell.Key
	switch keyStr {
	case "q":
		key = tcell.KeyQ
	case "esc", "escape":
		key = tcell.KeyEscape
	default:
		return 0, 0, fmt.Errorf("failed to parse %s as tcell key", keyStr)
	}

	return key, mod, nil
}

func (k Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func makeKeysFromKeyBinding(bindings []defaults.KeyBinding) []Key {
	keys := make([]Key, len(bindings))
	for i, b := range bindings {
		keys[i] = Key{Key: b.Key, Modifiers: b.Modifiers}
	}
	return keys
}
