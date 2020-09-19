# Part 3

Part three of the implementation is very similar to [part2](../part2/) in
the sense that we allow execution words which are defined in terms of
pre-existing words.

Our `Word`-structure looks the same as in the previous part:

* A word can contain a reference to a (go) function.
  * If present this is called, as before.
* If there is no (go) function reference we do something different.
  * We call existing words, based on their directory index

The difference is in our evaluator itself:

* Now we have a notion of being in a "compiling state".
  * When `:` is encountered we flip that switch on.
  * When `;` is encountered we flip that switch off.
* When compiling we use a _temporary_ word, and instead of executing words once we've found them we merely append their offsets to the list of words in this new definition.
  * There's a bit of magic to setup the name of the new word too

## Building

To build, and run this version:

```
go build .
./part3
> : square dup * ;
> 3 square .
9.000000
^D
```


## Implementation

* We have a simple stack implemented in [stack.go](stack.go)
  * This allows storing a series of `float64` objects.
  * This is identical to the previous version.
* We've defined a `Word` structure to hold known-commands/words.
  * Each word has a name, a pointer to a (go) function, and also a series of word-indexes.
  * Since we store our words as an array we can access them by number, as well as name.
* Our interpreter, [eval.go](eval.go) has been updated to set/clear the `compiling` flag, when it sees `:` and `;`
* When not in compiling-mode things work as before
  * If the token matches the name of one of the defined Words
    * And there is a (go) function-pointer, then call it.
    * Otherwise assume there is a list of `Word` indexes.  Call each one in turn.
    * Otherwise assume the input is a number, and push to the stack.
* When in compiling mode we instead lookup each input word and just store the offset of the found `Word` in the internal space for the new Word
  * Once we hit `;` we add this new word to the dictionary of known-words
* [main.go](main.go) just reads STDIN, line by line, and passes to the evaluator

There is one omission here - if we're in compiling mode we cannot handle numbers.  Recall that when we're executing code we add any number to the stack, as we encounter it.

At the point we're compiling words we don't have any way of using the stack though - that only comes into play when we're _executing_ the word.

This problem will be resolved in [part4](../part4/).
