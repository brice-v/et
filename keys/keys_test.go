package keys

import (
	"github.com/brice-v/et/config"
	"testing"

	"github.com/gdamore/tcell/v3"
)

func TestNormalizeKeyCtrlRange(t *testing.T) {
	tests := []struct {
		name      string
		key       tcell.Key
		keyAsRune string
		mod       tcell.ModMask
		wantKey   tcell.Key
		wantMod   tcell.ModMask
	}{
		{"CtrlA", tcell.KeyCtrlA, "", tcell.ModNone, tcell.Key('a'), tcell.ModCtrl},
		{"CtrlZ", tcell.KeyCtrlZ, "", tcell.ModNone, tcell.Key('z'), tcell.ModCtrl},
		{"CtrlM", tcell.KeyCtrlM, "", tcell.ModNone, tcell.Key('m'), tcell.ModCtrl},
		{"CtrlQ", tcell.KeyCtrlQ, "", tcell.ModNone, tcell.Key('q'), tcell.ModCtrl},
		{"CtrlShiftQ", tcell.KeyCtrlQ, "", tcell.ModShift, tcell.Key('q'), tcell.ModCtrl | tcell.ModShift},
		{"CtrlAltQ", tcell.KeyCtrlQ, "", tcell.ModAlt, tcell.Key('q'), tcell.ModCtrl | tcell.ModAlt},
		{"CtrlShiftAltQ", tcell.KeyCtrlQ, "", tcell.ModShift | tcell.ModAlt, tcell.Key('q'), tcell.ModCtrl | tcell.ModShift | tcell.ModAlt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotMod := NormalizeKey(tt.key, tt.keyAsRune, tt.mod)
			if gotKey != tt.wantKey {
				t.Errorf("key = %v, want %v", gotKey, tt.wantKey)
			}
			if gotMod != tt.wantMod {
				t.Errorf("mod = %v, want %v", gotMod, tt.wantMod)
			}
		})
	}
}

func TestNormalizeKeyRune(t *testing.T) {
	tests := []struct {
		name      string
		key       tcell.Key
		keyAsRune string
		mod       tcell.ModMask
		wantKey   tcell.Key
		wantMod   tcell.ModMask
	}{
		{"lowercase_q", tcell.KeyRune, "q", tcell.ModNone, tcell.Key('q'), tcell.ModNone},
		{"uppercase_Q", tcell.KeyRune, "Q", tcell.ModNone, tcell.Key('q'), tcell.ModNone},
		{"ctrl_q", tcell.KeyRune, "q", tcell.ModCtrl, tcell.Key('q'), tcell.ModCtrl},
		{"ctrl_shift_Q", tcell.KeyRune, "Q", tcell.ModCtrl | tcell.ModShift, tcell.Key('q'), tcell.ModCtrl | tcell.ModShift},
		{"space", tcell.KeyRune, " ", tcell.ModNone, tcell.Key(' '), tcell.ModNone},
		{"digit_1", tcell.KeyRune, "1", tcell.ModNone, tcell.Key('1'), tcell.ModNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotMod := NormalizeKey(tt.key, tt.keyAsRune, tt.mod)
			if gotKey != tt.wantKey {
				t.Errorf("key = %v, want %v", gotKey, tt.wantKey)
			}
			if gotMod != tt.wantMod {
				t.Errorf("mod = %v, want %v", gotMod, tt.wantMod)
			}
		})
	}
}

func TestNormalizeKeyPassthrough(t *testing.T) {
	tests := []struct {
		name      string
		key       tcell.Key
		keyAsRune string
		mod       tcell.ModMask
		wantKey   tcell.Key
		wantMod   tcell.ModMask
	}{
		{"escape", tcell.KeyEscape, "", tcell.ModNone, tcell.KeyEscape, tcell.ModNone},
		{"escape_with_shift", tcell.KeyEscape, "", tcell.ModShift, tcell.KeyEscape, tcell.ModShift},
		{"enter", tcell.KeyEnter, "", tcell.ModNone, tcell.KeyEnter, tcell.ModNone},
		{"tab", tcell.KeyTab, "", tcell.ModNone, tcell.KeyTab, tcell.ModNone},
		{"backspace", tcell.KeyBackspace, "", tcell.ModNone, tcell.KeyBackspace, tcell.ModNone},
		{"rune_multi_char", tcell.KeyRune, "ab", tcell.ModNone, tcell.KeyRune, tcell.ModNone},
		{"rune_empty", tcell.KeyRune, "", tcell.ModNone, tcell.KeyRune, tcell.ModNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotMod := NormalizeKey(tt.key, tt.keyAsRune, tt.mod)
			if gotKey != tt.wantKey {
				t.Errorf("key = %v, want %v", gotKey, tt.wantKey)
			}
			if gotMod != tt.wantMod {
				t.Errorf("mod = %v, want %v", gotMod, tt.wantMod)
			}
		})
	}
}

func TestIsKeyAny(t *testing.T) {
	bindings := []config.Key{
		{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl | tcell.ModShift},
		{Key: tcell.KeyQ, Modifiers: tcell.ModCtrl},
		{Key: tcell.KeyQ, Modifiers: tcell.ModNone},
		{Key: tcell.KeyEscape, Modifiers: tcell.ModNone},
	}

	tests := []struct {
		name      string
		key       tcell.Key
		keyAsRune string
		mod       tcell.ModMask
		want      bool
	}{
		{"ctrl_shift_q_legacy", tcell.KeyCtrlQ, "", tcell.ModShift, true},
		{"ctrl_q_legacy", tcell.KeyCtrlQ, "", tcell.ModNone, true},
		{"q_rune", tcell.KeyRune, "q", tcell.ModNone, true},
		{"q_legacy", tcell.KeyQ, "", tcell.ModNone, true},
		{"ctrl_q_advanced", tcell.KeyRune, "q", tcell.ModCtrl, true},
		{"ctrl_shift_q_advanced", tcell.KeyRune, "Q", tcell.ModCtrl | tcell.ModShift, true},
		{"esc_legacy", tcell.KeyEscape, "", tcell.ModNone, true},
		{"ctrl_alt_q_no_binding", tcell.KeyRune, "q", tcell.ModCtrl | tcell.ModAlt, false},
		{"shift_q_no_binding", tcell.KeyRune, "Q", tcell.ModShift, false},
		{"alt_q_no_binding", tcell.KeyRune, "q", tcell.ModAlt, false},
		{"enter_no_binding", tcell.KeyEnter, "", tcell.ModNone, false},
		{"ctrl_a_no_binding", tcell.KeyCtrlA, "", tcell.ModNone, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsKeyAny(tt.key, tt.keyAsRune, tt.mod, bindings)
			if got != tt.want {
				t.Errorf("IsKeyAny(%v, %q, %v) = %v, want %v", tt.key, tt.keyAsRune, tt.mod, got, tt.want)
			}
		})
	}
}

func TestIsKeySemicolonColon(t *testing.T) {
	tests := []struct {
		name      string
		key       tcell.Key
		keyAsRune string
		mod       tcell.ModMask
		cfgKey    config.Key
		want      bool
	}{
		// Terminal sends base keycode with mods (kitty protocol)
		{"base ; match ;", tcell.Key(';'), "", tcell.ModCtrl | tcell.ModShift, config.Key{Key: tcell.Key(';'), Modifiers: tcell.ModCtrl | tcell.ModShift}, true},
		// Terminal sends shifted char as rune
		{"shifted : match ;", tcell.KeyRune, ":", tcell.ModCtrl | tcell.ModShift, config.Key{Key: tcell.Key(';'), Modifiers: tcell.ModCtrl | tcell.ModShift}, false},
		// Cross-match fails (IsKey doesn't do ;/: swapping)
		{"base ; match :", tcell.Key(';'), "", tcell.ModCtrl | tcell.ModShift, config.Key{Key: tcell.Key(':'), Modifiers: tcell.ModCtrl | tcell.ModShift}, false},
		// Terminal sends shifted char, config uses :
		{"shifted : match :", tcell.KeyRune, ":", tcell.ModCtrl | tcell.ModShift, config.Key{Key: tcell.Key(':'), Modifiers: tcell.ModCtrl | tcell.ModShift}, true},
		// ctrl+; without shift
		{"ctrl ; no shift", tcell.KeyRune, ";", tcell.ModCtrl, config.Key{Key: tcell.Key(';'), Modifiers: tcell.ModCtrl}, true},
		// ctrl+t (default toggle)
		{"ctrl t legacy", tcell.KeyCtrlT, "", tcell.ModNone, config.Key{Key: tcell.Key('t'), Modifiers: tcell.ModCtrl}, true},
		{"ctrl t advanced", tcell.KeyRune, "t", tcell.ModCtrl, config.Key{Key: tcell.Key('t'), Modifiers: tcell.ModCtrl}, true},
		{"ctrl T advanced", tcell.KeyRune, "T", tcell.ModCtrl, config.Key{Key: tcell.Key('t'), Modifiers: tcell.ModCtrl}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsKey(tt.key, tt.keyAsRune, tt.mod, tt.cfgKey)
			if got != tt.want {
				t.Errorf("IsKey(%v, %q, %v, %v) = %v, want %v", tt.key, tt.keyAsRune, tt.mod, tt.cfgKey, got, tt.want)
			}
		})
	}
}

func TestNormalizeShiftTab(t *testing.T) {
	tests := []struct {
		name      string
		key       tcell.Key
		keyAsRune string
		mod       tcell.ModMask
		wantKey   tcell.Key
		wantMod   tcell.ModMask
	}{
		{"legacy backtab", tcell.KeyBacktab, "", tcell.ModNone, tcell.KeyBacktab, tcell.ModNone},
		{"advanced shift+tab", tcell.KeyTab, "", tcell.ModShift, tcell.KeyBacktab, tcell.ModNone},
		{"plain tab", tcell.KeyTab, "", tcell.ModNone, tcell.KeyTab, tcell.ModNone},
		{"ctrl+tab", tcell.KeyTab, "", tcell.ModCtrl, tcell.KeyTab, tcell.ModCtrl},
		{"shift+ctrl+tab", tcell.KeyTab, "", tcell.ModShift | tcell.ModCtrl, tcell.KeyBacktab, tcell.ModCtrl},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotMod := NormalizeKey(tt.key, tt.keyAsRune, tt.mod)
			if gotKey != tt.wantKey {
				t.Errorf("key = %v, want %v", gotKey, tt.wantKey)
			}
			if gotMod != tt.wantMod {
				t.Errorf("mod = %v, want %v", gotMod, tt.wantMod)
			}
		})
	}
}

func TestIsKeyAnyEmptyBindings(t *testing.T) {
	if IsKeyAny(tcell.KeyQ, "", tcell.ModNone, nil) {
		t.Error("IsKeyAny with nil bindings should be false")
	}
	if IsKeyAny(tcell.KeyQ, "", tcell.ModNone, []config.Key{}) {
		t.Error("IsKeyAny with empty bindings should be false")
	}
}
