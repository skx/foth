// This file contains the built-in facilities we have hard-coded.
//
// That means the implementation for "+", "-", "/", "*", and "print".
//
// We've added `emit` here, to output the value at the top of the stack
// as an ASCII character, as well as "do" (nop) and "loop".

package eval

import (
	"fmt"
	"sort"
	"strings"
)

func (e *Eval) binOp(op func(float64, float64) float64) func() error {
	return func() error {
		a, err := e.Stack.Pop()
		if err != nil {
			return err
		}

		b, err := e.Stack.Pop()
		if err != nil {
			return err
		}

		e.Stack.Push(op(a, b))
		return nil
	}
}

func (e *Eval) add() error {
	return e.binOp(func(n float64, m float64) float64 { return n + m })()
}

func (e *Eval) clearStack() error {
	for !e.Stack.IsEmpty() {
		e.Stack.Pop()
	}
	return nil
}

func (e *Eval) debugSet() error {

	v, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	if v == 0 {
		e.debug = false
	} else {
		e.debug = true
	}

	return nil
}

func (e *Eval) debugp() error {
	if e.debug {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}
	return nil
}

func (e *Eval) div() error {
	return e.binOp(func(n float64, m float64) float64 { return m / n })()
}

func (e *Eval) drop() error {
	_, err := e.Stack.Pop()
	return err
}

func (e *Eval) dump() error {
	a, err := e.Stack.Pop()
	if err != nil {
		return err
	}

	e.dumpWord(int(a))
	return nil
}

func (e *Eval) dup() error {
	a, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	e.Stack.Push(a)
	e.Stack.Push(a)

	return nil
}

func (e *Eval) emit() error {
	a, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	e.printString(string(rune(a)))
	return nil
}

func (e *Eval) eq() error {
	return e.binOp(func(n float64, m float64) float64 {
		if m == n {
			return 1
		}
		return 0
	})()
}

func (e *Eval) getVar() error {

	offset, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	val := e.vars[int(offset)]
	e.Stack.Push(val.Value)
	return nil
}

func (e *Eval) gt() error {
	return e.binOp(func(n float64, m float64) float64 {
		if m > n {
			return 1
		}
		return 0
	})()
}

func (e *Eval) gtEq() error {
	return e.binOp(func(n float64, m float64) float64 {
		if m >= n {
			return 1
		}
		return 0
	})()
}

func (e *Eval) i() error {
	if len(e.loops) > 0 {
		i := e.loops[len(e.loops)-1].Current
		e.Stack.Push(float64(i))
		return nil
	}
	return fmt.Errorf("you cannot access 'i' outside a loop-body")
}

func (e *Eval) invert() error {
	v, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	if v == 0 {
		e.Stack.Push(1)
	} else {
		e.Stack.Push(0)
	}

	return nil
}

func (e *Eval) loop() error {
	return nil
}

func (e *Eval) lt() error {
	return e.binOp(func(n float64, m float64) float64 {
		if m < n {
			return 1
		}
		return 0
	})()
}

func (e *Eval) ltEq() error {
	return e.binOp(func(n float64, m float64) float64 {
		if m <= n {
			return 1
		}
		return 0
	})()
}

func (e *Eval) max() error {
	return e.binOp(func(n float64, m float64) float64 {
		if m > n {
			return m
		}
		return n
	})()
}

func (e *Eval) min() error {
	return e.binOp(func(n float64, m float64) float64 {
		if m < n {
			return m
		}
		return n
	})()
}

func (e *Eval) m() error {
	if len(e.loops) > 0 {
		m := e.loops[len(e.loops)-1].Max
		e.Stack.Push(float64(m))
		return nil
	}

	return fmt.Errorf("you cannot access 'm' outside a loop-body")
}

func (e *Eval) mod() error {
	return e.binOp(func(n float64, m float64) float64 {
		return float64(int(m) % int(n))
	})()
}

func (e *Eval) mul() error {
	return e.binOp(func(n float64, m float64) float64 { return m * n })()
}

func (e *Eval) nop() error {
	return nil
}

func (e *Eval) over() error {
	a, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	b, err := e.Stack.Pop()
	if err != nil {
		return err
	}

	e.Stack.Push(b)
	e.Stack.Push(a)
	e.Stack.Push(b)
	return nil
}
func (e *Eval) print() error {
	n, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	e.printNumber(n)
	e.printString("\n")
	return nil
}

func (e *Eval) setVar() error {
	offset, err := e.Stack.Pop()
	if err != nil {
		return err
	}
	value, err2 := e.Stack.Pop()
	if err2 != nil {
		return err2
	}
	e.vars[int(offset)].Value = value
	return nil
}

func (e *Eval) stackDump() error {
	l := e.Stack.Len()
	e.printString(fmt.Sprintf("<len:%d> ", l))

	c := 0
	for c < l {
		e.printNumber(e.Stack.At(c))
		e.printString(" ")
		c++
	}
	e.printString("\n")
	return nil

}

// startDefinition moves us into compiling-mode
//
// Note the interpreter handles removing this when it sees ";"
func (e *Eval) startDefinition() error {
	e.compiling = true
	return nil
}

// strings
func (e *Eval) stringCount() error {
	// Return the number of strings we've seen
	e.Stack.Push(float64(len(e.strings)))
	return nil
}

// strlen
func (e *Eval) strlen() error {
	addr, err := e.Stack.Pop()
	if err != nil {
		return err
	}

	i := int(addr)

	if i < len(e.strings) {
		str := e.strings[i]
		e.Stack.Push(float64(len(str)))
		return nil
	}

	return fmt.Errorf("invalid stack offset for string reference")

}

// strprn - string printing
func (e *Eval) strprn() error {
	addr, err := e.Stack.Pop()
	if err != nil {
		return err
	}

	i := int(addr)

	if i < len(e.strings) {
		str := e.strings[i]
		e.printString(str)
		return nil
	}

	return fmt.Errorf("invalid stack offset for string reference")
}

func (e *Eval) sub() error {
	return e.binOp(func(n float64, m float64) float64 { return m - n })()
}

func (e *Eval) swap() error {
	var a, b float64
	var err error
	b, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	a, err = e.Stack.Pop()
	if err != nil {
		return err
	}
	e.Stack.Push(b)
	e.Stack.Push(a)

	return nil
}

func (e *Eval) variable() error {
	e.defining = true
	return nil
}

func (e *Eval) words() error {
	known := []string{}

	for _, entry := range e.Dictionary {

		// Skip any word that contains a " " in its name,
		// this covers "$ $" which is a hack to execute
		// immediately-compiled words
		if !strings.Contains(entry.Name, " ") {
			known = append(known, entry.Name)
		}
	}

	sort.Strings(known)
	fmt.Printf("%s\n", strings.Join(known, " "))

	return nil
}

func (e *Eval) wordLen() error {
	known := []string{}

	for _, entry := range e.Dictionary {

		// Skip any word that contains a " " in its name,
		// this covers "$ $" which is a hack to execute
		// immediately-compiled words
		if !strings.Contains(entry.Name, " ") {
			known = append(known, entry.Name)
		}
	}

	e.Stack.Push(float64(len(known)))
	return nil
}
