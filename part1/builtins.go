// This file contains the built-in facilities we have hard-coded.
//
// That means the implementation for "+", "-", "/", "*", and "print".

package main

import "fmt"

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
