# Part 1

Part one of the implementation only deals with hard-coded execution
of "words".  It only supports the basic mathematical operations, along
with the ability to print the top-most entry of the stack:

## Building

To build, and run this version:

```
go build .
./part1
> 2 3 + 4 5 + * print
45.000000
^D
```


## Implementation

* We have a simple stack implemented in [stack.go](stack.go)
  * This allows storing a series of `float64` objects.
* We've defined a `Word` structure to hold known-commands/words
  * Each word has a name, and a pointer to the function which contains the implementation.
* Our interpreter, [eval.go](eval.go) has a list of hard-coded Words defined.
  * Each token is executed in the same way:
    * If the token matches the name of one of the defined Words, then call the function associated with it.
    * Otherwise assume the input is a number, and push to the stack.
* [main.go](main.go) just reads STDIN, line by line, and passes to the evaluator
