# README for [fgg](https://github.com/rhu1/fgg)

---

This `fgg` package is a minimal and basic prototype of **Featherweight Go** and
**Featherweight Generic Go**, as presented in:

> Featherweight Go  
> *Robert Griesemer, Raymond Hu, Wen Kokke, Julien Lange, Ian Lance Taylor,  
> Bernardo Toninho, Philip Wadler and Nobuko Yoshida*  
> https://arxiv.org/abs/2005.11710

Currently, many aspects of the code are quite primitive, mainly for the
convenience of quick experimentation alongside the above paper.  For example,
types/functions/variables are not well named, except for some correspondence
with the formal definitions.  The tool is also not particularly user-friendly:

- it offers only the small (but meaningful) subset of Go as formalised in the
  paper;
- it does not support *any* syntactic sugar -- e.g., empty parentheses and type
  lists, and various separators (`;`), all need to be written out explicitly;
- most type errors are reported as panics, though an error message may be given
  at the top of the stack trace.

We plan to improve some of this in the near future.  Contact [Raymond
Hu](https://go.herts.ac.uk/raymond_hu) for issues related to this repository.

See this [Go blog post](https://blog.golang.org/generics-next-step) for
information about the [generics design
draft](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md)
by the Go team, and links to their generic-to-ordinary Go translation tool
(including an online playground) based on that draft.

- Currently, the `fgg` tool supports a few features (within its fragment of Go)
  that their translation tool does not.  These include type parameters for
  methods (for now, the latter has type parameters for types and functions
  only), _nomono_ (monomorphisability) checking, and covariant method receiver
  bounds.

[Featherweight-go-gen](https://github.com/wenkokke/featherweight-go-gen) is
tool that enumerates FGG programs and integrates with `fgg` for
testing.


---

### Summary.

This package includes:

* An **FG** static type checker and interpreter.
* An **FGG** static type checker and interpreter.
* An FGG static _nomono_ (i.e., "is/not monomorphisable") checker and FGG-to-FG
  **monomorphiser**.

Package organisation:
* `Makefile` -- install, and for running tests and examples.
* `main.go` -- main file.
* `fg` -- FG AST, typing and evaluation.
* `fgg` -- FGG AST, typing, evaluation, nomono check, and monomorphisation.
* `examples`
  * `examples/fg` -- FG examples.
    * `examples/fg/oopsla20` -- The FG examples from the paper (i.e., Figs. 1,
      2 and 7).
  * `examples/fgg` -- FGG examples.
    * `examples/fgg/oopsla20` -- The FGG examples from the paper (i.e., Figs.
      3--6).
* `parser` -- FG/FGG parsers (generated using ANTLR).
  * `parser/FG.g4` -- FG ANTLR grammar.
  * `parser/FGG.g4` -- FGG ANTLR grammar.

Use the [main](https://github.com/rhu1/fgg/tree/main) branch for the
latest working version.

**Syntax.**  The best source would be the formal grammars in the paper, or else
see the above ANTLR grammars.


---

### Install.

We assume a standard Go set up.  That is:

* Go (version 1.11+);
* a Go workspace, at `$GOPATH`;
* a `src` directory in the workspace.

You will also need the ANTLR v4 runtime for Go; e.g., see "Installing ANTLR v4"
in this
[tutorial](https://blog.gopheracademy.com/advent-2017/parsing-with-antlr4-and-go/).

Clone the `fgg` repo into the `src` directory of your Go workspace, i.e.,
`$GOPATH/src` (or use `go get`).  It should end up located at
`src/github.com/rhu1/fgg`.

Then, either copy over the pre-generated parser files and install by

- `make install-pregen-parser`  
  (generated using ANTLR 4.7.1)

or generate the parsers yourself using ANTLR and install by

- (assuming some suitable `antlr4` command; e.g., `java -jar [antlr-4.7.1-complete.jar]`)  
`antlr4 -Dlanguage=Go -o parser/fg parser/FG.g4`  
  `antlr4 -Dlanguage=Go -o parser/fgg parser/FGG.g4`  
  `make install`

To test the install -- inside the `github.com/rhu1/fgg` directory, this command
should work:

- `go run github.com/rhu1/fgg -eval=-1 -v examples/fg/oopsla20/fig1/functions.go`

Afer installing, you can also use the resulting `fgg` binary directly instead
of `go run`.

This package has been tested using Go version 1.11.5+ on:

* MacOS Catalina
* Cygwin/Windows 10

---

### A note on syntax.

Two points:

* FG and FGG _always_ need explicit `;` separators between type/method
  decls, field decls, etc., even across new lines.  E.g., in FG

  ```
  package main;
  type A struct {};
  type B struct {
    f1 A;
    f2 A
  };
  type IA interface {
    m1() B;
    m2() B
  };
  func (x0 B) foo() B { return x0 };
  func main() { _ = B{A{}, A{}}.foo() }
  ```

  A rule of thumb is, write all FG/FGG code as if you were writing Go _without_ line breaks.

* FGG does not support any syntactic sugar -- this means empty type
  declarations and type argument lists must always be written out in full.
  E.g., the FGG equivalent to the above is:

  ```
  package main;
  type A(type ) struct {};
  type B(type ) struct {
    f1 A();
    f2 A()
  };
  type IA(type ) interface {
    m1(type )() B();
    m2(type )() B()
  };
  func (x0 B(type )) foo(type )() B() { return x0 };
  func main() { _ = B(){A(){}, A(){}}.foo()() }
  ```


---

### FG/FGG program syntax.

This package implements the grammars defined in the paper.  E.g., a basic FGG
program has the form:

```
package main;
/* Type and method decls -- semicolon separated */
func main () { _ = /* main has this specific form */ }
```

For testing purposes, the package supports this additional form:

```
package main;
import "fmt";  // Only fmt is allowed
/* Type and method decls -- semicolon separated */
func main () { fmt.Printf("%#v", /* This specific Printf, and only in main */ ) }
```

Notes:

  * This package additionally supports interface embedding for both FG and
    FGG.  E.g.,  
    `type A(type a Any()) interface { B(a) }  // Any, B are interfaces`

  * The `var` declarations used for readability in some of the examples in the paper are not supported.


---

### Example run commands.

*Warning*:  Type checking and _nomono_ errors raise a panic -- basic error
messages can be found at the top of the stack trace.

The following commands can be run from the `github.com/rhu1/fgg` directory.

* **FG type check and evaluate**, with verbose printing.

  `go run github.com/rhu1/fgg -eval=-1 -v examples/fg/oopsla20/fig1/functions.go`

    * The argument to `-eval` is the number of steps to execute. `-1` means
      run to termination (either a value, or a panic due to a failed type
      assertion.)
    * `-eval` includes a dynamic type preservation check (an error raises a panic).

* **FGG type check and evaluate**, with verbose printing.  (Note the `-fgg`
  flag.)

  `go run github.com/rhu1/fgg -fgg -eval=-1 -v examples/fgg/oopsla20/fig4/functions.fgg`

* **FGG type check, nonomo check and monomorphisation**, with verbose printing.

  `go run github.com/rhu1/fgg -fgg -monomc=-- -v examples/fgg/oopsla20/fig4/functions.fgg`

    * The argument to `-monomc` is a file location for the FG output.  `--`
      means print the output.

* **Simulate FGG against its FG monomorphisation**, with verbose printing.

  `go run github.com/rhu1/fgg -test-monom -v examples/fgg/oopsla20/fig4/functions.fgg`

    * This includes dynamic checking of type preservation checking at both
      levels, and of the monomorphisation correspondence at every evaluation
      step.


---

### Example `Makefile` tests.

The following commands assume `make install`, and that the resulting `fgg`
binary is on the `$PATH`.

Running from the `github.com/rhu1/fgg` directory:

* `make test-monom-against-go`

  * Type checks and evaluates a series of FGG programs;
  * _nomono_ checks and monomorphises each to an FG program;
  * type checks and evaluates the FG program using `fgg`;
  * compiles and executes the FG program using `go`;
  * compares the results from `fgg` and `go`.

* `make simulate-monom`

  * Simulates a series of FGG programs against their FG monomorphisations.

* `make test-all`

  * Run all tests.
