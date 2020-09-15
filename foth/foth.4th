#
# This file is loaded on-startup, if it is present.
#
# NOTE: Lines having a "#"-prefix will be skipped.
#
#       This is not a standard approach to FORTH comments, but it makes
#       sense for this particular implementation.
#


#
# CR: Output a carrige return (newline).
#
: cr 10 emit ;

#
# Star: Output a star to the console.
#
# Here 42 is the ASCII code for the "*" character.
#
: star 42 emit ;


#
# Stars: Show the specified number of stars.
#
#        e.g. "3 stars"
#
: stars 0 do star loop 10 emit ;


#
# square: Square a number
#
: square dup * ;

#
# cube: cube a number
#
: cube dup square * ;

#
# 1+: add one to a number
#
: 1+ 1 + ;

#
# boot: output a message on-startup
#
: bootup 87 emit 101 emit 108 emit 99 emit 111 emit 109 emit 101 emit 32 emit 116 emit 111 emit 32 emit 102 emit 111 emit 116 emit 104 emit 33 emit 10 emit ;
bootup


# output "Hot"
: hot 72 emit 111 emit 116 emit 10 emit ;

# Output "Cod"
: cold 67 emit 111 emit 108 emit 100 emit 10 emit ;

: test_hot  0 > if hot then ;
: test_cold  0 <= if cold then ;
: temp dup test_hot test_cold ;
