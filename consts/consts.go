package consts

const (
	LogFileName         = "et.log"
	Version             = "0.0.1"
	welcomeMessageLine1 = "et - (e)dit (t)ext"
	welcomeMessageLine2 = "version " + Version
	welcomeMessageLine3 = "by Brice Vadnais"
)

var WelcomeMessages = [3]string{
	welcomeMessageLine1,
	welcomeMessageLine2,
	welcomeMessageLine3,
}
