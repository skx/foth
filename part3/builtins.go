// This file contains the built-in facilities we have hard-coded.
//
// That means the implementation for "+", "-", "/", "*", and "print".
//
// Compared to part3 we've added a word for ":" to allow defining words
// via manipulation of the "compiling" flag.

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

// startDefinition moves us into compiling-mode.
//
// Note the interpreter handles removing this when it sees ";"
func (e *Eval) startDefinition() {
	e.compiling = true
}
