package keys

import (
	"et/config"
	"unicode"

	"github.com/gdamore/tcell/v3"
)

func NormalizeKey(key tcell.Key, keyAsRune string, mod tcell.ModMask) (tcell.Key, tcell.ModMask) {
	if key >= tcell.KeyCtrlA && key <= tcell.KeyCtrlZ {
		return tcell.Key('a' + int(key-tcell.KeyCtrlA)), mod | tcell.ModCtrl
	}
	if key == tcell.KeyRune && len(keyAsRune) == 1 {
		r := []rune(keyAsRune)[0]
		return tcell.Key(unicode.ToLower(r)), mod
	}
	return key, mod
}

func IsKeyAny(key tcell.Key, keyAsRune string, mod tcell.ModMask, keys []config.Key) bool {
	nk, nm := NormalizeKey(key, keyAsRune, mod)
	for _, k := range keys {
		if nk == k.Key && nm == k.Modifiers {
			return true
		}
	}
	return false
}
