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
	inlineSrc   string // use content of this as source
	strictParse bool   // use strict parsing mode

	evalSteps int  // number of steps to evaluate
	verbose   bool // verbose mode
	printf    bool // use ToGoString for output (e.g., "main." type prefix)
)

func init() {
	flag.StringVar(&inlineSrc, "inline", "",
		`-inline="[FG src]", use inline input as source`)
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

	fgg [options] path/to/file.fg
	fgg [options] -inline "package main; type ...; func main() { ... }"

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

	intrpFg := frontend.NewFGInterp(verbose, src, strictParse)
	if evalSteps > frontend.NO_EVAL {
		intrpFg.Eval(evalSteps)
		frontend.PrintResult(printf, intrpFg.GetProgram())
	}
}
