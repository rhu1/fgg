//rhu@HZHL4 MINGW64 ~/code/go/src/github.com/rhu1/fgg
//$ go test github.com/rhu1/fgg/fgg
//$ go test github.com/rhu1/fgg/fgg -run Test001

package fgg_test // Separate package, can test "API"

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rhu1/fgg/fgg"
)

/* Harness funcs */

func parseAndCheckOk(prog string) fgg.FGGProgram {
	var adptr fgg.FGGAdaptor
	ast := adptr.Parse(true, prog)
	allowStupid := false
	ast.Ok(allowStupid)
	return ast
}

func parseAndOkGood(t *testing.T, elems ...string) fgg.FGGProgram {
	prog := fgg.MakeFggProgram(elems...)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: " + fmt.Sprintf("%v", r) + "\n" +
				prog)
		}
	}()
	return parseAndCheckOk(prog)
}

// N.B. do not use to check for bad *syntax* -- see the "[Parser]" panic check
func parseAndOkBad(t *testing.T, msg string, elems ...string) fgg.FGGProgram {
	prog := fgg.MakeFggProgram(elems...)
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
	return parseAndCheckOk(prog)
}

/*
// Pre: parseAndOkGood
func evalAndOkGood(t *testing.T, p fgg.FGGProgram, steps int) fgg.FGGProgram {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: " + fmt.Sprintf("%v", r) + "\n" +
				p.String())
		}
	}()
	allowStupid := true
	for i := 0; i < steps; i++ {
		p, _ = p.Eval() // CHECKME: check rule names as part of test?
		p.Ok(allowStupid)
	}
	return p
}

// Pre: parseAndOkGood
func evalAndOkBad(t *testing.T, p fgg.FGGProgram, msg string, steps int) fgg.FGGProgram {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but none occurred: " + msg + "\n" +
				p.String())
		} else {
			// [Parser] panic should be already checked by parseAndOkGood
			// TODO FIXME: check panic more specifically
		}
	}()
	allowStupid := true
	for i := 0; i < steps; i++ {
		p, _ = p.Eval()
		p.Ok(allowStupid)
	}
	return p
}
*/

/* Syntax and typing */

// TOOD: classify FG-compatible subset compare results to -fg

// Initial FGG test
func Test001(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	B := "type B(type a Any()) struct { f a }"
	e := "B(A()){A(){}}"
	//type IA(type ) interface { m1(type )() Any };
	//type A1(type ) struct { };
	parseAndOkGood(t, Any, A, B, e)
}

func Test001b(t *testing.T) {
	Any := "type Any(type ) interface {}"
	A := "type A(type ) struct {}"
	A1 := "type A1(type ) struct {}"
	B := "type B(type a Any()) struct { f a }"
	e := "B(A()){A1(){}}"
	parseAndOkBad(t, "A1 is not an A", Any, A, A1, B, e)
}

// Testing StructLit typing, t_S OK
func Test002(t *testing.T) {
	IA := "type IA(type ) interface { m1(type )() A() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type )() A() { return x0 }"
	B := "type B(type a IA()) struct { f a }"
	e := "B(A()){A(){}}"
	parseAndOkGood(t, IA, A, Am1, B, e)
}

func Test002b(t *testing.T) {
	IA := "type IA(type ) interface { m1(type )() A() }"
	A := "type A(type ) struct {}"
	B := "type B(type a IA()) struct { f a }"
	e := "B(A()){A(){}}"
	parseAndOkBad(t, "A is not an A1", IA, A, B, e)
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
	parseAndOkGood(t, Any, IA, A, Am1, A1, B, e)
}

func Test003b(t *testing.T) {
	Any := "type Any(type ) interface {}"
	IA := "type IA(type ) interface { m1(type )() Any() }"
	A := "type A(type ) struct {}"
	Am1 := "func (x0 A(type )) m1(type )() Any() { return x0 }"
	A1 := "type A1(type ) struct { }"
	B := "type B(type a IA()) struct { f a }"
	e := "B(A()){A1(){}}"
	parseAndOkBad(t, "A1 is not an A", Any, IA, A, Am1, A1, B, e)
}

/* Eval */

// TOOD: classify FG-compatible subset compare results to -fg
