# Really Ez Lisp
a toy lisp compiler started for a very stupid reason

# WARNING!!
this project is being rewritten on another branch.
if you want to see the fancy features look at Really-Ez-Lisp2

## goals
- compile down to more native like go
- make web assembly and MySQL targets happen
- make error checker? debugging is really painful
- file io
- named function arguments and named references

## how to compile
`go run reallyezylisp.go -i <input file>`outputs to out.go  
`go run reallyezylisp.go -i <input file> -o <outputname>`outputs to <outputname>.go  
## how to run
`go run <whatever>.go`


## things you can do
- `+ - * /`
- unary / is same as 1/x
- unary -
- `(if (condition) (what to return if true) (what to return if not true))`
- `(use something)` copy paste contents from file named something. if its not go code its very likely it wont compile
- `(import something)` parse Really Ez Lisp from file called something.
- `(fn main (something))` this something is all that will actually be done
- `(fn yourFuncHere (do something) (actually return something))` func names have to be valid go names, sorry ¯\_(ツ)_/¯. also only last statement is returned
- `(fn example (+ $0 1)))` $n gives the nth argument supplied to the function. (terrible!!!)
- `(list 1 2 3 4)` make list [1 2 3 4]
- `(get 0 (list 1 2))` gets element 0
- `(get 1 2 (list 0 1 2 3))` gets elements 1-2
- `(get 1 2 (list 0 1 2 3)(list 0 1 2 3))` returns list of sub-lists you asked for. I'm very proud of this one!!!
- `(len (list 1 2 3))` this should be obvious, but this language does have this feature
- `(print 123)` prints 123
- `(printc 127829)` print whatever character. in this case 🍕
- `(printc '\n')` print newline character
- `(printc "Hello, World!\n")` no one has ever done this before
- `(printc (list 72 101 108 108 111 44 32 87 111 114 108 100 33 10))` same effect as above, but no one has wanted to do this before
- `(head (list 1 2 3))` returns the first element (1)
- `(head (list 1 2 3)(list 2 3 4)(list 3 4 5))` returns all the first elements (1 2 3 4)
- `(tail (list 1 2 3))` returns everything but head
- `(tail (list 1 1 2 3) (list 1 1 2 3) (list 1 1 2 3))` tail all the things! ((1 2 3)(1 2 3)(1 2 3))
- `(and 0 1)` or `(&& 0 1 0)` can be just 2 or many
- `(or 1 0 1)` or `(|| 0 1)` same as last  but different

## currently busted
- probably a lot of things ¯\_(ツ)_/¯ but i cant think of any atm

