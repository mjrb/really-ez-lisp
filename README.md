# Really Ez Lisp2
Childhood dreams are made possible by gophers and aliens.

# WARNING!!!
this is part of a rewrite and is pretty broken right now.
don't expect this to compile, but I promise the cool features that work, work.

## Goals
- generate more native code for Golang target
  + use native Go types
  + use Go built ins over ridiculous wrapped functions
- dead code elimination
- analysis phase so that debugging wont cause loss of sanity
- duck typing
- MySQL target
- let statement and named parameters

## Progress
### generate more native code for Golang target
- use native Go types
yes! although there an abusive amount of type assertions and interface{} :(  
the alternative may be to check types at compile time and use unsafe at run time :(  
- use Go built ins over ridiculous wrapped functions
the groundwork is there, currently only print and println are implemented.
all functions other than main are forced to return interface{} so if a function
ends in a built in call, it get wrapped in a closure and looks weird. may be fixed  
with duck typing
### dead code elimination
yes! uncalled functions never go into source, so hello world is no longer 300 line output.### analysis phase so that debugging wont cause loss of sanity
yes! gives compile time errors with a line number, instead of passive aggressive warnings
### duck typing
not quite there yet, may be able to build backwards off of last call in functions
### MySQL target
no progress. list will be implemented a table, but i may have to do garbage collection or
make persistent immutable lists in SQL. no grantees if that will turn out good
### let statement and named parameters
yes! now it is use-able and doesn't feel like "what if bash was lisp"


## future goals
- LLVM
  + LLVM ir target is possible
  + actually using LLVM to generate bin is hard. i cant get the go extensions to compile. they may have rotted
- Hygienic macros
  + I probably cant call it a lisp if it doesn't have these
- Web Assembly target
  + may be solved with LLVM. however the structure of Web Assembly is similar to how
    programs are represented internally.
- Closures
  + i can't believe it still doesn't have these
- add library and go extension support back