# Part 2

Part two of the implementation is very similar to [part1](../part1/) in
the sense that it only allows the execution of hard-coded words.

However the implementation has become more flexible, because we've updated
our `Word`-structure to allow words to be executed in two ways:

* A word can contain a reference to a (go) function.
  * If present this is called, as before.
* If there is no (go) function reference we do something different.
  * We call existing words, based on their directory index

## Building

To build, and run this version:

```
go build .
./part2
> 2 dup * print
4.000000
> 4 square print
16.000000
^D
```


## Implementation

* We have a simple stack implemented in [stack.go](stack.go)
  * This allows storing a series of `float64` objects.
  * This is identical to the previous version.
* We've defined a `Word` structure to hold known-commands/words
  * Each word has a name, a pointer to a (go) function, and also a series of word-indexes.
  * Since we store our words as an array we can access them by number, as well as name.
* Our interpreter, [eval.go](eval.go) has a list of hard-coded Words defined.
  * Each token is executed in the same way:
    * If the token matches the name of one of the defined Words
      * And there is a (go) function-pointer, then call it.
      * Otherwise assume there is a list of `Word` indexes.  Call each one in turn.
    * Otherwise assume the input is a number, and push to the stack.
* [main.go](main.go) just reads STDIN, line by line, and passes to the evaluator
