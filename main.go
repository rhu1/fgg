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
//go:generate antlr4 -Dlanguage=Go -o parser/fg parser/FG.g4
//go:generate antlr4 -Dlanguage=Go -o parser/fgg parser/FGG.g4

// FGG gotchas:
// type B(type a Any) struct { f a }; // Any parsed as a TParam -- currently not permitted
// Node(Nat){...} // fgg.FGGNode (Nat) is fgg.TParam, not fgg.TName
// type IA(type ) interface { m1() };  // m1() parsed as a TName (an invalid Spec) -- N.B. ret missing anyway

/* TODO
- WF: repeat type decl

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
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

var EVAL_TO_VAL = -1 // Must be < 0
var NO_EVAL = -2     // Must be < EVAL_TO_VAL

var verbose bool = false

func main() {
	// N.B. flags (e.g., -internal=true) must be supplied before any non-flag args
	evalPtr := flag.Int("eval", NO_EVAL,
		"-eval=n, evaluate n (>=0) steps; or -steps=-1, evaluate to value (or panic)")
	compilePtr := flag.String("compile", "",
		"-compile=\"out.go\", [WIP] monomorphise FGG source to FG (ignored if -fgg not set);"+
			" specify \"--\" to print to stdout")
	fgPtr := flag.Bool("fg", false,
		"-fg=false, interpret input as FG (defaults to true if neither -fg/-fgg set)")
	fggPtr := flag.Bool("fgg", false, "-fgg=true, interpret input as FGG")
	internalPtr := flag.Bool("internal", false,
		"-internal=true, use \"internal\" input as source")
	inlinePtr := flag.String("inline", "",
		"-inline=\"[FG/FGG src]\", use inline input as source")
	monomPtr := flag.Bool("monom", false,
		"-monom=true, [WIP] monomorphise FGG source using formal notation (ignored if -fgg not set)")
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
		var a fg.FGAdaptor
		interp(&a, src, *strictParsePtr, *evalPtr, false, "")
	} else if *fggPtr {
		var a fgg.FGGAdaptor
		interp(&a, src, *strictParsePtr, *evalPtr, *monomPtr, *compilePtr)
	}
}

// Pre: monom==true || compile != "" => -fgg is set
func interp(a base.Adaptor, src string, strict bool, steps int, monom bool,
	compile string) {
	vPrintln("\nParsing AST:")
	prog := a.Parse(strict, src) // AST (FGProgram root)
	vPrintln(prog.String())

	vPrintln("\nChecking source program OK:")
	allowStupid := false
	prog.Ok(allowStupid)

	if steps > NO_EVAL {
		eval(prog, steps)
	}

	if monom || compile != "" {
		/*var gamma fgg.ClosedEnv
		omega := make(fgg.WMap)
		fgg.MakeWMap(prog.GetDecls(), gamma, prog.GetExpr().(fgg.Expr), omega)
		for _, v := range omega {
			vPrintln(v.GetTName().String() + " |-> " + string(v.GetMonomId()))
			gs := fgg.GetParameterisedSigs(v)
			if len(gs) > 0 {
				vPrintln("Instantiations of parameterised methods: (i.e., those that had \"additional method params\")")
				for _, g := range gs {
					vPrintln("\t" + g.String())
				}
			}
		}*/
		p_mono := fgg.Monomorph(prog.(fgg.FGGProgram)) // TODO: reformat (e.g., "<...>") to make an actual FG program
		if monom {
			vPrintln("\nMonomorphising, formal notation: [Warning] WIP [Warning]")
			vPrintln(p_mono.String())
		}
		if compile != "" {
			vPrintln("\nMonomorphising, FG output: [Warning] WIP [Warning]")
			out := p_mono.String()
			out = strings.Replace(out, ",,", "", -1)
			out = strings.Replace(out, "<", "", -1)
			out = strings.Replace(out, ">", "", -1)
			if compile == "--" {
				vPrintln(out)
			} else {
				vPrintln("Writing output to: " + compile)
				d1 := []byte(out)
				err := ioutil.WriteFile(compile, d1, 0644)
				checkErr(err)
			}
		}
	}
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
	vPrintln(fmt.Sprintf("%6d: %8s %v", 0, "", p.GetExpr())) // Initial prog OK already checked

	done := steps > EVAL_TO_VAL || // Ignore 'done' if num steps fixed (set true, for ||!done below)
		p.GetExpr().IsValue() // O/w evaluate until a val -- here, check if init expr is already a val
	var rule string
	for i := 1; i <= steps || !done; i++ {
		p, rule = p.Eval()
		vPrintln(fmt.Sprintf("%6d: %8s %v", i, "["+rule+"]", p.GetExpr()))
		vPrintln("Checking OK:") // TODO: maybe disable by default, enable by flag
		p.Ok(allowStupid)
		if !done && p.GetExpr().IsValue() {
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
