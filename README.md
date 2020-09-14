* [foth](#foth)
  * [Implementation Overview](#implementation-overview)
    * [Part 1](#part-1)
    * [Part 2](#part-2)


# foth

A simple implementation of a forth-like (hence foth) language, as
described in the following Hacker News comment:

* https://news.ycombinator.com/item?id=13082825


## Implementation Overview

Each subdirectory gets a bit further down the comment-chain.

There is a common "stack.go" which contains a simple `float64` stack,
there is no support for strings or similar.

Built-ins are documented in `builtins.go` in each directory.


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


### Part 3

Part three allows the user to define their own words, right from within the
REPL!

This means we've removed the `square` implementation, because you can add your own:

     cd part3
     go build .
     ./part3
     > : square dup * ;
     > 3 square
     9.000000
     ^D

**NOTE**: We don't support using numbers in definitions, yet.  That will come in part4!
