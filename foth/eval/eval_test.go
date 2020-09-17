package eval

import (
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

