// part6 - allow if, and implement more built-ins to make that useful.
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

	// offset of the argument to any IF
	ifOffset int
}

// NewEval returns a simple evaluator
func NewEval() *Eval {

	// Empty structure
	e := &Eval{}

	// Populate the dictionary of words we have implemented
	// which are hard-coded.
	e.Dictionary = []Word{
		{Name: "*", Function: e.mul},
		{Name: "+", Function: e.add},
		{Name: "-", Function: e.sub},
		{Name: ".", Function: e.print},
		{Name: "/", Function: e.div},
		{Name: ":", Function: e.startDefinition},
		{Name: "<", Function: e.lt},
		{Name: "<=", Function: e.ltEq},
		{Name: "=", Function: e.eq},
		{Name: "==", Function: e.eq},
		{Name: ">", Function: e.gt},
		{Name: ">=", Function: e.gtEq},
		{Name: "do", Function: e.do},  // NOP
		{Name: "if", Function: e.iff}, // NOP
		{Name: "invert", Function: e.invert},
		{Name: "drop", Function: e.drop},
		{Name: "dup", Function: e.dup},
		{Name: "emit", Function: e.emit},
		{Name: "loop", Function: e.loop},
		{Name: "print", Function: e.print},
		{Name: "swap", Function: e.swap},
		{Name: "then", Function: e.then}, // NOP
		{Name: "words", Function: e.words},
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

				// If the word was a "if"
				if tok == "if" {
					// we add the conditional-jump opcode
					e.tmp.Words = append(e.tmp.Words, -3)
					// placeholder jump-offset
					e.tmp.Words = append(e.tmp.Words, 99)

					// save the address of our stub,
					// so we can back-patch
					e.ifOffset = len(e.tmp.Words)
				}

				if tok == "then" {
					// back - patch the jump offset to the position of this word
					e.tmp.Words[e.ifOffset-1] = float64(len(e.tmp.Words) - 1)
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
// * Functions might contain a pointer to a function implemented in go.
//
//   If so we just call that pointer.
//
// * Functions will otherwise have lists of numbers, which point to
//   previously defined words.
//
//   In addition to the pointers to previously-defined words there are
//   also some special values:
//
//    "-1" means the next value is a number
//
//    "-2" is an unconditional jump, which will change our IP.
//
//    "-3" is a conditional-jump, which will change our IP if
//         the topmost item on the stack is "0".
//
func (e *Eval) evalWord(index int) {

	// Lookup the word
	word := e.Dictionary[index]

	// Is this implemented in golang?  If so just invoke the function
	// and we're done.
	if word.Function != nil {
		word.Function()
		return
	}

	//
	// TODO: Improve the way these special cases are handled.
	//
	// The reason this is handled like this, is to avoid poking the
	// indexes directly and risking array-overflow on malformed
	// word-lists.
	//
	// (i.e. When we see "[1, 2, -1]" the last instruction should add
	// the following number to the stack - but it is missing.  We want
	// to avoid messing around with the index to avoid that.)
	//

	// Adding a number?
	addNum := false

	// jumping?
	jump := false

	// jumping if the stack has a false-value?
	condJump := false

	// We need to allow control-jumps now, so we
	// have to store our index manually.
	ip := 0
	for ip < len(word.Words) {

		// the current opcode
		opcode := word.Words[ip]

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
				ip = int(opcode)

				// decrement as it'll get bumped at
				// the foot of the loop
				ip--
			}

			jump = false
		} else if condJump {
			// Jump only if 0 is on the top of the stack.
			//
			// i.e. This is an "if" test.
			val := e.Stack.Pop()
			if val == 0 {
				// change opcode
				ip = int(opcode)
				// decrement as it'll get bumped at
				// the foot of the loop
				ip--
			}
			condJump = false
		} else {

			// if we see -1 we're adding a number
			if opcode == -1 {
				addNum = true
			} else if opcode == -2 {
				// -2 is a jump
				jump = true
			} else if opcode == -3 {
				// -3 is a conditional-jump
				condJump = true
			} else {

				// otherwise we evaluate
				// otherwise eval as usual
				e.evalWord(int(opcode))
			}
		}

		// next instruction
		ip++
	}
}

// findWords returns the index in our dictionary of the entry for the
// given-name.  Returns -1 if the word cannot be found.
func (e *Eval) findWord(name string) int {
	for index, word := range e.Dictionary {
		if name == word.Name {
			return index
		}
	}
	return -1
}
