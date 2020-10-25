/* See copyright.txt for copyright.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/rhu1/fgg/internal/base"
	"github.com/rhu1/fgg/internal/fg"
	"github.com/rhu1/fgg/internal/fgg"
	"github.com/rhu1/fgg/internal/fgr"
	"github.com/rhu1/fgg/internal/frontend"
)

// Command line parameters/flags
var (
	monomtest bool
	oblittest bool

	inlineSrc   string // use content of this as source
	strictParse bool   // use strict parsing mode

	evalSteps int  // number of steps to evaluate
	verbose   bool // verbose mode
	printf    bool // use ToGoString for output (e.g., "main." type prefix)
)

func init() {
	flag.BoolVar(&monomtest, "monom", false, `Test monom correctness`)
	flag.BoolVar(&oblittest, "oblit", false, `[WIP] Test oblit correctness`)

	// Parsing options
	flag.StringVar(&inlineSrc, "inline", "",
		`-inline="[FG/FGG src]", use inline input as source`)
	flag.BoolVar(&strictParse, "strict", true,
		"strict parsing (default true, means don't attempt recovery on parsing errors)")

	flag.IntVar(&evalSteps, "eval", frontend.NO_EVAL,
		" N ⇒ evaluate N (≥ 0) steps; or\n-1 ⇒ evaluate to value (or panic)")
	flag.BoolVar(&verbose, "v", false,
		"enable verbose printing")
	frontend.Verbose = verbose
	flag.BoolVar(&printf, "printf", false,
		"use Go style output type name prefixes")
}

var usage = func() {
	fmt.Fprintf(os.Stderr, `Usage:

	fggsim [options] path/to/file.fg
	fggsim [options] -inline "package main; type ...; func main() { ... }"

Options:

`)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Determine source
	var src string
	switch {
	case inlineSrc != "": // -inline overrules src file
		src = inlineSrc
	default:
		if flag.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Input error: need a source .go file (or an -inline program)")
			flag.Usage()
		}
		src = frontend.ReadSourceFile(flag.Arg(0))
	}

	// Currently hacked
	if monomtest {
		testMonom(printf, verbose, src, evalSteps)
	}
	if oblittest {
		testOblit(verbose, src)
		//testOblit(verbose, src, evalSteps)  // TODO: "weak" oblit simulation
	}
}

/* monom simulation check */

// TODO: refactor to cmd dir
func testMonom(printf bool, verbose bool, src string, steps int) {
	intrp_fgg := frontend.NewFGGInterp(verbose, src, true)
	p_fgg := intrp_fgg.GetProgram().(fgg.FGGProgram)
	u := p_fgg.Ok(false).(fgg.Type) // TNamed, except TParam for primitives (string)
	frontend.VPrintln(verbose, "\nFGG expr: "+p_fgg.GetMain().String())

	if ok, msg := fgg.IsMonomOK(p_fgg); !ok {
		frontend.VPrintln(verbose, "\nAborting simulation: Cannot monomorphise (nomono detected):\n\t"+msg)
		return
	}

	// (Initial) left-vertical arrow
	//p_mono := fgg.Monomorph(p_fgg)
	ds_fgg := p_fgg.GetDecls()
	omega := fgg.GetOmega(ds_fgg, p_fgg.GetMain().(fgg.FGGExpr))
	p_mono := fgg.ApplyOmega(p_fgg, omega) // TODO: can just monom expr (ground main) directly
	frontend.VPrintln(verbose, "Monom expr: "+p_mono.GetMain().String())
	t := p_mono.Ok(false).(fg.Type)
	ds_mono := p_mono.GetDecls()
	u_fg := fgg.ToMonomId(u)
	if !t.Equals(u_fg) {
		panic("-test-monom failed: types do not match\n\tFGG type=" + u.String() +
			" -> " + u_fg.String() + "\n\tmono=" + t.String())
	}

	done := steps > frontend.EVAL_TO_VAL
	var main_fgg base.Expr
	var main_mono base.Expr
	for i := 0; i < steps || !done; i++ {
		main_fgg = p_fgg.GetMain()
		main_mono = p_mono.GetMain()
		if main_fgg.IsValue() { // N.B. IsValue -- not CanEval (checked below)
			if !main_mono.IsValue() { // TODO: add to -test-oblit
				panic("FGG is value but monom is not:\n\tfgg = " + main_fgg.String() +
					"\n\tmonom=" + main_mono.String())
			}
			break // Both are values
		} else if main_mono.IsValue() {
			panic("Monom is value but FGG is not:\n\tfgg = " + main_fgg.String() +
				"\n\tmonom=" + main_mono.String())
		}
		// Both non-values, check for stuck (e.g., bad asserts -- though panic is technically not stuck)
		if main_fgg.CanEval(ds_fgg) {
			if !main_mono.CanEval(ds_mono) {
				panic("FGG is stuck but monom is not:\n\tfgg = " + main_fgg.String() +
					"\n\tmonom=" + main_mono.String())
			}
		} else {
			if main_mono.CanEval(ds_mono) {
				panic("Monom is stuck but FGG is not:\n\tfgg = " + main_fgg.String() +
					"\n\tmonom=" + main_mono.String())
			}
			if _, ok := main_fgg.(fgg.Assert); ok {
				if _, ok1 := main_mono.(fg.Assert); ok1 {
					break // Both stuck on bad assert
				}
			}
		}

		// Repeat: horizontal arrows and right-vertical arrow
		p_fgg, u, p_mono = testMonomStep(verbose, omega, p_fgg, u, p_mono)
	}
	frontend.VPrintln(verbose, "\nFinished:\n\tfgg="+p_fgg.GetMain().String()+
		"\n\tmono="+p_mono.GetMain().String())
}

// Pre: u = p_fgg.Ok(), t = p_mono.Ok(), both CanEval
func testMonomStep(verbose bool, omega fgg.Omega, p_fgg fgg.FGGProgram,
	u fgg.Type, p_mono fg.FGProgram) (fgg.FGGProgram, fgg.Type,
	fg.FGProgram) {

	// Upper-horizontal arrow
	p1_fgg, _ := p_fgg.Eval()
	frontend.VPrintln(verbose, "\nEval FGG one step: "+p1_fgg.GetMain().String())
	u1 := p1_fgg.Ok(true).(fgg.Type)    // TNamed, except TParam for primitives (string)
	if !u1.Impls(p_fgg.GetDecls(), u) { // TODO: factor out with Frontend.eval
		panic("-test-monom failed: type not preserved\n\tprev=" + u.String() +
			"\n\tnext=" + u1.String())
	}

	// Lower-horizontal arrow
	p1_mono, _ := p_mono.Eval()
	frontend.VPrintln(verbose, "Eval monom one step: "+p1_mono.GetMain().String())
	t1 := p1_mono.Ok(true).(fg.Type)
	u1_fg := fgg.ToMonomId(u1)
	if !t1.Equals(u1_fg) { // CHECKME: needed? or just do monom-level type preservation?
		panic("-test-monom failed: types do not match\n\tFGG type=" + u1.String() +
			" -> " + u1_fg.String() + "\n\tmono=" + t1.String())
	}

	// Right-vertical arrow
	//res := fgg.Monomorph(p1_fgg.(fgg.FGGProgram))
	res := fgg.ApplyOmega(p1_fgg.(fgg.FGGProgram), omega)
	e_fgg := res.GetMain() // N.B. the monom'd FGG expr (i.e., an FGExpr)
	e_mono := p1_mono.GetMain()
	frontend.VPrintln(verbose, "Monom of one step'd FGG: "+e_fgg.String())

	_, string_fgg := e_fgg.(fg.StringLit)
	_, string_mono := e_mono.(fg.StringLit)
	if !(string_fgg && string_mono) {
		// Replaced parens by angle bracks -- hack for StringLit (cf. fgg/examples/ooplsa20/fig6/expression.fgg)
		hacked := testMonomStringHack(e_fgg.(fg.FGExpr)).String()
		if hacked != e_mono.String() {
			panic("-test-monom failed: exprs do not match\n\tFGG expr=" + hacked +
				"\n\tmono    =" + e_mono.String())
		}
	}

	return p1_fgg.(fgg.FGGProgram), u1, p1_mono.(fg.FGProgram)
}

func testMonomStringHack(e1 fg.FGExpr) fg.FGExpr {
	switch e := e1.(type) {
	case fg.Variable:
		return e
	case fg.StructLit:
		elems := e.GetElems()
		es := make([]fg.FGExpr, len(elems))
		for i, v := range elems {
			es[i] = testMonomStringHack(v)
		}
		return fg.NewStructLit(e.GetType(), es)
	case fg.Select:
		return fg.NewSelect(testMonomStringHack(e.GetExpr()), e.GetField())
	case fg.Call:
		e_recv := testMonomStringHack(e.GetReceiver())
		args := e.GetArgs()
		es := make([]fg.FGExpr, len(args))
		for i, v := range args {
			es[i] = testMonomStringHack(v)
		}
		return fg.NewCall(e_recv, e.GetMethod(), es)
	case fg.Assert:
		return fg.NewAssert(testMonomStringHack(e.GetExpr()), e.GetType())
	case fg.StringLit:
		// HACK: currently works because users cannot write string literals, specifically '(' and ')'
		msg := e.GetValue()
		msg = strings.Replace(msg, "(", "<", -1)
		msg = strings.Replace(msg, ")", ">", -1)
		return fg.NewString(msg)
	case fg.Sprintf:
		args := e.GetArgs()
		es := make([]fg.FGExpr, len(args))
		for i, v := range args {
			es[i] = testMonomStringHack(v)
		}
		return fg.NewSprintf(e.GetFormat(), es)
	default:
		panic("Unknown FGExpr type: " + reflect.TypeOf(e1).String() + "\n\t" + e1.String())
	}
}

// For convenient quick testing -- via flag "-internal"
func internalSrc() string {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	e := "ToAny{1}"                        // CHECKME: `1` skipped by parser?
	return fg.MakeFgProgram(Any, ToAny, e) // FIXME: hardcoded FG
}

/*














 */

/*
 * TODO -- not currently up to date
 */

/* oblit "weak" simulation check */

// TODO: update following latest -test-monom
func testOblit(verbose bool, src string) {
	intrp_fgg := frontend.NewFGGInterp(verbose, src, true)
	p_fgg := intrp_fgg.GetProgram().(fgg.FGGProgram)
	u := p_fgg.Ok(false).(fgg.TNamed) // Ground
	frontend.VPrintln(verbose, "\nFGG expr: "+p_fgg.GetMain().String())

	// (Initial) left-vertical arrow
	p_oblit := fgr.Obliterate(intrp_fgg.GetSource().(fgg.FGGProgram))
	frontend.VPrintln(verbose, "Oblit expr: "+p_oblit.GetMain().String())
	t := p_oblit.Ok(false).(fgr.Type)
	t_fgr := fgr.ToFgrTypeFromBounds(make(fgg.Delta), u)
	if !t.Equals(t_fgr) {
		panic("-test-oblit failed: types do not match\n\tFGG type=" + u.String() +
			" -> " + t_fgr.String() + "\n\toblit=" + t.String())
	}

	// Horizontal+ arrows
	t1 := frontend.Eval(intrp_fgg, frontend.EVAL_TO_VAL).(fgg.Type)
	t1_fgr := fgr.ToFgrTypeFromBounds(make(fgg.Delta), t1)
	intrp_oblit := frontend.NewFGRInterp(verbose, p_oblit)
	t1_oblit := frontend.Eval(intrp_oblit, frontend.EVAL_TO_VAL)
	if !t1_oblit.Equals(t1_fgr) {
		panic("-test-oblit failed: types do not match\n\tFGG type=" + u.String() +
			" -> " + t1_fgr.String() + "\n\toblit=" + t.String())
	}

	// (Final) right-vertical arrow
	p1_fgg := intrp_fgg.GetProgram().(fgg.FGGProgram)
	e1_fgg := p1_fgg.GetMain()
	p1_oblit := intrp_oblit.GetProgram().(fgr.FGRProgram)
	e1_oblit := p1_oblit.GetMain()
	p1_fgr := fgr.Obliterate(p1_fgg)
	e1_fgr := p1_fgr.GetMain()
	if e1_fgr.String() != e1_oblit.String() {
		panic("-test-oblit failed: exprs do not correspond\n\tFGG expr=" + e1_fgg.String() +
			"\n\toblit=" + e1_oblit.String())
	}

	frontend.VPrintln(verbose, "\nFinished:\n\tfgg="+e1_fgg.String()+
		"\n\toblit="+e1_oblit.String())
}
