// Code generated from parser/FGG.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // FGG

import "github.com/antlr/antlr4/runtime/Go/antlr"

// FGGListener is a complete listener for a parse tree produced by FGGParser.
type FGGListener interface {
	antlr.ParseTreeListener

	// EnterTypeParam is called when entering the TypeParam production.
	EnterTypeParam(c *TypeParamContext)

	// EnterTypeName is called when entering the TypeName production.
	EnterTypeName(c *TypeNameContext)

	// EnterTyps is called when entering the typs production.
	EnterTyps(c *TypsContext)

	// EnterTypeFormals is called when entering the typeFormals production.
	EnterTypeFormals(c *TypeFormalsContext)

	// EnterTypeFDecls is called when entering the typeFDecls production.
	EnterTypeFDecls(c *TypeFDeclsContext)

	// EnterTypeFDecl is called when entering the typeFDecl production.
	EnterTypeFDecl(c *TypeFDeclContext)

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterDecls is called when entering the decls production.
	EnterDecls(c *DeclsContext)

	// EnterTypeDecl is called when entering the typeDecl production.
	EnterTypeDecl(c *TypeDeclContext)

	// EnterMethDecl is called when entering the methDecl production.
	EnterMethDecl(c *MethDeclContext)

	// EnterStructTypeLit is called when entering the StructTypeLit production.
	EnterStructTypeLit(c *StructTypeLitContext)

	// EnterInterfaceTypeLit is called when entering the InterfaceTypeLit production.
	EnterInterfaceTypeLit(c *InterfaceTypeLitContext)

	// EnterFieldDecls is called when entering the fieldDecls production.
	EnterFieldDecls(c *FieldDeclsContext)

	// EnterFieldDecl is called when entering the fieldDecl production.
	EnterFieldDecl(c *FieldDeclContext)

	// EnterSpecs is called when entering the specs production.
	EnterSpecs(c *SpecsContext)

	// EnterSigSpec is called when entering the SigSpec production.
	EnterSigSpec(c *SigSpecContext)

	// EnterInterfaceSpec is called when entering the InterfaceSpec production.
	EnterInterfaceSpec(c *InterfaceSpecContext)

	// EnterSig is called when entering the sig production.
	EnterSig(c *SigContext)

	// EnterParams is called when entering the params production.
	EnterParams(c *ParamsContext)

	// EnterParamDecl is called when entering the paramDecl production.
	EnterParamDecl(c *ParamDeclContext)

	// EnterCall is called when entering the Call production.
	EnterCall(c *CallContext)

	// EnterVariable is called when entering the Variable production.
	EnterVariable(c *VariableContext)

	// EnterAssert is called when entering the Assert production.
	EnterAssert(c *AssertContext)

	// EnterSelect is called when entering the Select production.
	EnterSelect(c *SelectContext)

	// EnterStructLit is called when entering the StructLit production.
	EnterStructLit(c *StructLitContext)

	// EnterExprs is called when entering the exprs production.
	EnterExprs(c *ExprsContext)

	// ExitTypeParam is called when exiting the TypeParam production.
	ExitTypeParam(c *TypeParamContext)

	// ExitTypeName is called when exiting the TypeName production.
	ExitTypeName(c *TypeNameContext)

	// ExitTyps is called when exiting the typs production.
	ExitTyps(c *TypsContext)

	// ExitTypeFormals is called when exiting the typeFormals production.
	ExitTypeFormals(c *TypeFormalsContext)

	// ExitTypeFDecls is called when exiting the typeFDecls production.
	ExitTypeFDecls(c *TypeFDeclsContext)

	// ExitTypeFDecl is called when exiting the typeFDecl production.
	ExitTypeFDecl(c *TypeFDeclContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitDecls is called when exiting the decls production.
	ExitDecls(c *DeclsContext)

	// ExitTypeDecl is called when exiting the typeDecl production.
	ExitTypeDecl(c *TypeDeclContext)

	// ExitMethDecl is called when exiting the methDecl production.
	ExitMethDecl(c *MethDeclContext)

	// ExitStructTypeLit is called when exiting the StructTypeLit production.
	ExitStructTypeLit(c *StructTypeLitContext)

	// ExitInterfaceTypeLit is called when exiting the InterfaceTypeLit production.
	ExitInterfaceTypeLit(c *InterfaceTypeLitContext)

	// ExitFieldDecls is called when exiting the fieldDecls production.
	ExitFieldDecls(c *FieldDeclsContext)

	// ExitFieldDecl is called when exiting the fieldDecl production.
	ExitFieldDecl(c *FieldDeclContext)

	// ExitSpecs is called when exiting the specs production.
	ExitSpecs(c *SpecsContext)

	// ExitSigSpec is called when exiting the SigSpec production.
	ExitSigSpec(c *SigSpecContext)

	// ExitInterfaceSpec is called when exiting the InterfaceSpec production.
	ExitInterfaceSpec(c *InterfaceSpecContext)

	// ExitSig is called when exiting the sig production.
	ExitSig(c *SigContext)

	// ExitParams is called when exiting the params production.
	ExitParams(c *ParamsContext)

	// ExitParamDecl is called when exiting the paramDecl production.
	ExitParamDecl(c *ParamDeclContext)

	// ExitCall is called when exiting the Call production.
	ExitCall(c *CallContext)

	// ExitVariable is called when exiting the Variable production.
	ExitVariable(c *VariableContext)

	// ExitAssert is called when exiting the Assert production.
	ExitAssert(c *AssertContext)

	// ExitSelect is called when exiting the Select production.
	ExitSelect(c *SelectContext)

	// ExitStructLit is called when exiting the StructLit production.
	ExitStructLit(c *StructLitContext)

	// ExitExprs is called when exiting the exprs production.
	ExitExprs(c *ExprsContext)
}
