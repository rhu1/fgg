// Pre (1): ANTLR4
// E.g., antlr-4.7.1-complete.jar
// (See go:generate below)

// Pre (2): ANTLR4 Runtime for Go
//$ go get github.com/antlr/antlr4/runtime/Go/antlr
// Optional:
//$ cd $CYGHOME/code/go/src/github.com/antlr/antlr4
//$ git checkout -b antlr-go-runtime tags/4.7.1  // Match antlr-4.7.1-complete.jar -- but unnecessary

//rhu@HZHL4 MINGW64 ~/code/go/src/github.com/rhu1/fgg
//$ go run . tmp/scratch.go
//$ go run . -inline="package main; type A struct {}; func main() { _ = A{} }"
// or
//$ go install
//$ /c/Users/rhu/code/go/bin/fgg.exe ...

// N.B. GoInstall installs to $CYGHOME/code/go/bin (not $WINHOME)

// Assuming "antlr4" alias for (e.g.): java -jar ~/code/java/lib/antlr-4.7.1-complete.jar
//go:generate antlr4 -Dlanguage=Go -o parser FG.g4

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"github.com/rhu1/fgg/fg"
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

// For convenient quick testing -- via flag "-internal=true"
func makeInternalSrc() string {
	Any := "type Any interface {}"
	ToAny := "type ToAny struct { any Any }"
	A := "type A struct {}"
	e := "ToAny{A{}}.any.(A)"
	return fg.MakeFgProgram(Any, ToAny, A, e)
}

// N.B. flags (e.g., -internal=true) must be supplied before any non-flag args
func main() {
	strictParsePtr := flag.Bool("strict", true,
		"Set strict parsing (panic on error, no recovery)")
	internalPtr := flag.Bool("internal", false, "Use \"internal\" input as source")
	inlinePtr := flag.String("inline", "", "Use inline input as source")
	evalPtr := flag.Int("eval", -1, "Number of steps to evaluate")
	flag.Parse()

	var src string
	if *internalPtr { // First priority
		src = makeInternalSrc()
	} else if *inlinePtr != "" { // Second priority, i.e., -inline overrules src file arg
		src = *inlinePtr
	} else {
		if len(os.Args) < 2 {
			fmt.Println("Input error: need source go file (or an -inline program)")
		}
		bs, err := ioutil.ReadFile(os.Args[1])
		checkErr(err)
		src = string(bs)
	}

	fmt.Println("\nParsing AST:")
	var adptr fg.FGAdaptor
	prog := adptr.Parse(*strictParsePtr, src) // AST (FGProgram root)
	fmt.Println(prog)

	fmt.Println("\nChecking source program OK:")
	allowStupid := false
	prog.Ok(allowStupid)

	if *evalPtr < 0 {
		return
	}

	fmt.Println("\nEntering Eval loop:")
	fmt.Println("Decls:")
	for _, v := range prog.GetDecls() {
		fmt.Println("\t" + v.String())
	}
	fmt.Println("Eval steps:")
	allowStupid = true
	fmt.Printf("%6d: %v\n", 0, prog.GetExpr())
	for i := 1; i <= *evalPtr; i++ {
		prog = prog.Eval()
		fmt.Printf("%6d: %v\n", i, prog.GetExpr())
		fmt.Println("Checking OK:")
		prog.Ok(allowStupid)
	}
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

/* TODO
- WF: repeat type decl

	//b.WriteString("type B struct { f t };\n")  // TODO: unknown type
	//b.WriteString("type B struct { b B };\n")  // TODO: recursive struct
*/
