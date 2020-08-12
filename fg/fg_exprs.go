/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fg

import "fmt"
import "strings"

/* "Exported" constructors for fgg (monomorph) */

func NewVariable(id Name) Variable                    { return Variable{id} }
func NewStructLit(t Type, es []FGExpr) StructLit      { return StructLit{t, es} }
func NewSelect(e FGExpr, f Name) Select               { return Select{e, f} }
func NewCall(e FGExpr, m Name, es []FGExpr) Call      { return Call{e, m, es} }
func NewAssert(e FGExpr, t Type) Assert               { return Assert{e, t} }
func NewString(v string) StringLit                    { return StringLit{v} }
func NewSprintf(format string, args []FGExpr) Sprintf { return Sprintf{format, args} }

/* Variable */

type Variable struct {
	name Name
}

var _ FGExpr = Variable{}

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

func (x Variable) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	res, ok := gamma[x.name]
	if !ok {
		panic("Var not in env: " + x.String())
	}
	return res
}

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
	elems []FGExpr
}

var _ FGExpr = StructLit{}

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

func (s StructLit) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	if !isTypeOk(ds, s.t_S) {
		panic("Unknown type: " + string(s.t_S) + "\n\t" + s.String())
	}
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
	for i, v := range s.elems {
		t := v.Typing(ds, gamma, allowStupid)
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
	e_S   FGExpr
	field Name
}

var _ FGExpr = Select{}

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
	panic("Field not found: " + s.field + "\n\t" + s.String())
}

func (s Select) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	t := s.e_S.Typing(ds, gamma, allowStupid)
	if !isStructType(ds, t) {
		panic("Illegal select on expr of non-struct type: " + string(t) +
			"\n\t" + s.String())
	}
	fds := fields(ds, t)
	for _, v := range fds {
		if v.name == s.field {
			return v.t
		}
	}
	panic("Field " + s.field + " not found in type: " + t.String() +
		"\n\t" + s.String())
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
	return s.e_S.String() + "." + s.field
}

func (s Select) ToGoString(ds []Decl) string {
	return s.e_S.ToGoString(ds) + "." + s.field
}

/* Call */

type Call struct {
	e_recv FGExpr
	meth   Name
	args   []FGExpr
}

var _ FGExpr = Call{}

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
		b.WriteString("]\n\t")
		b.WriteString(c.String())
		panic(b.String())
	}
	for i, a := range c.args {
		t := a.Typing(ds, gamma, allowStupid)
		if !t.Impls(ds, g.pDecls[i].t) {
			panic("Arg expr type must implement param type: arg=" + t.String() +
				", param=" + g.pDecls[i].t.String() + "\n\t" + c.String())
		}
	}
	return g.t_ret
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
		if md, ok := d.(MethDecl); ok &&
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
	e_I    FGExpr
	t_cast Type
}

var _ FGExpr = Assert{}

func (a Assert) GetExpr() FGExpr { return a.e_I }
func (a Assert) GetType() Type   { return a.t_cast }

func (a Assert) Subs(subs map[Variable]FGExpr) FGExpr {
	return Assert{a.e_I.Subs(subs), a.t_cast}
}

func (a Assert) Eval(ds []Decl) (FGExpr, string) {
	if !a.e_I.IsValue() {
		e, rule := a.e_I.Eval(ds)
		return Assert{e.(FGExpr), a.t_cast}, rule
	}
	t_S := a.e_I.(StructLit).t_S
	if !isStructType(ds, t_S) {
		panic("Non struct type found in struct lit: " + t_S)
	}
	if t_S.Impls(ds, a.t_cast) {
		return a.e_I, "Assert"
	}
	panic("Cannot reduce: " + a.String())
}

func (a Assert) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	t := a.e_I.Typing(ds, gamma, allowStupid)
	if !isTypeOk(ds, a.t_cast) {
		panic("Unknown type: " + string(a.t_cast) + "\n\t" + a.String())
	}
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

/* StringLit, fmt.Sprintf */

type StringLit struct {
	val string
}

var _ FGExpr = StringLit{}

func (s StringLit) GetValue() string { return s.val }

func (s StringLit) Subs(subs map[Variable]FGExpr) FGExpr {
	return s
}

func (s StringLit) Eval(ds []Decl) (FGExpr, string) {
	panic("Cannot reduce: " + s.String())
}

func (s StringLit) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	return STRING_TYPE
}

// From base.Expr
func (s StringLit) IsValue() bool {
	return true
}

func (s StringLit) CanEval(ds []Decl) bool {
	return false
}

func (s StringLit) String() string {
	return "\"" + s.val + "\""
}

func (s StringLit) ToGoString(ds []Decl) string {
	return "\"" + s.val + "\""
}

type Sprintf struct {
	format string // Includes surrounding quotes
	args   []FGExpr
}

var _ FGExpr = Sprintf{}

func (s Sprintf) GetFormat() string { return s.format }
func (s Sprintf) GetArgs() []FGExpr { return s.args }

func (s Sprintf) Subs(subs map[Variable]FGExpr) FGExpr {
	args := make([]FGExpr, len(s.args))
	for i := 0; i < len(args); i++ {
		args[i] = s.args[i].Subs(subs)
	}
	return Sprintf{s.format, args}
}

func (s Sprintf) Eval(ds []Decl) (FGExpr, string) {
	args := make([]FGExpr, len(s.args))
	done := false
	var rule string
	for i := 0; i < len(s.args); i++ {
		v := s.args[i]
		if !done && !v.IsValue() {
			v, rule = v.Eval(ds)
			done = true
		}
		args[i] = v
	}
	if done {
		return Sprintf{s.format, args}, rule
	} else {
		cast := make([]interface{}, len(args))
		for i := range args {
			cast[i] = args[i] // N.B. inside fgg this is, e.g., a StructLit (not the struct itself, as in native Go)
		}
		template := s.format[1 : len(s.format)-1] // Remove surrounding quote chars
		str := fmt.Sprintf(template, cast...)
		str = strings.ReplaceAll(str, "\"", "") // HACK because StringLit.String() includes quotes
		// FIXME: currently, user templates cannot include explicit quote chars
		return StringLit{str}, "Sprintf"
	}
}

// TODO: [Warning] not "fully" type checked, cf. MISSING/EXTRA
func (s Sprintf) Typing(ds []Decl, gamma Gamma, allowStupid bool) Type {
	for i := 0; i < len(s.args); i++ {
		s.args[i].Typing(ds, gamma, allowStupid)
	}
	return STRING_TYPE
}

// From base.Expr
func (s Sprintf) IsValue() bool {
	return false
}

func (s Sprintf) CanEval(ds []Decl) bool {
	return true
}

func (s Sprintf) String() string {
	var b strings.Builder
	b.WriteString("fmt.Sprintf(")
	b.WriteString(s.format)
	if len(s.args) > 0 {
		b.WriteString(", ")
		writeExprs(&b, s.args)
	}
	b.WriteString(")")
	return b.String()
}

func (s Sprintf) ToGoString(ds []Decl) string {
	var b strings.Builder
	b.WriteString("fmt.Sprintf(")
	b.WriteString(s.format)
	if len(s.args) > 0 {
		b.WriteString(", ")
		writeToGoExprs(ds, &b, s.args)
	}
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

func writeToGoExprs(ds []Decl, b *strings.Builder, es []FGExpr) {
	if len(es) > 0 {
		b.WriteString(es[0].ToGoString(ds))
		for _, v := range es[1:] {
			b.WriteString(", ")
			b.WriteString(v.ToGoString(ds))
		}
	}
}
