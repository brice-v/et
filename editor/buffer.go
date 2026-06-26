package editor

import (
	"log/slog"
	"os"
	"strings"
)

type Buffer struct {
	lines [][]rune
}

func NewBuffer(fileName string) *Buffer {
	if fileName == "" {
		return &Buffer{}
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
