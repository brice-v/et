package config

import (
	"encoding/json"
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
	if keyName, ok := tcell.KeyNames[k.Key]; ok {
		base = strings.ToLower(keyName)
	} else if k.Key >= tcell.Key('a') && k.Key <= tcell.Key('z') {
		base = string(rune(k.Key))
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

var keyNamesReverse map[string]tcell.Key

func initKeyNamesReverse() {
	keyNamesReverse = make(map[string]tcell.Key, len(tcell.KeyNames))
	for k, v := range tcell.KeyNames {
		keyNamesReverse[strings.ToLower(v)] = k
	}
}

func parseKeyBinding(s string) (tcell.Key, tcell.ModMask, error) {
	parts := strings.Split(s, "+")
	if len(parts) == 0 {
		return 0, 0, fmt.Errorf("empty key binding")
	}

	var mod tcell.ModMask
	var keyParts []string
	for _, p := range parts {
		switch strings.ToLower(p) {
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

	keyStr := strings.ToLower(strings.TrimSpace(keyParts[0]))
	if keyStr == "escape" {
		keyStr = "esc"
	}

	if keyNamesReverse == nil {
		initKeyNamesReverse()
	}

	if k, ok := keyNamesReverse[keyStr]; ok {
		return k, mod, nil
	}

	if len(keyStr) == 1 {
		return tcell.Key(keyStr[0]), mod, nil
	}

	return 0, 0, fmt.Errorf("failed to parse %s as key binding", s)
}

func (k Key) MarshalJSON() ([]byte, error) {
	s := k.String()
	if s == "" {
		return nil, fmt.Errorf("cannot marshal zero-value Key to JSON")
	}
	return json.Marshal(s)
}

func CursorStyleFromString(s string) tcell.CursorStyle {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "blinking_block":
		return tcell.CursorStyleBlinkingBlock
	case "steady_block":
		return tcell.CursorStyleSteadyBlock
	case "blinking_underline":
		return tcell.CursorStyleBlinkingUnderline
	case "steady_underline":
		return tcell.CursorStyleSteadyUnderline
	case "blinking_bar":
		return tcell.CursorStyleBlinkingBar
	case "steady_bar":
		return tcell.CursorStyleSteadyBar
	default:
		return tcell.CursorStyleDefault
	}
}
