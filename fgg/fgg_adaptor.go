package fgg

import (
	"fmt"
	"reflect"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/rhu1/fgg/parser/fgg"
	"github.com/rhu1/fgg/parser/util"
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

// strictParse means panic upon any parsing error -- o/w error recovery is attempted
func (a *FGGAdaptor) Parse(strictParse bool, input string) FGGProgram {
	is := antlr.NewInputStream(input)
	var lexer antlr.Lexer
	if strictParse { // https://stackoverflow.com/questions/51683104/how-to-catch-minor-errors
		lexer = FGGBailLexer{parser.NewFGGLexer(is)} // FIXME: not working -- e.g., incr{1}, bad token
	} else {
		lexer = parser.NewFGGLexer(is)
	}
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewFGGParser(stream)
	if strictParse {
		p.RemoveErrorListeners()
		p.SetErrorHandler(&util.StrictErrorStrategy{})
	}
	antlr.ParseTreeWalkerDefault.Walk(a, p.Program())
	return a.pop().(FGGProgram)
}

/* #Typeparam ("typ"), #TypeName ("typ"), "typeFormals", "typeFDecls", "typeFDecl" */

func (a *FGGAdaptor) ExitTypeParam(ctx *parser.TypeParamContext) {
	b := TParam(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	a.push(b)
}

func (a *FGGAdaptor) ExitTypeName(ctx *parser.TypeNameContext) {
	t := Name(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	var us []Type
	if ctx.GetChildCount() > 3 { // typs "helper" Context, cf. exprs
		nus := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., u1 ',' u2 ',' u3
		us = make([]Type, nus)
		for i := nus - 1; i >= 0; i-- {
			us[i] = a.pop().(Type) // Adding backwards
		}
	}
	a.push(TName{t, us})
}

func (a *FGGAdaptor) ExitTypeFormals(ctx *parser.TypeFormalsContext) {
	var tfs []TFormal
	if ctx.GetChildCount() > 2 {
		ntfs := (ctx.GetChild(1).GetChildCount() + 1) / 2 // e.g., tf ',' tf ',' tf
		tfs = make([]TFormal, ntfs)
		for i := ntfs - 1; i >= 0; i-- {
			tfs[i] = a.pop().(TFormal) // Adding backwards
		}
	}
	a.push(TFormals{tfs})
}

func (a *FGGAdaptor) ExitTypeFDecl(ctx *parser.TypeFDeclContext) {
	u := a.pop().(Type) // CHECKME: TName? (\tau_I)
	b := a.pop().(TParam)
	a.push(TFormal{b, u})
}

/* "program" */

// Same as FG
func (a *FGGAdaptor) ExitProgram(ctx *parser.ProgramContext) {
	body := a.pop().(Expr)
	var ds []Decl
	if ctx.GetChildCount() > 13 {
		nds := ctx.GetChild(3).GetChildCount() / 2 // (decl ';')+ -- i.e, includes trailing ';'
		ds = make([]Decl, nds)
		for i := nds - 1; i >= 0; i-- {
			ds[i] = a.pop().(Decl) // Adding backwards
		}
	}
	a.push(FGGProgram{ds, body})
}

/* "typeDecl" */

// Children: 1=NAME, 2=typeFormals, 3=typeLit
func (a *FGGAdaptor) ExitTypeDecl(ctx *parser.TypeDeclContext) {
	t := Name(ctx.GetChild(1).(*antlr.TerminalNodeImpl).GetText())
	td := a.pop().(TDecl)
	psi := a.pop().(TFormals)
	if s, ok := td.(STypeLit); ok { // N.B. s is a *copy* of td
		s.t = t
		s.psi = psi
		a.push(s)
		/*} else if c, ok := td.(ITypeLit); ok {
		c.t = t
		c.psi = psi
		a.push(c)*/
	} else {
		panic("Unknown type decl: " + reflect.TypeOf(td).String())
	}
}

/* #StructTypeLit ("typeLit"), "fieldDecls", "fieldDecl" */

// Children: 2=fieldDecls
func (a *FGGAdaptor) ExitStructTypeLit(ctx *parser.StructTypeLitContext) {
	var fds []FieldDecl
	if ctx.GetChildCount() > 3 {
		nfds := (ctx.GetChild(2).GetChildCount() + 1) / 2 // fieldDecl (';' fieldDecl)*
		fds = make([]FieldDecl, nfds)
		for i := nfds - 1; i >= 0; i-- {
			fd := a.pop().(FieldDecl)
			fds[i] = fd // Adding backwards
		}
	}
	a.push(STypeLit{"^", TFormals{}, fds}) // "^" and TFormals{} to be overwritten in ExitTypeDecl
}

func (a *FGGAdaptor) ExitFieldDecl(ctx *parser.FieldDeclContext) {
	f := Name(ctx.GetField().GetText())
	//typ := Type(ctx.GetChild(1).GetText())
	u := a.pop().(Type)
	a.push(FieldDecl{f, u})
}

/* "expr": #Variable, #StructLit, #Select, #Call, #Assert */

// Same as FG
func (a *FGGAdaptor) ExitVariable(ctx *parser.VariableContext) {
	id := Name(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	a.push(Variable{id})
}

// Children: 0=typ (*antlr.TerminalNodeImpl), 1='{', 2=exprs (*parser.ExprsContext), 3='}'
func (a *FGGAdaptor) ExitStructLit(ctx *parser.StructLitContext) {
	var es []Expr
	if ctx.GetChildCount() > 3 {
		nes := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., 'x' ',' 'y' ',' 'z'
		es = make([]Expr, nes)
		for i := nes - 1; i >= 0; i-- {
			es[i] = a.pop().(Expr) // Adding backwards
		}
	}
	u := a.pop().(TName) // N.B. \tau_S, means "of the form t_S(~\tau)" (so a TName) -- i.e., not \alpha
	a.push(StructLit{u, es})
}
