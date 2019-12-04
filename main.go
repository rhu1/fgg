// Pre (1): ANTLR4
/*HERE
- getTypeReps
- add meth-param RepDecls
- FGR eval
*/
// E.g., antlr-4.7.1-complete.jar
// (See go:generate below)

// Pre (2): ANTLR4 Runtime for Go
//$ go get github.com/antlr/antlr4/runtime/Go/antlr
// Optional:
//$ cd $CYGHOME/code/go/src/github.com/antlr/antlr4
//$ git checkout -b antlr-go-runtime tags/4.7.1  // Match antlr-4.7.1-complete.jar -- but unnecessary

//rhu@HZHL4 MINGW64 ~/code/go/src/
//$ go run github.com/rhu1/fgg -v -eval=10 fg/examples/hello/hello.go
//$ go run github.com/rhu1/fgg -v -inline="package main; type A struct {}; func main() { _ = A{} }"
// or
//$ go install
//$ /c/Users/rhu/code/go/bin/fgg.exe ...

// N.B. GoInstall installs to $CYGHOME/code/go/bin (not $WINHOME)

// Assuming "antlr4" alias for (e.g.): java -jar ~/code/java/lib/antlr-4.7.1-complete.jar
//go:generate antlr4 -Dlanguage=Go -o parser/fg parser/FG.g4
//go:generate antlr4 -Dlanguage=Go -o parser/fgg parser/FGG.g4

// FGG gotchas:
// type B(type a Any) struct { f a }; // Any parsed as a TParam -- currently not permitted
// Node(Nat){...} // fgg.FGGNode (Nat) is fgg.TParam, not fgg.TName
// type IA(type ) interface { m1() };  // m1() parsed as a TName (an invalid Spec) -- N.B. ret missing anyway

/* TODO
- WF: repeat type decl
- Compact getters/constructors
- factor out fgg_monom, MDecl and ITypeLit
- fix type preservation check
- factor out more into base

	//b.WriteString("type B struct { f t };\n")  // TODO: unknown type
	//b.WriteString("type B struct { b B };\n")  // TODO: recursive struct
*/

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

var _ = reflect.TypeOf
var _ = strconv.Itoa
var _ = fgr.NewCall

const (
	EVAL_TO_VAL = -1 // Must be < 0
	NO_EVAL     = -2 // Must be < EVAL_TO_VAL
)

// Command line parameters/flags
var (
	interpFG  bool // parse FG
	interpFGG bool // parse FGG

	monom  bool   // parse FGG and monomorphise FGG source -- paper notation (angle bracks)
	monomc string // output filename of monomorphised FGG; "--" for stdout -- Go output (no angle bracks)
	// TODO refactor naming between "monomc", "compile" and "oblitc"

	oblitc         string // output filename of FGR compilation via oblit; "--" for stdout
	oblitEvalSteps int    // TODO: Need an actual FGR syntax, for oblitc to concrete output

	useInternalSrc bool   // use internal source
	inlineSrc      string // use content of this as source
	strictParse    bool   // use strict parsing mode

	evalSteps int  // number of steps to evaluate
	verbose   bool // verbose mode

	printf bool // Use fmt.Printf for main expr  // CHECKME: currently unused?
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

	// Parsing options
	flag.BoolVar(&useInternalSrc, "internal", false,
		`use "internal" input as source`)
	flag.StringVar(&inlineSrc, "inline", "",
		`-inline="[FG/FGG src]", use inline input as source`)
	flag.BoolVar(&strictParse, "strict", true,
		"strict parsing (don't attempt recovery on parsing errors)")

	flag.IntVar(&evalSteps, "eval", NO_EVAL,
		" N ⇒ evaluate N (≥ 0) steps; or\n-1 ⇒ evaluate to value (or panic)")
	flag.BoolVar(&verbose, "v", false,
		"enable verbose printing")

	flag.BoolVar(&printf, "printf", false,
		"Use fmt.Printf for main expr")
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
		src = internalSrc()
	case inlineSrc != "": // Second priority, i.e. -inline overrules src file
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

	switch { // Pre: !(interpFG && interpFGG)
	case interpFG:
		var a fg.FGAdaptor
		interp(&a, src, strictParse, evalSteps)
		// monom implicitly disabled
	case interpFGG:
		var a fgg.FGGAdaptor
		prog := interp(&a, src, strictParse, evalSteps)

		// TODO: refactor
		doMonom(prog, monom, monomc)
		//doWrappers(prog, wrapperc)
		doOblit(prog, oblitc)
	}
}

func interp(a base.Adaptor, src string, strict bool, steps int) base.Program {
	vPrintln("\nParsing AST:")
	prog := a.Parse(strict, src) // AST (Program root)
	vPrintln(prog.String())

	vPrintln("\nChecking source program OK:")
	allowStupid := false
	prog.Ok(allowStupid)

	if steps > NO_EVAL {
		eval(prog, steps)
	}

	return prog
}

// N.B. currently FG panic comes out implicitly as an underlying run-time panic
// TODO: add explicit FG panics
// If steps == EVAL_TO_VAL, then eval to value
func eval(p base.Program, steps int) {
	allowStupid := true
	vPrintln("\nEntering Eval loop:")
	vPrintln("Decls:")
	for _, v := range p.GetDecls() {
		vPrintln("\t" + v.String() + ";")
	}
	vPrintln("Eval steps:")
	vPrintln(fmt.Sprintf("%6d: %8s %v", 0, "", p.GetMain())) // Initial prog OK already checked

	done := steps > EVAL_TO_VAL || // Ignore 'done' if num steps fixed (set true, for `||!done` below)
		p.GetMain().IsValue() // O/w evaluate until a val -- here, check if init expr is already a val
	var rule string
	for i := 1; i <= steps || !done; i++ {
		p, rule = p.Eval()
		vPrintln(fmt.Sprintf("%6d: %8s %v", i, "["+rule+"]", p.GetMain()))
		vPrintln("Checking OK:") // TODO: maybe disable by default, enable by flag
		// TODO FIXME: check actual type preservation (not just typeability)
		p.Ok(allowStupid)
		if !done && p.GetMain().IsValue() {
			done = true
		}
	}
	fmt.Println(p.GetMain().ToGoString()) // Final result
}

// Pre: (monom == true || compile != "") => -fgg is set
// TODO: rename
func doMonom(prog base.Program, monom bool, compile string) {
	if !monom && compile == "" {
		return
	}
	p_mono := fgg.Monomorph(prog.(fgg.FGGProgram)) // TODO: reformat (e.g., "<...>") to make an actual FG program
	if monom {
		vPrintln("\nMonomorphising, formal notation: [Warning] WIP [Warning]")
		fmt.Println(p_mono.String())
	}
	if compile != "" {
		vPrintln("\nMonomorphising, FG output: [Warning] WIP [Warning]")
		out := p_mono.String()
		out = strings.Replace(out, ",,", "", -1)
		out = strings.Replace(out, "<", "", -1)
		out = strings.Replace(out, ">", "", -1)
		if compile == "--" {
			fmt.Println(out)
		} else {
			vPrintln(out)
			vPrintln("Writing output to: " + compile)
			bs := []byte(out)
			err := ioutil.WriteFile(compile, bs, 0644)
			checkErr(err)
		}
	}
}

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

func doOblit(prog base.Program, compile string) {
	if compile == "" {
		return
	}
	vPrintln("\nTranslating FGG to FG(R) using Obliteration: [Warning] WIP [Warning]")
	p_fgr := fgr.Obliterate(prog.(fgg.FGGProgram))
	out := p_fgr.String()
	// TODO: factor out with -monomc
	if compile == "--" {
		fmt.Println(out)
	} else {
		vPrintln(out)
		vPrintln("Writing output to: " + compile)
		bs := []byte(out)
		err := ioutil.WriteFile(compile, bs, 0644)
		checkErr(err)
	}

	// cf. interp -- TODO: refactor
	p_fgr.Ok(false)
	if oblitEvalSteps > NO_EVAL {
		vPrint("\nEvaluating FGR:") // eval prints a leading "\n"
		eval(p_fgr, oblitEvalSteps)
	}
}

// For convenient quick testing -- via flag "-internal=true"
func internalSrc() string {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	e := "ToAny{1}"
	return fg.MakeFgProgram(Any, ToAny, e)
}

/* Helpers */

// ECheckErr
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func vPrint(x string) {
	if verbose {
		fmt.Print(x)
	}
}

func vPrintln(x string) {
	if verbose {
		fmt.Println(x)
	}
}
