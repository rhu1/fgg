// Pre(1): ANTLR4 Runtime for Go
//$ go get github.com/antlr/antlr4/runtime/Go/antlr
//
// Optional:
//$ cd [GOHOME]/src/github.com/antlr/antlr4
//$ git checkout -b antlr-go-runtime tags/4.7.1  // Match antlr-4.7.1-complete.jar -- but unnecessary

// Pre(2):
// [GOHOME]/src/github.com/rhu1/fgg
// $ mkdir parser/fg
// $ cp parser/pregren/fg/* parser/fg
// $ mkdir parser/fgg
// $ cp parser/pregren/fgg/* parser/fgg

// Run examples:
//$ go run github.com/rhu1/fgg -v -eval=10 fg/examples/hello/hello.go
//$ go run github.com/rhu1/fgg -v -inline="package main; type A struct {}; func main() { _ = A{} }"

// Optional alternative to Pre(2): ANTLR4 -- e.g., antlr-4.7.1-complete.jar
// Assuming "antlr4" alias for (e.g.): java -jar ~/code/java/lib/antlr-4.7.1-complete.jar
//$ go generate
// Cf. below:
//go:generate antlr4 -Dlanguage=Go -o parser/fg parser/FG.g4
//go:generate antlr4 -Dlanguage=Go -o parser/fgg parser/FGG.g4

// FGG gotchas:
// type B(type a Any) struct { f a }; // Any parsed as a TParam -- currently not permitted
// Node(Nat){...} // fgg.FGGNode (Nat) is fgg.TParam, not fgg.TName
// type IA(type ) interface { m1() };  // m1() parsed as a TName (an invalid Spec) -- N.B. ret missing anyway

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	//"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
	"github.com/rhu1/fgg/fgr"
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

// Command line parameters/flags
var (
	interpFG  bool // parse FG
	interpFGG bool // parse FGG

	monom  bool   // parse FGG and monomorphise FGG source -- paper notation (angle bracks)
	monomc string // output filename of monomorphised FGG; "--" for stdout -- Go output (no angle bracks)
	// TODO refactor naming between "monomc", "compile" and "oblitc"

	oblitc         string // output filename of FGR compilation via oblit; "--" for stdout
	oblitEvalSteps int    // TODO: Need an actual FGR syntax, for oblitc to concrete output

	monomtest bool
	oblittest bool

	useInternalSrc bool   // use internal source
	inlineSrc      string // use content of this as source
	strictParse    bool   // use strict parsing mode

	evalSteps int  // number of steps to evaluate
	verbose   bool // verbose mode
)

func init() {
	// FG or FGG
	flag.BoolVar(&interpFG, "fg", false,
		"interpret input as FG (defaults to true if neither -fg/-fgg set)")
	flag.BoolVar(&interpFGG, "fgg", false,
		"interpret input as FGG")

	// Erasure by monomorphisation -- implicitly disabled if not -fgg
	flag.BoolVar(&monom, "monom", false,
		"[WIP] monomorphise FGG source using paper notation, i.e., angle bracks (ignored if -fgg not set)")
	flag.StringVar(&monomc, "monomc", "", // Empty string for "false"
		"[WIP] monomorphise FGG source to (Go-compatible) FG, i.e., no angle bracks (ignored if -fgg not set)\n"+
			"specify '--' to print to stdout")

	// Erasure(?) by translation based on type reps -- FGG vs. FGR?
	flag.StringVar(&oblitc, "oblitc", "", // Empty string for "false"
		"[WIP] compile FGG source to FGR (ignored if -fgg not set)\n"+
			"specify '--' to print to stdout")
	flag.IntVar(&oblitEvalSteps, "oblit-eval", NO_EVAL,
		" N ⇒ evaluate N (≥ 0) steps; or\n-1 ⇒ evaluate to value (or panic)")

	// WIP
	flag.BoolVar(&monomtest, "test-monom", false, `[WIP] Test monom correctness`)
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

// TODO
// - refactor functionality into cmd dir
// - add type pres to monom test -- DONE
// - add tests for interface omega building
// - fix embedding monom -- DONE
// - fix monom name mangling -- partial: fix "commas"
// - fix parser nil vs. empty creation
// - WF check for duplicate decl names
// - WF recursive structs check
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

	// WIP
	if monomtest {
		testMonom(verbose, src, evalSteps)
		return // FIXME
	} else if oblittest {
		testOblit(verbose, src)
		//testOblit(verbose, src, evalSteps)  // TODO: "weak" oblit simulation
		return
	}

	switch { // Pre: !(interpFG && interpFGG)
	case interpFG:
		//var a fg.FGAdaptor
		//interp(&a, src, strictParse, evalSteps)
		intrp_fg := NewFGInterp(verbose, src, strictParse)
		if evalSteps > NO_EVAL {
			intrp_fg.Eval(evalSteps)
			fmt.Println(intrp_fg.GetProgram().GetMain())
		}
		// monom implicitly disabled
	case interpFGG:
		//var a fgg.FGGAdaptor
		//prog := interp(&a, src, strictParse, evalSteps)
		intrp_fgg := NewFGGInterp(verbose, src, strictParse)

		if evalSteps > NO_EVAL {
			intrp_fgg.Eval(evalSteps)
			fmt.Println(intrp_fgg.GetProgram().GetMain())
		}

		// TODO: further refactoring (cf. Frontend, Interp)
		intrp_fgg.Monom(monom, monomc)
		intrp_fgg.Oblit(oblitc)
		////doWrappers(prog, wrapperc)
	}
}

/* monom simulation check */

// TODO: refactor to cmd dir
func testMonom(verbose bool, src string, steps int) {
	intrp_fgg := NewFGGInterp(verbose, src, true)
	p_fgg := intrp_fgg.GetProgram().(fgg.FGGProgram)
	u := p_fgg.Ok(false).(fgg.TNamed)
	vPrintln(verbose, "\nFGG expr: "+p_fgg.GetMain().String())

	// (Initial) left-vertical arrow
	//p_mono := fgg.Monomorph(p_fgg)
	omega := fgg.GetOmega(p_fgg.GetDecls(), p_fgg.GetMain().(fgg.FGGExpr))
	p_mono := fgg.ApplyOmega(p_fgg, omega)
	vPrintln(verbose, "Monom expr: "+p_mono.GetMain().String())
	t := p_mono.Ok(false).(fg.Type)
	u_fg := fgg.ToMonomId(u)
	if !t.Equals(u_fg) {
		panic("-test-monom failed: types do not match\n\tFGG type=" + u.String() +
			" -> " + u_fg.String() + "\n\tmono=" + t.String())
	}

	done := steps > EVAL_TO_VAL
	for i := 0; i < steps || !done; i++ {
		if p_fgg.GetMain().IsValue() {
			break
		}
		// Repeat: horizontal arrows and right-vertical arrow
		p_fgg, u, p_mono = testMonomStep(verbose, omega, p_fgg, u, p_mono)
	}
	vPrintln(verbose, "\nFinished:\n\tfgg="+p_fgg.GetMain().String()+
		"\n\tmono="+p_mono.GetMain().String())
}

// Pre: u = p_fgg.Ok(), t = p_mono.Ok()
func testMonomStep(verbose bool, omega fgg.Omega, p_fgg fgg.FGGProgram,
	u fgg.TNamed, p_mono fg.FGProgram) (fgg.FGGProgram, fgg.TNamed,
	fg.FGProgram) {

	// Upper-horizontal arrow
	p1_fgg, _ := p_fgg.Eval()
	vPrintln(verbose, "\nEval FGG one step: "+p1_fgg.GetMain().String())
	u1 := p1_fgg.Ok(true).(fgg.TNamed)
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
	e_fgg := res.GetMain()
	e_mono := p1_mono.GetMain()
	vPrintln(verbose, "Monom of one step'd FGG: "+e_fgg.String())
	if e_fgg.String() != e_mono.String() {
		panic("-test-monom failed: exprs do not match\n\tFGG expr=" + e_fgg.String() +
			"\n\tmono=" + e_mono.String())
	}

	return p1_fgg.(fgg.FGGProgram), u1, p1_mono.(fg.FGProgram)
}

/* oblit "weak" simulation check */

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

// For convenient quick testing -- via flag "-internal"
func internalSrc() string {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	e := "ToAny{1}"                        // FIXME: `1` skipped by parser?
	return fg.MakeFgProgram(Any, ToAny, e) // FIXME: hardcoded FG
}

/* Helpers */

// ECheckErr
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

/**
TODO:
- mutual-poly-rec should blow up when ismonom check off -- omega sigs => t.m pairs
- struct-poly-rec should be monomable -- more aggressive method dropping in omega *building*; need to distinguish actual receiver types from other seen types, for applying omega to mdecls
- WF: e.g., repeat type decl
- add monom-eval commutativity check
- factor out more into base

	//b.WriteString("type B struct { f t };\n")  // TODO: unknown type
	//b.WriteString("type B struct { b B };\n")  // TODO: recursive struct
*/

// Alternative Run:
//$ go install
//$ $GOPATH/bin/fgg.exe ...
// N.B. GoInstall installs to $CYGHOME/code/go/bin (not $WINHOME)
