# fgg

20190716 Ray

Added FG examples from the paper:
https://github.com/rhu1/fgg/tree/master/fg/examples.  All working -- except
for "whoopsie", as the check for "bad recursive struct decls" is not done yet.
Example commands to run them are in comments at the top of each.

Note: FG (or at least this implementation) requires mandatory ";" between all
decls (e.g., types and methods, and also fields), to make the language
white-space insensitive, unlike actual Go.

Extension to FGG on top of the current codebase would probably take a couple
of evenings, though I may not be able to do so immediately.

20190715 Ray

Added type assertions, and some CL flag options (e.g., -eval=n,
-inline="...").  N.B. flags must be given before any non-flag args.

A first HelloWorld run can be done by:
go run github.com/rhu1/fgg -eval=10 fg/examples/hello/hello.go

20190714 Ray

Mini prototype of FG/FGG/FGR in Go for quick testing.

So far, there is an almost-done FG -- need to add type assertions, test
evaluation, and more general testing.  (Knocked it up during a bit of free
time in the week end.)

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
