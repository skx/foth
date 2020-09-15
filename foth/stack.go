package main

import (
	"fmt"
)

// Stack holds our numbers.
type Stack []float64

// IsEmpty checks if our stack is empty.
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push adds a new number to the stack
func (s *Stack) Push(x float64) {
	*s = append(*s, x)
}

// Pop removes and returns the top element of stack.
func (s *Stack) Pop() (float64, error) {
	if s.IsEmpty() {
		return 0, fmt.Errorf("stack underflow")
	}

	i := len(*s) - 1
	x := (*s)[i]
	*s = (*s)[:i]

	return x, nil
}
