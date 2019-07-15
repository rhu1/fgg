# fgg

20190715

Added type assertions, and some CL flag options (e.g., -eval=n,
-inline="...").

A first HelloWorld run can be done by:
go run github.com/rhu1/fgg -eval=10 fg/examples/hello/hello.go

20190714

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

Exiting tests are here:
https://github.com/rhu1/fgg/blob/master/fg/fg_test.go#L51.  
Mainly syntax/typing tests.  
Can run them from CL by: go test github.com/rhu1/fgg/fg.

"High-level" FGG (i.e., without monomorph) probably won't take long (a couple
of days) on top of that -- overall code structure is the same, only relatively
straightforward extensions of existing stuff.

Monomorph then on top of that.

Finally FGR, and "generic" FGG-FGR translation.
