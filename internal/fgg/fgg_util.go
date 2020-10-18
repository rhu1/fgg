package fgg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/rhu1/fgg/parser/fgg/parser"
)

var _ = fmt.Errorf

// Pre: len(elems) > 1
// Pre: elems[:len(elems)-1] -- type/meth decls; elems[len(elems)-1] -- "main" func body expression
func MakeFggProgram(elems ...string) string {
	if len(elems) == 0 {
		panic("Bad empty args: must supply at least body expression for \"main\"")
	}
	var b strings.Builder
	b.WriteString("package main;\n")
	for _, v := range elems[:len(elems)-1] {
		b.WriteString(v)
		b.WriteString(";\n")
	}
	b.WriteString("func main() { _ = " + elems[len(elems)-1] + " }")
	return b.String()
}

/* For "strict" parsing, *lexer* errors */

type FGGBailLexer struct {
	*parser.FGGLexer
}

// FIXME: not working -- e.g., incr{1}, bad token
// Want to "override" *BaseLexer.Recover -- XXX that's not how Go works (because BaseLexer is a struct, not interface)
func (b *FGGBailLexer) Recover(re antlr.RecognitionException) {
	message := "lex error after token " + re.GetOffendingToken().GetText() +
		" at position " + strconv.Itoa(re.GetOffendingToken().GetStart())
	panic(message)
}
