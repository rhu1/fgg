package fgg

import (
	"fmt"
	"reflect"

	//"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/rhu1/fgg/parser/fgg"
)

var _ = fmt.Errorf
var _ = reflect.Append

// Convert ANTLR generated CST to an FGNode AST
type FGGAdaptor struct {
	*parser.BaseFGGListener
	stack []FGGNode // Because Listener methods don't return...
}

func (a *FGGAdaptor) push(n FGGNode) {
	a.stack = append(a.stack, n)
}

func (a *FGGAdaptor) pop() FGGNode {
	if len(a.stack) < 1 {
		panic("Stack is empty")
	}
	res := a.stack[len(a.stack)-1]
	a.stack = a.stack[:len(a.stack)-1]
	return res
}
