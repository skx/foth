package main

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
