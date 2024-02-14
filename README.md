[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/foth@v0.3.0/foth?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/foth)](https://goreportcard.com/report/github.com/skx/foth)
[![license](https://img.shields.io/github/license/skx/foth.svg)](https://github.com/skx/foth/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/skx/foth.svg)](https://github.com/skx/foth/releases/latest)

* [foth](#foth)
  * [Features](#features)
  * [Installation](#installation)
  * [Embedded Usage](#embedded-usage)
  * [Anti-Features](#anti-features)
  * [Implementation Approach](#implementation-approach)
  * [Implementation Overview](#implementation-overview)
    * [Part 1](#part-1) - Minimal initial-implementation.
    * [Part 2](#part-2) - Hard-coded recursive word definitions.
    * [Part 3](#part-3) - Allow defining minimal words via the REPL.
    * [Part 4](#part-4) - Allow defining improved words via the REPL.
    * [Part 5](#part-5) - Allow executing loops via `do`/`loop`.
    * [Part 6](#part-6) - Allow conditional execution via `if`/`then`.
    * [Part 7](#part-7) - Added minimal support for strings.
    * [Final Revision](#final-revision) - Idiomatic Go, test-cases, and many new words
  * [BUGS](#bugs)
    * [loops](#loops) - zero expected-iterations actually runs once
  * [See Also](#see-also)
  * [Github Setup](#github-setup)




# foth

A simple implementation of a FORTH-like language, hence _foth_ which is
close to _forth_.

If you're new to FORTH then [the wikipedia page](https://en.wikipedia.org/wiki/Forth_(programming_language)) is a good starting point, and there are more good reads online such as:

* [Forth in 7 easy steps](https://jeelabs.org/article/1612b/)
  * Just ignore any mention of the return-stack!
* [Starting FORTH](https://www.forth.com/starting-forth/)
  * A complete book, but the navigation of this site is non-obvious.

In brief FORTH is a stack-based language, which uses Reverse Polish notation.   The basic _thing_ in Forth is the "word", which is a named data item, subroutine, or operator.  Programming consists largely of defining new words, which are stored in a so-called "dictionary", in terms of existing ones.  Iteratively building up a local DSL suited to your particular task.

This repository was created by following the brief tutorial posted within the following Hacker News thread, designed to demonstrate how you could implement something _like_ FORTH, in a series of simple steps:

* https://news.ycombinator.com/item?id=13082825

The comment-thread shows example-code and pseudo-code in C, of course this repository is written in Go.



## Features

The end-result of this work is a simple scripting-language which you could easily embed within your golang application, allowing users to write simple FORTH-like scripts.  We implement the kind of features a FORTH-user would expect:

* Comments between `(` and `)` are ignored, as expected.
  * Single-line comments `\` to the end of the line are also supported.
* Support for floating-point numbers (anything that will fit inside a `float64`).
* Reverse-Polish mathematical operations.
  * Including support for `abs`, `min`, `max`, etc.
* Support for printing the top-most stack element (`.`, or `print`).
* Support for outputting ASCII characters (`emit`).
* Support for outputting strings (`." Hello, World "`).
* Support for basic stack operations (`clearstack`, `drop`, `dup`, `over`, `swap`, `.s`)
* Support for loops, via `do`/`loop`.
* Support for conditional-execution, via `if`, `else`, and `then`.
* Support for declaring variables with `variable`, and getting/setting their values with `@` and `!` respectively.
* Execute files specified on the command-line.
  * If no arguments are supplied run a simple REPL instead.
* A standard library is loaded, from the present directory, if it is present.
  * See what we load by default in [foth/foth.4th](foth/foth.4th).
* The use of recursive definitions, for example:
  * `: factorial recursive  dup 1 >  if  dup 1 -  factorial *  then  ;`



## Installation

You can find binary releases of the final-version upon the [project release page](https://github.com/skx/foth/releases), but if you prefer you can install from source easily.

Either run this to download and install the binary:

```
$ go get github.com/skx/foth/foth@v0.4.0

```

Or clone this repository, and build the executable like so:

```
cd foth
go build .
./foth
```

The executable will try to load [foth.4th](foth/foth.4th) from the current-directory, so you'll want to fetch that too.  But otherwise it should work as you'd expect - the startup-file defines several useful words, so running without it is a little annoying but it isn't impossible.



## Embedded Usage

Although this is a minimal interpreter it _can_ be embedded within a Golang host-application, allowing users to write scripts to control it.

As an example of this I put together a simple demo:

* [https://github.com/skx/turtle](https://github.com/skx/turtle)

This embeds the interpreter within an application, and defines some new words to allow the user to create graphics - in the style of [turtle](https://en.wikipedia.org/wiki/Turtle_graphics).



## Anti-Features

The obvious omission from this implementation is support for strings in the general case (string support is limited to outputting a constant-string).

We also lack the meta-programming facilities that FORTH users would expect, in a FORTH system it is possible to implement new control-flow systems, for example, by working with words and the control-flow directly.  Instead in this system these things are unavailable, and the implementation of IF/DO/LOOP/ELSE/THEN are handled in the golang-code in a way users cannot modify.

Basically we ignore the common FORTH-approach of using a return-stack, and implementing a VM with "cells".  Instead we just emulate the _behaviour_ of the more advanced words:

* So we implement `if` or `do`/`loop` in a hard-coded fashion.
  * That means we can't allow a user to define `while`, or similar.
  * But otherwise our language is flexible enough to allow _real_ work to be done with it.



## Implementation Approach

The code evolves through a series of simple steps, [contained in the comment-thread](https://news.ycombinator.com/item?id=13082825), ultimately ending with a [final revision](#final-revision) which is actually useful, usable, and pretty flexible.

While it would certainly be possible to further improve the implementation I'm going to declare this project as "almost complete" for my own tastes:

* I'll make minor changes, as they occur to me.
* Comments, test-cases, and similar are fair game.
* Outright crashes will be resolved, if I spot any.
* But no major new features will be added.

If **you** wanted to extend things further then there are some obvious things to work upon:

* Adding more of the "standard" FORTH-words.
  * For example we're missing `pow`, etc.
* Simplify the special-case handling of string-support.
* Simplify the conditional/loop handling.
  * Both of these probably involve using a proper return-stack.
  * This would have the side-effect of allowing new control-flow primitives to be added.
  * As well as more meta-programming.

Pull-requests adding additional functionality will be accepted with thanks.



## Implementation Overview

Each subdirectory within this repository gets a bit further down the comment-chain.

In terms of implementation two files are _largely_ unchanged in each example:

* `stack.go`, which contains a simple stack of `float64` numbers.
* `main.go`, contains a simple REPL/driver.
  * The final few examples will also allow loading a startup-file, if present.

Each example builds upon the previous ones, with a pair of implementation files that change:

* `builtins.go` contains the forth-words implemented in golang.
* `eval.go` is the workhorse which implements to FORTH-like interpreter.
  * This allows executing existing words, and defining new ones.


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

This part adds `do` and `loop`, allowing simple loops, and `emit` which outputs the ASCII character stored in the topmost stack-entry.

Sample usage would look like this:

```
> : cr 10 emit ;
> : star 42 emit ;
> : stars 0 do star loop cr ;
> 4 stars
****
> 5 stars
*****
> 1 stars
*
> 10 stars
**********
^D
```

Here we've defined two new words `cr` to print a return, and `star` to output a single star.

We then defined the `stars` word to use a loop to print the given number of stars.


(Note that the character `*` has the ASCII code 42).

`do` and `loop` are pretty basic, allowing only loops to be handled which increment by one each iteration.  You cannot use the standard `i` token to get the current index, instead you can see them on the stack:

* Top-most entry is the current index.
* Second entry is the limit.

So to write out numbers you could try something like this, using `dup` to duplicate the current offset within the loop:

     > : l 10 0 do dup . loop ;
     > l
     0.000000
     1.000000
     2.000000
     ..
     8.000000
     9.000000

     > : nums 10 0 do dup 48 + emit loop ;
     > nums
     0123456789>

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


### Part 7

This update adds a basic level of support for strings.

* When a string is encountered it is stored in "memory".
* The address of the string is pushed to the stack.
* Two new words are added:
  * `strlen` show the length of the string at the given address.
  * `strprn` print the string at the given address.

Sample usage:

    cd part7
    go build .
    ./part7
    > : steve "steve" ;
    > steve strlen .
    5
    > steve strprn .
    steve
    ^D

See [part7/](part7/) for the code.


### Final Revision

The final version, stored beneath [foth/](foth/), is pretty similar to the previous part from an end-user point of view, however there have been a lot of changes behind the scenes:

* We've added near 100% test-coverage.
* We've added a simple [lexer](foth/lexer/) to tokenize our input.
  * This was required to allow us to ignore comments, and handle string literals.
  * Merely splitting input-strings at whitespace characters would have made either of those impossible to handle correctly.
* The `if` handling has been updated to support an `else`-branch, the general form is now:
  * `$COND IF word1 [ .. wordN ] else alt_word1 [.. altN] then [more_word1 more_word2 ..]`
* It is now possible to use `if`, `else`, `then`, `do`, and `loop` outside word-definitions.
  * i.e. Immediately in the REPL.
* `do`/`loop` loops can be nested.
  * And the new words `i` and `m` used to return the current index and maximum index, respectively.
* There were many new words defined in the go-core:
  * `.s` to show the stack-contents.
  * `clearstack` to clear the stack.
  * `debug` to change the debug-flag.
  * `debug?` to reveal the status.
  * `dump` dumps the compiled form of the given word.
    * You can view the definitions of all available words this:
    * `#words 0 do i dump loop`
  * `#words` to return the number of defined words.
  * Variables can be declared, by name, with `variable`, and the value of the variable can be set/retrieved with `@` and `!` respectively.
    * See this demonstrated in the [standard library](foth/foth.4th)
* There were some new words defined in the [standard library](foth/foth.4th)
  * e.g. `abs`, `even?`, `negate`, `odd?`,
* Removed all calls to `os.Exit()`
  * We now return `error` objects where appropriate, allowing the caller to detect problems.
* It is now possible to redefine existing words.
* Execute any files specified on the command line.
  * If no files are specified run the REPL.
* We've added support for recursive definitions, in #16 for example allowing:
  * `: factorial recursive  dup 1 >  if  dup 1 -  factorial *  then  ;`

See [foth/](foth/) for the implementation.



## BUGS

A brief list of known-issues:


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

**NOTE**: In `gforth` the result of `0 0 do ... loop` is actually an __infinite__ loop, which is perhaps worse!

In our `stars` definition we handle this case by explicitly testing the loop
value before we proceed, only running the loop if the value is non-zero.




# See Also

This repository was put together after [experimenting with a scripting language](https://github.com/skx/monkey/), an [evaluation engine](https://github.com/skx/evalfilter/), putting together a [TCL-like scripting language](https://github.com/skx/critical), writing a [BASIC interpreter](https://github.com/skx/gobasic) and creating [yet another lisp](https://github.com/skx/yal).

I've also played around with a couple of compilers which might be interesting to refer to:

* Brainfuck compiler:
  * [https://github.com/skx/bfcc/](https://github.com/skx/bfcc/)
* A math-compiler:
  * [https://github.com/skx/math-compiler](https://github.com/skx/math-compiler)




# Github Setup

This repository is configured to run tests upon every commit, and when pull-requests are created/updated.  The testing is carried out via [.github/run-tests.sh](.github/run-tests.sh) which is used by the [github-action-tester](https://github.com/skx/github-action-tester) action.
