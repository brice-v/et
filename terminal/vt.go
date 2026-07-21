package terminal

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v3"
	"github.com/mattn/go-runewidth"
)

type (
	column int
	row    int
)

// ptyBackend abstracts the PTY or alternative backend for the VT terminal.
type ptyBackend interface {
	io.Reader
	io.Writer
	io.Closer
	io.StringWriter
	Resize(cols, rows uint16) error
}

// VT models a virtual terminal
type VT struct {
	// If true, OSC8 enables the output of OSC8 strings. Otherwise, any OSC8
	// sequences will be stripped
	OSC8 bool
	// Set the TERM environment variable to be passed to the command's
	// environment. If not set, xterm-256color will be used
	TERM string

	mu sync.Mutex

	activeScreen  [][]cell
	altScreen     [][]cell
	primaryScreen [][]cell

	charsets charsets
	cursor   cursor
	margin   margin
	mode     mode
	tabStop  []column
	// lastCol is a flag indicating we printed in the last col
	lastCol bool

	primaryState cursorState
	altState     cursorState

	dirty        bool
	eventHandler func(tcell.Event)
	parser       *Parser
	pty          ptyBackend
	surface      Surface
	events       chan tcell.Event

	mouseBtn tcell.ButtonMask
}

type cursorState struct {
	cursor   cursor
	decawm   bool
	decom    bool
	charsets charsets
}

type margin struct {
	top    row
	bottom row
	left   column
	right  column
}

func New() *VT {
	tabs := []column{}
	for i := 7; i < (50 * 7); i += 8 {
		tabs = append(tabs, column(i))
	}
	return &VT{
		OSC8: true,
		charsets: charsets{
			designations: map[charsetDesignator]charset{
				g0: ascii,
				g1: ascii,
				g2: ascii,
				g3: ascii,
			},
		},
		mode: decawm | dectcem,
		primaryState: cursorState{
			charsets: charsets{
				designations: map[charsetDesignator]charset{
					g0: ascii,
					g1: ascii,
					g2: ascii,
					g3: ascii,
				},
			},
			decawm: true,
		},
		altState: cursorState{
			charsets: charsets{
				designations: map[charsetDesignator]charset{
					g0: ascii,
					g1: ascii,
					g2: ascii,
					g3: ascii,
				},
			},
			decawm: true,
		},
		tabStop:      tabs,
		eventHandler: func(ev tcell.Event) {},
		// Buffering to 2 events. If there is ever a case where one
		// sequence can trigger two events, this should be increased
		events: make(chan tcell.Event, 2),
	}
}

// Start starts the terminal with the specified command. Start returns when the
// command has been successfully started.
func (vt *VT) Start(cmd *exec.Cmd) error {
	if cmd == nil {
		return fmt.Errorf("no command to run")
	}
	vt.mu.Lock()
	w, h := vt.surface.Size()
	vt.mu.Unlock()

	if vt.TERM == "" {
		vt.TERM = "xterm-256color"
	}

	env := os.Environ()
	if cmd.Env != nil {
		env = cmd.Env
	}
	cmd.Env = append(env, "TERM="+vt.TERM)

	var err error
	vt.pty, err = newPtyBackend(cmd, uint16(w), uint16(h))
	if err != nil {
		return err
	}

	vt.Resize(w, h)
	vt.parser = NewParser(vt.pty)
	go func() {
		defer vt.recover()
		for {
			select {
			case ev := <-vt.events:
				vt.eventHandler(ev)
			default:
				seq := vt.parser.Next()
				switch seq := seq.(type) {
				case EOF:
					vt.eventHandler(&EventClosed{
						EventTerminal: newEventTerminal(vt),
					})
					return
				default:
					vt.update(seq)
				}
			}
		}
	}()
	return nil
}

func (vt *VT) update(seq Sequence) {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	switch seq := seq.(type) {
	case Print:
		vt.print(rune(seq))
	case C0:
		vt.c0(rune(seq))
	case ESC:
		esc := append(seq.Intermediate, seq.Final)
		vt.esc(string(esc))
	case CSI:
		csi := append(seq.Intermediate, seq.Final)
		vt.csi(string(csi), seq.Parameters)
	case OSC:
		vt.osc(string(seq.Payload))
	case DCS:
	case DCSData:
	case DCSEndOfData:
	}
	if !vt.dirty {
		vt.dirty = true
		vt.postEvent(&EventRedraw{
			EventTerminal: newEventTerminal(vt),
		})
	}
}

func (vt *VT) String() string {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	str := strings.Builder{}
	for row := range vt.activeScreen {
		for col := range vt.activeScreen[row] {
			_, _ = str.WriteRune(vt.activeScreen[row][col].rune())
			for _, comb := range vt.activeScreen[row][col].combining {
				_, _ = str.WriteRune(comb)
			}
		}
		if row < vt.height()-1 {
			str.WriteRune('\n')
		}
	}
	return str.String()
}

func (vt *VT) recover() {
	err := recover()
	if err == nil {
		return
	}
	ret := strings.Builder{}
	fmt.Fprintf(&ret, "cursor row=%d col=%d\n", vt.cursor.row, vt.cursor.col)
	fmt.Fprintf(&ret, "margin left=%d right=%d\n", vt.margin.left, vt.margin.right)
	fmt.Fprintf(&ret, "%s\n", err)
	ret.Write(debug.Stack())

	vt.postEvent(&EventPanic{
		EventTerminal: newEventTerminal(vt),
		Error:         fmt.Errorf("%s", ret.String()),
	})
	vt.Close()
}

// row, col, style, vis
func (vt *VT) Cursor() (int, int, tcell.CursorStyle, bool) {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	vis := vt.mode&dectcem > 0
	return int(vt.cursor.row), int(vt.cursor.col), vt.cursor.style, vis
}

func (vt *VT) Resize(w int, h int) {
	primary := vt.primaryScreen
	vt.altScreen = make([][]cell, h)
	vt.primaryScreen = make([][]cell, h)
	for i := range vt.altScreen {
		vt.altScreen[i] = make([]cell, w)
		vt.primaryScreen[i] = make([]cell, w)
	}
	last := vt.cursor.row
	vt.margin.bottom = row(h) - 1
	vt.margin.right = column(w) - 1
	vt.cursor.row = 0
	vt.cursor.col = 0
	vt.lastCol = false
	vt.activeScreen = vt.primaryScreen

	// transfer primary to new, skipping the last row
	for row := range len(primary) {
		if row == int(last) {
			break
		}
		wrapped := false
		for col := range len(primary[0]) {
			cell := primary[row][col]
			vt.cursor.attrs = cell.attrs
			vt.print(cell.content)
			wrapped = cell.wrapped
		}
		if !wrapped {
			vt.nel()
		}
	}
	switch vt.mode & smcup {
	case 0:
		vt.activeScreen = vt.primaryScreen
	default:
		vt.activeScreen = vt.altScreen
	}

	if vt.pty != nil {
		if err := vt.pty.Resize(uint16(w), uint16(h)); err != nil {
			slog.Warn("error setting terminal backend size", "err", err)
		}
	}
}

func (vt *VT) width() int {
	if len(vt.activeScreen) > 0 {
		return len(vt.activeScreen[0])
	}
	return 0
}

func (vt *VT) height() int {
	return len(vt.activeScreen)
}

// print sets the current cell contents to the given rune. The attributes will
// be copied from the current cursor attributes
func (vt *VT) print(r rune) {
	if vt.charsets.designations[vt.charsets.selected] == decSpecialAndLineDrawing {
		shifted, ok := decSpecial[r]
		if ok {
			r = shifted
		}
	}

	// If we are single-shifted, move the previous charset into the current
	if vt.charsets.singleShift {
		vt.charsets.selected = vt.charsets.saved
	}

	if vt.cursor.col == vt.margin.right && vt.lastCol {
		col := vt.cursor.col
		rw := vt.cursor.row
		vt.activeScreen[rw][col].wrapped = true
		vt.nel()
	}

	col := vt.cursor.col
	rw := vt.cursor.row
	w := runewidth.RuneWidth(r)

	if vt.mode&irm != 0 {
		line := vt.activeScreen[rw]
		for i := vt.margin.right; i > col; i -= 1 {
			line[i] = line[i-column(w)]
		}
	}
	if col > column(vt.width())-1 {
		col = column(vt.width()) - 1
	}
	if rw > row(vt.height()-1) {
		rw = row(vt.height() - 1)
	}

	if w == 0 {
		if col-1 < 0 {
			return
		}
		vt.activeScreen[rw][col-1].combining = append(vt.activeScreen[rw][col-1].combining, r)
		return
	}
	cell := cell{
		content: r,
		width:   w,
		attrs:   vt.cursor.attrs,
	}

	vt.activeScreen[rw][col] = cell

	// Set trailing cells to a space if wide rune
	for i := column(1); i < column(w); i++ {
		if col+i > vt.margin.right {
			break
		}
		vt.activeScreen[rw][col+i].content = ' '
		vt.activeScreen[rw][col+i].attrs = vt.cursor.attrs
	}

	switch {
	case vt.mode&decawm != 0 && col == vt.margin.right:
		vt.lastCol = true
	case col == vt.margin.right:
		// don't move the cursor
	default:
		vt.cursor.col += column(w)
	}
}

// scrollUp shifts all text upward by n rows. Semantically, this is backwards -
// usually scroll up would mean you shift rows down
func (vt *VT) scrollUp(n int) {
	for row := range vt.activeScreen {
		if row > int(vt.margin.bottom) {
			continue
		}
		if row < int(vt.margin.top) {
			continue
		}
		if row+n > int(vt.margin.bottom) {
			for col := vt.margin.left; col <= vt.margin.right; col += 1 {
				vt.activeScreen[row][col].erase(vt.cursor.attrs)
			}
			continue
		}
		copy(vt.activeScreen[row], vt.activeScreen[row+n])
	}
}

// scrollDown shifts all lines down by n rows.
func (vt *VT) scrollDown(n int) {
	for r := vt.margin.bottom; r >= vt.margin.top; r -= 1 {
		if r-row(n) < vt.margin.top {
			for col := vt.margin.left; col <= vt.margin.right; col += 1 {
				vt.activeScreen[r][col].erase(vt.cursor.attrs)
			}
			continue
		}
		copy(vt.activeScreen[r], vt.activeScreen[r-row(n)])
	}
}

func (vt *VT) Close() {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	if vt.pty != nil {
		if err := vt.pty.Close(); err != nil {
			slog.Warn("error closing terminal backend", "err", err)
		}
	}
}

func (vt *VT) Attach(fn func(ev tcell.Event)) {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	vt.eventHandler = fn
}

func (vt *VT) Detach() {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	vt.eventHandler = func(ev tcell.Event) {}
}

func (vt *VT) postEvent(ev tcell.Event) {
	vt.events <- ev
}

func (vt *VT) SetSurface(srf Surface) {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	vt.surface = srf
}

func (vt *VT) Draw() {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	vt.dirty = false
	if vt.surface == nil {
		return
	}
	for row := 0; row < vt.height(); row += 1 {
		for col := 0; col < vt.width(); {
			cell := vt.activeScreen[row][col]
			w := cell.width
			vt.surface.SetContent(col, row, cell.content, cell.combining, cell.attrs)
			if w == 0 {
				w = 1
			}
			col += w
		}
	}
}

func (vt *VT) HandleEvent(e tcell.Event) bool {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	switch e := e.(type) {
	case *tcell.EventKey:
		if _, err := vt.pty.Write([]byte(keyCode(e))); err != nil {
			slog.Warn("error writing key event to backend", "err", err)
		}
		return true
	case *tcell.EventPaste:
		switch {
		case vt.mode&paste == 0:
			return false
		case e.Start():
			if _, err := vt.pty.Write([]byte(info.PasteStart)); err != nil {
				slog.Warn("error writing paste start to backend", "err", err)
			}
			return true
		case e.End():
			if _, err := vt.pty.Write([]byte(info.PasteEnd)); err != nil {
				slog.Warn("error writing paste end to backend", "err", err)
			}
			return true
		}
	case *tcell.EventMouse:
		str := vt.handleMouse(e)
		if _, err := vt.pty.Write([]byte(str)); err != nil {
			slog.Warn("error writing mouse event to backend", "err", err)
		}
	}
	return false
}
