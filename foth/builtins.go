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

func (e *Eval) add() error {
	var a, b float64
	var err error

	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	e.Stack.Push(a + b)
	return nil
}

func (e *Eval) div() error {
	var a, b float64
	var err error

	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}

	e.Stack.Push(b / a)
	return nil
}

func (e *Eval) do() error {
	// nop
	return nil
}

func (e *Eval) drop() error {
	_, err := e.Stack.Pop()
	return err
}

func (e *Eval) dup() error {
	a, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	e.Stack.Push(a)
	e.Stack.Push(a)

	return nil
}

func (e *Eval) emit() error {
	a, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	fmt.Printf("%c", rune(a))
	return nil
}

func (e *Eval) eq() error {
	var a, b float64
	var err error
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	if a == b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
	return nil
}

func (e *Eval) gt() error {
	var a, b float64
	var err error
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	if a > b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
	return nil
}

func (e *Eval) gtEq() error {
	var a, b float64
	var err error
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	if a >= b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
	return nil
}

func (e *Eval) iff() error {
	// nop
	return nil
}

func (e *Eval) invert() error {
	v, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	if v == 0 {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}

	return nil
}

func (e *Eval) loop() error {
	var cur, max float64
	var err error
	cur, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	max, err = e.Stack.Pop()
	if err != nil {
		return err
	}

	cur++

	e.Stack.Push(max)
	e.Stack.Push(cur)

	return nil
}

func (e *Eval) lt() error {
	var a, b float64
	var err error
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}

	if a < b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}

	return nil
}

func (e *Eval) ltEq() error {
	var a, b float64
	var err error
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}

	if a <= b {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}

	return nil
}

func (e *Eval) mul() error {
	var a, b float64
	var err error
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	e.Stack.Push(a * b)
	return nil
}

func (e *Eval) print() error {
	a, err := e.Stack.Pop()
	if err != nil {
		return err
	}

	// If the value on the top of the stack is an integer
	// then show it as one - i.e. without any ".00000".
	if float64(int(a)) == a {
		fmt.Printf("%d\n", int(a))
		return nil
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
	return nil
}

// startDefinition moves us into compiling-mode
//
// Note the interpreter handles removing this when it sees ";"
func (e *Eval) startDefinition() error {
	e.compiling = true
	return nil
}

func (e *Eval) sub() error {
	var a, b float64
	var err error
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	e.Stack.Push(b - a)
	return nil
}

func (e *Eval) swap() error {
	var a, b float64
	var err error
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	e.Stack.Push(a)
	e.Stack.Push(b)

	return nil
}

func (e *Eval) then() error {
	// nop
	return nil
}

func (e *Eval) words() error {
	known := []string{}

	for _, entry := range e.Dictionary {
		known = append(known, entry.Name)
	}

	sort.Strings(known)
	fmt.Printf("%s\n", strings.Join(known, " "))

	return nil
}
