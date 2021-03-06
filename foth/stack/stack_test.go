package stack

import (
	"testing"
)

func TestEmpty(t *testing.T) {

	var s Stack

	if !s.IsEmpty() {
		t.Fatalf("new stack should be empty")
	}

	s.Push(33)
	if s.IsEmpty() {
		t.Fatalf("populated stack should not be empty")
	}
	if s.Len() != 1 {
		t.Fatalf("stack length was wrong")
	}
	if s.At(0) != 33 {
		t.Fatalf("stack value was wrong")
	}
	s.Pop()
	if !s.IsEmpty() {
		t.Fatalf("stack should be back to new")
	}

}

func TestUnderflow(t *testing.T) {

	var s Stack

	_, err := s.Pop()

	if err == nil {
		t.Fatalf("stack was not covered")
	}
}
