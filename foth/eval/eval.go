// Package eval contains our simple forth-like interpreter.
package eval

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"foth/lexer"
	"foth/stack"
)

// Word is the structure for a single word.
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

	// Does this word switch us into immediate-mode?
	StartImmediate bool

	// Does this word end us from immediate-mode?
	EndImmediate bool
}

// Eval is our evaluation structure
type Eval struct {

	// Stack holds our operands.
	Stack stack.Stack

	// Dictionary entries
	Dictionary []Word

	// Are we in a compiling mode?
	compiling bool

	// Are we in debug-mode?
	debug bool

	// Are we in immediate-mode?
	immediate int

	// Temporary word we're compiling
	tmp Word

	// open of the last do
	doOpen int

	// When we generate IF-statements we have to patch one or two
	// offsets, depending on whether there is an ELSE branch present
	// or not.
	//
	// Here we keep track of them
	ifOffset1 int
	ifOffset2 int
}

// New returns a simple evaluator, which will allow executing forth-like words.
func New() *Eval {

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
		Word{Name: ";", Function: e.nop},
		Word{Name: "<", Function: e.lt},
		Word{Name: "<=", Function: e.ltEq},
		Word{Name: "=", Function: e.eq},
		Word{Name: "==", Function: e.eq},
		Word{Name: ">", Function: e.gt},
		Word{Name: ">=", Function: e.gtEq},
		Word{Name: "dump", Function: e.dump},
		Word{Name: "debug", Function: e.debugSet},
		Word{Name: "debug?", Function: e.debugp},
		Word{Name: "do", Function: e.nop, StartImmediate: true},
		Word{Name: "drop", Function: e.drop},
		Word{Name: "dup", Function: e.dup},
		Word{Name: "else", Function: e.nop},
		Word{Name: "emit", Function: e.emit},
		Word{Name: "if", Function: e.nop, StartImmediate: true},
		Word{Name: "invert", Function: e.invert},
		Word{Name: "loop", Function: e.loop, EndImmediate: true},
		Word{Name: "print", Function: e.print},
		Word{Name: "swap", Function: e.swap},
		Word{Name: "then", Function: e.nop, EndImmediate: true},
		Word{Name: "words", Function: e.words},
		Word{Name: "#words", Function: e.wordLen},
	}

	return e
}

// Eval evaluates the given expression.
//
// This is invoked by our repl with a line of input at the time.
func (e *Eval) Eval(input string) error {

	// Lex our input string into a series of tokens.
	//
	// This is done for two reasons:
	//
	//  1.  We want to remove comments
	//
	//  2.  We support inline strings, such as the following:
	//
	//       ." foo bar "
	//
	//      Blindly splitting on whitespace would screw those up
	//
	l := lexer.New(input)
	args, err := l.Tokens()
	if err != nil {
		return err
	}

	//
	// For each token..
	//
	for _, token := range args {

		// Get the name of the token.
		//
		// The name is the only thing we care about, except
		// in the case of string-literals
		tok := token.Name

		// Trim the leading/trailing spaces,
		// and skip any empty tokens
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		// Are we in compiling mode?
		if e.compiling || (e.immediate > 0) {

			// If so compile the token
			err := e.compileToken(tok)
			if err != nil {
				return err
			}

			// And loop around.
			continue
		}

		// Lookup this word from our dictionary
		idx := e.findWord(tok)
		if idx != -1 {

			// Are we starting immediate mode?
			if !e.compiling && e.Dictionary[idx].StartImmediate {
				e.immediate++

				err := e.compileToken(tok)
				if err != nil {
					return err
				}
			} else {

				err := e.evalWord(idx)
				if err != nil {
					return err
				}
			}
		} else {

			// If we didn't handle this as a word, then
			// assume it is a number.
			i, err := strconv.ParseFloat(tok, 64)
			if err != nil {
				return fmt.Errorf("11 failed to convert %s to number %s", tok, err.Error())
			}

			e.Stack.Push(i)
		}
	}

	return nil
}

// compileToken is called with a new token, when we're in compiling-mode.
func (e *Eval) compileToken(tok string) error {

	// Did we start in immediate-mode?
	imm := (!e.compiling && e.immediate > 0)

	if imm {

		// In immediate mode we just setup a bogus
		// word-name, which can never be legal.
		//
		// We'll then execute it immediately post-definition.
		e.tmp.Name = "$ $"
	}

	// If we don't have a name
	if e.tmp.Name == "" {

		// is the name used?  If so remove it
		idx := e.findWord(tok)
		if idx != -1 {
			e.Dictionary[idx].Name = ""
		}

		// save the name
		e.tmp.Name = tok
		return nil
	}

	// End of a definition?
	if tok == ";" {

		// Save the word to our dictionary
		e.Dictionary = append(e.Dictionary, e.tmp)

		// Show what we compiled each new definition
		// to, when running in debug-mode
		if e.debug {
			e.dumpWord(len(e.Dictionary) - 1)
		}

		// reset for the next definition
		e.tmp.Name = ""
		e.tmp.Words = []float64{}
		e.compiling = false
		return nil
	}

	// Is the user adding an existing word?
	idx := e.findWord(tok)
	if idx >= 0 {

		// Found the word, add to the end.
		e.tmp.Words = append(e.tmp.Words, float64(idx))

		//
		// Now some special cases.
		//
		// Horrid
		//
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

		//
		// Conditional support is a bit nasty.
		//
		// Basically we expect to allow someone to writ
		// something like:
		//
		//  : foo 0 < if neg else pos then
		//
		// That translates to:
		//
		//  if foo < 0 {
		//    neg
		//  } else {
		//    pos
		// }
		//
		// We want to convert that to:
		//
		//  COND
		//  IF
		//     CONDJUMP xx
		//     NEG CODE
		//     JMP end
		//   ELSE
		//  xx:
		//     POS CODE
		//   THEN
		//  end:
		//
		// In short we have to insert a conditional-jump, and
		// an unconditional one.
		//
		// We then use back-patching to fixup the offsets.
		//
		if tok == "if" {
			// reset both possible places to patch
			e.ifOffset1 = 0
			e.ifOffset2 = 0

			// we add the conditional-jump opcode
			e.tmp.Words = append(e.tmp.Words, -3)
			// placeholder jump-offset
			e.tmp.Words = append(e.tmp.Words, 99)

			// The offset of the last instruction is
			// the 99-byte we just added as a placeholder
			// record that.
			e.ifOffset1 = len(e.tmp.Words)
		}

		if tok == "else" {
			e.tmp.Words[e.ifOffset1-1] = float64(len(e.tmp.Words) + 2)

			// before we compile the end we have to
			// add a jump to after the THEN
			e.tmp.Words = append(e.tmp.Words, -4)
			e.tmp.Words = append(e.tmp.Words, 999)

			e.ifOffset2 = len(e.tmp.Words)
		}

		if tok == "then" {

			// back - patch the jump offset to the position of this word
			if e.ifOffset2 > 0 {
				e.tmp.Words[e.ifOffset2-1] = float64(len(e.tmp.Words))
			} else {
				e.tmp.Words[e.ifOffset1-1] = float64(len(e.tmp.Words) - 1)
			}
			// Reset for future
			e.ifOffset2 = 0
			e.ifOffset1 = 0
		}

		// Did we just end immediate mode?
		if e.Dictionary[idx].EndImmediate {

			// If we're inside an immediate
			if !e.compiling && e.immediate > 0 {
				e.immediate--
			}

			if e.immediate == 0 && imm {

				// We've compiled the word.
				e.Dictionary = append(e.Dictionary, e.tmp)

				if e.debug {
					fmt.Printf("Completed the temporary word - '$ $'\n")
					e.dumpWord(len(e.Dictionary) - 1)
				}

				// reset for the next definition
				e.tmp.Name = ""
				e.tmp.Words = []float64{}
				// Run it.
				e.evalWord(len(e.Dictionary) - 1)
				return nil
			}
		}

		return nil
	}

	// At this point we assume the user entered a number
	// so we save a magic "-1" flag in our
	// definition, and then the number itself
	e.tmp.Words = append(e.tmp.Words, -1)

	// Convert to float
	val, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return fmt.Errorf("22 failed to convert %s to number %s", tok, err.Error())
	}
	e.tmp.Words = append(e.tmp.Words, val)

	return nil
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
//    "-2" is an loop jump, which will change our IP.
//
//    "-3" is a conditional-jump, which will change our IP if
//         the topmost item on the stack is "0".
//
//    "-4" is an unconditional jump, which will change our IP
//
func (e *Eval) evalWord(index int) error {

	// Lookup the word in our dictionary.
	word := e.Dictionary[index]

	// Is this implemented in golang?  If so just invoke the function
	// and we're done.
	if word.Function != nil {

		if e.debug {
			fmt.Printf(" calling built-in word %s\n", word.Name)
		}
		err := word.Function()
		return err
	}

	if e.debug {
		fmt.Printf(" calling dynamic stuff\n")
	}

	//
	// We use a simple state-machine to handle some of our
	// "opcodes".  Opcodes are basically indexes into our
	// dictionary of previously-defined words, and also some
	// special codes (which are less than zero).
	//
	// The reason this is handled like this, is to avoid poking the
	// indexes directly and risking array-overflow on malformed
	// word-lists.
	//
	// (i.e. When we see "[1, 2, -1]" the last instruction should add
	// the following number to the stack - but it is missing.  We want
	// to avoid messing around with the index to avoid that.)
	//

	state := "default"

	// We need to allow control-jumps now, so we
	// have to store our index manually.
	ip := 0
	for ip < len(word.Words) {

		// the current opcode
		opcode := word.Words[ip]

		// adding a number?
		if state == "add-number" {
			if e.debug {
				fmt.Printf(" storing %f on stack\n", opcode)
			}

			// add to stack
			e.Stack.Push(opcode)

			state = "default"

		} else if state == "loop-jump" {
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

			state = "default"
		} else if state == "cond-jump" {
			// Jump only if 0 is on the top of the stack.
			//
			// i.e. This is an "if" test.
			val, err := e.Stack.Pop()
			if err != nil {
				return err
			}

			if val == 0 {
				if e.debug {

					fmt.Printf(" popped %f from stack - jumping to %f\n", val, opcode)
				}
				// change opcode
				ip = int(opcode)
				// decrement as it'll get bumped at
				// the foot of the loop
				ip--
			} else {
				if e.debug {
					fmt.Printf(" popped %f from stack - not making conditional jump\n", val)
				}
			}
			state = "default"
		} else if state == "jump" {
			// change opcode
			ip = int(opcode)
			// decrement as it'll get bumped at
			// the foot of the loop
			ip--
			state = "default"
		} else if state == "default" {

			// if we see -1 we're adding a number
			switch opcode {
			case -1:
				state = "add-number"
			case -2:
				state = "loop-jump"
			case -3:
				state = "cond-jump"
			case -4:
				state = "jump"
			default:
				err := e.evalWord(int(opcode))
				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("unknown state '%s'", state)
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

// dumpWord dumps the definition of the given word.
func (e *Eval) dumpWord(idx int) {

	word := e.Dictionary[idx]

	codes := []string{}

	// Show the names of the codes, as well as their indexes
	off := 0
	for off < len(word.Words) {

		v := word.Words[int(off)]

		if v == -1 {
			codes = append(codes, fmt.Sprintf("%d: store %f", off, word.Words[off+1]))
			off++
		} else if v == -2 {
			codes = append(codes, fmt.Sprintf("%d: [loop-jmp %f]", off, word.Words[off+1]))
			off++
		} else if v == -3 {
			codes = append(codes, fmt.Sprintf("%d: [cond-jmp %f]", off, word.Words[off+1]))
			off++
		} else if v == -4 {
			codes = append(codes, fmt.Sprintf("%d: [jmp %f]", off, word.Words[off+1]))
			off++
		} else {
			codes = append(codes, fmt.Sprintf("%d: %s", off, e.Dictionary[int(v)].Name))
		}
		off++
	}

	// Didn't decompile?  Then it was a native-word
	if len(codes) == 0 {
		fmt.Printf("Word '%s' - native\n", word.Name)
	} else {
		// Otherwise show the bytecode.
		fmt.Printf("Word '%s'\n %s\n", word.Name, strings.Join(codes, "\n "))
	}
}
