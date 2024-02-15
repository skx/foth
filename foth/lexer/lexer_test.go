package lexer

import (
	"strings"
	"testing"
)

// Test single-character constants
func TestCharacter(t *testing.T) {

	var out []Token
	var err error

	// unterminated character
	l := New(` '*`)
	_, err = l.Tokens()
	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if !strings.Contains(err.Error(), "unterminated") {
		t.Fatalf("got an error, but the wrong kind")
	}

	// single character
	l = New(` '*' `)
	out, err = l.Tokens()
	if err != nil {
		t.Fatalf("error lexing: %s", err)
	}
	if len(out) != 1 {
		t.Fatalf("wrong number of tokens")
	}
	if out[0].Name != "42" {
		t.Fatalf("unexpected result: %v", out[0].Name)
	}

	// unclosed character
	l = New(` '** `)
	_, err = l.Tokens()
	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if !strings.Contains(err.Error(), "syntax error") {
		t.Fatalf("got an error, but the wrong kind")
	}

}

// Comments should be removed
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

// Escape characters
func TestEscapeCharacters(t *testing.T) {

	l := New("\"hello\\n\\r\\t\\r\\\"\\\\steve\"")

	toks, err := l.Tokens()
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}

	expect := "hello\n\r\t\r\"\\steve"
	if toks[0].Value != expect {
		t.Fatalf("unexpected value for string; got '%s' not '%s'", toks[0].Value, expect)
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
	if out[1].Value != " foo bar baz " {
		t.Fatalf("got bad string: '%s'", out[1].Value)
	}
	if out[2].Name != "end" {
		t.Fatalf("got bad suffix")
	}

}

func TestString2(t *testing.T) {

	l := New("start \" foo bar baz \" end")
	out, err := l.Tokens()

	if err != nil {
		t.Fatalf("error lexing")
	}

	if out[0].Name != "start" {
		t.Fatalf("got bad prefix")
	}
	if out[1].Name != "\"" {
		t.Fatalf("got bad string")
	}
	if out[1].Value != " foo bar baz " {
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

// Unterminated strings are a bug
func TestStringUnterminated2(t *testing.T) {

	l := New("  \" string here ")

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
