package editor

import (
	"et/config"
	"et/consts"
	"log/slog"
	"os"
	"strings"

	"github.com/gdamore/tcell/v3"
)

type Editor struct {
	s tcell.Screen
	// sw, sh screen width and height, calculated every Draw
	sw, sh int
	// lPad is the padding needed for line numbers or tilde on the left
	lPad int
	// sbh is the status bar height (always drawn at bottom but including
	//  this so it can be more than 1 high)
	// (defaults to 1)
	sbh        int
	baseStyle  tcell.Style
	hl1Style   tcell.Style
	hl2Style   tcell.Style
	hl3Style   tcell.Style
	hlStrStyle tcell.Style
	hlSpcStyle tcell.Style
	hlComStyle tcell.Style
	// highlight database, the keywords/tokens mapped to the corresponding color
	hldb        map[string]consts.HlStyleType
	hlOperators string

	commentToken string
	stringTokens []string

	// cx, cy cursor x and y position
	cx, cy int
	// vScrollOffset is the first visible line in the viewport
	vScrollOffset int
	// hScrollOffset is the first visible column in the viewport
	hScrollOffset int
	// stickyCol is the file column for vertical movement that gets "stuck"
	stickyCol int

	cfg *config.Config

	fileName         string
	fileContentLines [][]rune
	fileExtension    string

	// Exit is a flag to trigger exit
	Exit bool
}

func New(s tcell.Screen, cfg *config.Config, fileName string) *Editor {
	fcl := getFileContent(fileName)
	baseStyle := tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(cfg.Colors.Foreground.Color)
	splitFilename := strings.Split(fileName, ".")
	fileExtension := ""
	if len(splitFilename) > 0 {
		fileExtension = splitFilename[len(splitFilename)-1]
	}
	e := &Editor{
		s:                s,
		sbh:              1,
		baseStyle:        baseStyle,
		cfg:              cfg,
		fileName:         fileName,
		fileContentLines: fcl,
		fileExtension:    fileExtension,
	}
	e.setupHlStyles()
	return e
}

func (e *Editor) setupHlStyles() {
	if e.fileExtension == "" {
		return
	}
	fileType, ok := e.cfg.FileExtensions[e.fileExtension]
	if !ok {
		slog.Warn("no syntax highlighting available for file extension", "fileExtension", e.fileExtension)
		return
	}
	colorMap, ok := e.cfg.Colors.Languages[fileType]
	if !ok {
		slog.Warn("no color map for highlighting found for fileType", "fileType", fileType)
		return
	}
	e.hl1Style = tcell.StyleDefault.Background(e.cfg.Colors.Background.Color).Foreground(colorMap.Color1.Color)
	e.hl2Style = tcell.StyleDefault.Background(e.cfg.Colors.Background.Color).Foreground(colorMap.Color2.Color)
	e.hl3Style = tcell.StyleDefault.Background(e.cfg.Colors.Background.Color).Foreground(colorMap.Color3.Color)
	e.hlStrStyle = tcell.StyleDefault.Background(e.cfg.Colors.Background.Color).Foreground(colorMap.ColorString.Color)
	e.hlSpcStyle = tcell.StyleDefault.Background(e.cfg.Colors.Background.Color).Foreground(colorMap.SpecialColor.Color)
	e.hlComStyle = tcell.StyleDefault.Background(e.cfg.Colors.Background.Color).Foreground(colorMap.CommentColor.Color)
	e.hldb = make(map[string]consts.HlStyleType)
	for _, kw := range colorMap.Keywords1 {
		e.hldb[kw] = consts.Hl1
	}
	for _, kw := range colorMap.Keywords2 {
		e.hldb[kw] = consts.Hl2
	}
	for _, kw := range colorMap.Keywords3 {
		e.hldb[kw] = consts.Hl3
	}
	for _, t := range colorMap.StringTokens {
		e.hldb[t] = consts.HlStr
	}
	for _, t := range colorMap.SpecialTokens {
		e.hldb[t] = consts.HlSpc
	}
	e.hlOperators = colorMap.Operators
	for _, ch := range e.hlOperators {
		e.hldb[string(ch)] = consts.Hl1
	}

	e.commentToken = colorMap.CommentToken
	e.stringTokens = colorMap.StringTokens
}

func (e *Editor) getHighlightStyle(styleType consts.HlStyleType) (*tcell.Style, bool) {
	var hlStyle *tcell.Style
	switch styleType {
	case consts.Hl1:
		hlStyle = &e.hl1Style
	case consts.Hl2:
		hlStyle = &e.hl2Style
	case consts.Hl3:
		hlStyle = &e.hl3Style
	case consts.HlStr:
		hlStyle = &e.hlStrStyle
	case consts.HlSpc:
		hlStyle = &e.hlSpcStyle
	case consts.HlCom:
		hlStyle = &e.hlComStyle
	}
	return hlStyle, true
}

func getFileContent(fileName string) [][]rune {
	var fcl [][]rune = nil
	if fileName != "" {
		data, err := os.ReadFile(fileName)
		if err != nil {
			slog.Warn("could not read file", "err", err)
			return fcl
		}
		lines := strings.Split(string(data), "\n")
		fcl = make([][]rune, len(lines))
		for i, line := range lines {
			fcl[i] = []rune(line)
		}
	}
	return fcl
}
