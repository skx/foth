// Package stack allows a stack of float64 to be maintained.
package stack

import (
	"fmt"
)

// Stack holds our numbers.
type Stack []float64

// At returns the value at the given offset.
func (s *Stack) At(offset int) float64 {
	return (*s)[offset]
}

// IsEmpty checks whether the stack is empty.
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Len returns the length of the stack.
func (s *Stack) Len() int {
	return len(*s)
}

// Push adds a new number to the top of the stack.
func (s *Stack) Push(x float64) {
	*s = append(*s, x)
}

// Pop removes, and returns, the top element of stack.
func (s *Stack) Pop() (float64, error) {
	if s.IsEmpty() {
		return 0, fmt.Errorf("stack underflow")
	}

	i := len(*s) - 1
	x := (*s)[i]
	*s = (*s)[:i]

	return x, nil
}
