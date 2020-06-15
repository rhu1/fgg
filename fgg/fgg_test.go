//$ go test github.com/rhu1/fgg/fgg
//$ go test github.com/rhu1/fgg/fgg -run Test001

package fgg_test // Separate package, can test "API"

import (
	"fmt"
	"testing"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/base/testutils"
	"github.com/rhu1/fgg/fgg"
)

/* Harness funcs */

func fggParseAndOkGood(t *testing.T, elems ...string) base.Program {
	var adptr fgg.FGGAdaptor
	p := testutils.ParseAndOkGood(t, &adptr,
		fgg.MakeFggProgram(elems...))
	return p
}

func fggParseAndOkMonomGood(t *testing.T, elems ...string) base.Program {
	p := fggParseAndOkGood(t, elems...).(fgg.FGGProgram)
	if ok, msg := fgg.IsMonomOK(p); !ok {
		t.Errorf("Unexpected nomono rejection:\n\t" + msg + "\n" +
			p.String())
	}
	return fgg.Monomorph(p)
}

// N.B. do not use to check for bad *syntax* -- see the PARSER_PANIC_PREFIX panic check in base.ParseAndOkBad
func fggParseAndOkBad(t *testing.T, msg string, elems ...string) base.Program {
	var adptr fgg.FGGAdaptor
	return testutils.ParseAndOkBad(t, msg, &adptr, fgg.MakeFggProgram(elems...))
}

// Based on testutils.EvalAndOkGood
// Pre: parseAndOkGood
func NomonoGood(t *testing.T, p fgg.FGGProgram) fgg.FGGProgram {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: " + fmt.Sprintf("%v", r) + "\n" +
				p.String())
		}
	}()
	if ok, msg := fgg.IsMonomOK(p); !ok {
		t.Errorf("Unexpected nomono rejection:\n\t" + msg + "\n" +
			p.String())
	}
	return p
}

// Based on testutils.EvalAndOkBad
// Pre: parseAndOkGood
func NomonoBad(t *testing.T, p fgg.FGGProgram, msg string) fgg.FGGProgram {
	if ok, _ := fgg.IsMonomOK(p); ok {
		t.Errorf("Expected nomono violation, but none occurred: " + msg + "\n" +
			p.String())
	}
	return p
}

/* Syntax and typing */

// TOOD: classify FG-compatible subset compare results to -fg

// Initial FGG test

// Initial FGG test
func Test001(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	B := "type B(type a Any()) struct { f a }"
	e := "B(A()){A(){}}"
	//type IA(type ) interface { m1(type )() Any };
	//type A1(type ) struct { };
	fggParseAndOkMonomGood(t, Any, A, B, e)
}

func Test001b(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	A1 := "type A1(type ) struct {}"
	B := "type B(type a Any()) struct { f a }"
	e := "B(A()){A1(){}}"
	fggParseAndOkBad(t, "A1() is not an A()", Any, A, A1, B, e)
}

// Testing StructLit typing, t_S OK
func Test002(t *testing.T) {
	IA := "type IA(type ) interface { m1(type )() A() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type )() A() { return x0 }"
	B := "type B(type a IA()) struct { f a }"
	e := "B(A()){A(){}}"
	fggParseAndOkMonomGood(t, IA, A, Am1, B, e)
}

func Test002b(t *testing.T) {
	IA := "type IA(type ) interface { m1(type )() A() }"
	A := "type A(type ) struct {}"
	B := "type B(type a IA()) struct { f a }"
	e := "B(A()){A(){}}"
	fggParseAndOkBad(t, "A() is not an A1()", IA, A, B, e)
}

// Testing fields (and t-args subs)
func Test003(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type )() Any() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type )() Any() { return x0 }"
	A1 := "type A1(type ) struct { }"
	B := "type B(type a IA()) struct { f a }"
	e := "B(A()){A(){}}"
	fggParseAndOkMonomGood(t, Any, IA, A, Am1, A1, B, e)
}

func Test003b(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type )() Any() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type )() Any() { return x0 }"
	A1 := "type A1(type ) struct { }"
	B := "type B(type a IA()) struct { f a }"
	e := "B(A()){A1(){}}"
	fggParseAndOkBad(t, "A1() is not an A()", Any, IA, A, Am1, A1, B, e)
}

// Initial testing for select on parameterised struct
func Test004(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct { fA Any() }"
	Am1 := "func (x0 A(type )) m1(type )() Any() { return x0 }"
	A1 := "type A1(type ) struct { }"
	B := "type B(type a Any()) struct { fB a }"
	e := "B(A()){A(){A1(){}}}.fB.fA"
	fggParseAndOkMonomGood(t, Any, A, Am1, A1, B, e)
}

func Test004b(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct { fA Any() }"
	Am1 := "func (x0 A(type )) m1(type )() Any() { return x0 }"
	A1 := "type A1(type ) struct { }"
	B := "type B(type a Any()) struct { fB a }"
	e := "B(A1()){A1(){}}.fB.fA"
	fggParseAndOkBad(t, "A1 has no field fA", Any, A, Am1, A1, B, e)
}

// Initial testing for call
func Test005(t *testing.T) {
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type )() A() { return x0 }"
	e := "A(){}.m1()()"
	fggParseAndOkMonomGood(t, A, Am1, e)
}

func Test006(t *testing.T) {
	IA := "type IA(type ) interface { m1(type a IA())() A() }"
	A := "type A(type ) struct {}"
	e := "A(){}"
	fggParseAndOkMonomGood(t, IA, A, e)
}

func Test006b(t *testing.T) {
	IA := "type IA(type ) interface { m1(type a A())() A() }"
	A := "type A(type ) struct {}"
	e := "A(){}"
	fggParseAndOkBad(t, "A() invalid upper bound", IA, A, e)
}

func Test007(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type a IA())() A() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() A() { return x0 }"
	A1 := "type A1(type ) struct {}"
	e := "A(){}.m1(A())()"
	fggParseAndOkMonomGood(t, Any, IA, A, Am1, A1, e)
}

func Test007b(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type a IA())() A() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() A() { return x0 }"
	A1 := "type A1(type ) struct {}"
	e := "A(){}.m1()()"
	fggParseAndOkBad(t, "Missing type actual", Any, IA, A, Am1, A1, e)
}

func Test007c(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type a IA())() A() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() A() { return x0 }"
	A1 := "type A1(type ) struct {}"
	e := "A(){}.m1(A1())()"
	fggParseAndOkBad(t, "A1() is not an IA()", Any, IA, A, Am1, A1, e)
}

func Test007d(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type a IA())() IA() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() IA() { return x0 }"
	A1 := "type A1(type ) struct {}"
	e := "A(){}.m1(A())()"
	fggParseAndOkMonomGood(t, Any, IA, A, Am1, A1, e)
}

// Testing Sig parsing
func Test008(t *testing.T) {
	IA := "type IA(type ) interface { m1(type a IA())() IA() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() IA() { return x0 }"
	B := "type B(type a IA()) struct {}"
	Bm2 := "func (x0 B(type a IA())) m2(type )(x1 a) B(a) { return x0 }"
	e := "A(){}"
	fggParseAndOkMonomGood(t, IA, A, Am1, B, Bm2, e)
}

// Testing calls on parameterised struct
func Test009(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type a IA())() IA() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() IA() { return x0 }"
	B := "type B(type a IA()) struct {}"
	Bm2 := "func (x0 B(type a IA())) m2(type )(x1 a) B(a) { return x0 }"
	e := "B(A()){}.m2()(A(){})"
	fggParseAndOkMonomGood(t, Any, IA, A, Am1, B, Bm2, e)
}

func Test009b(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type a IA())() IA() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() IA() { return x0 }"
	A1 := "type A1(type ) struct {}"
	B := "type B(type a IA()) struct {}"
	Bm2 := "func (x0 B(type a IA())) m2(type )(x1 a) B(a) { return x0 }"
	e := "B(A()){}.m2()(A1(){})"
	fggParseAndOkBad(t, "A1() is not an A()", Any, IA, A, Am1, A1, B, Bm2, e)
}

// Initial test for generic type assertion
func Test010(t *testing.T) {
	Any := "type Any(type ) interface {}"
	ToAny := "type ToAny(type ) struct { any Any() }"
	IA := "type IA(type ) interface { m1(type a IA())() IA() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() IA() { return x0 }"
	B := "type B(type a IA()) struct {}"
	Bm2 := "func (x0 B(type a IA())) m2(type )(x1 a) Any() { return x1 }" // Unnecessary
	e := "ToAny(){B(A()){}}.any.(B(A()))"
	fggParseAndOkMonomGood(t, Any, ToAny, IA, A, Am1, B, Bm2, e)
}

func Test011(t *testing.T) {
	IA := "type IA(type ) interface { m1(type a IA())() IA() }"
	ToIA := "type ToIA(type ) struct { upcast IA() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() IA() { return x0 }"
	e := "ToIA(){A(){}}.upcast.(A())"
	fggParseAndOkMonomGood(t, IA, ToIA, A, Am1, e)
}

func Test011b(t *testing.T) {
	IA := "type IA(type ) interface { m1(type a IA())() IA() }"
	ToIA := "type ToIA(type ) struct { upcast IA() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type a IA())() IA() { return x0 }"
	A1 := "type A1(type ) struct {}"
	e := "ToIA(){A(){}}.upcast.(A1())"
	fggParseAndOkBad(t, "A1() is not an IA", IA, ToIA, A, Am1, A1, e)
}

func Test011c(t *testing.T) {
	Any := "type Any(type ) interface {}"
	ToAny := "type ToAny(type ) struct { any Any() }"
	B := "type B(type ) struct {}"
	Bm3 := "func (x0 B(type )) m3(type b Any())(x1 b) Any() { return x1 }"
	e := "ToAny(){B(){}}"
	fggParseAndOkMonomGood(t, Any, ToAny, B, Bm3, e)
}

// Testing parsing for Call with both targ and arg
func Test012(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	B := "type B(type ) struct {}"
	Bm := "func (x0 B(type )) m(type a Any())(x1 a) a { return x1 }"
	e := "B(){}.m(A())(A(){})"
	fggParseAndOkMonomGood(t, Any, A, B, Bm, e)
}

// Testing Call typing, meth-tparam TSubs of result
func Test013(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	B := "type B(type a Any()) struct { f a }"
	Bm := "func (x0 B(type a Any())) m(type b Any())(x1 b) b { return x1 }"
	e := "B(A()){A(){}}.m(B(A()))(B(A()){A(){}}).f"
	fggParseAndOkMonomGood(t, Any, A, B, Bm, e)
}

// Testing u <: a, i.e., upper is open type param
func Test014(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	B := "type B(type a Any()) struct { f a }"
	Bm := "func (x0 B(type a Any())) m(type b Any())() b { return A(){} }"
	e := "B(A()){A(){}}.m(B(A()))(B(A()){A(){}}).f" // Eval would break type preservation, see TestEval001
	fggParseAndOkBad(t, Any, A, B, Bm, e)
}

// testing sigAlphaEquals
func Test015(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) interface { m(type a Any())(x a) Any() }"
	B := "type B(type ) interface { m(type b Any())(x b) Any() }"
	C := "type C(type ) struct {}"
	Cm := "func (x0 C(type )) m(type b Any())(x b) Any() { return x0 }"
	D := "type D(type ) struct {}"
	Dm := "func (x0 D(type )) foo(type )(x A()) Any() { return x0 }"
	e := "D(){}.foo()(C(){})"
	fggParseAndOkBad(t, Any, A, B, C, Cm, D, Dm, e)
}

// testing covariant receiver bounds (MDecl.OK) -- cf. map.fgg (memberBr)
func Test016(t *testing.T) {
	Any := "Any(type ) interface {}"
	A := "type A(type a Any()) interface { m(type )(x a) Any() }" // param must occur in a meth sig
	B := "type B(type a A(a)) struct {}"                          // must have recursive param
	Bm := "func (x0 B(type b A(b))) m(type )(x b) Any() { return x0 }"
	D := "type D(type ) struct{}"
	e := "D(){}"
	fggParseAndOkBad(t, Any, A, B, Bm, D, e)
}

func Test017(t *testing.T) {
	Any := "type Any(type ) interface {}"
	I := "type I(type ) interface { bar(type )() Any() }"
	A := "type A(type a Any()) struct { }"
	Afoo := "func (x0 A(type a I())) foo(type )() Any() { return x0 }"
	D := "type D(type ) struct { }"
	e := "A(D()){}.foo()()"
	fggParseAndOkBad(t, Any, I, A, Afoo, D, e)
}

func Test017b(t *testing.T) {
	Any := "type Any(type ) interface {}"
	I := "type I(type ) interface { bar(type )() Any() }"
	A := "type A(type a Any()) struct { }"
	Afoo := "func (x0 A(type a I())) foo(type )() Any() { return x0 }"
	D := "type D(type ) struct { }"
	Dbar := "func (x0 D(type )) bar(type )() Any() { return x0 }"
	e := "A(D()){}.foo()()"
	fggParseAndOkGood(t, Any, I, A, Afoo, D, Dbar, e)
}

/* Monom */

// TODO: isMonomorphisable -- should fail that check
/*
func TestMonom001(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	B := "type B(type a Any()) struct { f a }"
	Bm := "func (x0 B(type a Any())) m(type )() Any() { return B(B(a)){x0}.m()() }"
	e := "B(A()){A(){}}.m()()"
	parseAndOkBad(t, "Polymorphic recursion on the receiver type", Any, A, B, Bm, e)
}
//*/

//TODO: add -monom compose.fgg bug -- missing field type collection when visiting struct lits (e.g., Compose f, g types)
//TODO: add -monom map.fgg bug -- missing add-meth-param instans collection for interface type receivers (e.g., Bool().Cond(Bool())(...))

/* Eval */

// TOOD: classify FG-compatible subset compare results to -fg

func TestEval001(t *testing.T) {
	Any := "type Any(type ) interface {}"
	ToAny := "type ToAny(type ) struct { any Any() }"
	A := "type A(type ) struct {}"
	B := "type B(type a Any()) struct { f a }"
	Bm := "func (x0 B(type a Any())) m(type b Any())(x1 b) b { return ToAny(){A(){}}.any.(b) }"
	e := "B(A()){A(){}}.m(B(A()))(B(A()){A(){}}).f"
	prog := fggParseAndOkMonomGood(t, Any, ToAny, A, B, Bm, e)
	testutils.EvalAndOkBad(t, prog, "Cannot cast A() to B(A())", 3)
}

/* fmt.Sprintf */

func TestEval002(t *testing.T) {
	imp := "import \"fmt\""
	A := "type A(type ) struct {}"
	e := "fmt.Sprintf(\"\")"
	prog := fggParseAndOkMonomGood(t, imp, A, e)
	testutils.EvalAndOkGood(t, prog, 1)
}

func TestEval003(t *testing.T) {
	imp := "import \"fmt\""
	A := "type A(type ) struct {}"
	e := "fmt.Sprintf(\"%v ,_()+- %v\", A(){}, A(){})"
	prog := fggParseAndOkMonomGood(t, imp, A, e)
	testutils.EvalAndOkGood(t, prog, 1)
}

/* Nomono */

func TestNomono001(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type a Any()) struct {}"
	ma1 := "func (x0 A(type a Any())) ma1(type )() Any() { return x0.ma1()() }"
	e := "A(Any()){}.ma1()()"
	prog := fggParseAndOkGood(t, Any, A, ma1, e).(fgg.FGGProgram)
	NomonoGood(t, prog)
}

func TestNomono002(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type a Any()) struct {}"
	ma1 := "func (x0 A(type a Any())) ma1(type )() Any() { return A(A(a)){}.ma1()() }"
	e := "A(Any()){}.ma1()()" // Tests nomono collectExpr Delta (Call receiver)
	prog := fggParseAndOkGood(t, Any, A, ma1, e).(fgg.FGGProgram)
	NomonoBad(t, prog, "ma1 receiver polymorphic recursion, a -> A(A(a))")
}

func TestNomono003(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	ma1 := "func (x0 A(type )) ma1(type a Any())() Any() { return A(){}.ma1(a)() }"
	e := "A(){}.ma1(Any())()"
	prog := fggParseAndOkGood(t, Any, A, ma1, e).(fgg.FGGProgram)
	NomonoGood(t, prog)
}

func TestNomono004(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	ma1 := "func (x0 A(type )) ma1(type a Any())() Any() { return A(){}.ma1(A(a))() }"
	e := "A(){}.ma1(Any())()"
	prog := fggParseAndOkGood(t, Any, A, ma1, e).(fgg.FGGProgram)
	NomonoBad(t, prog, "ma1 meth-param polymorphic recursion, a -> A(a)")
}

func TestNomono005(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type a Any()) struct {}"
	ma1 := "func (x0 A(type a Any())) ma1(type )() Any() { return B(a){}.mb1()() }"
	B := "type B(type b Any()) struct {}"
	mb1 := "func (x0 B(type b Any())) mb1(type )() Any() { return A(b){}.ma1()() }"
	e := "A(Any()){}.ma1()()"
	prog := fggParseAndOkGood(t, Any, A, ma1, B, mb1, e).(fgg.FGGProgram)
	NomonoGood(t, prog)
}

func TestNomono006(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type a Any()) struct {}"
	ma1 := "func (x0 A(type a Any())) ma1(type )() Any() { return B(a){}.mb1()() }"
	B := "type B(type b Any()) struct {}"
	mb1 := "func (x0 B(type b Any())) mb1(type )() Any() { return A(A(b)){}.ma1()() }"
	e := "A(Any()){}.ma1()()"
	prog := fggParseAndOkGood(t, Any, A, ma1, B, mb1, e).(fgg.FGGProgram)
	NomonoBad(t, prog, "ma1 receiver polymorphic mutual recursion, a -> b -> A(A(b))")
}

func TestNomono007(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	ma1 := "func (x0 A(type )) ma1(type a Any())() Any() { return B(){}.mb1(a)() }"
	B := "type B(type ) struct {}"
	mb1 := "func (x0 B(type )) mb1(type b Any())() Any() { return A(){}.ma1(b)() }"
	e := "A(){}.ma1(Any())()"
	prog := fggParseAndOkGood(t, Any, A, ma1, B, mb1, e).(fgg.FGGProgram)
	NomonoGood(t, prog)
}

func TestNomono008(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	ma1 := "func (x0 A(type )) ma1(type a Any())() Any() { return B(){}.mb1(a)() }"
	B := "type B(type ) struct {}"
	mb1 := "func (x0 B(type )) mb1(type b Any())() Any() { return A(){}.ma1(A(b))() }"
	e := "A(){}.ma1(Any())()"
	//e := "A(){}"  // TODO: nomono conservative (ma1 not called)
	prog := fggParseAndOkGood(t, Any, A, ma1, B, mb1, e).(fgg.FGGProgram)
	NomonoBad(t, prog, "ma1+mb1 meth-param polymorphic mutual recursion, a -> b -> A(b)")
}

func TestNomono009(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { ma1(type )() Any() }"
	A := "type A(type ) struct {}"
	ma1 := "func (x0 A(type )) ma1(type )() Any() { return B(){}.mb1()() }"
	ma2 := "func (x0 A(type )) ma2(type )() Any() { return B(){}.mb3()() }"
	B := "type B(type ) struct {}"
	mb1 := "func (x0 B(type )) mb1(type )() Any() { return x0 }"
	mb2 := "func (x0 B(type )) mb2(type )() Any() { return x0.mb2()() }"
	mb3 := "func (x0 B(type )) mb3(type )() Any() { return A(){}.ma2()() }"
	C := "type C(type ) struct {}"
	mc1 := "func (x0 C(type )) ma1(type )() Any() { return x0 }" // C <: IA
	D := "type D(type ) struct {}"
	foo := "func (x0 D(type )) foo(type )(x IA()) Any() { return x.ma1()() }"
	e := "A(){}.ma1()()"
	prog := fggParseAndOkGood(t, Any, IA, A, ma1, ma2, B, mb1, mb2, mb3, C, mc1, D, foo, e).(fgg.FGGProgram)
	NomonoGood(t, prog)
}
