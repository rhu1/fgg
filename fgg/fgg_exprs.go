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

func (x Variable) Subs(m map[Variable]Expr) Expr {
	res, ok := m[x]
	if !ok {
		panic("Unknown var: " + x.String())
	}
	return res
}

func (x Variable) Eval(ds []Decl) (Expr, string) {
	panic("Cannot evaluate free variable: " + x.id)
}

func (x Variable) Typing(ds []Decl, delta TEnv, gamma Env,
	allowStupid bool) Type {
	res, ok := gamma[x.id]
	if !ok {
		panic("Var not in env: " + x.String())
	}
	return res
}

func (x Variable) String() string {
	return x.id
}

/* StructLit */

type StructLit struct {
	u  TName // u.t is a t_S
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
	s.u.Ok(ds, delta)
	fs := fields(ds, s.u)
	if len(s.es) != len(fs) {
		var b strings.Builder
		b.WriteString("Arity mismatch: args=[")
		writeExprs(&b, s.es)
		b.WriteString("], fields=[")
		writeFieldDecls(&b, fs)
		b.WriteString("]\n\t")
		b.WriteString(s.String())
		panic(b.String())
	}
	for i := 0; i < len(s.es); i++ {
		u := s.es[i].Typing(ds, delta, gamma, allowStupid)
		r := fs[i].u
		if !u.Impls(ds, delta, r) {
			panic("Arg expr must implement field type: arg=" + u.String() +
				", field=" + r.String() + "\n\t" + s.String())
		}
	}
	return s.u
}

func (s StructLit) String() string {
	var b strings.Builder
	b.WriteString(s.u.String())
	b.WriteString("{")
	writeExprs(&b, s.es)
	b.WriteString("}")
	return b.String()
}

/* Aux, helpers */

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

func writeExprs(b *strings.Builder, es []Expr) {
	if len(es) > 0 {
		b.WriteString(es[0].String())
		for _, v := range es[1:] {
			b.WriteString(", " + v.String())
		}
	}
}
