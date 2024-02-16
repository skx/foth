// Package lexer parses a string as FORTH
package lexer

import (
	"fmt"
)

// Type holds the type of a node
type Type string

// Types of tokens we return
const (

	// STRING is a string literal, enclosed in quotes.
	// e.g. "Steve"
	STRING = "string"

	// PSTRING is a string literal, enclosed in quotes,
	// which is printed to STDOUT when encountered.
	// e.g. ."Steve"
	PSTRING = "pstring"

	// WORD is anything which is not a STRING or PSTRING.
	//
	// This includes words, numbers, and other symbols which
	// are valid.
	//
	// The interpreter will internally convert the token `3.14'
	// from a string to a number, when appropriate, for example,
	// so we don't need an INTEGER/NUMBER type.
	//
	WORD = "word"
)

// Token is a single token.
//
// All our tokens have a name and type, only our two string-types have a value
// which is used for anything.
type Token struct {

	// Name holds the name of the token.
	Name string

	// Type holds the token-type
	Type Type

	// For the case of a string-literal we store the value here.
	Value string
}

// Lexer holds our state.
type Lexer struct {

	// input is the string we were given
	input string

	// offset points to the character we're currently looking at
	offset int
}

// New creates a new lexer which allows parsing a string of FORTH tokens
// into an array of tokens that can be interpreted.
func New(input string) *Lexer {
	return &Lexer{input: input}
}

// Tokens returns all the tokens from the given input-string.
func (l *Lexer) Tokens() ([]Token, error) {

	var res []Token

	// We walk the input from start to finish
	l.offset = 0

	// Value of the current token - built up character by character.
	cur := ""

	for l.offset < len(l.input) {

		c := l.input[l.offset]
		switch string(c) {

		case " ", "\n", "\r", "\t":

			// If we've built up a word then we save it away.
			if len(cur) != 0 {
				res = append(res, Token{Name: cur, Type: WORD})
				cur = ""
			}

		case "\\":
			// Comment to the end of the line
			for l.offset < len(l.input) {
				if l.input[l.offset] == '\n' {
					break
				}
				l.offset++
			}

		case "'":
			// We parse 'x' as the ASCII code of the character x.

			// can we peek ahead two characters?
			if l.offset+2 < len(l.input) {

				// confirm we have a close
				if string(l.input[l.offset+2]) == "'" {

					c = l.input[l.offset+1]
					d := int(c)
					s := fmt.Sprintf("%d", d)
					res = append(res, Token{Name: s, Type: WORD})
					l.offset += 2
				} else {
					return res, fmt.Errorf("syntax error")
				}
			} else {
				return res, fmt.Errorf("unterminated single-character constant")
			}

		case "(":

			// skip the "("
			l.offset++

			// Eat the comment - which is everything
			// between the "(" and ")" (inclusive)
			//
			// NOTE: Nested comments are prohibited
			closed := false
			for l.offset < len(l.input) {
				if l.input[l.offset] == ')' {
					closed = true
					break
				}
				if l.input[l.offset] == '(' {
					return res, fmt.Errorf("nested comments are illegal")
				}
				l.offset++
			}
			if !closed {
				return res, fmt.Errorf("unterminated comment")
			}

			// This is for strings
		case "\"":

			// skip the opening """
			l.offset++

			// Read the string, handling control-characters & etc
			str, err := l.readString()
			if err != nil {
				return nil, err
			}

			res = append(res, Token{Name: "\"", Value: str, Type: STRING})

			// This is for ." xxx "
		case ".":

			// ensure we don't walk off the array
			if l.offset+1 < len(l.input) {

				// next character is a string?
				if l.input[l.offset+1] == '"' {

					l.offset += 2

					// Read the string, handling control-characters & etc
					str, err := l.readString()
					if err != nil {
						return nil, err
					}

					res = append(res, Token{Name: ".\"", Value: str, Type: PSTRING})
				} else {
					cur = cur + "."
				}
			} else {
				cur = cur + "."
			}
		default:
			cur = cur + string(c)
		}
		l.offset++
	}

	// end token?
	if cur != "" {
		res = append(res, Token{Name: cur, Type: WORD})
	}

	// All done.
	return res, nil
}

// readString is called to read until a close of the string
// is encountered.  (i.e. ").
func (l *Lexer) readString() (string, error) {

	// We're now inside a string
	closed := false
	val := ""

	for l.offset < len(l.input) {

		c := l.input[l.offset]

		if c == '"' {
			closed = true
			l.offset++
			break
		}

		// Handle \n, etc.
		if c == '\\' {

			// if there is another character
			if l.offset+1 < len(l.input) {

				// look at what it is
				l.offset++
				c := l.input[l.offset]

				if c == 'n' {
					c = '\n'
				}
				if c == 'r' {
					c = '\r'
				}
				if c == 't' {
					c = '\t'
				}
				if c == '"' {
					c = '"'
				}
				if c == '\\' {
					c = '\\'
				}

				val += string(c)
				l.offset++
				continue
			}
		}

		// default
		val += string(l.input[l.offset])
		l.offset++
	}

	// Failed to close the string?
	if !closed {
		return val, fmt.Errorf("unterminated string")
	}

	// Returned okay
	return val, nil
}
