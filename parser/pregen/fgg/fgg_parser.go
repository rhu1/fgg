// Code generated from parser/FGG.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // FGG

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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 28, 252,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 3, 2, 3, 2, 3, 2, 3, 2, 5, 2, 45, 10, 2,
	3, 2, 5, 2, 48, 10, 2, 3, 3, 3, 3, 3, 3, 7, 3, 53, 10, 3, 12, 3, 14, 3,
	56, 11, 3, 3, 4, 3, 4, 3, 4, 5, 4, 61, 10, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3,
	5, 7, 5, 68, 10, 5, 12, 5, 14, 5, 71, 11, 5, 3, 6, 3, 6, 3, 6, 3, 7, 3,
	7, 3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 82, 10, 7, 3, 7, 5, 7, 85, 10, 7, 3, 7,
	3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7,
	3, 7, 3, 7, 3, 7, 3, 7, 5, 7, 104, 10, 7, 3, 7, 3, 7, 3, 7, 3, 8, 3, 8,
	5, 8, 111, 10, 8, 3, 8, 3, 8, 6, 8, 115, 10, 8, 13, 8, 14, 8, 116, 3, 9,
	3, 9, 3, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 10,
	3, 10, 3, 10, 3, 10, 3, 10, 3, 10, 3, 11, 3, 11, 3, 11, 5, 11, 139, 10,
	11, 3, 11, 3, 11, 3, 11, 3, 11, 5, 11, 145, 10, 11, 3, 11, 5, 11, 148,
	10, 11, 3, 12, 3, 12, 3, 12, 7, 12, 153, 10, 12, 12, 12, 14, 12, 156, 11,
	12, 3, 13, 3, 13, 3, 13, 3, 14, 3, 14, 3, 14, 7, 14, 164, 10, 14, 12, 14,
	14, 14, 167, 11, 14, 3, 15, 3, 15, 5, 15, 171, 10, 15, 3, 16, 3, 16, 3,
	16, 3, 16, 5, 16, 177, 10, 16, 3, 16, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17,
	7, 17, 185, 10, 17, 12, 17, 14, 17, 188, 11, 17, 3, 18, 3, 18, 3, 18, 3,
	19, 3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 198, 10, 19, 3, 19, 3, 19, 3, 19,
	3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 7, 19, 209, 10, 19, 12, 19, 14,
	19, 212, 11, 19, 3, 19, 5, 19, 215, 10, 19, 3, 19, 3, 19, 3, 19, 3, 19,
	3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 225, 10, 19, 3, 19, 3, 19, 3, 19, 5,
	19, 230, 10, 19, 3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 3, 19, 7, 19,
	239, 10, 19, 12, 19, 14, 19, 242, 11, 19, 3, 20, 3, 20, 3, 20, 7, 20, 247,
	10, 20, 12, 20, 14, 20, 250, 11, 20, 3, 20, 2, 3, 36, 21, 2, 4, 6, 8, 10,
	12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 2, 3, 4, 2, 11,
	11, 28, 28, 2, 261, 2, 47, 3, 2, 2, 2, 4, 49, 3, 2, 2, 2, 6, 57, 3, 2,
	2, 2, 8, 64, 3, 2, 2, 2, 10, 72, 3, 2, 2, 2, 12, 75, 3, 2, 2, 2, 14, 114,
	3, 2, 2, 2, 16, 118, 3, 2, 2, 2, 18, 123, 3, 2, 2, 2, 20, 147, 3, 2, 2,
	2, 22, 149, 3, 2, 2, 2, 24, 157, 3, 2, 2, 2, 26, 160, 3, 2, 2, 2, 28, 170,
	3, 2, 2, 2, 30, 172, 3, 2, 2, 2, 32, 181, 3, 2, 2, 2, 34, 189, 3, 2, 2,
	2, 36, 214, 3, 2, 2, 2, 38, 243, 3, 2, 2, 2, 40, 48, 7, 24, 2, 2, 41, 42,
	7, 24, 2, 2, 42, 44, 7, 3, 2, 2, 43, 45, 5, 4, 3, 2, 44, 43, 3, 2, 2, 2,
	44, 45, 3, 2, 2, 2, 45, 46, 3, 2, 2, 2, 46, 48, 7, 4, 2, 2, 47, 40, 3,
	2, 2, 2, 47, 41, 3, 2, 2, 2, 48, 3, 3, 2, 2, 2, 49, 54, 5, 2, 2, 2, 50,
	51, 7, 5, 2, 2, 51, 53, 5, 2, 2, 2, 52, 50, 3, 2, 2, 2, 53, 56, 3, 2, 2,
	2, 54, 52, 3, 2, 2, 2, 54, 55, 3, 2, 2, 2, 55, 5, 3, 2, 2, 2, 56, 54, 3,
	2, 2, 2, 57, 58, 7, 3, 2, 2, 58, 60, 7, 19, 2, 2, 59, 61, 5, 8, 5, 2, 60,
	59, 3, 2, 2, 2, 60, 61, 3, 2, 2, 2, 61, 62, 3, 2, 2, 2, 62, 63, 7, 4, 2,
	2, 63, 7, 3, 2, 2, 2, 64, 69, 5, 10, 6, 2, 65, 66, 7, 5, 2, 2, 66, 68,
	5, 10, 6, 2, 67, 65, 3, 2, 2, 2, 68, 71, 3, 2, 2, 2, 69, 67, 3, 2, 2, 2,
	69, 70, 3, 2, 2, 2, 70, 9, 3, 2, 2, 2, 71, 69, 3, 2, 2, 2, 72, 73, 7, 24,
	2, 2, 73, 74, 5, 2, 2, 2, 74, 11, 3, 2, 2, 2, 75, 76, 7, 16, 2, 2, 76,
	77, 7, 15, 2, 2, 77, 81, 7, 6, 2, 2, 78, 79, 7, 20, 2, 2, 79, 80, 7, 28,
	2, 2, 80, 82, 7, 6, 2, 2, 81, 78, 3, 2, 2, 2, 81, 82, 3, 2, 2, 2, 82, 84,
	3, 2, 2, 2, 83, 85, 5, 14, 8, 2, 84, 83, 3, 2, 2, 2, 84, 85, 3, 2, 2, 2,
	85, 86, 3, 2, 2, 2, 86, 87, 7, 13, 2, 2, 87, 88, 7, 15, 2, 2, 88, 89, 7,
	3, 2, 2, 89, 90, 7, 4, 2, 2, 90, 103, 7, 7, 2, 2, 91, 92, 7, 8, 2, 2, 92,
	93, 7, 9, 2, 2, 93, 104, 5, 36, 19, 2, 94, 95, 7, 21, 2, 2, 95, 96, 7,
	10, 2, 2, 96, 97, 7, 22, 2, 2, 97, 98, 7, 3, 2, 2, 98, 99, 7, 11, 2, 2,
	99, 100, 7, 5, 2, 2, 100, 101, 5, 36, 19, 2, 101, 102, 7, 4, 2, 2, 102,
	104, 3, 2, 2, 2, 103, 91, 3, 2, 2, 2, 103, 94, 3, 2, 2, 2, 104, 105, 3,
	2, 2, 2, 105, 106, 7, 12, 2, 2, 106, 107, 7, 2, 2, 3, 107, 13, 3, 2, 2,
	2, 108, 111, 5, 16, 9, 2, 109, 111, 5, 18, 10, 2, 110, 108, 3, 2, 2, 2,
	110, 109, 3, 2, 2, 2, 111, 112, 3, 2, 2, 2, 112, 113, 7, 6, 2, 2, 113,
	115, 3, 2, 2, 2, 114, 110, 3, 2, 2, 2, 115, 116, 3, 2, 2, 2, 116, 114,
	3, 2, 2, 2, 116, 117, 3, 2, 2, 2, 117, 15, 3, 2, 2, 2, 118, 119, 7, 19,
	2, 2, 119, 120, 7, 24, 2, 2, 120, 121, 5, 6, 4, 2, 121, 122, 5, 20, 11,
	2, 122, 17, 3, 2, 2, 2, 123, 124, 7, 13, 2, 2, 124, 125, 7, 3, 2, 2, 125,
	126, 7, 24, 2, 2, 126, 127, 7, 24, 2, 2, 127, 128, 5, 6, 4, 2, 128, 129,
	7, 4, 2, 2, 129, 130, 5, 30, 16, 2, 130, 131, 7, 7, 2, 2, 131, 132, 7,
	17, 2, 2, 132, 133, 5, 36, 19, 2, 133, 134, 7, 12, 2, 2, 134, 19, 3, 2,
	2, 2, 135, 136, 7, 18, 2, 2, 136, 138, 7, 7, 2, 2, 137, 139, 5, 22, 12,
	2, 138, 137, 3, 2, 2, 2, 138, 139, 3, 2, 2, 2, 139, 140, 3, 2, 2, 2, 140,
	148, 7, 12, 2, 2, 141, 142, 7, 14, 2, 2, 142, 144, 7, 7, 2, 2, 143, 145,
	5, 26, 14, 2, 144, 143, 3, 2, 2, 2, 144, 145, 3, 2, 2, 2, 145, 146, 3,
	2, 2, 2, 146, 148, 7, 12, 2, 2, 147, 135, 3, 2, 2, 2, 147, 141, 3, 2, 2,
	2, 148, 21, 3, 2, 2, 2, 149, 154, 5, 24, 13, 2, 150, 151, 7, 6, 2, 2, 151,
	153, 5, 24, 13, 2, 152, 150, 3, 2, 2, 2, 153, 156, 3, 2, 2, 2, 154, 152,
	3, 2, 2, 2, 154, 155, 3, 2, 2, 2, 155, 23, 3, 2, 2, 2, 156, 154, 3, 2,
	2, 2, 157, 158, 7, 24, 2, 2, 158, 159, 5, 2, 2, 2, 159, 25, 3, 2, 2, 2,
	160, 165, 5, 28, 15, 2, 161, 162, 7, 6, 2, 2, 162, 164, 5, 28, 15, 2, 163,
	161, 3, 2, 2, 2, 164, 167, 3, 2, 2, 2, 165, 163, 3, 2, 2, 2, 165, 166,
	3, 2, 2, 2, 166, 27, 3, 2, 2, 2, 167, 165, 3, 2, 2, 2, 168, 171, 5, 30,
	16, 2, 169, 171, 5, 2, 2, 2, 170, 168, 3, 2, 2, 2, 170, 169, 3, 2, 2, 2,
	171, 29, 3, 2, 2, 2, 172, 173, 7, 24, 2, 2, 173, 174, 5, 6, 4, 2, 174,
	176, 7, 3, 2, 2, 175, 177, 5, 32, 17, 2, 176, 175, 3, 2, 2, 2, 176, 177,
	3, 2, 2, 2, 177, 178, 3, 2, 2, 2, 178, 179, 7, 4, 2, 2, 179, 180, 5, 2,
	2, 2, 180, 31, 3, 2, 2, 2, 181, 186, 5, 34, 18, 2, 182, 183, 7, 5, 2, 2,
	183, 185, 5, 34, 18, 2, 184, 182, 3, 2, 2, 2, 185, 188, 3, 2, 2, 2, 186,
	184, 3, 2, 2, 2, 186, 187, 3, 2, 2, 2, 187, 33, 3, 2, 2, 2, 188, 186, 3,
	2, 2, 2, 189, 190, 7, 24, 2, 2, 190, 191, 5, 2, 2, 2, 191, 35, 3, 2, 2,
	2, 192, 193, 8, 19, 1, 2, 193, 215, 7, 24, 2, 2, 194, 195, 5, 2, 2, 2,
	195, 197, 7, 7, 2, 2, 196, 198, 5, 38, 20, 2, 197, 196, 3, 2, 2, 2, 197,
	198, 3, 2, 2, 2, 198, 199, 3, 2, 2, 2, 199, 200, 7, 12, 2, 2, 200, 215,
	3, 2, 2, 2, 201, 202, 7, 21, 2, 2, 202, 203, 7, 10, 2, 2, 203, 204, 7,
	23, 2, 2, 204, 205, 7, 3, 2, 2, 205, 210, 9, 2, 2, 2, 206, 209, 7, 5, 2,
	2, 207, 209, 5, 36, 19, 2, 208, 206, 3, 2, 2, 2, 208, 207, 3, 2, 2, 2,
	209, 212, 3, 2, 2, 2, 210, 208, 3, 2, 2, 2, 210, 211, 3, 2, 2, 2, 211,
	213, 3, 2, 2, 2, 212, 210, 3, 2, 2, 2, 213, 215, 7, 4, 2, 2, 214, 192,
	3, 2, 2, 2, 214, 194, 3, 2, 2, 2, 214, 201, 3, 2, 2, 2, 215, 240, 3, 2,
	2, 2, 216, 217, 12, 6, 2, 2, 217, 218, 7, 10, 2, 2, 218, 239, 7, 24, 2,
	2, 219, 220, 12, 5, 2, 2, 220, 221, 7, 10, 2, 2, 221, 222, 7, 24, 2, 2,
	222, 224, 7, 3, 2, 2, 223, 225, 5, 4, 3, 2, 224, 223, 3, 2, 2, 2, 224,
	225, 3, 2, 2, 2, 225, 226, 3, 2, 2, 2, 226, 227, 7, 4, 2, 2, 227, 229,
	7, 3, 2, 2, 228, 230, 5, 38, 20, 2, 229, 228, 3, 2, 2, 2, 229, 230, 3,
	2, 2, 2, 230, 231, 3, 2, 2, 2, 231, 239, 7, 4, 2, 2, 232, 233, 12, 4, 2,
	2, 233, 234, 7, 10, 2, 2, 234, 235, 7, 3, 2, 2, 235, 236, 5, 2, 2, 2, 236,
	237, 7, 4, 2, 2, 237, 239, 3, 2, 2, 2, 238, 216, 3, 2, 2, 2, 238, 219,
	3, 2, 2, 2, 238, 232, 3, 2, 2, 2, 239, 242, 3, 2, 2, 2, 240, 238, 3, 2,
	2, 2, 240, 241, 3, 2, 2, 2, 241, 37, 3, 2, 2, 2, 242, 240, 3, 2, 2, 2,
	243, 248, 5, 36, 19, 2, 244, 245, 7, 5, 2, 2, 245, 247, 5, 36, 19, 2, 246,
	244, 3, 2, 2, 2, 247, 250, 3, 2, 2, 2, 248, 246, 3, 2, 2, 2, 248, 249,
	3, 2, 2, 2, 249, 39, 3, 2, 2, 2, 250, 248, 3, 2, 2, 2, 29, 44, 47, 54,
	60, 69, 81, 84, 103, 110, 116, 138, 144, 147, 154, 165, 170, 176, 186,
	197, 208, 210, 214, 224, 229, 238, 240, 248,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'('", "')'", "','", "';'", "'{'", "'_'", "'='", "'.'", "'\"%#v\"'",
	"'}'", "'func'", "'interface'", "'main'", "'package'", "'return'", "'struct'",
	"'type'", "'import'", "'fmt'", "'Printf'", "'Sprintf'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "FUNC", "INTERFACE", "MAIN",
	"PACKAGE", "RETURN", "STRUCT", "TYPE", "IMPORT", "FMT", "PRINTF", "SPRINTF",
	"NAME", "WHITESPACE", "COMMENT", "LINE_COMMENT", "STRING",
}

var ruleNames = []string{
	"typ", "typs", "typeFormals", "typeFDecls", "typeFDecl", "program", "decls",
	"typeDecl", "methDecl", "typeLit", "fieldDecls", "fieldDecl", "specs",
	"spec", "sig", "params", "paramDecl", "expr", "exprs",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type FGGParser struct {
	*antlr.BaseParser
}

func NewFGGParser(input antlr.TokenStream) *FGGParser {
	this := new(FGGParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "FGG.g4"

	return this
}

// FGGParser tokens.
const (
	FGGParserEOF          = antlr.TokenEOF
	FGGParserT__0         = 1
	FGGParserT__1         = 2
	FGGParserT__2         = 3
	FGGParserT__3         = 4
	FGGParserT__4         = 5
	FGGParserT__5         = 6
	FGGParserT__6         = 7
	FGGParserT__7         = 8
	FGGParserT__8         = 9
	FGGParserT__9         = 10
	FGGParserFUNC         = 11
	FGGParserINTERFACE    = 12
	FGGParserMAIN         = 13
	FGGParserPACKAGE      = 14
	FGGParserRETURN       = 15
	FGGParserSTRUCT       = 16
	FGGParserTYPE         = 17
	FGGParserIMPORT       = 18
	FGGParserFMT          = 19
	FGGParserPRINTF       = 20
	FGGParserSPRINTF      = 21
	FGGParserNAME         = 22
	FGGParserWHITESPACE   = 23
	FGGParserCOMMENT      = 24
	FGGParserLINE_COMMENT = 25
	FGGParserSTRING       = 26
)

// FGGParser rules.
const (
	FGGParserRULE_typ         = 0
	FGGParserRULE_typs        = 1
	FGGParserRULE_typeFormals = 2
	FGGParserRULE_typeFDecls  = 3
	FGGParserRULE_typeFDecl   = 4
	FGGParserRULE_program     = 5
	FGGParserRULE_decls       = 6
	FGGParserRULE_typeDecl    = 7
	FGGParserRULE_methDecl    = 8
	FGGParserRULE_typeLit     = 9
	FGGParserRULE_fieldDecls  = 10
	FGGParserRULE_fieldDecl   = 11
	FGGParserRULE_specs       = 12
	FGGParserRULE_spec        = 13
	FGGParserRULE_sig         = 14
	FGGParserRULE_params      = 15
	FGGParserRULE_paramDecl   = 16
	FGGParserRULE_expr        = 17
	FGGParserRULE_exprs       = 18
)

// ITypContext is an interface to support dynamic dispatch.
type ITypContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypContext differentiates from other interfaces.
	IsTypContext()
}

type TypContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypContext() *TypContext {
	var p = new(TypContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_typ
	return p
}

func (*TypContext) IsTypContext() {}

func NewTypContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypContext {
	var p = new(TypContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_typ

	return p
}

func (s *TypContext) GetParser() antlr.Parser { return s.parser }

func (s *TypContext) CopyFrom(ctx *TypContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *TypContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type TypeNameContext struct {
	*TypContext
}

func NewTypeNameContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TypeNameContext {
	var p = new(TypeNameContext)

	p.TypContext = NewEmptyTypContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TypContext))

	return p
}

func (s *TypeNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeNameContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, 0)
}

func (s *TypeNameContext) Typs() ITypsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypsContext)
}

func (s *TypeNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterTypeName(s)
	}
}

func (s *TypeNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitTypeName(s)
	}
}

type TypeParamContext struct {
	*TypContext
}

func NewTypeParamContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TypeParamContext {
	var p = new(TypeParamContext)

	p.TypContext = NewEmptyTypContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TypContext))

	return p
}

func (s *TypeParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeParamContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, 0)
}

func (s *TypeParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterTypeParam(s)
	}
}

func (s *TypeParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitTypeParam(s)
	}
}

func (p *FGGParser) Typ() (localctx ITypContext) {
	localctx = NewTypContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, FGGParserRULE_typ)
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

	p.SetState(45)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		localctx = NewTypeParamContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(38)
			p.Match(FGGParserNAME)
		}

	case 2:
		localctx = NewTypeNameContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(39)
			p.Match(FGGParserNAME)
		}
		{
			p.SetState(40)
			p.Match(FGGParserT__0)
		}
		p.SetState(42)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FGGParserNAME {
			{
				p.SetState(41)
				p.Typs()
			}

		}
		{
			p.SetState(44)
			p.Match(FGGParserT__1)
		}

	}

	return localctx
}

// ITypsContext is an interface to support dynamic dispatch.
type ITypsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypsContext differentiates from other interfaces.
	IsTypsContext()
}

type TypsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypsContext() *TypsContext {
	var p = new(TypsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_typs
	return p
}

func (*TypsContext) IsTypsContext() {}

func NewTypsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypsContext {
	var p = new(TypsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_typs

	return p
}

func (s *TypsContext) GetParser() antlr.Parser { return s.parser }

func (s *TypsContext) AllTyp() []ITypContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITypContext)(nil)).Elem())
	var tst = make([]ITypContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITypContext)
		}
	}

	return tst
}

func (s *TypsContext) Typ(i int) ITypContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITypContext)
}

func (s *TypsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterTyps(s)
	}
}

func (s *TypsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitTyps(s)
	}
}

func (p *FGGParser) Typs() (localctx ITypsContext) {
	localctx = NewTypsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, FGGParserRULE_typs)
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
		p.SetState(47)
		p.Typ()
	}
	p.SetState(52)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGGParserT__2 {
		{
			p.SetState(48)
			p.Match(FGGParserT__2)
		}
		{
			p.SetState(49)
			p.Typ()
		}

		p.SetState(54)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ITypeFormalsContext is an interface to support dynamic dispatch.
type ITypeFormalsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeFormalsContext differentiates from other interfaces.
	IsTypeFormalsContext()
}

type TypeFormalsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeFormalsContext() *TypeFormalsContext {
	var p = new(TypeFormalsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_typeFormals
	return p
}

func (*TypeFormalsContext) IsTypeFormalsContext() {}

func NewTypeFormalsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeFormalsContext {
	var p = new(TypeFormalsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_typeFormals

	return p
}

func (s *TypeFormalsContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeFormalsContext) TYPE() antlr.TerminalNode {
	return s.GetToken(FGGParserTYPE, 0)
}

func (s *TypeFormalsContext) TypeFDecls() ITypeFDeclsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeFDeclsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeFDeclsContext)
}

func (s *TypeFormalsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeFormalsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeFormalsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterTypeFormals(s)
	}
}

func (s *TypeFormalsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitTypeFormals(s)
	}
}

func (p *FGGParser) TypeFormals() (localctx ITypeFormalsContext) {
	localctx = NewTypeFormalsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, FGGParserRULE_typeFormals)
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
		p.SetState(55)
		p.Match(FGGParserT__0)
	}
	{
		p.SetState(56)
		p.Match(FGGParserTYPE)
	}
	p.SetState(58)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FGGParserNAME {
		{
			p.SetState(57)
			p.TypeFDecls()
		}

	}
	{
		p.SetState(60)
		p.Match(FGGParserT__1)
	}

	return localctx
}

// ITypeFDeclsContext is an interface to support dynamic dispatch.
type ITypeFDeclsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeFDeclsContext differentiates from other interfaces.
	IsTypeFDeclsContext()
}

type TypeFDeclsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeFDeclsContext() *TypeFDeclsContext {
	var p = new(TypeFDeclsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_typeFDecls
	return p
}

func (*TypeFDeclsContext) IsTypeFDeclsContext() {}

func NewTypeFDeclsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeFDeclsContext {
	var p = new(TypeFDeclsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_typeFDecls

	return p
}

func (s *TypeFDeclsContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeFDeclsContext) AllTypeFDecl() []ITypeFDeclContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITypeFDeclContext)(nil)).Elem())
	var tst = make([]ITypeFDeclContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITypeFDeclContext)
		}
	}

	return tst
}

func (s *TypeFDeclsContext) TypeFDecl(i int) ITypeFDeclContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeFDeclContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITypeFDeclContext)
}

func (s *TypeFDeclsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeFDeclsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeFDeclsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterTypeFDecls(s)
	}
}

func (s *TypeFDeclsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitTypeFDecls(s)
	}
}

func (p *FGGParser) TypeFDecls() (localctx ITypeFDeclsContext) {
	localctx = NewTypeFDeclsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, FGGParserRULE_typeFDecls)
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
		p.SetState(62)
		p.TypeFDecl()
	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGGParserT__2 {
		{
			p.SetState(63)
			p.Match(FGGParserT__2)
		}
		{
			p.SetState(64)
			p.TypeFDecl()
		}

		p.SetState(69)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ITypeFDeclContext is an interface to support dynamic dispatch.
type ITypeFDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeFDeclContext differentiates from other interfaces.
	IsTypeFDeclContext()
}

type TypeFDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeFDeclContext() *TypeFDeclContext {
	var p = new(TypeFDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_typeFDecl
	return p
}

func (*TypeFDeclContext) IsTypeFDeclContext() {}

func NewTypeFDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeFDeclContext {
	var p = new(TypeFDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_typeFDecl

	return p
}

func (s *TypeFDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeFDeclContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, 0)
}

func (s *TypeFDeclContext) Typ() ITypContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypContext)
}

func (s *TypeFDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeFDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeFDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterTypeFDecl(s)
	}
}

func (s *TypeFDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitTypeFDecl(s)
	}
}

func (p *FGGParser) TypeFDecl() (localctx ITypeFDeclContext) {
	localctx = NewTypeFDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, FGGParserRULE_typeFDecl)

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
		p.SetState(70)
		p.Match(FGGParserNAME)
	}
	{
		p.SetState(71)
		p.Typ()
	}

	return localctx
}

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
	p.RuleIndex = FGGParserRULE_program
	return p
}

func (*ProgramContext) IsProgramContext() {}

func NewProgramContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgramContext {
	var p = new(ProgramContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_program

	return p
}

func (s *ProgramContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgramContext) PACKAGE() antlr.TerminalNode {
	return s.GetToken(FGGParserPACKAGE, 0)
}

func (s *ProgramContext) AllMAIN() []antlr.TerminalNode {
	return s.GetTokens(FGGParserMAIN)
}

func (s *ProgramContext) MAIN(i int) antlr.TerminalNode {
	return s.GetToken(FGGParserMAIN, i)
}

func (s *ProgramContext) FUNC() antlr.TerminalNode {
	return s.GetToken(FGGParserFUNC, 0)
}

func (s *ProgramContext) EOF() antlr.TerminalNode {
	return s.GetToken(FGGParserEOF, 0)
}

func (s *ProgramContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ProgramContext) FMT() antlr.TerminalNode {
	return s.GetToken(FGGParserFMT, 0)
}

func (s *ProgramContext) PRINTF() antlr.TerminalNode {
	return s.GetToken(FGGParserPRINTF, 0)
}

func (s *ProgramContext) IMPORT() antlr.TerminalNode {
	return s.GetToken(FGGParserIMPORT, 0)
}

func (s *ProgramContext) STRING() antlr.TerminalNode {
	return s.GetToken(FGGParserSTRING, 0)
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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterProgram(s)
	}
}

func (s *ProgramContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitProgram(s)
	}
}

func (p *FGGParser) Program() (localctx IProgramContext) {
	localctx = NewProgramContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, FGGParserRULE_program)
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
		p.SetState(73)
		p.Match(FGGParserPACKAGE)
	}
	{
		p.SetState(74)
		p.Match(FGGParserMAIN)
	}
	{
		p.SetState(75)
		p.Match(FGGParserT__3)
	}
	p.SetState(79)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FGGParserIMPORT {
		{
			p.SetState(76)
			p.Match(FGGParserIMPORT)
		}
		{
			p.SetState(77)
			p.Match(FGGParserSTRING)
		}
		{
			p.SetState(78)
			p.Match(FGGParserT__3)
		}

	}
	p.SetState(82)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(81)
			p.Decls()
		}

	}
	{
		p.SetState(84)
		p.Match(FGGParserFUNC)
	}
	{
		p.SetState(85)
		p.Match(FGGParserMAIN)
	}
	{
		p.SetState(86)
		p.Match(FGGParserT__0)
	}
	{
		p.SetState(87)
		p.Match(FGGParserT__1)
	}
	{
		p.SetState(88)
		p.Match(FGGParserT__4)
	}
	p.SetState(101)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FGGParserT__5:
		{
			p.SetState(89)
			p.Match(FGGParserT__5)
		}
		{
			p.SetState(90)
			p.Match(FGGParserT__6)
		}
		{
			p.SetState(91)
			p.expr(0)
		}

	case FGGParserFMT:
		{
			p.SetState(92)
			p.Match(FGGParserFMT)
		}
		{
			p.SetState(93)
			p.Match(FGGParserT__7)
		}
		{
			p.SetState(94)
			p.Match(FGGParserPRINTF)
		}
		{
			p.SetState(95)
			p.Match(FGGParserT__0)
		}
		{
			p.SetState(96)
			p.Match(FGGParserT__8)
		}
		{
			p.SetState(97)
			p.Match(FGGParserT__2)
		}
		{
			p.SetState(98)
			p.expr(0)
		}
		{
			p.SetState(99)
			p.Match(FGGParserT__1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(103)
		p.Match(FGGParserT__9)
	}
	{
		p.SetState(104)
		p.Match(FGGParserEOF)
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
	p.RuleIndex = FGGParserRULE_decls
	return p
}

func (*DeclsContext) IsDeclsContext() {}

func NewDeclsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclsContext {
	var p = new(DeclsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_decls

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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterDecls(s)
	}
}

func (s *DeclsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitDecls(s)
	}
}

func (p *FGGParser) Decls() (localctx IDeclsContext) {
	localctx = NewDeclsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, FGGParserRULE_decls)

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
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			p.SetState(108)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case FGGParserTYPE:
				{
					p.SetState(106)
					p.TypeDecl()
				}

			case FGGParserFUNC:
				{
					p.SetState(107)
					p.MethDecl()
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}
			{
				p.SetState(110)
				p.Match(FGGParserT__3)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(114)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
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
	p.RuleIndex = FGGParserRULE_typeDecl
	return p
}

func (*TypeDeclContext) IsTypeDeclContext() {}

func NewTypeDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeDeclContext {
	var p = new(TypeDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_typeDecl

	return p
}

func (s *TypeDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeDeclContext) TYPE() antlr.TerminalNode {
	return s.GetToken(FGGParserTYPE, 0)
}

func (s *TypeDeclContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, 0)
}

func (s *TypeDeclContext) TypeFormals() ITypeFormalsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeFormalsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeFormalsContext)
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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterTypeDecl(s)
	}
}

func (s *TypeDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitTypeDecl(s)
	}
}

func (p *FGGParser) TypeDecl() (localctx ITypeDeclContext) {
	localctx = NewTypeDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, FGGParserRULE_typeDecl)

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
		p.SetState(116)
		p.Match(FGGParserTYPE)
	}
	{
		p.SetState(117)
		p.Match(FGGParserNAME)
	}
	{
		p.SetState(118)
		p.TypeFormals()
	}
	{
		p.SetState(119)
		p.TypeLit()
	}

	return localctx
}

// IMethDeclContext is an interface to support dynamic dispatch.
type IMethDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetRecv returns the recv token.
	GetRecv() antlr.Token

	// GetTypn returns the typn token.
	GetTypn() antlr.Token

	// SetRecv sets the recv token.
	SetRecv(antlr.Token)

	// SetTypn sets the typn token.
	SetTypn(antlr.Token)

	// IsMethDeclContext differentiates from other interfaces.
	IsMethDeclContext()
}

type MethDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	recv   antlr.Token
	typn   antlr.Token
}

func NewEmptyMethDeclContext() *MethDeclContext {
	var p = new(MethDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_methDecl
	return p
}

func (*MethDeclContext) IsMethDeclContext() {}

func NewMethDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MethDeclContext {
	var p = new(MethDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_methDecl

	return p
}

func (s *MethDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *MethDeclContext) GetRecv() antlr.Token { return s.recv }

func (s *MethDeclContext) GetTypn() antlr.Token { return s.typn }

func (s *MethDeclContext) SetRecv(v antlr.Token) { s.recv = v }

func (s *MethDeclContext) SetTypn(v antlr.Token) { s.typn = v }

func (s *MethDeclContext) FUNC() antlr.TerminalNode {
	return s.GetToken(FGGParserFUNC, 0)
}

func (s *MethDeclContext) TypeFormals() ITypeFormalsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeFormalsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeFormalsContext)
}

func (s *MethDeclContext) Sig() ISigContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISigContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISigContext)
}

func (s *MethDeclContext) RETURN() antlr.TerminalNode {
	return s.GetToken(FGGParserRETURN, 0)
}

func (s *MethDeclContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MethDeclContext) AllNAME() []antlr.TerminalNode {
	return s.GetTokens(FGGParserNAME)
}

func (s *MethDeclContext) NAME(i int) antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, i)
}

func (s *MethDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MethDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterMethDecl(s)
	}
}

func (s *MethDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitMethDecl(s)
	}
}

func (p *FGGParser) MethDecl() (localctx IMethDeclContext) {
	localctx = NewMethDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, FGGParserRULE_methDecl)

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
		p.SetState(121)
		p.Match(FGGParserFUNC)
	}
	{
		p.SetState(122)
		p.Match(FGGParserT__0)
	}
	{
		p.SetState(123)

		var _m = p.Match(FGGParserNAME)

		localctx.(*MethDeclContext).recv = _m
	}
	{
		p.SetState(124)

		var _m = p.Match(FGGParserNAME)

		localctx.(*MethDeclContext).typn = _m
	}
	{
		p.SetState(125)
		p.TypeFormals()
	}
	{
		p.SetState(126)
		p.Match(FGGParserT__1)
	}
	{
		p.SetState(127)
		p.Sig()
	}
	{
		p.SetState(128)
		p.Match(FGGParserT__4)
	}
	{
		p.SetState(129)
		p.Match(FGGParserRETURN)
	}
	{
		p.SetState(130)
		p.expr(0)
	}
	{
		p.SetState(131)
		p.Match(FGGParserT__9)
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
	p.RuleIndex = FGGParserRULE_typeLit
	return p
}

func (*TypeLitContext) IsTypeLitContext() {}

func NewTypeLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeLitContext {
	var p = new(TypeLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_typeLit

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
	return s.GetToken(FGGParserSTRUCT, 0)
}

func (s *StructTypeLitContext) FieldDecls() IFieldDeclsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFieldDeclsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFieldDeclsContext)
}

func (s *StructTypeLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterStructTypeLit(s)
	}
}

func (s *StructTypeLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
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
	return s.GetToken(FGGParserINTERFACE, 0)
}

func (s *InterfaceTypeLitContext) Specs() ISpecsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISpecsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISpecsContext)
}

func (s *InterfaceTypeLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterInterfaceTypeLit(s)
	}
}

func (s *InterfaceTypeLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitInterfaceTypeLit(s)
	}
}

func (p *FGGParser) TypeLit() (localctx ITypeLitContext) {
	localctx = NewTypeLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, FGGParserRULE_typeLit)
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

	p.SetState(145)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case FGGParserSTRUCT:
		localctx = NewStructTypeLitContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(133)
			p.Match(FGGParserSTRUCT)
		}
		{
			p.SetState(134)
			p.Match(FGGParserT__4)
		}
		p.SetState(136)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FGGParserNAME {
			{
				p.SetState(135)
				p.FieldDecls()
			}

		}
		{
			p.SetState(138)
			p.Match(FGGParserT__9)
		}

	case FGGParserINTERFACE:
		localctx = NewInterfaceTypeLitContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(139)
			p.Match(FGGParserINTERFACE)
		}
		{
			p.SetState(140)
			p.Match(FGGParserT__4)
		}
		p.SetState(142)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FGGParserNAME {
			{
				p.SetState(141)
				p.Specs()
			}

		}
		{
			p.SetState(144)
			p.Match(FGGParserT__9)
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
	p.RuleIndex = FGGParserRULE_fieldDecls
	return p
}

func (*FieldDeclsContext) IsFieldDeclsContext() {}

func NewFieldDeclsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldDeclsContext {
	var p = new(FieldDeclsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_fieldDecls

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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterFieldDecls(s)
	}
}

func (s *FieldDeclsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitFieldDecls(s)
	}
}

func (p *FGGParser) FieldDecls() (localctx IFieldDeclsContext) {
	localctx = NewFieldDeclsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, FGGParserRULE_fieldDecls)
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
		p.SetState(147)
		p.FieldDecl()
	}
	p.SetState(152)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGGParserT__3 {
		{
			p.SetState(148)
			p.Match(FGGParserT__3)
		}
		{
			p.SetState(149)
			p.FieldDecl()
		}

		p.SetState(154)
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

	// SetField sets the field token.
	SetField(antlr.Token)

	// IsFieldDeclContext differentiates from other interfaces.
	IsFieldDeclContext()
}

type FieldDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	field  antlr.Token
}

func NewEmptyFieldDeclContext() *FieldDeclContext {
	var p = new(FieldDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_fieldDecl
	return p
}

func (*FieldDeclContext) IsFieldDeclContext() {}

func NewFieldDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldDeclContext {
	var p = new(FieldDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_fieldDecl

	return p
}

func (s *FieldDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldDeclContext) GetField() antlr.Token { return s.field }

func (s *FieldDeclContext) SetField(v antlr.Token) { s.field = v }

func (s *FieldDeclContext) Typ() ITypContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypContext)
}

func (s *FieldDeclContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, 0)
}

func (s *FieldDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterFieldDecl(s)
	}
}

func (s *FieldDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitFieldDecl(s)
	}
}

func (p *FGGParser) FieldDecl() (localctx IFieldDeclContext) {
	localctx = NewFieldDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, FGGParserRULE_fieldDecl)

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
		p.SetState(155)

		var _m = p.Match(FGGParserNAME)

		localctx.(*FieldDeclContext).field = _m
	}
	{
		p.SetState(156)
		p.Typ()
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
	p.RuleIndex = FGGParserRULE_specs
	return p
}

func (*SpecsContext) IsSpecsContext() {}

func NewSpecsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SpecsContext {
	var p = new(SpecsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_specs

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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterSpecs(s)
	}
}

func (s *SpecsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitSpecs(s)
	}
}

func (p *FGGParser) Specs() (localctx ISpecsContext) {
	localctx = NewSpecsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, FGGParserRULE_specs)
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
		p.SetState(158)
		p.Spec()
	}
	p.SetState(163)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGGParserT__3 {
		{
			p.SetState(159)
			p.Match(FGGParserT__3)
		}
		{
			p.SetState(160)
			p.Spec()
		}

		p.SetState(165)
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
	p.RuleIndex = FGGParserRULE_spec
	return p
}

func (*SpecContext) IsSpecContext() {}

func NewSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SpecContext {
	var p = new(SpecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_spec

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

func (s *InterfaceSpecContext) Typ() ITypContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypContext)
}

func (s *InterfaceSpecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterInterfaceSpec(s)
	}
}

func (s *InterfaceSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterSigSpec(s)
	}
}

func (s *SigSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitSigSpec(s)
	}
}

func (p *FGGParser) Spec() (localctx ISpecContext) {
	localctx = NewSpecContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, FGGParserRULE_spec)

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

	p.SetState(168)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSigSpecContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(166)
			p.Sig()
		}

	case 2:
		localctx = NewInterfaceSpecContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(167)
			p.Typ()
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

	// SetMeth sets the meth token.
	SetMeth(antlr.Token)

	// IsSigContext differentiates from other interfaces.
	IsSigContext()
}

type SigContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	meth   antlr.Token
}

func NewEmptySigContext() *SigContext {
	var p = new(SigContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_sig
	return p
}

func (*SigContext) IsSigContext() {}

func NewSigContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SigContext {
	var p = new(SigContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_sig

	return p
}

func (s *SigContext) GetParser() antlr.Parser { return s.parser }

func (s *SigContext) GetMeth() antlr.Token { return s.meth }

func (s *SigContext) SetMeth(v antlr.Token) { s.meth = v }

func (s *SigContext) TypeFormals() ITypeFormalsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypeFormalsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypeFormalsContext)
}

func (s *SigContext) Typ() ITypContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypContext)
}

func (s *SigContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, 0)
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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterSig(s)
	}
}

func (s *SigContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitSig(s)
	}
}

func (p *FGGParser) Sig() (localctx ISigContext) {
	localctx = NewSigContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, FGGParserRULE_sig)
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
		p.SetState(170)

		var _m = p.Match(FGGParserNAME)

		localctx.(*SigContext).meth = _m
	}
	{
		p.SetState(171)
		p.TypeFormals()
	}
	{
		p.SetState(172)
		p.Match(FGGParserT__0)
	}
	p.SetState(174)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == FGGParserNAME {
		{
			p.SetState(173)
			p.Params()
		}

	}
	{
		p.SetState(176)
		p.Match(FGGParserT__1)
	}
	{
		p.SetState(177)
		p.Typ()
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
	p.RuleIndex = FGGParserRULE_params
	return p
}

func (*ParamsContext) IsParamsContext() {}

func NewParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamsContext {
	var p = new(ParamsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_params

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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterParams(s)
	}
}

func (s *ParamsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitParams(s)
	}
}

func (p *FGGParser) Params() (localctx IParamsContext) {
	localctx = NewParamsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, FGGParserRULE_params)
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
		p.SetState(179)
		p.ParamDecl()
	}
	p.SetState(184)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGGParserT__2 {
		{
			p.SetState(180)
			p.Match(FGGParserT__2)
		}
		{
			p.SetState(181)
			p.ParamDecl()
		}

		p.SetState(186)
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

	// SetVari sets the vari token.
	SetVari(antlr.Token)

	// IsParamDeclContext differentiates from other interfaces.
	IsParamDeclContext()
}

type ParamDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	vari   antlr.Token
}

func NewEmptyParamDeclContext() *ParamDeclContext {
	var p = new(ParamDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = FGGParserRULE_paramDecl
	return p
}

func (*ParamDeclContext) IsParamDeclContext() {}

func NewParamDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamDeclContext {
	var p = new(ParamDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_paramDecl

	return p
}

func (s *ParamDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamDeclContext) GetVari() antlr.Token { return s.vari }

func (s *ParamDeclContext) SetVari(v antlr.Token) { s.vari = v }

func (s *ParamDeclContext) Typ() ITypContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypContext)
}

func (s *ParamDeclContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, 0)
}

func (s *ParamDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterParamDecl(s)
	}
}

func (s *ParamDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitParamDecl(s)
	}
}

func (p *FGGParser) ParamDecl() (localctx IParamDeclContext) {
	localctx = NewParamDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, FGGParserRULE_paramDecl)

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
		p.SetState(187)

		var _m = p.Match(FGGParserNAME)

		localctx.(*ParamDeclContext).vari = _m
	}
	{
		p.SetState(188)
		p.Typ()
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
	p.RuleIndex = FGGParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_expr

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
	recv  IExprContext
	targs ITypsContext
	args  IExprsContext
}

func NewCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CallContext {
	var p = new(CallContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *CallContext) GetRecv() IExprContext { return s.recv }

func (s *CallContext) GetTargs() ITypsContext { return s.targs }

func (s *CallContext) GetArgs() IExprsContext { return s.args }

func (s *CallContext) SetRecv(v IExprContext) { s.recv = v }

func (s *CallContext) SetTargs(v ITypsContext) { s.targs = v }

func (s *CallContext) SetArgs(v IExprsContext) { s.args = v }

func (s *CallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CallContext) NAME() antlr.TerminalNode {
	return s.GetToken(FGGParserNAME, 0)
}

func (s *CallContext) Expr() IExprContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *CallContext) Typs() ITypsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypsContext)
}

func (s *CallContext) Exprs() IExprsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprsContext)
}

func (s *CallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterCall(s)
	}
}

func (s *CallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
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
	return s.GetToken(FGGParserNAME, 0)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
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

func (s *AssertContext) Typ() ITypContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypContext)
}

func (s *AssertContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterAssert(s)
	}
}

func (s *AssertContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
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
	return s.GetToken(FGGParserFMT, 0)
}

func (s *SprintfContext) SPRINTF() antlr.TerminalNode {
	return s.GetToken(FGGParserSPRINTF, 0)
}

func (s *SprintfContext) STRING() antlr.TerminalNode {
	return s.GetToken(FGGParserSTRING, 0)
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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterSprintf(s)
	}
}

func (s *SprintfContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
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
	return s.GetToken(FGGParserNAME, 0)
}

func (s *SelectContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterSelect(s)
	}
}

func (s *SelectContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
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

func (s *StructLitContext) Typ() ITypContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITypContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITypContext)
}

func (s *StructLitContext) Exprs() IExprsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExprsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExprsContext)
}

func (s *StructLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterStructLit(s)
	}
}

func (s *StructLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitStructLit(s)
	}
}

func (p *FGGParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *FGGParser) expr(_p int) (localctx IExprContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 34
	p.EnterRecursionRule(localctx, 34, FGGParserRULE_expr, _p)
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
	p.SetState(212)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		localctx = NewVariableContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(191)
			p.Match(FGGParserNAME)
		}

	case 2:
		localctx = NewStructLitContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(192)
			p.Typ()
		}
		{
			p.SetState(193)
			p.Match(FGGParserT__4)
		}
		p.SetState(195)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == FGGParserFMT || _la == FGGParserNAME {
			{
				p.SetState(194)
				p.Exprs()
			}

		}
		{
			p.SetState(197)
			p.Match(FGGParserT__9)
		}

	case 3:
		localctx = NewSprintfContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(199)
			p.Match(FGGParserFMT)
		}
		{
			p.SetState(200)
			p.Match(FGGParserT__7)
		}
		{
			p.SetState(201)
			p.Match(FGGParserSPRINTF)
		}
		{
			p.SetState(202)
			p.Match(FGGParserT__0)
		}
		{
			p.SetState(203)
			_la = p.GetTokenStream().LA(1)

			if !(_la == FGGParserT__8 || _la == FGGParserSTRING) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(208)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<FGGParserT__2)|(1<<FGGParserFMT)|(1<<FGGParserNAME))) != 0 {
			p.SetState(206)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case FGGParserT__2:
				{
					p.SetState(204)
					p.Match(FGGParserT__2)
				}

			case FGGParserFMT, FGGParserNAME:
				{
					p.SetState(205)
					p.expr(0)
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(210)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(211)
			p.Match(FGGParserT__1)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(238)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(236)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
			case 1:
				localctx = NewSelectContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, FGGParserRULE_expr)
				p.SetState(214)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
				}
				{
					p.SetState(215)
					p.Match(FGGParserT__7)
				}
				{
					p.SetState(216)
					p.Match(FGGParserNAME)
				}

			case 2:
				localctx = NewCallContext(p, NewExprContext(p, _parentctx, _parentState))
				localctx.(*CallContext).recv = _prevctx

				p.PushNewRecursionContext(localctx, _startState, FGGParserRULE_expr)
				p.SetState(217)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(218)
					p.Match(FGGParserT__7)
				}
				{
					p.SetState(219)
					p.Match(FGGParserNAME)
				}
				{
					p.SetState(220)
					p.Match(FGGParserT__0)
				}
				p.SetState(222)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == FGGParserNAME {
					{
						p.SetState(221)

						var _x = p.Typs()

						localctx.(*CallContext).targs = _x
					}

				}
				{
					p.SetState(224)
					p.Match(FGGParserT__1)
				}
				{
					p.SetState(225)
					p.Match(FGGParserT__0)
				}
				p.SetState(227)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == FGGParserFMT || _la == FGGParserNAME {
					{
						p.SetState(226)

						var _x = p.Exprs()

						localctx.(*CallContext).args = _x
					}

				}
				{
					p.SetState(229)
					p.Match(FGGParserT__1)
				}

			case 3:
				localctx = NewAssertContext(p, NewExprContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, FGGParserRULE_expr)
				p.SetState(230)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(231)
					p.Match(FGGParserT__7)
				}
				{
					p.SetState(232)
					p.Match(FGGParserT__0)
				}
				{
					p.SetState(233)
					p.Typ()
				}
				{
					p.SetState(234)
					p.Match(FGGParserT__1)
				}

			}

		}
		p.SetState(240)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())
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
	p.RuleIndex = FGGParserRULE_exprs
	return p
}

func (*ExprsContext) IsExprsContext() {}

func NewExprsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprsContext {
	var p = new(ExprsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = FGGParserRULE_exprs

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
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.EnterExprs(s)
	}
}

func (s *ExprsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(FGGListener); ok {
		listenerT.ExitExprs(s)
	}
}

func (p *FGGParser) Exprs() (localctx IExprsContext) {
	localctx = NewExprsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, FGGParserRULE_exprs)
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
		p.SetState(241)
		p.expr(0)
	}
	p.SetState(246)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == FGGParserT__2 {
		{
			p.SetState(242)
			p.Match(FGGParserT__2)
		}
		{
			p.SetState(243)
			p.expr(0)
		}

		p.SetState(248)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

func (p *FGGParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 17:
		var t *ExprContext = nil
		if localctx != nil {
			t = localctx.(*ExprContext)
		}
		return p.Expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *FGGParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
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
