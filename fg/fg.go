package fg

import "strings"

type Name = string

type Expr interface {
	Subs(map[Variable]Expr) Expr
	Eval() Expr
	String() string
}

var _ Expr = Variable{}
var _ Expr = StructLit{}

type Variable struct {
	n Name
}

func (this Variable) Subs(m map[Variable]Expr) Expr {
	res, ok := m[this]
	if !ok {
		panic("Unknown var: " + this.String())
	}
	return res
}

func (this Variable) Eval() Expr {
	panic(this.n)
}

func (this Variable) String() string {
	return this.n
}

type StructLit struct {
	t  Name
	es []Expr
}

func (this StructLit) Subs(m map[Variable]Expr) Expr {
	return this
}

func (this StructLit) Eval() Expr {
	return this
}

func (this StructLit) String() string {
	var sb strings.Builder
	sb.WriteString(this.t)
	sb.WriteString("{")
	if len(this.es) > 0 {
		sb.WriteString(this.es[0].String())
		for _, v := range this.es[1:] {
			sb.WriteString(", ")
			sb.WriteString(v.String())
		}
		sb.WriteString("}")
	}
	return sb.String()
}

/*
type Select struct {
	e Expr
	f Name
}

func (this Select) Eval() Expr {

}

type Call struct {
	recv Expr
	m    Name
	args []Expr
}

func (this Call) Eval() Expr {

}

type Assert struct {
	e Expr
	t Name
}

func (this assert) Eval() Expr {

}
*/
