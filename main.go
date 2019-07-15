// Pre:
//$ go get github.com/antlr/antlr4/runtime/Go/antlr
//$ cd $CYGHOME/code/go/src/github.com/antlr/antlr4
//$ (git checkout -b antlr-go-runtime tags/4.7.1)  // Match antlr-4.7.1-complete.jar -- unnecessary

//rhu@HZHL4 MINGW64 ~/code/go/src/temp/antlr/fgg
//$ go install
//$ /c/Users/rhu/code/go/bin/antlr01.exe
// or
//$ go run .

// N.B. GoInstall installs to $CYGHOME/code/go/bin (not $WINHOME)

//go:generate antlr4 -Dlanguage=Go -o parser FG.g4

package main

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/rhu1/fgg/fg"
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

func main() {
	fmt.Println("Source:")
	IA := "type IA interface { m0() A }"
	A := "type A struct {}"
	A1m := "func (x0 A) m1(x1 IA) A { return x0 }"
	e := "A{}.m1(A{})"
	prog := fg.MakeFgProgram(IA, A, A1m, e)
	fmt.Println(prog)

	fmt.Println("\nParsing AST:")
	var adptr fg.FGAdaptor
	strictParse := false
	ast := adptr.Parse(strictParse, prog)
	fmt.Println(ast)

	fmt.Println("\nChecking program OK:")
	ast.Ok()
}

/* TODO
- WF: repeat type decl

	//b.WriteString("type B struct { f t };\n")  // TODO: unknown type
	//b.WriteString("type B struct { b B };\n")  // TODO: recursive struct
*/
