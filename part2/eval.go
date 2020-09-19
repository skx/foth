// part2 - allow defining words in terms of others
package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Word is the structure for a single word
type Word struct {
	// Name is the name of the function "+", "print", etc.
	Name string

	// Function is the function-pointer to call to invoke it.
	//
	// If this is nil then instead we interpret the known-codes
	// from previously defined words.
	Function func()

	// Words holds the words we execute if the function-pointer
	// is empty.
	//
	// The indexes here are relative to the Dictionary our evaluator
	// holds/maintains
	Words []int
}

// Eval is our evaluation structure
type Eval struct {

	// Internal stack
	Stack Stack

	// Dictionary entries
	Dictionary []Word
}

// NewEval returns a simple evaluator
func NewEval() *Eval {

	// Empty structure
	e := &Eval{}

	// Populate the dictionary of words we have implemented
	// which are hard-coded.
	//
	// All of these are coded in go, except for the new function
	// square - which is implemented as just "dup *".  The implementation
	// just has the offsets in this array of the words used to invoke.
	e.Dictionary = []Word{
		{Name: "+", Function: e.add},                        // 0
		{Name: "-", Function: e.sub},                        // 1
		{Name: "*", Function: e.mul},                        // 2
		{Name: "/", Function: e.div},                        // 3
		{Name: "print", Function: e.print},                  // 4
		{Name: ".", Function: e.print},                      // 5
		{Name: "dup", Function: e.dup},                      // 6
		{Name: "square", Function: nil, Words: []int{6, 2}}, // 7
	}
	return e
}

// Eval processes a list of tokens.
//
// This is invoked by our repl with a line of input at the time.
func (e *Eval) Eval(args []string) {

	for _, tok := range args {

		// Trim the leading/trailing spaces,
		// and skip any empty tokens
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}

		// Did we handle this as a dictionary item?
		handled := false
		for index, word := range e.Dictionary {
			if tok == word.Name {
				e.evalWord(index)
				handled = true
			}
		}

		// If we didn't handle this as a word, then
		// assume it is a number.
		if !handled {
			i, err := strconv.ParseFloat(tok, 64)
			if err != nil {
				fmt.Printf("%s: %s\n", tok, err.Error())
				return
			}

			e.Stack.Push(i)
		}
	}
}

// evalWord evaluates a word, by index from the dictionary
func (e *Eval) evalWord(index int) {

	word := e.Dictionary[index]
	if word.Function != nil {
		word.Function()
	} else {
		for _, offset := range word.Words {
			e.evalWord(offset)
		}
	}
}
