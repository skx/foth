// Package lexer parses a string as FORTH
package lexer

import (
	"fmt"
	"strings"
)

// Token is a single token.
//
// All our tokens have a name, only string-literals have a value which
// is used for anything.
type Token struct {

	// Name holds the name of the token.
	Name string

	// For the case of a string-literal we store the value here.
	Value string
}

// Lexer holds our state.
type Lexer struct {

	// input is the string we were given
	input string
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
	offset := 0

	// Value of the current token - built up character by character.
	cur := ""

	for offset < len(l.input) {

		c := l.input[offset]
		switch string(c) {

		case " ", "\n", "\r", "\t":

			// If we've built up a word then we save it away.
			if len(cur) != 0 {
				res = append(res, Token{Name: cur})
				cur = ""
			}

		case "\\":
			// Comment to the end of the line
			for offset < len(l.input) {
				if l.input[offset] == '\n' {
					break
				}
				offset++
			}

		case "'":
			// We parse 'x' as the ASCII code of the character x.

			// can we peek ahead two characters?
			if offset+2 < len(l.input) {

				// confirm we have a close
				if string(l.input[offset+2]) == "'" {

					c := l.input[offset+1]
					d := int(c)
					s := fmt.Sprintf("%d", d)
					res = append(res, Token{Name: s})
					offset += 2
				} else {
					return res, fmt.Errorf("syntax error")
				}
			} else {
				return res, fmt.Errorf("unterminated single-character constant")
			}

		case "(":

			// skip the "("
			offset++

			// Eat the comment - which is everything
			// between the "(" and ")" (inclusive)
			//
			// NOTE: Nested comments are prohibited
			closed := false
			for offset < len(l.input) {
				if l.input[offset] == ')' {
					closed = true
					break
				}
				if l.input[offset] == '(' {
					return res, fmt.Errorf("nested comments are illegal")
				}
				offset++
			}
			if !closed {
				return res, fmt.Errorf("unterminated comment")
			}

		case ".":

			// ensure we don't walk off the array
			if offset+1 < len(l.input) {

				// next character is a string?
				if l.input[offset+1] == '"' {

					// skip the "."
					offset++

					// skip the opening """
					offset++

					// We're now inside the string
					closed := false
					val := ""
					for offset < len(l.input) {
						if l.input[offset] == '"' {
							closed = true
							break
						} else {
							val += string(l.input[offset])
						}
						offset++
					}

					// Failed to close the string?
					if !closed {
						return res, fmt.Errorf("unterminated string")
					}

					// Otherwise save it away
					val = strings.TrimSpace(val)
					res = append(res, Token{Name: ".\"", Value: val})
				} else {
					cur = cur + "."
				}
			} else {
				cur = cur + "."
			}
		default:
			cur = cur + string(c)
		}
		offset++
	}

	// end token?
	if cur != "" {
		res = append(res, Token{Name: cur})
	}

	// All done.
	return res, nil
}
