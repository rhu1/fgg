# fgg -- README

---

This package includes:

* An **FG** static type checker and interpreter.
* An **FGG** static type checker and interpreter.
* An FGG static _nomono_ (i.e., "is/not monomorphisable") checker and FGG-to-FG **monomorphiser**.

Package organisation:
* `Makefile` -- for running tests and examples.
* `main.go` -- main file.
* `fg` -- FG AST, typing and evaluation.
  * `fg/examples` -- FG examples.
    * `fg/examples/oopsla20` -- The FG examples from the paper (i.e., Figs. 1,
      2 and 7).
* `fgg` -- FGG AST, typing, evaluation, nomono check, and monomorphisation.
  * `fgg/examples` -- FGG examples.
    * `fgg/examples/oopsla20` -- The FGG examples from the paper (i.e., Figs.
      3--6).
* `parser` -- FG/FGG parsers (generated using ANTLR).
  * `parser/FG.g4` -- FG ANTLR grammar.
  * `parser/FGG.g4` -- FGG ANTLR grammar.


---

### Install.

We assume a standard Go set up.  That is:

* Go (version 1.11.5+);
* a Go workspace, at `$GOPATH`;
* a `src` directory in the workspace.

Simply extract the contents of the zip (including the ANTLR library) directly into the `src` directory of
your Go workspace (i.e., `$GOPATH/src`).  

  * Then in the `github.com/rhu1/fgg` directory just extracted, the following command should work:

    `go run github.com/rhu1/fgg -eval=-1 -v fg/examples/oopsla20/fig1/functions.go`

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

* FGG does _not_ support any syntactic sugar -- this means empty type
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
/*Type and method decls -- semicolon separated*/
func main () { _ = /*main has this specific form*/ }
```

For testing purposes, the package supports this additional form:

```
package main;
import "fmt";  // Only fmt is allowed
/*Type and method decls -- semicolon separated*/
func main () { fmt.Printf("%#v", /* This specific Printf, and only in main*/) }
```

Notes:

  * This package additionally supports interface embedding for both FG and
    FGG.  E.g.,  
    `type A(type a Any()) interface { B(a) }  // Any, B are interfaces`

  * The `var` declarations used for readability in some of the examples in the paper are not supported.

---

### Example run commands.

Warning:  Type checking and _nomono_ errors raise a panic -- basic error
messages can be found at the top of the stack trace.

The following commands are run from the `github.com/rhu1/fgg` directory.

* **FG type check and evaluate**, with verbose printing.

  `go run github.com/rhu1/fgg -eval=-1 -v fg/examples/oopsla20/fig1/functions.go`

    * The argument to `-eval` is the number of steps to execute. `-1` means
      run to termination (either a value, or a panic due to a failed type
      assertion.)
    * `-eval` includes a dynamic type preservation check (an error raises a panic).

* **FGG type check and evaluate**, with verbose printing.  (Note the `-fgg`
  flag.)

  `go run github.com/rhu1/fgg -fgg -eval=-1 -v fgg/examples/oopsla20/fig3/functions.fgg`

* **FGG type check, nonomo check and monomorphisation**, with verbose printing.

  `go run github.com/rhu1/fgg -fgg -monomc=-- -v fgg/examples/oopsla20/fig3/functions.fgg`

    * The argument to `-monomc` is a file location for the FG output.  `--`
      means print the output.

* **Simulate FGG against its FG monomorphisation**, with verbose printing.

  `go run github.com/rhu1/fgg -test-monom -v
  fgg/examples/oopsla20/fig3/functions.fgg`

    * This includes dynamic checking of type preservation checking at both
      levels, and of the monomorphisation correspondence at every evaluation
      step.


---

### Example `Makefile` tests.

The following commands are run from the `github.com/rhu1/fgg` directory.

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
