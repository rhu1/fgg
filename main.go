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

func main() {
	var adptr fg.FGAdaptor // i.e., a zero-value constructor? -- stack is empty and rest only needs methods?

	ast := adptr.Parse("tS{x, y, tt{z}}")

	fmt.Println(ast)
}
