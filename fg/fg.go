package fg

import "reflect"
import "strings"

type Name = string // Type alias (cf. definition)
type Env map[Name]Type

type Type Name // Type definition (cf. alias)

var _ Spec = Type("")

// Pre: t0, t are known types
// t0 <: t
func (t0 Type) Impls(ds []Decl, t Type) bool {
	if isStructType(ds, t) {
		return isStructType(ds, t0) && t0 == t
	}

	m := methods(ds, t)   // t is a t_I
	m0 := methods(ds, t0) // t0 may be any
	for k, s := range m {
		s0, ok := m0[k]
		if !ok || !s.EqExceptVars(s0) {
			return false
		}
	}
	return true
}

// t_I is a Spec, but not t_S -- this aspect is currently "dynamically typed"
func (t Type) GetSigs(ds []Decl) []Sig {
	if !isInterfaceType(ds, t) { // isStructType would be more efficient
		panic("Cannot use non-interface type as a Spec: " + t.String() +
			" is a " + reflect.TypeOf(t).String())
	}
	td := getTDecl(ds, t).(ITypeLit)
	var res []Sig
	for _, s := range td.ss {
		res = append(res, s.GetSigs(ds)...)
	}
	return res
}

func (t Type) String() string {
	return string(t)
}

func isStructType(ds []Decl, t Type) bool {
	for _, v := range ds {
		d, ok := v.(STypeLit)
		if ok && d.t == t {
			return true
		}
	}
	return false
}

func isInterfaceType(ds []Decl, t Type) bool {
	for _, v := range ds {
		d, ok := v.(ITypeLit)
		if ok && d.t == t {
			return true
		}
	}
	return false
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
	GetSigs(ds []Decl) []Sig
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

func (s Sig) GetSigs(_ []Decl) []Sig {
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
	//CanEval(ds []Decl) bool // Should only panic if badly typed, o/w return false if stuck
	//IsValue() bool
	Eval(ds []Decl) Expr
	//IsPanic() bool  // TODO
	Typing(ds []Decl, gamma Env) Type
	// N.B. gamma should be effectively immutable (and ds, of course)
	// (No typing rule adds to gamma, except T-Func bootstrap)
}
