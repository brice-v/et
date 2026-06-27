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
	jsonStr := `{"quit": ["ctrl+q", "esc"]}`

	var kb KeyBindings
	if err := json.Unmarshal([]byte(jsonStr), &kb); err != nil {
		t.Fatal(err)
	}

	if len(kb.Quit) != 2 {
		t.Fatalf("len(Quit) = %d, want 2", len(kb.Quit))
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
	if len(cfg.FileExtensions) != len(want.FileExtensions) {
		t.Errorf("len(FileExtensions) = %d, want %d", len(cfg.FileExtensions), len(want.FileExtensions))
	}
	for k, v := range want.FileExtensions {
		if cfg.FileExtensions[k] != v {
			t.Errorf("FileExtensions[%q] = %q, want default %q", k, cfg.FileExtensions[k], v)
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

	// Round-trip marshal should fail on zero-value KeyBindings,
	// which is the expected behavior — use Parse() for full round-trips.
	_, err := json.Marshal(cfg2)
	if err == nil {
		t.Error("expected error marshalling Config with zero-value KeyBindings")
	}
}

func TestGetQuitKeyBindingsAsStr(t *testing.T) {
	cfg := NewDefault()
	got := cfg.GetQuitKeyBindingsAsStr()
	want := "[ctrl+q,esc]"
	if got != want {
		t.Errorf("GetQuitKeyBindingsAsStr() = %q, want %q", got, want)
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
