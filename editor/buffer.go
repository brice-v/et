package editor

import (
	"log/slog"
	"os"
	"strings"
)

type Buffer struct {
	lines      [][]rune
	dirty      bool
	lineEnding string
}

func NewBuffer(fileName string) *Buffer {
	if fileName == "" {
		return &Buffer{lines: [][]rune{{}}, lineEnding: "lf"}
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		slog.Warn("could not read file", "err", err)
		return &Buffer{lineEnding: "lf"}
	}
	lineEnding := detectLineEnding(string(data))
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")
	fcl := make([][]rune, len(lines))
	for i, line := range lines {
		fcl[i] = []rune(line)
	}
	return &Buffer{lines: fcl, lineEnding: lineEnding}
}

func detectLineEnding(content string) string {
	if strings.Contains(content, "\r\n") {
		return "crlf"
	}
	return "lf"
}

func (b *Buffer) NumLines() int {
	return len(b.lines)
}

func (b *Buffer) Line(n int) []rune {
	return b.lines[n]
}

func (b *Buffer) IsOpen() bool {
	return b.lines != nil
}

func (b *Buffer) IsDirty() bool {
	return b.dirty
}

func (b *Buffer) LineEnding() string {
	return b.lineEnding
}

func (b *Buffer) SetLineEnding(le string) {
	b.lineEnding = le
	b.dirty = true
}

func (b *Buffer) ToggleLineEnding() {
	if b.lineEnding == "lf" {
		b.lineEnding = "crlf"
	} else {
		b.lineEnding = "lf"
	}
	b.dirty = true
}

func (b *Buffer) Bytes() []byte {
	var sb strings.Builder
	sep := "\n"
	if b.lineEnding == "crlf" {
		sep = "\r\n"
	}
	for i, line := range b.lines {
		if i > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(string(line))
	}
	return []byte(sb.String())
}

func (b *Buffer) InsertRune(lineNum int, col int, r rune) {
	if len(b.lines) == 0 {
		b.lines = [][]rune{{}}
	}
	for lineNum >= len(b.lines) {
		b.lines = append(b.lines, []rune{})
	}
	line := b.lines[lineNum]
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}
	newLine := make([]rune, len(line)+1)
	copy(newLine, line[:col])
	newLine[col] = r
	copy(newLine[col+1:], line[col:])
	b.lines[lineNum] = newLine
	b.dirty = true
}

func (b *Buffer) DeleteRune(lineNum int, col int) {
	line := b.lines[lineNum]
	b.lines[lineNum] = append(line[:col], line[col+1:]...)
	b.dirty = true
}

func (b *Buffer) SplitLine(lineNum int, col int) {
	line := b.lines[lineNum]
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}
	right := make([]rune, len(line)-col)
	copy(right, line[col:])
	b.lines[lineNum] = b.lines[lineNum][:col]

	b.lines = append(b.lines, nil)
	copy(b.lines[lineNum+2:], b.lines[lineNum+1:])
	b.lines[lineNum+1] = right
	b.dirty = true
}

func (b *Buffer) JoinLine(lineNum int) {
	b.lines[lineNum] = append(b.lines[lineNum], b.lines[lineNum+1]...)
	b.lines = append(b.lines[:lineNum+1], b.lines[lineNum+2:]...)
	b.dirty = true
}
