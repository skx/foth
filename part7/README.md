# Part 7

Part seven of the implementation is very similar to [part6](../part6/), the difference is that we've added **basic** support for strings.

The specific problem we have is that our stack, and word definitions, only allow support for storing floating-point numbers.  So we cannot store a string on the stack, which means we must be indirect:

* When we see a string we store it in an array of strings.
* We then push the offset of the new string entry onto the stack.
* This allows it to be referenced and used.

However because we don't have arbitrary read/write to RAM opcodes/words we can't do much more than that.  We've added two new string-specific words as a proof of concept though:

* `strlen` - Return the length of a string.
* `strprn` - Print a string.
* `strings` - Return the maximum string index we've seen.



## Building

To build and run this version:

```
go build .
./part7
> : steve "steve" ;
> steve strlen .
5
> steve strprn .
steve
^D
```



## Implementation

The implementation here is pretty simple again, as suits a tutorial-code.

The interpreter got a new string-storing area:

```
// Eval is our evaluation structure
type Eval struct {

..
	// strings contains string-storage
	strings []string
}
```

In the past when we saw input we didn't recognize we assumed it was a number, and parsed that with `strconv.ParseFloat`, but now we test if the first character of the token is a `"` character.  If it is we :

* Strip the leading/trailing `"` from it.
* Append the string to our storage array.
* Push the offset of the new entry.
  * Using the magic -1 value, if we're in compiling mode.

From there there is no special support.  The primitives just read from the string area, for example `strlen` looks like this:

```
// strlen
func (e *Eval) strlen() {
	addr := e.Stack.Pop()
	i := int(addr)

	if i < len(e.strings) {
		str := e.strings[i]
		e.Stack.Push(float64(len(str)))
	} else {
		e.Stack.Push(-1)
	}
}
```



## Bugs

Because we don't have a decent lexer we can only handle strings without spaces, or newlines.
