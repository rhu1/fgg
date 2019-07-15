# fgg
Mini prototype of FG/FGG/FGR in Go for quick testing.

So far, there is an almost-done FG -- need to add type assertions, test evaluation, and more general testing.

Parser is generated using ANTLR4.  CL incantation from repo root dir is: antlr4 -Dlanguage=Go -o parser FG -- where "antlr4" is, e.g., an alias for java -jar ~/code/java/lib/antlr-4.7.1-complete.jar (i.e., ANTLR4 installation):w
.  Grammar is here: https://github.com/rhu1/fgg/blob/master/FG.g4.

Exiting tests are here: https://github.com/rhu1/fgg/blob/master/fg/fg_test.go#L51.  Mainly syntax/typing tests.  Can run them from CL by: go test github.com/rhu1/fgg/fg.

"High-level" FGG (i.e., without monomorph) probably won't take long (a couple of days) on top of that -- overall code structure is the same, only relatively straightforward extensions of existing stuff.

Monomorph then on top of that.

Finally FGR, and "generic" FGG-FGR translation.
