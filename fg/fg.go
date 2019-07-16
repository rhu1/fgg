package fg

import "reflect"

/* Name, Env, Type */

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

/* AST base intefaces: FGNode, Decl, TDecl, Spec, Expr */

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

type Expr interface {
	FGNode
	Subs(map[Variable]Expr) Expr

	// CHECKME: resulting Exprs are not "parsed" from source, OK?
	// string is the type name of the actual evaluated expr (within the eval context)
	Eval(ds []Decl) (Expr, string)

	//IsPanic() bool  // TODO "explicit" FG panic -- cf. underlying runtime panic

	// N.B. gamma should be effectively immutable (and ds, of course)
	// (No typing rule modifies gamma, except the T-Func bootstrap)
	Typing(ds []Decl, gamma Env, allowStupid bool) Type
}
