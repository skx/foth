package lexer

import (
	"strings"
	"testing"
)

// Comments should be reoved
func TestComment(t *testing.T) {

	l := New(` to \ This is a comment
( comment here )fo
`)

	out, err := l.Tokens()
	if err != nil {
		t.Fatalf("error lexing: %s", err)
	}
	if len(out) != 2 {
		t.Fatalf("Unexpected output, got: %v", out)
	}
	if out[0].Name != "to" {
		t.Fatalf("got bad prefix")
	}
	if out[1].Name != "fo" {
		t.Fatalf("got bad suffix")
	}
}

// Nested comments are a bug
func TestCommentNested(t *testing.T) {

	l := New("  ( comment ( here ) ) ")

	_, err := l.Tokens()
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if !strings.Contains(err.Error(), "nested comments") {
		t.Fatalf("got an error, but the wrong one")
	}
}

// Unterminated comments are a bug
func TestCommentUnterminated(t *testing.T) {

	l := New("  ( comment here ")

	_, err := l.Tokens()
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if !strings.Contains(err.Error(), "unterminated comment") {
		t.Fatalf("got an error, but the wrong one")
	}
}

// Empty input should give empty output
func TestEmpty(t *testing.T) {
	l := New(" \t \r \n")
	out, err := l.Tokens()
	if err != nil {
		t.Fatalf("error lexing")
	}

	if len(out) != 0 {
		t.Fatalf("Unexpected output, got: %v", out)
	}

}

func TestString(t *testing.T) {

	l := New("start .\" foo bar baz \" end")
	out, err := l.Tokens()

	if err != nil {
		t.Fatalf("error lexing")
	}

	if out[0].Name != "start" {
		t.Fatalf("got bad prefix")
	}
	if out[1].Name != ".\"" {
		t.Fatalf("got bad string")
	}
	if out[1].Value != "foo bar baz" {
		t.Fatalf("got bad string: '%s'", out[1].Value)
	}
	if out[2].Name != "end" {
		t.Fatalf("got bad suffix")
	}

}

// Unterminated strings are a bug
func TestStringUnterminated(t *testing.T) {

	l := New("  .\" string here ")

	_, err := l.Tokens()
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if !strings.Contains(err.Error(), "unterminated string") {
		t.Fatalf("got an error, but the wrong one")
	}
}

// Not strings
func TestNotString(t *testing.T) {

	l := New("  .-  .")

	out, err := l.Tokens()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if out[0].Name != ".-" {
		t.Fatalf("got bad prefix")
	}
	if out[1].Name != "." {
		t.Fatalf("got bad suffix")
	}

}
