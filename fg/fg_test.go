//rhu@HZHL4 MINGW64 ~/code/go/src/temp/antlr/antlr04
//$ go test temp/antlr/antlr04/fg
//$ go test temp/antlr/antlr04/fg -run Test001

package fg

import (
	"testing"
)

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
