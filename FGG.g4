//rhu@HZHL4 MINGW64 ~/code/go/src/github.com/rhu1/fgg
//$ antlr4 -Dlanguage=Go -o parser FGG.g4

// Cf. https://github.com/antlr/grammars-v4/blob/master/golang/Golang.g4
// (This grammar is not at all based on that one, mention for ref only)

// FGGt.g4
grammar FGG;


/* Keywords */

FUNC      : 'func' ;
INTERFACE : 'interface' ;
MAIN      : 'main' ;
PACKAGE   : 'package' ;
RETURN    : 'return' ;
STRUCT    : 'struct' ;
TYPE      : 'type' ;


/* Tokens */

fragment LETTER : ('a' .. 'z') | ('A' .. 'Z') ;
fragment DIGIT  : ('0' .. '9') ;
NAME            : LETTER (LETTER | DIGIT | '_')* ;
WHITESPACE      : [ \r\n\t]+ -> skip ;
COMMENT         : '/*' .*? '*/' -> channel(HIDDEN) ;
LINE_COMMENT    : '//' ~[\r\n]* -> channel(HIDDEN) ;


/* Rules */

// Conventions:
// "tag=" to distinguish repeat productions within a rule: comes out in field/getter names
// "#tag" for cases within a rule: comes out as Context names (i.e., types)
// "plurals", e.g., decls, used for sequences: comes out as "helper" Context..
// ..nodes that group up actual children underneath -- makes "adapting" easier

typ        : NAME                 # TypeParam
           | NAME '(' typs? ')'   # TypeName
           ;
typs       : typ (',' typ)*  ;
typeFormal : TYPE typeForms? ;
typeForms  : typeForm (',' typeForm)* ;
typeForm   : NAME typ ;

program    : PACKAGE MAIN ';' decls? FUNC MAIN '(' ')' '{' '_' '=' expr '}' EOF ;
decls      : ((typeDecl | methDecl) ';')+ ;
typeDecl   : TYPE NAME typeLit ;  // TODO: tag id=NAME, better for adapting (vs., index constants)
methDecl   : FUNC '(' paramDecl ')' sig '{' RETURN expr '}' ;
typeLit    : STRUCT '{' fieldDecls? '}'             # StructTypeLit
           | INTERFACE '{' specs? '}'               # InterfaceTypeLit ;
fieldDecls : fieldDecl (';' fieldDecl)* ;
fieldDecl  : field=NAME typ ;
specs      : spec (';' spec)* ;
spec       : sig                                    # SigSpec
           | typ                                    # InterfaceSpec
           ;
sig        : meth=NAME '(' params? ')' ret=typ ;  // TODO: meth-tparams
params     : paramDecl (',' paramDecl)* ;
paramDecl  : vari=NAME typ ;
expr       : NAME                                   # Variable
           | typ '{' exprs? '}'                     # StructLit
           | expr '.' NAME                          # Select
           | recv=expr '.' NAME '(' args=exprs? ')' # Call  // TODO: meth-targs
           | expr '.' '(' typ ')'                   # Assert
           ;
exprs      : expr (',' expr)* ;

