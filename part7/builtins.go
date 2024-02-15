// This file contains the built-in facilities we have hard-coded.
//
// That means the implementation for "+", "-", "/", "*", and "print".
//
// We've added `emit` here, to output the value at the top of the stack
// as an ASCII character, as well as "do" (nop) and "loop".

package main

import (
	"fmt"
	"sort"
	"strings"
)

func (e *Eval) add() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(a + b)
}

func (e *Eval) div() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(b / a)
}

func (e *Eval) do() {
	// nop
}

func (e *Eval) drop() {
	e.Stack.Pop()
}

func (e *Eval) dup() {
	a := e.Stack.Pop()
	e.Stack.Push(a)
	e.Stack.Push(a)
}

func (e *Eval) emit() {
	a := e.Stack.Pop()
	fmt.Printf("%c", rune(a))
}

func (e *Eval) eq() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	if a == b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
}

func (e *Eval) gt() {
	b := e.Stack.Pop()
	a := e.Stack.Pop()
	if a > b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
}

func (e *Eval) gtEq() {
	b := e.Stack.Pop()
	a := e.Stack.Pop()
	if a >= b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
}

func (e *Eval) iff() {
	// nop
}

func (e *Eval) invert() {
	v := e.Stack.Pop()
	if v == 0 {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
}

func (e *Eval) loop() {
	cur := e.Stack.Pop()
	max := e.Stack.Pop()

	cur++

	e.Stack.Push(max)
	e.Stack.Push(cur)
}

func (e *Eval) lt() {
	b := e.Stack.Pop()
	a := e.Stack.Pop()
	if a < b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
}

func (e *Eval) ltEq() {
	b := e.Stack.Pop()
	a := e.Stack.Pop()
	if a <= b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
}

func (e *Eval) mul() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(a * b)
}

func (e *Eval) print() {
	a := e.Stack.Pop()

	// If the value on the top of the stack is an integer
	// then show it as one - i.e. without any ".00000".
	if float64(int(a)) == a {
		fmt.Printf("%d\n", int(a))
		return
	}

	// OK we have a floating-point result.  Show it, but
	// remove any trailing "0".
	//
	// This means we get 1.25 instead of 1.2500000 shown
	// when the user runs `5 4 / .`.
	//
	output := fmt.Sprintf("%f", a)
	for strings.HasSuffix(output, "0") {
		output = strings.TrimSuffix(output, "0")
	}
	fmt.Printf("%s\n", output)

}

// startDefinition moves us into compiling-mode
//
// Note the interpreter handles removing this when it sees ";"
func (e *Eval) startDefinition() {
	e.compiling = true
}

// strings
func (e *Eval) stringCount() {
	// Return the number of strings we've seen
	e.Stack.Push(float64(len(e.strings)))
}

// strlen
func (e *Eval) strlen() {
	addr := e.Stack.Pop()
	i := int(addr)

	if i < len(e.strings) {
		str := e.strings[i]
		e.Stack.Push(float64(len(str)))
	} else {
		e.Stack.Push(-1)
	}
}

// strprn - string printing
func (e *Eval) strprn() {
	addr := e.Stack.Pop()
	i := int(addr)

	if i < len(e.strings) {
		str := e.strings[i]
		fmt.Printf("%s", str)
	}
}

func (e *Eval) sub() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(b - a)
}

func (e *Eval) swap() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(a)
	e.Stack.Push(b)
}

func (e *Eval) then() {
	// nop
}

func (e *Eval) words() {
	known := []string{}

	for _, entry := range e.Dictionary {
		known = append(known, entry.Name)
	}

	sort.Strings(known)
	fmt.Printf("%s\n", strings.Join(known, " "))
}
