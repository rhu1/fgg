package fgg

import "reflect"
import "strings"

var _ = reflect.Append

/* Name, Type, Type param, Type name -- !!! submission version, "Type name" overloaded */

type Name = string

type Type interface {
	//Subs(map[TParam]Type)
	//Impls()
	String() string
}

type TParam Name

func (a TParam) String() string {
	return string(a)
}

type TName struct {
	t    Name
	typs []TName
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
}
