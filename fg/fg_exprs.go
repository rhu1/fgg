/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base (abstract) types, interfaces, etc. are in fg.go.
 */

package fg

import "fmt"
import "strings"

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

/*func (v Variable) CanEval(ds []Decl) bool {
	return false
}*/

/*func (v Variable) IsValue() bool {
	return false
}*/

func (v Variable) Eval(ds []Decl) Expr {
	panic("Cannot evaluate free variable: " + v.id)
}

func (v Variable) Typing(ds []Decl, gamma Env) Type {
	res, ok := gamma[v.id]
	if !ok {
		panic("Var not in env: " + v.String())
	}
	return res
}

func (v Variable) String() string {
	return v.id
}

type StructLit struct {
	t  Type
	es []Expr
}

var _ Expr = StructLit{}

func (s StructLit) Subs(m map[Variable]Expr) Expr {
	return s
}

/*func (s StructLit) CanEval(ds []Decl) bool {
	for _, v := range s.es {
		if v.CanEval(ds) {
			return true
		}
	}
	return false
}*/

/*func (s StructLit) IsValue() bool {
	return true
}*/

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
	return StructLit{s.t, es}
}

func (s StructLit) Typing(ds []Decl, gamma Env) Type {
	fs := fields(ds, s.t)
	if len(s.es) != len(fs) {
		tmp := ""
		if len(fs) > 0 {
			tmp = fs[0].String()
			for _, v := range fs[1:] {
				tmp = tmp + ", " + v.String()
			}
		}
		panic("Arity mismatch: found=" +
			strings.Join(strings.Split(fmt.Sprint(s.es), " "), ", ") +
			", expected=[" + tmp + "]" + "\n\t" + s.String())
	}
	for i := 0; i < len(s.es); i++ {
		t := s.es[i].Typing(ds, gamma)
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

type Select struct {
	e Expr
	f Name
}

func (s Select) Subs(m map[Variable]Expr) Expr {
	return Select{s.e.Subs(m), s.f}
}

/*func (s Select) CanEval(ds []Decl) bool {
	if _, ok := s.e.(StructLit); !ok {
		return false
	}
	v := s.e.(StructLit)
	td := getTDecl(ds, v.t)
	if t_S, ok := td.(STypeLit); !ok { // Unnecessary?
		return false
	} else {
		for i := 0; i < len(t_S.fds); i++ {
			if t_S.fds[i].f == s.f {
				return true
			}
		}
		return false
	}
}*/

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

func (s Select) Typing(ds []Decl, gamma Env) Type {
	t := s.e.Typing(ds, gamma)
	if !isStructType(ds, t) {
		panic("Illegal select on non-struct type expr: " + t)
	}
	td := getTDecl(ds, t).(STypeLit)
	for _, v := range td.fds {
		if v.f == s.f {
			return v.t
		}
	}
	panic("Field not found: " + s.f + " in\n\t" + td.String())
}

func (s Select) String() string {
	return s.e.String() + "." + s.f
}

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

/*func (c Call) CanEval(ds []Decl) bool {
	if c.e.CanEval(ds) {
		return true
	}
	for _, v := range c.args {
		if v.CanEval(ds) {
			return true
		}
	}
	// c.e and c.args cannot eval -- TODO: but not sure if good or bad
	if s, ok := e.(StructLit); ok {
		body(ds, s.t, c.m)
		return true
	}
	return false
}*/

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
	x0, xs, e := body(ds, s.t, c.m)
	subs := make(map[Variable]Expr)
	subs[Variable{x0}] = c.e
	for i := 0; i < len(xs); i++ {
		subs[Variable{xs[i]}] = c.args[i]
	}
	return e.Subs(subs) // N.B. slightly different to R-Call
}

func (c Call) Typing(ds []Decl, gamma Env) Type {
	t0 := c.e.Typing(ds, gamma)
	var s Sig
	if tmp, ok := methods(ds, t0)[c.m]; !ok {
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
			tmp)
	}
	for i := 0; i < len(c.args); i++ {
		t := c.args[i].Typing(ds, gamma)
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

/*
type Assert struct {
	e Expr
	t Name
}

func (a Assert) Subs(m map[Variable]Expr) Expr {
}

func (a Assert) Eval(ds[]Decl) Expr {
}

func (a Assert) Typing(ds []Decl, gamma Env) Type {
}

func (a Assert) String() string {
}
*/

func isValue(e Expr) bool {
	if _, ok := e.(StructLit); ok {
		return true
	}
	return false
}
