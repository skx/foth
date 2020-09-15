package eval

import (
	"testing"
)

func TestAdd(t *testing.T) {

	e := New()

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

	e := New()

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

	e := New()
	if e.do() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestDrop(t *testing.T) {

	e := New()
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

	e := New()
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

	e := New()
	if e.emit() == nil {
		t.Fatalf("expected error, got none")
	}

	e.Stack.Push(12)
	if e.emit() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestEq(t *testing.T) {

	e := New()

	// empty stack
	err := e.eq()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.eq()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(1)
	e.Stack.Push(1)
	err = e.eq()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// equal
	x, _ := e.Stack.Pop()
	if x != 1 {
		t.Fatalf("eq() failed")
	}

	// two items
	e.Stack.Push(12)
	e.Stack.Push(1)
	err = e.eq()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// non-equal
	x, _ = e.Stack.Pop()
	if x != 0 {
		t.Fatalf("eq() failed")
	}
}

func TestGt(t *testing.T) {

	e := New()

	// empty stack
	err := e.gt()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.gt()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(10)
	e.Stack.Push(1)
	err = e.gt()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// gt
	x, _ := e.Stack.Pop()
	if x != 1 {
		t.Fatalf("gt() failed")
	}

	// two items
	e.Stack.Push(1)
	e.Stack.Push(1)
	err = e.gt()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// not-gt
	x, _ = e.Stack.Pop()
	if x != 0 {
		t.Fatalf("gt() failed")
	}
}

func TestGtEq(t *testing.T) {

	e := New()

	// empty stack
	err := e.gtEq()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.gtEq()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(1)
	e.Stack.Push(1)
	err = e.gtEq()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// gt
	x, _ := e.Stack.Pop()
	if x != 1 {
		t.Fatalf(">=() failed")
	}

	// two items
	e.Stack.Push(-1)
	e.Stack.Push(1)
	err = e.gtEq()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// not->=
	x, _ = e.Stack.Pop()
	if x != 0 {
		t.Fatalf("gt() failed")
	}
}

func TestInvert(t *testing.T) {

	e := New()

	// empty stack
	err := e.invert()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// 0 -> 1
	e.Stack.Push(0)
	e.invert()
	out, err2 := e.Stack.Pop()
	if err2 != nil {
		t.Errorf("unexpected error")
	}
	if out != 1 {
		t.Errorf("unexpected result")
	}

	// 10 -> 0
	e.Stack.Push(10)
	e.invert()
	out, err2 = e.Stack.Pop()
	if err2 != nil {
		t.Errorf("unexpected error")
	}
	if out != 0 {
		t.Errorf("unexpected result")
	}

}
func TestIff(t *testing.T) {

	e := New()
	if e.iff() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestThen(t *testing.T) {

	e := New()
	if e.then() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestWords(t *testing.T) {

	e := New()
	if e.words() != nil {
		t.Fatalf("unexpected error")
	}
}
