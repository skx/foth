\
\ This file is loaded on-startup, if it is present.
\
\ NOTE: Lines having a "\"-prefix will be skipped.
\
\       This is not a standard approach to FORTH comments, but it makes
\       sense for this particular implementation.
\


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
