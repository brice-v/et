package consts

import (
	"math"
)

const (
	LogFileName         = "et.log"
	Version             = "0.0.5"
	StickyColMax        = math.MaxInt64
	welcomeMessageLine1 = "et - (e)dit (t)ext"
	welcomeMessageLine2 = "version " + Version
	welcomeMessageLine3 = "by Brice Vadnais"
)

var WelcomeMessages = [3]string{
	welcomeMessageLine1,
	welcomeMessageLine2,
	welcomeMessageLine3,
}

type HlStyleType int

const (
	HlBase HlStyleType = iota
	Hl1
	Hl2
	Hl3
	HlStr
	HlSpc
	HlCom
)
