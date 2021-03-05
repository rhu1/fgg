/* See copyright.txt for copyright.
 */

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rhu1/fgg/internal/frontend"
)

// Command line flags/parameters
var (
	monomc string // output filename of monomorphised FGG; "--" for stdout -- Go output (no angle bracks)

	oblitc         string // output filename of FGR compilation via oblit; "--" for stdout
	oblitEvalSteps int    // TODO: A concrete FGR syntax, for oblitc to output

	inlineSrc   string // use content of this as source
	strictParse bool   // use strict parsing mode

	evalSteps int  // number of steps to evaluate
	verbose   bool // verbose mode
	printf    bool // use ToGoString for output (e.g., "main." type prefix)
)

func init() {

	flag.StringVar(&monomc, "monomc", "", // Empty string for "false"
		"monomorphise FGG source to (Go-compatible) FG, i.e., no angle bracks (ignored if -fgg not set)\n"+
			"specify '--' to print to stdout")

	// Erasure(?) by translation based on type reps -- FGG vs. FGR?
	flag.StringVar(&oblitc, "oblitc", "", // Empty string for "false"
		"[WIP] compile FGG source to FGR (ignored if -fgg not set)\n"+
			"specify '--' to print to stdout")
	flag.IntVar(&oblitEvalSteps, "oblit-eval", -2,
		" N ⇒ evaluate N (≥ 0) steps; or\n-1 ⇒ evaluate to value (or panic)")

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

	fgg [options] path/to/file.fgg
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

	frontend.FGGmain(verbose, src, strictParse, evalSteps, printf, false, monomc, oblitc)
}
