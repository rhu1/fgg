/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fgr

import (
	"fmt"
	"strings"

	"github.com/rhu1/fgg/internal/fgg"
	"github.com/rhu1/fgg/internal/parser"
)

var _ = fmt.Errorf

/* "Exported" constructors for fgg (translation) */

func NewVariable(id Name) Variable                 { return Variable{id} }
func NewStructLit(t Type, es []FGRExpr) StructLit  { return StructLit{t, es} }
func NewSelect(e FGRExpr, f Name) Select           { return Select{e, f} }
func NewCall(e FGRExpr, m Name, es []FGRExpr) Call { return Call{e, m, es} }
func NewAssert(e FGRExpr, t Type) Assert           { return Assert{e, t} }
func NewSynthAssert(e FGRExpr, t Type) SynthAssert { return SynthAssert{e, t} }

/* Variable */

type Variable struct {
	name Name
}

var _ FGRExpr = Variable{}

func (x Variable) Subs(subs map[Variable]FGRExpr) FGRExpr {
	res, ok := subs[x]
	if !ok {
		panic("Unknown var: " + x.String())
	}
	return res
}

func (x Variable) Eval(ds []Decl) (FGRExpr, string) {
	panic("Cannot reduce free variable: " + x.name)
}

func (x Variable) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	res, ok := gamma[x.name]
	if !ok {
		panic("Var not in env: " + x.String())
	}
	return res
}

func (x Variable) DropSynthAsserts(ds []Decl) FGRExpr {
	return x
}

// From base.Expr
func (x Variable) IsValue() bool {
	return false
}

func (x Variable) CanEval(ds []Decl) bool {
	return false
}

func (x Variable) String() string {
	return x.name
}

func (x Variable) ToGoString(ds []Decl) string {
	return x.name
}

/* StructLit */

type StructLit struct {
	t_S   Type
	elems []FGRExpr
}

var _ FGRExpr = StructLit{}

func (s StructLit) GetType() Type       { return s.t_S }
func (s StructLit) GetElems() []FGRExpr { return s.elems }

func (s StructLit) Subs(subs map[Variable]FGRExpr) FGRExpr {
	es := make([]FGRExpr, len(s.elems))
	for i := 0; i < len(s.elems); i++ {
		es[i] = s.elems[i].Subs(subs)
	}
	return StructLit{s.t_S, es}
}

func (s StructLit) Eval(ds []Decl) (FGRExpr, string) {
	es := make([]FGRExpr, len(s.elems))
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

func (s StructLit) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
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

func (s StructLit) DropSynthAsserts(ds []Decl) FGRExpr {
	es := make([]FGRExpr, len(s.elems))
	for i := 0; i < len(s.elems); i++ {
		es[i] = s.elems[i].DropSynthAsserts(ds)
	}
	return StructLit{s.t_S, es}
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

func (s StructLit) CanEval(ds []Decl) bool {
	for _, v := range s.elems {
		if v.CanEval(ds) {
			return true
		} else if !v.IsValue() {
			return false
		}
	}
	return false
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

func (s StructLit) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString("main.")
	b.WriteString(s.t_S.String())
	b.WriteString("{")
	td := getTDecl(ds, s.t_S).(STypeLit)
	if len(s.elems) > 0 {
		b.WriteString(td.fDecls[0].name)
		b.WriteString(":")
		b.WriteString(s.elems[0].ToGoString(ds))
		for i, v := range s.elems[1:] {
			b.WriteString(", ")
			b.WriteString(td.fDecls[i+1].name)
			b.WriteString(":")
			b.WriteString(v.ToGoString(ds))
		}
	}
	b.WriteString("}")
	return b.String()
}

/* Select */

type Select struct {
	e_S   FGRExpr
	field Name
}

var _ FGRExpr = Select{}

func (s Select) GetExpr() FGRExpr { return s.e_S }
func (s Select) GetField() Name   { return s.field }

func (s Select) Subs(subs map[Variable]FGRExpr) FGRExpr {
	return Select{s.e_S.Subs(subs), s.field}
}

func (s Select) Eval(ds []Decl) (FGRExpr, string) {
	if !s.e_S.IsValue() {
		e, rule := s.e_S.Eval(ds)
		return Select{e.(FGRExpr), s.field}, rule
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

func (s Select) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
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

func (s Select) DropSynthAsserts(ds []Decl) FGRExpr {
	if s.CanEval(ds) { // !!!
		e2, _ := s.Eval(ds)
		if v, ok := e2.(TRep); ok { // !!! cf. nomono.fgg
			return v
		}
	}
	return Select{s.e_S.DropSynthAsserts(ds), s.field}
}

// From base.Expr
func (s Select) IsValue() bool {
	return false
}

func (s Select) CanEval(ds []Decl) bool {
	if s.e_S.CanEval(ds) {
		return true
	} else if !s.e_S.IsValue() {
		return false
	}
	for _, v := range fields(ds, s.e_S.(StructLit).t_S) { // N.B. "purely operational", no typing aspect
		if v.name == s.field {
			return true
		}
	}
	return false
}

func (s Select) String() string {
	var b strings.Builder
	b.WriteString(s.e_S.String())
	b.WriteString(".")
	b.WriteString(s.field)
	return b.String()
}

func (s Select) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString(s.e_S.ToGoString(ds))
	b.WriteString(".")
	b.WriteString(s.field)
	return b.String()
}

/* Call */

type Call struct {
	e_recv FGRExpr
	meth   Name
	args   []FGRExpr
}

var _ FGRExpr = Call{}

func (c Call) GetReceiver() FGRExpr { return c.e_recv }
func (c Call) GetMethod() Name      { return c.meth }
func (c Call) GetArgs() []FGRExpr   { return c.args }

func (c Call) Subs(subs map[Variable]FGRExpr) FGRExpr {
	e := c.e_recv.Subs(subs)
	args := make([]FGRExpr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].Subs(subs)
	}
	return Call{e, c.meth, args}
}

func (c Call) Eval(ds []Decl) (FGRExpr, string) {
	if !c.e_recv.IsValue() {
		e, rule := c.e_recv.Eval(ds)
		return Call{e.(FGRExpr), c.meth, c.args}, rule
	}
	args := make([]FGRExpr, len(c.args))
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
	subs := make(map[Variable]FGRExpr)
	subs[Variable{x0}] = c.e_recv
	for i := 0; i < len(xs); i++ {
		subs[Variable{xs[i]}] = c.args[i]
	}
	return e.Subs(subs), "Call" // N.B. single combined substitution map slightly different to R-Call
}

func (c Call) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
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

func (c Call) DropSynthAsserts(ds []Decl) FGRExpr {
	e := c.e_recv.DropSynthAsserts(ds)
	args := make([]FGRExpr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].DropSynthAsserts(ds)
	}
	return Call{e, c.meth, args}
}

// From base.Expr
func (c Call) IsValue() bool {
	return false
}

func (c Call) CanEval(ds []Decl) bool {
	if c.e_recv.CanEval(ds) {
		return true
	} else if !c.e_recv.IsValue() {
		return false
	}
	for _, v := range c.args {
		if v.CanEval(ds) {
			return true
		} else if !v.IsValue() {
			return false
		}
	}
	t_S := c.e_recv.(StructLit).t_S
	for _, d := range ds { // TODO: factor out GetMethDecl
		if md, ok := d.(MDecl); ok &&
			md.recv.t == t_S && md.name == c.meth { // i.e., Impls, Cf. typing, aux methods
			return len(md.pDecls) == len(c.args) // Needed?
		}
	}
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

func (c Call) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString(c.e_recv.ToGoString(ds))
	b.WriteString(".")
	b.WriteString(c.meth)
	b.WriteString("(")
	writeToGoExprs(ds, &b, c.args)
	b.WriteString(")")
	return b.String()
}

/* Assert */

type Assert struct {
	e_I    FGRExpr
	t_cast Type
}

var _ FGRExpr = Assert{}

func (a Assert) GetExpr() FGRExpr { return a.e_I }
func (a Assert) GetType() Type    { return a.t_cast }

func (a Assert) Subs(subs map[Variable]FGRExpr) FGRExpr {
	return Assert{a.e_I.Subs(subs), a.t_cast}
}

func (a Assert) Eval(ds []Decl) (FGRExpr, string) {
	if !a.e_I.IsValue() {
		e, rule := a.e_I.Eval(ds)
		return Assert{e.(FGRExpr), a.t_cast}, rule
	}
	t_S := a.e_I.(StructLit).t_S
	if !isStructType(ds, t_S) {
		panic("Non struct type found in struct lit: " + t_S.String())
	}
	if t_S.Impls(ds, a.t_cast) {
		return a.e_I, "Assert"
	}
	panic("Cannot reduce: " + a.String())
}

// Typing ...
func (a Assert) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	t := a.e_I.Typing(ds, gamma, allowStupid)
	if isStructType(ds, t) {
		if allowStupid {
			return a.t_cast
		} else {
			panic("Expr must be an interface type (in a non-stupid context): found " +
				t.String() + " for\n\t" + a.String())
		}
	}
	// t is an interface type
	if isInterfaceType(ds, a.t_cast) {
		return a.t_cast // No further checks -- N.B., Robert said they are looking to refine this
	}
	// a.t is a struct type
	if a.t_cast.Impls(ds, t) {
		return a.t_cast
	}
	panic("Struct type assertion must implement expr type: asserted=" +
		a.t_cast.String() + ", expr=" + t.String())
}

func (a Assert) DropSynthAsserts(ds []Decl) FGRExpr {
	return Assert{a.e_I.DropSynthAsserts(ds), a.t_cast}
}

// From base.Expr
func (a Assert) IsValue() bool {
	return false
}

func (a Assert) CanEval(ds []Decl) bool {
	if a.e_I.CanEval(ds) {
		return true
	} else if !a.e_I.IsValue() {
		return false
	}
	return a.e_I.(StructLit).t_S.Impls(ds, a.t_cast)
}

func (a Assert) String() string {
	var b strings.Builder
	b.WriteString(a.e_I.String())
	b.WriteString(".(")
	b.WriteString(a.t_cast.String())
	b.WriteString(")")
	return b.String()
}

func (a Assert) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString(a.e_I.ToGoString(ds))
	b.WriteString(".(main.")
	b.WriteString(a.t_cast.String())
	b.WriteString(")")
	return b.String()
}

/* Synth assert -- duplicated from Assert */

type SynthAssert struct {
	e_I    FGRExpr
	t_cast Type
}

var _ FGRExpr = SynthAssert{}

func (a SynthAssert) GetExpr() FGRExpr { return a.e_I }
func (a SynthAssert) GetType() Type    { return a.t_cast }

// Subs from fgr.FGRExpr
func (a SynthAssert) Subs(subs map[Variable]FGRExpr) FGRExpr {
	return SynthAssert{a.e_I.Subs(subs), a.t_cast}
}

// Eval from base.Expr
func (a SynthAssert) Eval(ds []Decl) (FGRExpr, string) {
	if !a.e_I.IsValue() {
		e, rule := a.e_I.Eval(ds)
		return SynthAssert{e.(FGRExpr), a.t_cast}, rule
	}
	t_S := a.e_I.(StructLit).t_S
	if !isStructType(ds, t_S) {
		panic("Non struct type found in struct lit: " + t_S.String())
	}
	if t_S.Impls(ds, a.t_cast) {
		return a.e_I, "SynthAssert"
	}
	panic("Cannot reduce: " + a.String())
}

// Typing from fgr.FGRExpr
func (a SynthAssert) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	t := a.e_I.Typing(ds, gamma, allowStupid)
	if isStructType(ds, t) {
		if allowStupid {
			return a.t_cast
		} else {
			panic("Expr must be an interface type (in a non-stupid context): found " +
				t.String() + " for\n\t" + a.String())
		}
	}
	// t is an interface type
	if isInterfaceType(ds, a.t_cast) {
		return a.t_cast // No further checks -- N.B., Robert said they are looking to refine this
	}
	// a.t is a struct type
	if a.t_cast.Impls(ds, t) {
		return a.t_cast
	}
	panic("Struct type assertion must implement expr type: asserted=" +
		a.t_cast.String() + ", expr=" + t.String())
}

func (a SynthAssert) DropSynthAsserts(ds []Decl) FGRExpr {
	return a.e_I.DropSynthAsserts(ds)
}

// IsValue from base.Expr
func (a SynthAssert) IsValue() bool {
	return false
}

// CanEval from base.Expr
func (a SynthAssert) CanEval(ds []Decl) bool {
	if a.e_I.CanEval(ds) {
		return true
	} else if !a.e_I.IsValue() {
		return false
	}
	return a.e_I.(StructLit).t_S.Impls(ds, a.t_cast)
}

func (a SynthAssert) String() string {
	var b strings.Builder
	b.WriteString(a.e_I.String())
	b.WriteString(".((")
	b.WriteString(a.t_cast.String())
	b.WriteString("))")
	return b.String()
}

// ToGoString from base.Expr
func (a SynthAssert) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString(a.e_I.ToGoString(ds))
	b.WriteString(".((main.")
	b.WriteString(a.t_cast.String())
	b.WriteString("))")
	return b.String()
}

/* Panic */

type Panic struct{}

var _ FGRExpr = Panic{}

func (p Panic) Subs(subs map[Variable]FGRExpr) FGRExpr {
	return p
}

func (p Panic) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	panic("TODO: " + p.String()) // !!! FIXME: allow any t
}

func (p Panic) Eval(ds []Decl) (FGRExpr, string) {
	panic("Cannot reduce panic.")
}

func (p Panic) DropSynthAsserts(ds []Decl) FGRExpr {
	return p
}

// From base.Expr
func (p Panic) IsValue() bool {
	return true
}

func (p Panic) CanEval(ds []Decl) bool {
	return false
}

func (p Panic) String() string {
	return "panic"
}

func (p Panic) ToGoString(ds []Decl) string {
	return "panic"
}

/* IfThenElse */

// IfThenElse represents type rep comparaisons
type IfThenElse struct {
	e1 FGRExpr // Cannot hardcode as Call, needs to be a general eval context
	e2 FGRExpr // TRep (or TmpTParam (Variable) for "wrappers")
	e3 FGRExpr
	//rho Map[fgg.Type]([]fgg.Sig)  // !!!
	src string // Original FGG source  // TODO store as a top-level comment or so?
}

var _ FGRExpr = IfThenElse{}

func (c IfThenElse) Subs(subs map[Variable]FGRExpr) FGRExpr {
	return IfThenElse{c.e1.Subs(subs), c.e2.Subs(subs),
		c.e3.Subs(subs), c.src}
}

func (c IfThenElse) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	if t1 := c.e1.Typing(ds, gamma, allowStupid); t1 != RepType {
		panic("IfThenElse comparison LHS must be of type " + RepType.String() +
			": found " + t1.String())
	}
	if t2 := c.e2.Typing(ds, gamma, allowStupid); t2 != RepType {
		panic("IfThenElse comparison RHS must be of type " + RepType.String() +
			": found " + t2.String())
	}
	t3 := c.e3.Typing(ds, gamma, allowStupid)
	// !!! no explicit e4 -- should always be panic? (panic typing is TODO)
	return t3
}

func (c IfThenElse) Eval(ds []Decl) (FGRExpr, string) {
	if !c.e1.IsValue() {
		e, rule := c.e1.Eval(ds)
		return IfThenElse{e.(FGRExpr), c.e2, c.e3, c.src}, rule
	}
	if !c.e2.IsValue() {
		e, rule := c.e2.Eval(ds)
		return IfThenElse{c.e1, e.(FGRExpr), c.e3, c.src}, rule
	}

	// TODO: refactor
	var a parser.FGGAdaptor
	p_fgg := a.Parse(true, c.src).(fgg.FGGProgram)
	ds_fgg := p_fgg.GetDecls()

	r1 := c.e1.(TRep)
	r2 := c.e2.(TRep)
	if r1.Reify().ImplsDelta(ds_fgg, make(fgg.Delta), r2.Reify()) {
		return c.e3, "If-true"
	} else {
		return Panic{}, "If-false"
	}
}

func (c IfThenElse) DropSynthAsserts(ds []Decl) FGRExpr {
	return IfThenElse{c.e1, c.e2, c.e3.DropSynthAsserts(ds), c.src}
}

// From base.Expr
func (c IfThenElse) IsValue() bool {
	return false
}

func (c IfThenElse) CanEval(ds []Decl) bool {
	if c.e1.CanEval(ds) {
		return true
	} else if !c.e1.IsValue() {
		return false
	}
	if c.e2.CanEval(ds) {
		return true
	} else if !c.e2.IsValue() {
		return false
	}
	return c.e3.CanEval(ds)
}

func (c IfThenElse) String() string {
	var b strings.Builder
	b.WriteString("(if ")
	b.WriteString(c.e1.String())
	b.WriteString(" << ")
	b.WriteString(c.e2.String())
	b.WriteString(" then ")
	b.WriteString(c.e3.String())
	b.WriteString(" else panic)") // !!! hardcoded else-panic
	return b.String()
}

func (c IfThenElse) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString("(if ")
	b.WriteString(c.e1.ToGoString(ds))
	b.WriteString(" << ")
	b.WriteString(c.e2.ToGoString(ds))
	b.WriteString(" then ")
	b.WriteString(c.e3.ToGoString(ds))
	b.WriteString(" else panic)") // !!! hardcoded else-panic
	return b.String()
}

/* Let */

// Let represents: let x = e1 in e2
type Let struct {
	x  Variable
	e1 FGRExpr
	e2 FGRExpr
}

var _ FGRExpr = Let{}

// GetDef public getter
func (e Let) GetDef() FGRExpr { return e.e1 }

// GetBody public getter
func (e Let) GetBody() FGRExpr { return e.e2 }

// Subs from fgr.FGRExpr
func (e Let) Subs(subs map[Variable]FGRExpr) FGRExpr {
	subs1 := make(map[Variable]FGRExpr)
	for k, v := range subs {
		subs1[k] = v
	}
	subs1[e.x] = e.x
	return Let{e.x, e.e1.Subs(subs1), e.e2.Subs(subs1)}
}

// Typing from fgr.FGRExpr
func (e Let) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	t1 := e.e1.Typing(ds, gamma, allowStupid)
	gamma1 := make(Gamma)
	for k, v := range gamma {
		gamma1[k] = v
	}
	gamma1[e.x.name] = t1
	return e.e2.Typing(ds, gamma1, allowStupid)
}

// Eval from base.Expr
func (e Let) Eval(ds []Decl) (FGRExpr, string) {
	if !e.e1.IsValue() {
		e11, rule := e.e1.Eval(ds)
		return Let{e.x, e11, e.e2}, rule
	}
	subs := make(map[Variable]FGRExpr)
	subs[e.x] = e.e1
	return e.e2.Subs(subs), "Let"
}

func (e Let) DropSynthAsserts(ds []Decl) FGRExpr {
	return Let{e.x, e.e1.DropSynthAsserts(ds), e.e2.DropSynthAsserts(ds)}
}

// IsValue from base.Expr
func (e Let) IsValue() bool {
	return false
}

// CanEval from base.Expr
func (e Let) CanEval(ds []Decl) bool {
	if e.e1.CanEval(ds) {
		return true
	} else if !e.e1.IsValue() {
		return false
	}
	// Here: e.e1.IsValue()
	return true
}

func (e Let) String() string {
	var b strings.Builder
	b.WriteString("let ")
	b.WriteString(e.x.String())
	b.WriteString("=")
	b.WriteString(e.e1.String())
	b.WriteString(" in ")
	b.WriteString(e.e2.String())
	return b.String()
}

// ToGoString from base.Expr
func (e Let) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString("let ")
	b.WriteString(e.x.ToGoString(ds))
	b.WriteString("=")
	b.WriteString(e.e1.ToGoString(ds))
	b.WriteString(" in ")
	b.WriteString(e.e2.ToGoString(ds))
	return b.String()
}

/* TRep -- the result of mkRep, i.e., an FGR expr/value (of type RepType) that represents a (parameterised) FGG type */

type TRep struct {
	t_name Name
	args   []FGRExpr // TRep or TmpTParam -- CHECKME: TmpTParam still needed? ("wrappers" only?)
	// CHECKME: factor out TArg?
}

var _ FGRExpr = TRep{}

// GetArgs public getter
func (r TRep) GetArgs() []FGRExpr { return r.args }

func (r TRep) Reify() fgg.TNamed {
	if !r.IsValue() {
		panic("Cannot refiy non-ground TRep: " + r.String())
	}
	us := make([]fgg.Type, len(r.args)) // All TName
	for i := 0; i < len(us); i++ {
		us[i] = r.args[i].(TRep).Reify() // CHECKME: guaranteed TRep?
	}
	return fgg.NewTName(r.t_name, us)
}

func (r TRep) Subs(subs map[Variable]FGRExpr) FGRExpr {
	es := make([]FGRExpr, len(r.args))
	for i := 0; i < len(es); i++ {
		es[i] = r.args[i].Subs(subs)
	}
	return TRep{r.t_name, es}
}

// !!! TRep evaluation contexts vs. reify aux?
func (r TRep) Eval(ds []Decl) (FGRExpr, string) {
	// Cf. StructLit.Eval
	es := make([]FGRExpr, len(r.args))
	done := false
	var rule string
	for i := 0; i < len(es); i++ {
		v := r.args[i]
		if !done && !v.IsValue() {
			v, rule = v.Eval(ds)
			done = true
		}
		es[i] = v
	}
	if done {
		return TRep{r.t_name, es}, rule
	} else {
		panic("Cannot reduce: " + r.String())
	}
}

func (r TRep) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	return RepType
}

func (r TRep) DropSynthAsserts(ds []Decl) FGRExpr {
	return r
}

// From base.Expr
func (r TRep) IsValue() bool {
	for _, v := range r.args {
		if !v.IsValue() {
			return false
		}
	}
	return true
}

func (r TRep) CanEval(ds []Decl) bool {
	for _, v := range r.args {
		if v.CanEval(ds) {
			return true
		} else if !v.IsValue() {
			return false
		}
	}
	return false
}

func (r TRep) String() string {
	var b strings.Builder
	b.WriteString(r.t_name)
	b.WriteString("[[")
	writeExprs(&b, r.args)
	b.WriteString("]]")
	return b.String()
}

func (r TRep) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString("main.")
	b.WriteString(r.t_name)
	b.WriteString("[[")
	writeToGoExprs(ds, &b, r.args)
	b.WriteString("]]")
	return b.String()
}

/* Intermediate TParam -- for WIP "wrappers" (fgr_translation), not oblit */

// Cf. Variable
type TmpTParam struct {
	id Name
}

var _ FGRExpr = TmpTParam{}

func (tmp TmpTParam) Subs(subs map[Variable]FGRExpr) FGRExpr {
	a := NewVariable(tmp.id) // FIXME -- should just make Variable earlier? -- or make a Disamb pass?
	if e, ok := subs[a]; ok {
		return e
	}
	return a // FIXME -- should not depend on calling Subs to disamb?
}

func (tmp TmpTParam) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	panic("TODO: " + tmp.String()) // CHECKME?
}

func (tmp TmpTParam) Eval(ds []Decl) (FGRExpr, string) {
	panic("Shouldn't get in here: " + tmp.String())
}

func (tmp TmpTParam) DropSynthAsserts(ds []Decl) FGRExpr {
	panic("Shouldn't get in here: " + tmp.String())
}

func (tmp TmpTParam) IsValue() bool {
	panic("Shouldn't get in here: " + tmp.String())
}

func (tmp TmpTParam) CanEval(ds []Decl) bool {
	panic("Shouldn't get in here: " + tmp.String())
}

func (tmp TmpTParam) String() string {
	return tmp.id
}

func (tmp TmpTParam) ToGoString(ds []Decl) string {
	return tmp.id
}

/* Aux, helpers */

func writeExprs(b *strings.Builder, es []FGRExpr) {
	if len(es) > 0 {
		b.WriteString(es[0].String())
		for _, v := range es[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
}

func writeToGoExprs(ds []Decl, b *strings.Builder, es []FGRExpr) {
	if len(es) > 0 {
		b.WriteString(es[0].ToGoString(ds))
		for _, v := range es[1:] {
			b.WriteString(", ")
			b.WriteString(v.ToGoString(ds))
		}
	}
}
