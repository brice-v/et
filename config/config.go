package config

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Key struct {
	tcell.Key
	Modifiers tcell.ModMask
}

type KeyChord struct {
	Prefix Key `json:"prefix"`
	Suffix Key `json:"suffix"`
}

type Color struct {
	color.Color
}

type ColorMap struct {
	Keywords1     []string `json:"keywords1"`
	Color1        Color    `json:"color1"`
	Keywords2     []string `json:"keywords2"`
	Color2        Color    `json:"color2"`
	Keywords3     []string `json:"keywords3"`
	Color3        Color    `json:"color3"`
	StringTokens  []string `json:"string_tokens"`
	ColorString   Color    `json:"color_string"`
	Operators     string   `json:"operators"`
	SpecialTokens []string `json:"special_tokens"`
	SpecialColor  Color    `json:"special_color"`
	CommentToken  string   `json:"comment_token"`
	CommentColor  Color    `json:"comment_color"`
}

type Colors struct {
	Foreground            Color               `json:"foreground"`
	Background            Color               `json:"background"`
	StatusBar             Color               `json:"status_bar"`
	MatchHighlight        Color               `json:"match_highlight"`
	CurrentMatchHighlight Color               `json:"current_match_highlight"`
	Languages             map[string]ColorMap `json:"languages"`
}

type KeyBindings struct {
	Quit                 Key      `json:"quit"`
	ExitPrompt           Key      `json:"exit_prompt"`
	Find                 Key      `json:"find"`
	FindNext             Key      `json:"find_next"`
	FindPrevious         Key      `json:"find_previous"`
	FindSecondary1Chord  KeyChord `json:"find_secondary1_chord"`
	FindSecondary2Chord  KeyChord `json:"find_secondary2_chord"`
	ToggleTerminal       KeyChord `json:"toggle_terminal"`
	ToggleLineEnding     Key      `json:"toggle_line_ending"`
	ToggleExpandTabs     Key      `json:"toggle_expand_tabs"`
	TerminalIncreaseChord KeyChord `json:"terminal_increase_chord"`
	TerminalDecreaseChord KeyChord `json:"terminal_decrease_chord"`
}

type Config struct {
	Colors              Colors            `json:"colors"`
	KeyBindings         KeyBindings       `json:"keybindings"`
	TabWidth            int               `json:"tab_width"`
	LeftPadString       string            `json:"left_pad_string"`
	ShowLineNumbers     bool              `json:"show_line_numbers"`
	FileExtensions      map[string]string `json:"file_extensions"`
	DisableHighlighting bool              `json:"disable_highlighting"`
	CursorStyle         string            `json:"cursor_style"`
	CursorColor         Color             `json:"cursor_color"`
	DefaultLineEnding   string            `json:"default_line_ending"`
	ExpandTabs          bool              `json:"expand_tabs"`
	TerminalHeightPercentage float64       `json:"terminal_height_percentage"`
}

func NewDefault() *Config {
	return &Config{
		Colors: Colors{
			Foreground:            DefaultColorForeground(),
			Background:            DefaultColorBackground(),
			StatusBar:             DefaultColorStatusBar(),
			MatchHighlight:        DefaultColorMatchHighlight(),
			CurrentMatchHighlight: DefaultColorCurrentMatchHighlight(),
			Languages:             DefaultLanguagesColorMap(),
		},
		KeyBindings: KeyBindings{
			Quit:                 DefaultKeyBindingQuit(),
			ExitPrompt:           DefaultKeyBindingExitPrompt(),
			Find:                 DefaultKeyBindingFind(),
			FindNext:             DefaultKeyBindingFindNext(),
			FindPrevious:         DefaultKeyBindingFindPrevious(),
			FindSecondary1Chord:  DefaultKeyBindingFindSecondary1Chord(),
			FindSecondary2Chord:  DefaultKeyBindingFindSecondary2Chord(),
			ToggleTerminal:       DefaultKeyBindingToggleTerminalChord(),
			ToggleLineEnding:     DefaultKeyBindingToggleLineEnding(),
			ToggleExpandTabs:     DefaultKeyBindingToggleExpandTabs(),
			TerminalIncreaseChord: DefaultKeyBindingTerminalIncreaseChord(),
			TerminalDecreaseChord: DefaultKeyBindingTerminalDecreaseChord(),
		},
		TabWidth:            DefaultTabWidth(),
		LeftPadString:       DefaultLeftPadString(),
		ShowLineNumbers:     DefaultShowLineNumbers(),
		FileExtensions:      DefaultFileExtensions(),
		DisableHighlighting: DefaultDisableHighlighting(),
		CursorStyle:         DefaultCursorStyle(),
		CursorColor:         DefaultCursorColor(),
		DefaultLineEnding:   DefaultLineEnding(),
		ExpandTabs:          DefaultExpandTabs(),
		TerminalHeightPercentage: DefaultTerminalHeightPercentage(),
	}
}

func Parse(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			slog.Warn("parse close file failed", "err", err)
		}
	}()

	cfg := NewDefault()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
