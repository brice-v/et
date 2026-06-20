package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func TestParseDefaultJSON(t *testing.T) {
	jsonStr := `{
		"colors": {
			"foreground": "white",
			"background": "black",
			"status_bar": "darkcyan"
		}
	}`

	var cfg Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		t.Fatal(err)
	}

	want := NewDefault()

	if cfg.Colors.Foreground.Color != want.Colors.Foreground.Color {
		t.Errorf("Foreground = %v, want %v", cfg.Colors.Foreground, want.Colors.Foreground)
	}
	if cfg.Colors.Background.Color != want.Colors.Background.Color {
		t.Errorf("Background = %v, want %v", cfg.Colors.Background, want.Colors.Background)
	}
	if cfg.Colors.StatusBar.Color != want.Colors.StatusBar.Color {
		t.Errorf("StatusBar = %v, want %v", cfg.Colors.StatusBar, want.Colors.StatusBar)
	}
}

func TestRoundTrip(t *testing.T) {
	want := NewDefault()

	got, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}

	var cfg Config
	if err := json.Unmarshal(got, &cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.Colors.Foreground.Color != want.Colors.Foreground.Color {
		t.Errorf("Foreground = %v, want %v", cfg.Colors.Foreground, want.Colors.Foreground)
	}
	if cfg.Colors.Background.Color != want.Colors.Background.Color {
		t.Errorf("Background = %v, want %v", cfg.Colors.Background, want.Colors.Background)
	}
	if cfg.Colors.StatusBar.Color != want.Colors.StatusBar.Color {
		t.Errorf("StatusBar = %v, want %v", cfg.Colors.StatusBar, want.Colors.StatusBar)
	}
}

func TestKeyBindingsRoundTrip(t *testing.T) {
	jsonStr := `{"quit": ["ctrl+q", "q", "esc"]}`

	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}

	if len(kb.Quit) != 3 {
		t.Fatalf("len(Quit) = %d, want 3", len(kb.Quit))
	}

	want := NewDefault().KeyBindings
	for i := range want.Quit {
		if kb.Quit[i].Key != want.Quit[i].Key {
			t.Errorf("Quit[%d].Key = %v, want %v", i, kb.Quit[i].Key, want.Quit[i].Key)
		}
		if kb.Quit[i].Modifiers != want.Quit[i].Modifiers {
			t.Errorf("Quit[%d].Modifiers = %v, want %v", i, kb.Quit[i].Modifiers, want.Quit[i].Modifiers)
		}
	}
}

func TestKeyString(t *testing.T) {
	tests := []struct {
		name string
		key  Key
		want string
	}{
		{"q", Key{Key: tcell.KeyQ, Modifiers: tcell.ModNone}, "q"},
		{"ctrl+q", Key{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl}, "ctrl+q"},
		{"shift+q", Key{Key: tcell.KeyQ, Modifiers: tcell.ModShift}, "shift+q"},
		{"ctrl+shift+q", Key{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl | tcell.ModShift}, "ctrl+shift+q"},
		{"ctrl+alt+q", Key{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl | tcell.ModAlt}, "ctrl+alt+q"},
		{"ctrl+shift+alt+q", Key{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl | tcell.ModShift | tcell.ModAlt}, "ctrl+shift+alt+q"},
		{"esc", Key{Key: tcell.KeyEscape, Modifiers: tcell.ModNone}, "esc"},
		{"ctrl+esc", Key{Key: tcell.KeyEscape, Modifiers: tcell.ModCtrl}, "ctrl+esc"},
		{"unknown", Key{Key: tcell.KeyEnter, Modifiers: tcell.ModNone}, ""},
		{"unknown_with_mod", Key{Key: tcell.KeyEnter, Modifiers: tcell.ModCtrl}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.key.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseKeyBindingErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"unknown key", "foo"},
		{"unknown modifier", "super+q"},
		{"only modifiers", "ctrl+shift"},
		{"multiple key parts", "ctrl+q+z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseKeyBinding(tt.input)
			if err == nil {
				t.Errorf("parseKeyBinding(%q) expected error", tt.input)
			}
		})
	}
}

func TestKeyUnmarshalJSONErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid JSON", "{"},
		{"unknown key", `"foo"`},
		{"unknown modifier", `"super+q"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var k Key
			err := k.UnmarshalJSON([]byte(tt.input))
			if err == nil {
				t.Errorf("UnmarshalJSON(%q) expected error", tt.input)
			}
		})
	}
}

func TestColorJSON(t *testing.T) {
	c := Color{color.GetColor("red")}
	got, err := c.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `"red"` {
		t.Errorf("MarshalJSON = %s, want \"red\"", string(got))
	}

	var c2 Color
	if err := c2.UnmarshalJSON([]byte(`"blue"`)); err != nil {
		t.Fatal(err)
	}
	if c2.Color != color.GetColor("blue") {
		t.Errorf("UnmarshalJSON got %v, want blue", c2.Color)
	}
}

func TestParseModifierKeyBindings(t *testing.T) {
	tests := []struct {
		input string
		key   tcell.Key
		mod   tcell.ModMask
	}{
		{"ctrl+q", tcell.KeyQ, tcell.ModCtrl},
		{"shift+q", tcell.KeyQ, tcell.ModShift},
		{"ctrl+shift+q", tcell.KeyQ, tcell.ModCtrl | tcell.ModShift},
		{"ctrl+alt+q", tcell.KeyQ, tcell.ModCtrl | tcell.ModAlt},
		{"shift+alt+q", tcell.KeyQ, tcell.ModShift | tcell.ModAlt},
		{"ctrl+shift+alt+q", tcell.KeyQ, tcell.ModCtrl | tcell.ModShift | tcell.ModAlt},
		{"q", tcell.KeyQ, tcell.ModNone},
		{"esc", tcell.KeyEscape, tcell.ModNone},
		{"escape", tcell.KeyEscape, tcell.ModNone},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, mod, err := parseKeyBinding(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if key != tt.key {
				t.Errorf("key = %v, want %v", key, tt.key)
			}
			if mod != tt.mod {
				t.Errorf("mod = %v, want %v", mod, tt.mod)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	cfg, err := Parse("test_et_config.json")
	if err != nil {
		t.Fatal(err)
	}

	want := NewDefault()

	if cfg.Colors.Foreground.Color != want.Colors.Foreground.Color {
		t.Errorf("Foreground = %v, want %v", cfg.Colors.Foreground, want.Colors.Foreground)
	}
	if cfg.Colors.Background.Color != want.Colors.Background.Color {
		t.Errorf("Background = %v, want %v", cfg.Colors.Background, want.Colors.Background)
	}
	if cfg.Colors.StatusBar.Color != want.Colors.StatusBar.Color {
		t.Errorf("StatusBar = %v, want %v", cfg.Colors.StatusBar, want.Colors.StatusBar)
	}
	if cfg.TabWidth != want.TabWidth {
		t.Errorf("TabWidth = %d, want %d", cfg.TabWidth, want.TabWidth)
	}
	if cfg.LeftPadString != want.LeftPadString {
		t.Errorf("LeftPadString = %q, want %q", cfg.LeftPadString, want.LeftPadString)
	}
	if cfg.ShowLineNumbers != want.ShowLineNumbers {
		t.Errorf("ShowLineNumbers = %t, want %t", cfg.ShowLineNumbers, want.ShowLineNumbers)
	}
	if len(cfg.KeyBindings.Quit) != len(want.KeyBindings.Quit) {
		t.Fatalf("len(Quit) = %d, want %d", len(cfg.KeyBindings.Quit), len(want.KeyBindings.Quit))
	}
	for i := range want.KeyBindings.Quit {
		if cfg.KeyBindings.Quit[i].Key != want.KeyBindings.Quit[i].Key {
			t.Errorf("Quit[%d].Key = %v, want %v", i, cfg.KeyBindings.Quit[i].Key, want.KeyBindings.Quit[i].Key)
		}
		if cfg.KeyBindings.Quit[i].Modifiers != want.KeyBindings.Quit[i].Modifiers {
			t.Errorf("Quit[%d].Modifiers = %v, want %v", i, cfg.KeyBindings.Quit[i].Modifiers, want.KeyBindings.Quit[i].Modifiers)
		}
	}
}

func TestParseMinimalDefaults(t *testing.T) {
	cfg, err := Parse("test_et_config_minimal.json")
	if err != nil {
		t.Fatal(err)
	}

	want := NewDefault()

	if cfg.TabWidth != 2 {
		t.Errorf("TabWidth = %d, want 2", cfg.TabWidth)
	}
	if cfg.ShowLineNumbers != false {
		t.Errorf("ShowLineNumbers = %t, want false", cfg.ShowLineNumbers)
	}
	if cfg.Colors.Foreground.Color != want.Colors.Foreground.Color {
		t.Errorf("Foreground = %v, want default %v", cfg.Colors.Foreground, want.Colors.Foreground)
	}
	if cfg.Colors.Background.Color != want.Colors.Background.Color {
		t.Errorf("Background = %v, want default %v", cfg.Colors.Background, want.Colors.Background)
	}
	if cfg.Colors.StatusBar.Color != want.Colors.StatusBar.Color {
		t.Errorf("StatusBar = %v, want default %v", cfg.Colors.StatusBar, want.Colors.StatusBar)
	}
	if cfg.LeftPadString != want.LeftPadString {
		t.Errorf("LeftPadString = %q, want default %q", cfg.LeftPadString, want.LeftPadString)
	}
	if len(cfg.KeyBindings.Quit) != len(want.KeyBindings.Quit) {
		t.Fatalf("len(Quit) = %d, want %d", len(cfg.KeyBindings.Quit), len(want.KeyBindings.Quit))
	}
	for i := range want.KeyBindings.Quit {
		if cfg.KeyBindings.Quit[i].Key != want.KeyBindings.Quit[i].Key {
			t.Errorf("Quit[%d].Key = %v, want default %v", i, cfg.KeyBindings.Quit[i].Key, want.KeyBindings.Quit[i].Key)
		}
		if cfg.KeyBindings.Quit[i].Modifiers != want.KeyBindings.Quit[i].Modifiers {
			t.Errorf("Quit[%d].Modifiers = %v, want default %v", i, cfg.KeyBindings.Quit[i].Modifiers, want.KeyBindings.Quit[i].Modifiers)
		}
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := Parse("nonexistent.json")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestParseInvalidJSON(t *testing.T) {
	path := "test_invalid.json"
	if err := os.WriteFile(path, []byte("{invalid}"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestKeyBindingsRoundTripWithModifier(t *testing.T) {
	jsonStr := `{"quit": ["ctrl+shift+q", "ctrl+q", "q", "esc"]}`

	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}

	if len(kb.Quit) != 4 {
		t.Fatalf("len(Quit) = %d, want 4", len(kb.Quit))
	}

	want := []Key{
		{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl | tcell.ModShift},
		{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl},
		{Key: tcell.KeyQ, Modifiers: tcell.ModNone},
		{Key: tcell.KeyEscape, Modifiers: tcell.ModNone},
	}

	for i := range want {
		if kb.Quit[i].Key != want[i].Key {
			t.Errorf("Quit[%d].Key = %v, want %v", i, kb.Quit[i].Key, want[i].Key)
		}
		if kb.Quit[i].Modifiers != want[i].Modifiers {
			t.Errorf("Quit[%d].Modifiers = %v, want %v", i, kb.Quit[i].Modifiers, want[i].Modifiers)
		}
	}
}
