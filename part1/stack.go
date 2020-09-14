package main

import (
	"fmt"
	"os"
)

// Stack holds our numbers.
type Stack []float64

// IsEmpty checks if the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push adds a new number to the stack
func (s *Stack) Push(x float64) {
	*s = append(*s, x)
}

// Pop removes and returns the top element of stack.
func (s *Stack) Pop() float64 {
	if s.IsEmpty() {
		fmt.Printf("stack underflow\n")
		os.Exit(1)
	}

	i := len(*s) - 1
	x := (*s)[i]
	*s = (*s)[:i]

	return x
}
