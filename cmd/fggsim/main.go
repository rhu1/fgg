/* See copyright.txt for copyright.
 */

package main

import (
	"flag"
	"fmt"
	"os"

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
		frontend.TestMonom(printf, verbose, src, evalSteps)
	}
	if oblittest {
		frontend.TestOblit(verbose, src, evalSteps)
		//testOblit(verbose, src, evalSteps)  // TODO: "weak" oblit simulation
	}
}
