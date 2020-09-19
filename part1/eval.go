// part1 - only allow a fixed number of hard-coded words
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
	Function func()
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
	e.Dictionary = []Word{
		{"+", e.add},       // 0
		{"-", e.sub},       // 1
		{"*", e.mul},       // 2
		{"/", e.div},       // 3
		{"print", e.print}, // 4
		{".", e.print},     // 5
	}
	return e
}

// Eval processes a list of tokens.
//
// This is invoked by our repl with a line of input at the time.
func (e *Eval) Eval(args []string) {

	//
	// token = NextToken()
	// if token is in dictionary
	//   call function from that dict entry
	// else if token is a number
	//  push that number onto the data stack
	for _, tok := range args {

		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		var i float64
		var err error

		// Dictionary item?
		for _, word := range e.Dictionary {
			if tok == word.Name {
				word.Function()
				goto end
			}
		}

		// is this a number?
		i, err = strconv.ParseFloat(tok, 64)
		if err != nil {
			fmt.Printf("%s: %s\n", tok, err.Error())
			return
		}

		e.Stack.Push(i)

	end:
	}
}
