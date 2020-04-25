package fg

import (
	"fmt"
	"reflect"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/parser/fg"
	"github.com/rhu1/fgg/parser/util"
)

var _ = fmt.Errorf
var _ = reflect.Append

// Convert ANTLR generated CST to an FGNode AST
type FGAdaptor struct {
	*parser.BaseFGListener
	stack []FGNode // Because Listener methods don't return...
}

var _ base.Adaptor = &FGAdaptor{}

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
func (a *FGAdaptor) Parse(strictParse bool, input string) base.Program {
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
		p.SetErrorHandler(&util.StrictErrorStrategy{})
	}
	antlr.ParseTreeWalkerDefault.Walk(a, p.Program())
	return a.pop().(FGProgram)
}

/* "program" */

func (a *FGAdaptor) ExitProgram(ctx *parser.ProgramContext) {
	body := a.pop().(FGExpr)
	var ds []Decl
	offset := 0 // TODO: refactor
	printf := false
	c3 := ctx.GetChild(3)
	if _, ok := c3.GetPayload().(*antlr.CommonToken); ok { // IMPORT
		//c3.GetPayload().(*antlr.CommonToken).GetText() == "import" {
		offset = 5
		printf = true
	}
	//if ctx.GetChildCount() > offset+13 {  // well-typed program must have at least one decl?
	tmp := ctx.GetChild(offset + 3)
	if _, ok := tmp.GetPayload().(*antlr.BaseParserRuleContext); ok { // If no decls, then *antlr.CommonToken, 'func'
		nds := ctx.GetChild(offset+3).GetChildCount() / 2 // (decl ';')+ -- i.e, includes trailing ';'
		ds = make([]Decl, nds)
		for i := nds - 1; i >= 0; i-- {
			ds[i] = a.pop().(Decl) // Adding backwards
		}
	}
	a.push(FGProgram{ds, body, printf})
}

/* "typeDecl" */

// Children: 1=NAME, 2=typeLit
func (a *FGAdaptor) ExitTypeDecl(ctx *parser.TypeDeclContext) {
	t := Type(ctx.GetChild(1).(*antlr.TerminalNodeImpl).GetText())
	td := a.pop().(TDecl)
	if s, ok := td.(STypeLit); ok { // N.B. s is a *copy* of td
		s.t_S = t
		a.push(s)
	} else if c, ok := td.(ITypeLit); ok {
		c.t_I = t
		a.push(c)
	} else {
		panic("Unknown type decl: " + reflect.TypeOf(td).String())
	}
}

/* #StructTypeLit ("typeLit"), "fieldDecls", "fieldDecl" */

// Children: 2=fieldDecls
func (a *FGAdaptor) ExitStructTypeLit(ctx *parser.StructTypeLitContext) {
	var fds []FieldDecl
	if ctx.GetChildCount() > 3 {
		nfds := (ctx.GetChild(2).GetChildCount() + 1) / 2 // fieldDecl (';' fieldDecl)*
		fds = make([]FieldDecl, nfds)
		for i := nfds - 1; i >= 0; i-- {
			fd := a.pop().(FieldDecl)
			fds[i] = fd // Adding backwards
		}
	}
	a.push(STypeLit{"^", fds}) // "^" to be overwritten in ExitTypeDecl
}

func (a *FGAdaptor) ExitFieldDecl(ctx *parser.FieldDeclContext) {
	f := Name(ctx.GetField().GetText())
	t := Type(ctx.GetTyp().GetText())
	a.push(FieldDecl{f, t})
}

/* "methDecl", "paramDecl" */

func (a *FGAdaptor) ExitMethDecl(ctx *parser.MethDeclContext) {
	// Reverse order
	e := a.pop().(FGExpr)
	g := a.pop().(Sig)
	recv := a.pop().(ParamDecl)
	a.push(MDecl{recv, g.meth, g.pDecls, g.t_ret, e})
}

// Cf. ExitFieldDecl
func (a *FGAdaptor) ExitParamDecl(ctx *parser.ParamDeclContext) {
	x := ctx.GetVari().GetText()
	t := Type(ctx.GetTyp().GetText())
	a.push(ParamDecl{x, t})
}

/* #InterfaceTypeLit ("typeLit"), "specs", #SigSpec ("spec"), #InterfaceSpec ("spec"), "sig" */

// Cf. ExitStructTypeLit
func (a *FGAdaptor) ExitInterfaceTypeLit(ctx *parser.InterfaceTypeLitContext) {
	var ss []Spec
	if ctx.GetChildCount() > 3 {
		nss := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., s ';' s ';' s
		ss = make([]Spec, nss)
		for i := nss - 1; i >= 0; i-- {
			s := a.pop().(Spec)
			ss[i] = s // Adding backwards
		}
	}
	a.push(ITypeLit{"^", ss}) // "^" to be overwritten in ExitTypeDecl
}

func (a *FGAdaptor) ExitSigSpec(ctx *parser.SigSpecContext) {
	// No action -- Sig is at a.stack[len(a.stack)-1]
}

func (a *FGAdaptor) ExitInterfaceSpec(ctx *parser.InterfaceSpecContext) {
	t := Type(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	a.push(t)
}

func (a *FGAdaptor) ExitSig(ctx *parser.SigContext) {
	m := ctx.GetMeth().GetText()
	// Reverse order
	t := Type(ctx.GetRet().GetText())
	var pds []ParamDecl
	if ctx.GetChildCount() > 4 {
		npds := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., pd ',' pd ',' pd
		pds = make([]ParamDecl, npds)
		for i := npds - 1; i >= 0; i-- {
			pds[i] = a.pop().(ParamDecl) // Adding backwards
		}
	}
	a.push(Sig{m, pds, t})
}

/* "expr": #Variable, #StructLit, #Select, #Call, #Assert */

func (a *FGAdaptor) ExitVariable(ctx *parser.VariableContext) {
	id := Name(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	a.push(Variable{id})
}

// Children: 0=typ (*antlr.TerminalNodeImpl), 1='{', 2=exprs (*parser.ExprsContext), 3='}'
// N.B. ExprsContext is a "helper" Context, actual exprs are its children
func (a *FGAdaptor) ExitStructLit(ctx *parser.StructLitContext) {
	t := Type(ctx.GetChild(0).(*antlr.TerminalNodeImpl).GetText())
	var es []FGExpr
	if ctx.GetChildCount() > 3 {
		nes := (ctx.GetChild(2).GetChildCount() + 1) / 2 // e.g., 'x' ',' 'y' ',' 'z'
		es = make([]FGExpr, nes)
		for i := nes - 1; i >= 0; i-- {
			es[i] = a.pop().(FGExpr) // Adding backwards
		}
	}
	a.push(StructLit{t, es})
}

func (a *FGAdaptor) ExitSelect(ctx *parser.SelectContext) {
	e := a.pop().(FGExpr)
	f := Name(ctx.GetChild(2).(*antlr.TerminalNodeImpl).GetText())
	a.push(Select{e, f})
}

func (a *FGAdaptor) ExitCall(ctx *parser.CallContext) {
	var args []FGExpr
	if ctx.GetChildCount() > 5 { // TODO: refactor as ctx.GetArgs() != nil -- and child-count-checks elsewhere
		nargs := (ctx.GetChild(4).GetChildCount() + 1) / 2 // e.g., e ',' e ',' e
		args = make([]FGExpr, nargs)
		for i := nargs - 1; i >= 0; i-- {
			args[i] = a.pop().(FGExpr) // Adding backwards
		}
	}
	m := Name(ctx.GetChild(2).(*antlr.TerminalNodeImpl).GetText())
	e := a.pop().(FGExpr)
	a.push(Call{e, m, args})
}

func (a *FGAdaptor) ExitAssert(ctx *parser.AssertContext) {
	t := Type(ctx.GetChild(3).(*antlr.TerminalNodeImpl).GetText())
	e := a.pop().(FGExpr)
	a.push(Assert{e, t})
}

func (a *FGAdaptor) ExitSprintf(ctx *parser.SprintfContext) {
	var format string = ctx.GetChild(4).(*antlr.TerminalNodeImpl).GetText()
	nargs := (ctx.GetChildCount() - 6) / 2 // Because of the comma
	args := make([]FGExpr, nargs)
	for i := 0; i < nargs; i++ {
		tmp := a.pop()
		args[i] = tmp.(FGExpr)
	}
	a.push(Sprintf{format, args})
}
