// This file contains the built-in facilities we have hard-coded.
//
// That means the implementation for "+", "-", "/", "*", and "print".
//
// We've added `emit` here, to output the value at the top of the stack
// as an ASCII character, as well as "do" (nop) and "loop".

package main

import (
	"fmt"
)

func (e *Eval) add() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(a + b)
}

func (e *Eval) sub() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(b - a)
}

func (e *Eval) mul() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(a * b)
}

func (e *Eval) div() {
	a := e.Stack.Pop()
	b := e.Stack.Pop()
	e.Stack.Push(b / a)
}

func (e *Eval) print() {
	a := e.Stack.Pop()
	fmt.Printf("%f\n", a)
}

func (e *Eval) dup() {
	a := e.Stack.Pop()
	e.Stack.Push(a)
	e.Stack.Push(a)
}

// startDefinition moves us into compiling-mode
//
// Note the interpreter handles removing this when it sees ";"
func (e *Eval) startDefinition() {
	e.compiling = true
}

func (e *Eval) emit() {
	a := e.Stack.Pop()
	fmt.Printf("%c", rune(a))
}

func (e *Eval) do() {
	// nop
}

func (e *Eval) loop() {
	cur := e.Stack.Pop()
	max := e.Stack.Pop()

	cur += 1

	e.Stack.Push(max)
	e.Stack.Push(cur)
}
