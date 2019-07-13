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

func (a *FGAdaptor) Parse(input string) FGNode {
	is := antlr.NewInputStream(input)
	lexer := parser.NewFGLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewFGParser(stream)
	antlr.ParseTreeWalkerDefault.Walk(a, p.Program())
	return a.pop()
}

func (l *FGAdaptor) push(i FGNode) {
	l.stack = append(l.stack, i)
}

func (l *FGAdaptor) pop() FGNode {
	if len(l.stack) < 1 {
		panic("stack is empty unable to pop")
	}
	result := l.stack[len(l.stack)-1]
	l.stack = l.stack[:len(l.stack)-1]
	return result
}

// Children: 3=type_decl start (if any)
func (a *FGAdaptor) ExitProgram(ctx *parser.ProgramContext) {

	body := a.pop().(Expr)
	fmt.Println("body: ", body)

	fmt.Println("program children: ", ctx.GetChildCount())
	c := ctx.GetChild(3)
	fmt.Println(c, reflect.TypeOf(c), c.GetChildCount()) // Type_declsContext

	var numDecls int
	if ctx.GetChildCount() <= 13 {
		numDecls = 0
	} else {
		numDecls = ctx.GetChild(3).GetChildCount() / 2 // e.g., decl ';' decl ';' decl ';'
	}
	typeDecls := make([]TypeLit, numDecls)
	for i := numDecls - 1; i >= 0; i-- {
		popped := a.pop()
		fmt.Println("popped: ", popped)
		typeDecls[numDecls-i-1] = popped.(TypeLit)
	}
	fmt.Println("numtds", typeDecls)

	a.push(FGProgram{typeDecls, body})
}

// Children: 0=typ, 1=name, 2=typelit
func (a *FGAdaptor) ExitType_decl(ctx *parser.Type_declContext) {
	fmt.Println("td children", ctx.GetChildCount())
	name := ctx.GetName().GetText()
	td := a.pop().(TStruct)
	td.typ = name
	a.push(td)
	fmt.Println("pused td: ", td)
}

// Children: 2=field_decls
func (a *FGAdaptor) ExitStruct(ctx *parser.StructContext) {

	fmt.Println("struct children: ", ctx.GetChildCount())
	var numDecls int
	if ctx.GetChildCount() <= 3 {
		numDecls = 0
	} else {
		numDecls = (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., fd ';' fd ';' fd
	}
	fmt.Println("td struct children", numDecls, ctx.GetChild(2).GetChildCount())
	//fmt.Println("aaa: ", ctx.GetChild(2).GetChild(0), reflect.TypeOf(ctx.GetChild(2).GetChild(0)))
	//fmt.Println("bbb: ", ctx.GetChild(2).GetChild(1), reflect.TypeOf(ctx.GetChild(2).GetChild(1)))
	//elems := make(map[Name]Name)         // N.B. lost ordering -- fine?
	elems := make([]FieldDecl, numDecls)
	for i := numDecls - 1; i >= 0; i-- { // N.B. ordering doesn't currently matter, stored in a map
		fd := a.pop().(FieldDecl)
		elems[numDecls-i-1] = fd
	}
	a.push(TStruct{"^", elems})
	fmt.Println("pused tstruct: ", TStruct{"^", elems})
}

func (a *FGAdaptor) ExitField_decl(ctx *parser.Field_declContext) {
	field := ctx.GetField().GetText()
	typ := ctx.GetTyp().GetText()
	a.push(FieldDecl{field, typ})
}

func (a *FGAdaptor) EnterCall(ctx *parser.CallContext) {}

func (a *FGAdaptor) ExitCall(ctx *parser.CallContext) {}

func (a *FGAdaptor) ExitVariable(ctx *parser.VariableContext) {
	a.push(Variable{ctx.GetVariable().GetText()})
}

// Children: 0=typ (*antlr.TerminalNodeImpl), 1='{', 2=exprs (*parser.ExprsContext), 3='}'
// N.B. ExprsContext is a "helper" Context, actual exprs are its children
func (a *FGAdaptor) ExitLit(ctx *parser.LitContext) {
	typ := ctx.GetChild(0).(*antlr.TerminalNodeImpl)
	name := typ.GetText()
	ctx.GetArgs()

	fmt.Println("lit children: ", ctx.GetChildCount())
	var numExprs int
	if ctx.GetChildCount() <= 3 {
		numExprs = 0
	} else {
		numExprs = (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., 'x' ',' 'y' ',' 'z'
	}
	es := make([]Expr, numExprs)
	for i := numExprs - 1; i >= 0; i-- {
		es[numExprs-i-1] = a.pop().(Expr)
	}
	a.push(StructLit{name, es})
}

func (a *FGAdaptor) ExitSelect(ctx *parser.SelectContext) {}

func (a *FGAdaptor) ExitAssertion(ctx *parser.AssertionContext) {}
