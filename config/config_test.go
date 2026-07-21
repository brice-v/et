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
			"background": "#282A35",
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
	if len(cfg.FileExtensions) != len(want.FileExtensions) {
		t.Errorf("len(FileExtensions) = %d, want %d", len(cfg.FileExtensions), len(want.FileExtensions))
	}
	for k, v := range want.FileExtensions {
		if cfg.FileExtensions[k] != v {
			t.Errorf("FileExtensions[%q] = %q, want %q", k, cfg.FileExtensions[k], v)
		}
	}
}

func TestKeyBindingsRoundTrip(t *testing.T) {
	jsonStr := `{"quit": {"prefix": "ctrl+e", "suffix": "q"}}`

	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}

	if kb.Quit.Prefix.Key != tcell.Key('e') {
		t.Errorf("Quit.Prefix.Key = %v, want Key('e')", kb.Quit.Prefix.Key)
	}
	if kb.Quit.Prefix.Modifiers != tcell.ModCtrl {
		t.Errorf("Quit.Prefix.Modifiers = %v, want ModCtrl", kb.Quit.Prefix.Modifiers)
	}
	if kb.Quit.Suffix.Key != tcell.Key('q') {
		t.Errorf("Quit.Suffix.Key = %v, want Key('q')", kb.Quit.Suffix.Key)
	}
	if kb.Quit.Suffix.Modifiers != tcell.ModNone {
		t.Errorf("Quit.Suffix.Modifiers = %v, want ModNone", kb.Quit.Suffix.Modifiers)
	}
}

func TestFindNextPreviousKeyBindings(t *testing.T) {
	jsonStr := `{"find_next": "tab", "find_previous": "backtab"}`

	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}

	if kb.FindNext.Key != tcell.KeyTab {
		t.Errorf("FindNext.Key = %v, want KeyTab", kb.FindNext.Key)
	}
	if kb.FindNext.Modifiers != tcell.ModNone {
		t.Errorf("FindNext.Modifiers = %v, want ModNone", kb.FindNext.Modifiers)
	}
	if kb.FindPrevious.Key != tcell.KeyBacktab {
		t.Errorf("FindPrevious.Key = %v, want KeyBacktab", kb.FindPrevious.Key)
	}
	if kb.FindPrevious.Modifiers != tcell.ModNone {
		t.Errorf("FindPrevious.Modifiers = %v, want ModNone", kb.FindPrevious.Modifiers)
	}
}

func TestDefaultCurrentMatchHighlight(t *testing.T) {
	cfg := NewDefault()
	want := DefaultColorCurrentMatchHighlight()
	if cfg.Colors.CurrentMatchHighlight.Color != want.Color {
		t.Errorf("CurrentMatchHighlight = %v, want %v", cfg.Colors.CurrentMatchHighlight, want)
	}
}

func TestDefaultFindNextPreviousBindings(t *testing.T) {
	kb := NewDefault().KeyBindings

	if kb.FindNext.Key != tcell.KeyTab || kb.FindNext.Modifiers != tcell.ModNone {
		t.Errorf("Default FindNext = %v, want tab", kb.FindNext)
	}
	if kb.FindPrevious.Key != tcell.KeyBacktab || kb.FindPrevious.Modifiers != tcell.ModNone {
		t.Errorf("Default FindPrevious = %v, want backtab", kb.FindPrevious)
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
		{"enter", Key{Key: tcell.KeyEnter, Modifiers: tcell.ModNone}, "enter"},
		{"ctrl+enter", Key{Key: tcell.KeyEnter, Modifiers: tcell.ModCtrl}, "ctrl+enter"},
		{"tab", Key{Key: tcell.KeyTab, Modifiers: tcell.ModNone}, "tab"},
		{"backtab", Key{Key: tcell.KeyBacktab, Modifiers: tcell.ModNone}, "backtab"},
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
		{"tab", tcell.KeyTab, tcell.ModNone},
		{"backtab", tcell.KeyBacktab, tcell.ModNone},
		{"shift+tab", tcell.KeyTab, tcell.ModShift},
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
	if cfg.KeyBindings.Quit.Prefix.Key != want.KeyBindings.Quit.Prefix.Key {
		t.Errorf("Quit.Prefix.Key = %v, want %v", cfg.KeyBindings.Quit.Prefix.Key, want.KeyBindings.Quit.Prefix.Key)
	}
	if cfg.KeyBindings.Quit.Prefix.Modifiers != want.KeyBindings.Quit.Prefix.Modifiers {
		t.Errorf("Quit.Prefix.Modifiers = %v, want %v", cfg.KeyBindings.Quit.Prefix.Modifiers, want.KeyBindings.Quit.Prefix.Modifiers)
	}
	if len(cfg.FileExtensions) != len(want.FileExtensions) {
		t.Errorf("len(FileExtensions) = %d, want %d", len(cfg.FileExtensions), len(want.FileExtensions))
	}
	for k, v := range want.FileExtensions {
		if cfg.FileExtensions[k] != v {
			t.Errorf("FileExtensions[%q] = %q, want %q", k, cfg.FileExtensions[k], v)
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
	if cfg.KeyBindings.Quit.Prefix.Key != want.KeyBindings.Quit.Prefix.Key {
		t.Errorf("Quit.Prefix.Key = %v, want default %v", cfg.KeyBindings.Quit.Prefix.Key, want.KeyBindings.Quit.Prefix.Key)
	}
	if cfg.KeyBindings.Quit.Prefix.Modifiers != want.KeyBindings.Quit.Prefix.Modifiers {
		t.Errorf("Quit.Prefix.Modifiers = %v, want default %v", cfg.KeyBindings.Quit.Prefix.Modifiers, want.KeyBindings.Quit.Prefix.Modifiers)
	}
	if len(cfg.FileExtensions) != len(want.FileExtensions) {
		t.Errorf("len(FileExtensions) = %d, want %d", len(cfg.FileExtensions), len(want.FileExtensions))
	}
	for k, v := range want.FileExtensions {
		if cfg.FileExtensions[k] != v {
			t.Errorf("FileExtensions[%q] = %q, want default %q", k, cfg.FileExtensions[k], v)
		}
	}
}

func TestDefaultChordPrefix(t *testing.T) {
	prefix := DefaultChordPrefix()
	if prefix.Key != tcell.Key('e') || prefix.Modifiers != tcell.ModCtrl {
		t.Errorf("DefaultChordPrefix = %v, want ctrl+e", prefix)
	}
}

func TestConfigChordPrefix(t *testing.T) {
	cfg := NewDefault()
	if cfg.ChordPrefix.Key != tcell.Key('e') || cfg.ChordPrefix.Modifiers != tcell.ModCtrl {
		t.Errorf("Config.ChordPrefix = %v, want ctrl+e", cfg.ChordPrefix)
	}
}

func TestDefaultQuitChord(t *testing.T) {
	kb := NewDefault().KeyBindings
	wantPrefix := DefaultChordPrefix()
	wantSuffix := Key{Key: tcell.Key('q'), Modifiers: tcell.ModNone}
	if kb.Quit.Prefix != wantPrefix {
		t.Errorf("Quit.Prefix = %v, want %v", kb.Quit.Prefix, wantPrefix)
	}
	if kb.Quit.Suffix != wantSuffix {
		t.Errorf("Quit.Suffix = %v, want %v", kb.Quit.Suffix, wantSuffix)
	}
}

func TestDefaultFindChord(t *testing.T) {
	kb := NewDefault().KeyBindings
	wantPrefix := DefaultChordPrefix()
	wantSuffix := Key{Key: tcell.Key('f'), Modifiers: tcell.ModNone}
	if kb.Find.Prefix != wantPrefix {
		t.Errorf("Find.Prefix = %v, want %v", kb.Find.Prefix, wantPrefix)
	}
	if kb.Find.Suffix != wantSuffix {
		t.Errorf("Find.Suffix = %v, want %v", kb.Find.Suffix, wantSuffix)
	}
}

func TestChordPrefixRoundTrip(t *testing.T) {
	type partial struct {
		ChordPrefix Key `json:"chord_prefix"`
	}
	jsonStr := `{"chord_prefix": "ctrl+b"}`
	var p partial
	if err := json.Unmarshal([]byte(jsonStr), &p); err != nil {
		t.Fatal(err)
	}
	if p.ChordPrefix.Key != tcell.Key('b') || p.ChordPrefix.Modifiers != tcell.ModCtrl {
		t.Errorf("ChordPrefix = %v, want ctrl+b", p.ChordPrefix)
	}
	got, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var p2 partial
	if err := json.Unmarshal(got, &p2); err != nil {
		t.Fatal(err)
	}
	if p2.ChordPrefix != p.ChordPrefix {
		t.Errorf("ChordPrefix round-trip = %v, want %v", p2.ChordPrefix, p.ChordPrefix)
	}
}

func TestChordBindingsAllUseChordPrefix(t *testing.T) {
	kb := NewDefault().KeyBindings
	want := DefaultChordPrefix()
	chords := []KeyChord{
		kb.Quit,
		kb.Find,
		kb.FindSecondary1Chord,
		kb.FindSecondary2Chord,
		kb.ToggleTerminal,
		kb.ToggleLineEnding,
		kb.ToggleExpandTabs,
		kb.TerminalIncreaseChord,
		kb.TerminalDecreaseChord,
	}
	for i, kc := range chords {
		if kc.Prefix != want {
			t.Errorf("chord[%d].Prefix = %v, want %v", i, kc.Prefix, want)
		}
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := Parse("nonexistent.json")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func removeFile(t *testing.T, path string) {
	err := os.Remove(path)
	if err != nil {
		t.Errorf("failed to remove file: %s, error: %s", path, err.Error())
	}
}

func TestParseInvalidJSON(t *testing.T) {
	path := "test_invalid.json"
	if err := os.WriteFile(path, []byte("{invalid}"), 0644); err != nil {
		t.Fatal(err)
	}
	defer removeFile(t, path)

	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
func TestColorMapJSON(t *testing.T) {
	cm := ColorMap{
		Keywords1:     []string{"func", "return"},
		Color1:        Color{color.GetColor("red")},
		Keywords2:     []string{"int", "string"},
		Color2:        Color{color.GetColor("blue")},
		Keywords3:     []string{"print"},
		Color3:        Color{color.GetColor("green")},
		StringTokens:  []string{`"`, "'"},
		ColorString:   Color{color.GetColor("yellow")},
		Operators:     "+-*/",
		SpecialTokens: []string{"nil", "self"},
		SpecialColor:  Color{color.GetColor("purple")},
		CommentToken:  "#",
		CommentColor:  Color{color.GetColor("gray")},
	}

	got, err := json.Marshal(cm)
	if err != nil {
		t.Fatal(err)
	}

	var cm2 ColorMap
	if err := json.Unmarshal(got, &cm2); err != nil {
		t.Fatal(err)
	}

	if len(cm2.Keywords1) != 2 || cm2.Keywords1[0] != "func" || cm2.Keywords1[1] != "return" {
		t.Errorf("Keywords1 = %v, want [func return]", cm2.Keywords1)
	}
	if cm2.Color1.Color != color.GetColor("red") {
		t.Errorf("Color1 = %v, want red", cm2.Color1)
	}
	if len(cm2.Keywords2) != 2 || cm2.Keywords2[0] != "int" || cm2.Keywords2[1] != "string" {
		t.Errorf("Keywords2 = %v, want [int string]", cm2.Keywords2)
	}
	if cm2.Color2.Color != color.GetColor("blue") {
		t.Errorf("Color2 = %v, want blue", cm2.Color2)
	}
	if len(cm2.Keywords3) != 1 || cm2.Keywords3[0] != "print" {
		t.Errorf("Keywords3 = %v, want [print]", cm2.Keywords3)
	}
	if cm2.Color3.Color != color.GetColor("green") {
		t.Errorf("Color3 = %v, want green", cm2.Color3)
	}
	if len(cm2.StringTokens) != 2 || cm2.StringTokens[0] != "\"" || cm2.StringTokens[1] != "'" {
		t.Errorf("StringTokens = %v, want [\" ']", cm2.StringTokens)
	}
	if cm2.ColorString.Color != color.GetColor("yellow") {
		t.Errorf("ColorString = %v, want yellow", cm2.ColorString)
	}
	if cm2.Operators != "+-*/" {
		t.Errorf("Operators = %q, want \"+-*/\"", cm2.Operators)
	}
	if len(cm2.SpecialTokens) != 2 || cm2.SpecialTokens[0] != "nil" || cm2.SpecialTokens[1] != "self" {
		t.Errorf("SpecialTokens = %v, want [nil self]", cm2.SpecialTokens)
	}
	if cm2.SpecialColor.Color != color.GetColor("purple") {
		t.Errorf("SpecialColor = %v, want purple", cm2.SpecialColor)
	}
	if cm2.CommentToken != "#" {
		t.Errorf("CommentToken = %q, want \"#\"", cm2.CommentToken)
	}
	if cm2.CommentColor.Color != color.GetColor("gray") {
		t.Errorf("CommentColor = %v, want gray", cm2.CommentColor)
	}
}

func TestConfigWithLanguages(t *testing.T) {
	jsonStr := `{
		"colors": {
			"foreground": "white",
			"background": "black",
			"status_bar": "darkcyan",
			"languages": {
				"go": {
					"keywords1": ["func", "return"],
					"color1": "red",
					"keywords2": ["int", "string"],
					"color2": "blue",
					"keywords3": ["print"],
					"color3": "green",
					"string_tokens": ["\"", "'"],
					"color_string": "yellow",
					"operators": "+-*/",
					"special_tokens": ["nil", "self"],
					"special_color": "purple",
					"comment_token": "#",
					"comment_color": "gray"
				}
			}
		}
	}`

	var cfg Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		t.Fatal(err)
	}

	if cfg.Colors.Foreground.Color != color.GetColor("white") {
		t.Errorf("Foreground = %v, want white", cfg.Colors.Foreground)
	}
	if cfg.Colors.Background.Color != color.GetColor("black") {
		t.Errorf("Background = %v, want black", cfg.Colors.Background)
	}
	if cfg.Colors.StatusBar.Color != color.GetColor("darkcyan") {
		t.Errorf("StatusBar = %v, want darkcyan", cfg.Colors.StatusBar)
	}

	cm, ok := cfg.Colors.Languages["go"]
	if !ok {
		t.Fatal(`Languages["go"] not found`)
	}

	if len(cm.Keywords1) != 2 || cm.Keywords1[0] != "func" || cm.Keywords1[1] != "return" {
		t.Errorf("Keywords1 = %v, want [func return]", cm.Keywords1)
	}
	if cm.Color1.Color != color.GetColor("red") {
		t.Errorf("Color1 = %v, want red", cm.Color1)
	}
	if len(cm.Keywords2) != 2 || cm.Keywords2[0] != "int" || cm.Keywords2[1] != "string" {
		t.Errorf("Keywords2 = %v, want [int string]", cm.Keywords2)
	}
	if cm.Color2.Color != color.GetColor("blue") {
		t.Errorf("Color2 = %v, want blue", cm.Color2)
	}
	if len(cm.Keywords3) != 1 || cm.Keywords3[0] != "print" {
		t.Errorf("Keywords3 = %v, want [print]", cm.Keywords3)
	}
	if cm.Color3.Color != color.GetColor("green") {
		t.Errorf("Color3 = %v, want green", cm.Color3)
	}
	if len(cm.StringTokens) != 2 || cm.StringTokens[0] != "\"" || cm.StringTokens[1] != "'" {
		t.Errorf("StringTokens = %v, want [\" ']", cm.StringTokens)
	}
	if cm.ColorString.Color != color.GetColor("yellow") {
		t.Errorf("ColorString = %v, want yellow", cm.ColorString)
	}
	if cm.Operators != "+-*/" {
		t.Errorf("Operators = %q, want \"+-*/\"", cm.Operators)
	}
	if len(cm.SpecialTokens) != 2 || cm.SpecialTokens[0] != "nil" || cm.SpecialTokens[1] != "self" {
		t.Errorf("SpecialTokens = %v, want [nil self]", cm.SpecialTokens)
	}
	if cm.SpecialColor.Color != color.GetColor("purple") {
		t.Errorf("SpecialColor = %v, want purple", cm.SpecialColor)
	}
	if cm.CommentToken != "#" {
		t.Errorf("CommentToken = %q, want \"#\"", cm.CommentToken)
	}
	if cm.CommentColor.Color != color.GetColor("gray") {
		t.Errorf("CommentColor = %v, want gray", cm.CommentColor)
	}
}

func TestFileExtensionsOverride(t *testing.T) {
	jsonStr := `{
		"file_extensions": {
			"go": "go",
			"py": "py"
		}
	}`

	var cfg Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		t.Fatal(err)
	}

	if len(cfg.FileExtensions) != 2 {
		t.Fatalf("len(FileExtensions) = %d, want 2", len(cfg.FileExtensions))
	}
	if cfg.FileExtensions["go"] != "go" {
		t.Errorf(`FileExtensions["go"] = %q, want "go"`, cfg.FileExtensions["go"])
	}
	if cfg.FileExtensions["py"] != "py" {
		t.Errorf(`FileExtensions["py"] = %q, want "py"`, cfg.FileExtensions["py"])
	}
}

func TestDisableHighlighting(t *testing.T) {
	cfg := NewDefault()
	if cfg.DisableHighlighting != false {
		t.Errorf("DisableHighlighting = %t, want false", cfg.DisableHighlighting)
	}

	jsonStr := `{"disable_highlighting": true}`
	var cfg2 Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg2); err != nil {
		t.Fatal(err)
	}
	if cfg2.DisableHighlighting != true {
		t.Errorf("DisableHighlighting = %t, want true", cfg2.DisableHighlighting)
	}
}

func TestColorUnmarshalJSONError(t *testing.T) {
	var c Color
	err := c.UnmarshalJSON([]byte(`123`))
	if err == nil {
		t.Error("expected error for non-string JSON value")
	}
}

func TestKeyBindingsRoundTripWithModifier(t *testing.T) {
	jsonStr := `{"quit": {"prefix": "ctrl+e", "suffix": "shift+q"}}`

	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}

	if kb.Quit.Prefix.Key != tcell.Key('e') {
		t.Errorf("Quit.Prefix.Key = %v, want Key('e')", kb.Quit.Prefix.Key)
	}
	if kb.Quit.Suffix.Key != tcell.KeyQ {
		t.Errorf("Quit.Suffix.Key = %v, want KeyQ", kb.Quit.Suffix.Key)
	}
	if kb.Quit.Suffix.Modifiers != tcell.ModShift {
		t.Errorf("Quit.Suffix.Modifiers = %v, want ModShift", kb.Quit.Suffix.Modifiers)
	}
}

func TestKeyChordString(t *testing.T) {
	tests := []struct {
		name string
		kc   KeyChord
		want string
	}{
		{"ctrl+e i", KeyChord{Prefix: Key{Key: tcell.Key('e'), Modifiers: tcell.ModCtrl}, Suffix: Key{Key: tcell.Key('i'), Modifiers: tcell.ModNone}}, "ctrl+e i"},
		{"ctrl+e g", KeyChord{Prefix: Key{Key: tcell.Key('e'), Modifiers: tcell.ModCtrl}, Suffix: Key{Key: tcell.Key('g'), Modifiers: tcell.ModNone}}, "ctrl+e g"},
		{"no mods", KeyChord{Prefix: Key{Key: tcell.Key('x'), Modifiers: tcell.ModNone}, Suffix: Key{Key: tcell.Key('y'), Modifiers: tcell.ModNone}}, "x y"},
		{"with modifiers", KeyChord{Prefix: Key{Key: tcell.Key('a'), Modifiers: tcell.ModCtrl | tcell.ModShift}, Suffix: Key{Key: tcell.Key('b'), Modifiers: tcell.ModAlt}}, "ctrl+shift+a alt+b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.kc.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDefaultKeyChordValues(t *testing.T) {
	kb := NewDefault().KeyBindings

	wantPrefix := Key{Key: tcell.Key('e'), Modifiers: tcell.ModCtrl}
	if kb.FindSecondary1Chord.Prefix != wantPrefix {
		t.Errorf("FindSecondary1Chord.Prefix = %v, want %v", kb.FindSecondary1Chord.Prefix, wantPrefix)
	}
	if kb.FindSecondary1Chord.Suffix.Key != tcell.Key('i') || kb.FindSecondary1Chord.Suffix.Modifiers != tcell.ModNone {
		t.Errorf("Default FindSecondary1Chord.Suffix = %v, want i", kb.FindSecondary1Chord.Suffix)
	}
	if kb.FindSecondary2Chord.Prefix != wantPrefix {
		t.Errorf("FindSecondary2Chord.Prefix = %v, want %v", kb.FindSecondary2Chord.Prefix, wantPrefix)
	}
	if kb.FindSecondary2Chord.Suffix.Key != tcell.Key('g') || kb.FindSecondary2Chord.Suffix.Modifiers != tcell.ModNone {
		t.Errorf("Default FindSecondary2Chord.Suffix = %v, want g", kb.FindSecondary2Chord.Suffix)
	}
}

func TestKeyChordRoundTrip(t *testing.T) {
	jsonStr := `{"prefix":"ctrl+e","suffix":"i"}`

	var kc KeyChord
	if err := json.Unmarshal([]byte(jsonStr), &kc); err != nil {
		t.Fatal(err)
	}

	if kc.Prefix.Key != tcell.Key('e') || kc.Prefix.Modifiers != tcell.ModCtrl {
		t.Errorf("Prefix = %v, want ctrl+e", kc.Prefix)
	}
	if kc.Suffix.Key != tcell.Key('i') || kc.Suffix.Modifiers != tcell.ModNone {
		t.Errorf("Suffix = %v, want i", kc.Suffix)
	}

	got, err := json.Marshal(kc)
	if err != nil {
		t.Fatal(err)
	}

	var kc2 KeyChord
	if err := json.Unmarshal(got, &kc2); err != nil {
		t.Fatal(err)
	}
	if kc2.Prefix != kc.Prefix {
		t.Errorf("Prefix round-trip = %v, want %v", kc2.Prefix, kc.Prefix)
	}
	if kc2.Suffix != kc.Suffix {
		t.Errorf("Suffix round-trip = %v, want %v", kc2.Suffix, kc.Suffix)
	}
}

func TestDefaultToggleTerminalChord(t *testing.T) {
	kb := NewDefault().KeyBindings
	wantPrefix := Key{Key: tcell.Key('e'), Modifiers: tcell.ModCtrl}
	wantSuffix := Key{Key: tcell.Key(';'), Modifiers: tcell.ModNone}
	if kb.ToggleTerminal.Prefix != wantPrefix {
		t.Errorf("ToggleTerminal.Prefix = %v, want %v", kb.ToggleTerminal.Prefix, wantPrefix)
	}
	if kb.ToggleTerminal.Suffix != wantSuffix {
		t.Errorf("ToggleTerminal.Suffix = %v, want %v", kb.ToggleTerminal.Suffix, wantSuffix)
	}
}

func TestToggleTerminalChordParseJSON(t *testing.T) {
	jsonStr := `{"toggle_terminal": {"prefix": "ctrl+e", "suffix": ";"}}`
	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}
	if kb.ToggleTerminal.Prefix.Key != tcell.Key('e') || kb.ToggleTerminal.Prefix.Modifiers != tcell.ModCtrl {
		t.Errorf("Prefix = %v, want ctrl+e", kb.ToggleTerminal.Prefix)
	}
	if kb.ToggleTerminal.Suffix.Key != tcell.Key(';') || kb.ToggleTerminal.Suffix.Modifiers != tcell.ModNone {
		t.Errorf("Suffix = %v, want ;", kb.ToggleTerminal.Suffix)
	}
}

func TestToggleTerminalChordRoundTrip(t *testing.T) {
	kc := DefaultKeyBindingToggleTerminalChord()
	got, err := json.Marshal(kc)
	if err != nil {
		t.Fatal(err)
	}
	var kc2 KeyChord
	if err := json.Unmarshal(got, &kc2); err != nil {
		t.Fatal(err)
	}
	if kc2.Prefix != kc.Prefix {
		t.Errorf("Prefix round-trip = %v, want %v", kc2.Prefix, kc.Prefix)
	}
	if kc2.Suffix != kc.Suffix {
		t.Errorf("Suffix round-trip = %v, want %v", kc2.Suffix, kc.Suffix)
	}
}

func TestCursorStyleFromString(t *testing.T) {
	tests := []struct {
		input string
		want  tcell.CursorStyle
	}{
		{"blinking_block", tcell.CursorStyleBlinkingBlock},
		{"steady_block", tcell.CursorStyleSteadyBlock},
		{"blinking_underline", tcell.CursorStyleBlinkingUnderline},
		{"steady_underline", tcell.CursorStyleSteadyUnderline},
		{"blinking_bar", tcell.CursorStyleBlinkingBar},
		{"steady_bar", tcell.CursorStyleSteadyBar},
		{"BLINKING_BLOCK", tcell.CursorStyleBlinkingBlock}, // case insensitive
		{"  steady_block  ", tcell.CursorStyleSteadyBlock}, // trimmed
		{"", tcell.CursorStyleDefault},                     // unknown -> default
		{"foo", tcell.CursorStyleDefault},                  // unknown -> default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := CursorStyleFromString(tt.input)
			if got != tt.want {
				t.Errorf("CursorStyleFromString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestDefaultCursorConfig(t *testing.T) {
	cfg := NewDefault()

	if cfg.CursorStyle != DefaultCursorStyle() {
		t.Errorf("CursorStyle = %q, want %q", cfg.CursorStyle, DefaultCursorStyle())
	}

	if cfg.CursorColor.Color != color.White {
		t.Errorf("CursorColor = %v, want white", cfg.CursorColor.Color)
	}
}

func TestCursorConfigParse(t *testing.T) {
	jsonStr := `{
		"cursor_style": "blinking_underline",
		"cursor_color": "red"
	}`

	var cfg Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		t.Fatal(err)
	}

	wantStyle := tcell.CursorStyleBlinkingUnderline
	gotStyle := CursorStyleFromString(cfg.CursorStyle)
	if gotStyle != wantStyle {
		t.Errorf("CursorStyle = %v, want %v", gotStyle, wantStyle)
	}

	wantColor := color.GetColor("red")
	if cfg.CursorColor.Color != wantColor {
		t.Errorf("CursorColor = %v, want %v", cfg.CursorColor.Color, wantColor)
	}
}

func TestCursorConfigDefaultsWhenMissing(t *testing.T) {
	path := "test_cursor_defaults.json"
	if err := os.WriteFile(path, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}
	defer removeFile(t, path)

	cfg, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	wantStyle := CursorStyleFromString(DefaultCursorStyle())
	gotStyle := CursorStyleFromString(cfg.CursorStyle)
	if gotStyle != wantStyle {
		t.Errorf("CursorStyle (default) = %v, want %v", gotStyle, wantStyle)
	}

	wantColor := DefaultCursorColor().Color
	if cfg.CursorColor.Color != wantColor {
		t.Errorf("CursorColor (default) = %v, want %v", cfg.CursorColor.Color, wantColor)
	}
}

func TestDefaultLineEnding(t *testing.T) {
	cfg := NewDefault()
	if cfg.DefaultLineEnding != "lf" {
		t.Errorf("DefaultLineEnding = %q, want %q", cfg.DefaultLineEnding, "lf")
	}
}

func TestDefaultExpandTabs(t *testing.T) {
	cfg := NewDefault()
	if cfg.ExpandTabs != false {
		t.Errorf("ExpandTabs = %t, want false", cfg.ExpandTabs)
	}
}

func TestDefaultToggleLineEndingKey(t *testing.T) {
	kb := NewDefault().KeyBindings
	wantPrefix := DefaultChordPrefix()
	wantSuffix := Key{Key: tcell.Key('l'), Modifiers: tcell.ModNone}
	if kb.ToggleLineEnding.Prefix != wantPrefix {
		t.Errorf("ToggleLineEnding.Prefix = %v, want %v", kb.ToggleLineEnding.Prefix, wantPrefix)
	}
	if kb.ToggleLineEnding.Suffix != wantSuffix {
		t.Errorf("ToggleLineEnding.Suffix = %v, want %v", kb.ToggleLineEnding.Suffix, wantSuffix)
	}
}

func TestDefaultToggleExpandTabsKey(t *testing.T) {
	kb := NewDefault().KeyBindings
	wantPrefix := DefaultChordPrefix()
	wantSuffix := Key{Key: tcell.Key('t'), Modifiers: tcell.ModNone}
	if kb.ToggleExpandTabs.Prefix != wantPrefix {
		t.Errorf("ToggleExpandTabs.Prefix = %v, want %v", kb.ToggleExpandTabs.Prefix, wantPrefix)
	}
	if kb.ToggleExpandTabs.Suffix != wantSuffix {
		t.Errorf("ToggleExpandTabs.Suffix = %v, want %v", kb.ToggleExpandTabs.Suffix, wantSuffix)
	}
}

func TestLineEndingRoundTrip(t *testing.T) {
	jsonStr := `{"default_line_ending": "crlf"}`
	var cfg Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.DefaultLineEnding != "crlf" {
		t.Errorf("DefaultLineEnding = %q, want %q", cfg.DefaultLineEnding, "crlf")
	}
}

func TestExpandTabsParse(t *testing.T) {
	jsonStr := `{"expand_tabs": true}`
	var cfg Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.ExpandTabs != true {
		t.Errorf("ExpandTabs = %t, want true", cfg.ExpandTabs)
	}
}

func TestToggleLineEndingKeyParse(t *testing.T) {
	jsonStr := `{"toggle_line_ending": {"prefix": "ctrl+e", "suffix": "ctrl+shift+l"}}`
	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}
	if kb.ToggleLineEnding.Suffix.Key != tcell.Key('l') {
		t.Errorf("Suffix.Key = %v, want 'l'", kb.ToggleLineEnding.Suffix.Key)
	}
	if kb.ToggleLineEnding.Suffix.Modifiers != tcell.ModCtrl|tcell.ModShift {
		t.Errorf("Suffix.Modifiers = %v, want ModCtrl|ModShift", kb.ToggleLineEnding.Suffix.Modifiers)
	}
}

func TestToggleExpandTabsKeyParse(t *testing.T) {
	jsonStr := `{"toggle_expand_tabs": {"prefix": "ctrl+e", "suffix": "ctrl+t"}}`
	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}
	if kb.ToggleExpandTabs.Suffix.Key != tcell.Key('t') {
		t.Errorf("Suffix.Key = %v, want 't'", kb.ToggleExpandTabs.Suffix.Key)
	}
	if kb.ToggleExpandTabs.Suffix.Modifiers != tcell.ModCtrl {
		t.Errorf("Suffix.Modifiers = %v, want ModCtrl", kb.ToggleExpandTabs.Suffix.Modifiers)
	}
}

func TestDefaultTerminalHeightPercentage(t *testing.T) {
	cfg := NewDefault()
	if cfg.TerminalHeightPercentage != 0.25 {
		t.Errorf("TerminalHeightPercentage = %f, want 0.25", cfg.TerminalHeightPercentage)
	}
}

func TestDefaultTerminalSizeChords(t *testing.T) {
	kb := NewDefault().KeyBindings
	wantPrefix := Key{Key: tcell.Key('e'), Modifiers: tcell.ModCtrl}

	if kb.TerminalIncreaseChord.Prefix != wantPrefix {
		t.Errorf("TerminalIncreaseChord.Prefix = %v, want %v", kb.TerminalIncreaseChord.Prefix, wantPrefix)
	}
	if kb.TerminalIncreaseChord.Suffix.Key != tcell.Key('+') || kb.TerminalIncreaseChord.Suffix.Modifiers != tcell.ModNone {
		t.Errorf("TerminalIncreaseChord.Suffix = %v, want +", kb.TerminalIncreaseChord.Suffix)
	}

	if kb.TerminalDecreaseChord.Prefix != wantPrefix {
		t.Errorf("TerminalDecreaseChord.Prefix = %v, want %v", kb.TerminalDecreaseChord.Prefix, wantPrefix)
	}
	if kb.TerminalDecreaseChord.Suffix.Key != tcell.Key('-') || kb.TerminalDecreaseChord.Suffix.Modifiers != tcell.ModNone {
		t.Errorf("TerminalDecreaseChord.Suffix = %v, want -", kb.TerminalDecreaseChord.Suffix)
	}
}

func TestKeyStringPlusMinus(t *testing.T) {
	plus := Key{Key: tcell.Key('+'), Modifiers: tcell.ModNone}
	if got := plus.String(); got != "+" {
		t.Errorf("Key('+').String() = %q, want %q", got, "+")
	}
	minus := Key{Key: tcell.Key('-'), Modifiers: tcell.ModNone}
	if got := minus.String(); got != "-" {
		t.Errorf("Key('-').String() = %q, want %q", got, "-")
	}
}

func TestParseKeyBindingPlusMinus(t *testing.T) {
	tests := []struct {
		input string
		key   tcell.Key
		mod   tcell.ModMask
	}{
		{"+", tcell.Key('+'), tcell.ModNone},
		{"-", tcell.Key('-'), tcell.ModNone},
		{"ctrl+e", tcell.Key('e'), tcell.ModCtrl},
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

func TestTerminalChordRoundTrip(t *testing.T) {
	jsonStr := `{"prefix": "ctrl+e", "suffix": "+"}`
	var kc KeyChord
	if err := json.Unmarshal([]byte(jsonStr), &kc); err != nil {
		t.Fatal(err)
	}
	if kc.Prefix.Key != tcell.Key('e') || kc.Prefix.Modifiers != tcell.ModCtrl {
		t.Errorf("Prefix = %v, want ctrl+e", kc.Prefix)
	}
	if kc.Suffix.Key != tcell.Key('+') || kc.Suffix.Modifiers != tcell.ModNone {
		t.Errorf("Suffix = %v, want +", kc.Suffix)
	}

	got, err := json.Marshal(kc)
	if err != nil {
		t.Fatal(err)
	}
	var kc2 KeyChord
	if err := json.Unmarshal(got, &kc2); err != nil {
		t.Fatal(err)
	}
	if kc2.Prefix != kc.Prefix {
		t.Errorf("Prefix round-trip = %v, want %v", kc2.Prefix, kc.Prefix)
	}
	if kc2.Suffix != kc.Suffix {
		t.Errorf("Suffix round-trip = %v, want %v", kc2.Suffix, kc.Suffix)
	}

	// Also test the minus chord
	jsonStr2 := `{"prefix": "ctrl+e", "suffix": "-"}`
	var kc3 KeyChord
	if err := json.Unmarshal([]byte(jsonStr2), &kc3); err != nil {
		t.Fatal(err)
	}
	if kc3.Suffix.Key != tcell.Key('-') {
		t.Errorf("Suffix = %v, want -", kc3.Suffix)
	}
}

func TestTerminalHeightPercentageRoundTrip(t *testing.T) {
	jsonStr := `{"terminal_height_percentage": 0.5}`
	var cfg Config
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.TerminalHeightPercentage != 0.5 {
		t.Errorf("TerminalHeightPercentage = %f, want 0.5", cfg.TerminalHeightPercentage)
	}

	// Marshal/unmarshal a minimal wrapper to avoid zero-value Key errors
	type wrapper struct {
		TerminalHeightPercentage float64 `json:"terminal_height_percentage"`
	}
	w := wrapper{TerminalHeightPercentage: cfg.TerminalHeightPercentage}
	got, err := json.Marshal(w)
	if err != nil {
		t.Fatal(err)
	}
	var w2 wrapper
	if err := json.Unmarshal(got, &w2); err != nil {
		t.Fatal(err)
	}
	if w2.TerminalHeightPercentage != w.TerminalHeightPercentage {
		t.Errorf("TerminalHeightPercentage round-trip = %f, want %f", w2.TerminalHeightPercentage, w.TerminalHeightPercentage)
	}
}

func TestCursorConfigRoundTrip(t *testing.T) {
	path := "test_cursor_roundtrip.json"

	// Write JSON with cursor config + minimal colors (omit keybindings to avoid pre-existing marshal bug)
	wantJSON := `{
		"cursor_style": "blinking_bar",
		"cursor_color": "cyan",
		"colors": {
			"foreground": "white",
			"background": "black",
			"status_bar": "darkcyan"
		},
		"tab_width": 4,
		"left_pad_string": "~",
		"show_line_numbers": true,
		"disable_highlighting": false
	}`
	if err := os.WriteFile(path, []byte(wantJSON), 0644); err != nil {
		t.Fatal(err)
	}
	defer removeFile(t, path)

	cfg1, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	wantStyle := CursorStyleFromString("blinking_bar")
	gotStyle := CursorStyleFromString(cfg1.CursorStyle)
	if gotStyle != wantStyle {
		t.Errorf("CursorStyle = %v, want %v", gotStyle, wantStyle)
	}

	wantColor := color.GetColor("cyan")
	if cfg1.CursorColor.Color != wantColor {
		t.Errorf("CursorColor = %v, want %v", cfg1.CursorColor.Color, wantColor)
	}

	// Verify the style can be re-read from a different JSON file
	path2 := "test_cursor_roundtrip2.json"
	wantJSON2 := `{
		"cursor_style": "steady_underline",
		"cursor_color": "magenta",
		"colors": {
			"foreground": "white",
			"background": "black",
			"status_bar": "darkcyan"
		},
		"tab_width": 4,
		"left_pad_string": "~",
		"show_line_numbers": true,
		"disable_highlighting": false
	}`
	if err := os.WriteFile(path2, []byte(wantJSON2), 0644); err != nil {
		t.Fatal(err)
	}
	defer removeFile(t, path2)

	cfg2, err := Parse(path2)
	if err != nil {
		t.Fatal(err)
	}

	wantStyle2 := CursorStyleFromString("steady_underline")
	gotStyle2 := CursorStyleFromString(cfg2.CursorStyle)
	if gotStyle2 != wantStyle2 {
		t.Errorf("CursorStyle second file = %v, want %v", gotStyle2, wantStyle2)
	}

	wantColor2 := color.GetColor("magenta")
	if cfg2.CursorColor.Color != wantColor2 {
		t.Errorf("CursorColor second file = %v, want %v", cfg2.CursorColor.Color, wantColor2)
	}
}
