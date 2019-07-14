//rhu@HZHL4 MINGW64 ~/code/go/src/temp/antlr/antlr04
//$ go test temp/antlr/antlr04/fg
//$ go test temp/antlr/antlr04/fg -run Test001

package fg_test // Separate package, can test "API"

import (
	"fmt"
	"strings"
	"testing"

	"temp/antlr/antlr04/fg"
)

/* Harness funcs */

func parseAndCheckOk(prog string) {
	var adptr fg.FGAdaptor
	ast := adptr.Parse(true, prog)
	ast.Ok()
}

func parseAndOkGood(t *testing.T, elems ...string) {
	prog := fg.MakeFgProgram(elems...)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: " + fmt.Sprintf("%v", r) + "\n" + prog)
		}
	}()
	parseAndCheckOk(prog)
}

// N.B. do not use to check for bad *syntax* -- see below, "[Parser]" panic check
func parseAndOkBad(t *testing.T, msg string, elems ...string) {
	prog := fg.MakeFgProgram(elems...)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but none occurred: " + msg + "\n" +
				prog)
		} else {
			rec := fmt.Sprintf("%v", r)
			if strings.HasPrefix(rec, "[Parser]") {
				t.Errorf("Unexpected panic: " + rec + "\n" + prog)
			}
			// TODO FIXME: check panic more specifically
		}
	}()
	parseAndCheckOk(prog)
}

/* Syntax and typing */

func Test001(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	Am2 := "func (x0 A) m2(x1 A) A { return x1 }"
	Am3 := "func (x0 A) m3(x1 A, x2 A) A { return x2 }"
	B := "type B struct { a A }"
	e := "B{A{}}"
	parseAndOkGood(t, A, Am1, Am2, Am3, B, e)
}

func Test002(t *testing.T) {
	//parseAndOkGood(t, "A{}") // Testing parseAndOkGood
	parseAndOkGood(t, "type A struct {}", "A{}")
}

func Test002b(t *testing.T) {
	//parseAndOkBad(t, "type A not declared", "type A struct{}", "A{}")  // Testing parseAndOkBad
	parseAndOkBad(t, "type A not declared", "A{}")
}

func Test002c(t *testing.T) {
	parseAndOkBad(t, "A doesn't take anything", "type A struct {}", "A{A{}}")
}

func Test003(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { a A }"
	e := "B{A{}}"
	parseAndOkGood(t, A, B, e)
}

func Test003b(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { a A }"
	e := "B{B{A{}}}"
	parseAndOkBad(t, "B takes an A, not a B", A, B, e)
}

func Test003c(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { a A }"
	e := "B{}"
	parseAndOkBad(t, "B takes an A", A, B, e)
}

func Test004(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { a A }"
	C := "type C struct { a A; b B }"
	e := "C{A{}, B{A{}}}"
	parseAndOkGood(t, A, B, C, e)
}

func Test005(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	e := "A{}"
	parseAndOkGood(t, A, Am1, e)
}

func Test005b(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return A{} }"
	e := "A{}"
	parseAndOkGood(t, A, Am1, e)
}

func Test005c(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	Am2 := "func (x0 A) m2(x1 A) A { return x1 }"
	e := "A{}"
	parseAndOkGood(t, A, Am1, Am2, e)
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
	B := "type B struct { a A }"
	e := "A{}"
	parseAndOkGood(t, A, Am1, Am2, Am3, B, e)
}

func Test006(t *testing.T) {
	Any := "type Any interface {}"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() Any { return x0 }"
	e := "A{}"
	parseAndOkGood(t, Any, A, Am1, e)
}

func Test007(t *testing.T) {
	IA := "type IA interface { m1(x1 A) A }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	parseAndOkGood(t, IA, A, Am1, B, e)
}

func Test007b(t *testing.T) {
	IA := "type IA interface { m1(x1 A) A }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m2(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	parseAndOkBad(t, "A is not an IA", IA, A, Am1, B, e)
}

func Test007c(t *testing.T) {
	IA := "type IA interface { m1(x1 A) A }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	parseAndOkBad(t, "A is not an IA", IA, A, Am1, B, e)
}

func Test007d(t *testing.T) {
	Any := "type Any interface {}"
	IA := "type IA interface { m1(x1 A) A }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) Any { return x0 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	parseAndOkBad(t, "A is not an IA", Any, IA, A, Am1, B, e)
}

func Test008(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return foo }"
	e := "A{}"
	parseAndOkBad(t, "foo is not bound", A, Am1, e)
}

func Test009(t *testing.T) {
	Any := "type Any interface { }"
	IA := "type IA interface { m1(x1 A) A; Any }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	parseAndOkGood(t, Any, IA, A, Am1, B, e)
}

func Test009b(t *testing.T) {
	Any := "type Foo interface { foo(a A) A }"
	IA := "type IA interface { m1(x1 A) A; Foo }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) A { return x1 }"
	Afoo := "func (x0 A) foo(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	parseAndOkGood(t, Any, IA, A, Am1, Afoo, B, e)
}

func Test010b(t *testing.T) {
	Any := "type Foo interface { foo(a A) A }"
	IA := "type IA interface { m1(x1 A) A; Foo }"
	A := "type A struct {}"
	Am1 := "func (x0 A) m1(x1 A) A { return x1 }"
	B := "type B struct { f IA }"
	e := "B{A{}}"
	parseAndOkBad(t, "A is not an IA", Any, IA, A, Am1, B, e)
}

// Testing bad return
func Test011(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return B{A{}} }"
	B := "type B struct { a A }"
	e := "B{A{}}"
	parseAndOkBad(t, "Cannot return a B as an A", A, Am1, B, e)
}

// Initial testing for select
func Test012(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { a A }"
	e := "B{A{}}.a"
	parseAndOkGood(t, A, B, e)
}

func Test012b(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { a A }"
	e := "B{A{}}.b"
	parseAndOkBad(t, "B does not have a \"b\" field", A, B, e)
}

func Test012c(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { a A }"
	e := "B{B{A{}}.a}"
	parseAndOkGood(t, A, B, e)
}

// Initial testing for call
func Test013(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1() A { return x0 }"
	e := "A{}.m1()"
	parseAndOkGood(t, A, A1m, e)
}

func Test013b(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1() A { return x0.m1() }"
	e := "A{}.m1()"
	parseAndOkGood(t, A, A1m, e)
}

func Test013c(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 A) A { return x1 }"
	e := "A{}.m1(A{})"
	parseAndOkGood(t, A, A1m, e)
}

func Test013d(t *testing.T) {
	fmt.Println("Source:")
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 A) A { return x1.m1(x0) }"
	e := "A{}.m1(A{}.m1(A{}))"
	parseAndOkGood(t, A, A1m, e)
}

func Test013e(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 A) A { return x0 }"
	e := "A{}.m1(A{}.m1())"
	parseAndOkBad(t, "(Nested) m1 call missing arg", A, A1m, e)
}

func Test013f(t *testing.T) {
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 A) A { return x0 }"
	e := "A{}.m1(A{}.m1(A{}, A{}))"
	parseAndOkBad(t, "(Nested) m1 call too many args", A, A1m, e)
}

func Test013g(t *testing.T) {
	fmt.Println("Source:")
	A := "type A struct {}"
	B := "type B struct { a A }"
	A1m := "func (x0 A) m1(x1 A) A { return x0 }"
	e := "A{}.m1(A{}.m1(B{A{}}))"
	parseAndOkBad(t, "(Nested) m1 call given a B, expecting an A", A, A1m, B, e)
}

/* Eval */

//func TestEval001(t *testing.T) { }
