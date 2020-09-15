package main

import (
	"testing"
)

func TestAdd(t *testing.T) {

	e := NewEval()

	// empty stack
	err := e.add()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.add()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(1.2)
	e.Stack.Push(1.1)
	err = e.add()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	x, _ := e.Stack.Pop()
	if x != 2.3 {
		t.Fatalf("wrong result for add")
	}
}

func TestDiv(t *testing.T) {

	e := NewEval()

	// empty stack
	err := e.div()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.div()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(9)
	e.Stack.Push(3)
	err = e.div()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	x, _ := e.Stack.Pop()
	if x != 3 {
		t.Fatalf("wrong result for add")
	}
}

func TestDo(t *testing.T) {

	e := NewEval()
	if e.do() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestDrop(t *testing.T) {

	e := NewEval()
	if e.drop() == nil {
		t.Fatalf("expected error, got none")
	}

	e.Stack.Push(12)
	if e.drop() != nil {
		t.Fatalf("unexpected error")
	}
	if !e.Stack.IsEmpty() {
		t.Fatalf("stack should be empty now")
	}
}

func TestDup(t *testing.T) {

	e := NewEval()
	if e.dup() == nil {
		t.Fatalf("expected error, got none")
	}

	e.Stack.Push(12)
	if e.dup() != nil {
		t.Fatalf("unexpected error")
	}

	// should now have two entries
	_, err := e.Stack.Pop()
	if err != nil {
		t.Errorf("unexpected error")
	}
	_, err2 := e.Stack.Pop()
	if err2 != nil {
		t.Errorf("unexpected error")
	}
	if !e.Stack.IsEmpty() {
		t.Fatalf("stack should be empty now")
	}
}

func TestEmit(t *testing.T) {

	e := NewEval()
	if e.emit() == nil {
		t.Fatalf("expected error, got none")
	}

	e.Stack.Push(12)
	if e.emit() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestIff(t *testing.T) {

	e := NewEval()
	if e.iff() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestThen(t *testing.T) {

	e := NewEval()
	if e.then() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestWords(t *testing.T) {

	e := NewEval()
	if e.words() != nil {
		t.Fatalf("unexpected error")
	}
}
