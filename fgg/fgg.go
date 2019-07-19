package fgg

import "reflect"
import "strings"

var _ = reflect.Append

/* Name, Type, Type param, Type name -- !!! submission version, "Type name" overloaded */

type Name = string

type Type interface {
	Subs(subs map[TParam]Type) Type
	Impls(ds []Decl, delta TEnv, u Type) bool
	Equals(u Type) bool
	String() string
}

type TParam Name

var _ Type = TParam("")

func (a TParam) Subs(subs map[TParam]Type) Type {
	res, ok := subs[a]
	if !ok {
		panic("Unknown param: " + a.String())
	}
	return res
}

// u0 <: u
func (a TParam) Impls(ds []Decl, delta TEnv, u Type) bool {
	return a == u || a.Impls(ds, delta, bounds(delta, u))
}

func (a TParam) Equals(u Type) bool {
	if b, ok := u.(TParam); ok {
		return a == b
	}
	return false
}

func (a TParam) String() string {
	return string(a)
}

type TName struct {
	t  Name
	us []Type
}

var _ Type = TName{}

func (u0 TName) Subs(subs map[TParam]Type) Type {
	us := make([]Type, len(u0.us))
	for i := 0; i < len(us); i++ {
		us[i] = u0.us[i].Subs(subs)
	}
	return TName{u0.t, us}
}

// u0 <: 1
func (u0 TName) Impls(ds []Decl, delta TEnv, u Type) bool {
	if isStructType(ds, u) {
		return isStructType(ds, u0) && u0.Equals(u) // Asks equality of nested TParam
	}

	gs := methods(ds, u)   // u is a t_I
	gs0 := methods(ds, u0) // t0 may be any
	for k, g := range gs {
		g0, ok := gs0[k]
		if !ok || !g.EqExceptTParamsAndVars(g0) {
			return false
		}
	}
	return true
}

func (u0 TName) Equals(u Type) bool {
	if _, ok := u.(TName); !ok {
		return false
	}
	u1 := u.(TName)
	if u0.t != u1.t || len(u0.us) != len(u1.us) {
		return false
	}
	for i := 0; i < len(u0.us); i++ {
		if !u0.us[i].Equals(u1.us[i]) { // Asks equality of nested TParam
			return false
		}
	}
	return true
}

func (u TName) String() string {
	var b strings.Builder
	b.WriteString(string(u.t))
	b.WriteString("(")
	writeTypes(&b, u.us)
	b.WriteString(")")
	return b.String()
}

/* Context, Type context */

type Env map[Variable]Type

type TEnv map[TParam]Type

/* AST base intefaces: FGGNode, Decl, TDecl, Spec, Expr */

type FGGNode interface {
	String() string
}

type Decl interface {
	FGGNode
	GetName() Name
}

type TDecl interface {
	Decl
	//GetType() Type // == Type(GetName())
}

type Spec interface {
	FGGNode
	GetSigs(ds []Decl) []Sig
}

type Expr interface {
	FGGNode
	Subs(subs map[Variable]Expr) Expr
	Eval(ds []Decl) (Expr, string)
	// Like gamma, delta is effectively immutable
	Typing(ds []Decl, delta TEnv, gamma Env, allowStupid bool) Type
}

/* Helpers */

func isStructType(ds []Decl, u Type) bool {
	if u1, ok := u.(TName); ok {
		for _, v := range ds {
			d, ok := v.(STypeLit)
			if ok && d.t == u1.t {
				return true
			}
		}
	}
	return false
}

func isInterfaceType(ds []Decl, u Type) bool {
	/*if u1, ok := u.(TName); ok {
		for _, v := range ds {
			d, ok := v.(ITypeLit)
			if ok && d.t == u1.t {
				return true
			}
		}
	}*/
	panic("[TODO]: ")
}

func writeTypes(b *strings.Builder, us []Type) {
	if len(us) > 0 {
		b.WriteString(us[0].String())
		for _, v := range us[1:] {
			b.WriteString(", " + v.String())
		}
	}
}
