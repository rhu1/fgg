package fgg

import "reflect"
import "strings"

var _ = reflect.Append

/* Name, Type, Type param, Type name -- !!! submission version, "Type name" overloaded */

type Name = string

type Type interface {
	//Subs(map[TParam]Type)
	Impls(ds []Decl, delta TEnv, u Type) bool
	String() string
}

type TParam Name

func (u0 TParam) Impls(ds []Decl, delta TEnv, u Type) bool {
	return u0 == u || u0.Impls(ds, delta, bounds(delta, u))
}

func (a TParam) String() string {
	return string(a)
}

type TName struct {
	t    Name
	typs []TName
}

func (u0 TName) Impls(ds []Decl, delta TEnv, u Type) bool {
	return true // TODO FIXME
}

func (t TName) String() string {
	var b strings.Builder
	b.WriteString(string(t.t))
	if len(t.typs) > 0 {
		b.WriteString(t.typs[0].String())
		for _, v := range t.typs[1:] {
			b.WriteString(", " + v.String())
		}
	}
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
	GetType() Type // == Type(GetName())
}

type Expr interface {
	FGGNode
	Subs(subs map[Variable]Expr) Expr
	Eval(ds []Decl) (Expr, string)
	// Like gamma, delta is effectively immutable
	Typing(ds []Decl, delta TEnv, gamma Env, allowStupid bool) Type
}
