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
\        This is only supported here because we process single-line
\        comments before the other kind of comments.


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
\ CR: Output a carrige return (newline).
\
: cr 10 emit ;

\
\ Star: Output a star to the console.
\
\ Here 42 is the ASCII code for the "*" character.
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
\ Squares: Draw a box
\
\          e.g. 10 squares
\
: squares 0 do
   over stars cr
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
\ Here we repeat our work because we don't have support for ELSE when we
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
: frozen 102 emit 114 emit 111 emit 122 emit 101 emit 110 emit 10 emit ;

\ Output "NOT frozen\n"
: non_frozen 78 emit 79 emit 84 emit 32 emit 102 emit 114 emit 111 emit 122 emit  101 emit 110 emit 10 emit ;

\ Output one or other of the messages?
: frozen? 0 <= if frozen else non_frozen then cr ;

\ All in one.
: frozen2? 0 <= if ." frozen " else ." not frozen " then cr ;


\
\ We have to ensure we have a newline on the last line, or it will be
\ ignored
\
