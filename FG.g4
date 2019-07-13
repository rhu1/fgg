//rhu@HZHL4 MINGW64 ~/code/go/src/temp/antlr/antlr04
//$ antlr4 -Dlanguage=Go -o parser FG.g4

// Cf. https://github.com/antlr/grammars-v4/blob/master/golang/Golang.g4

// FG.g4
grammar FG;


// Keywords

PACKAGE: 'package';
MAIN: 'main';
STRUCT: 'struct';
INTERFACE: 'interface';
FUNC: 'func';
RETURN: 'return';
TYPE: 'type';


// Tokens

fragment NAME_START  // LETTER
   : ('a' .. 'z')
   | ('A' .. 'Z')
   /*| '+'
   | '-'
   | '*'
   | '/'
   | '.'*/
   ;

NAME
   : NAME_START (NAME_START | DIGIT | '_')*
   ;

fragment DIGIT
   : ('0' .. '9')
   ;

WHITESPACE: [ \r\n\t]+ -> skip;
COMMENT:            '/*' .*? '*/'    -> channel(HIDDEN);
LINE_COMMENT:       '//' ~[\r\n]*    -> channel(HIDDEN);

// Rules
program : PACKAGE MAIN ';' type_decls? FUNC MAIN '(' ')' '{' '_' '=' body=expression '}' EOF ;

type_decls : (type_decl ';')+ ;
type_decl: TYPE name=NAME type_lit ;

type_lit : STRUCT '{' elems=field_decls? '}' # Struct;

field_decls : field_decl (';' field_decl)* ;  // Makes adapting easier, helper context with actual children below
field_decl : field=NAME typ=NAME ;

expression
    : variable=NAME                            # Variable
    | typ=NAME '{' args=exprs? '}'  # Lit
    | expr=expression '.' field=NAME      # Select
    | recv=expression '.' meth=NAME '(' args=exprs* ')'  # Call
    | expr=expression '.' '(' typ=NAME ')'        # Assertion
    ;

exprs : expression (',' expression)* ;

//meth_sig : meth=name '(' parms=formal* ')' ret=name # MethSig ;

//formal : var=name type=name # ParamDecl ;

