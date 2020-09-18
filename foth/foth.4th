\
\ This file is loaded on-startup, if it is present.
\
\ In FORTH there are two kinds of comments:
\
\  1. Anything between \ and a newline will be skipped.
\
\  2. Anything between brackets `(` + `)` will be skipped.
\     These comments can span lines, although they typically do not.
\
\  NOTE: Nesting `(` and `)` comments is a syntax error, so the following
\        is explicitly denied:
\
\        ( comment ( comment again ) )
\

\
\ Declare a variable named `PI`
\
variable PI

\
\ Set the value of PI to be the expected constant.
\
\ The following web-reference is useful reading for variable-access, even
\ though our support is slightly different:
\
\  https://www.forth.com/starting-forth/8-variables-constants-arrays/
\
3.14 PI !

\
\ Variables can be retrieved, and displayed, like so:
\
\    PI @ .


\
\ Of course you can modify them inside words.
\
\ The following example sets a variable "val", and allows showing/incrementing
\ it
\
variable val

\
\ A simple helper to get the content of a variable and show it.
\
: ?  @ . ;

\
\ Get the value of the variable named "val", and show it
\
: show_val val ? ;

\
\ Get the value, increase by one, and store it back
\
: inc_val val @ 1 + val ! ;

\
\ Get the value, decrement, and store.
\
: dec_val val @ 1 - val ! ;


\
\ bell: make some noise
\
: bell 7 emit ;

\
\ gnome-terminal won't re-ring the bell, unless there is a delay.
\
: bells
    3 0 do
        bell
        3000000 0 do
            nop
        loop
    loop
;

\
\ We define `=` (and `==`) by default, but we do not have a built-in
\ function for not-equals.  We can fix that now:
\
: != = invert ;

\
\ Similarly we do not have a built-in method for abs(N), so we
\ can resolve that here
\
: abs dup 0 > if 1 else -1 then * ;


\
\ CR: Output a carrige return (newline).
\
: cr 10 emit ;

\
\ Star: Output a star to the console.
\
\ Here 42 is the ASCII code for the "*" character.
\
\ If you prefer you could use:
\
\   : star '*' emit ;
\
: star 42 emit ;


\
\ Stars: Show the specified number of stars.
\
\        e.g. "3 stars"
\
\ We add a test here to make sure that the user enters > 0
\ as their argument
\
: stars dup 0 > if 0 do star loop else drop then ;

\
\ Squares: Draw a square of stars
\
\          e.g. 10 squares
\
\ Here we define a loop that runs from N to 0, and we use "m" as
\ the maximum-value of the loop.
\
\ Inside loop-bodies we can access two variables like that:
\
\    i  -> The current value of the loop
\
\    m  -> The maximum value of the loop (which will terminate it).
\
: squares 0 do
   m stars cr
  loop ;


\
\ square: Square a number
\
: square dup * ;

\
\ cube: cube a number
\
: cube dup square * ;


\
\ 1+: add one to a number
\
: 1+ 1 + ;


\
\ boot: output a message on-startup
\
: bootup ." Welcome to foth!\n " ;
bootup

\
\ IF test
\
\ This section of the startup-file outputs either "hot" or "cold" depending
\ on whether a number is <0 or not.
\
\ Here we repeat our work because we did't have support for ELSE when we
\ added this example.
\

\ output a hot/cold message
: hot  ." Hot\n  " ;
: cold ." Cold\n " ;

: test_hot   0 >  if hot then ;
: test_cold  0 <= if cold then ;
: temp? dup test_hot test_cold ;


\
\ Here we have a word that uses IF and ELSE, allowing our test to
\ be simplified compared to the previous version.
\

\ Output "frozen\n"
: frozen ." frozen\n " ;

\ Output "NOT frozen\n"
: non_frozen ." NOT frozen\n " ;

\ Output the appropiate message.
: frozen? 0 <= if frozen else non_frozen then cr ;

\
\ Or we could have written the following word to do everything
\ in one-step:
\
: frozen2? 0 <= if ." frozen " else ." NOT frozen " then cr ;


\
\ We have to ensure we have a newline on the last line, or it will be
\ ignored
\
