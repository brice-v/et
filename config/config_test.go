package config

import (
	"encoding/json"
	"testing"
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
			t.Errorf("Quit[%d] = %v, want %v", i, kb.Quit[i].Key, want.Quit[i].Key)
		}
	}
}
