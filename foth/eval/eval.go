// Package eval contains our simple forth-like interpreter.
package eval

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/skx/foth/foth/lexer"
	"github.com/skx/foth/foth/stack"
)

// Loop holds the state of a running do/while loop.
type Loop struct {
	// Start is the starting number our loop begins from.
	Start int

	// Max holds the terminating number our loop finishes at.
	Max int

	// Current holds the current number of the iteration.
	Current int
}

// Variable is the structure for storing variable names, and contents
type Variable struct {
	// Name is the name of the variable
	Name string

	// Value is the value we store within it.
	Value float64
}

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

// Eval is our evaluation structure, which holds state of where
// we're executing code from.
//
// The state includes the variables, inline-strings, etc.
type Eval struct {

	// Stack holds our operands.
	Stack stack.Stack

	// Dictionary entries
	Dictionary []Word

	// STDOUT is the writer used for `.`, `print`, and `emit` words
	STDOUT *bufio.Writer

	// Private details

	// Literal strings.  As encountered in our program.
	strings []string

	// Have we already bumped our immediate-count?
	//
	// This is a gross-hack to account for the fact that
	// compileToken needs to bump the immediate-count when
	// it sees a nested do, or similar token.
	bumped bool

	// Are we in a compiling mode?
	compiling bool

	// Are we in debug-mode?
	debug bool

	// Are we in immediate-mode?
	immediate int

	// Temporary word we're compiling
	tmp Word

	// We keep a stack of the last time we saw a `do` token,
	// so we can pair it with the appropriate matching `loop`.
	doOpen []int

	// When we generate IF-statements we have to patch one or two
	// offsets, depending on whether there is an ELSE branch present
	// or not.
	//
	// Here we keep track of them
	ifOffset1 int
	ifOffset2 int

	// Loops stores loops which are currently open.
	//
	// When we compile `do` we add a new one, when the `loop`
	// word is completed we remove one.
	loops []Loop

	// Variables
	vars []Variable

	// Are we currently defining a variable?
	defining bool
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
		// comparisons
		{Name: "<", Function: e.lt},
		{Name: "<=", Function: e.ltEq},
		{Name: "=", Function: e.eq},
		{Name: "==", Function: e.eq}, // synonym
		{Name: ">", Function: e.gt},
		{Name: ">=", Function: e.gtEq},

		// conditionals
		{Name: "else", Function: e.nop},
		{Name: "if", Function: e.nop, StartImmediate: true},
		{Name: "then", Function: e.nop, EndImmediate: true},

		// debug-handling
		{Name: "debug", Function: e.debugSet},
		{Name: "debug?", Function: e.debugp},

		// I/O
		{Name: ".", Function: e.print},
		{Name: ".\"", Function: e.nop},
		{Name: "emit", Function: e.emit},
		{Name: "print", Function: e.print},

		// loop-handling
		{Name: "do", Function: e.nop, StartImmediate: true},
		{Name: "i", Function: e.i},
		{Name: "loop", Function: e.loop, EndImmediate: true},
		{Name: "m", Function: e.m},

		// mathematical
		{Name: "*", Function: e.mul},
		{Name: "+", Function: e.add},
		{Name: "-", Function: e.sub},
		{Name: "/", Function: e.div},
		{Name: "max", Function: e.max},
		{Name: "min", Function: e.min},
		{Name: "mod", Function: e.mod},

		// misc
		{Name: "nop", Function: e.nop},

		// stack-related
		{Name: ".s", Function: e.stackDump},
		{Name: "clearstack", Function: e.clearStack},
		{Name: "drop", Function: e.drop},
		{Name: "dup", Function: e.dup},
		{Name: "invert", Function: e.invert},
		{Name: "over", Function: e.over},
		{Name: "swap", Function: e.swap},

		// variable-handling
		{Name: "!", Function: e.setVar},
		{Name: "@", Function: e.getVar},
		{Name: "variable", Function: e.variable},

		// word-handling
		{Name: "#words", Function: e.wordLen},
		{Name: ":", Function: e.startDefinition},
		{Name: ";", Function: e.nop},
		{Name: "dump", Function: e.dump},
		{Name: "words", Function: e.words},
	}

	return e
}

// Eval evaluates the given expression.
//
// This is the main public-facing the user of this library would be expected
// to use.
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

		// Are we defining a variable?
		if e.defining {
			if e.debug {
				fmt.Printf("defining variable %s\n", tok)
			}
			e.vars = append(e.vars, Variable{Name: tok})
			e.defining = false
			continue
		}

		// Are we in compiling mode?
		if e.compiling || (e.immediate > 0) {

			// If so compile the token
			err := e.compileToken(token)
			if err != nil {
				return err
			}

			// And loop around.
			continue
		}

		// Is this an immediate print?  If so do it.
		if tok == ".\"" {
			e.printString(token.Value)
			continue
		}

		// Lookup this word from our dictionary
		idx := e.findWord(tok)
		if idx != -1 {

			// Are we starting immediate mode?
			if !e.compiling && e.Dictionary[idx].StartImmediate {
				e.immediate++
				e.bumped = true

				// Errors here can't happen.
				//
				// We only compile at the start
				// of conditionals, and they can't
				// happen.
				e.compileToken(token)
			} else {

				err := e.evalWord(idx)
				if err != nil {
					return err
				}
			}
		} else {

			// Is this a variable?  If so push the variable offset
			idx = e.findVariable(tok)
			if idx >= 0 {
				e.Stack.Push(float64(idx))
				continue
			}

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

// GetVariable returns the contents of the specified variable.
//
// This is designed to be used by host-applications which embed
// this library.
func (e *Eval) GetVariable(name string) (float64, error) {

	idx := e.findVariable(name)
	if idx >= 0 {
		return e.vars[idx].Value, nil
	}

	return 0, fmt.Errorf("variable %s not found", name)
}

// Reset updates the internal state of the evaluator, which is a useful
// thing to do if `Eval` returns an error
func (e *Eval) Reset() {

	// Clear the stack
	for !e.Stack.IsEmpty() {
		e.Stack.Pop()
	}

	// reset our state
	e.defining = false
	e.immediate = 0
	e.compiling = false

	// we're not defining anything
	e.tmp.Name = ""
	e.tmp.Words = []float64{}

	// we're not in a do/loop
	e.doOpen = []int{}
	e.loops = []Loop{}

	// we're not in a conditional
	e.ifOffset1 = 0
	e.ifOffset2 = 0
}

// SetVariable stores the specified value in the variable of the given
// name.
//
// This is designed to be used by host-applications which embed
// this library.
func (e *Eval) SetVariable(name string, value float64) {

	idx := e.findVariable(name)
	if idx >= 0 {
		e.vars[idx].Value = value
		return
	}

	e.vars = append(e.vars, Variable{Name: name, Value: value})
}

// SetWriter allows you to setup a special writer for all STDOUT
// messages this application will produce.
//
// i.e. Writes from `.`, `emit`, `print`, and string-immediates will
// go there.
func (e *Eval) SetWriter(writer *bufio.Writer) {
	e.STDOUT = writer
}

// compileToken is called with a new token, when we're in compiling-mode.
//
// This is called in two ways:
//
//  1.  To compile a new word.
//
//  2.  In "immediate" mode.  Where we compile a fake word,
//     with an impossible name (`$ $`), with the expectation
//     we'll then immediately execute it.
//
func (e *Eval) compileToken(token lexer.Token) error {

	tok := token.Name

	// Did we start in immediate-mode?
	//
	// We keep track of this here, because we might
	// disable immediate-mode later and need to know.
	imm := (!e.compiling && e.immediate > 0)

	if imm {

		// In immediate mode we are going to compile
		// a word which has an illegal-name, which ensures
		// we never overwrite a valid user-word
		//
		// We'll then execute it immediately post-definition.
		e.tmp.Name = "$ $"
	}

	// If we don't yet have a name
	if e.tmp.Name == "" {

		// is the name used?  If so remove it
		idx := e.findWord(tok)
		if idx != -1 {
			e.Dictionary[idx].Name = ""
		}

		// Set the name for this word.
		e.tmp.Name = tok
		return nil
	}

	// End of a definition?
	if tok == ";" {

		// Save the word to our dictionary
		e.tmp.Name = strings.ToLower(e.tmp.Name)
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

	// Is the user adding an existing word to the definition?
	idx := e.findWord(tok)
	if idx >= 0 {

		//
		// We have to do some juggling here because
		// when we find a word with `EndImmediate` we
		// terminate.
		//
		// So we have to make sure we don't terminate
		// early if something else opened.
		//
		// i.e. "loop" would usually terminate immediate-mode,
		// but we can't stop there for definitions that use
		// nested loops
		//
		if imm && e.Dictionary[idx].StartImmediate {
			if e.bumped {
				e.immediate--
				e.bumped = false
			}
			e.immediate++
		}

		// Found the word, add to the end.
		e.tmp.Words = append(e.tmp.Words, float64(idx))

		//
		// Now some special cases.
		//
		// Horrid
		//

		// If the word was a "DO"
		if tok == "do" {

			// keep track of where we are
			e.doOpen = append(e.doOpen, len(e.tmp.Words)-1)

			// we compile this into a "new-loop" instruction
			e.tmp.Words = append(e.tmp.Words, -10)
			e.tmp.Words = append(e.tmp.Words, 99) // dull
		}

		// if the word was a "LOOP"
		if tok == "loop" {

			// We load the loop, increment, etc.
			e.tmp.Words = append(e.tmp.Words, -11)
			e.tmp.Words = append(e.tmp.Words, 99) // dull

			// We've bumped the instance, and pushed
			// a result onto the stack now.
			//
			// So we jump back to repeat if we must.
			e.tmp.Words = append(e.tmp.Words, -3)
			e.tmp.Words = append(e.tmp.Words, float64(e.doOpen[len(e.doOpen)-1]+3))

			// We've matched the do-loop pair - drop the
			// open-reference
			e.doOpen = e.doOpen[:len(e.doOpen)-1]

		}

		// output a string, in compiled form
		if token.Name == ".\"" {
			e.strings = append(e.strings, token.Value)
			e.tmp.Words = append(e.tmp.Words, -5)
			e.tmp.Words = append(e.tmp.Words, float64(len(e.strings))-1)
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

	// OK so we're compiling something - and it wasn't a word
	// was it a variable?
	idx = e.findVariable(tok)
	if idx >= 0 {
		// compile this into something that will push
		// the offset of the variable onto the stack
		e.tmp.Words = append(e.tmp.Words, -1)
		e.tmp.Words = append(e.tmp.Words, float64(idx))
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

// dumpWord dumps the definition of the given word.
func (e *Eval) dumpWord(idx int) {

	// Lookup the word
	word := e.Dictionary[idx]

	// Store temporary data here
	codes := []string{}

	// Walk over the opcodes in the word-definition
	off := 0
	for off < len(word.Words) {

		// Get the actual byte
		//
		// Values >=0 are references to other words.
		//
		// Values <0 are "magic", and were created via
		// the "compilation" process.
		v := word.Words[int(off)]

		if v == -1 {
			codes = append(codes, fmt.Sprintf("%d: store %f", off, word.Words[off+1]))
			off++
		} else if v == -3 {
			codes = append(codes, fmt.Sprintf("%d: [cond-jmp %f]", off, word.Words[off+1]))
			off++
		} else if v == -4 {
			codes = append(codes, fmt.Sprintf("%d: [jmp %f]", off, word.Words[off+1]))
			off++
		} else if v == -5 {
			codes = append(codes, fmt.Sprintf("%d: [print-string %f (\"%s\")]", off, word.Words[off+1], e.strings[int(word.Words[off+1])]))
			off++
		} else if v == -10 {
			codes = append(codes, fmt.Sprintf("%d: [new-loop]", off))
			off++
		} else if v == -11 {
			codes = append(codes, fmt.Sprintf("%d: [loop-test]", off))
			off++
		} else {
			codes = append(codes, fmt.Sprintf("%d: %s", off, e.Dictionary[int(v)].Name))
		}

		// keep going for further words
		off++
	}

	// Didn't decompile?  Then it was a native-word
	if len(codes) == 0 {
		fmt.Printf("Word '%s' - [Native]\n", word.Name)
	} else {
		// Otherwise show the bytecode.
		fmt.Printf("Word '%s'\n %s\n", word.Name, strings.Join(codes, "\n "))
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
//    "-3" is a conditional-jump, which will change our IP if
//         the topmost item on the stack is "0".
//
//    "-4" is an unconditional jump, which will change our IP
//
//    "-5" prints a string, stored in our literal-area.
//         Dynamic strings are not supported.
//
//    "-10" creates a new Loop structure.
//          (i.e. `do`).
//
//    "-11" handles the test/termination of a loop condition.
//          (i.e. `loop`).
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
		} else if state == "string-print" {
			// print a string
			e.printString(e.strings[int(opcode)])
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
		} else if state == "new-loop" {

			// given the two-values on the stack
			// create and save a new Loop structure
			// to describe this loop.
			cur, err := e.Stack.Pop()
			if err != nil {
				return err
			}

			max, err2 := e.Stack.Pop()
			if err2 != nil {
				return err2
			}

			// new loop
			l := Loop{
				Start:   int(cur),
				Max:     int(max),
				Current: int(cur),
			}

			// save it away
			e.loops = append(e.loops, l)
			state = "default"
		} else if state == "loop-test" {

			// we've working with the last loop
			l := len(e.loops) - 1

			// bump the count
			e.loops[l].Current++

			// test to see if the loop is over
			if e.loops[l].Current >= e.loops[l].Max {
				e.Stack.Push(1)

				// loop is over now
				e.loops = e.loops[:len(e.loops)-1]
			} else {
				e.Stack.Push(0)
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

			switch opcode {
			case -1:
				state = "add-number"
			case -3:
				state = "cond-jump"
			case -4:
				state = "jump"
			case -5:
				state = "string-print"
			case -10:
				state = "new-loop"
			case -11:
				state = "loop-test"
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

// findVariable returns the index of the specified variable in our list
// of variables.
//
// Returns -1 if the variable cannot be found.
//
// Yes we store these in an array, rather than a map.  That's because
// we want to differentiate between an undefined variable and one with
// no value.
func (e *Eval) findVariable(name string) int {
	for i, v := range e.vars {
		if v.Name == name {
			return i
		}
	}
	return -1
}

// findWords returns the index of the specified word in our dictionary.
//
// Returns -1 if the word cannot be found.
func (e *Eval) findWord(name string) int {
	name = strings.ToLower(name)

	for index, word := range e.Dictionary {
		if name == word.Name {
			return index
		}
	}
	return -1
}

// printNumber - outputs a floating-point number.  However if the
// value is actually an integer then that is displayed instead.
func (e *Eval) printNumber(n float64) {

	// If the value on the top of the stack is an integer
	// then show it as one - i.e. without any ".00000".
	if float64(int(n)) == n {
		e.printString(fmt.Sprintf("%d", int(n)))
		return
	}

	// OK we have a floating-point result.  Show it, but
	// remove any trailing "0".
	//
	// This means we get 1.25 instead of 1.2500000 shown
	// when the user runs `5 4 / .`.
	//
	output := fmt.Sprintf("%f", n)

	for strings.HasSuffix(output, "0") {
		output = strings.TrimSuffix(output, "0")
	}
	e.printString(fmt.Sprintf("%s", output))
}

// printString outputs a string - replacing "\n", etc, with the
// real codes.
func (e *Eval) printString(str string) {

	str = strings.ReplaceAll(str, "\\n", "\n")
	str = strings.ReplaceAll(str, "\\t", "\t")
	str = strings.ReplaceAll(str, "\\r", "\r")

	if e.STDOUT == nil {
		e.STDOUT = bufio.NewWriter(os.Stdout)
	}

	e.STDOUT.WriteString(str)
	e.STDOUT.Flush()
}
