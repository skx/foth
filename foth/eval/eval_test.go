package eval

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {

	e := New()
	e.debug = true

	// stack-underflow
	out := e.Eval("  . ")
	if out == nil {
		t.Fatalf("expected error, got none")
	}

	// nested-comments
	out = e.Eval(" ( one ( two ) ) ")
	if out == nil {
		t.Fatalf("expected error, got none")
	}

	// no error
	out = e.Eval(".\" Hello, World \" ")
	if out != nil {
		t.Fatalf("unexpected")
	}
}

func TestClearWords(t *testing.T) {

	// create instance
	e := New()
	e.debug = true

	// Push some stuff
	err := e.Eval("1 3 4 5")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// Ensure it is non-empty
	if e.Stack.IsEmpty() {
		t.Fatalf("unexpected stack")
	}
	if e.Stack.Len() != 4 {
		t.Fatalf("unexpected stack")
	}

	// Clear the stack
	err = e.Eval(".s clearstack")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// Ensure it is non-empty
	if !e.Stack.IsEmpty() {
		t.Fatalf("unexpected stack")
	}
	if e.Stack.Len() != 0 {
		t.Fatalf("unexpected stack")
	}

}
func TestDumpWords(t *testing.T) {

	// dummy test
	os.Setenv("DEBUG", "1")
	e := New()

	if e.debug != true {
		t.Fatalf("putenv didn't enable debugging")
	}

	// test definitions
	tests := []string{": star 42 emit ;",
		": stars 0 do star loop 10 emit ;",
		": test_hot  0 > if star then star ;",
		": tests 0 0 = if 1 else 2 then ;",
		": tests 0 0 = if .\" test \" else .\" ok\" ;",
		"0 0 = if .\" test \" else .\" ok\"",
	}

	for _, str := range tests {
		e.Eval(str)
	}

	e.dumpWord(0)
	os.Setenv("DEBUG", "")
}

func TestError(t *testing.T) {

	// Some things that will generate errors
	tests := []string{": foo . ; foo",

		// two stack-items are expected for `do`
		": foo do i emit loop ; foo",
		": foo 2 do i emit loop ; foo",

		// i & m only within a loop-body
		": foo i ; foo",
		": foo m ; foo",
	}

	for _, test := range tests {
		e := New()
		e.debug = true

		err := e.Eval(test)
		if err == nil {
			t.Fatalf("expected error, got none for: %s", test)
		}
	}

}

// Try running one of each of our test-cases
func TestEvalWord(t *testing.T) {

	// dummy test
	e := New()
	e.debug = true

	// test definitions
	tests := []string{": star 42 emit ;",
		": stars 0 do star loop 10 emit ;",
		": test_hot  0 > if star then ;",

		// deliberately redefine a word
		": test_hot  0 >= if star then ;",
		"10 stars",
		"star",
		"10 test_hot",
		"-1 test_hot"}

	for _, str := range tests {
		e.Eval(str)
	}

}

func TestFloatFail(t *testing.T) {

	tests := []string{": foo. 3.2.1.2 emit ; foo.",
		"6.7.8.9 ."}

	for _, str := range tests {
		e := New()
		e.debug = true
		err := e.Eval(str)
		if err == nil {
			t.Fatalf("expected error processing '%s', got none", str)
		}
		if !strings.Contains(err.Error(), "failed to convert") {
			t.Fatalf("got an error, but the wrong one: %s", err.Error())
		}
	}
}

func TestIfThenElse(t *testing.T) {

	type Test struct {
		input  string
		result float64
	}

	tests := []Test{
		Test{input: "3 3 = if 1 then", result: 1},
		Test{input: "3 3 = if .\" ok \" 1 then", result: 1},
		Test{input: ": f 3 3 = if 1 then ; f", result: 1},
		Test{input: ": f 3 3 = if .\" ok \" 1 then ; f", result: 1},

		Test{input: "3 3 = invert if 1 else 2 then", result: 2},
		Test{input: ": f 3 3 = invert if 1 else 2 then ; f", result: 2},

		Test{input: "3 31 = if 0 else 3 then", result: 3},
		Test{input: ": f 3 31 = if 0 else 3 then ; f ", result: 3},

		Test{input: "3 31 = if 1 else 12 then", result: 12},
		Test{input: ": f 3 31 = if 1 else 12 then ; f", result: 12},

		Test{input: "3 21 = invert if 221 else 112 then", result: 221},
		Test{input: ": ff 3 21 = invert if 221 else 112 then ; ff ", result: 221},
	}

	for _, test := range tests {

		e := New()
		e.debug = true
		err := e.Eval(test.input)
		if err != nil {
			t.Fatalf("unexpected error processing '%s': %s", test.input, err.Error())
		}

		ret, err2 := e.Stack.Pop()
		if err2 != nil {
			t.Fatalf("failed to get stack value from %s", test.input)
		}
		if !e.Stack.IsEmpty() {
			t.Fatalf("%s: expected stack to be empty", test.input)
		}
		if ret != test.result {
			t.Fatalf("%s: %f got %f", test.input, test.result, ret)
		}
	}

	// Test of an if with an empty stack
	e := New()
	e.debug = true
	err := e.Eval(": fail if . then ; fail")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestMaxMin(t *testing.T) {

	errors := []string{
		// no entry
		"max",
		"min",

		// one too few
		"3 max",
		"3 min"}

	for _, txt := range errors {

		// create instance
		e := New()
		e.debug = true

		err := e.Eval(txt)
		if err == nil {
			t.Fatalf("expected an error, got none")
		}
		if !strings.Contains(err.Error(), "underflow") {
			t.Fatalf("found wrong error: %s", err.Error())
		}
	}

	type TestCase struct {
		Input  string
		Result float64
	}

	tests := []TestCase{
		{Input: "3 4 max", Result: 4},
		{Input: "4 3 max", Result: 4},
		{Input: "3 3 max", Result: 3},
		{Input: "-4 -30 max", Result: -4},

		{Input: "3 4 min", Result: 3},
		{Input: "4 3 min", Result: 3},
		{Input: "2 2 min", Result: 2},
		{Input: "-4 -30 min", Result: -30},
	}

	for _, test := range tests {

		// create instance
		e := New()
		e.debug = true

		err := e.Eval(test.Input)
		if err != nil {
			t.Fatalf("unexpected error")
		}

		ret, err2 := e.Stack.Pop()
		if err2 != nil {
			t.Fatalf("failed to get stack value")
		}
		if !e.Stack.IsEmpty() {
			t.Fatalf("expected stack to be empty, it wasn't")
		}
		if ret != test.Result {
			t.Fatalf("unexpected result for %s -> %f", test.Input, ret)
		}
	}
}

func TestReset(t *testing.T) {

	// create instance
	e := New()
	e.debug = true

	// Trigger error-state
	err := e.Eval(": foo 1 3 + .;")
	if err == nil {
		t.Fatalf("expected an error, got none")
	}

	// Ensure something is on the stack
	e.Stack.Push(33.1)

	// Now reset and run something else to confirm it worked
	e.Reset()

	// The stack should be empty
	if !e.Stack.IsEmpty() {
		t.Fatalf("expected stack to be empty, it wasn't")
	}

	err = e.Eval(": foo 1 3 + ; foo ")
	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}

	ret, err2 := e.Stack.Pop()
	if err2 != nil {
		t.Fatalf("failed to get stack value")
	}
	if !e.Stack.IsEmpty() {
		t.Fatalf("expected stack to be empty, it wasn't")
	}
	if ret != 4 {
		t.Fatalf("unexpected result, post-recovery")
	}
}

func TestVariables(t *testing.T) {

	// create instance
	e := New()
	e.debug = true

	var v float64
	var err error

	// fetching a variable will fail, as it is not present
	_, err = e.GetVariable("unset")
	if err == nil {
		t.Fatalf("expected error accessing a missing variable")
	}

	// Now set it
	e.SetVariable("unset", 22)
	e.SetVariable("unset", 33)

	v, err = e.GetVariable("unset")
	if err != nil {
		t.Fatalf("unexpected error accessing variable")
	}
	if v != 33 {
		t.Fatalf("variable has wrong value")
	}

	// Run a script to change the variable
	err = e.Eval("12 unset !")
	if err != nil {
		t.Fatalf("error running script")
	}

	// Get the value and confirm it is updated
	v, err = e.GetVariable("unset")
	if err != nil {
		t.Fatalf("unexpected error accessing variable")
	}
	if v != 12 {
		t.Fatalf("variable has wrong value")
	}

	// Finally declare a variable and set the value
	err = e.Eval("variable meow")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	err = e.Eval("3 meow !")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	// Get the value and confirm it is updated
	v, err = e.GetVariable("meow")
	if err != nil {
		t.Fatalf("unexpected error accessing variable")
	}
	if v != 3 {
		t.Fatalf("variable has wrong value")
	}

	// Double the value of the variable
	err = e.Eval(": double meow @ 2 * meow ! ; double")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	v, err = e.GetVariable("meow")
	if err != nil {
		t.Fatalf("unexpected error accessing variable")
	}
	if v != 6 {
		t.Fatalf("variable has wrong value")
	}
}

func TestSetWriter(t *testing.T) {

	var b bytes.Buffer
	out := bufio.NewWriter(&b)

	e := New()
	e.debug = true
	e.SetWriter(out)

	// write something simple
	err := e.Eval("42 emit 42 emit")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if b.String() != "**" {
		t.Fatalf("STDOUT didn't match")
	}

	// write something more complex
	b.Reset()

	// i is the current loop index
	// m is the max
	//
	// so we're outputting "1/10", "2/10", etc.
	//
	err = e.Eval("10 0 do i 48 + emit 47 emit m . loop")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if b.String() != "0/10\n1/10\n2/10\n3/10\n4/10\n5/10\n6/10\n7/10\n8/10\n9/10\n" {
		t.Fatalf("STDOUT didn't match, got '%s' for ", b.String())
	}

	// Finally a string literal
	b.Reset()
	err = e.Eval(".\" Steve\nKemp \"")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if b.String() != "Steve\nKemp" {
		t.Fatalf("STDOUT didn't match, got '%s'", b.String())
	}
}
