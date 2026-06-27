package editor

import (
	"log/slog"
	"os"
	"strings"
)

type Buffer struct {
	lines [][]rune
	dirty bool
}

func NewBuffer(fileName string) *Buffer {
	if fileName == "" {
		return &Buffer{lines: [][]rune{{}}}
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		slog.Warn("could not read file", "err", err)
		return &Buffer{}
	}
	lines := strings.Split(string(data), "\n")
	fcl := make([][]rune, len(lines))
	for i, line := range lines {
		fcl[i] = []rune(line)
	}
	return &Buffer{lines: fcl}
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
