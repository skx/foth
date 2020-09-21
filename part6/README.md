# Part 6

Part six of the implementation is very similar to [part5](../part5/), the difference is that we've added support for `if`/`then` instructions.

**NOTE** We didn't add support for `else`.

In FORTH `if` looks like this:

```
: one? 1 = if 42 emit 10 emit then ;
```

`if` will execute if the top-most word on the stack is `1`, otherwise it will skip to the word after the `then`.

(Of course if there were `else` support then that would be jumped to instead!)



## Building

To build and run this version:

```
go build .
./part6
> : one? 1 = if 42 emit 10 emit then ;
> 1 one?
*
> 0 one?
> 3 one?
> 1 one?
*
> ^D
```

Here we've defined a word which takes a number from the stack, compares it with `1` and if equal outputs a star (and newline).

You'll see that we've implemented a whole bunch of new primitives, specifically to allow new conditional things:

* `=` is an equality test.  It pops two values off the stack.
  * If equal it pushes `1`.
  * Otherwise it pushes `0`.
* `<` is a less-than test.  It pops two values off the stack.
  * If the test passes it pushes `1`.
  * Otherwise it pushes `0`.
* etc.


## Implementation

The implementation here is pretty simple again, as suits a tutorial-code.

When we see an `if` we need to add a conditional-jump opcode, which will skip over the loop body at run-time if the topmost stack value is not 0.

In C we might write something like this:

    if ( a < b ) {
       do stuff
    }

When we implement that we have to have this in our definition:

    a
    b
    <
    if
    -3           \ -3 is a JUMP opcode, which conditionally
    XXX          \ jumps to the specified offset.
     do
     stuff
  XXX:

Here we're using `-3` as the magic opcode to mean "pop a value from the stack, and if the value is zero we jump to the specified offset in our word-list.

Hopefully that is clear!
