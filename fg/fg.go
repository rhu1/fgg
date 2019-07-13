package fg

import "strings"

type Name = string

type FGNode interface {
	String() string
}

type FGProgram struct {
	decls []TypeLit
	body  Expr
}

func (p FGProgram) String() string {
	var b strings.Builder
	b.WriteString("package main;\n")
	for _, v := range p.decls {
		b.WriteString(v.String())
		b.WriteString("\n")
	}
	b.WriteString("func main() {\n")
	b.WriteString("\t_ = ")
	b.WriteString(p.body.String())
	b.WriteString("\n")
	b.WriteString("}")
	return b.String()
}

var _ TypeLit = TStruct{}
var _ FGNode = FieldDecl{}

type TypeLit interface {
	FGNode
	GetType() Name
}

type TStruct struct {
	typ Name
	//elems map[Name]Name // N.B. Unordered -- OK?
	elems []FieldDecl
}

func (s TStruct) GetType() Name {
	return s.typ
}

func (s TStruct) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(s.typ)
	b.WriteString(" struct {")
	if len(s.elems) > 0 {
		b.WriteString(" ")
		b.WriteString(s.elems[0].String())
		for _, v := range s.elems[1:] {
			b.WriteString("; ")
			b.WriteString(v.String())
		}
		b.WriteString(" ")
	}
	b.WriteString("}")
	return b.String()
}

type FieldDecl struct {
	field Name
	typ   Name
}

func (fd FieldDecl) String() string {
	return fd.field + " " + fd.typ
}

type Expr interface {
	FGNode
	Subs(map[Variable]Expr) Expr
	Eval() Expr
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
	}
	sb.WriteString("}")
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
