package config

import (
	"encoding/json"
	"fmt"

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
	switch k.Key {
	case tcell.KeyCtrlQ:
		return "ctrl+q"
	case tcell.KeyQ:
		return "q"
	case tcell.KeyEscape:
		return "esc"
	default:
		return ""
	}
}

func (k *Key) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	kk, err := parseKeyBindingAsKey(s)
	if err != nil {
		return err
	}
	k.Key = kk
	return nil
}

func parseKeyBindingAsKey(s string) (tcell.Key, error) {
	switch s {
	case "ctrl+q":
		return tcell.KeyCtrlQ, nil
	case "q":
		return tcell.KeyQ, nil
	case "esc", "escape":
		return tcell.KeyEscape, nil
	}
	return tcell.Key0, fmt.Errorf("failed to parse %s as tcell key", s)
}

func (k Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func makeKeysFromTcellKeys(keys []tcell.Key) []Key {
	ks := make([]Key, len(keys))
	for i := range len(keys) {
		ks[i] = Key{keys[i]}
	}
	return ks
}
