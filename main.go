// Pre (1): ANTLR4
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
//go:generate antlr4 -Dlanguage=Go -o parser/fg FG.g4
//go:generate antlr4 -Dlanguage=Go -o parser/fgg FGG.g4

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

var EVAL_TO_VAL = -1 // Must be < 0
var NO_EVAL = -2     // Must be < EVAL_TO_VAL

var verbose bool = false

// Gotchas:
// type B(type a Any) struct { f a }; // Any parsed as a TParam -- currently not permitted
// type IA(type ) interface { m1() };  // m1() parsed as a TName (an invalid Spec) -- N.B. ret missing anyway

// N.B. flags (e.g., -internal=true) must be supplied before any non-flag args
func main() {
	evalPtr := flag.Int("eval", NO_EVAL,
		"-steps=n, evaluate n (>=0) steps; or -steps=-1, evaluate to value (or panic)")
	fgPtr := flag.Bool("fg", false, "-fg=false, interpret input as FG (defaults to true if neither -fg/-fgg set)")
	fggPtr := flag.Bool("fgg", false, "-fgg=true, interpret input as FGG")
	internalPtr := flag.Bool("internal", false,
		"-internal=true, use \"internal\" input as source")
	inlinePtr := flag.String("inline", "",
		"-inline=true, use inline input as source")
	strictParsePtr := flag.Bool("strict", true,
		"-strict=false, disable strict parsing (attempt recovery on parsing errors)")
	verbosePtr := flag.Bool("v", false, "-v=true, enable verbose printing")
	flag.Parse()
	if !*fgPtr && !*fggPtr {
		*fgPtr = true
	}
	verbose = *verbosePtr

	var src string
	if *internalPtr { // First priority
		src = makeInternalSrc()
	} else if *inlinePtr != "" { // Second priority, i.e., -inline overrules src file arg
		src = *inlinePtr
	} else {
		if len(os.Args) < 2 {
			fmt.Println("Input error: need a source .go file (or an -inline program)")
		}
		bs, err := ioutil.ReadFile(os.Args[len(os.Args)-1])
		checkErr(err)
		src = string(bs)
	}

	if *fgPtr {
		interpFg(src, *strictParsePtr, *evalPtr)
	} else if *fggPtr {
		interpFgg(src, *strictParsePtr, *evalPtr)
	}
}

func interpFg(src string, strict bool, eval int) {
	vPrintln("\nParsing AST:")
	var adptr fg.FGAdaptor
	prog := adptr.Parse(strict, src) // AST (FGProgram root)
	vPrintln(prog.String())

	vPrintln("\nChecking source program OK:")
	allowStupid := false
	prog.Ok(allowStupid)

	if eval > NO_EVAL {
		evalFg(prog, eval)
	}
}

func interpFgg(src string, strict bool, eval int) {
	vPrintln("\nParsing AST:")
	var adptr fgg.FGGAdaptor
	prog := adptr.Parse(strict, src) // AST (FGProgram root)
	vPrintln(prog.String())

	vPrintln("\nChecking source program OK:")
	allowStupid := false
	prog.Ok(allowStupid)

	/*if eval > NO_EVAL {
		eval(prog, eval)
	}*/
}

// N.B. currently FG panic comes out implicitly as an underlying run-time panic
// TODO: add explicit FG panics
// If steps == EVAL_TO_VAL, then eval to value
func evalFg(p fg.FGProgram, steps int) {
	allowStupid := true
	vPrintln("\nEntering Eval loop:")
	vPrintln("Decls:")
	for _, v := range p.GetDecls() {
		vPrintln("\t" + v.String() + ";")
	}
	vPrintln("Eval steps:")
	vPrintln(fmt.Sprintf("%6d: %8s %v", 0, "", p.GetExpr())) // Initial prog OK already checked

	done := steps > EVAL_TO_VAL || // Ignore 'done' if num steps fixed (set true, for ||!done below)
		fg.IsValue(p.GetExpr()) // O/w evaluate until a val -- here, check if init expr is already a val
	var rule string
	for i := 1; i <= steps || !done; i++ {
		p, rule = p.Eval()
		vPrintln(fmt.Sprintf("%6d: %8s %v", i, "["+rule+"]", p.GetExpr()))
		vPrintln("Checking OK:") // TODO: maybe disable by default, enable by flag
		p.Ok(allowStupid)
		if !done && fg.IsValue(p.GetExpr()) {
			done = true
		}
	}
	fmt.Println(p.GetExpr().String()) // Final result
}

// For convenient quick testing -- via flag "-internal=true"
func makeInternalSrc() string {
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

func vPrintln(x string) {
	if verbose {
		fmt.Println(x)
	}
}

/* TODO
- WF: repeat type decl

	//b.WriteString("type B struct { f t };\n")  // TODO: unknown type
	//b.WriteString("type B struct { b B };\n")  // TODO: recursive struct
*/
