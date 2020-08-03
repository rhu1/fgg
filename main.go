/* See README.md for install and run examples.
 *
 * Copyright 2020 The fgg authors.
 */

//go:generate antlr4 -Dlanguage=Go -o parser/fg parser/FG.g4
//go:generate antlr4 -Dlanguage=Go -o parser/fgg parser/FGG.g4

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
	"github.com/rhu1/fgg/fgr"
)

var _ = reflect.Append
var _ = strconv.AppendBool

// Command line parameters/flags
var (
	interpFG  bool // parse FG
	interpFGG bool // parse FGG

	monom  bool   // parse FGG and monomorphise FGG source -- paper notation (angle bracks)
	monomc string // output filename of monomorphised FGG; "--" for stdout -- Go output (no angle bracks)
	// TODO: fix naming between "monomc", "compile" and "oblitc"

	oblitc         string // output filename of FGR compilation via oblit; "--" for stdout
	oblitEvalSteps int    // TODO: A concrete FGR syntax, for oblitc to output

	monomtest bool
	oblittest bool

	useInternalSrc bool   // use internal source
	inlineSrc      string // use content of this as source
	strictParse    bool   // use strict parsing mode

	evalSteps int  // number of steps to evaluate
	verbose   bool // verbose mode
	printf    bool // use ToGoString for output (e.g., "main." type prefix)
)

func init() {
	// FG or FGG
	flag.BoolVar(&interpFG, "fg", false,
		"interpret input as FG (defaults to true if neither -fg/-fgg set)")
	flag.BoolVar(&interpFGG, "fgg", false,
		"interpret input as FGG")

	// Erasure by monomorphisation -- implicitly disabled if not -fgg
	flag.BoolVar(&monom, "monom", false,
		"monomorphise FGG source using paper notation, i.e., angle bracks (ignored if -fgg not set)")
	flag.StringVar(&monomc, "monomc", "", // Empty string for "false"
		"monomorphise FGG source to (Go-compatible) FG, i.e., no angle bracks (ignored if -fgg not set)\n"+
			"specify '--' to print to stdout")

	// Erasure(?) by translation based on type reps -- FGG vs. FGR?
	flag.StringVar(&oblitc, "oblitc", "", // Empty string for "false"
		"[WIP] compile FGG source to FGR (ignored if -fgg not set)\n"+
			"specify '--' to print to stdout")
	flag.IntVar(&oblitEvalSteps, "oblit-eval", NO_EVAL,
		" N ⇒ evaluate N (≥ 0) steps; or\n-1 ⇒ evaluate to value (or panic)")

	flag.BoolVar(&monomtest, "test-monom", false, `Test monom correctness`)
	flag.BoolVar(&oblittest, "test-oblit", false, `[WIP] Test oblit correctness`)

	// Parsing options
	flag.BoolVar(&useInternalSrc, "internal", false,
		`use "internal" input as source`)
	flag.StringVar(&inlineSrc, "inline", "",
		`-inline="[FG/FGG src]", use inline input as source`)
	flag.BoolVar(&strictParse, "strict", true,
		"strict parsing (default true, means don't attempt recovery on parsing errors)")

	flag.IntVar(&evalSteps, "eval", NO_EVAL,
		" N ⇒ evaluate N (≥ 0) steps; or\n-1 ⇒ evaluate to value (or panic)")
	flag.BoolVar(&verbose, "v", false,
		"enable verbose printing")
	flag.BoolVar(&printf, "printf", false,
		"use Go style output type name prefixes")
}

var usage = func() {
	fmt.Fprintf(os.Stderr, `Usage:

	fgg [options] -fg  path/to/file.fg
	fgg [options] -fgg path/to/file.fgg
	fgg [options] -internal
	fgg [options] -inline "package main; type ...; func main() { ... }"

Options:

`)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Determine (default) mode
	if interpFG {
		if interpFGG { // -fg "overrules" -fgg
			interpFGG = false
		}
	} else if !interpFGG {
		interpFG = true // -fg default
	}

	// Determine source
	var src string
	switch {
	case useInternalSrc: // First priority
		src = internalSrc() // FIXME: hardcoded to FG
	case inlineSrc != "": // Second priority, i.e., -inline overrules src file
		src = inlineSrc
	default:
		if flag.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Input error: need a source .go file (or an -inline program)")
			flag.Usage()
		}
		b, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			checkErr(err)
		}
		src = string(b)
	}

	// Currently hacked
	if monomtest {
		testMonom(printf, verbose, src, evalSteps)
		return // FIXME
	} else if oblittest {
		testOblit(verbose, src)
		//testOblit(verbose, src, evalSteps)  // TODO: "weak" oblit simulation
		return
	}

	switch { // Pre: !(interpFG && interpFGG)
	case interpFG:
		intrp_fg := NewFGInterp(verbose, src, strictParse)
		if evalSteps > NO_EVAL {
			intrp_fg.Eval(evalSteps)
			printResult(printf, intrp_fg.GetProgram())
		}
		// monom implicitly disabled
	case interpFGG:
		intrp_fgg := NewFGGInterp(verbose, src, strictParse)

		if evalSteps > NO_EVAL {
			intrp_fgg.Eval(evalSteps)
			printResult(printf, intrp_fgg.GetProgram())
		}

		// TODO: refactor (cf. Frontend, Interp)
		intrp_fgg.Monom(monom, monomc)
		intrp_fgg.Oblit(oblitc)
		////doWrappers(prog, wrapperc)
	}
}

func printResult(printf bool, p base.Program) {
	res := p.GetMain()
	if printf {
		fmt.Println(res.ToGoString(p.GetDecls()))
	} else {
		fmt.Println(res)
	}
}

/* monom simulation check */

// TODO: refactor to cmd dir
func testMonom(printf bool, verbose bool, src string, steps int) {
	intrp_fgg := NewFGGInterp(verbose, src, true)
	p_fgg := intrp_fgg.GetProgram().(fgg.FGGProgram)
	u := p_fgg.Ok(false).(fgg.Type) // TNamed, except TParam for primitives (string)
	vPrintln(verbose, "\nFGG expr: "+p_fgg.GetMain().String())

	if ok, msg := fgg.IsMonomOK(p_fgg); !ok {
		vPrintln(verbose, "\nAborting simulation: Cannot monomorphise (nomono detected):\n\t"+msg)
		return
	}

	// (Initial) left-vertical arrow
	//p_mono := fgg.Monomorph(p_fgg)
	ds_fgg := p_fgg.GetDecls()
	omega := fgg.GetOmega(ds_fgg, p_fgg.GetMain().(fgg.FGGExpr))
	p_mono := fgg.ApplyOmega(p_fgg, omega) // TODO: can just monom expr (ground main) directly
	vPrintln(verbose, "Monom expr: "+p_mono.GetMain().String())
	t := p_mono.Ok(false).(fg.Type)
	ds_mono := p_mono.GetDecls()
	u_fg := fgg.ToMonomId(u)
	if !t.Equals(u_fg) {
		panic("-test-monom failed: types do not match\n\tFGG type=" + u.String() +
			" -> " + u_fg.String() + "\n\tmono=" + t.String())
	}

	done := steps > EVAL_TO_VAL
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
		} else if main_mono.CanEval(ds_mono) {
			panic("Monom is stuck but FGG is not:\n\tfgg = " + main_fgg.String() +
				"\n\tmonom=" + main_mono.String())
		}

		// Repeat: horizontal arrows and right-vertical arrow
		p_fgg, u, p_mono = testMonomStep(verbose, omega, p_fgg, u, p_mono)
	}
	vPrintln(verbose, "\nFinished:\n\tfgg="+p_fgg.GetMain().String()+
		"\n\tmono="+p_mono.GetMain().String())
}

// Pre: u = p_fgg.Ok(), t = p_mono.Ok(), both CanEval
func testMonomStep(verbose bool, omega fgg.Omega, p_fgg fgg.FGGProgram,
	u fgg.Type, p_mono fg.FGProgram) (fgg.FGGProgram, fgg.Type,
	fg.FGProgram) {

	// Upper-horizontal arrow
	p1_fgg, _ := p_fgg.Eval()
	vPrintln(verbose, "\nEval FGG one step: "+p1_fgg.GetMain().String())
	u1 := p1_fgg.Ok(true).(fgg.Type)    // TNamed, except TParam for primitives (string)
	if !u1.Impls(p_fgg.GetDecls(), u) { // TODO: factor out with Frontend.eval
		panic("-test-monom failed: type not preserved\n\tprev=" + u.String() +
			"\n\tnext=" + u1.String())
	}

	// Lower-horizontal arrow
	p1_mono, _ := p_mono.Eval()
	vPrintln(verbose, "Eval monom one step: "+p1_mono.GetMain().String())
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
	vPrintln(verbose, "Monom of one step'd FGG: "+e_fgg.String())

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

/* Helpers */

// ECheckErr
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

/*









































 */

/*
 * TODO -- not currently up to date
 */

/* oblit "weak" simulation check */

// TODO: update following latest -test-monom
func testOblit(verbose bool, src string) {
	intrp_fgg := NewFGGInterp(verbose, src, true)
	p_fgg := intrp_fgg.GetProgram().(fgg.FGGProgram)
	u := p_fgg.Ok(false).(fgg.TNamed) // Ground
	vPrintln(verbose, "\nFGG expr: "+p_fgg.GetMain().String())

	// (Initial) left-vertical arrow
	p_oblit := fgr.Obliterate(intrp_fgg.GetSource().(fgg.FGGProgram))
	vPrintln(verbose, "Oblit expr: "+p_oblit.GetMain().String())
	t := p_oblit.Ok(false).(fgr.Type)
	t_fgr := fgr.ToFgrTypeFromBounds(make(fgg.Delta), u)
	if !t.Equals(t_fgr) {
		panic("-test-oblit failed: types do not match\n\tFGG type=" + u.String() +
			" -> " + t_fgr.String() + "\n\toblit=" + t.String())
	}

	// Horizontal+ arrows
	t1 := eval(intrp_fgg, EVAL_TO_VAL).(fgg.Type)
	t1_fgr := fgr.ToFgrTypeFromBounds(make(fgg.Delta), t1)
	intrp_oblit := NewFGRInterp(verbose, p_oblit)
	t1_oblit := eval(intrp_oblit, EVAL_TO_VAL)
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

	vPrintln(verbose, "\nFinished:\n\tfgg="+e1_fgg.String()+
		"\n\toblit="+e1_oblit.String())
}

//*/

/*
func testOblit(verbose bool, src string, steps int) {
	intrp_fgg := NewFGGInterp(verbose, src, true)
	p_fgg := intrp_fgg.GetProgram().(fgg.FGGProgram)
	u := p_fgg.Ok(false).(fgg.TNamed) // Ground
	vPrintln(verbose, "\nFGG expr: "+p_fgg.GetMain().String())

	// (Initial) left-vertical arrow
	p_oblit := fgr.Obliterate(intrp_fgg.GetSource().(fgg.FGGProgram))
	vPrintln(verbose, "Oblit expr: "+p_oblit.GetMain().String())
	t := p_oblit.Ok(false).(fgr.Type)
	if !t.Equals(fgr.ToFgrTypeFromBounds(make(fgg.Delta), u)) {
		panic("-test-oblit failed: types do not match\n\tFGG type=" + u.String() +
			" -> " + fgg.ToMonomId(u).String() + "\n\toblit=" + t.String())
	}

	done := steps > EVAL_TO_VAL
	for i := 0; i < steps || !done; i++ {
		if p_fgg.GetMain().IsValue() {
			break
		}
		// Repeat: horizontal arrows and right-vertical arrow
		p_fgg, u, p_oblit = testOblitStep(verbose, p_fgg, u, p_oblit)
	}
	vPrintln(verbose, "\nFinished:\n\tfgg="+p_fgg.GetMain().String()+
		"\n\toblit="+p_oblit.GetMain().String())
}

// Pre: u = p_fgg.Ok(), t = p_fgr.Ok()
func testOblitStep(verbose bool, p_fgg fgg.FGGProgram, u fgg.TNamed,
	p_oblit fgr.FGRProgram) (fgg.FGGProgram, fgg.TNamed, fgr.FGRProgram) {

	// Upper-horizontal arrow
	p1_fgg, _ := p_fgg.Eval()
	vPrintln(verbose, "\nEval FGG one step: "+p1_fgg.GetMain().String())
	u1 := p1_fgg.Ok(true).(fgg.TNamed)  // Ground
	if !u1.Impls(p_fgg.GetDecls(), u) { // TODO: factor out with Frontend.eval
		panic("-test-oblit failed: type not preserved\n\tprev=" + u.String() +
			"\n\tnext=" + u1.String())
	}

	// Lower-horizontal arrow -- FIXME: need to greedily do "weak" inserted asserts
	p1_oblit, _ := p_oblit.Eval()
	vPrintln(verbose, "Eval oblit one step: "+p1_oblit.GetMain().String())
	t1 := p1_oblit.Ok(true).(fgr.Type)
	if !t1.Equals(fgr.ToFgrTypeFromBounds(make(fgg.Delta), u1)) { // CHECKME: needed? or just do monom-level type preservation?
		panic("-test-oblit failed: types do not match\n\tFGG type=" + u1.String() +
			" -> " + fgg.ToMonomId(u1).String() + "\n\toblit=" + t1.String())
	}

	// Right-vertical arrow
	res := fgr.Obliterate(p1_fgg.(fgg.FGGProgram))
	e_fgg := res.GetMain()
	e_fgr := p1_oblit.GetMain()
	if e_fgg.IsValue() { // FIXME failed hack, number of eval steps don't correspond
		vPrintln(verbose, "Oblit of one step'd FGG: "+e_fgg.String())
		if e_fgg.String() != e_fgr.String() {
			panic("-test-oblit failed: exprs do not match\n\tFGG expr=" + e_fgg.String() +
				"\n\toblit=" + e_fgr.String())
		}
	} else if e_fgr.IsValue() {
		panic("-test-oblit failed: exprs do not match\n\tFGG expr=" + e_fgg.String() +
			"\n\toblit=" + e_fgr.String())
	}

	return p1_fgg.(fgg.FGGProgram), u1, p1_oblit.(fgr.FGRProgram)
}
//*/

/* [WIP] TODO -- not functional yet
func doWrappers(prog base.Program, compile string) {
	if compile == "" {
		return
	}
	vPrintln("\nTranslating FGG to FG(R) using Wrappers: [Warning] WIP [Warning]")
	//p_fgr := fgg.FgAdptrTranslate(prog.(fgg.FGGProgram))
	//p_fgr := fgg.FgrTranslate(prog.(fgg.FGGProgram))
	p_fgr := fgr.Translate(prog.(fgg.FGGProgram))
	out := p_fgr.String()
	// TODO: factor out with -monomc
	if compile == "--" {
		fmt.Println(out)
	} else {
		vPrintln("Writing output to: " + compile)
		bs := []byte(out)
		err := ioutil.WriteFile(compile, bs, 0644)
		checkErr(err)
	}
}
//*/
