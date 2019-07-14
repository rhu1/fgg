package fg

import "strings"

type Name = string
type Env map[Name]Type

type Type Name

// t0 <: t
func (t0 Type) Impls(ds []Decl, t Type) bool {
	if isStructType(ds, t) {
		return isStructType(ds, t0) && t0 == t
	}

	m := methods(ds, t)   // t is a t_I
	m0 := methods(ds, t0) // t0 may be any
	for k, md := range m {
		md0, ok := m0[k]
		if !ok || !md.ToSig().EqExceptVars(md0.ToSig()) {
			return false
		}
	}
	return true
}

func (t Type) String() string {
	return string(t)
}

// Base interface for all AST nodes
type FGNode interface {
	String() string
}

type Decl interface {
	FGNode
	GetName() Name
}

type TDecl interface {
	Decl
	GetType() Type // == Type(GetName())
}

type Spec interface {
	FGNode
	GetSigs() []Sig
}

type Sig struct {
	m  Name
	ps []ParamDecl
	t  Type
}

var _ Spec = Sig{}

// !!! Sig in FG (also, Go spec) includes ~x, which breaks "impls"
func (s0 Sig) EqExceptVars(s Sig) bool {
	if len(s0.ps) != len(s.ps) {
		return false
	}
	for i := 0; i < len(s0.ps); i++ {
		if s0.ps[i].t != s.ps[i].t {
			return false
		}
	}
	return s0.m == s.m && s0.t == s.t
}

func (s Sig) GetSigs() []Sig {
	return []Sig{s}
}

func (s Sig) String() string {
	var b strings.Builder
	b.WriteString(s.m)
	b.WriteString("(")
	if len(s.ps) > 0 {
		b.WriteString(s.ps[0].String())
		for _, v := range s.ps[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
	b.WriteString(") ")
	b.WriteString(s.t.String())
	return b.String()
}

type Expr interface {
	FGNode
	Subs(map[Variable]Expr) Expr
	Eval() Expr
	//IsPanic() bool
	Typing(ds []Decl, gamma Env) Type
}
