//$ go test github.com/rhu1/fgg/fg
//$ go test github.com/rhu1/fgg/fg -run Test001

package fg_test // Separate package, can test "API"

import (
	"fmt"
	"testing"

	"github.com/rhu1/fgg/internal/base"
	"github.com/rhu1/fgg/internal/base/testutils"
	"github.com/rhu1/fgg/internal/fg"
)

/* Harness funcs */

func fgParseAndOkGood(t *testing.T, elems ...string) base.Program {
	var adptr fg.FGAdaptor
	return testutils.ParseAndOkGood(t, &adptr, fg.MakeFgProgram(elems...))
}

// N.B. do not use to check for bad *syntax* -- see the PARSER_PANIC_PREFIX panic check in base.ParseAndOkBad
func fgParseAndOkBad(t *testing.T, msg string, elems ...string) base.Program {
	var adptr fg.FGAdaptor
	return testutils.ParseAndOkBad(t, msg, &adptr, fg.MakeFgProgram(elems...))
}

/* Syntax and typing */

// TOOD: make translation to FGG and compare results to -fgg

func Test001(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	Am2 := "func (x0 A) m2(x1 A) A { return x1 }"
	Am3 := "func (x0 A) m3(x1 A, x2 A) A { return x2 }"
	B := "type B struct { f A }"
	e := "B{A{}}"
	fgParseAndOkGood(t, A, Am1, Am2, Am3, B, e)
}

func Test002(t *testing.T) {
	//parseAndOkGood(t, "A{}") // Testing parseAndOkGood
	fgParseAndOkGood(t, "type A struct {}", "A{}")
}

func Test002b(t *testing.T) {
	//parseAndOkBad(t, "type A not declared", "type A struct{}", "A{}")  // Testing parseAndOkBad
	fgParseAndOkBad(t, "type A not declared", "A{}")
}

func Test002c(t *testing.T) {
	fgParseAndOkBad(t, "A doesn't take anything", "type A struct {}", "A{A{}}")
}

func Test003(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	e := "B{A{}}"
	fgParseAndOkGood(t, A, B, e)
}

func Test003b(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	e := "B{B{A{}}}"
	fgParseAndOkBad(t, "B takes an A, not a B", A, B, e)
}

func Test003c(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	e := "B{}"
	fgParseAndOkBad(t, "B takes an A", A, B, e)
}

func Test004(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	C := "type C struct { f1 A; f2 B }"
	e := "C{A{}, B{A{}}}"
	fgParseAndOkGood(t, A, B, C, e)
}

func Test005(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	e := "A{}"
	fgParseAndOkGood(t, A, Am1, e)
}

func Test005b(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return A{} }"
	e := "A{}"
	fgParseAndOkGood(t, A, Am1, e)
}

func Test005c(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	Am2 := "func (x0 A) m2(x1 A) A { return x1 }"
	e := "A{}"
	fgParseAndOkGood(t, A, Am1, Am2, e)
}

/*func Test005d(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	Am2 := "func (x0 A) m2(x1 A) A { return x0.m1() }"  // TODO
	e := "A{}"
	parseAndOkGood(t, A, Am1, Am2, e)
}*/

func Test005e(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	Am2 := "func (x0 A) m2(x1 A) A { return x1 }"
	Am3 := "func (x0 A) m3(x1 A, x2 B) B { return x2 }"
	B := "type B struct { f A }"
	e := "A{}"
	fgParseAndOkGood(t, A, Am1, Am2, Am3, B, e)
}

func Test006(t *testing.T) {
	Any := "type Any interface {}"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() Any { return x0 }"
	e := "A{}"
	fgParseAndOkGood(t, Any, A, Am1, e)
}

func Test007(t *testing.T) {
	IA := "type IA interface { m1(x1 A) A }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	fgParseAndOkGood(t, IA, A, Am1, B, e)
}

func Test007b(t *testing.T) {
	IA := "type IA interface { m1(x1 A) A }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m2(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	fgParseAndOkBad(t, "A is not an IA", IA, A, Am1, B, e)
}

func Test007c(t *testing.T) {
	IA := "type IA interface { m1(x1 A) A }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	fgParseAndOkBad(t, "A is not an IA", IA, A, Am1, B, e)
}

func Test007d(t *testing.T) {
	Any := "type Any interface {}"
	IA := "type IA interface { m1(x1 A) A }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) Any { return x0 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	fgParseAndOkBad(t, "A is not an IA", Any, IA, A, Am1, B, e)
}

func Test008(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return foo }"
	e := "A{}"
	fgParseAndOkBad(t, "foo is not bound", A, Am1, e)
}

func Test009(t *testing.T) {
	Any := "type Any interface { }"
	IA := "type IA interface { m1(x1 A) A; Any }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	fgParseAndOkGood(t, Any, IA, A, Am1, B, e)
}

func Test009b(t *testing.T) {
	Any := "type Foo interface { foo(x A) A }"
	IA := "type IA interface { m1(x1 A) A; Foo }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) A { return x1 }"
	Afoo := "func (x0 A) foo(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	fgParseAndOkGood(t, Any, IA, A, Am1, Afoo, B, e)
}

func Test010b(t *testing.T) {
	Any := "type Foo interface { foo(x A) A }"
	IA := "type IA interface { m1(x1 A) A; Foo }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	fgParseAndOkBad(t, "A is not an IA", Any, IA, A, Am1, B, e)
}

// Testing bad return
func Test011(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return B{A{}} }"
	B := "type B struct { f A }"
	e := "B{A{}}"
	fgParseAndOkBad(t, "Cannot return a B as an A", A, Am1, B, e)
}

// Initial testing for select
func Test012(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	e := "B{A{}}.f"
	fgParseAndOkGood(t, A, B, e)
}

func Test012b(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	e := "B{A{}}.f1"
	fgParseAndOkBad(t, "B does not have a \"f1\" field", A, B, e)
}

func Test012c(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	e := "B{B{A{}}.f}"
	fgParseAndOkGood(t, A, B, e)
}

// Initial testing for call
func Test013(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1() A { return x0 }"
	e := "A{}.m1()"
	fgParseAndOkGood(t, A, A1m, e)
}

func Test013b(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1() A { return x0.m1() }"
	e := "A{}.m1()"
	fgParseAndOkGood(t, A, A1m, e)
}

func Test013c(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 A) A { return x1 }"
	e := "A{}.m1(A{})"
	fgParseAndOkGood(t, A, A1m, e)
}

func Test013d(t *testing.T) {
	fmt.Println("Source:")
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 A) A { return x1.m1(x0) }"
	e := "A{}.m1(A{}.m1(A{}))"
	fgParseAndOkGood(t, A, A1m, e)
}

func Test013e(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 A) A { return x0 }"
	e := "A{}.m1(A{}.m1())"
	fgParseAndOkBad(t, "(Nested) m1 call missing arg", A, A1m, e)
}

func Test013f(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 A) A { return x0 }"
	e := "A{}.m1(A{}.m1(A{}, A{}))"
	fgParseAndOkBad(t, "(Nested) m1 call too many args", A, A1m, e)
}

func Test013g(t *testing.T) {
	fmt.Println("Source:")
	A := "type A struct {}"
	B := "type B struct { f A }"
	A1m := "func (x0 A) m1(x1 A) A { return x0 }"
	e := "A{}.m1(A{}.m1(B{A{}}))"
	fgParseAndOkBad(t, "(Nested) m1 call given a B, expecting an A", A, A1m, B, e)
}

// Fixed bug in methods, md.t => md.recv.t
func Test014(t *testing.T) {
	fmt.Println("Source:")
	A := "type A struct {}"
	e := "A{}.m1()"
	fgParseAndOkBad(t, "A has no method m1", A, e)
}

func Test015(t *testing.T) {
	Any := "type Any interface {}"
	A := "type A struct {}"
	B := "type B struct { f A }"
	A1m := "func (x0 A) m1(x1 Any) Any { return B{x0} }"
	e := "A{}.m1(B{A{}})"
	fgParseAndOkGood(t, Any, A, A1m, B, e)
}

func Test015b(t *testing.T) {
	fmt.Println("Source:")
	IA := "type IA interface { m0() A }"
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 IA) A { return x0 }"
	e := "A{}.m1(A{})"
	fgParseAndOkBad(t, "A is a not an IA", IA, A, A1m, e)
}

// Initial testing for assert
func Test016(t *testing.T) {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	A := "type A struct {}"
	e := "ToAny{A{}}.any.(A)"
	fgParseAndOkGood(t, Any, ToAny, A, e)
}

func Test016b(t *testing.T) {
	A := "type A struct {}"
	e := "A{}.(A)"
	fgParseAndOkBad(t, "Stupid cast on A struct lit", A, e)
}

// FIXME: should be a parser panic (lexing error, bad token), but currently caught as a typing panic
func Test017(t *testing.T) {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	e := "ToAny{1}" // ANTLR "warning token recognition error at: '1'" -- need to escalate to strict
	fgParseAndOkBad(t, "Bad token, \"1\"", Any, ToAny, e)
}

// Testing OK check for multiple declarations of a type/method name
func Test018(t *testing.T) {
	A := "type A struct {}"
	e := "A{}"
	fgParseAndOkBad(t, "Multiple declarations of type name 'A'", A, A, e)
}

func Test018b(t *testing.T) {
	A := "type A struct {}"
	Am := "func (x0 A) m() A { return x0 }"
	e := "A{}"
	fgParseAndOkBad(t, "Multiple declarations of method name 'm' for receiver A",
		A, Am, Am, e)
}

/* Eval */

// TODO: run all the above Good tests using -eval=-1
// TODO: put these tests through actual Go and compare the results
// TOOD: and make translation to FGG and compare results to -fgg

func TestEval001(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	e := "B{A{}}.f"
	prog := fgParseAndOkGood(t, A, B, e)
	testutils.EvalAndOkGood(t, prog, 1)
}

func TestEval002(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0.m1() }"
	e := "A{}.m1()"
	prog := fgParseAndOkGood(t, A, Am1, e)
	testutils.EvalAndOkGood(t, prog, 10)
}

func TestEval003(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() B { return B{x0} }"
	B := "type B struct { f A }"
	e := "A{}.m1().f"
	prog := fgParseAndOkGood(t, A, Am1, B, e)
	testutils.EvalAndOkGood(t, prog, 2)
}

// Initial testing for assert -- Cf. Test016
func TestEval004(t *testing.T) {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	A := "type A struct {}"
	e := "ToAny{A{}}.any.(A)"
	prog := fgParseAndOkGood(t, Any, ToAny, A, e)
	testutils.EvalAndOkGood(t, prog, 2)
}

// Testing isValue on StructLit
func TestEval005(t *testing.T) {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	A := "type A struct {}"
	e := "ToAny{ToAny{ToAny{A{}}.any.(A)}}"
	prog := fgParseAndOkGood(t, Any, ToAny, A, e)
	testutils.EvalAndOkGood(t, prog, 2)
}

// //TODO: test -eval=-1 -- test is currently added as -eval=0
func TestEval006(t *testing.T) {
	A := "type A struct {}"
	e := "A{}"
	prog := fgParseAndOkGood(t, A, e)
	testutils.EvalAndOkGood(t, prog, 0)
}

/* fmt.Sprintf */

func TestEval007(t *testing.T) {
	imp := "import \"fmt\""
	A := "type A struct {}"
	e := "fmt.Sprintf(\"\")"
	prog := fgParseAndOkGood(t, imp, A, e)
	testutils.EvalAndOkGood(t, prog, 1)
}

func TestEval008(t *testing.T) {
	imp := "import \"fmt\""
	A := "type A struct {}"
	e := "fmt.Sprintf(\"%v ,_()+- %v\", A{}, A{})"
	prog := fgParseAndOkGood(t, imp, A, e)
	testutils.EvalAndOkGood(t, prog, 1)
}
