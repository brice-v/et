package terminal

import "github.com/gdamore/tcell/v3"

type cell struct {
	content   rune
	combining []rune
	width     int
	attrs     tcell.Style
	wrapped   bool
}

func (c *cell) rune() rune {
	if c.content == rune(0) {
		return ' '
	}
	return c.content
}

func (c *cell) erase(s tcell.Style) {
	bg := s.GetBackground()
	c.content = 0
	c.attrs = tcell.StyleDefault.Background(bg)
}

// selectiveErase removes the cell content, but keeps the attributes
func (c *cell) selectiveErase() {
	c.content = 0
}
