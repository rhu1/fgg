/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fgg

import "fmt"
import "strings"

var _ = fmt.Errorf
var _ = strings.Compare

/* Variable */

type Variable struct {
	id Name
}

var _ Expr = Variable{}

func (v Variable) Subs(m map[Variable]Expr) Expr {
	res, ok := m[v]
	if !ok {
		panic("Unknown var: " + v.String())
	}
	return res
}

func (v Variable) Eval(ds []Decl) (Expr, string) {
	panic("Cannot evaluate free variable: " + v.id)
}

func (v Variable) Typing(ds []Decl, gamma Env, delta TEnv, allowStupid bool) Type {
	res, ok := gamma[v]
	if !ok {
		panic("Var not in env: " + v.String())
	}
	return res
}

func (v Variable) String() string {
	return v.id
}
