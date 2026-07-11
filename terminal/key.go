package terminal

import (
	"strings"

	"github.com/gdamore/tcell/v3"
)

func keyCode(ev *tcell.EventKey) string {
	key := ev.Key()
	mod := ev.Modifiers()

	if mod == tcell.ModNone {
		if key == tcell.KeyRune {
			return ev.Str()
		}
		if s, ok := keyCodes[key]; ok {
			return s
		}
		return string(rune(key))
	}

	if mod == tcell.ModAlt && key == tcell.KeyRune {
		return "\x1b" + ev.Str()
	}

	if mod == tcell.ModCtrl && key >= tcell.KeyCtrlA && key <= tcell.KeyCtrlZ {
		return string(byte(key - tcell.KeyCtrlA + 1))
	}

	if mod&tcell.ModMeta != 0 {
		return metaKeyCode(ev)
	}

	if s, ok := modLookup[mod][key]; ok {
		return s
	}
	return ""
}

func metaKeyCode(ev *tcell.EventKey) string {
	paramNum, ok := metaParamNums[ev.Modifiers()]
	if !ok {
		return ""
	}
	switch ev.Key() {
	case tcell.KeyUp, tcell.KeyDown, tcell.KeyRight, tcell.KeyLeft,
		tcell.KeyHome, tcell.KeyEnd, tcell.KeyInsert, tcell.KeyDelete,
		tcell.KeyPgUp, tcell.KeyPgDn,
		tcell.KeyF1, tcell.KeyF2, tcell.KeyF3, tcell.KeyF4,
		tcell.KeyF5, tcell.KeyF6, tcell.KeyF7, tcell.KeyF8,
		tcell.KeyF9, tcell.KeyF10, tcell.KeyF11, tcell.KeyF12:
	default:
		return ""
	}
	kc := keyCodes[ev.Key()]
	switch {
	case strings.HasSuffix(kc, "~"):
		return strings.TrimSuffix(kc, "~") + ";" + paramNum + "~"
	default:
		return "\x1b[1;" + paramNum + strings.TrimPrefix(kc, "\x1bO")
	}
}

type modMap map[tcell.Key]string

var metaParamNums = map[tcell.ModMask]string{
	tcell.ModMeta:                                                 "9",
	tcell.ModMeta | tcell.ModShift:                                "10",
	tcell.ModMeta | tcell.ModAlt:                                  "11",
	tcell.ModMeta | tcell.ModAlt | tcell.ModShift:                 "12",
	tcell.ModMeta | tcell.ModCtrl:                                 "13",
	tcell.ModMeta | tcell.ModCtrl | tcell.ModShift:                "14",
	tcell.ModMeta | tcell.ModCtrl | tcell.ModAlt:                  "15",
	tcell.ModMeta | tcell.ModCtrl | tcell.ModAlt | tcell.ModShift: "16",
}

var modLookup = map[tcell.ModMask]modMap{
	tcell.ModShift: {
		tcell.KeyUp:     info.KeyShfUp,
		tcell.KeyDown:   info.KeyShfDown,
		tcell.KeyRight:  info.KeyShfRight,
		tcell.KeyLeft:   info.KeyShfLeft,
		tcell.KeyHome:   info.KeyShfHome,
		tcell.KeyEnd:    info.KeyShfEnd,
		tcell.KeyInsert: info.KeyShfInsert,
		tcell.KeyDelete: info.KeyShfDelete,
		tcell.KeyPgUp:   info.KeyShfPgUp,
		tcell.KeyPgDn:   info.KeyShfPgDn,
		tcell.KeyF1:     info.KeyF13,
		tcell.KeyF2:     info.KeyF14,
		tcell.KeyF3:     info.KeyF15,
		tcell.KeyF4:     info.KeyF16,
		tcell.KeyF5:     info.KeyF17,
		tcell.KeyF6:     info.KeyF18,
		tcell.KeyF7:     info.KeyF19,
		tcell.KeyF8:     info.KeyF20,
		tcell.KeyF9:     info.KeyF21,
		tcell.KeyF10:    info.KeyF22,
		tcell.KeyF11:    info.KeyF23,
		tcell.KeyF12:    info.KeyF24,
	},
	tcell.ModAlt: {
		tcell.KeyUp:     info.KeyAltUp,
		tcell.KeyDown:   info.KeyAltDown,
		tcell.KeyRight:  info.KeyAltRight,
		tcell.KeyLeft:   info.KeyAltLeft,
		tcell.KeyHome:   info.KeyAltHome,
		tcell.KeyEnd:    info.KeyAltEnd,
		tcell.KeyInsert: extendedInfo.KeyAltInsert,
		tcell.KeyDelete: extendedInfo.KeyAltDelete,
		tcell.KeyPgUp:   extendedInfo.KeyAltPgUp,
		tcell.KeyPgDn:   extendedInfo.KeyAltPgDown,
		tcell.KeyF1:     info.KeyF49,
		tcell.KeyF2:     info.KeyF50,
		tcell.KeyF3:     info.KeyF51,
		tcell.KeyF4:     info.KeyF53,
		tcell.KeyF5:     info.KeyF54,
		tcell.KeyF6:     info.KeyF55,
		tcell.KeyF7:     info.KeyF56,
		tcell.KeyF8:     info.KeyF57,
		tcell.KeyF9:     info.KeyF58,
		tcell.KeyF10:    info.KeyF59,
		tcell.KeyF11:    info.KeyF60,
		tcell.KeyF12:    info.KeyF61,
	},
	tcell.ModCtrl: {
		tcell.KeyUp:     info.KeyCtrlUp,
		tcell.KeyDown:   info.KeyCtrlDown,
		tcell.KeyRight:  info.KeyCtrlRight,
		tcell.KeyLeft:   info.KeyCtrlLeft,
		tcell.KeyHome:   info.KeyCtrlHome,
		tcell.KeyEnd:    info.KeyCtrlEnd,
		tcell.KeyInsert: extendedInfo.KeyCtrlInsert,
		tcell.KeyDelete: extendedInfo.KeyCtrlDelete,
		tcell.KeyPgUp:   extendedInfo.KeyCtrlPgUp,
		tcell.KeyPgDn:   extendedInfo.KeyCtrlPgDown,
		tcell.KeyF1:     info.KeyF25,
		tcell.KeyF2:     info.KeyF26,
		tcell.KeyF3:     info.KeyF27,
		tcell.KeyF4:     info.KeyF28,
		tcell.KeyF5:     info.KeyF29,
		tcell.KeyF6:     info.KeyF30,
		tcell.KeyF7:     info.KeyF31,
		tcell.KeyF8:     info.KeyF32,
		tcell.KeyF9:     info.KeyF33,
		tcell.KeyF10:    info.KeyF34,
		tcell.KeyF11:    info.KeyF35,
		tcell.KeyF12:    info.KeyF36,
	},
	tcell.ModCtrl | tcell.ModShift: {
		tcell.KeyUp:     info.KeyCtrlShfUp,
		tcell.KeyDown:   info.KeyCtrlShfDown,
		tcell.KeyRight:  info.KeyCtrlShfRight,
		tcell.KeyLeft:   info.KeyCtrlShfLeft,
		tcell.KeyHome:   info.KeyCtrlShfHome,
		tcell.KeyEnd:    info.KeyCtrlShfEnd,
		tcell.KeyInsert: extendedInfo.KeyCtrlShfInsert,
		tcell.KeyDelete: extendedInfo.KeyCtrlShfDelete,
		tcell.KeyPgUp:   extendedInfo.KeyCtrlShfPgUp,
		tcell.KeyPgDn:   extendedInfo.KeyCtrlShfPgDown,
		tcell.KeyF1:     info.KeyF37,
		tcell.KeyF2:     info.KeyF38,
		tcell.KeyF3:     info.KeyF39,
		tcell.KeyF4:     info.KeyF40,
		tcell.KeyF5:     info.KeyF41,
		tcell.KeyF6:     info.KeyF42,
		tcell.KeyF7:     info.KeyF43,
		tcell.KeyF8:     info.KeyF44,
		tcell.KeyF9:     info.KeyF45,
		tcell.KeyF10:    info.KeyF46,
		tcell.KeyF11:    info.KeyF47,
		tcell.KeyF12:    info.KeyF48,
	},
	tcell.ModAlt | tcell.ModShift: {
		tcell.KeyUp:     info.KeyAltShfUp,
		tcell.KeyDown:   info.KeyAltShfDown,
		tcell.KeyRight:  info.KeyAltShfRight,
		tcell.KeyLeft:   info.KeyAltShfLeft,
		tcell.KeyHome:   info.KeyAltShfHome,
		tcell.KeyEnd:    info.KeyAltShfEnd,
		tcell.KeyInsert: extendedInfo.KeyAltShfInsert,
		tcell.KeyDelete: extendedInfo.KeyAltShfDelete,
		tcell.KeyPgUp:   extendedInfo.KeyAltShfPgUp,
		tcell.KeyPgDn:   extendedInfo.KeyAltShfPgDown,
		tcell.KeyF1:     info.KeyF61,
		tcell.KeyF2:     info.KeyF62,
		tcell.KeyF3:     info.KeyF63,
		tcell.KeyF4:     info.KeyF64,
	},
	tcell.ModAlt | tcell.ModCtrl: {
		tcell.KeyUp:     extendedInfo.KeyCtrlAltUp,
		tcell.KeyDown:   extendedInfo.KeyCtrlAltDown,
		tcell.KeyRight:  extendedInfo.KeyCtrlAltRight,
		tcell.KeyLeft:   extendedInfo.KeyCtrlAltLeft,
		tcell.KeyHome:   extendedInfo.KeyCtrlAltHome,
		tcell.KeyEnd:    extendedInfo.KeyCtrlAltEnd,
		tcell.KeyInsert: extendedInfo.KeyCtrlAltInsert,
		tcell.KeyDelete: extendedInfo.KeyCtrlAltDelete,
		tcell.KeyPgUp:   extendedInfo.KeyCtrlAltPgUp,
		tcell.KeyPgDn:   extendedInfo.KeyCtrlAltPgDown,
	},
	tcell.ModAlt | tcell.ModCtrl | tcell.ModShift: {
		tcell.KeyUp:     extendedInfo.KeyCtrlAltShfUp,
		tcell.KeyDown:   extendedInfo.KeyCtrlAltShfDown,
		tcell.KeyRight:  extendedInfo.KeyCtrlAltShfRight,
		tcell.KeyLeft:   extendedInfo.KeyCtrlAltShfLeft,
		tcell.KeyHome:   extendedInfo.KeyCtrlAltShfHome,
		tcell.KeyEnd:    extendedInfo.KeyCtrlAltShfEnd,
		tcell.KeyInsert: extendedInfo.KeyCtrlAltShfInsert,
		tcell.KeyDelete: extendedInfo.KeyCtrlAltShfDelete,
		tcell.KeyPgUp:   extendedInfo.KeyCtrlAltShfPgUp,
		tcell.KeyPgDn:   extendedInfo.KeyCtrlAltShfPgDown,
	},
}

var keyCodes = map[tcell.Key]string{
	tcell.KeyBackspace: info.KeyBackspace,
	tcell.KeyF1:        info.KeyF1,
	tcell.KeyF2:        info.KeyF2,
	tcell.KeyF3:        info.KeyF3,
	tcell.KeyF4:        info.KeyF4,
	tcell.KeyF5:        info.KeyF5,
	tcell.KeyF6:        info.KeyF6,
	tcell.KeyF7:        info.KeyF7,
	tcell.KeyF8:        info.KeyF8,
	tcell.KeyF9:        info.KeyF9,
	tcell.KeyF10:       info.KeyF10,
	tcell.KeyF11:       info.KeyF11,
	tcell.KeyF12:       info.KeyF12,
	tcell.KeyF13:       info.KeyF13,
	tcell.KeyF14:       info.KeyF14,
	tcell.KeyF15:       info.KeyF15,
	tcell.KeyF16:       info.KeyF16,
	tcell.KeyF17:       info.KeyF17,
	tcell.KeyF18:       info.KeyF18,
	tcell.KeyF19:       info.KeyF19,
	tcell.KeyF20:       info.KeyF20,
	tcell.KeyF21:       info.KeyF21,
	tcell.KeyF22:       info.KeyF22,
	tcell.KeyF23:       info.KeyF23,
	tcell.KeyF24:       info.KeyF24,
	tcell.KeyF25:       info.KeyF25,
	tcell.KeyF26:       info.KeyF26,
	tcell.KeyF27:       info.KeyF27,
	tcell.KeyF28:       info.KeyF28,
	tcell.KeyF29:       info.KeyF29,
	tcell.KeyF30:       info.KeyF30,
	tcell.KeyF31:       info.KeyF31,
	tcell.KeyF32:       info.KeyF32,
	tcell.KeyF33:       info.KeyF33,
	tcell.KeyF34:       info.KeyF34,
	tcell.KeyF35:       info.KeyF35,
	tcell.KeyF36:       info.KeyF36,
	tcell.KeyF37:       info.KeyF37,
	tcell.KeyF38:       info.KeyF38,
	tcell.KeyF39:       info.KeyF39,
	tcell.KeyF40:       info.KeyF40,
	tcell.KeyF41:       info.KeyF41,
	tcell.KeyF42:       info.KeyF42,
	tcell.KeyF43:       info.KeyF43,
	tcell.KeyF44:       info.KeyF44,
	tcell.KeyF45:       info.KeyF45,
	tcell.KeyF46:       info.KeyF46,
	tcell.KeyF47:       info.KeyF47,
	tcell.KeyF48:       info.KeyF48,
	tcell.KeyF49:       info.KeyF49,
	tcell.KeyF50:       info.KeyF50,
	tcell.KeyF51:       info.KeyF51,
	tcell.KeyF52:       info.KeyF52,
	tcell.KeyF53:       info.KeyF53,
	tcell.KeyF54:       info.KeyF54,
	tcell.KeyF55:       info.KeyF55,
	tcell.KeyF56:       info.KeyF56,
	tcell.KeyF57:       info.KeyF57,
	tcell.KeyF58:       info.KeyF58,
	tcell.KeyF59:       info.KeyF59,
	tcell.KeyF60:       info.KeyF60,
	tcell.KeyF61:       info.KeyF61,
	tcell.KeyF62:       info.KeyF62,
	tcell.KeyF63:       info.KeyF63,
	tcell.KeyF64:       info.KeyF64,
	tcell.KeyInsert:    info.KeyInsert,
	tcell.KeyDelete:    info.KeyDelete,
	tcell.KeyHome:      info.KeyHome,
	tcell.KeyEnd:       info.KeyEnd,
	tcell.KeyHelp:      info.KeyHelp,
	tcell.KeyPgUp:      info.KeyPgUp,
	tcell.KeyPgDn:      info.KeyPgDn,
	tcell.KeyUp:        info.KeyUp,
	tcell.KeyDown:      info.KeyDown,
	tcell.KeyLeft:      info.KeyLeft,
	tcell.KeyRight:     info.KeyRight,
	tcell.KeyBacktab:   info.KeyBacktab,
	tcell.KeyExit:      info.KeyExit,
	tcell.KeyClear:     info.KeyClear,
	tcell.KeyPrint:     info.KeyPrint,
	tcell.KeyCancel:    info.KeyCancel,
}
