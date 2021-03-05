/* See README.md for install and run examples.
 * See copyright.txt for copyright.
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

	"github.com/rhu1/fgg/internal/fg"
	"github.com/rhu1/fgg/internal/frontend"
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
	flag.IntVar(&oblitEvalSteps, "oblit-eval", frontend.NO_EVAL,
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
	frontend.OblitEvalSteps = oblitEvalSteps

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
			frontend.CheckErr(err)
		}
		src = string(b)
	}

	// Currently hacked
	if monomtest {
		frontend.TestMonom(printf, verbose, src, evalSteps)
		//testMonom(printf bool, verbose bool, src string, steps int) {
		return // FIXME
	} else if oblittest {
		frontend.TestOblit(verbose, src, evalSteps)
		//testOblit(verbose, src, evalSteps)  // TODO: "weak" oblit simulation
		return
	}

	switch { // Pre: !(interpFG && interpFGG)
	case interpFG:
		frontend.FGmain(verbose, src, strictParse, evalSteps, printf)
		// monom implicitly disabled
	case interpFGG:
		frontend.FGGmain(verbose, src, strictParse, evalSteps, printf, monom, monomc, oblitc)
	}
}

// For convenient quick testing -- via flag "-internal"
func internalSrc() string {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	e := "ToAny{1}"                        // CHECKME: `1` skipped by parser?
	return fg.MakeFgProgram(Any, ToAny, e) // FIXME: hardcoded FG
}
