package base

type Adaptor interface {
	Parse(strictParse bool, input string) Program
}

type Name = string // Type alias (cf. definition)

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
	GetExpr() Expr
	Ok(allowStupid bool)
	Eval() (Program, string) // Eval one step; string is the name of the (innermost) applied rule
}

type Expr interface {
	AstNode
	IsValue() bool
}
