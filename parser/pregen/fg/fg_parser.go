// Code generated from parser/FG.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // FG

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 29, 197,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2,
	3, 2, 5, 2, 38, 10, 2, 3, 2, 5, 2, 41, 10, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3,
	2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3,
	2, 5, 2, 60, 10, 2, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 5, 3, 67, 10, 3, 3, 3,
	3, 3, 6, 3, 71, 10, 3, 13, 3, 14, 3, 72, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5,
	3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6,
	5, 6, 92, 10, 6, 3, 6, 3, 6, 3, 6, 3, 6, 5, 6, 98, 10, 6, 3, 6, 5, 6, 101,
	10, 6, 3, 7, 3, 7, 3, 7, 7, 7, 106, 10, 7, 12, 7, 14, 7, 109, 11, 7, 3,
	8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 7, 9, 117, 10, 9, 12, 9, 14, 9, 120, 11,
	9, 3, 10, 3, 10, 5, 10, 124, 10, 10, 3, 11, 3, 11, 3, 11, 5, 11, 129, 10,
	11, 3, 11, 3, 11, 3, 11, 3, 12, 3, 12, 3, 12, 7, 12, 137, 10, 12, 12, 12,
	14, 12, 140, 11, 12, 3, 13, 3, 13, 3, 13, 3, 14, 3, 14, 3, 14, 3, 14, 3,
	14, 5, 14, 150, 10, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14,
	3, 14, 7, 14, 160, 10, 14, 12, 14, 14, 14, 163, 11, 14, 3, 14, 5, 14, 166,
	10, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 5, 14,
	176, 10, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 7, 14, 184, 10,
	14, 12, 14, 14, 14, 187, 11, 14, 3, 15, 3, 15, 3, 15, 7, 15, 192, 10, 15,
	12, 15, 14, 15, 195, 11, 15, 3, 15, 2, 3, 26, 16, 2, 4, 6, 8, 10, 12, 14,
	16, 18, 20, 22, 24, 26, 28, 2, 2, 2, 205, 2, 30, 3, 2, 2, 2, 4, 70, 3,
	2, 2, 2, 6, 74, 3, 2, 2, 2, 8, 78, 3, 2, 2, 2, 10, 100, 3, 2, 2, 2, 12,
	102, 3, 2, 2, 2, 14, 110, 3, 2, 2, 2, 16, 113, 3, 2, 2, 2, 18, 123, 3,
	2, 2, 2, 20, 125, 3, 2, 2, 2, 22, 133, 3, 2, 2, 2, 24, 141, 3, 2, 2, 2,
	26, 165, 3, 2, 2, 2, 28, 188, 3, 2, 2, 2, 30, 31, 7, 17, 2, 2, 31, 32,
	7, 16, 2, 2, 32, 37, 7, 3, 2, 2, 33, 34, 7, 21, 2, 2, 34, 35, 7, 4, 2,
	2, 35, 36, 7, 22, 2, 2, 36, 38, 7, 4, 2, 2, 37, 33, 3, 2, 2, 2, 37, 38,
	3, 2, 2, 2, 38, 40, 3, 2, 2, 2, 39, 41, 5, 4, 3, 2, 40, 39, 3, 2, 2, 2,
	40, 41, 3, 2, 2, 2, 41, 42, 3, 2, 2, 2, 42, 43, 7, 14, 2, 2, 43, 44, 7,
	16, 2, 2, 44, 45, 7, 5, 2, 2, 45, 46, 7, 6, 2, 2, 46, 59, 7, 7, 2, 2, 47,
	48, 7, 8, 2, 2, 48, 49, 7, 9, 2, 2, 49, 60, 5, 26, 14, 2, 50, 51, 7, 22,
	2, 2, 51, 52, 7, 10, 2, 2, 52, 53, 7, 23, 2, 2, 53, 54, 7, 5, 2, 2, 54,
	55, 7, 11, 2, 2, 55, 56, 7, 12, 2, 2, 56, 57, 5, 26, 14, 2, 57, 58, 7,
	6, 2, 2, 58, 60, 3, 2, 2, 2, 59, 47, 3, 2, 2, 2, 59, 50, 3, 2, 2, 2, 60,
	61, 3, 2, 2, 2, 61, 62, 7, 13, 2, 2, 62, 63, 7, 2, 2, 3, 63, 3, 3, 2, 2,
	2, 64, 67, 5, 6, 4, 2, 65, 67, 5, 8, 5, 2, 66, 64, 3, 2, 2, 2, 66, 65,
	3, 2, 2, 2, 67, 68, 3, 2, 2, 2, 68, 69, 7, 3, 2, 2, 69, 71, 3, 2, 2, 2,
	70, 66, 3, 2, 2, 2, 71, 72, 3, 2, 2, 2, 72, 70, 3, 2, 2, 2, 72, 73, 3,
	2, 2, 2, 73, 5, 3, 2, 2, 2, 74, 75, 7, 20, 2, 2, 75, 76, 7, 25, 2, 2, 76,
	77, 5, 10, 6, 2, 77, 7, 3, 2, 2, 2, 78, 79, 7, 14, 2, 2, 79, 80, 7, 5,
	2, 2, 80, 81, 5, 24, 13, 2, 81, 82, 7, 6, 2, 2, 82, 83, 5, 20, 11, 2, 83,
	84, 7, 7, 2, 2, 84, 85, 7, 18, 2, 2, 85, 86, 5, 26, 14, 2, 86, 87, 7, 13,
	2, 2, 87, 9, 3, 2, 2, 2, 88, 89, 7, 19, 2, 2, 89, 91, 7, 7, 2, 2, 90, 92,
	5, 12, 7, 2, 91, 90, 3, 2, 2, 2, 91, 92, 3, 2, 2, 2, 92, 93, 3, 2, 2, 2,
	93, 101, 7, 13, 2, 2, 94, 95, 7, 15, 2, 2, 95, 97, 7, 7, 2, 2, 96, 98,
	5, 16, 9, 2, 97, 96, 3, 2, 2, 2, 97, 98, 3, 2, 2, 2, 98, 99, 3, 2, 2, 2,
	99, 101, 7, 13, 2, 2, 100, 88, 3, 2, 2, 2, 100, 94, 3, 2, 2, 2, 101, 11,
	3, 2, 2, 2, 102, 107, 5, 14, 8, 2, 103, 104, 7, 3, 2, 2, 104, 106, 5, 14,
	8, 2, 105, 103, 3, 2, 2, 2, 106, 109, 3, 2, 2, 2, 107, 105, 3, 2, 2, 2,
	107, 108, 3, 2, 2, 2, 108, 13, 3, 2, 2, 2, 109, 107, 3, 2, 2, 2, 110, 111,
	7, 25, 2, 2, 111, 112, 7, 25, 2, 2, 112, 15, 3, 2, 2, 2, 113, 118, 5, 18,
	10, 2, 114, 115, 7, 3, 2, 2, 115, 117, 5, 18, 10, 2, 116, 114, 3, 2, 2,
	2, 117, 120, 3, 2, 2, 2, 118, 116, 3, 2, 2, 2, 118, 119, 3, 2, 2, 2, 119,
	17, 3, 2, 2, 2, 120, 118, 3, 2, 2, 2, 121, 124, 5, 20, 11, 2, 122, 124,
	7, 25, 2, 2, 123, 121, 3, 2, 2, 2, 123, 122, 3, 2, 2, 2, 124, 19, 3, 2,
	2, 2, 125, 126, 7, 25, 2, 2, 126, 128, 7, 5, 2, 2, 127, 129, 5, 22, 12,
	2, 128, 127, 3, 2, 2, 2, 128, 129, 3, 2, 2, 2, 129, 130, 3, 2, 2, 2, 130,
	131, 7, 6, 2, 2, 131, 132, 7, 25, 2, 2, 132, 21, 3, 2, 2, 2, 133, 138,
	5, 24, 13, 2, 134, 135, 7, 12, 2, 2, 135, 137, 5, 24, 13, 2, 136, 134,
	3, 2, 2, 2, 137, 140, 3, 2, 2, 2, 138, 136, 3, 2, 2, 2, 138, 139, 3, 2,
	2, 2, 139, 23, 3, 2, 2, 2, 140, 138, 3, 2, 2, 2, 141, 142, 7, 25, 2, 2,
	142, 143, 7, 25, 2, 2, 143, 25, 3, 2, 2, 2, 144, 145, 8, 14, 1, 2, 145,
	166, 7, 25, 2, 2, 146, 147, 7, 25, 2, 2, 147, 149, 7, 7, 2, 2, 148, 150,
	5, 28, 15, 2, 149, 148, 3, 2, 2, 2, 149, 150, 3, 2, 2, 2, 150, 151, 3,
	2, 2, 2, 151, 166, 7, 13, 2, 2, 152, 153, 7, 22, 2, 2, 153, 154, 7, 10,
	2, 2, 154, 155, 7, 24, 2, 2, 155, 156, 7, 5, 2, 2, 156, 161, 7, 29, 2,
	2, 157, 160, 7, 12, 2, 2, 158, 160, 5, 26, 14, 2, 159, 157, 3, 2, 2, 2,
	159, 158, 3, 2, 2, 2, 160, 163, 3, 2, 2, 2, 161, 159, 3, 2, 2, 2, 161,
	162, 3, 2, 2, 2, 162, 164, 3, 2, 2, 2, 163, 161, 3, 2, 2, 2, 164, 166,
	7, 6, 2, 2, 165, 144, 3, 2, 2, 2, 165, 146, 3, 2, 2, 2, 165, 152, 3, 2,
	2, 2, 166, 185, 3, 2, 2, 2, 167, 168, 12, 6, 2, 2, 168, 169, 7, 10, 2,
	2, 169, 184, 7, 25, 2, 2, 170, 171, 12, 5, 2, 2, 171, 172, 7, 10, 2, 2,
	172, 173, 7, 25, 2, 2, 173, 175, 7, 5, 2, 2, 174, 176, 5, 28, 15, 2, 175,
	174, 3, 2, 2, 2, 175, 176, 3, 2, 2, 2, 176, 177, 3, 2, 2, 2, 177, 184,
	7, 6, 2, 2, 178, 179, 12, 4, 2, 2, 179, 180, 7, 10, 2, 2, 180, 181, 7,
	5, 2, 2, 181, 182, 7, 25, 2, 2, 182, 184, 7, 6, 2, 2, 183, 167, 3, 2, 2,
	2, 183, 170, 3, 2, 2, 2, 183, 178, 3, 2, 2, 2, 184, 187, 3, 2, 2, 2, 185,
	183, 3, 2, 2, 2, 185, 186, 3, 2, 2, 2, 186, 27, 3, 2, 2, 2, 187, 185, 3,
	2, 2, 2, 188, 193, 5, 26, 14, 2, 189, 190, 7, 12, 2, 2, 190, 192, 5, 26,
	14, 2, 191, 189, 3, 2, 2, 2, 192, 195, 3, 2, 2, 2, 193, 191, 3, 2, 2, 2,
	193, 194, 3, 2, 2, 2, 194, 29, 3, 2, 2, 2, 195, 193, 3, 2, 2, 2, 23, 37,
	40, 59, 66, 72, 91, 97, 100, 107, 118, 123, 128, 138, 149, 159, 161, 165,
	175, 183, 185, 193,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "';'", "'\"'", "'('", "')'", "'{'", "'_'", "'='", "'.'", "'\"%#v\"'",
	"','", "'}'", "'func'", "'interface'", "'main'", "'package'", "'return'",
	"'struct'", "'type'", "'import'", "'fmt'", "'Printf'", "'Sprintf'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "", "FUNC", "INTERFACE", "MAIN",
	"PACKAGE", "RETURN", "STRUCT", "TYPE", "IMPORT", "FMT", "PRINTF", "SPRINTF",
	"NAME", "WHITESPACE", "COMMENT", "LINE_COMMENT", "SPRINTF_HACK",
}

var ruleNames = []string{
	"program", "decls", "typeDecl", "methDecl", "typeLit", "fieldDecls", "fieldDecl",
	"specs", "spec", "sig", "params", "paramDecl", "expr", "exprs",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type FGParser struct {
	*antlr.BaseParser
}

func NewFGParser(input antlr.TokenStream) *FGParser {
	this := new(FGParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "FG.g4"

	return this
}

// FGParser tokens.
const (
	FGParserEOF          = antlr.TokenEOF
	FGParserT__0         = 1
	FGParserT__1         = 2
	FGParserT__2         = 3
	FGParserT__3         = 4
	FGParserT__4         = 5
	FGParserT__5         = 6
	FGParserT__6         = 7
	FGParserT__7         = 8
	FGParserT__8         = 9
	FGParserT__9         = 10
	FGParserT__10        = 11
	FGParserFUNC         = 12
	FGParserINTERFACE    = 13
	FGParserMAIN         = 14
	FGParserPACKAGE      = 15
	FGParserRETURN       = 16
	FGParserSTRUCT       = 17
	FGParserTYPE         = 18
	FGParserIMPORT       = 19
	FGParserFMT          = 20
	FGParserPRINTF       = 21
	FGParserSPRINTF      = 22
	FGParserNAME         = 23
	FGParserWHITESPACE   = 24
	FGParserCOMMENT      = 25
	FGParserLINE_COMMENT = 26
	FGParserSPRINTF_HACK = 27
)

// FGParser rules.
const (
	FGParserRULE_program    = 0
	FGParserRULE_decls      = 1
	FGParserRULE_typeDecl   = 2
	FGParserRULE_methDecl   = 3
	FGParserRULE_typeLit    = 4
	FGParserRULE_fieldDecls = 5
	FGParserRULE_fieldDecl  = 6
	FGParserRULE_specs      = 7
	FGParserRULE_spec       = 8
	FGParserRULE_sig        = 9
	FGParserRULE_params     = 10
	FGParserRULE_paramDecl  = 11
	FGParserRULE_expr       = 12
	FGParserRULE_exprs      = 13
)

// IProgramContext is an interface to support dynamic dispatch.
type IProgramContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsProgramContext differentiates from other interfaces.
	IsProgramContext()
}

type ProgramContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgramContext() *ProgramContext {
	var p = new(ProgramContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_program
	return p
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) PACKAGE() antlr.TerminalNode {
	return s.GetToken(FGParserPACKAGE, 0)
}

func (s *ProgramContext) AllMAIN() []antlr.TerminalNode {
	return s.GetTokens(FGParserMAIN)
}

func (s *ProgramContext) MAIN(i int) antlr.TerminalNode {
	return s.GetToken(FGParserMAIN, i)
}

func (s *ProgramContext) FUNC() antlr.TerminalNode {
	return s.GetToken(FGParserFUNC, 0)
}

func (s *ProgramContext) EOF() antlr.TerminalNode {
	return s.GetToken(FGParserEOF, 0)
}

func (s *ProgramContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ProgramContext) AllFMT() []antlr.TerminalNode {
	return s.GetTokens(FGParserFMT)
}

func (s *ProgramContext) FMT(i int) antlr.TerminalNode {
	return s.GetToken(FGParserFMT, i)
}

func (s *ProgramContext) PRINTF() antlr.TerminalNode {
	return s.GetToken(FGParserPRINTF, 0)
}

func (s *ProgramContext) IMPORT() antlr.TerminalNode {
	return s.GetToken(FGParserIMPORT, 0)
}

func (s *ProgramContext) Decls() IDeclsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDeclsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDeclsContext)
}

func (s *ProgramContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgramContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgramContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (p *FGParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, FGParserRULE_program)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(28)
		p.Match(FGParserPACKAGE)
	}
	{
		p.SetState(29)
		p.Match(FGParserMAIN)
	}
	{
		p.SetState(30)
		p.Match(FGParserT__0)
	}
	p.SetState(35)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FGParserIMPORT {
		{
			p.SetState(31)
			p.Match(FGParserIMPORT)
		}
		{
			p.SetState(32)
			p.Match(FGParserT__1)
		}
		{
			p.SetState(33)
			p.Match(FGParserFMT)
		}
		{
			p.SetState(34)
			p.Match(FGParserT__1)
		}

	}
	p.SetState(38)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(37)
			p.Decls()
		}

	}
	{
		p.SetState(40)
		p.Match(FGParserFUNC)
	}
	{
		p.SetState(41)
		p.Match(FGParserMAIN)
	}
	{
		p.SetState(42)
		p.Match(FGParserT__2)
	}
	{
		p.SetState(43)
		p.Match(FGParserT__3)
	}
	{
		p.SetState(44)
		p.Match(FGParserT__4)
	}
	p.SetState(57)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FGParserT__5:
		{
			p.SetState(45)
			p.Match(FGParserT__5)
		}
		{
			p.SetState(46)
			p.Match(FGParserT__6)
		}
		{
			p.SetState(47)
			p.expr(0)
		}

	case FGParserFMT:
		{
			p.SetState(48)
			p.Match(FGParserFMT)
		}
		{
			p.SetState(49)
			p.Match(FGParserT__7)
		}
		{
			p.SetState(50)
			p.Match(FGParserPRINTF)
		}
		{
			p.SetState(51)
			p.Match(FGParserT__2)
		}
		{
			p.SetState(52)
			p.Match(FGParserT__8)
		}
		{
			p.SetState(53)
			p.Match(FGParserT__9)
		}
		{
			p.SetState(54)
			p.expr(0)
		}
		{
			p.SetState(55)
			p.Match(FGParserT__3)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(59)
		p.Match(FGParserT__10)
	}
	{
		p.SetState(60)
		p.Match(FGParserEOF)
	}

	return localctx
}

// IDeclsContext is an interface to support dynamic dispatch.
type IDeclsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDeclsContext differentiates from other interfaces.
	IsDeclsContext()
}

type DeclsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDeclsContext() *DeclsContext {
	var p = new(DeclsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_decls
	return p
}

func (*DeclsContext) IsDeclsContext() {}

func NewDeclsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclsContext {
	var p = new(DeclsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_decls

	return p
}

func (s *DeclsContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclsContext) AllTypeDecl() []ITypeDeclContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITypeDeclContext)(nil)).Elem())
	var tst = make([]ITypeDeclContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITypeDeclContext)
		}
	}

	return tst
}

func (s *DeclsContext) TypeDecl(i int) ITypeDeclContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeDeclContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITypeDeclContext)
}

func (s *DeclsContext) AllMethDecl() []IMethDeclContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IMethDeclContext)(nil)).Elem())
	var tst = make([]IMethDeclContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IMethDeclContext)
		}
	}

	return tst
}

func (s *DeclsContext) MethDecl(i int) IMethDeclContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMethDeclContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IMethDeclContext)
}

func (s *DeclsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeclsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterDecls(s)
	}
}

func (s *DeclsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitDecls(s)
	}
}

func (p *FGParser) Decls() (localctx IDeclsContext) {
	localctx = NewDeclsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, FGParserRULE_decls)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(68)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(64)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case FGParserTYPE:
				{
					p.SetState(62)
					p.TypeDecl()
				}

			case FGParserFUNC:
				{
					p.SetState(63)
					p.MethDecl()
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}
			{
				p.SetState(66)
				p.Match(FGParserT__0)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(70)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext())
	}

	return localctx
}

// ITypeDeclContext is an interface to support dynamic dispatch.
type ITypeDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeDeclContext differentiates from other interfaces.
	IsTypeDeclContext()
}

type TypeDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeDeclContext() *TypeDeclContext {
	var p = new(TypeDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_typeDecl
	return p
}

func (*TypeDeclContext) IsTypeDeclContext() {}

func NewTypeDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeDeclContext {
	var p = new(TypeDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_typeDecl

	return p
}

func (s *TypeDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeDeclContext) TYPE() antlr.TerminalNode {
	return s.GetToken(FGParserTYPE, 0)
}

func (s *TypeDeclContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGParserNAME, 0)
}

func (s *TypeDeclContext) TypeLit() ITypeLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeLitContext)
}

func (s *TypeDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterTypeDecl(s)
	}
}

func (s *TypeDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitTypeDecl(s)
	}
}

func (p *FGParser) TypeDecl() (localctx ITypeDeclContext) {
	localctx = NewTypeDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, FGParserRULE_typeDecl)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(72)
		p.Match(FGParserTYPE)
	}
	{
		p.SetState(73)
		p.Match(FGParserNAME)
	}
	{
		p.SetState(74)
		p.TypeLit()
	}

	return localctx
}

// IMethDeclContext is an interface to support dynamic dispatch.
type IMethDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMethDeclContext differentiates from other interfaces.
	IsMethDeclContext()
}

type MethDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMethDeclContext() *MethDeclContext {
	var p = new(MethDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_methDecl
	return p
}

func (*MethDeclContext) IsMethDeclContext() {}

func NewMethDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MethDeclContext {
	var p = new(MethDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_methDecl

	return p
}

func (s *MethDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *MethDeclContext) FUNC() antlr.TerminalNode {
	return s.GetToken(FGParserFUNC, 0)
}

func (s *MethDeclContext) ParamDecl() IParamDeclContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamDeclContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamDeclContext)
}

func (s *MethDeclContext) Sig() ISigContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISigContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISigContext)
}

func (s *MethDeclContext) RETURN() antlr.TerminalNode {
	return s.GetToken(FGParserRETURN, 0)
}

func (s *MethDeclContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MethDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MethDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterMethDecl(s)
	}
}

func (s *MethDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitMethDecl(s)
	}
}

func (p *FGParser) MethDecl() (localctx IMethDeclContext) {
	localctx = NewMethDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, FGParserRULE_methDecl)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(76)
		p.Match(FGParserFUNC)
	}
	{
		p.SetState(77)
		p.Match(FGParserT__2)
	}
	{
		p.SetState(78)
		p.ParamDecl()
	}
	{
		p.SetState(79)
		p.Match(FGParserT__3)
	}
	{
		p.SetState(80)
		p.Sig()
	}
	{
		p.SetState(81)
		p.Match(FGParserT__4)
	}
	{
		p.SetState(82)
		p.Match(FGParserRETURN)
	}
	{
		p.SetState(83)
		p.expr(0)
	}
	{
		p.SetState(84)
		p.Match(FGParserT__10)
	}

	return localctx
}

// ITypeLitContext is an interface to support dynamic dispatch.
type ITypeLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeLitContext differentiates from other interfaces.
	IsTypeLitContext()
}

type TypeLitContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeLitContext() *TypeLitContext {
	var p = new(TypeLitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_typeLit
	return p
}

func (*TypeLitContext) IsTypeLitContext() {}

func NewTypeLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeLitContext {
	var p = new(TypeLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_typeLit

	return p
}

func (s *TypeLitContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeLitContext) CopyFrom(ctx *TypeLitContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *TypeLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type StructTypeLitContext struct {
	*TypeLitContext
}

func NewStructTypeLitContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StructTypeLitContext {
	var p = new(StructTypeLitContext)

	p.TypeLitContext = NewEmptyTypeLitContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TypeLitContext))

	return p
}

func (s *StructTypeLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StructTypeLitContext) STRUCT() antlr.TerminalNode {
	return s.GetToken(FGParserSTRUCT, 0)
}

func (s *StructTypeLitContext) FieldDecls() IFieldDeclsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldDeclsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldDeclsContext)
}

func (s *StructTypeLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterStructTypeLit(s)
	}
}

func (s *StructTypeLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitStructTypeLit(s)
	}
}

type InterfaceTypeLitContext struct {
	*TypeLitContext
}

func NewInterfaceTypeLitContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InterfaceTypeLitContext {
	var p = new(InterfaceTypeLitContext)

	p.TypeLitContext = NewEmptyTypeLitContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TypeLitContext))

	return p
}

func (s *InterfaceTypeLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InterfaceTypeLitContext) INTERFACE() antlr.TerminalNode {
	return s.GetToken(FGParserINTERFACE, 0)
}

func (s *InterfaceTypeLitContext) Specs() ISpecsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISpecsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISpecsContext)
}

func (s *InterfaceTypeLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterInterfaceTypeLit(s)
	}
}

func (s *InterfaceTypeLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitInterfaceTypeLit(s)
	}
}

func (p *FGParser) TypeLit() (localctx ITypeLitContext) {
	localctx = NewTypeLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, FGParserRULE_typeLit)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(98)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FGParserSTRUCT:
		localctx = NewStructTypeLitContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(86)
			p.Match(FGParserSTRUCT)
		}
		{
			p.SetState(87)
			p.Match(FGParserT__4)
		}
		p.SetState(89)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FGParserNAME {
			{
				p.SetState(88)
				p.FieldDecls()
			}

		}
		{
			p.SetState(91)
			p.Match(FGParserT__10)
		}

	case FGParserINTERFACE:
		localctx = NewInterfaceTypeLitContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(92)
			p.Match(FGParserINTERFACE)
		}
		{
			p.SetState(93)
			p.Match(FGParserT__4)
		}
		p.SetState(95)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FGParserNAME {
			{
				p.SetState(94)
				p.Specs()
			}

		}
		{
			p.SetState(97)
			p.Match(FGParserT__10)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IFieldDeclsContext is an interface to support dynamic dispatch.
type IFieldDeclsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldDeclsContext differentiates from other interfaces.
	IsFieldDeclsContext()
}

type FieldDeclsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldDeclsContext() *FieldDeclsContext {
	var p = new(FieldDeclsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_fieldDecls
	return p
}

func (*FieldDeclsContext) IsFieldDeclsContext() {}

func NewFieldDeclsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldDeclsContext {
	var p = new(FieldDeclsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_fieldDecls

	return p
}

func (s *FieldDeclsContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldDeclsContext) AllFieldDecl() []IFieldDeclContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFieldDeclContext)(nil)).Elem())
	var tst = make([]IFieldDeclContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFieldDeclContext)
		}
	}

	return tst
}

func (s *FieldDeclsContext) FieldDecl(i int) IFieldDeclContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldDeclContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFieldDeclContext)
}

func (s *FieldDeclsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldDeclsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldDeclsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterFieldDecls(s)
	}
}

func (s *FieldDeclsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitFieldDecls(s)
	}
}

func (p *FGParser) FieldDecls() (localctx IFieldDeclsContext) {
	localctx = NewFieldDeclsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, FGParserRULE_fieldDecls)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(100)
		p.FieldDecl()
	}
	p.SetState(105)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGParserT__0 {
		{
			p.SetState(101)
			p.Match(FGParserT__0)
		}
		{
			p.SetState(102)
			p.FieldDecl()
		}

		p.SetState(107)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IFieldDeclContext is an interface to support dynamic dispatch.
type IFieldDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetField returns the field token.
	GetField() antlr.Token

	// GetTyp returns the typ token.
	GetTyp() antlr.Token

	// SetField sets the field token.
	SetField(antlr.Token)

	// SetTyp sets the typ token.
	SetTyp(antlr.Token)

	// IsFieldDeclContext differentiates from other interfaces.
	IsFieldDeclContext()
}

type FieldDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	field  antlr.Token
	typ    antlr.Token
}

func NewEmptyFieldDeclContext() *FieldDeclContext {
	var p = new(FieldDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_fieldDecl
	return p
}

func (*FieldDeclContext) IsFieldDeclContext() {}

func NewFieldDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldDeclContext {
	var p = new(FieldDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_fieldDecl

	return p
}

func (s *FieldDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldDeclContext) GetField() antlr.Token { return s.field }

func (s *FieldDeclContext) GetTyp() antlr.Token { return s.typ }

func (s *FieldDeclContext) SetField(v antlr.Token) { s.field = v }

func (s *FieldDeclContext) SetTyp(v antlr.Token) { s.typ = v }

func (s *FieldDeclContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(FGParserNAME)
}

func (s *FieldDeclContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(FGParserNAME, i)
}

func (s *FieldDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterFieldDecl(s)
	}
}

func (s *FieldDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitFieldDecl(s)
	}
}

func (p *FGParser) FieldDecl() (localctx IFieldDeclContext) {
	localctx = NewFieldDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, FGParserRULE_fieldDecl)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(108)

		var _m = p.Match(FGParserNAME)

		localctx.(*FieldDeclContext).field = _m
	}
	{
		p.SetState(109)

		var _m = p.Match(FGParserNAME)

		localctx.(*FieldDeclContext).typ = _m
	}

	return localctx
}

// ISpecsContext is an interface to support dynamic dispatch.
type ISpecsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSpecsContext differentiates from other interfaces.
	IsSpecsContext()
}

type SpecsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySpecsContext() *SpecsContext {
	var p = new(SpecsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_specs
	return p
}

func (*SpecsContext) IsSpecsContext() {}

func NewSpecsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SpecsContext {
	var p = new(SpecsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_specs

	return p
}

func (s *SpecsContext) GetParser() antlr.Parser { return s.parser }

func (s *SpecsContext) AllSpec() []ISpecContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISpecContext)(nil)).Elem())
	var tst = make([]ISpecContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISpecContext)
		}
	}

	return tst
}

func (s *SpecsContext) Spec(i int) ISpecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISpecContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISpecContext)
}

func (s *SpecsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SpecsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SpecsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterSpecs(s)
	}
}

func (s *SpecsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitSpecs(s)
	}
}

func (p *FGParser) Specs() (localctx ISpecsContext) {
	localctx = NewSpecsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, FGParserRULE_specs)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(111)
		p.Spec()
	}
	p.SetState(116)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGParserT__0 {
		{
			p.SetState(112)
			p.Match(FGParserT__0)
		}
		{
			p.SetState(113)
			p.Spec()
		}

		p.SetState(118)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISpecContext is an interface to support dynamic dispatch.
type ISpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSpecContext differentiates from other interfaces.
	IsSpecContext()
}

type SpecContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySpecContext() *SpecContext {
	var p = new(SpecContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_spec
	return p
}

func (*SpecContext) IsSpecContext() {}

func NewSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SpecContext {
	var p = new(SpecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_spec

	return p
}

func (s *SpecContext) GetParser() antlr.Parser { return s.parser }

func (s *SpecContext) CopyFrom(ctx *SpecContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *SpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type InterfaceSpecContext struct {
	*SpecContext
}

func NewInterfaceSpecContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InterfaceSpecContext {
	var p = new(InterfaceSpecContext)

	p.SpecContext = NewEmptySpecContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SpecContext))

	return p
}

func (s *InterfaceSpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InterfaceSpecContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGParserNAME, 0)
}

func (s *InterfaceSpecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterInterfaceSpec(s)
	}
}

func (s *InterfaceSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitInterfaceSpec(s)
	}
}

type SigSpecContext struct {
	*SpecContext
}

func NewSigSpecContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SigSpecContext {
	var p = new(SigSpecContext)

	p.SpecContext = NewEmptySpecContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SpecContext))

	return p
}

func (s *SigSpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SigSpecContext) Sig() ISigContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISigContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISigContext)
}

func (s *SigSpecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterSigSpec(s)
	}
}

func (s *SigSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitSigSpec(s)
	}
}

func (p *FGParser) Spec() (localctx ISpecContext) {
	localctx = NewSpecContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, FGParserRULE_spec)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(121)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSigSpecContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(119)
			p.Sig()
		}

	case 2:
		localctx = NewInterfaceSpecContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(120)
			p.Match(FGParserNAME)
		}

	}

	return localctx
}

// ISigContext is an interface to support dynamic dispatch.
type ISigContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetMeth returns the meth token.
	GetMeth() antlr.Token

	// GetRet returns the ret token.
	GetRet() antlr.Token

	// SetMeth sets the meth token.
	SetMeth(antlr.Token)

	// SetRet sets the ret token.
	SetRet(antlr.Token)

	// IsSigContext differentiates from other interfaces.
	IsSigContext()
}

type SigContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	meth   antlr.Token
	ret    antlr.Token
}

func NewEmptySigContext() *SigContext {
	var p = new(SigContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_sig
	return p
}

func (*SigContext) IsSigContext() {}

func NewSigContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SigContext {
	var p = new(SigContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_sig

	return p
}

func (s *SigContext) GetParser() antlr.Parser { return s.parser }

func (s *SigContext) GetMeth() antlr.Token { return s.meth }

func (s *SigContext) GetRet() antlr.Token { return s.ret }

func (s *SigContext) SetMeth(v antlr.Token) { s.meth = v }

func (s *SigContext) SetRet(v antlr.Token) { s.ret = v }

func (s *SigContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(FGParserNAME)
}

func (s *SigContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(FGParserNAME, i)
}

func (s *SigContext) Params() IParamsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IParamsContext)
}

func (s *SigContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SigContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SigContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterSig(s)
	}
}

func (s *SigContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitSig(s)
	}
}

func (p *FGParser) Sig() (localctx ISigContext) {
	localctx = NewSigContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, FGParserRULE_sig)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(123)

		var _m = p.Match(FGParserNAME)

		localctx.(*SigContext).meth = _m
	}
	{
		p.SetState(124)
		p.Match(FGParserT__2)
	}
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FGParserNAME {
		{
			p.SetState(125)
			p.Params()
		}

	}
	{
		p.SetState(128)
		p.Match(FGParserT__3)
	}
	{
		p.SetState(129)

		var _m = p.Match(FGParserNAME)

		localctx.(*SigContext).ret = _m
	}

	return localctx
}

// IParamsContext is an interface to support dynamic dispatch.
type IParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParamsContext differentiates from other interfaces.
	IsParamsContext()
}

type ParamsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamsContext() *ParamsContext {
	var p = new(ParamsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_params
	return p
}

func (*ParamsContext) IsParamsContext() {}

func NewParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamsContext {
	var p = new(ParamsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_params

	return p
}

func (s *ParamsContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamsContext) AllParamDecl() []IParamDeclContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IParamDeclContext)(nil)).Elem())
	var tst = make([]IParamDeclContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IParamDeclContext)
		}
	}

	return tst
}

func (s *ParamsContext) ParamDecl(i int) IParamDeclContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IParamDeclContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IParamDeclContext)
}

func (s *ParamsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterParams(s)
	}
}

func (s *ParamsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitParams(s)
	}
}

func (p *FGParser) Params() (localctx IParamsContext) {
	localctx = NewParamsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, FGParserRULE_params)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(131)
		p.ParamDecl()
	}
	p.SetState(136)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGParserT__9 {
		{
			p.SetState(132)
			p.Match(FGParserT__9)
		}
		{
			p.SetState(133)
			p.ParamDecl()
		}

		p.SetState(138)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IParamDeclContext is an interface to support dynamic dispatch.
type IParamDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVari returns the vari token.
	GetVari() antlr.Token

	// GetTyp returns the typ token.
	GetTyp() antlr.Token

	// SetVari sets the vari token.
	SetVari(antlr.Token)

	// SetTyp sets the typ token.
	SetTyp(antlr.Token)

	// IsParamDeclContext differentiates from other interfaces.
	IsParamDeclContext()
}

type ParamDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	vari   antlr.Token
	typ    antlr.Token
}

func NewEmptyParamDeclContext() *ParamDeclContext {
	var p = new(ParamDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_paramDecl
	return p
}

func (*ParamDeclContext) IsParamDeclContext() {}

func NewParamDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamDeclContext {
	var p = new(ParamDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_paramDecl

	return p
}

func (s *ParamDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamDeclContext) GetVari() antlr.Token { return s.vari }

func (s *ParamDeclContext) GetTyp() antlr.Token { return s.typ }

func (s *ParamDeclContext) SetVari(v antlr.Token) { s.vari = v }

func (s *ParamDeclContext) SetTyp(v antlr.Token) { s.typ = v }

func (s *ParamDeclContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(FGParserNAME)
}

func (s *ParamDeclContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(FGParserNAME, i)
}

func (s *ParamDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterParamDecl(s)
	}
}

func (s *ParamDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitParamDecl(s)
	}
}

func (p *FGParser) ParamDecl() (localctx IParamDeclContext) {
	localctx = NewParamDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, FGParserRULE_paramDecl)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(139)

		var _m = p.Match(FGParserNAME)

		localctx.(*ParamDeclContext).vari = _m
	}
	{
		p.SetState(140)

		var _m = p.Match(FGParserNAME)

		localctx.(*ParamDeclContext).typ = _m
	}

	return localctx
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) CopyFrom(ctx *ExprContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type CallContext struct {
	*ExprContext
	recv IExprContext
	args IExprsContext
}

func NewCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CallContext {
	var p = new(CallContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *CallContext) GetRecv() IExprContext { return s.recv }

func (s *CallContext) GetArgs() IExprsContext { return s.args }

func (s *CallContext) SetRecv(v IExprContext) { s.recv = v }

func (s *CallContext) SetArgs(v IExprsContext) { s.args = v }

func (s *CallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CallContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGParserNAME, 0)
}

func (s *CallContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *CallContext) Exprs() IExprsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprsContext)
}

func (s *CallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterCall(s)
	}
}

func (s *CallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitCall(s)
	}
}

type VariableContext struct {
	*ExprContext
}

func NewVariableContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *VariableContext {
	var p = new(VariableContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGParserNAME, 0)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitVariable(s)
	}
}

type AssertContext struct {
	*ExprContext
}

func NewAssertContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AssertContext {
	var p = new(AssertContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *AssertContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssertContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *AssertContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGParserNAME, 0)
}

func (s *AssertContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterAssert(s)
	}
}

func (s *AssertContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitAssert(s)
	}
}

type SprintfContext struct {
	*ExprContext
}

func NewSprintfContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SprintfContext {
	var p = new(SprintfContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *SprintfContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SprintfContext) FMT() antlr.TerminalNode {
	return s.GetToken(FGParserFMT, 0)
}

func (s *SprintfContext) SPRINTF() antlr.TerminalNode {
	return s.GetToken(FGParserSPRINTF, 0)
}

func (s *SprintfContext) SPRINTF_HACK() antlr.TerminalNode {
	return s.GetToken(FGParserSPRINTF_HACK, 0)
}

func (s *SprintfContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *SprintfContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SprintfContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterSprintf(s)
	}
}

func (s *SprintfContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitSprintf(s)
	}
}

type SelectContext struct {
	*ExprContext
}

func NewSelectContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectContext {
	var p = new(SelectContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *SelectContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SelectContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGParserNAME, 0)
}

func (s *SelectContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterSelect(s)
	}
}

func (s *SelectContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitSelect(s)
	}
}

type StructLitContext struct {
	*ExprContext
}

func NewStructLitContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StructLitContext {
	var p = new(StructLitContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *StructLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StructLitContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGParserNAME, 0)
}

func (s *StructLitContext) Exprs() IExprsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprsContext)
}

func (s *StructLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterStructLit(s)
	}
}

func (s *StructLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitStructLit(s)
	}
}

func (p *FGParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *FGParser) expr(_p int) (localctx IExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 24
	p.EnterRecursionRule(localctx, 24, FGParserRULE_expr, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(163)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 16, p.GetParserRuleContext()) {
	case 1:
		localctx = NewVariableContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(143)
			p.Match(FGParserNAME)
		}

	case 2:
		localctx = NewStructLitContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(144)
			p.Match(FGParserNAME)
		}
		{
			p.SetState(145)
			p.Match(FGParserT__4)
		}
		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FGParserFMT || _la == FGParserNAME {
			{
				p.SetState(146)
				p.Exprs()
			}

		}
		{
			p.SetState(149)
			p.Match(FGParserT__10)
		}

	case 3:
		localctx = NewSprintfContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(150)
			p.Match(FGParserFMT)
		}
		{
			p.SetState(151)
			p.Match(FGParserT__7)
		}
		{
			p.SetState(152)
			p.Match(FGParserSPRINTF)
		}
		{
			p.SetState(153)
			p.Match(FGParserT__2)
		}
		{
			p.SetState(154)
			p.Match(FGParserSPRINTF_HACK)
		}
		p.SetState(159)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FGParserT__9)|(1<<FGParserFMT)|(1<<FGParserNAME))) != 0 {
			p.SetState(157)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case FGParserT__9:
				{
					p.SetState(155)
					p.Match(FGParserT__9)
				}

			case FGParserFMT, FGParserNAME:
				{
					p.SetState(156)
					p.expr(0)
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(161)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(162)
			p.Match(FGParserT__3)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(183)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(181)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 18, p.GetParserRuleContext()) {
			case 1:
				localctx = NewSelectContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, FGParserRULE_expr)
				p.SetState(165)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(166)
					p.Match(FGParserT__7)
				}
				{
					p.SetState(167)
					p.Match(FGParserNAME)
				}

			case 2:
				localctx = NewCallContext(p, NewExprContext(p, _parentctx, _parentState))
				localctx.(*CallContext).recv = _prevctx

				p.PushNewRecursionContext(localctx, _startState, FGParserRULE_expr)
				p.SetState(168)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(169)
					p.Match(FGParserT__7)
				}
				{
					p.SetState(170)
					p.Match(FGParserNAME)
				}
				{
					p.SetState(171)
					p.Match(FGParserT__2)
				}
				p.SetState(173)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == FGParserFMT || _la == FGParserNAME {
					{
						p.SetState(172)

						var _x = p.Exprs()

						localctx.(*CallContext).args = _x
					}

				}
				{
					p.SetState(175)
					p.Match(FGParserT__3)
				}

			case 3:
				localctx = NewAssertContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, FGParserRULE_expr)
				p.SetState(176)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(177)
					p.Match(FGParserT__7)
				}
				{
					p.SetState(178)
					p.Match(FGParserT__2)
				}
				{
					p.SetState(179)
					p.Match(FGParserNAME)
				}
				{
					p.SetState(180)
					p.Match(FGParserT__3)
				}

			}

		}
		p.SetState(185)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())
	}

	return localctx
}

// IExprsContext is an interface to support dynamic dispatch.
type IExprsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprsContext differentiates from other interfaces.
	IsExprsContext()
}

type ExprsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprsContext() *ExprsContext {
	var p = new(ExprsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGParserRULE_exprs
	return p
}

func (*ExprsContext) IsExprsContext() {}

func NewExprsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprsContext {
	var p = new(ExprsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGParserRULE_exprs

	return p
}

func (s *ExprsContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprsContext) AllExpr() []IExprContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExprContext)(nil)).Elem())
	var tst = make([]IExprContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExprContext)
		}
	}

	return tst
}

func (s *ExprsContext) Expr(i int) IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.EnterExprs(s)
	}
}

func (s *ExprsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGListener); ok {
		listenerT.ExitExprs(s)
	}
}

func (p *FGParser) Exprs() (localctx IExprsContext) {
	localctx = NewExprsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, FGParserRULE_exprs)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(186)
		p.expr(0)
	}
	p.SetState(191)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGParserT__9 {
		{
			p.SetState(187)
			p.Match(FGParserT__9)
		}
		{
			p.SetState(188)
			p.expr(0)
		}

		p.SetState(193)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

func (p *FGParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 12:
		var t *ExprContext = nil
		if localctx != nil {
			t = localctx.(*ExprContext)
		}
		return p.Expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *FGParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
