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

func NewVariable(id Name) Variable {
	return Variable{id}
}

/* Variable */

type Variable struct {
	id Name
}

var _ Expr = Variable{}

// TODO refactor
func (x Variable) GetName() Name {
	return x.id
}

func (x Variable) Subs(m map[Variable]Expr) Expr {
	res, ok := m[x]
	if !ok {
		panic("Unknown var: " + x.String())
	}
	return res
}

func (x Variable) TSubs(subs map[TParam]Type) Expr {
	return x
}

func (x Variable) Eval(ds []Decl) (Expr, string) {
	panic("Cannot evaluate free variable: " + x.id)
}

// TODO: refactor Typing and StupidTyping (clearer than bool param)
func (x Variable) Typing(ds []Decl, delta TEnv, gamma Env,
	allowStupid bool) Type {
	res, ok := gamma[x.id]
	if !ok {
		panic("Var not in env: " + x.String())
	}
	return res
}

func (x Variable) IsValue() bool {
	return false
}

func (x Variable) String() string {
	return x.id
}

func (x Variable) ToGoString() string {
	return x.id
}

/* StructLit */

type StructLit struct {
	u  TName // u.t is a t_S
	es []Expr
}

var _ Expr = StructLit{}

// TODO refactor
func (s StructLit) GetTName() TName {
	return s.u
}

// TODO refactor
func (s StructLit) GetArgs() []Expr {
	return s.es
}

func (s StructLit) Subs(subs map[Variable]Expr) Expr {
	es := make([]Expr, len(s.es))
	for i := 0; i < len(s.es); i++ {
		es[i] = s.es[i].Subs(subs)
	}
	return StructLit{s.u, es}
}

func (s StructLit) TSubs(subs map[TParam]Type) Expr {
	es := make([]Expr, len(s.es))
	for i := 0; i < len(s.es); i++ {
		es[i] = s.es[i].TSubs(subs)
	}
	return StructLit{s.u.TSubs(subs).(TName), es}
}

func (s StructLit) Eval(ds []Decl) (Expr, string) {
	es := make([]Expr, len(s.es))
	done := false
	var rule string
	for i := 0; i < len(s.es); i++ {
		v := s.es[i]
		if !done && !v.IsValue() {
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

func (s StructLit) IsValue() bool {
	for _, v := range s.es {
		if !v.IsValue() {
			return false
		}
	}
	return true
}

func (s StructLit) String() string {
	var b strings.Builder
	b.WriteString(s.u.String())
	b.WriteString("{")
	writeExprs(&b, s.es)
	b.WriteString("}")
	return b.String()
}

func (s StructLit) ToGoString() string {
	var b strings.Builder
	b.WriteString(s.u.ToGoString())
	b.WriteString("{")
	writeToGoExprs(&b, s.es)
	b.WriteString("}")
	return b.String()
}

/* Select */

type Select struct {
	e Expr
	f Name
}

var _ Expr = Select{}

// TODO refactor
func (s Select) GetExpr() Expr {
	return s.e
}

// TODO refactor
func (s Select) GetName() Name {
	return s.f
}

func (s Select) Subs(subs map[Variable]Expr) Expr {
	return Select{s.e.Subs(subs), s.f}
}

func (s Select) TSubs(subs map[TParam]Type) Expr {
	return Select{s.e.TSubs(subs), s.f}
}

func (s Select) Eval(ds []Decl) (Expr, string) {
	if !s.e.IsValue() {
		e, rule := s.e.Eval(ds)
		return Select{e, s.f}, rule
	}
	v := s.e.(StructLit)
	fds := fields(ds, v.u)
	for i := 0; i < len(fds); i++ {
		if fds[i].f == s.f {
			return v.es[i], "Select"
		}
	}
	panic("Field not found: " + s.f)
}

func (s Select) Typing(ds []Decl, delta TEnv, gamma Env,
	allowStupid bool) Type {
	u := s.e.Typing(ds, delta, gamma, allowStupid)
	if !isStructTName(ds, u) {
		panic("Illegal select on non-struct type expr: " + u.String())
	}
	fds := fields(ds, u.(TName))
	for _, v := range fds {
		if v.f == s.f {
			return v.u
		}
	}
	panic("Field not found: " + s.f + " in " + u.String())
}

func (s Select) IsValue() bool {
	return false
}

func (s Select) String() string {
	return s.e.String() + "." + s.f
}

func (s Select) ToGoString() string {
	return s.e.ToGoString() + "." + s.f
}

/* Call */

type Call struct {
	e     Expr
	m     Name
	targs []Type
	args  []Expr
}

var _ Expr = Call{}

// TODO refactor
func (c Call) GetRecv() Expr    { return c.e }
func (c Call) GetName() Name    { return c.m }
func (c Call) GetTArgs() []Type { return c.targs }
func (c Call) GetArgs() []Expr  { return c.args }

func (c Call) Subs(subs map[Variable]Expr) Expr {
	e := c.e.Subs(subs)
	args := make([]Expr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].Subs(subs)
	}
	return Call{e, c.m, c.targs, args}
}

func (c Call) TSubs(subs map[TParam]Type) Expr {
	targs := make([]Type, len(c.targs))
	for i := 0; i < len(c.targs); i++ {
		targs[i] = c.targs[i].TSubs(subs)
	}
	args := make([]Expr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].TSubs(subs)
	}
	return Call{c.e.TSubs(subs), c.m, targs, args}
}

func (c Call) Eval(ds []Decl) (Expr, string) {
	if !c.e.IsValue() {
		e, rule := c.e.Eval(ds)
		return Call{e, c.m, c.targs, c.args}, rule
	}
	args := make([]Expr, len(c.args))
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
		return Call{c.e, c.m, c.targs, args}, rule
	}
	// c.e and c.args all values
	s := c.e.(StructLit)
	x0, xs, e := body(ds, s.u, c.m, c.targs) // panics if method not found
	subs := make(map[Variable]Expr)
	subs[Variable{x0}] = c.e
	for i := 0; i < len(xs); i++ {
		subs[Variable{xs[i]}] = c.args[i]
	}
	return e.Subs(subs), "Call" // N.B. single combined substitution map slightly different to R-Call
}

func (c Call) Typing(ds []Decl, delta TEnv, gamma Env, allowStupid bool) Type {
	u0 := c.e.Typing(ds, delta, gamma, allowStupid)
	var g Sig
	if tmp, ok := methods(ds, bounds(delta, u0))[c.m]; !ok { // !!! submission version had "methods(m)"
		panic("Method not found: " + c.m + " in " + u0.String())
	} else {
		g = tmp
	}
	if len(c.targs) != len(g.psi.tfs) {
		var b strings.Builder
		b.WriteString("Arity mismatch: type actuals=[")
		writeTypes(&b, c.targs)
		b.WriteString("], formals=[")
		b.WriteString(g.psi.String())
		b.WriteString("]\n\t")
		b.WriteString(c.String())
		panic(b.String())
	}
	if len(c.args) != len(g.pds) {
		var b strings.Builder
		b.WriteString("Arity mismatch: args=[")
		writeExprs(&b, c.args)
		b.WriteString("], params=[")
		writeParamDecls(&b, g.pds)
		b.WriteString("]\n\t")
		b.WriteString(c.String())
		panic(b.String())
	}
	subs := make(map[TParam]Type) // CHECKME: applying this subs vs. adding to a new delta?
	for i := 0; i < len(c.targs); i++ {
		subs[g.psi.tfs[i].a] = c.targs[i]
	}
	for i := 0; i < len(c.targs); i++ {
		u := g.psi.tfs[i].u.TSubs(subs)
		if !c.targs[i].Impls(ds, delta, u) {
			panic("Type actual must implement type formal: actual=" +
				c.targs[i].String() + ", param=" + u.String())
		}
	}
	for i := 0; i < len(c.args); i++ {
		// CHECKME: submission version's notation, (~\tau :> ~\rho)[subs], slightly unclear
		u_a := c.args[i].Typing(ds, delta, gamma, allowStupid)
		//.TSubs(subs)  // !!! submission version, subs also applied to ~tau, ..
		// ..falsely captures "repeat" var occurrences in recursive calls, ..
		// ..e.g., bad monomorph (Box) example.
		// The ~beta morally do not occur in ~tau, they only bind ~rho
		u_p := g.pds[i].u.TSubs(subs)
		if !u_a.Impls(ds, delta, u_p) {
			panic("Arg expr type must implement param type: arg=" + u_a.String() +
				", param=" + u_p.String() + "\n\t" + c.String())
		}
	}
	return g.u.TSubs(subs) // subs necessary, c.psi info (i.e., bounds) will be "lost" after leaving this context
}

func (c Call) IsValue() bool {
	return false
}

func (c Call) String() string {
	var b strings.Builder
	b.WriteString(c.e.String())
	b.WriteString(".")
	b.WriteString(c.m)
	b.WriteString("(")
	writeTypes(&b, c.targs)
	b.WriteString(")(")
	writeExprs(&b, c.args)
	b.WriteString(")")
	return b.String()
}

func (c Call) ToGoString() string {
	var b strings.Builder
	b.WriteString(c.e.ToGoString())
	b.WriteString(".")
	b.WriteString(c.m)
	b.WriteString("(")
	writeToGoTypes(&b, c.targs)
	b.WriteString(")(")
	writeToGoExprs(&b, c.args)
	b.WriteString(")")
	return b.String()
}

/* Assert */

type Assert struct {
	e Expr
	u Type
}

func (a Assert) GetExpr() Expr { return a.e }
func (a Assert) GetType() Type { return a.u }

func (a Assert) Subs(subs map[Variable]Expr) Expr {
	return Assert{a.e.Subs(subs), a.u}
}

func (a Assert) TSubs(subs map[TParam]Type) Expr {
	return Assert{a.e.TSubs(subs), a.u.TSubs(subs)}
}

func (a Assert) Eval(ds []Decl) (Expr, string) {
	if !a.e.IsValue() {
		e, rule := a.e.Eval(ds)
		return Assert{e, a.u}, rule
	}
	u_S := typ(ds, a.e.(StructLit))                // panics if StructLit.u is not a TName u_S
	if u_S.Impls(ds, make(map[TParam]Type), a.u) { // Empty Delta -- not super clear in submission version
		return a.e, "Assert"
	}
	panic("Cannot reduce: " + a.String())
}

func (a Assert) Typing(ds []Decl, delta TEnv, gamma Env, allowStupid bool) Type {
	u := a.e.Typing(ds, delta, gamma, allowStupid)
	if isStructTName(ds, u) {
		if allowStupid {
			return a.u
		} else {
			panic("Expr must be an interface type (in a non-stupid context): found " +
				u.String() + " for\n\t" + a.String())
		}
	}
	// u is a TParam or an interface type TName
	if _, ok := a.u.(TParam); ok || isInterfaceTName(ds, a.u) {
		return a.u // No further checks -- N.B., Robert said they are looking to refine this
	}
	// a.u is a struct type TName
	if a.u.Impls(ds, delta, u) {
		return a.u
	}
	panic("Struct type assertion must implement expr type: asserted=" +
		a.u.String() + ", expr=" + u.String())
}

func (a Assert) IsValue() bool {
	return false
}

func (a Assert) String() string {
	var b strings.Builder
	b.WriteString(a.e.String())
	b.WriteString(".(")
	b.WriteString(a.u.String())
	b.WriteString(")")
	return b.String()
}

func (a Assert) ToGoString() string {
	var b strings.Builder
	b.WriteString(a.e.ToGoString())
	b.WriteString(".(")
	b.WriteString(a.u.ToGoString())
	b.WriteString(")")
	return b.String()
}

/* Aux, helpers */

func writeExprs(b *strings.Builder, es []Expr) {
	if len(es) > 0 {
		b.WriteString(es[0].String())
		for _, v := range es[1:] {
			b.WriteString(", " + v.String())
		}
	}
}

func writeToGoExprs(b *strings.Builder, es []Expr) {
	if len(es) > 0 {
		b.WriteString(es[0].ToGoString())
		for _, v := range es[1:] {
			b.WriteString(", " + v.ToGoString())
		}
	}
}
