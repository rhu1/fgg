//rhu@HZHL4 MINGW64 ~/code/go/src/temp/antlr/antlr04
//$ antlr4 -Dlanguage=Go -o parser FG.g4

// Cf. https://github.com/antlr/grammars-v4/blob/master/golang/Golang.g4

// FG.g4
grammar FG;


/* Keywords */

FUNC      : 'func';
INTERFACE : 'interface';
MAIN      : 'main';
PACKAGE   : 'package';
RETURN    : 'return';
STRUCT    : 'struct';
TYPE      : 'type';


/* Tokens */

// i.e., LETTER
fragment NAME_START : ('a' .. 'z') | ('A' .. 'Z') ;
fragment DIGIT      : ('0' .. '9') ;
NAME                : NAME_START (NAME_START | DIGIT | '_')* ;

WHITESPACE   : [ \r\n\t]+ -> skip;
COMMENT      : '/*' .*? '*/'    -> channel(HIDDEN);
LINE_COMMENT : '//' ~[\r\n]*    -> channel(HIDDEN);

/* Rules */

// Conventions:
// "tag=" to distinguish repeat productions within a rule: comes out in field/getter names
// "#tag" for cases within a rule: comes out as Context names (i.e., types)
// "plurals", e.g., decls, used for sequences: comes out as "helper" Context...
// ...nodes that group up actual children underneath, makes adapting easier

program : PACKAGE MAIN ';' decls? FUNC MAIN '(' ')' '{' '_' '=' expr '}' EOF ;

decls     : ((typeDecl | methDecl) ';')+ ;
typeDecl  : TYPE NAME typeLit ;
methDecl  : FUNC '(' paramDecl ')' meth=NAME '(' params? ')' ret=NAME '{' RETURN expr '}' ;
params    : paramDecl (',' paramDecl)* ;
paramDecl : vari=NAME typ=NAME ;

typeLit : STRUCT '{' fieldDecls? '}' # StructTypeLit ;
        //| INTERFACE '{' specs? '}' # InterfaceTypeLit ;

fieldDecls : fieldDecl (';' fieldDecl)* ;
fieldDecl  : field=NAME typ=NAME ;

expr  : NAME                                   # Variable
      | NAME '{' exprs? '}'                    # StructLit
      | expr '.' field=NAME                    # Select
      | recv=expr '.' NAME '(' args=exprs* ')' # Call
      | expr '.' '(' NAME ')'                  # Assert
      ;
exprs : expr (',' expr)* ;

//meth_sig : meth=name '(' parms=formal* ')' ret=name # MethSig ;

//formal : var=name type=name # ParamDecl ;

