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

func TestDebug(t *testing.T) {

	e := New()

	// debug is off
	err := e.debugp()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	var val float64
	val, err = e.Stack.Pop()
	if err != nil {
		t.Errorf("unexpected error")
	}
	if val != 0 {
		t.Fatalf("debug value is wrong")
	}

	// empty stack
	err = e.debugSet()
	if err == nil {
		t.Fatalf("expected error with empty stack; got none")
	}

	// set debug
	e.Stack.Push(1)
	err = e.debugSet()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// Getting it should confirm it is set.
	err = e.debugp()
	if err != nil {
		t.Fatalf("unexpected error")
	}
	val, err = e.Stack.Pop()
	if err != nil {
		t.Errorf("unexpected error")
	}
	if val != 1 {
		t.Fatalf("debug value is wrong")
	}

	// Now set it off
	e.Stack.Push(0)
	err = e.debugSet()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// and it should be off
	err = e.debugp()
	if err != nil {
		t.Fatalf("unexpected error")
	}
	val, err = e.Stack.Pop()
	if err != nil {
		t.Errorf("unexpected error")
	}
	if val != 0 {
		t.Fatalf("debug value is wrong")
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

func TestDump(t *testing.T) {

	e := New()

	// test with empty stack
	err := e.dump()
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	// count the words
	if e.wordLen() != nil {
		t.Fatalf("unexpected error")
	}

	count, err := e.Stack.Pop()
	if err != nil {
		t.Fatalf("stack error")
	}

	// for each word
	cur := 0
	for cur < int(count) {

		e.Stack.Push(float64(cur))
		e.dump()

		if !e.Stack.IsEmpty() {
			t.Fatalf("stack surprise")
		}
		cur++
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

func TestLoop(t *testing.T) {

	e := New()

	// empty stack
	err := e.loop()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	e.Stack.Push(10)
	e.loop()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// Two values
	e.Stack.Push(4)
	e.Stack.Push(54)

	e.loop()
	a, err := e.Stack.Pop()
	if a != 55 {
		t.Fatalf("unexpected result, got %f", a)
	}
	if err != nil {
		t.Fatalf("unexpected error")
	}
	b, err2 := e.Stack.Pop()
	if b != 4 {
		t.Fatalf("unexpected result, got %f", b)
	}
	if err2 != nil {
		t.Fatalf("unexpected error")
	}
}

func TestLt(t *testing.T) {

	e := New()

	// empty stack
	err := e.lt()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.lt()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(10)
	e.Stack.Push(22)
	err = e.lt()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// lt
	x, _ := e.Stack.Pop()
	if x != 1 {
		t.Fatalf("lt() failed")
	}

	// two items
	e.Stack.Push(1)
	e.Stack.Push(1)
	err = e.lt()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// not-lt
	x, _ = e.Stack.Pop()
	if x != 0 {
		t.Fatalf("lt() failed")
	}
}

func TestLtEq(t *testing.T) {

	e := New()

	// empty stack
	err := e.ltEq()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.ltEq()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(1)
	e.Stack.Push(1)
	err = e.ltEq()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// lt
	x, _ := e.Stack.Pop()
	if x != 1 {
		t.Fatalf("<=() failed")
	}

	// two items
	e.Stack.Push(10)
	e.Stack.Push(-22)
	err = e.ltEq()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	// not-<=
	x, _ = e.Stack.Pop()
	if x != 0 {
		t.Fatalf("<=() failed")
	}
}

func TestMul(t *testing.T) {

	e := New()

	// empty stack
	err := e.mul()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.mul()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(6)
	e.Stack.Push(7)
	err = e.mul()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	x, _ := e.Stack.Pop()
	if x != 42 {
		t.Fatalf("wrong result for mul")
	}
}

func TestNop(t *testing.T) {

	e := New()
	if e.nop() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestOver(t *testing.T) {

	e := New()

	// empty stack
	err := e.over()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// int
	e.Stack.Push(3)
	err = e.over()
	if err == nil {
		t.Fatalf("expected underflow")
	}

	// OK now we add 2 1
	e.Stack.Push(2)
	e.Stack.Push(1)
	err = e.over()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// stack should now be: 2 1 2
	v, er := e.Stack.Pop()
	if er != nil {
		t.Fatalf("unexpected error")
	}
	if v != 2 {
		t.Fatalf("unexpected error")
	}

	v, er = e.Stack.Pop()
	if er != nil {
		t.Fatalf("unexpected error")
	}
	if v != 1 {
		t.Fatalf("unexpected error")
	}

	v, er = e.Stack.Pop()
	if er != nil {
		t.Fatalf("unexpected error")
	}
	if v != 2 {
		t.Fatalf("unexpected error")
	}

	if !e.Stack.IsEmpty() {
		t.Fatalf("unexpected stack content")
	}

}
func TestPrint(t *testing.T) {
	e := New()

	// empty stack
	err := e.print()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// int
	e.Stack.Push(3)
	err = e.print()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if !e.Stack.IsEmpty() {
		t.Fatalf("stack should be empty now")

	}

	f := 5 / 4.0

	// float
	e.Stack.Push(f)
	a, _ := e.Stack.Pop()
	if a != 1.25 {
		t.Fatalf("not a float?: %f", a)
	}
	e.Stack.Push(f)
	err = e.print()
	if err != nil {
		t.Fatalf("unexpected error")
	}

}

func TestSub(t *testing.T) {

	e := New()

	// empty stack
	err := e.sub()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.sub()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(6)
	e.Stack.Push(4)
	err = e.sub()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	x, _ := e.Stack.Pop()
	if x != 2 {
		t.Fatalf("wrong result for sub, got %f", x)
	}
}

func TestSwap(t *testing.T) {

	e := New()

	// empty stack
	err := e.swap()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.swap()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// two items
	e.Stack.Push(3)
	e.Stack.Push(1)
	err = e.swap()
	if err != nil {
		t.Fatalf("expected no error, but got one")
	}

	a, err := e.Stack.Pop()
	if err != nil {
		t.Fatalf("unexpected error")
	}
	b, err2 := e.Stack.Pop()
	if err2 != nil {
		t.Fatalf("unexpected error")
	}
	if a != 3 {
		t.Fatalf("unexpected error, got %f", a)
	}
	if b != 1 {
		t.Fatalf("unexpected error, got %f", b)
	}
}

func TestWords(t *testing.T) {

	e := New()
	if e.words() != nil {
		t.Fatalf("unexpected error")
	}
}

func TestWordLen(t *testing.T) {

	e := New()
	if e.wordLen() != nil {
		t.Fatalf("unexpected error")
	}

	if e.Stack.IsEmpty() {
		t.Fatalf("expected stack entry")
	}

	count, err := e.Stack.Pop()
	if err != nil {
		t.Fatalf("stack error")
	}
	if int(count) > len(e.Dictionary) {
		t.Fatalf("too many wrds")
	}
}
