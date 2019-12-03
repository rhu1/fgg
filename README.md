# fgg -- README

Mini prototype of FG/FGG/FGR written in Go for quick testing.  
The primary aim is to keep the implementation "very close" to the formalisms
to assist understanding and experimentation.

FG and FGG are running.  Monomorphisation is WIP.  FGR is TODO.  
I mainly hacked it together the last couple of weekends, so apologies for
bugs, etc.

Go version 1.11 or later is needed.  Download this package to `github.com/rhu1/fgg` under 
your `$GOPATH/src` directory.

The [FG](https://github.com/rhu1/fgg/blob/master/parser/FG.g4) and
[FGG](https://github.com/rhu1/fgg/blob/master/parser/FGG.g4) grammars are
written using ANTLR 4.  

To bypass using ANTLR to generate the parsers, copy all contents from
[`parser/pregen/fg`](https://github.com/rhu1/fgg/tree/master/parser/pregen/fg)
to `parser/fg`, for FG; similarly for FGG.  
Otherwise, ANTLR can be installed following (e.g.) the instructions
[here](https://blog.gopheracademy.com/advent-2017/parsing-with-antlr4-and-go/).  
Then, use `go generate`, or from the repo root dir: `antlr4 -Dlanguage=Go -o parser/fg parser/FG.g4`,
assuming `antlr4` is an alias for, e.g.,
`java -jar [snip]/antlr-4.7-complete.jar`;
similarly for FGG.

The ANTLR4 Runtime for Go is needed to run the interpreters.  
Example download instructions are
[here](https://github.com/antlr/antlr4/blob/master/doc/go-target.md) (mainly
the second bullet).

Warning: the FG grammar is white-space insensitive while remaining a subset of
Go, and thus requires separators like ';' to be written regardless of
newlines.  
Similarly for FGG.

FG example usages: (`-eval=-1` means evaluate until a value is reached)

- `go run github.com/rhu1/fgg -eval=-1 -v -inline="package main; type A struct {}; func main() { _ = A{} }"`
- `go run github.com/rhu1/fgg -eval=10 -v fg/examples/hello/hello.go`
- See [`fg/examples`](https://github.com/rhu1/fgg/tree/master/fg/examples) for
more examples, including all the examples from the paper submission.
- To run the existing FG tests: `go test github.com/rhu1/fgg/fg`

FGG example usages:

- `go run github.com/rhu1/fgg -fgg -eval=-1 -v -inline="package main; type A(type ) struct {}; func main() { _ = A(){} }"`
- `go run github.com/rhu1/fgg -fgg -eval=10 -v fgg/examples/hello/hello.fgg`
- See [`fgg/examples`](https://github.com/rhu1/fgg/tree/master/fgg/examples) for
more examples, including all the examples from the paper submission.
- To run the existing FGG tests: `go test github.com/rhu1/fgg/fgg`

I started a naive first hack at a monomorphisation function -- it is still
WIP, the reported results are not yet complete and have been in no way checked
for correctness.  
It can be experimented with using the `-monom` flag (use it in conjunction
with `-fgg` and `-v`).  
E.g., `go run github.com/rhu1/fgg -fgg -monom -v fgg/examples/popl20/monomorph/monomorph.fgg`  
The above prints out the monomorphised code using something close to the
notation used in the paper.  
The output can also be formatted to be valid Go using the `-monomc` flag  
E.g., `go run github.com/rhu1/fgg -fgg -monomc=tmp/out.go -v fgg/examples/popl20/monomorph/monomorph.fgg`  
Use `-monomc=--` to print to stdout instead of writing to a file.


