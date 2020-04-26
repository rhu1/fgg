// Code generated from parser/FG.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // FG

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseFGListener is a complete listener for a parse tree produced by FGParser.
type BaseFGListener struct{}

var _ FGListener = &BaseFGListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseFGListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseFGListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseFGListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseFGListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BaseFGListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BaseFGListener) ExitProgram(ctx *ProgramContext) {}

// EnterDecls is called when production decls is entered.
func (s *BaseFGListener) EnterDecls(ctx *DeclsContext) {}

// ExitDecls is called when production decls is exited.
func (s *BaseFGListener) ExitDecls(ctx *DeclsContext) {}

// EnterTypeDecl is called when production typeDecl is entered.
func (s *BaseFGListener) EnterTypeDecl(ctx *TypeDeclContext) {}

// ExitTypeDecl is called when production typeDecl is exited.
func (s *BaseFGListener) ExitTypeDecl(ctx *TypeDeclContext) {}

// EnterMethDecl is called when production methDecl is entered.
func (s *BaseFGListener) EnterMethDecl(ctx *MethDeclContext) {}

// ExitMethDecl is called when production methDecl is exited.
func (s *BaseFGListener) ExitMethDecl(ctx *MethDeclContext) {}

// EnterStructTypeLit is called when production StructTypeLit is entered.
func (s *BaseFGListener) EnterStructTypeLit(ctx *StructTypeLitContext) {}

// ExitStructTypeLit is called when production StructTypeLit is exited.
func (s *BaseFGListener) ExitStructTypeLit(ctx *StructTypeLitContext) {}

// EnterInterfaceTypeLit is called when production InterfaceTypeLit is entered.
func (s *BaseFGListener) EnterInterfaceTypeLit(ctx *InterfaceTypeLitContext) {}

// ExitInterfaceTypeLit is called when production InterfaceTypeLit is exited.
func (s *BaseFGListener) ExitInterfaceTypeLit(ctx *InterfaceTypeLitContext) {}

// EnterFieldDecls is called when production fieldDecls is entered.
func (s *BaseFGListener) EnterFieldDecls(ctx *FieldDeclsContext) {}

// ExitFieldDecls is called when production fieldDecls is exited.
func (s *BaseFGListener) ExitFieldDecls(ctx *FieldDeclsContext) {}

// EnterFieldDecl is called when production fieldDecl is entered.
func (s *BaseFGListener) EnterFieldDecl(ctx *FieldDeclContext) {}

// ExitFieldDecl is called when production fieldDecl is exited.
func (s *BaseFGListener) ExitFieldDecl(ctx *FieldDeclContext) {}

// EnterSpecs is called when production specs is entered.
func (s *BaseFGListener) EnterSpecs(ctx *SpecsContext) {}

// ExitSpecs is called when production specs is exited.
func (s *BaseFGListener) ExitSpecs(ctx *SpecsContext) {}

// EnterSigSpec is called when production SigSpec is entered.
func (s *BaseFGListener) EnterSigSpec(ctx *SigSpecContext) {}

// ExitSigSpec is called when production SigSpec is exited.
func (s *BaseFGListener) ExitSigSpec(ctx *SigSpecContext) {}

// EnterInterfaceSpec is called when production InterfaceSpec is entered.
func (s *BaseFGListener) EnterInterfaceSpec(ctx *InterfaceSpecContext) {}

// ExitInterfaceSpec is called when production InterfaceSpec is exited.
func (s *BaseFGListener) ExitInterfaceSpec(ctx *InterfaceSpecContext) {}

// EnterSig is called when production sig is entered.
func (s *BaseFGListener) EnterSig(ctx *SigContext) {}

// ExitSig is called when production sig is exited.
func (s *BaseFGListener) ExitSig(ctx *SigContext) {}

// EnterParams is called when production params is entered.
func (s *BaseFGListener) EnterParams(ctx *ParamsContext) {}

// ExitParams is called when production params is exited.
func (s *BaseFGListener) ExitParams(ctx *ParamsContext) {}

// EnterParamDecl is called when production paramDecl is entered.
func (s *BaseFGListener) EnterParamDecl(ctx *ParamDeclContext) {}

// ExitParamDecl is called when production paramDecl is exited.
func (s *BaseFGListener) ExitParamDecl(ctx *ParamDeclContext) {}

// EnterCall is called when production Call is entered.
func (s *BaseFGListener) EnterCall(ctx *CallContext) {}

// ExitCall is called when production Call is exited.
func (s *BaseFGListener) ExitCall(ctx *CallContext) {}

// EnterVariable is called when production Variable is entered.
func (s *BaseFGListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production Variable is exited.
func (s *BaseFGListener) ExitVariable(ctx *VariableContext) {}

// EnterAssert is called when production Assert is entered.
func (s *BaseFGListener) EnterAssert(ctx *AssertContext) {}

// ExitAssert is called when production Assert is exited.
func (s *BaseFGListener) ExitAssert(ctx *AssertContext) {}

// EnterSprintf is called when production Sprintf is entered.
func (s *BaseFGListener) EnterSprintf(ctx *SprintfContext) {}

// ExitSprintf is called when production Sprintf is exited.
func (s *BaseFGListener) ExitSprintf(ctx *SprintfContext) {}

// EnterSelect is called when production Select is entered.
func (s *BaseFGListener) EnterSelect(ctx *SelectContext) {}

// ExitSelect is called when production Select is exited.
func (s *BaseFGListener) ExitSelect(ctx *SelectContext) {}

// EnterStructLit is called when production StructLit is entered.
func (s *BaseFGListener) EnterStructLit(ctx *StructLitContext) {}

// ExitStructLit is called when production StructLit is exited.
func (s *BaseFGListener) ExitStructLit(ctx *StructLitContext) {}

// EnterExprs is called when production exprs is entered.
func (s *BaseFGListener) EnterExprs(ctx *ExprsContext) {}

// ExitExprs is called when production exprs is exited.
func (s *BaseFGListener) ExitExprs(ctx *ExprsContext) {}
