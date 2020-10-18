//rhu@HZHL4 MINGW64 ~/code/go/src/github.com/rhu1/fgg
//$ go test github.com/rhu1/fgg/fg
//$ go test github.com/rhu1/fgg/fg -run Test001

package fgr_test // Separate package, can test "API"

/*import (
	"fmt"
	"testing"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/base/testutils"
	"github.com/rhu1/fgg/fgr"
)*/

/* Harness funcs * /

func parseAndOkGood(t *testing.T, elems ...string) base.Program {
	var adptr fg.FGAdaptor
	return testutils.ParseAndOkGood(t, &adptr, fg.MakeFgProgram(elems...))
}

// N.B. do not use to check for bad *syntax* -- see the "[Parser]" panic check in base.ParseAndOkBad
func parseAndOkBad(t *testing.T, msg string, elems ...string) base.Program {
	var adptr fg.FGAdaptor
	return testutils.ParseAndOkBad(t, msg, &adptr, fg.MakeFgProgram(elems...))
}
//*/

/* Syntax and typing */

// TOOD: make translation to FGG and compare results to -fgg

/*func Test001(t *testing.T) {
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	Am2 := "func (x0 A) m2(x1 A) A { return x1 }"
	Am3 := "func (x0 A) m3(x1 A, x2 A) A { return x2 }"
	B := "type B struct { f A }"
	e := "B{A{}}"
	parseAndOkGood(t, A, Am1, Am2, Am3, B, e)
}*/

/* Eval */

/*func TestEval001(t *testing.T) {
	A := "type A struct {}"
	B := "type B struct { f A }"
	e := "B{A{}}.f"
	prog := parseAndOkGood(t, A, B, e)
	testutils.EvalAndOkGood(t, prog, 1)
}*/
