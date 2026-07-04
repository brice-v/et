package editor

import (
	"github.com/brice-v/et/config"
	"github.com/brice-v/et/consts"
	"github.com/brice-v/et/lexer"
	"log/slog"

	"github.com/gdamore/tcell/v3"
)

type HighlightState struct {
	hl1Style   tcell.Style
	hl2Style   tcell.Style
	hl3Style   tcell.Style
	hlStrStyle tcell.Style
	hlSpcStyle tcell.Style
	hlComStyle tcell.Style

	hldb         map[string]consts.HlStyleType
	hlOperators  string
	commentToken string
	stringTokens []string
}

func NewHighlightState(cfg *config.Config, fileExtension string) *HighlightState {
	hs := &HighlightState{}
	hs.setup(cfg, fileExtension)
	return hs
}

func (hs *HighlightState) setup(cfg *config.Config, fileExtension string) {
	if fileExtension == "" {
		return
	}
	fileType, ok := cfg.FileExtensions[fileExtension]
	if !ok {
		slog.Warn("no syntax highlighting available for file extension", "fileExtension", fileExtension)
		return
	}
	colorMap, ok := cfg.Colors.Languages[fileType]
	if !ok {
		slog.Warn("no color map for highlighting found for fileType", "fileType", fileType)
		return
	}
	hs.hl1Style = tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(colorMap.Color1.Color)
	hs.hl2Style = tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(colorMap.Color2.Color)
	hs.hl3Style = tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(colorMap.Color3.Color)
	hs.hlStrStyle = tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(colorMap.ColorString.Color)
	hs.hlSpcStyle = tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(colorMap.SpecialColor.Color)
	hs.hlComStyle = tcell.StyleDefault.Background(cfg.Colors.Background.Color).Foreground(colorMap.CommentColor.Color)
	hs.hldb = make(map[string]consts.HlStyleType)
	for _, kw := range colorMap.Keywords1 {
		hs.hldb[kw] = consts.Hl1
	}
	for _, kw := range colorMap.Keywords2 {
		hs.hldb[kw] = consts.Hl2
	}
	for _, kw := range colorMap.Keywords3 {
		hs.hldb[kw] = consts.Hl3
	}
	for _, t := range colorMap.StringTokens {
		hs.hldb[t] = consts.HlStr
	}
	for _, t := range colorMap.SpecialTokens {
		hs.hldb[t] = consts.HlSpc
	}
	hs.hlOperators = colorMap.Operators
	for _, ch := range hs.hlOperators {
		hs.hldb[string(ch)] = consts.Hl1
	}
	hs.commentToken = colorMap.CommentToken
	hs.stringTokens = colorMap.StringTokens
}

func (hs *HighlightState) getStyle(styleType consts.HlStyleType) *tcell.Style {
	switch styleType {
	case consts.Hl1:
		return &hs.hl1Style
	case consts.Hl2:
		return &hs.hl2Style
	case consts.Hl3:
		return &hs.hl3Style
	case consts.HlStr:
		return &hs.hlStrStyle
	case consts.HlSpc:
		return &hs.hlSpcStyle
	case consts.HlCom:
		return &hs.hlComStyle
	}
	return nil
}

func (hs *HighlightState) newLexer(line []rune) *lexer.Lexer {
	return lexer.New(line, hs.hldb, hs.hlOperators, hs.commentToken, hs.stringTokens)
}
