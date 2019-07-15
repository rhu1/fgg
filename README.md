# fgg
Mini prototype of FG/FGG/FGR in Go for quick testing.

So far, there is an almost-done FG -- need to add type assertions, test evaluation, and more general testing.

Exiting tests are in fg/fg_test.go.  Mainly syntax/typing tests.  Can run them from CL by: go test github.com/rhu1/fgg/fg.

"High-level" FGG (i.e., without monomorph) probably won't take long (a couple of days) on top of that -- overall code structure is the same, only relatively straightforward extensions of existing stuff.

Monomorph then on top of that.

Finally FGR, and "generic" FGG-FGR translation.
