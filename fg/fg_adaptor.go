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
	stack []FGNode // Because Listener methods don't return...
}

func (a *FGAdaptor) Parse(input string) FGProgram {
	is := antlr.NewInputStream(input)
	lexer := parser.NewFGLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewFGParser(stream)
	antlr.ParseTreeWalkerDefault.Walk(a, p.Program())
	return a.pop().(FGProgram)
}

func (a *FGAdaptor) ExitProgram(ctx *parser.ProgramContext) {
	body := a.pop().(Expr)
	var ds []Decl
	if ctx.GetChildCount() > 13 {
		nds := ctx.GetChild(3).GetChildCount() / 2 // e.g., decl ';' decl ';' decl ';' -- includes trailing ';'
		ds = make([]Decl, nds)
		for i := nds - 1; i >= 0; i-- {
			ds[i] = a.pop().(Decl) // Adding backwards
		}
	}
	a.push(FGProgram{ds, body})
}

// Children: 0=typ, 1=name, 2=typelit
func (a *FGAdaptor) ExitType_decl(ctx *parser.Type_declContext) {
	td := a.pop().(TStruct)
	td.t = Type(ctx.GetName().GetText())
	a.push(td)
}

func (a *FGAdaptor) ExitMeth_decl(ctx *parser.Meth_declContext) {
	m := ctx.GetMeth().GetText()
	// Reverse order
	e := a.pop().(Expr)
	t := Type(ctx.GetRet().GetText())
	var ps []ParamDecl
	if ctx.GetChildCount() > 12 {
		nps := (ctx.GetChild(6).GetChildCount() + 1) / 2 // e.g., pd ',' pd ',' pd
		ps = make([]ParamDecl, nps)
		for i := nps - 1; i >= 0; i-- {
			ps[i] = a.pop().(ParamDecl) // Adding backwards
		}
	}
	recv := a.pop().(ParamDecl)
	a.push(MDecl{recv, m, ps, t, e})
}

// Children: 2=field_decls
func (a *FGAdaptor) ExitStruct(ctx *parser.StructContext) {
	var elems []FieldDecl
	if ctx.GetChildCount() > 3 {
		nelems := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., fd ';' fd ';' fd
		elems = make([]FieldDecl, nelems)
		for i := nelems - 1; i >= 0; i-- { // N.B. ordering doesn't currently matter, stored in a map
			fd := a.pop().(FieldDecl)
			elems[i] = fd // Adding backwards
		}
	}
	a.push(TStruct{"^", elems}) // "^" to be overwritten in ExitType_decl
}

func (a *FGAdaptor) ExitField_decl(ctx *parser.Field_declContext) {
	field := ctx.GetField().GetText()
	typ := Type(ctx.GetTyp().GetText())
	a.push(FieldDecl{field, typ})
}

// Cf. ExitField_decl
func (a *FGAdaptor) ExitParamdecl(ctx *parser.ParamdeclContext) {
	x := ctx.GetVari().GetText()
	t := Type(ctx.GetTyp().GetText())
	a.push(ParamDecl{x, t})
}

func (a *FGAdaptor) EnterCall(ctx *parser.CallContext) {}

func (a *FGAdaptor) ExitCall(ctx *parser.CallContext) {}

func (a *FGAdaptor) ExitVariable(ctx *parser.VariableContext) {
	a.push(Variable{ctx.GetVariable().GetText()})
}

// Children: 0=typ (*antlr.TerminalNodeImpl), 1='{', 2=exprs (*parser.ExprsContext), 3='}'
// N.B. ExprsContext is a "helper" Context, actual exprs are its children
func (a *FGAdaptor) ExitLit(ctx *parser.LitContext) {
	typ := Type(ctx.GetTyp().GetText())
	var es []Expr
	if ctx.GetChildCount() > 3 {
		numExprs := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., 'x' ',' 'y' ',' 'z'
		es = make([]Expr, numExprs)
		for i := numExprs - 1; i >= 0; i-- {
			es[i] = a.pop().(Expr) // Adding backwards
		}
	}
	a.push(StructLit{typ, es})
}

func (a *FGAdaptor) ExitSelect(ctx *parser.SelectContext) {}

func (a *FGAdaptor) ExitAssertion(ctx *parser.AssertionContext) {}

func (a *FGAdaptor) push(i FGNode) {
	a.stack = append(a.stack, i)
}

func (a *FGAdaptor) pop() FGNode {
	if len(a.stack) < 1 {
		panic("stack is empty unable to pop")
	}
	result := a.stack[len(a.stack)-1]
	a.stack = a.stack[:len(a.stack)-1]
	return result
}
