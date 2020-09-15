// part5 - allow loops
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Loop holds a loop - as these can be nested
type Loop struct {
	cur int
	max int
}

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
	Words []float64
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

	// open of the last do
	doOpen int
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
		Word{Name: "emit", Function: e.emit},         // 8
		Word{Name: "do", Function: e.do},             // 9
		Word{Name: "loop", Function: e.loop},         // 10
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
				e.tmp.Words = []float64{}
				e.compiling = false

				continue

			}

			// OK we have a name, so lookup the word definition
			// for it.
			idx := e.findWord(tok)
			if idx >= 0 {
				// Found it
				e.tmp.Words = append(e.tmp.Words, float64(idx))

				// If the word was a "DO"
				if tok == "do" {
					e.doOpen = len(e.tmp.Words) - 1
				}
				// if the word was a "LOOP"
				if tok == "loop" {

					// offset of do must be present
					e.tmp.Words = append(e.tmp.Words, -2)
					e.tmp.Words = append(e.tmp.Words, float64(e.doOpen))
				}
			} else {

				// OK we assume the user entered a number
				// so we save a magic "-1" flag in our
				// definition, and then the number itself
				e.tmp.Words = append(e.tmp.Words, -1)

				// Convert to float
				val, err := strconv.ParseFloat(tok, 64)
				if err != nil {
					fmt.Printf("%s: %s\n", tok, err.Error())
					return
				}
				e.tmp.Words = append(e.tmp.Words, val)
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
//
// * Functions might have go-pointers
//   if so we just call the pointer.
//
// * Functions might have lists of numbers, which point to previous
//   definitions.  There are two special-cases, "-1" means the next
//   value is a number, and "-2" is a jump opcode to move around
//   in our bytecode array.
//
func (e *Eval) evalWord(index int) {

	// Lookup the word
	word := e.Dictionary[index]
	if word.Function != nil {
		word.Function()
		return
	}

	// Adding a number?
	addNum := false

	// jumping?
	jump := false

	// We need to allow control-jumps now, so we
	// have to store our index manually.
	inst := 0
	for inst < len(word.Words) {

		// the current opcode
		opcode := word.Words[inst]

		// adding a number?
		if addNum {
			// add to stack
			e.Stack.Push(opcode)
			addNum = false
		} else if jump {
			// If the two top-most entries
			// are not equal, then jump
			cur := e.Stack.Pop()
			max := e.Stack.Pop()

			if max > cur {
				// put them back
				e.Stack.Push(max)
				e.Stack.Push(cur)

				// change opcode
				inst = int(opcode)
				// decrement as it'll get bumped at
				// the foot of the loop
				inst--
			}

			jump = false
		} else {

			// if we see -1 we're adding a number
			if opcode == -1 {
				addNum = true
			} else if opcode == -2 {
				// -2 is a jump
				jump = true
			} else {

				// otherwise we evaluate
				// otherwise eval as usual
				e.evalWord(int(opcode))
			}
		}

		// next instruction
		inst++
	}
}

// findWords returns the index in our dictionary of the entry for the
// given-name.  Returns -1 on error.
func (e *Eval) findWord(name string) int {
	for index, word := range e.Dictionary {
		if name == word.Name {
			return index
		}
	}
	return -1
}
