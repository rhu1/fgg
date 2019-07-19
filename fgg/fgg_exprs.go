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

func (v Variable) Typing(ds []Decl, delta TEnv, gamma Env, allowStupid bool) Type {
	res, ok := gamma[v]
	if !ok {
		panic("Var not in env: " + v.String())
	}
	return res
}

func (v Variable) String() string {
	return v.id
}

/* StructLit */

type StructLit struct {
	u  Type // u is a TName, and u.(TName).t is a t_S
	es []Expr
}

var _ Expr = StructLit{}

func (s StructLit) Subs(m map[Variable]Expr) Expr {
	es := make([]Expr, len(s.es))
	for i := 0; i < len(s.es); i++ {
		es[i] = s.es[i].Subs(m)
	}
	return StructLit{s.u, es}
}

func (s StructLit) Eval(ds []Decl) (Expr, string) {
	es := make([]Expr, len(s.es))
	done := false
	var rule string
	for i := 0; i < len(s.es); i++ {
		v := s.es[i]
		if !done && !IsValue(v) {
			v, rule = v.Eval(ds)
			done = true
		}
		es[i] = v
	}
	if done {
		return StructLit{s.u, es}, rule
	} else {
		panic("Cannot reduce: " + s.String())
	}
}

func (s StructLit) Typing(ds []Decl, delta TEnv, gamma Env,
	allowStupid bool) Type {
	fs := fields(ds, s.u.(TName))
	if len(s.es) != len(fs) {
		tmp := ""
		if len(fs) > 0 {
			tmp = fs[0].String()
			for _, v := range fs[1:] {
				tmp = tmp + ", " + v.String()
			}
		}
		panic("Arity mismatch: args=" +
			strings.Join(strings.Split(fmt.Sprint(s.es), " "), ", ") +
			", fields=[" + tmp + "]" + "\n\t" + s.String())
	}
	for i := 0; i < len(s.es); i++ {
		u := s.es[i].Typing(ds, delta, gamma, allowStupid)
		r := fs[i].u
		if !u.Impls(ds, delta, r) {
			panic("Arg expr must impl field type: arg=" + u.String() + ", field=" +
				r.String() + "\n\t" + s.String())
		}
	}
	return s.u
}

func (s StructLit) String() string {
	var b strings.Builder
	b.WriteString(s.u.String())
	b.WriteString("{")
	//b.WriteString(strings.Trim(strings.Join(strings.Split(fmt.Sprint(s.es), " "), ", "), "[]"))
	// ^ No: broken for nested structs
	if len(s.es) > 0 {
		b.WriteString(s.es[0].String())
		for _, v := range s.es[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
	b.WriteString("}")
	return b.String()
}

/* Helper */

// Cf. checkErr
func IsValue(e Expr) bool {
	if s, ok := e.(StructLit); ok {
		for _, v := range s.es {
			if !IsValue(v) {
				return false
			}
		}
		return true
	}
	return false
}
