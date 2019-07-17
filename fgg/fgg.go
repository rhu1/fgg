package fgg

import "reflect"

var _ = reflect.Append

/* Name, Env, Type */

type Name = string // Type alias (cf. definition)

type Env map[Name]Type

type Type Name // Type definition (cf. alias)

func (t Type) String() string {
	return string(t)
}

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
	String() string
}
