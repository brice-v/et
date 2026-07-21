package editor

import (
	"strings"

	"github.com/brice-v/et/config"
	"github.com/brice-v/et/keys"
	"github.com/gdamore/tcell/v3"
)

// ActionID identifies a named editor action.
type ActionID int

const (
	ActionQuit ActionID = iota
	ActionFind
	ActionFindIgnoreCase
	ActionFindRegex
	ActionToggleTerminal
	ActionToggleLineEnding
	ActionToggleExpandTabs
	ActionTerminalIncrease
	ActionTerminalDecrease
	ActionExitPrompt
	ActionFindNext
	ActionFindPrevious
	ActionMoveUp
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionEnter
	ActionBackspace
	ActionDelete
	ActionInsertTab

	numActions
)

var actionNames = [...]string{
	ActionQuit:              "quit",
	ActionFind:              "find",
	ActionFindIgnoreCase:    "find_ignore_case",
	ActionFindRegex:         "find_regex",
	ActionToggleTerminal:    "toggle_terminal",
	ActionToggleLineEnding:  "toggle_line_ending",
	ActionToggleExpandTabs:  "toggle_expand_tabs",
	ActionTerminalIncrease:  "terminal_increase",
	ActionTerminalDecrease:  "terminal_decrease",
	ActionExitPrompt:        "exit_prompt",
	ActionFindNext:          "find_next",
	ActionFindPrevious:      "find_previous",
	ActionMoveUp:            "move_up",
	ActionMoveDown:          "move_down",
	ActionMoveLeft:          "move_left",
	ActionMoveRight:         "move_right",
	ActionEnter:             "enter",
	ActionBackspace:         "backspace",
	ActionDelete:            "delete",
	ActionInsertTab:         "insert_tab",
}

func (id ActionID) String() string {
	if int(id) < len(actionNames) {
		return actionNames[id]
	}
	return "unknown"
}

// Action is a function that performs an editor action.
type Action func(e *Editor)

// ActionEntry describes a registered action with its name.
type ActionEntry struct {
	ID   ActionID
	Name string
	Action
}

// AllActions returns all registered actions.
func AllActions() []ActionEntry {
	entries := make([]ActionEntry, numActions)
	for i := range entries {
		id := ActionID(i)
		entries[i] = ActionEntry{
			ID:     id,
			Name:   id.String(),
			Action: globalActions[i],
		}
	}
	return entries
}

// globalActions maps each ActionID to its Action handler.
var globalActions [numActions]Action

func init() {
	globalActions = [numActions]Action{
		ActionQuit:              func(e *Editor) { e.Exit = true },
		ActionFind:              func(e *Editor) { e.activateFind() },
		ActionFindIgnoreCase:    func(e *Editor) { e.activateFindMode(findModeIgnoreCase) },
		ActionFindRegex:         func(e *Editor) { e.activateFindMode(findModeRegex) },
		ActionToggleTerminal:    func(e *Editor) { e.ToggleTerminal() },
		ActionToggleLineEnding:  func(e *Editor) { e.buffer.ToggleLineEnding() },
		ActionToggleExpandTabs:  func(e *Editor) { e.ToggleExpandTabs() },
		ActionTerminalIncrease:  func(e *Editor) { e.IncreaseTerminalHeight() },
		ActionTerminalDecrease:  func(e *Editor) { e.DecreaseTerminalHeight() },
		ActionExitPrompt:        func(e *Editor) { e.exitPrompt() },
		ActionFindNext:          func(e *Editor) { e.findNextMatch() },
		ActionFindPrevious:      func(e *Editor) { e.findPreviousMatch() },
		ActionMoveUp:            func(e *Editor) { e.handleMoveUp() },
		ActionMoveDown:          func(e *Editor) { e.handleMoveDown() },
		ActionMoveLeft:          func(e *Editor) { e.handleMoveLeft() },
		ActionMoveRight:         func(e *Editor) { e.handleMoveRight() },
		ActionEnter:             func(e *Editor) { e.handleEnter() },
		ActionBackspace:         func(e *Editor) { e.handleBackspace() },
		ActionDelete:            func(e *Editor) { e.handleDelete() },
		ActionInsertTab:         func(e *Editor) { e.insertTab() },
	}
}

// chordBinding links a config key to an action for chord dispatch.
type chordBinding struct {
	key      config.Key
	action   Action
	findOnly bool
}

// chordBindings returns the list of chord-action bindings from the editor's config.
func (e *Editor) chordBindings() []chordBinding {
	return []chordBinding{
		{key: e.cfg.KeyBindings.Quit.Suffix, action: globalActions[ActionQuit]},
		{key: e.cfg.KeyBindings.Find.Suffix, action: globalActions[ActionFind]},
		{key: e.cfg.KeyBindings.ToggleLineEnding.Suffix, action: globalActions[ActionToggleLineEnding]},
		{key: e.cfg.KeyBindings.ToggleExpandTabs.Suffix, action: globalActions[ActionToggleExpandTabs]},
		{key: e.cfg.KeyBindings.ToggleTerminal.Suffix, action: globalActions[ActionToggleTerminal]},
		{key: e.cfg.KeyBindings.TerminalIncreaseChord.Suffix, action: globalActions[ActionTerminalIncrease]},
		{key: e.cfg.KeyBindings.TerminalDecreaseChord.Suffix, action: globalActions[ActionTerminalDecrease]},
		{key: e.cfg.KeyBindings.FindSecondary1Chord.Suffix, action: globalActions[ActionFindIgnoreCase], findOnly: true},
		{key: e.cfg.KeyBindings.FindSecondary2Chord.Suffix, action: globalActions[ActionFindRegex], findOnly: true},
	}
}

// activateFind starts or resets find mode.
func (e *Editor) activateFind() {
	if e.promptMode == promptModeFind {
		e.Find.Mode = findModeExact
		e.Find.LastSearchTerm = ""
		e.updatePromptLabel(e.getPromptFindLabel())
	} else {
		e.promptMode = promptModeFind
		e.prompt(e.getPromptFindLabel())
	}
}

// activateFindMode sets the find mode and resets search state (for chord-switching in find mode).
func (e *Editor) activateFindMode(mode findMode) {
	e.Find.Mode = mode
	e.Find.LastSearchTerm = ""
	e.updatePromptLabel(e.getPromptFindLabel())
}

// insertTab inserts either spaces or a tab character based on expandTabs setting.
func (e *Editor) insertTab() {
	if e.expandTabs {
		e.handleInsertRune(strings.Repeat(" ", e.cfg.TabWidth))
	} else {
		e.handleInsertRune("\t")
	}
}

// DoAction performs the action identified by id.
func (e *Editor) DoAction(id ActionID) {
	if int(id) < len(globalActions) && globalActions[id] != nil {
		globalActions[id](e)
	}
}

// IsChordAction reports whether an action ID is dispatched through the chord prefix.
func (id ActionID) IsChordAction() bool {
	switch id {
	case ActionQuit, ActionFind, ActionFindIgnoreCase, ActionFindRegex,
		ActionToggleTerminal, ActionToggleLineEnding, ActionToggleExpandTabs,
		ActionTerminalIncrease, ActionTerminalDecrease:
		return true
	}
	return false
}

// IsDirectAction reports whether an action ID is dispatched directly (not through a chord).
func (id ActionID) IsDirectAction() bool {
	switch id {
	case ActionMoveUp, ActionMoveDown, ActionMoveLeft, ActionMoveRight,
		ActionEnter, ActionBackspace, ActionDelete, ActionInsertTab,
		ActionExitPrompt, ActionFindNext, ActionFindPrevious:
		return true
	}
	return false
}

// handleChordSuffix attempts to match a chord suffix key to a binding.
// It returns true if the key was consumed by a chord action.
func (e *Editor) handleChordSuffix(key tcell.Key, keyAsRune string, k *tcell.EventKey) bool {
	for _, cb := range e.chordBindings() {
		if cb.findOnly && e.promptMode != promptModeFind {
			continue
		}
		if keys.IsKey(key, keyAsRune, k.Modifiers(), cb.key) {
			cb.action(e)
			return true
		}
	}
	e.chordInvalidSuffix = "invalid suffix " + k.Name()
	return true
}
