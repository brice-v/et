package editor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewBufferEmpty(t *testing.T) {
	b := NewBuffer("")
	if b.LineEnding() != "lf" {
		t.Errorf("LineEnding() = %q, want %q", b.LineEnding(), "lf")
	}
	if b.NumLines() != 1 {
		t.Errorf("NumLines() = %d, want 1", b.NumLines())
	}
}

func TestNewBufferLF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lf.txt")
	if err := os.WriteFile(path, []byte("hello\nworld\n"), 0644); err != nil {
		t.Fatal(err)
	}
	b := NewBuffer(path)
	if b.LineEnding() != "lf" {
		t.Errorf("LineEnding() = %q, want %q", b.LineEnding(), "lf")
	}
}

func TestNewBufferCRLF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "crlf.txt")
	if err := os.WriteFile(path, []byte("hello\r\nworld\r\n"), 0644); err != nil {
		t.Fatal(err)
	}
	b := NewBuffer(path)
	if b.LineEnding() != "crlf" {
		t.Errorf("LineEnding() = %q, want %q", b.LineEnding(), "crlf")
	}
}

func TestNewBufferMixedCRLF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mixed.txt")
	if err := os.WriteFile(path, []byte("line1\r\nline2\nline3\r\n"), 0644); err != nil {
		t.Fatal(err)
	}
	b := NewBuffer(path)
	if b.LineEnding() != "crlf" {
		t.Errorf("LineEnding() = %q, want %q (first line has CRLF)", b.LineEnding(), "crlf")
	}
}

func TestToggleLineEnding(t *testing.T) {
	b := NewBuffer("")
	if b.LineEnding() != "lf" {
		t.Errorf("expected lf, got %q", b.LineEnding())
	}
	b.ToggleLineEnding()
	if b.LineEnding() != "crlf" {
		t.Errorf("expected crlf, got %q", b.LineEnding())
	}
	if !b.IsDirty() {
		t.Error("expected dirty after toggling line ending")
	}
	b.ToggleLineEnding()
	if b.LineEnding() != "lf" {
		t.Errorf("expected lf, got %q", b.LineEnding())
	}
}

func TestSetLineEnding(t *testing.T) {
	b := NewBuffer("")
	b.SetLineEnding("crlf")
	if b.LineEnding() != "crlf" {
		t.Errorf("expected crlf, got %q", b.LineEnding())
	}
	if !b.IsDirty() {
		t.Error("expected dirty after setting line ending")
	}
}

func TestBytesLF(t *testing.T) {
	b := NewBuffer("")
	b.InsertRune(0, 0, 'a')
	b.InsertRune(0, 1, 'b')
	got := string(b.Bytes())
	want := "ab"
	if got != want {
		t.Errorf("Bytes() = %q, want %q", got, want)
	}
}

func TestBytesCRLF(t *testing.T) {
	b := NewBuffer("")
	b.InsertRune(0, 0, 'a')
	b.SplitLine(0, 1)
	b.InsertRune(1, 0, 'b')
	b.SetLineEnding("crlf")
	got := string(b.Bytes())
	want := "a\r\nb"
	if got != want {
		t.Errorf("Bytes() = %q, want %q", got, want)
	}
}

func TestBytesMultiLineLF(t *testing.T) {
	b := NewBuffer("")
	b.InsertRune(0, 0, 'a')
	b.SplitLine(0, 1)
	b.InsertRune(1, 0, 'b')
	got := string(b.Bytes())
	want := "a\nb"
	if got != want {
		t.Errorf("Bytes() = %q, want %q", got, want)
	}
}
