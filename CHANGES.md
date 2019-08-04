# fgg -- CHANGELOG

TODO: proper formatting

20190804 Ray

Nick added `cmd/fg2fgg`.

I added a (temporary) Makefile for testing, including -monom and fg2fgg.


20190724 Ray

Refactored common Adaptor ("parser"), Name, Decl, Program and Expr interfaces
between FG and FGG.  
This reduces code duplication in main and tests, but makes the type hierarchy
slightly more abstract/complicated and introduces a few casts.

Added some more FGG examples in fgg/examples/popl20, e.g., irregular and
monomorph.

Started a naive first hack at a monomorphisation function -- is WIP.  
It can be experimented with using the -monom flag, e.g.,:  
go run github.com/rhu1/fgg -fgg -monom -v fgg/examples/popl20/monomorph/monomorph.fgg  
The results that have not been checked for correctness yet.  
The -monom flag is ignored unless -fgg is also set (and won't show anything if
-v is not set).  
Use the -compile-out.go flag instead of -monom to write the output as valid Go
to a file.


20190723 Ray

Added initial version of FGG, with examples from the paper.

To generate parser: antlr4 -Dlanguage=Go -o parser/fgg parser/FGG.g4  
(Or use go generate -- or simply copy all contents of parser/pregen/fgg into
parser/fgg.)  
Note: I moved the .g4 files inside the parser directory.

Warning: I did not implement any sugar (and don't intend to) -- all FGG must
be written strictly according to the formal grammar.  
(Unlike all the examples in the paper.)

First run (e.g.): go run github.com/rhu1/fgg -v -fgg -inline="package main;
type A(type ) struct {}; func main() { _ = A(){} }"  
Note: need the -fgg flag.  (-fg is still the implicit default.)

To run an existing example (e.g.): go run github.com/rhu1/fgg -fgg -eval=-1 -v
fgg/examples/popl20/compose/compose.fgg

To run tests: go test github.com/rhu1/fgg/fgg

Examples are here: https://github.com/rhu1/fgg/tree/master/fgg/examples

Have started a naive monomorphisation routine -- WIP.


20190716 Ray

Added FG examples from the paper:
https://github.com/rhu1/fgg/tree/master/fg/examples.  
All working -- except for "whoopsie", as the check for "bad recursive struct
decls" is not done yet.  
Example commands to run them are in comments at the top of each.

Note: FG (or at least this implementation) requires mandatory ";" between all
decls (e.g., types and methods, and also fields).  
This is to make it such that both (i) FG is white-space insensitive and (ii) every
FG program is an actual Go program.

Added -v (verbose printing, default=false), and -eval=-1 (evaluate until
value, or panic) CL flags.

Extension to FGG on top of the current codebase would probably take a couple
of evenings, though I may not be able to do so immediately.


20190715 Ray

Added type assertions, and some CL flag options (e.g., -eval=n,
-inline="...").  
N.B. flags must be given before any non-flag args.

A first HelloWorld run can be done by:
go run github.com/rhu1/fgg -eval=10 -v fg/examples/hello/hello.go


20190714 Ray

Mini prototype of FG/FGG/FGR in Go for quick testing.

So far, there is an almost-done FG -- need to add type assertions, test
evaluation, and more general testing.  
(Knocked it up during a bit of free time in the week end.)

Parser is generated using ANTLR4.  
CL incantation from repo root dir is: antlr4 -Dlanguage=Go -o parser FG --
where "antlr4" is, e.g., an alias for: java -jar
~/code/java/lib/antlr-4.7.1-complete.jar (i.e., an ANTLR4 installation).  
Grammar is here: https://github.com/rhu1/fgg/blob/master/FG.g4#L40.  
The ANTLR4 Runtime for Go is also needed:
https://github.com/antlr/antlr4/blob/master/doc/go-target.md.

Existing tests are here:
https://github.com/rhu1/fgg/blob/master/fg/fg_test.go#L51.  
Mainly syntax/typing tests.  
Can run them from CL by: go test github.com/rhu1/fgg/fg.

"High-level" FGG (i.e., without monomorph) probably won't take long (a couple
of days) on top of that -- overall code structure is the same, only relatively
straightforward extensions of existing stuff.

Monomorph then on top of that.

Finally FGR, and "generic" FGG-FGR translation.
