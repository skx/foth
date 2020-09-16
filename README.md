* [foth](#foth)
  * [Implementation Overview](#implementation-overview)
    * [Part 1](#part-1) - Minimal initial-implementation.
    * [Part 2](#part-2) - Hard-coded recursive word definitions.
    * [Part 3](#part-3) - Allow defining minimal words via the REPL.
    * [Part 4](#part-4) - Allow defining improved words via the REPL.
    * [Part 5](#part-5) - Allow executing loops via `do`/`loop`.
    * [Part 6](#part-6) - Allow conditional execution via `if`/`then`.
    * [Final Revision](#final-revision) - Idiomatic Go, test-cases, and many new words
  * [BUGS](#bugs)
    * [loops](#loops) - zero expected-iterations actually runs once
  * [Github Setup](#github-setup)


# foth

A simple implementation of a FORTH-like language, hence _foth_ which is
close to _forth_.

If you're new to FORTH then the following wikipedia page is a good starting point:

* [Forth](https://en.wikipedia.org/wiki/Forth_(programming_language))

This repository was created by following the brief tutorial posted within the following Hacker News thread:

* https://news.ycombinator.com/item?id=13082825

The end-result of this work is a simple scripting-language which you could easily embed within your golang application, allowing users to write simple FORTH-like scripts.  We have the kind of features you would expect from a minimal system:

* Reverse-Polish mathematical operations.
* Comments between `(` and `)` are ignored, as expected.
  * Single-line comments `\` to the end of the line are also supported.
* Support for floating-point numbers (anything that will fit inside a `float64`).
* Support for printing the top-most stack element (`.`, or `print`).
* Support for outputting ASCII characters (`emit`).
* Support for outputting strings (`." Hello, World "`).
* Support for basic stack operatoins (`dup`, `swap`, `drop`)
* Support for loops, via `do`/`loop`.
* Support for conditional-execution, via `if`, `else`, and `then`.
* Load any files specified on the command-line
  * If no arguments are included run the REPL
* A standard library is loaded, from the present directory, if it is present.
  * See what we load by default in [foth/foth.4th](foth/foth.4th).

The code evolves through a series of simple steps, guided by the comment-linked, ultimately ending with a featurefull [final revision](#final-revision) which is actually usable.

While it would be possible to further improve the implementation from the final stage I'm going to stop there, because I think I've done enough for the moment.

If you did want to extend further then there _are_ some obvious things to add:

* Adding more of the "standard" FORTH-words.
* Case-insensitive lookup of words.
  * e.g. "dup" should act the same way as "DUP".
* Pull-requests to add additional functionality to the [final revision](#final-revision) are most welcome.
* Simplify the special-cases of string-support, along with the conditional/loop handling.


## Implementation Overview

Each subdirectory gets a bit further down the comment-chain.

In terms of implementation two files are _largely_ unchanged in each example:

* `stack.go`, which contains a simple stack of `float64` numbers.
* `main.go`, contains a simple REPL/driver.
  * The final few examples will also allow loading a startup-file, if present.

Each example builds upon the previous ones, with a pair of implementation files that change:

* `builtins.go` contains the words implemented in golang.
* `eval.go` is the workhorse which implements to FORTH-like interpreter.


### Part 1

Part one of the implementation only deals with hard-coded execution
of "words".  It only supports the basic mathematical operations, along
with the ability to print the top-most entry of the stack:

     cd part1
     go build .
     ./part1
     > 2 3 + 4 5 + * print
     45.000000
     ^D

See [part1/](part1/) for details.



### Part 2

Part two allows the definition of new words in terms of existing ones,
which can even happen recursively.

We've added `dup` to pop an item off the stack, and push it back twice, which
has the ultimate effect of duplicating it.

To demonstrate the self-definition there is the new function `square` which
squares the number at the top of the stack.

     cd part2
     go build .
     ./part2
     > 3 square .
     9.000000
     > 3 dup + .
     6.000000
     ^D

See [part2/](part2/) for details.



### Part 3

Part three allows the user to define their own words, right from within the
REPL!

This means we've removed the `square` implementation, because you can add your own:

     cd part3
     go build .
     ./part3
     > : square dup * ;
     > : cube dup square * ;
     > 3 cube .
     27.000000
     > 25 square .
     625.000000
     ^D

See [part3/](part3/) for details.

**NOTE**: We don't support using numbers in definitions, yet.  That will come in part4!


### Part 4

Part four allows the user to define their own words, including the use of numbers, from within the REPL.  Here the magic is handling the input of numbers when in "compiling mode".

To support this we switched our `Words` array from `int` to `float64`, specifically to ensure that we could continue to support floating-point numbers.

     cd part4
     go build .
     ./part4
     > : add1 1 + ;
     > -100 add1 .
     -99.000000
     > 4 add1 .
     5.000000
     ^D

See [part4/](part4/) for details.



### Part 5

This part adds `do`, `emit`, and `loop`, allowing simple loops.

(Emit outputs the ASCII character stored in the topmost stack-entry)

Sample usage would look like this - note that the character `*` has the ASCII code 42:

    cd part5
    go build .
    ./part5
    > : star 42 emit ;
    > : stars 0 do star loop 10 emit ;
    > 10 stars
    **********
    > 3 stars
    ***
    ^D

Here `do` is a NOP, and the `loop` instruction handles a pair of values
on the stack.

See [part5/](part5/) for details.



### Part 6

This update adds a lot of new primitives to our dictionary of predefined words:

* `drop` - Removes an item from the stack.
* `swap` - Swaps the top-most two stack-items.
* `words` - Outputs a list of all defined words.
* `<`, `<=`, `=` (`==` as a synonym), `>`, `>=`
  * Remove two items from the stack, and compare them appropriately.
  * If the condition is true push `1` onto the stack, otherwise `0`.
* The biggest feature here is the support for using `if` & `then`, which allow conditional actions to be carried out.
  * (These are why we added the comparison operations.)

In addition to these new primitives the driver, `main.go`, was updated to load and evaluate [foth.4th](part6/foth.4th) on-startup if it is present.

Sample usage:

    cd part6
    go build .
    ./part6
    > : hot 72 emit 111 emit 116 emit 10 emit ;
    > : cold 67 emit 111 emit 108 emit 100 emit 10 emit ;
    > : test_hot  0 > if hot then ;
    > : test_cold  0 <= if cold then ;
    > : test dup test_hot test_cold ;
    > 10 test
    Hot
    > 0 test
    Cold
    > -1 test
    Cold
    > 10 test_hot
    Hot
    > 10 test_cold
    > -1 test_cold
    Cold
    ^D

See [part6/](part6/) for the code.

**NOTE**: The `if` handler allows:

    : foo $COND IF word1 [word2 .. wordN] then [more_word1 more_word2 ..] ;

This means if the condition is true then we run `word1`, `word2` .. and otherwise we skip them, and continue running after the `then` statement.  Specifically note there is **no support for `else`**.  That is why we call the `test_host` and `test_cold` words in our `test` definition.  Each word tests separately.

As an example:

    > : foo 0 > if star star then star star cr ;

If the test-passes, because you give a positive number, you'll see FOUR stars.  if it fails you just get TWO:

     > 2 foo
     ****
     > 1 foo
     ****
     > 0 foo
     **
     > -1 foo
     **

This is because the code is synonymous with the following C-code:

     if ( x > 0 ) {
        printf("*");
        printf("*");
     }
     printf("*");
     printf("*");
     printf("\n");

I found this page useful, it also documents `invert` which I added for completeness:

* https://www.forth.com/starting-forth/4-conditional-if-then-statements/




### Final Revision

The final version, stored beneath [foth/](foth/), is pretty similar to the previous part, however there have been a lot of changes behind the scenes:

* We've added a large number of test cases, with a good amount of test-coverage.
* We use a simple [lexer/](lexer/) to tokenize our input
  * This was required to allow us to ignore comments, and handle string literals.
  * Merely splitting on whitespace characters would have left either of those impossible.
* The `if` handling has been updated to support an `else`-branch, the general form is now:
  * `$COND IF word1 [ .. wordN ] else alt_word1 [.. altN] then [more_word1 more_word2 ..]`
* It is now possible to use `if`, `else`, `then`, `do`, and `loop` outside word-definitions.
  * i.e. Immediately.
* There were many new words defined:
  * `debug` to change the debug-flag.
  * `debug?` to reveal the status.
  * `dump` dump the compiled form of the given word
    * You can view all the definitions with something like this:
    * `#words 0 do dup dump loop`
  * `#words` to return the number of defined words.
* Removed all calls to `os.Exit()`
  * We return `error` objects where appropriate, allowing the caller to detect problems.
* Make redefining existing words possible.
  * Note that due to our implementation previously defined words remain unchanged, even if a word is replaced/updated.
* Load any files specified on the command line.
  * If no files are specified run the REPL.

See [foth/](foth/) for the implementation.



## BUGS

There are two known-issues at the moment:

### Loops

The handling of loops isn't correct when there should be zero-iterations:

```
     > : star 42 emit ;
     > : stars 0 do star loop 10 emit ;
     > 3 stars
     ***
     > 1 stars
     *
     > 0 stars
     *
     ^D
```

In our `stars` definition we handle this by explicitly testing the loop
value before we proceed - At the moment any loop of `0 0` will run once
so you'll need to add that test if we can't fix this for the general case.


# Github Setup

This repository is configured to run tests upon every commit, and when pull-requests are created/updated.  The testing is carried out via [.github/run-tests.sh](.github/run-tests.sh) which is used by the [github-action-tester](https://github.com/skx/github-action-tester) action.
