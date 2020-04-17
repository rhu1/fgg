/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fgg

import "fmt"
import "reflect"
import "strings"

var _ = fmt.Errorf
var _ = reflect.Append
var _ = strings.Compare

/* Public constructors */

func NewVariable(id Name) Variable { return Variable{id} }

/* Variable */

type Variable struct {
	name Name
}

var _ FGGExpr = Variable{}

func (x Variable) GetName() Name { return x.name }

func (x Variable) Subs(m map[Variable]FGGExpr) FGGExpr {
	res, ok := m[x]
	if !ok {
		panic("Unknown var: " + x.String())
	}
	return res
}

func (x Variable) TSubs(subs map[TParam]Type) FGGExpr {
	return x
}

func (x Variable) Eval(ds []Decl) (FGGExpr, string) {
	panic("Cannot evaluate free variable: " + x.name)
}

// TODO: refactor as Typing and StupidTyping? (clearer than bool param)
func (x Variable) Typing(ds []Decl, delta Delta, gamma Gamma,
	allowStupid bool) Type {
	res, ok := gamma[x.name]
	if !ok {
		panic("Var not in env: " + x.String())
	}
	return res
}

// From base.Expr
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

type StructLit struct {
	u_S   TNamed // u.t is a t_S
	elems []FGGExpr
}

var _ FGGExpr = StructLit{}

func (s StructLit) GetNamedType() TNamed { return s.u_S }
func (s StructLit) GetElems() []FGGExpr  { return s.elems }

func (s StructLit) Subs(subs map[Variable]FGGExpr) FGGExpr {
	es := make([]FGGExpr, len(s.elems))
	for i := 0; i < len(s.elems); i++ {
		es[i] = s.elems[i].Subs(subs)
	}
	return StructLit{s.u_S, es}
}

func (s StructLit) TSubs(subs map[TParam]Type) FGGExpr {
	es := make([]FGGExpr, len(s.elems))
	for i := 0; i < len(s.elems); i++ {
		es[i] = s.elems[i].TSubs(subs)
	}
	return StructLit{s.u_S.TSubs(subs).(TNamed), es}
}

func (s StructLit) Eval(ds []Decl) (FGGExpr, string) {
	es := make([]FGGExpr, len(s.elems))
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
		return StructLit{s.u_S, es}, rule
	} else {
		panic("Cannot reduce: " + s.String())
	}
}

func (s StructLit) Typing(ds []Decl, delta Delta, gamma Gamma,
	allowStupid bool) Type {
	s.u_S.Ok(ds, delta)
	fs := fields(ds, s.u_S)
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
		u := s.elems[i].Typing(ds, delta, gamma, allowStupid)
		r := fs[i].u
		if !u.ImplsDelta(ds, delta, r) {
			panic("Arg expr must implement field type: arg=" + u.String() +
				", field=" + r.String() + "\n\t" + s.String())
		}
	}
	return s.u_S
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
	b.WriteString(s.u_S.String())
	b.WriteString("{")
	writeExprs(&b, s.elems)
	b.WriteString("}")
	return b.String()
}

func (s StructLit) ToGoString() string {
	var b strings.Builder
	b.WriteString(s.u_S.ToGoString())
	b.WriteString("{")
	writeToGoExprs(&b, s.elems)
	b.WriteString("}")
	return b.String()
}

/* Select */

type Select struct {
	e_S   FGGExpr
	field Name
}

var _ FGGExpr = Select{}

func (s Select) GetExpr() FGGExpr { return s.e_S }
func (s Select) GetField() Name   { return s.field }

func (s Select) Subs(subs map[Variable]FGGExpr) FGGExpr {
	return Select{s.e_S.Subs(subs), s.field}
}

func (s Select) TSubs(subs map[TParam]Type) FGGExpr {
	return Select{s.e_S.TSubs(subs), s.field}
}

func (s Select) Eval(ds []Decl) (FGGExpr, string) {
	if !s.e_S.IsValue() {
		e, rule := s.e_S.Eval(ds)
		return Select{e, s.field}, rule
	}
	v := s.e_S.(StructLit)
	fds := fields(ds, v.u_S)
	for i := 0; i < len(fds); i++ {
		if fds[i].field == s.field {
			return v.elems[i], "Select"
		}
	}
	panic("Field not found: " + s.field)
}

func (s Select) Typing(ds []Decl, delta Delta, gamma Gamma,
	allowStupid bool) Type {
	u := s.e_S.Typing(ds, delta, gamma, allowStupid)
	if !IsStructType(ds, u) {
		panic("Illegal select on non-struct type expr: " + u.String())
	}
	fds := fields(ds, u.(TNamed))
	for _, v := range fds {
		if v.field == s.field {
			return v.u
		}
	}
	panic("Field not found: " + s.field + " in " + u.String())
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

type Call struct {
	e_recv FGGExpr
	meth   Name
	t_args []Type // Rename u_args?
	args   []FGGExpr
}

var _ FGGExpr = Call{}

func (c Call) GetRecv() FGGExpr   { return c.e_recv }
func (c Call) GetMethod() Name    { return c.meth }
func (c Call) GetTArgs() []Type   { return c.t_args }
func (c Call) GetArgs() []FGGExpr { return c.args }

func (c Call) Subs(subs map[Variable]FGGExpr) FGGExpr {
	e := c.e_recv.Subs(subs)
	args := make([]FGGExpr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].Subs(subs)
	}
	return Call{e, c.meth, c.t_args, args}
}

func (c Call) TSubs(subs map[TParam]Type) FGGExpr {
	targs := make([]Type, len(c.t_args))
	for i := 0; i < len(c.t_args); i++ {
		targs[i] = c.t_args[i].TSubs(subs)
	}
	args := make([]FGGExpr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].TSubs(subs)
	}
	return Call{c.e_recv.TSubs(subs), c.meth, targs, args}
}

func (c Call) Eval(ds []Decl) (FGGExpr, string) {
	if !c.e_recv.IsValue() {
		e, rule := c.e_recv.Eval(ds)
		return Call{e, c.meth, c.t_args, c.args}, rule
	}
	args := make([]FGGExpr, len(c.args))
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
		return Call{c.e_recv, c.meth, c.t_args, args}, rule
	}
	// c.e and c.args all values
	s := c.e_recv.(StructLit)
	x0, xs, e := body(ds, s.u_S, c.meth, c.t_args) // panics if method not found
	subs := make(map[Variable]FGGExpr)
	subs[Variable{x0}] = c.e_recv
	for i := 0; i < len(xs); i++ {
		subs[Variable{xs[i].name}] = c.args[i]
	}
	return e.Subs(subs), "Call" // N.B. single combined substitution map slightly different to R-Call
}

func (c Call) Typing(ds []Decl, delta Delta, gamma Gamma, allowStupid bool) Type {
	u0 := c.e_recv.Typing(ds, delta, gamma, allowStupid)
	var g Sig
	if tmp, ok := methods(ds, bounds(delta, u0))[c.meth]; !ok { // !!! submission version had "methods(m)"
		panic("Method not found: " + c.meth + " in " + u0.String())
	} else {
		g = tmp
	}
	if len(c.t_args) != len(g.psi.tFormals) {
		var b strings.Builder
		b.WriteString("Arity mismatch: type actuals=[")
		writeTypes(&b, c.t_args)
		b.WriteString("], formals=[")
		b.WriteString(g.psi.String())
		b.WriteString("]\n\t")
		b.WriteString(c.String())
		panic(b.String())
	}
	if len(c.args) != len(g.pDecls) {
		var b strings.Builder
		b.WriteString("Arity mismatch: args=[")
		writeExprs(&b, c.args)
		b.WriteString("], params=[")
		writeParamDecls(&b, g.pDecls)
		b.WriteString("]\n\t")
		b.WriteString(c.String())
		panic(b.String())
	}
	subs := make(map[TParam]Type) // CHECKME: applying this subs vs. adding to a new delta?
	for i := 0; i < len(c.t_args); i++ {
		subs[g.psi.tFormals[i].name] = c.t_args[i]
	}
	for i := 0; i < len(c.t_args); i++ {
		u := g.psi.tFormals[i].u_I.TSubs(subs)
		if !c.t_args[i].ImplsDelta(ds, delta, u) {
			panic("Type actual must implement type formal: actual=" +
				c.t_args[i].String() + ", param=" + u.String())
		}
	}
	for i := 0; i < len(c.args); i++ {
		// CHECKME: submission version's notation, (~\tau :> ~\rho)[subs], slightly unclear
		u_a := c.args[i].Typing(ds, delta, gamma, allowStupid)
		//.TSubs(subs)  // !!! submission version, subs also applied to ~tau, ..
		// ..falsely captures "repeat" var occurrences in recursive calls, ..
		// ..e.g., bad monomorph (Box) example.
		// The ~beta morally do not occur in ~tau, they only bind ~rho
		u_p := g.pDecls[i].u.TSubs(subs)
		if !u_a.ImplsDelta(ds, delta, u_p) {
			panic("Arg expr type must implement param type: arg=" + u_a.String() +
				", param=" + u_p.String() + "\n\t" + c.String())
		}
	}
	return g.u_ret.TSubs(subs) // subs necessary, c.psi info (i.e., bounds) will be "lost" after leaving this context
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
	writeTypes(&b, c.t_args)
	b.WriteString(")(")
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
	writeToGoTypes(&b, c.t_args)
	b.WriteString(")(")
	writeToGoExprs(&b, c.args)
	b.WriteString(")")
	return b.String()
}

/* Assert */

type Assert struct {
	e_I    FGGExpr
	u_cast Type
}

var _ FGGExpr = Assert{}

func (a Assert) GetExpr() FGGExpr { return a.e_I }
func (a Assert) GetType() Type    { return a.u_cast }

func (a Assert) Subs(subs map[Variable]FGGExpr) FGGExpr {
	return Assert{a.e_I.Subs(subs), a.u_cast}
}

func (a Assert) TSubs(subs map[TParam]Type) FGGExpr {
	return Assert{a.e_I.TSubs(subs), a.u_cast.TSubs(subs)}
}

func (a Assert) Eval(ds []Decl) (FGGExpr, string) {
	if !a.e_I.IsValue() {
		e, rule := a.e_I.Eval(ds)
		return Assert{e, a.u_cast}, rule
	}
	u_S := a.e_I.(StructLit).u_S
	if !IsStructType(ds, u_S) {
		panic("Non struct type found in struct lit: " + u_S.String())
	}
	if u_S.ImplsDelta(ds, make(map[TParam]Type), a.u_cast) { // Empty Delta -- not super clear in submission version
		return a.e_I, "Assert"
	}
	panic("Cannot reduce: " + a.String())
}

func (a Assert) Typing(ds []Decl, delta Delta, gamma Gamma, allowStupid bool) Type {
	u := a.e_I.Typing(ds, delta, gamma, allowStupid)
	if IsStructType(ds, u) {
		if allowStupid {
			return a.u_cast
		} else {
			panic("Expr must be an interface type (in a non-stupid context): found " +
				u.String() + " for\n\t" + a.String())
		}
	}
	// u is a TParam or an interface type TName
	if _, ok := a.u_cast.(TParam); ok || IsNamedIfaceType(ds, a.u_cast) {
		return a.u_cast // No further checks -- N.B., Robert said they are looking to refine this
	}
	// a.u is a struct type TName
	if a.u_cast.ImplsDelta(ds, delta, u) {
		return a.u_cast
	}
	panic("Struct type assertion must implement expr type: asserted=" +
		a.u_cast.String() + ", expr=" + u.String())
}

// From base.fgg
func (a Assert) IsValue() bool {
	return false
}

func (a Assert) String() string {
	var b strings.Builder
	b.WriteString(a.e_I.String())
	b.WriteString(".(")
	b.WriteString(a.u_cast.String())
	b.WriteString(")")
	return b.String()
}

func (a Assert) ToGoString() string {
	var b strings.Builder
	b.WriteString(a.e_I.ToGoString())
	b.WriteString(".(")
	b.WriteString(a.u_cast.ToGoString())
	b.WriteString(")")
	return b.String()
}

/* Aux, helpers */

func writeExprs(b *strings.Builder, es []FGGExpr) {
	if len(es) > 0 {
		b.WriteString(es[0].String())
		for _, v := range es[1:] {
			b.WriteString(", " + v.String())
		}
	}
}

func writeToGoExprs(b *strings.Builder, es []FGGExpr) {
	if len(es) > 0 {
		b.WriteString(es[0].ToGoString())
		for _, v := range es[1:] {
			b.WriteString(", " + v.ToGoString())
		}
	}
}
