package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"et/defaults"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func captureLog(t *testing.T) *bytes.Buffer {
	var buf bytes.Buffer
	old := cfgLog
	cfgLog = log.New(&buf, "config: ", 0)
	t.Cleanup(func() { cfgLog = old })
	return &buf
}

// parseKeySpec tests

func TestParseCtrlLetter(t *testing.T) {
	ks, err := parseKeySpec("CtrlQ")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyCtrlQ {
		t.Fatalf("expected KeyCtrlQ, got %v", ks.key)
	}
}

func TestParseCtrlAltShiftLetter(t *testing.T) {
	ks, err := parseKeySpec("CtrlAltShiftA")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyCtrlA {
		t.Fatal("expected KeyCtrlA")
	}
}

func TestParseCtrlShiftLetter(t *testing.T) {
	ks, err := parseKeySpec("CtrlShiftZ")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyCtrlZ {
		t.Fatal("expected KeyCtrlZ")
	}
}

func TestParseAltShiftLetter(t *testing.T) {
	ks, err := parseKeySpec("AltShiftX")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyRune || ks.str != "X" || ks.mods != tcell.ModAlt|tcell.ModShift {
		t.Fatalf("unexpected spec: key=%v str=%q mods=%v", ks.key, ks.str, ks.mods)
	}
}

func TestParseShiftKey(t *testing.T) {
	ks, err := parseKeySpec("ShiftQ")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyRune || ks.str != "Q" || ks.mods != tcell.ModShift {
		t.Fatalf("unexpected spec: key=%v str=%q mods=%v", ks.key, ks.str, ks.mods)
	}
}

func TestParseCtrlAlt(t *testing.T) {
	ks, err := parseKeySpec("CtrlAltC")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyCtrlC {
		t.Fatal("expected KeyCtrlC")
	}
}

func TestParseAltKey(t *testing.T) {
	ks, err := parseKeySpec("AltQ")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyRune || ks.str != "q" || ks.mods != tcell.ModAlt {
		t.Fatalf("unexpected spec: key=%v str=%q mods=%v", ks.key, ks.str, ks.mods)
	}
}

func TestParseNamedKeys(t *testing.T) {
	tests := []struct {
		name string
		key  tcell.Key
	}{
		{"Escape", tcell.KeyEscape},
		{"Enter", tcell.KeyEnter},
		{"Tab", tcell.KeyTab},
		{"Backspace", tcell.KeyBackspace},
		{"Delete", tcell.KeyDelete},
		{"Insert", tcell.KeyInsert},
		{"Home", tcell.KeyHome},
		{"End", tcell.KeyEnd},
		{"Up", tcell.KeyUp},
		{"Down", tcell.KeyDown},
		{"Left", tcell.KeyLeft},
		{"Right", tcell.KeyRight},
		{"PgUp", tcell.KeyPgUp},
		{"PgDn", tcell.KeyPgDn},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks, err := parseKeySpec(tt.name)
			if err != nil {
				t.Fatal(err)
			}
			if ks.key != tt.key {
				t.Fatalf("expected %v, got %v", tt.key, ks.key)
			}
		})
	}
}

func TestParseSpace(t *testing.T) {
	ks, err := parseKeySpec("Space")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyRune || ks.str != " " {
		t.Fatalf("unexpected spec: key=%v str=%q", ks.key, ks.str)
	}
}

func TestParseFunctionKeys(t *testing.T) {
	tests := []struct {
		name string
		idx  int
	}{
		{"F1", 1},
		{"F10", 10},
		{"F12", 12},
		{"F64", 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks, err := parseKeySpec(tt.name)
			if err != nil {
				t.Fatal(err)
			}
			expected := tcell.KeyF1 + tcell.Key(tt.idx-1)
			if ks.key != expected {
				t.Fatalf("expected F%d key, got %v", tt.idx, ks.key)
			}
		})
	}
}

func TestParseInvalidFunctionKey(t *testing.T) {
	_, err := parseKeySpec("F0")
	if err == nil {
		t.Fatal("expected error for F0")
	}
	_, err = parseKeySpec("F65")
	if err == nil {
		t.Fatal("expected error for F65")
	}
}

func TestParseSingleRune(t *testing.T) {
	ks, err := parseKeySpec("a")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyRune || ks.str != "a" || ks.mods != tcell.ModNone {
		t.Fatalf("unexpected spec: key=%v str=%q mods=%v", ks.key, ks.str, ks.mods)
	}
}

func TestParseSingleRuneWithMod(t *testing.T) {
	ks, err := parseKeySpec("Altx")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyRune || ks.str != "x" || ks.mods != tcell.ModAlt {
		t.Fatalf("unexpected spec: key=%v str=%q mods=%v", ks.key, ks.str, ks.mods)
	}
}

func TestParseUnknownKey(t *testing.T) {
	_, err := parseKeySpec("SuperKey")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
}

func TestParseEmptyKey(t *testing.T) {
	_, err := parseKeySpec("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestParseCtrlAltShiftSingleRune(t *testing.T) {
	ks, err := parseKeySpec("CtrlAltShiftM")
	if err != nil {
		t.Fatal(err)
	}
	if ks.key != tcell.KeyCtrlM {
		t.Fatal("expected KeyCtrlM")
	}
}

// keySpec.matches tests

func TestMatchesRuneKey(t *testing.T) {
	ks := keySpec{key: tcell.KeyRune, str: "q", mods: tcell.ModNone}
	e := tcell.NewEventKey(tcell.KeyRune, "q", tcell.ModNone)
	if !ks.matches(e) {
		t.Fatal("should match KeyRune 'q'")
	}
}

func TestMatchesRuneKeyWrongStr(t *testing.T) {
	ks := keySpec{key: tcell.KeyRune, str: "q", mods: tcell.ModNone}
	e := tcell.NewEventKey(tcell.KeyRune, "w", tcell.ModNone)
	if ks.matches(e) {
		t.Fatal("should not match different rune")
	}
}

func TestMatchesRuneKeyWrongMod(t *testing.T) {
	ks := keySpec{key: tcell.KeyRune, str: "q", mods: tcell.ModAlt}
	e := tcell.NewEventKey(tcell.KeyRune, "q", tcell.ModNone)
	if ks.matches(e) {
		t.Fatal("should not match without Alt modifier")
	}
}

func TestMatchesCtrlKey(t *testing.T) {
	ks := keySpec{key: tcell.KeyCtrlQ}
	e := tcell.NewEventKey(tcell.KeyCtrlQ, "", tcell.ModNone)
	if !ks.matches(e) {
		t.Fatal("should match KeyCtrlQ")
	}
}

func TestMatchesCtrlKeyWrongKey(t *testing.T) {
	ks := keySpec{key: tcell.KeyCtrlQ}
	e := tcell.NewEventKey(tcell.KeyEscape, "", tcell.ModNone)
	if ks.matches(e) {
		t.Fatal("should not match Escape when expecting CtrlQ")
	}
}

func TestMatchesNamedKey(t *testing.T) {
	ks := keySpec{key: tcell.KeyEscape}
	e := tcell.NewEventKey(tcell.KeyEscape, "", tcell.ModNone)
	if !ks.matches(e) {
		t.Fatal("should match Escape")
	}
}

func TestMatchesAltRune(t *testing.T) {
	ks := keySpec{key: tcell.KeyRune, str: "x", mods: tcell.ModAlt}
	e := tcell.NewEventKey(tcell.KeyRune, "x", tcell.ModAlt)
	if !ks.matches(e) {
		t.Fatal("should match Alt+x")
	}
}

func TestMatchesKeyRuneEventWithNonRuneSpec(t *testing.T) {
	ks := keySpec{key: tcell.KeyEscape}
	e := tcell.NewEventKey(tcell.KeyRune, "a", tcell.ModNone)
	if ks.matches(e) {
		t.Fatal("should not match rune event when expecting Escape")
	}
}

// findConfigPaths tests

func TestFindConfigPaths(t *testing.T) {
	local, global, err := findConfigPaths()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(local, configFileName) {
		t.Fatalf("local path should end with %s, got %s", configFileName, local)
	}
	if !strings.HasSuffix(global, configFileName) {
		t.Fatalf("global path should end with %s, got %s", configFileName, global)
	}
	if local == global {
		t.Fatal("local and global paths should differ")
	}
}

// defaultKeybindings tests

func TestDefaultKeybindings(t *testing.T) {
	b := defaults.Keybindings()
	if len(b) != 2 {
		t.Fatalf("expected 2 default bindings, got %d", len(b))
	}
	if b["CtrlQ"] != "quit" {
		t.Fatal("CtrlQ should map to quit")
	}
	if b["Escape"] != "quit" {
		t.Fatal("Escape should map to quit")
	}
}

// loadConfig tests

func TestLoadConfigNoFile(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig()
	if len(cfg.Keybindings) != 2 {
		t.Fatalf("expected 2 default bindings, got %d", len(cfg.Keybindings))
	}
}

func TestLoadConfigValidFile(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{"keybindings": {"CtrlQ": "quit"}}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig()
	if len(cfg.Keybindings) != 1 {
		t.Fatalf("expected 1 binding, got %d", len(cfg.Keybindings))
	}
	if cfg.Keybindings["CtrlQ"] != "quit" {
		t.Fatal("CtrlQ should map to quit")
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{bad json}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	buf := captureLog(t)
	cfg := loadConfig()
	if cfg.Colors.Foreground != "" {
		t.Fatal("expected empty foreground for invalid JSON")
	}
	if len(cfg.Keybindings) != 2 {
		t.Fatalf("expected 2 default keybindings, got %d", len(cfg.Keybindings))
	}
	if !strings.Contains(buf.String(), "could not parse") {
		t.Fatalf("expected parse warning, got: %s", buf.String())
	}
}

func TestLoadConfigInvalidKey(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{"keybindings": {"BadKeyName": "quit"}}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	buf := captureLog(t)
	cfg := loadConfig()
	if len(cfg.Keybindings) != 2 {
		t.Fatalf("expected 2 default keybindings, got %d", len(cfg.Keybindings))
	}
	if !strings.Contains(buf.String(), "invalid key") {
		t.Fatalf("expected invalid key warning, got: %s", buf.String())
	}
}

func TestLoadConfigUnknownAction(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{"keybindings": {"CtrlQ": "fly"}}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	buf := captureLog(t)
	cfg := loadConfig()
	if len(cfg.Keybindings) != 2 {
		t.Fatalf("expected 2 default keybindings, got %d", len(cfg.Keybindings))
	}
	if !strings.Contains(buf.String(), "unknown action") {
		t.Fatalf("expected unknown action warning, got: %s", buf.String())
	}
}

func TestLoadConfigNilKeybindings(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig()
	if len(cfg.Keybindings) != 2 {
		t.Fatalf("expected 2 default bindings with empty config, got %d", len(cfg.Keybindings))
	}
}

func TestLoadConfigUnreadableLocalFile(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(dir, configFileName), 0755); err != nil {
		t.Fatal(err)
	}

	buf := captureLog(t)
	cfg := loadConfig()
	if len(cfg.Keybindings) != 2 {
		t.Fatalf("expected 2 default keybindings, got %d", len(cfg.Keybindings))
	}
	if !strings.Contains(buf.String(), "could not read") {
		t.Fatalf("expected could not read warning, got: %s", buf.String())
	}
}

func TestLoadConfigGlobalFallback(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	globalDir := filepath.Join(configDir, "et")
	os.MkdirAll(globalDir, 0755)
	if err := os.WriteFile(filepath.Join(globalDir, configFileName), []byte(`{"keybindings": {"CtrlQ": "quit"}}`), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(globalDir)

	cfg := loadConfig()
	if len(cfg.Keybindings) != 1 {
		t.Fatalf("expected 1 binding from global fallback, got %d", len(cfg.Keybindings))
	}
	if _, ok := cfg.Keybindings["CtrlQ"]; !ok {
		t.Fatal("global CtrlQ binding should be present")
	}
}

// parseColor tests

func TestParseColorNamed(t *testing.T) {
	c, err := parseColor("darkcyan")
	if err != nil {
		t.Fatal(err)
	}
	if c != color.DarkCyan {
		t.Fatalf("expected DarkCyan, got %v", c)
	}
}

func TestParseColorHex(t *testing.T) {
	c, err := parseColor("#ff0000")
	if err != nil {
		t.Fatal(err)
	}
	if !c.IsRGB() {
		t.Fatalf("expected RGB color, got %v", c)
	}
}

func TestParseColorEmpty(t *testing.T) {
	c, err := parseColor("")
	if err != nil {
		t.Fatal(err)
	}
	if c != color.Default {
		t.Fatalf("expected Default, got %v", c)
	}
}

func TestParseColorUnknown(t *testing.T) {
	_, err := parseColor("blarg")
	if err == nil {
		t.Fatal("expected error for unknown color")
	}
}

func TestParseColorBlackWhite(t *testing.T) {
	c, err := parseColor("black")
	if err != nil {
		t.Fatal(err)
	}
	if c != color.Black {
		t.Fatalf("expected Black, got %v", c)
	}
	c, err = parseColor("white")
	if err != nil {
		t.Fatal(err)
	}
	if c != color.White {
		t.Fatalf("expected White, got %v", c)
	}
}

// loadConfig color tests

func TestLoadConfigWithColors(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{"colors": {"foreground": "red", "background": "black", "status": "blue"}, "keybindings": {"CtrlQ": "quit"}}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig()
	if cfg.Colors.Foreground != "red" {
		t.Fatalf("expected foreground red, got %q", cfg.Colors.Foreground)
	}
	if cfg.Colors.Background != "black" {
		t.Fatalf("expected background black, got %q", cfg.Colors.Background)
	}
	if cfg.Colors.StatusBG != "blue" {
		t.Fatalf("expected status blue, got %q", cfg.Colors.StatusBG)
	}
}

func TestLoadConfigInvalidColor(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{"colors": {"foreground": "hotdog"}}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	buf := captureLog(t)
	cfg := loadConfig()
	if cfg.Colors.Foreground != "" {
		t.Fatalf("expected empty foreground for invalid color, got %q", cfg.Colors.Foreground)
	}
	if !strings.Contains(buf.String(), "unknown color") {
		t.Fatalf("expected unknown color warning, got: %s", buf.String())
	}
}

func TestLoadConfigNoColorsUsesDefaults(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig()
	if cfg.Colors.Foreground != "" {
		t.Fatalf("expected empty foreground, got %q", cfg.Colors.Foreground)
	}
	if cfg.Colors.Background != "" {
		t.Fatalf("expected empty background, got %q", cfg.Colors.Background)
	}
	if cfg.Colors.StatusBG != "" {
		t.Fatalf("expected empty status, got %q", cfg.Colors.StatusBG)
	}
}

func TestLoadConfigBothPathsUnreadable(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(dir, configFileName), 0755); err != nil {
		t.Fatal(err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	globalDir := filepath.Join(configDir, "et")
	os.MkdirAll(globalDir, 0755)
	if err := os.RemoveAll(filepath.Join(globalDir, configFileName)); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(globalDir, configFileName), 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(globalDir)

	buf := captureLog(t)
	cfg := loadConfig()
	if len(cfg.Keybindings) != 2 {
		t.Fatalf("expected 2 default keybindings, got %d", len(cfg.Keybindings))
	}
	if !strings.Contains(buf.String(), "could not read") {
		t.Fatalf("expected could not read warnings, got: %s", buf.String())
	}
}

func TestLoadConfigLocalTakesPriority(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(`{"keybindings": {"Escape": "quit"}}`), 0644); err != nil {
		t.Fatal(err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	globalDir := filepath.Join(configDir, "et")
	os.MkdirAll(globalDir, 0755)
	if err := os.WriteFile(filepath.Join(globalDir, configFileName), []byte(`{"keybindings": {"CtrlQ": "quit", "Enter": "quit"}}`), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(globalDir)

	cfg := loadConfig()
	if len(cfg.Keybindings) != 1 {
		t.Fatalf("expected 1 binding from local, got %d", len(cfg.Keybindings))
	}
	if _, ok := cfg.Keybindings["Escape"]; !ok {
		t.Fatal("local Escape binding should be present")
	}
}

func TestDefaultTabWidth(t *testing.T) {
	if defaults.TabWidth() != 4 {
		t.Fatalf("expected default tab width 4, got %d", defaults.TabWidth())
	}
}

func TestLoadConfigTabWidth(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{"tab_width": 8}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig()
	if cfg.TabWidth != 8 {
		t.Fatalf("expected tab width 8, got %d", cfg.TabWidth)
	}
}

func TestLoadConfigInvalidTabWidth(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prev)
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	content := `{"tab_width": -1}`
	if err := os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	buf := captureLog(t)
	cfg := loadConfig()
	if cfg.TabWidth != 4 {
		t.Fatalf("expected tab width 4 default, got %d", cfg.TabWidth)
	}
	if !strings.Contains(buf.String(), "invalid tab_width") {
		t.Fatalf("expected invalid tab_width warning, got: %s", buf.String())
	}
}

func TestExpandTabs(t *testing.T) {
	tests := []struct {
		input    string
		tabWidth int
		want     string
	}{
		{"", 4, ""},
		{"hello", 4, "hello"},
		{"\t", 4, "    "},
		{"\t", 2, "  "},
		{"a\tb", 4, "a   b"},
		{"aa\tb", 4, "aa  b"},
		{"aaa\tb", 4, "aaa b"},
		{"aaaa\tb", 4, "aaaa    b"},
		{"\t\t", 4, "        "},
		{"a\tb\tc", 2, "a b c"},
	}
	for _, tt := range tests {
		got := expandTabs(tt.input, tt.tabWidth)
		if got != tt.want {
			t.Errorf("expandTabs(%q, %d) = %q, want %q", tt.input, tt.tabWidth, got, tt.want)
		}
	}
}
