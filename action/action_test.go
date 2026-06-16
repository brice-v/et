package action

import "testing"

func TestParseQuit(t *testing.T) {
	a, err := Parse("quit")
	if err != nil {
		t.Fatal(err)
	}
	if a != Quit {
		t.Fatalf("expected Quit, got %d", a)
	}
}

func TestParseUnknown(t *testing.T) {
	_, err := Parse("foobar")
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestParseEmpty(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Fatal("expected error for empty action")
	}
}
