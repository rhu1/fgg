package base

/* Name */

type Name = string // Type alias (cf. type def)

/* AST Nodes */

type AstNode interface {
	String() string
}

type Decl interface {
	AstNode
	GetName() Name
	Ok(ds []Decl)
}

type Program interface {
	AstNode
	GetDecls() []Decl
	GetMain() Expr
	Ok(allowStupid bool)     // Set false for source check
	Eval() (Program, string) // Eval one step; string is the name of the (innermost) applied rule
}

type Expr interface {
	AstNode
	IsValue() bool
	ToGoString() string // Basically, type T printed as main.T
}

/* ANTLR (parsing) */

type Adaptor interface {
	Parse(strictParse bool, input string) Program
}
