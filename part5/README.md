# Part 5

Part five of the implementation is very similar to [part4](../part4/), the difference is that we've added _very basic_ support for looping instructions.


## Building

To build, and run this version:

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

Here we've defined two new words `cr` to print a return, and `star` to output a single star.  We then use those in our `stars` word to allow showing the user the specified number of stars, followed by a newline.


## Implementation

The implementation here is as simple as the previous version, we look for `do` and when we find it we remember where it occurred in the input.  We then output another "special opcode" when we hit the `loop` word:

* We use `-2` to mark a loop-test.
  * At runtime we pop two items from the stack, and if they're not equal we jump back to the opcode of the `do` instruction.
  * This supports our loop.

This looping implementation is __very basic__.  It stores the `do`/`loop` state on the stack, as you can see here:

```
> : bar 5 1 do dup . loop ;
> bar
1.000000
2.000000
3.000000
4.000000
```
