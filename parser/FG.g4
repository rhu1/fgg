//rhu@HZHL4 MINGW64 ~/code/go/src/github.com/rhu1/fgg
//$ antlr4 -Dlanguage=Go -o parser/fg parser/FG.g4

// Cf. https://github.com/antlr/grammars-v4/blob/master/golang/Golang.g4
// (This grammar is not at all based on that one, mention for ref only)

// FG.g4
grammar FG;


/* Keywords */

FUNC      : 'func' ;
INTERFACE : 'interface' ;
MAIN      : 'main' ;
PACKAGE   : 'package' ;
RETURN    : 'return' ;
STRUCT    : 'struct' ;
TYPE      : 'type' ;

IMPORT    : 'import' ;
FMT       : 'fmt' ;
PRINTF    : 'Printf' ;


/* Tokens */

fragment LETTER : ('a' .. 'z') | ('A' .. 'Z') ;
fragment DIGIT  : ('0' .. '9') ;
//fragment HACK   : 'ᐸ' | 'ᐳ' ;  // Doesn't seem to work?
fragment HACK   : '\u1438' | '\u1433' | '\u1428' ;  // Hack for monom output
NAME            : (LETTER | '_' | HACK) (LETTER | '_' | DIGIT | HACK)* ;
WHITESPACE      : [ \r\n\t]+ -> skip ;
COMMENT         : '/*' .*? '*/' -> channel(HIDDEN) ;
LINE_COMMENT    : '//' ~[\r\n]* -> channel(HIDDEN) ;


/* Rules */

// Conventions:
// "tag=" to distinguish repeat productions within a rule: comes out in
// field/getter names.
// "#tag" for cases within a rule: comes out as Context names (i.e., types).
// "plurals", e.g., decls, used for sequences: comes out as "helper" Contexts,
// nodes that group up actual children underneath -- makes "adapting" easier.

program    : PACKAGE MAIN ';' 
             (IMPORT '"' FMT '"')?
             decls? FUNC MAIN '(' ')' '{' 
             ('_' '=' expr | FMT '.' PRINTF '(' '"' '%' '#' 'v' '"' ',' expr ')')  // TODO: too permissive re. whitespace
             '}' EOF ;
decls      : ((typeDecl | methDecl) ';')+ ;
typeDecl   : TYPE NAME typeLit ;  // TODO: tag id=NAME, better for adapting (vs., index constants)
methDecl   : FUNC '(' paramDecl ')' sig '{' RETURN expr '}' ;
typeLit    : STRUCT '{' fieldDecls? '}'             # StructTypeLit
           | INTERFACE '{' specs? '}'               # InterfaceTypeLit ;
fieldDecls : fieldDecl (';' fieldDecl)* ;
fieldDecl  : field=NAME typ=NAME ;
specs      : spec (';' spec)* ;
spec       : sig                                    # SigSpec
           | NAME                                   # InterfaceSpec
           ;
sig        : meth=NAME '(' params? ')' ret=NAME ;
params     : paramDecl (',' paramDecl)* ;
paramDecl  : vari=NAME typ=NAME ;
expr       : NAME                                   # Variable
           | NAME '{' exprs? '}'                    # StructLit
           | expr '.' NAME                          # Select
           | recv=expr '.' NAME '(' args=exprs? ')' # Call
           | expr '.' '(' NAME ')'                  # Assert
           ;
exprs      : expr (',' expr)* ;

