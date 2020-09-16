package eval

import (
	"os"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {

	e := New()
	out := e.Eval([]string{"", ".", " "})

	if out == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestDumpWords(t *testing.T) {

	// dummy test
	os.Setenv("DEBUG", "1")
	e := New()

	if e.debug != true {
		t.Fatalf("putenv didn't enable debugging")
	}
	os.Setenv("DEBUG", "")

	// test definitions
	tests := []string{": star 42 emit ;",
		": stars 0 do star loop 10 emit ;",
		": test_hot  0 > if star then star ;",
		": tests 0 0 = if 1 else 2 then ;"}

	for _, str := range tests {
		e.Eval(strings.Split(str, " "))
	}

	e.dumpWord(0)
}

// Try running one of each of our test-cases
func TestEvalWord(t *testing.T) {

	// dummy test
	e := New()

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
		e.Eval(strings.Split(str, " "))
	}

}

func TestFloatFail(t *testing.T) {

	tests := []string{": foo. 3.2.1.2 emit ; foo.",
		"6.7.8.9 ."}

	for _, str := range tests {
		e := New()
		err := e.Eval(strings.Split(str, " "))
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
		Test{input: ": f 3 3 = if 1 then ; f", result: 1},

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
		err := e.Eval(strings.Split(test.input, " "))
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

}
