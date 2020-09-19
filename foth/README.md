# foth

This is our final version of the code.  Compared to the previous revision it has been overhauled quite significantly, but it should still be obvious there is a common ancestry present:

* We've moved a bit of the code around, into sub-packages.
  * We've added test-cases to __everything__.
* We implement a lot of new primitives.
  * In [eval/builtins.go](eval/builtins.go).
  * And also in a new file loaded at startup: [foth.4th](foth.4th).
* We've reworked our implementation to be a bit more idiomatic go.
* We've improved a lot of implementation details:
  * `do`/`loop` now work in the REPL, not just when compiled into words.
    * `do`/`loop` can now be nested too.
  * `if`/`else`/`then` work now.
    * If you recall in [part6](../part6/) we didn't support an else-branch.


## Building

## Building

To build, and run this version:

```
go build .
./foth
```


## Implementation

No real surprises here, but we did make a bunch of changes to ease usage.

We also moved to using a real `go.mod`, so you could embed the interpreter in your application if you wished:

* [main.go](main.go) now does this.
 * It even defines a "secret word" to show how you could export your application-specific methods to the scripting-language.

For example the [stack/stack.go](stack/stack.go) gained a bunch of new methods, so that we can dump the stack at run-time we need to see the size, and peek at given offsets:

```
$ ./foth
Welcome to foth!
> 1 2 3
> .s
<3> 1.000000 2.000000 3.000000
> .
3
> .
2
> .s
<1> 1.000000
```

Here we see:

* We added the numbers `1`, `2`, and `3` to the stack.
* Then we used the new `.s` word to dump the stack-contents.
 *  That showed "<3>" (meaning three entries) then the values.
* We removed two numbers, by printing them.
* Then we dumped again.

Similarly we decided to add a simple [lexer](lexer/) to allow handling string-immediates, and skipping comments without bothering our interpreter with them.

The lexer also allowed us to simplify some things.  We used "output star" as a simple demo a lot of times, defined like this:

    : star 42 emit ;

With the lexer we can sneakily allow a different form:

    : star '*' emit ;

The lexer just turns `'X'` into the ASCII code for the character X.  Simple, but useful.

We added error-testing, at parse/run-time.  And allow the state to be dumped.

* As noted already we can see the list of words via `words`.
* We can see the count of words via `#words .`.
* And we can dump a word with `32 dump`.
  * Dumping all words is as simple as:
     * `#words 0 do i dump loop`
* We show the debug-state via `debug?` and allow it to be toggled via:
  * `debug? invert debug`
* We added variable-support, both defining then getting/setting.



## Where next?

Pull-requests for new primitives, cleanups, etc are most welcome.

But I'm probably done .. although it is quite addictive working with DSL like this!

Steve
--
