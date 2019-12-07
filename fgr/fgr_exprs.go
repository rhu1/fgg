/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fgr

import "fmt"
import "strings"

import "github.com/rhu1/fgg/fgg"

var _ = fmt.Errorf

/* "Exported" constructors for fgg (translation) */

func NewVariable(id Name) Variable                 { return Variable{id} }
func NewStructLit(t Type, es []FGRExpr) StructLit  { return StructLit{t, es} }
func NewSelect(e FGRExpr, f Name) Select           { return Select{e, f} }
func NewCall(e FGRExpr, m Name, es []FGRExpr) Call { return Call{e, m, es} }
func NewAssert(e FGRExpr, t Type) Assert           { return Assert{e, t} }

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

func (s Select) IsValue() bool {
	return false
}

func (s Select) String() string {
	var b strings.Builder
	b.WriteString(s.e_S.String())
	b.WriteString(".")
	b.WriteString(s.field)
	return b.String()
}

func (s Select) ToGoString() string {
	var b strings.Builder
	b.WriteString(s.e_S.ToGoString())
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

// From base.Expr
func (a Assert) IsValue() bool {
	return false
}

func (a Assert) String() string {
	var b strings.Builder
	b.WriteString(a.e_I.String())
	b.WriteString(".(")
	b.WriteString(a.t_cast.String())
	b.WriteString(")")
	return b.String()
}

func (a Assert) ToGoString() string {
	var b strings.Builder
	b.WriteString(a.e_I.ToGoString())
	b.WriteString(".(main.")
	b.WriteString(a.t_cast.String())
	b.WriteString(")")
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

// From base.Expr
func (p Panic) IsValue() bool {
	return true
}

func (p Panic) String() string {
	return "panic"
}

func (p Panic) ToGoString() string {
	return "panic"
}

/* IfThenElse */

type IfThenElse struct {
	e1 FGRExpr // Cannot hardcode as Call, needs to be a general eval context
	e2 FGRExpr // TmpTParam (Variable) or TypeTree
	e3 FGRExpr
	//rho Map[fgg.Type]([]fgg.Sig)  // !!!
	src string // Original FGG source
}

var _ FGRExpr = IfThenElse{}

func (c IfThenElse) Subs(subs map[Variable]FGRExpr) FGRExpr {
	return IfThenElse{c.e1.Subs(subs), c.e2.Subs(subs),
		c.e3.Subs(subs), c.src}
}

func (c IfThenElse) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	if t1 := c.e1.Typing(ds, gamma, allowStupid); t1 != FggType {
		panic("IfThenElse comparison LHS must be of type " + string(FggType) +
			": found " + t1.String())
	}
	if t2 := c.e2.Typing(ds, gamma, allowStupid); t2 != FggType {
		panic("IfThenElse comparison RHS must be of type " + string(FggType) +
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
	var a fgg.FGGAdaptor
	p_fgg := a.Parse(true, c.src).(fgg.FGGProgram)
	ds_fgg := p_fgg.GetDecls()

	tt1 := c.e1.(TRep)
	tt2 := c.e2.(TRep)
	if tt1.Reify().Impls(ds_fgg, make(fgg.Delta), tt2.Reify()) {
		return c.e3, "If-true"
	} else {
		return Panic{}, "If-false"
	}
}

// From base.Expr
func (c IfThenElse) IsValue() bool {
	return false
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

func (c IfThenElse) ToGoString() string {
	var b strings.Builder
	b.WriteString("(if ")
	b.WriteString(c.e1.ToGoString())
	b.WriteString(" << ")
	b.WriteString(c.e2.ToGoString())
	b.WriteString(" then ")
	b.WriteString(c.e3.ToGoString())
	b.WriteString(" else panic)") // !!! hardcoded else-panic
	return b.String()
}

/* TRep -- the result of mkRep, i.e., an FGR expr/value (of type FggType) that represents a (parameterised) FGG type */

type TRep struct {
	t  Type
	es []FGRExpr // TRep or TmpTParam -- CHECKME: TmpTParam still needed?
}

var _ FGRExpr = TRep{}

func (tt TRep) Reify() fgg.TNamed {
	if !tt.IsValue() {
		panic("Cannot refiy non-ground TypeTree: " + tt.String())
	}
	us := make([]fgg.Type, len(tt.es)) // All TName
	for i := 0; i < len(us); i++ {
		us[i] = tt.es[i].(TRep).Reify() // CHECKME: guaranteed TypeTree?
	}
	return fgg.NewTName(string(tt.t), us)
}

func (tt TRep) Subs(subs map[Variable]FGRExpr) FGRExpr {
	es := make([]FGRExpr, len(tt.es))
	for i := 0; i < len(es); i++ {
		es[i] = tt.es[i].Subs(subs)
	}
	return TRep{tt.t, es}
}

func (tt TRep) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	return FggType
}

// !!! TypeTree evaluation contexts vs. reify aux?
func (tt TRep) Eval(ds []Decl) (FGRExpr, string) {
	// Cf. StructLit.Eval
	es := make([]FGRExpr, len(tt.es))
	done := false
	var rule string
	for i := 0; i < len(es); i++ {
		v := tt.es[i]
		if !done && !v.IsValue() {
			v, rule = v.Eval(ds)
			done = true
		}
		es[i] = v
	}
	if done {
		return TRep{tt.t, es}, rule
	} else {
		panic("Cannot reduce: " + tt.String())
	}
}

func (tt TRep) IsValue() bool {
	for i := 0; i < len(tt.es); i++ {
		if !tt.es[i].IsValue() {
			return false
		}
	}
	return true
}

func (tt TRep) String() string {
	var b strings.Builder
	b.WriteString(string(tt.t))
	b.WriteString("[[")
	writeExprs(&b, tt.es)
	b.WriteString("]]")
	return b.String()
}

func (tt TRep) ToGoString() string {
	var b strings.Builder
	b.WriteString("main.")
	b.WriteString(string(tt.t))
	b.WriteString("[[")
	writeToGoExprs(&b, tt.es)
	b.WriteString("]]")
	return b.String()
}

/* Intermediate TParam */

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

func (tmp TmpTParam) IsValue() bool {
	panic("Shouldn't get in here: " + tmp.String())
}

func (tmp TmpTParam) String() string {
	return tmp.id
}

func (tmp TmpTParam) ToGoString() string {
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

func writeToGoExprs(b *strings.Builder, es []FGRExpr) {
	if len(es) > 0 {
		b.WriteString(es[0].ToGoString())
		for _, v := range es[1:] {
			b.WriteString(", ")
			b.WriteString(v.ToGoString())
		}
	}
}
