# Part 4

Part four of the implementation is very similar to [part3](../part3/):

* We still have only a small number of built-in functions.
* We allow you to define new ones.

The difference here is that we now support compiling words which contain
_numbers_, in addition to references to existing words.

If you recall from [part3](../part3/):

* We have a notion of being in a "compiling state".
  * When `:` is encountered we flip that switch on.
  * When `;` is encountered we flip that switch off.
* When compiling we use a _temporary_ storage, and instead of executing words once we've found them we merely append their offsets to the list of words in this new definition.

The problem we found is that when the user entered a number in their new definition there was nothing we could do with it - adding it to the stack would be the wrong thing to do, and we couldn't add it to the word-definition.

We've updated things now so that the list of references to pre-existing words can store numbers too!  We add a magic "opcode", -1, and then add the number itself.

## Building

To build, and run this version:

```
go build .
./part4
> : square dup * ;
> 3 square .
9.000000
> : +1 1 + ;
> 3 +1 .
4.000000
^D
```


## Implementation

* We have a simple stack implemented in [stack.go](stack.go)
  * This allows storing a series of `float64` objects.
  * This is identical to the previous version.
* We've defined a `Word` structure to hold known-commands/words.
  * Each word has a name, a pointer to a (go) function, and also a series of word-indexes.
  * Since we store our words as an array we can access them by number, as well as name.
* Our interpreter, still has the `compiling` flag, when it juggles when sees `:` and `;`
* When not in compiling-mode things work as before.
  * If the token matches the name of one of the defined Words
    * And there is a (go) function-pointer, then call it.
    * Otherwise assume there is a list of `Word` indexes:
      * But if we see the magic flag `-1` then the next number is a number to push to the stack.
  * Otherwise assume the input is a number, and push to the stack.
* When in compiling mode we instead lookup each input word and just store the offset of the found `Word` in the internal space for the new Word.
  * If that fails then we assume the user entered a number, and add it to the word-array (prefixed by `-1`).
  * Once we hit `;` we add this new word to the dictionary of known-words
* [main.go](main.go) just reads STDIN, line by line, and passes to the evaluator

**NOTE**: As a consequence of storing numbers inside the definition of the word, prefixed by the `-1` marker, we've had to change the word-array to be an array of `float64` values - not integers.



## More Details

If we skip ahead to the [final version](../foth/) we gain access to a `dump` command, which lets us dissassemble the dictionary-entries.

Using that we can show the compiled version of the `+1` word we defined in this README:

```
   $ cd ../foth
   $ go build .
   $ ./foth
   Welcome to foth!
   > #words .
   71
   > : +1 1 + ;
   > #words .
   72
   > 71 dump
   Word '+1'
     0: store 1.000000
     2: +
```

Here we used `#words` to get a count of all the defined words (`words` would show their names).  To start with we see there are 71 words (0-70).  Then we define our new one, and confirm the count is updated.

The final step is to dump the word, and we see that the array containing the implementation contains three entries:

* 0x00: -1
* 0x01: 1.0
* 0x02: +

(The -1 is implicit here, and you have to guess that two bytes have been taken by the offset/index in the left-column.)

When this word is executed the following thing happens:

* The -1 value is taken to mean "load-word"
* The next value, 1, is then pushed onto the stack
* Then the "+" word is executed.
* And this word is complete.
