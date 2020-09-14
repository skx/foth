// part3 - allow defining words in our own environment.
package main

import (
	"fmt"
	"os"
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

	// Are we in a compiling mode?
	compiling bool

	// Temporary word we're compiling
	tmp Word
}

// NewEval returns a simple evaluator
func NewEval() *Eval {

	// Empty structure
	e := &Eval{}

	// Populate the dictionary of words we have implemented
	// which are hard-coded.
	e.Dictionary = []Word{
		Word{Name: "+", Function: e.add},             // 0
		Word{Name: "-", Function: e.sub},             // 1
		Word{Name: "*", Function: e.mul},             // 2
		Word{Name: "/", Function: e.div},             // 3
		Word{Name: "print", Function: e.print},       // 4
		Word{Name: ".", Function: e.print},           // 5
		Word{Name: "dup", Function: e.dup},           // 6
		Word{Name: ":", Function: e.startDefinition}, // 7
		// Note we don't handle ";" here.
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

		// Trim the leading/trailing spaces,
		// and skip any empty tokens
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}

		// Are we in compiling mode?
		if e.compiling {

			// If we don't have a name
			if e.tmp.Name == "" {

				// is the name used?
				idx := e.findWord(tok)
				if idx != -1 {
					fmt.Printf("word %s already defined\n", tok)
					os.Exit(1)
				}

				// save the name
				e.tmp.Name = tok
				continue
			}

			// End of a definition?
			if tok == ";" {
				e.Dictionary = append(e.Dictionary, e.tmp)
				e.tmp.Name = ""
				e.tmp.Words = []int{}
				e.compiling = false
				continue

			}

			// OK we have a name, so lookup the word definition
			// for it.
			idx := e.findWord(tok)
			if idx >= 0 {
				// Found it
				e.tmp.Words = append(e.tmp.Words, idx)
			}

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

func (e *Eval) findWord(name string) int {
	for index, word := range e.Dictionary {
		if name == word.Name {
			return index
		}
	}
	return -1
}
