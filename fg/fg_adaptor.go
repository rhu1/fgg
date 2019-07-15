package fg

import (
	"fmt"
	"reflect"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/rhu1/fgg/parser"
)

var _ = fmt.Errorf
var _ = reflect.Append

// Convert ANTLR generated CST to an FGNode AST
type FGAdaptor struct {
	*parser.BaseFGListener
	stack []FGNode // Because Listener methods don't return...
}

func (a *FGAdaptor) push(n FGNode) {
	a.stack = append(a.stack, n)
}

func (a *FGAdaptor) pop() FGNode {
	if len(a.stack) < 1 {
		panic("Stack is empty")
	}
	res := a.stack[len(a.stack)-1]
	a.stack = a.stack[:len(a.stack)-1]
	return res
}

// strictParse means panic upon any parsing error -- o/w error recovery is attempted
func (a *FGAdaptor) Parse(strictParse bool, input string) FGProgram {
	is := antlr.NewInputStream(input)
	var lexer antlr.Lexer
	if strictParse { // https://stackoverflow.com/questions/51683104/how-to-catch-minor-errors
		lexer = FGBailLexer{parser.NewFGLexer(is)} // FIXME: not working -- e.g., incr{1}, bad token
	} else {
		lexer = parser.NewFGLexer(is)
	}
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewFGParser(stream)
	if strictParse {
		p.RemoveErrorListeners()
		p.SetErrorHandler(&StrictErrorStrategy{})
	}
	antlr.ParseTreeWalkerDefault.Walk(a, p.Program())
	return a.pop().(FGProgram)
}

/* "program", "decls" */

func (a *FGAdaptor) ExitProgram(ctx *parser.ProgramContext) {
	body := a.pop().(Expr)
	var ds []Decl
	if ctx.GetChildCount() > 13 {
		nds := ctx.GetChild(3).GetChildCount() / 2 // (decl ';')+ -- i.e, includes trailing ';'
		ds = make([]Decl, nds)
		for i := nds - 1; i >= 0; i-- {
			ds[i] = a.pop().(Decl) // Adding backwards
		}
	}
	a.push(FGProgram{ds, body})
}

/* "typeDecl" */

// Children: 1=NAME, 2=typeLit
func (a *FGAdaptor) ExitTypeDecl(ctx *parser.TypeDeclContext) {
	typ := Type(ctx.GetChild(1).(*antlr.TerminalNodeImpl).GetText())
	td := a.pop()
	if s, ok := td.(STypeLit); ok { // N.B. s is a *copy* of td
		s.t = typ
		a.push(s)
	} else if r, ok := td.(ITypeLit); ok {
		r.t = typ
		a.push(r)
	} else {
		panic("Unknown type decl: " + reflect.TypeOf(td).String())
	}
}

/* "methDecl", "paramDecl" */

func (a *FGAdaptor) ExitMethDecl(ctx *parser.MethDeclContext) {
	// Reverse order
	e := a.pop().(Expr)
	sig := a.pop().(Sig)
	recv := a.pop().(ParamDecl)
	a.push(MDecl{recv, sig.m, sig.ps, sig.t, e})
}

// Cf. ExitFieldDecl
func (a *FGAdaptor) ExitParamDecl(ctx *parser.ParamDeclContext) {
	x := ctx.GetVari().GetText()
	t := Type(ctx.GetTyp().GetText())
	a.push(ParamDecl{x, t})
}

/* #StructTypeLit ("typeLit"), "fieldDecls", "fieldDecl" */

// Children: 2=fieldDecls
func (a *FGAdaptor) ExitStructTypeLit(ctx *parser.StructTypeLitContext) {
	var fds []FieldDecl
	if ctx.GetChildCount() > 3 {
		nfds := (ctx.GetChild(2).GetChildCount() + 1) / 2 // fieldDecl (';' fieldDecl)*
		fds = make([]FieldDecl, nfds)
		for i := nfds - 1; i >= 0; i-- { // N.B. ordering doesn't currently matter, stored in a map
			fd := a.pop().(FieldDecl)
			fds[i] = fd // Adding backwards
		}
	}
	a.push(STypeLit{"^", fds}) // "^" to be overwritten in ExitTypeDecl
}

func (a *FGAdaptor) ExitFieldDecl(ctx *parser.FieldDeclContext) {
	field := Name(ctx.GetField().GetText())
	typ := Type(ctx.GetTyp().GetText())
	a.push(FieldDecl{field, typ})
}

/* #InterfaceTypeLit ("typeLit"), "specs", #SigSpec ("spec"), #InterfaceSpec ("spec"), "sig" */

// Cf. ExitStructTypeLit
func (a *FGAdaptor) ExitInterfaceTypeLit(ctx *parser.InterfaceTypeLitContext) {
	var ss []Spec
	if ctx.GetChildCount() > 3 {
		nss := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., fd ';' fd ';' fd
		ss = make([]Spec, nss)
		for i := nss - 1; i >= 0; i-- { // N.B. ordering doesn't currently matter, stored in a map
			s := a.pop().(Spec)
			ss[i] = s // Adding backwards
		}
	}
	a.push(ITypeLit{"^", ss}) // "^" to be overwritten in ExitTypeDecl
}

func (a *FGAdaptor) ExitSigSpec(ctx *parser.SigSpecContext) {
	// No action -- Sig is at a.stack.peek()
}

func (a *FGAdaptor) ExitInterfaceSpec(ctx *parser.InterfaceSpecContext) {
	typ := Type(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	a.push(typ)
}

func (a *FGAdaptor) ExitSig(ctx *parser.SigContext) {
	m := ctx.GetMeth().GetText()
	// Reverse order
	t := Type(ctx.GetRet().GetText())
	var ps []ParamDecl
	if ctx.GetChildCount() > 4 {
		nps := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., pd ',' pd ',' pd
		ps = make([]ParamDecl, nps)
		for i := nps - 1; i >= 0; i-- {
			ps[i] = a.pop().(ParamDecl) // Adding backwards
		}
	}
	a.push(Sig{m, ps, t})
}

/* "expr": #Variable, #StructLit, #Select, #Call, #Assert */

func (a *FGAdaptor) ExitVariable(ctx *parser.VariableContext) {
	id := Name(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	a.push(Variable{id})
}

// Children: 0=typ (*antlr.TerminalNodeImpl), 1='{', 2=exprs (*parser.ExprsContext), 3='}'
// N.B. ExprsContext is a "helper" Context, actual exprs are its children
func (a *FGAdaptor) ExitStructLit(ctx *parser.StructLitContext) {
	typ := Type(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	var es []Expr
	if ctx.GetChildCount() > 3 {
		nes := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., 'x' ',' 'y' ',' 'z'
		es = make([]Expr, nes)
		for i := nes - 1; i >= 0; i-- {
			es[i] = a.pop().(Expr) // Adding backwards
		}
	}
	a.push(StructLit{typ, es})
}

func (a *FGAdaptor) ExitSelect(ctx *parser.SelectContext) {
	e := a.pop().(Expr)
	f := Name(ctx.GetChild(2).(*antlr.TerminalNodeImpl).GetText())
	a.push(Select{e, f})
}

func (a *FGAdaptor) ExitCall(ctx *parser.CallContext) {
	var args []Expr
	if ctx.GetChildCount() > 5 {
		nargs := (ctx.GetChild(4).GetChildCount() + 1) / 2 // e.g., e ',' e ',' e
		args = make([]Expr, nargs)
		for i := nargs - 1; i >= 0; i-- {
			args[i] = a.pop().(Expr) // Adding backwards
		}
	}
	m := Name(ctx.GetChild(2).(*antlr.TerminalNodeImpl).GetText())
	e := a.pop().(Expr)
	a.push(Call{e, m, args})
}

func (a *FGAdaptor) ExitAssert(ctx *parser.AssertContext) {
	typ := Type(ctx.GetChild(3).(*antlr.TerminalNodeImpl).GetText())
	e := a.pop().(Expr)
	a.push(Assert{e, typ})
}
