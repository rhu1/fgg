/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fg

import "strings"

/* "Exported" constructors for fgg (monomorph) */

func NewVariable(id Name) Variable               { return Variable{id} }
func NewStructLit(t Type, es []FGExpr) StructLit { return StructLit{t, es} }
func NewSelect(e FGExpr, f Name) Select          { return Select{e, f} }
func NewCall(e FGExpr, m Name, es []FGExpr) Call { return Call{e, m, es} }
func NewAssert(e FGExpr, t Type) Assert          { return Assert{e, t} }

/* Variable */

var _ FGExpr = Variable{}

type Variable struct {
	name Name
}

func (x Variable) Subs(subs map[Variable]FGExpr) FGExpr {
	res, ok := subs[x]
	if !ok {
		panic("Unknown var: " + x.String())
	}
	return res
}

func (x Variable) Eval(ds []Decl) (FGExpr, string) {
	panic("Cannot evaluate free variable: " + x.name)
}

func (x Variable) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	res, ok := gamma[x.name]
	if !ok {
		panic("Var not in env: " + x.String())
	}
	return res
}

func (x Variable) IsValue() bool {
	return false
}

func (x Variable) String() string {
	return x.name
}

func (x Variable) ToGoString() string {
	return x.name
}

/* StructLit */

var _ FGExpr = StructLit{}

type StructLit struct {
	t_S   Type
	elems []FGExpr
}

func (s StructLit) GetType() Type      { return s.t_S }
func (s StructLit) GetElems() []FGExpr { return s.elems }

func (s StructLit) Subs(subs map[Variable]FGExpr) FGExpr {
	es := make([]FGExpr, len(s.elems))
	for i := 0; i < len(s.elems); i++ {
		es[i] = s.elems[i].Subs(subs)
	}
	return StructLit{s.t_S, es}
}

func (s StructLit) Eval(ds []Decl) (FGExpr, string) {
	es := make([]FGExpr, len(s.elems))
	done := false
	var rule string
	for i := 0; i < len(s.elems); i++ {
		v := s.elems[i]
		if !done && !v.IsValue() {
			v, rule = v.Eval(ds)
			done = true
		}
		es[i] = v
	}
	if done {
		return StructLit{s.t_S, es}, rule
	} else {
		panic("Cannot reduce: " + s.String())
	}
}

func (s StructLit) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	fs := fields(ds, s.t_S)
	if len(s.elems) != len(fs) {
		var b strings.Builder
		b.WriteString("Arity mismatch: args=[")
		writeExprs(&b, s.elems)
		b.WriteString("], fields=[")
		writeFieldDecls(&b, fs)
		b.WriteString("]\n\t")
		b.WriteString(s.String())
		panic(b.String())
	}
	for i := 0; i < len(s.elems); i++ {
		t := s.elems[i].Typing(ds, gamma, allowStupid)
		u := fs[i].t
		if !t.Impls(ds, u) {
			panic("Arg expr must implement field type: arg=" + t.String() +
				", field=" + u.String() + "\n\t" + s.String())
		}
	}
	return s.t_S
}

// From base.Expr
func (s StructLit) IsValue() bool {
	for _, v := range s.elems {
		if !v.IsValue() {
			return false
		}
	}
	return true
}

func (s StructLit) String() string {
	var b strings.Builder
	b.WriteString(s.t_S.String())
	b.WriteString("{")
	//b.WriteString(strings.Trim(strings.Join(strings.Split(fmt.Sprint(s.es), " "), ", "), "[]"))
	// ^ No: broken for nested structs
	writeExprs(&b, s.elems)
	b.WriteString("}")
	return b.String()
}

func (s StructLit) ToGoString() string {
	var b strings.Builder
	b.WriteString("main.")
	b.WriteString(s.t_S.String())
	b.WriteString("{")
	writeToGoExprs(&b, s.elems)
	b.WriteString("}")
	return b.String()
}

/* Select */

var _ FGExpr = Select{}

type Select struct {
	e_S   FGExpr
	field Name
}

func (s Select) GetExpr() FGExpr { return s.e_S }
func (s Select) GetField() Name  { return s.field }

func (s Select) Subs(subs map[Variable]FGExpr) FGExpr {
	return Select{s.e_S.Subs(subs), s.field}
}

func (s Select) Eval(ds []Decl) (FGExpr, string) {
	if !s.e_S.IsValue() {
		e, rule := s.e_S.Eval(ds)
		return Select{e.(FGExpr), s.field}, rule
	}
	v := s.e_S.(StructLit)
	fds := fields(ds, v.t_S)
	for i := 0; i < len(fds); i++ {
		if fds[i].name == s.field {
			return v.elems[i], "Select"
		}
	}
	panic("Field not found: " + s.field)
}

func (s Select) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	t := s.e_S.Typing(ds, gamma, allowStupid)
	if !isStructType(ds, t) {
		panic("Illegal select on non-struct type expr: " + t)
	}
	fds := fields(ds, t)
	for _, v := range fds {
		if v.name == s.field {
			return v.t
		}
	}
	panic("Field not found: " + s.field + " in " + t.String())
}

// From base.Expr
func (s Select) IsValue() bool {
	return false
}

func (s Select) String() string {
	return s.e_S.String() + "." + s.field
}

func (s Select) ToGoString() string {
	return s.e_S.ToGoString() + "." + s.field
}

/* Call */

var _ FGExpr = Call{}

type Call struct {
	e_recv FGExpr
	meth   Name
	args   []FGExpr
}

func (c Call) GetReceiver() FGExpr { return c.e_recv }
func (c Call) GetMethod() Name     { return c.meth }
func (c Call) GetArgs() []FGExpr   { return c.args }

func (c Call) Subs(subs map[Variable]FGExpr) FGExpr {
	e := c.e_recv.Subs(subs)
	args := make([]FGExpr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].Subs(subs)
	}
	return Call{e, c.meth, args}
}

func (c Call) Eval(ds []Decl) (FGExpr, string) {
	if !c.e_recv.IsValue() {
		e, rule := c.e_recv.Eval(ds)
		return Call{e.(FGExpr), c.meth, c.args}, rule
	}
	args := make([]FGExpr, len(c.args))
	done := false
	var rule string
	for i := 0; i < len(c.args); i++ {
		e := c.args[i]
		if !done && !e.IsValue() {
			e, rule = e.Eval(ds)
			done = true
		}
		args[i] = e
	}
	if done {
		return Call{c.e_recv, c.meth, args}, rule
	}
	// c.e and c.args all values
	s := c.e_recv.(StructLit)
	x0, xs, e := body(ds, s.t_S, c.meth) // panics if method not found
	subs := make(map[Variable]FGExpr)
	subs[Variable{x0}] = c.e_recv
	for i := 0; i < len(xs); i++ {
		subs[Variable{xs[i]}] = c.args[i]
	}
	return e.Subs(subs), "Call" // N.B. single combined substitution map slightly different to R-Call
}

func (c Call) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	t0 := c.e_recv.Typing(ds, gamma, allowStupid)
	var g Sig
	if tmp, ok := methods(ds, t0)[c.meth]; !ok { // !!! submission version had "methods(m)"
		panic("Method not found: " + c.meth + " in " + t0.String() + "\n\t" + c.String())
	} else {
		g = tmp
	}
	if len(c.args) != len(g.pDecls) {
		var b strings.Builder
		b.WriteString("Arity mismatch: args=[")
		writeExprs(&b, c.args)
		b.WriteString("], params=[")
		writeParamDecls(&b, g.pDecls)
		b.WriteString("]")
		panic(b.String())
	}
	for i := 0; i < len(c.args); i++ {
		t := c.args[i].Typing(ds, gamma, allowStupid)
		if !t.Impls(ds, g.pDecls[i].t) {
			panic("Arg expr type must implement param type: arg=" + t + ", param=" +
				g.pDecls[i].t)
		}
	}
	return g.t_ret
}

// From base.Expr
func (c Call) IsValue() bool {
	return false
}

func (c Call) String() string {
	var b strings.Builder
	b.WriteString(c.e_recv.String())
	b.WriteString(".")
	b.WriteString(c.meth)
	b.WriteString("(")
	writeExprs(&b, c.args)
	b.WriteString(")")
	return b.String()
}

func (c Call) ToGoString() string {
	var b strings.Builder
	b.WriteString(c.e_recv.ToGoString())
	b.WriteString(".")
	b.WriteString(c.meth)
	b.WriteString("(")
	writeToGoExprs(&b, c.args)
	b.WriteString(")")
	return b.String()
}

/* Assert */

var _ FGExpr = Assert{}

type Assert struct {
	e_I FGExpr
	t   Type
}

func (a Assert) GetExpr() FGExpr { return a.e_I }
func (a Assert) GetType() Type   { return a.t }

func (a Assert) Subs(subs map[Variable]FGExpr) FGExpr {
	return Assert{a.e_I.Subs(subs), a.t}
}

func (a Assert) Eval(ds []Decl) (FGExpr, string) {
	if !a.e_I.IsValue() {
		e, rule := a.e_I.Eval(ds)
		return Assert{e.(FGExpr), a.t}, rule
	}
	t := a.e_I.(StructLit).t_S
	if !isStructType(ds, t) {
		panic("Non struct type found in struct lit: " + t)
	}
	if t.Impls(ds, a.t) {
		return a.e_I, "Assert"
	}
	panic("Cannot reduce: " + a.String())
}

func (a Assert) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	t := a.e_I.Typing(ds, gamma, allowStupid)
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

// From base.Expr
func (a Assert) IsValue() bool {
	return false
}

func (a Assert) String() string {
	var b strings.Builder
	b.WriteString(a.e_I.String())
	b.WriteString(".(")
	b.WriteString(a.t.String())
	b.WriteString(")")
	return b.String()
}

func (a Assert) ToGoString() string {
	var b strings.Builder
	b.WriteString(a.e_I.ToGoString())
	b.WriteString(".(main.")
	b.WriteString(a.t.String())
	b.WriteString(")")
	return b.String()
}

/* Aux, helpers */

func writeExprs(b *strings.Builder, es []FGExpr) {
	if len(es) > 0 {
		b.WriteString(es[0].String())
		for _, v := range es[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
}

func writeToGoExprs(b *strings.Builder, es []FGExpr) {
	if len(es) > 0 {
		b.WriteString(es[0].ToGoString())
		for _, v := range es[1:] {
			b.WriteString(", ")
			b.WriteString(v.ToGoString())
		}
	}
}
