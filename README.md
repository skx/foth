* [foth](#foth)
  * [Implementation Overview](#implementation-overview)
    * [Part 1](#part-1)
    * [Part 2](#part-2)
    * [Part 3](#part-3)
    * [Part 4](#part-4)
    * [Part 5](#part-5)
  * [TODO](#todo)
  * [BUGS](#bugs)


# foth

A simple implementation of a FORTH-like language, hence _foth_ which is
close to _forth_.

This repository was implemented based upon the following Hacker News comment:

* https://news.ycombinator.com/item?id=13082825

The feature-set is pretty minimal:

* Reverse-Polish mathematical operations.
* Support for floating-point numbers (anything that will fit in a float64).
* Support for printing the top-most stack element (`.`, or `print`).
* Support for outputting ASCII characters (`emit`).
* Support for loops.


## Implementation Overview

Each subdirectory gets a bit further down the comment-chain.

In terms of implementation two files are identical in each example:

* `stack.go` contains a simple `float64` stack.
* `main.go` contains a simple driver.

Each example has a slightly improving set of built-in functions implemented
in golang, which you can see in `builtins.go`.


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

Part two allows defining new words in terms of others, internally we now
allow recursive use of previously-defined words, as well as the built-in
functions.

We've added `dup` to pop an item off the stack, and push it back twice - essentially duplicating it.

To demonstrate the self-definition there is the new function `square` which squares the
top number on the stack.

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

This part adds `do`, `emit`, and `loop`, allowing simple loops:

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



## TODO

- Control-Flow (i.e. "if").
- Test-cases.
- Declare and document a "final" version:
  - Allow executing a file-contents, not just a REPL.
  - Have a few more primitives.



## BUGS

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

Should test in the `do` maybe?  Before the first iteration?
