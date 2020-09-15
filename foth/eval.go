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
	Function func() error

	// Words holds the words we execute if the function-pointer
	// is empty.
	//
	// The indexes here are relative to the Dictionary our evaluator
	// holds/maintains
	//
	// NOTE: We specifically store `float64` here so that we can add
	// floats to the stack when compiling.
	Words []float64
}

// Eval is our evaluation structure
type Eval struct {

	// Stack holds our operands.
	Stack Stack

	// Dictionary entries
	Dictionary []Word

	// Are we in a compiling mode?
	compiling bool

	// Are we in debug-mode?
	debug bool

	// Temporary word we're compiling
	tmp Word

	// open of the last do
	doOpen int

	// offset of the argument to any IF
	ifOffset int
}

// NewEval returns a simple evaluator, which will allow executing forth-like words.
func NewEval() *Eval {

	// Empty structure
	e := &Eval{}

	// Are we debugging?
	if os.Getenv("DEBUG") != "" {
		e.debug = true
	}

	// Populate our built-in functions.
	e.Dictionary = []Word{
		Word{Name: "*", Function: e.mul},
		Word{Name: "+", Function: e.add},
		Word{Name: "-", Function: e.sub},
		Word{Name: ".", Function: e.print},
		Word{Name: "/", Function: e.div},
		Word{Name: ":", Function: e.startDefinition},
		Word{Name: "<", Function: e.lt},
		Word{Name: "<=", Function: e.lt_eq},
		Word{Name: "=", Function: e.eq},
		Word{Name: "==", Function: e.eq},
		Word{Name: ">", Function: e.gt},
		Word{Name: ">=", Function: e.gt_eq},
		Word{Name: "do", Function: e.do},
		Word{Name: "if", Function: e.iff},
		Word{Name: "invert", Function: e.invert},
		Word{Name: "drop", Function: e.drop},
		Word{Name: "dup", Function: e.dup},
		Word{Name: "emit", Function: e.emit},
		Word{Name: "loop", Function: e.loop},
		Word{Name: "print", Function: e.print},
		Word{Name: "swap", Function: e.swap},
		Word{Name: "then", Function: e.then},
		Word{Name: "words", Function: e.words},
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

					// The name is used, so we need to remove it
					e.Dictionary[idx].Name = ""
				}

				// save the name
				e.tmp.Name = tok

				continue
			}

			// End of a definition?
			if tok == ";" {
				e.Dictionary = append(e.Dictionary, e.tmp)

				// Show what we compiled each new definition
				// to, when running in debug-mode
				if e.debug {
					e.dumpWords()
				}

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
func (e *Eval) evalWord(index int) error {

	// Lookup the word in our dictionary.
	word := e.Dictionary[index]

	// Is this implemented in golang?  If so just invoke the function
	// and we're done.
	if word.Function != nil {
		word.Function()
		return nil
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
			cur, ee := e.Stack.Pop()
			if ee != nil {
				return ee
			}
			max, eee := e.Stack.Pop()
			if eee != nil {
				return eee
			}

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
			val, err := e.Stack.Pop()
			if err != nil {
				return err
			}
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
			switch opcode {
			case -1:
				addNum = true
			case -2:
				jump = true
			case -3:
				condJump = true
			default:
				e.evalWord(int(opcode))
			}
		}

		// next instruction
		ip++
	}

	return nil
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

func (e *Eval) dumpWords() {
	codes := []string{}

	// Show the names of the codes,
	// as well as their indexes
	off := 0
	for off < len(e.tmp.Words) {

		v := e.tmp.Words[int(off)]

		if v == -1 {
			codes = append(codes, fmt.Sprintf("%f", e.tmp.Words[off+1]))
			off++
		} else if v == -2 {
			codes = append(codes, fmt.Sprintf("[jmp %f]", e.tmp.Words[off+1]))
			off++
		} else if v == -3 {
			codes = append(codes, fmt.Sprintf("[cond-jmp %f]", e.tmp.Words[off+1]))
			off++
		} else {
			codes = append(codes, fmt.Sprintf("%s", e.Dictionary[int(v)].Name))
		}
		off++
	}
	fmt.Printf("Compiled word '%s' -> %s\n", e.tmp.Name, strings.Join(codes, " "))
}
