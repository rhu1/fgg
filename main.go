// Pre:
//$ go get github.com/antlr/antlr4/runtime/Go/antlr
//$ cd $CYGHOME/code/go/src/github.com/antlr/antlr4
//$ (git checkout -b antlr-go-runtime tags/4.7.1)  // Match antlr-4.7.1-complete.jar -- unnecessary

//rhu@HZHL4 MINGW64 ~/code/go/src/temp/antlr/antlr01
//$ go install
//$ /c/Users/rhu/code/go/bin/antlr01.exe

// N.B. GoInstall installs to $CYGHOME/code/go/bin (not win10-home)

package main

import (
	"fmt"
	"reflect"
	"strconv"

	"temp/antlr/antlr04/fg"
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

/* TODO
- WF: repeat type decl

	var adptr fg.FGAdaptor

	//e := "A{}"
	//e := "t_S{x, y, t_S{z}}"

	//b.WriteString("type t_S struct { };\n")
	//b.WriteString("func (x0 A) m1() A { return A{} };\n")
	//b.WriteString("type B struct { f t };\n")  // TODO: unknown type
	//b.WriteString("type B struct { b B };\n")  // TODO: recursive struct
	//b.WriteString("type t_S struct { f1 t; f2 t };\n")
*/

func main() {
	fmt.Println("Source:")
	A := "type A struct {}"
	Am1 := "func (x0 A) m1() A { return x0 }"
	Am2 := "func (x0 A) m2(x1 A) A { return x1 }"
	Am3 := "func (x0 A) m3(x1 A, x2 A) A { return x2 }"
	B := "type B struct { a A }"
	e := "B{A{}}"
	prog := fg.MakeFgProgram(A, Am1, Am2, Am3, B, e)
	fmt.Println(prog)

	fmt.Println("\nParsing AST:")
	var adptr fg.FGAdaptor
	ast := adptr.Parse(prog)
	fmt.Println(ast)

	fmt.Println("\nChecking program OK:")
	ast.Ok()
}
