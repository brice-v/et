package main

import "github.com/gdamore/tcell/v3"

type keySpec struct {
	key  tcell.Key
	str  string
	mods tcell.ModMask
}

func (ks keySpec) matches(e *tcell.EventKey) bool {
	if ks.key == tcell.KeyRune {
		return e.Key() == tcell.KeyRune && e.Str() == ks.str && e.Modifiers() == ks.mods
	}
	return e.Key() == ks.key
}
