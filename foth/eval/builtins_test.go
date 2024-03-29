package eval

import (
	"os"
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

func TestGetVar(t *testing.T) {

	e := New()

	// Need to have a variable set before it can be retrieved
	e.SetVariable("foo", 93.2)

	// empty stack
	err := e.getVar()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// get the first variable.
	e.Stack.Push(0)
	err = e.getVar()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// The value should match
	val, err := e.Stack.Pop()
	if err != nil {
		t.Fatalf("stack underflow")
	}

	if val != 93.2 {
		t.Fatalf("getvar('foo') had %f not %f", val, 93.2)
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

func TestMod(t *testing.T) {

	e := New()

	// empty stack
	err := e.mod()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.mod()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	type TestCase struct {
		in  float64
		out float64
	}

	tests := []TestCase{{in: 1, out: 1},
		{in: 2, out: 2},
		{in: 3, out: 3},
		{in: 4, out: 0},
		{in: 5, out: 1},
		{in: 6, out: 2},
		{in: 7, out: 3},
		{in: 8, out: 0},
		{in: 9, out: 1},
		{in: 10, out: 2},
	}

	for _, test := range tests {
		// two items
		e.Stack.Push(test.in)
		e.Stack.Push(4)
		err = e.mod()
		if err != nil {
			t.Fatalf("expected no error, but got one")
		}

		x, _ := e.Stack.Pop()
		if x != test.out {
			t.Fatalf("wrong result %f %% 4.  Got %f, not %f", test.in, x, test.out)
		}
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

func TestSetVar(t *testing.T) {

	e := New()

	// Setup a variable
	e.SetVariable("name", 6)

	// empty stack
	err := e.setVar()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// only one item
	e.Stack.Push(1)
	err = e.setVar()
	if err == nil {
		t.Fatalf("expected error with empty stack")
	}

	// Now set
	e.Stack.Push(32.1)
	e.Stack.Push(0)
	err = e.setVar()
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// confirm it worked
	v, err := e.GetVariable("name")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if v != 32.1 {
		t.Fatalf("value mismatch after setting variable")
	}

}

// strings counts the string literals, and will return
// a number on the stack
func TestStrings(t *testing.T) {

	e := New()

	// Empty stack on first start
	n := e.Stack.Len()
	if n != 0 {
		t.Fatalf("failing result for stringCount, got %d", n)
	}

	// call the function
	err := e.stringCount()
	if err != nil {
		t.Fatalf("unexpected error with stringCount %s", err.Error())
	}

	n = e.Stack.Len()
	if n != 1 {
		t.Fatalf("failing result for stringCount, got %d", n)
		os.Exit(1)
	}
}

func TestStrlen(t *testing.T) {

	e := New()
	e.strings = append(e.strings, "Steve")

	// call the function
	err := e.strlen()
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	// Empty stack on first start
	n := e.Stack.Len()
	if n != 0 {
		t.Fatalf("failing result for strlen, got %d", n)
	}

	// push an invalid string
	e.Stack.Push(100.0)

	// call the function
	err = e.strlen()
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	// Now try to get the length of Steve
	e.Stack.Push(0.0)
	err = e.strlen()
	if err != nil {
		t.Fatalf("unexpected error, calling strlen %s", err.Error())
	}

	// Is the result expected?
	x, _ := e.Stack.Pop()
	if x != 5 {
		t.Fatalf("wrong result for strlen, got %f", x)
	}
}

func TestStrPrn(t *testing.T) {

	e := New()

	// We want to avoid spamming stdout, so our string to print is "empty"
	e.strings = append(e.strings, "")

	// call the function
	err := e.strprn()
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	// Empty stack on first start
	n := e.Stack.Len()
	if n != 0 {
		t.Fatalf("failing result for strprn, got %d", n)
	}

	// push an invalid string
	e.Stack.Push(100.0)

	// call the function
	err = e.strprn()
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	// Now try to get the length of Steve
	e.Stack.Push(0.0)
	err = e.strprn()
	if err != nil {
		t.Fatalf("unexpected error, calling strprn %s", err.Error())
	}

	// Empty stack, still?
	n = e.Stack.Len()
	if n != 0 {
		t.Fatalf("failing result for strprn, got %d", n)
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
