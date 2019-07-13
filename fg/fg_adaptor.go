package fg

import (
	"fmt"
	"reflect"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"temp/antlr/antlr04/parser"
)

var _ = fmt.Errorf
var _ = reflect.Append

type FGAdaptor struct {
	*parser.BaseFGListener
	stack []Expr // Because Listener methods don't return...
}

func (this *FGAdaptor) Parse(input string) Expr {
	is := antlr.NewInputStream(input)
	lexer := parser.NewFGLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewFGParser(stream)
	antlr.ParseTreeWalkerDefault.Walk(this, p.Start())
	return this.pop()
}

func (l *FGAdaptor) push(i Expr) {
	l.stack = append(l.stack, i)
}

func (l *FGAdaptor) pop() Expr {
	if len(l.stack) < 1 {
		panic("stack is empty unable to pop")
	}
	result := l.stack[len(l.stack)-1]
	l.stack = l.stack[:len(l.stack)-1]
	return result
}

func (s *FGAdaptor) ExitStart(ctx *parser.StartContext) {}

func (s *FGAdaptor) EnterCall(ctx *parser.CallContext) {}

func (s *FGAdaptor) ExitCall(ctx *parser.CallContext) {}

func (s *FGAdaptor) ExitVariable(ctx *parser.VariableContext) {
	s.push(Variable{Name(ctx.GetVariable().GetText())})
}

// children: 0=typ (*antlr.TerminalNodeImpl), 1={, 2=exprs (*parser.ExprsContext), 3=}
func (s *FGAdaptor) ExitLit(ctx *parser.LitContext) {
	typ := ctx.GetChild(0).(*antlr.TerminalNodeImpl)
	name := typ.GetText()
	numExprs := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., 'x' ',' 'y' ',' 'z'
	es := make([]Expr, numExprs)
	for i := numExprs - 1; i >= 0; i-- {
		es[i] = s.pop()
	}
	s.push(StructLit{name, es})
}

func (s *FGAdaptor) ExitSelect(ctx *parser.SelectContext) {}

func (s *FGAdaptor) ExitAssertion(ctx *parser.AssertionContext) {}
