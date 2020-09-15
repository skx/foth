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
		": test_hot  0 > if star then star ;"}

	for _, str := range tests {
		e.Eval(strings.Split(str, " "))
	}

	e.dumpWords()
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
