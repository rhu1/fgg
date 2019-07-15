/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fg

import "fmt"
import "reflect"
import "strings"

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

func (v Variable) Eval(ds []Decl) Expr {
	panic("Cannot evaluate free variable: " + v.id)
}

func (v Variable) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	res, ok := gamma[v.id]
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
	t  Type
	es []Expr
}

var _ Expr = StructLit{}

func (s StructLit) Subs(m map[Variable]Expr) Expr {
	es := make([]Expr, len(s.es))
	for i := 0; i < len(s.es); i++ {
		es[i] = s.es[i].Subs(m)
	}
	return StructLit{s.t, es}
}

func (s StructLit) Eval(ds []Decl) Expr {
	done := false
	es := make([]Expr, len(s.es))
	for i := 0; i < len(s.es); i++ {
		v := s.es[i]
		if !done && !isValue(v) {
			v = v.Eval(ds)
			done = true
		}
		es[i] = v
	}
	if done {
		return StructLit{s.t, es}
	} else {
		for _, v := range s.es {
			fmt.Println("aaa: " + reflect.TypeOf(v).String() + "\n" + v.String())
		}
		panic("Cannot reduce: " + s.String())
	}
}

func (s StructLit) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	fs := fields(ds, s.t)
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
		t := s.es[i].Typing(ds, gamma, allowStupid)
		u := fs[i].t
		if !t.Impls(ds, u) {
			panic("Arg expr must impl field type: arg=" + t.String() + ", field=" +
				u.String() + "\n\t" + s.String())
		}
	}
	return s.t
}

func (s StructLit) String() string {
	var sb strings.Builder
	sb.WriteString(s.t.String())
	sb.WriteString("{")
	sb.WriteString(strings.Trim(strings.Join(strings.Split(fmt.Sprint(s.es), " "), ", "), "[]"))
	sb.WriteString("}")
	return sb.String()
}

/* Select */

type Select struct {
	e Expr
	f Name
}

func (s Select) Subs(m map[Variable]Expr) Expr {
	return Select{s.e.Subs(m), s.f}
}

func (s Select) Eval(ds []Decl) Expr {
	if !isValue(s.e) {
		e := s.e.Eval(ds)
		return Select{e, s.f}
	}
	v := s.e.(StructLit)
	fds := fields(ds, v.t)
	for i := 0; i < len(fds); i++ {
		if fds[i].f == s.f {
			return v.es[i]
		}
	}
	panic("Field not found: " + s.f)
}

func (s Select) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	t := s.e.Typing(ds, gamma, allowStupid)
	if !isStructType(ds, t) {
		panic("Illegal select on non-struct type expr: " + t)
	}
	fds := fields(ds, t)
	for _, v := range fds {
		if v.f == s.f {
			return v.t
		}
	}
	panic("Field not found: " + s.f + " in" + t.String())
}

func (s Select) String() string {
	return s.e.String() + "." + s.f
}

/* Call */

type Call struct {
	e    Expr
	m    Name
	args []Expr
}

func (c Call) Subs(m map[Variable]Expr) Expr {
	e := c.e.Subs(m)
	args := make([]Expr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].Subs(m)
	}
	return Call{e, c.m, args}
}

func (c Call) Eval(ds []Decl) Expr {
	if !isValue(c.e) {
		e := c.e.Eval(ds)
		return Call{e, c.m, c.args}
	}
	args := make([]Expr, len(c.args))
	done := false
	for i := 0; i < len(c.args); i++ {
		e := c.args[i]
		if !done && !isValue(e) {
			e = e.Eval(ds)
			done = true
		}
		args[i] = e
	}
	if done {
		return Call{c.e, c.m, args}
	}
	// c.e and c.args all values
	s := c.e.(StructLit)
	x0, xs, e := body(ds, s.t, c.m) // panics if method not found
	subs := make(map[Variable]Expr)
	subs[Variable{x0}] = c.e
	for i := 0; i < len(xs); i++ {
		subs[Variable{xs[i]}] = c.args[i]
	}
	return e.Subs(subs) // N.B. slightly different to R-Call
}

func (c Call) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	t0 := c.e.Typing(ds, gamma, allowStupid)
	var s Sig
	if tmp, ok := methods(ds, t0)[c.m]; !ok { // !!! submission version had "methods(m)"
		panic("Method not found: " + c.m + " in " + t0.String())
	} else {
		s = tmp
	}
	if len(c.args) != len(s.ps) {
		tmp := "" // TODO: factor out with StructLit.Typing
		if len(s.ps) > 0 {
			tmp = s.ps[0].String()
			for _, v := range s.ps[1:] {
				tmp = tmp + ", " + v.String()
			}
		}
		panic("Arity mismatch: args=" +
			strings.Join(strings.Split(fmt.Sprint(c.args), " "), ", ") + ", params=" +
			"[" + tmp + "]")
	}
	for i := 0; i < len(c.args); i++ {
		t := c.args[i].Typing(ds, gamma, allowStupid)
		if !t.Impls(ds, s.ps[i].t) {
			panic("Arg expr type must implement param type: arg=" + t + ", param=" +
				s.ps[i].t)
		}
	}
	return s.t
}

func (c Call) String() string {
	var b strings.Builder
	b.WriteString(c.e.String())
	b.WriteString(".")
	b.WriteString(c.m)
	b.WriteString("(")
	if len(c.args) > 0 {
		b.WriteString(c.args[0].String())
		for _, v := range c.args[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
	b.WriteString(")")
	return b.String()
}

/* Assert */

type Assert struct {
	e Expr
	t Type
}

func (a Assert) Subs(m map[Variable]Expr) Expr {
	return Assert{a.e.Subs(m), a.t}
}

func (a Assert) Eval(ds []Decl) Expr {
	if !isValue(a.e) {
		return Assert{a.e.Eval(ds), a.t}
	}
	t_S := typ(ds, a.e.(StructLit)) // panics if StructLit.t is not a t_S
	if t_S.Impls(ds, a.t) {
		return a.e
	}
	panic("Cannot reduce: " + a.String())
}

func (a Assert) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	t := a.e.Typing(ds, gamma, allowStupid)
	if isStructType(ds, t) {
		if allowStupid {
			return a.t
		} else {
			panic("Expr must be an interface type (in a non-stupid context): found " +
				t.String() + " for\n\t" + a.String())
		}
	}
	// t is an interface type
	if isInterfaceType(ds, a.t) {
		return a.t // No further checks -- N.B., Robert said they are looking to refine this
	}
	// a.t is a struct type
	if a.t.Impls(ds, t) {
		return a.t
	}
	panic("Struct type assertion must implement expr type: asserted=" +
		a.t.String() + ", expr=" + t.String())
}

func (a Assert) String() string {
	return a.e.String() + ".(" + a.t.String() + ")"
}

/* Helper */

// Cf. checkErr
func isValue(e Expr) bool {
	if s, ok := e.(StructLit); ok {
		for _, v := range s.es {
			if !isValue(v) {
				return false
			}
		}
		return true
	}
	return false
}
