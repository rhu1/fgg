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

func NewVariable(id Name) Variable             { return Variable{id} }
func NewStructLit(t Type, es []Expr) StructLit { return StructLit{t, es} }
func NewSelect(e Expr, f Name) Select          { return Select{e, f} }
func NewCall(e Expr, m Name, es []Expr) Call   { return Call{e, m, es} }
func NewAssert(e Expr, t Type) Assert          { return Assert{e, t} }

/* Variable */

type Variable struct {
	id Name
}

var _ Expr = Variable{}

func (x Variable) Subs(subs map[Variable]Expr) Expr {
	res, ok := subs[x]
	if !ok {
		panic("Unknown var: " + x.String())
	}
	return res
}

func (x Variable) Eval(ds []Decl) (Expr, string) {
	panic("Cannot reduce free variable: " + x.id)
}

func (x Variable) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
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
	t  Type
	es []Expr
}

func (s StructLit) Type() Type         { return s.t }
func (s StructLit) FieldExprs() []Expr { return s.es }

var _ Expr = StructLit{}

func (s StructLit) Subs(subs map[Variable]Expr) Expr {
	es := make([]Expr, len(s.es))
	for i := 0; i < len(s.es); i++ {
		es[i] = s.es[i].Subs(subs)
	}
	return StructLit{s.t, es}
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
		return StructLit{s.t, es}, rule
	} else {
		panic("Cannot reduce: " + s.String())
	}
}

func (s StructLit) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	fs := fields(ds, s.t)
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
		t := s.es[i].Typing(ds, gamma, allowStupid)
		u := fs[i].t
		if !t.Impls(ds, u) {
			panic("Arg expr must implement field type: arg=" + t.String() +
				", field=" + u.String() + "\n\t" + s.String())
		}
	}
	return s.t
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
	b.WriteString(s.t.String())
	b.WriteString("{")
	//b.WriteString(strings.Trim(strings.Join(strings.Split(fmt.Sprint(s.es), " "), ", "), "[]"))
	// ^ No: broken for nested structs
	writeExprs(&b, s.es)
	b.WriteString("}")
	return b.String()
}

func (s StructLit) ToGoString() string {
	var b strings.Builder
	b.WriteString("main.")
	b.WriteString(s.t.String())
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

func (s Select) Expr() Expr      { return s.e }
func (s Select) FieldName() Name { return s.f }

func (s Select) Subs(subs map[Variable]Expr) Expr {
	return Select{s.e.Subs(subs), s.f}
}

func (s Select) Eval(ds []Decl) (Expr, string) {
	if !s.e.IsValue() {
		e, rule := s.e.Eval(ds)
		return Select{e.(Expr), s.f}, rule
	}
	v := s.e.(StructLit)
	fds := fields(ds, v.t)
	for i := 0; i < len(fds); i++ {
		if fds[i].f == s.f {
			return v.es[i], "Select"
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
	panic("Field not found: " + s.f + " in " + t.String())
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
	e    Expr
	m    Name
	args []Expr
}

var _ Expr = Call{}

func (c Call) Expr() Expr       { return c.e }
func (c Call) MethodName() Name { return c.m }
func (c Call) Args() []Expr     { return c.args }

func (c Call) Subs(subs map[Variable]Expr) Expr {
	e := c.e.Subs(subs)
	args := make([]Expr, len(c.args))
	for i := 0; i < len(c.args); i++ {
		args[i] = c.args[i].Subs(subs)
	}
	return Call{e, c.m, args}
}

func (c Call) Eval(ds []Decl) (Expr, string) {
	if !c.e.IsValue() {
		e, rule := c.e.Eval(ds)
		return Call{e.(Expr), c.m, c.args}, rule
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
		return Call{c.e, c.m, args}, rule
	}
	// c.e and c.args all values
	s := c.e.(StructLit)
	x0, xs, e := body(ds, s.t, c.m) // panics if method not found
	subs := make(map[Variable]Expr)
	subs[Variable{x0}] = c.e
	for i := 0; i < len(xs); i++ {
		subs[Variable{xs[i]}] = c.args[i]
	}
	return e.Subs(subs), "Call" // N.B. single combined substitution map slightly different to R-Call
}

func (c Call) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	t0 := c.e.Typing(ds, gamma, allowStupid)
	var g Sig
	if tmp, ok := methods(ds, t0)[c.m]; !ok { // !!! submission version had "methods(m)"
		panic("Method not found: " + c.m + " in " + t0.String() + "\n\t" + c.String())
	} else {
		g = tmp
	}
	if len(c.args) != len(g.pds) {
		var b strings.Builder
		b.WriteString("Arity mismatch: args=[")
		writeExprs(&b, c.args)
		b.WriteString("], params=[")
		writeParamDecls(&b, g.pds)
		b.WriteString("]")
		panic(b.String())
	}
	for i := 0; i < len(c.args); i++ {
		t := c.args[i].Typing(ds, gamma, allowStupid)
		if !t.Impls(ds, g.pds[i].t) {
			panic("Arg expr type must implement param type: arg=" + t + ", param=" +
				g.pds[i].t)
		}
	}
	return g.t
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
	writeToGoExprs(&b, c.args)
	b.WriteString(")")
	return b.String()
}

/* Assert */

type Assert struct {
	e Expr
	t Type
}

var _ Expr = Assert{}

func (a Assert) Expr() Expr       { return a.e }
func (a Assert) AssertType() Type { return a.t }

func (a Assert) Subs(subs map[Variable]Expr) Expr {
	return Assert{a.e.Subs(subs), a.t}
}

func (a Assert) Eval(ds []Decl) (Expr, string) {
	if !a.e.IsValue() {
		e, rule := a.e.Eval(ds)
		return Assert{e.(Expr), a.t}, rule
	}
	t_S := typ(ds, a.e.(StructLit)) // panics if StructLit.t is not a t_S
	if t_S.Impls(ds, a.t) {
		return a.e, "Assert"
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

func (a Assert) IsValue() bool {
	return false
}

func (a Assert) String() string {
	var b strings.Builder
	b.WriteString(a.e.String())
	b.WriteString(".(")
	b.WriteString(a.t.String())
	b.WriteString(")")
	return b.String()
}

func (a Assert) ToGoString() string {
	var b strings.Builder
	b.WriteString(a.e.ToGoString())
	b.WriteString(".(main.")
	b.WriteString(a.t.String())
	b.WriteString(")")
	return b.String()
}

/* Panic */

type Panic struct{}

var _ Expr = Panic{}

func (p Panic) Subs(subs map[Variable]Expr) Expr {
	return p
}

func (p Panic) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	panic("TODO: " + p.String()) // FIXME: allow any t
}

func (p Panic) Eval(ds []Decl) (Expr, string) {
	panic("Cannot reduce panic.")
}

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
	e1 Expr // Cannot hardcode Call, needs to be a general eval context
	e2 Expr // TmpTParam (Variable) or TypeTree
	e3 Expr
	//rho Map[fgg.Type]([]fgg.Sig)
	src string // Original FGG source
}

var _ Expr = IfThenElse{}

func (c IfThenElse) Subs(subs map[Variable]Expr) Expr {
	return IfThenElse{c.e1.Subs(subs), c.e2.Subs(subs),
		c.e3.Subs(subs), c.src}
}

func (c IfThenElse) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	if t1 := c.e1.Typing(ds, gamma, allowStupid); t1 != TRep {
		panic("IfThenElse comparison LHS must be of type " + string(TRep) +
			": found " + t1.String())
	}
	if t2 := c.e2.Typing(ds, gamma, allowStupid); t2 != TRep {
		panic("IfThenElse comparison RHS must be of type " + string(TRep) +
			": found " + t2.String())
	}
	t3 := c.e3.Typing(ds, gamma, allowStupid)
	// !!! no explicit e4 -- should always be panic? (panic typing is TODO)
	return t3
}

func (c IfThenElse) Eval(ds []Decl) (Expr, string) {
	if !c.e1.IsValue() {
		e, rule := c.e1.Eval(ds)
		return IfThenElse{e.(Expr), c.e2, c.e3, c.src}, rule
	}
	if !c.e2.IsValue() {
		e, rule := c.e2.Eval(ds)
		return IfThenElse{c.e1, e.(Expr), c.e3, c.src}, rule
	}

	// TODO: refactor
	var a fgg.FGGAdaptor
	p_fgg := a.Parse(true, c.src).(fgg.FGGProgram)
	ds_fgg := p_fgg.GetDecls()

	tt1 := c.e1.(TypeTree)
	tt2 := c.e2.(TypeTree)
	if tt1.Reify().Impls(ds_fgg, make(fgg.TEnv), tt2.Reify()) {
		return c.e3, "If-true"
	} else {
		return Panic{}, "If-false"
	}
}

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

/* TypeTree -- the result of mkRep, i.e., run-time type rep value */

type TypeTree struct {
	t  Type
	es []Expr // TypeTree or TmpTParam -- CHECKME: TmpTParam still needed?
}

var _ Expr = TypeTree{}

func (tt TypeTree) Reify() fgg.TName {
	if !tt.IsValue() {
		panic("Cannot refiy non-ground TypeTree: " + tt.String())
	}
	us := make([]fgg.Type, len(tt.es)) // All TName
	for i := 0; i < len(us); i++ {
		us[i] = tt.es[i].(TypeTree).Reify() // CHECKME: guaranteed TypeTree?
	}
	return fgg.NewTName(string(tt.t), us)
}

func (tt TypeTree) Subs(subs map[Variable]Expr) Expr {
	es := make([]Expr, len(tt.es))
	for i := 0; i < len(es); i++ {
		es[i] = tt.es[i].Subs(subs)
	}
	return TypeTree{tt.t, es}
}

func (tt TypeTree) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	return TRep
}

// !!! TypeTree evaluation contexts vs. reify aux?
func (tt TypeTree) Eval(ds []Decl) (Expr, string) {
	// Cf. StructLit.Eval
	es := make([]Expr, len(tt.es))
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
		return TypeTree{tt.t, es}, rule
	} else {
		panic("Cannot reduce: " + tt.String())
	}
}

func (tt TypeTree) IsValue() bool {
	for i := 0; i < len(tt.es); i++ {
		if !tt.es[i].IsValue() {
			return false
		}
	}
	return true
}

func (tt TypeTree) String() string {
	var b strings.Builder
	b.WriteString(string(tt.t))
	b.WriteString("[[")
	writeExprs(&b, tt.es)
	b.WriteString("]]")
	return b.String()
}

func (tt TypeTree) ToGoString() string {
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

var _ Expr = TmpTParam{}

func (tmp TmpTParam) Subs(subs map[Variable]Expr) Expr {
	a := NewVariable(tmp.id) // FIXME -- should just make Variable earlier? -- or make a Disamb pass?
	if e, ok := subs[a]; ok {
		return e
	}
	return a // FIXME -- should not depend on calling Subs to disamb?
}

func (tmp TmpTParam) Typing(ds []Decl, gamma Env, allowStupid bool) Type {
	panic("TODO: " + tmp.String()) // CHECKME?
}

func (tmp TmpTParam) Eval(ds []Decl) (Expr, string) {
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

func writeExprs(b *strings.Builder, es []Expr) {
	if len(es) > 0 {
		b.WriteString(es[0].String())
		for _, v := range es[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
}

func writeToGoExprs(b *strings.Builder, es []Expr) {
	if len(es) > 0 {
		b.WriteString(es[0].ToGoString())
		for _, v := range es[1:] {
			b.WriteString(", ")
			b.WriteString(v.ToGoString())
		}
	}
}
