/*
 * This file contains defs for "concrete" syntax w.r.t. programs and decls.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fgr

import "fmt"
import "reflect"
import "strings"

import "github.com/rhu1/fgg/base"
import "github.com/rhu1/fgg/fgg"

/* "Exported" constructors for fgg (monomorph) */

// TODO: compact
func NewFGRProgram(ds []Decl, e Expr) FGRProgram {
	return FGRProgram{ds, e}
}

func NewSTypeLit(t Type /*rds []RepDecl,*/, fds []FieldDecl) STypeLit {
	return STypeLit{t /*rds,*/, fds}
}

func NewFieldDecl(f Name, t Type) FieldDecl {
	return FieldDecl{f, t}
}

func NewMDecl(recv ParamDecl, m Name, rds []RepDecl, pds []ParamDecl, t Type,
	e Expr) MDecl {
	return MDecl{recv, m, rds, pds, t, e}
}

func NewParamDecl(x Name, t Type) ParamDecl { // For fgg_util.MakeWMap
	return ParamDecl{x, t}
}

func NewITypeLit(t Type, ss []Spec) ITypeLit {
	return ITypeLit{t, ss}
}

func NewSig(m Name, pds []ParamDecl, t Type) Sig { // For fgg_util.MakeWMap
	return Sig{m, pds, t}
}

func (g Sig) GetMethName() Name { // Hack
	return g.m
}

/* Program */

type FGRProgram struct {
	ds []Decl
	e  Expr
}

var _ base.Program = FGRProgram{}
var _ FGRNode = FGRProgram{}

func (p FGRProgram) Ok(allowStupid bool) {
	if !allowStupid { // Hack, to print the following only for "top-level" programs (not during Eval)
		fmt.Println("[Warning] Type/method decl OK not fully checked yet " +
			"(e.g., distinct field/param names, etc.)")
	}
	tds := make(map[Type]TDecl)
	mds := make(map[string]MDecl) // Hack, string = string(md.recv.t) + "." + md.GetName()
	for _, v := range p.ds {
		switch d := v.(type) {
		case TDecl:
			d.Ok(p.ds) // Currently empty -- TODO: check, e.g., unique field names -- cf., above [Warning]
			// N.B. checks also omitted from submission version
			t := Type(d.GetName())
			if _, ok := tds[t]; ok {
				panic("Multiple declarations of type name: " + string(t) + "\n\t" +
					d.String())
			}
			tds[t] = d
		case MDecl:
			d.Ok(p.ds)
			n := d.GetName()
			hash := string(d.recv.t) + "." + n
			if _, ok := mds[hash]; ok {
				panic("Multiple declarations for receiver " + string(d.recv.t) +
					" of the method name: " + n + "\n\t" + d.String())
			}
			mds[hash] = d
		default:
			panic("Unknown decl: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}
	var gamma Env // Empty env for main
	p.e.Typing(p.ds, gamma, allowStupid)
}

// CHECKME: resulting FGRProgram is not parsed from source, OK? -- cf. Expr.Eval
// But doesn't affect FGRPprogam.Ok() (i.e., Expr.Typing)
func (p FGRProgram) Eval() (base.Program, string) {
	e, rule := p.e.Eval(p.ds)
	return FGRProgram{p.ds, e.(Expr)}, rule
}

func (p FGRProgram) GetDecls() []Decl {
	return p.ds // Returns a copy?
}

func (p FGRProgram) GetExpr() base.Expr {
	return p.e
}

func (p FGRProgram) String() string {
	var b strings.Builder
	b.WriteString("package main;\n")
	for _, v := range p.ds {
		b.WriteString(v.String())
		b.WriteString(";\n")
	}
	b.WriteString("func main() { _ = ")
	b.WriteString(p.e.String())
	b.WriteString(" }")
	return b.String()
}

/* STypeLit, RepDecl, FieldDecl */

type STypeLit struct {
	t Type
	//rds []RepDecl
	fds []FieldDecl
}

var _ TDecl = STypeLit{}

func (s STypeLit) GetType() Type       { return s.t }
func (s STypeLit) GetName() Name       { return Name(s.t) }
func (s STypeLit) Fields() []FieldDecl { return s.fds }

func (s STypeLit) Ok(ds []Decl) {
	// TODO
}

func (s STypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(s.t.String())
	b.WriteString(" struct {")
	/*if len(s.rds) > 0 {
		b.WriteString(" ")
		writeRepDecls(&b, s.rds)
		if len(s.fds) > 0 {
			b.WriteString(";")
		}
	}*/
	if len(s.fds) > 0 {
		b.WriteString(" ")
		writeFieldDecls(&b, s.fds)
	}
	b.WriteString(" ")
	b.WriteString("}")
	return b.String()
}

type RepDecl struct {
	a fgg.TParam
	r Rep // TODO: Rep shouldn't be parameterised
}

var _ FGRNode = RepDecl{}

func (rd RepDecl) String() string {
	return rd.a.String() + " " + rd.r.String()
}

type FieldDecl struct {
	f Name
	t Type
}

func (f FieldDecl) GetName() Name { return f.f }
func (f FieldDecl) GetType() Type { return f.t }

var _ FGRNode = FieldDecl{}

func (fd FieldDecl) String() string {
	return fd.f + " " + fd.t.String()
}

/* MDecl, ParamDecl */

type MDecl struct {
	recv ParamDecl
	m    Name // Not embedding Sig because Sig doesn't take xs
	rds  []RepDecl
	pds  []ParamDecl
	t    Type // Return
	e    Expr
}

var _ Decl = MDecl{}

func (md MDecl) Receiver() ParamDecl    { return md.recv }
func (md MDecl) MethodName() Name       { return md.m }
func (md MDecl) GetRepDecls() []RepDecl { return md.rds }

// MethodParams returns the non-receiver parameters
func (md MDecl) MethodParams() []ParamDecl { return md.pds }
func (md MDecl) ReturnType() Type          { return md.t }
func (md MDecl) Impl() Expr                { return md.e }

func (md MDecl) ToSig() Sig {
	return Sig{md.m, md.pds, md.t}
}

func (md MDecl) Ok(ds []Decl) {
	if !isStructType(ds, md.recv.t) {
		panic("Receiver must be a struct type: not " + md.recv.t.String() +
			"\n\t" + md.String())
	}
	env := Env{md.recv.x: md.recv.t}
	// TODO: rds
	for _, v := range md.pds {
		env[v.x] = v.t
	}
	allowStupid := false
	t := md.e.Typing(ds, env, allowStupid)
	if !t.Impls(ds, md.t) {
		panic("Method body type must implement declared return type: found=" +
			t.String() + ", expected=" + md.t.String() + "\n\t" + md.String())
	}
}

func (md MDecl) GetName() Name {
	return md.m
}

func (md MDecl) String() string {
	var b strings.Builder
	b.WriteString("func (")
	b.WriteString(md.recv.String())
	b.WriteString(") ")
	b.WriteString(md.m)
	b.WriteString("(")
	writeRepDecls(&b, md.rds)
	if len(md.rds) > 0 && len(md.pds) > 0 {
		b.WriteString("; ")
	}
	writeParamDecls(&b, md.pds)
	b.WriteString(") ")
	b.WriteString(md.t.String())
	b.WriteString(" { return ")
	b.WriteString(md.e.String())
	b.WriteString(" }")
	return b.String()
}

// Cf. FieldDecl
type ParamDecl struct {
	x Name // CHECKME: Variable? (also Env key)
	t Type
}

func (pd ParamDecl) GetName() Name { return pd.x }
func (pd ParamDecl) GetType() Type { return pd.t }

var _ FGRNode = ParamDecl{}

func (pd ParamDecl) String() string {
	return pd.x + " " + pd.t.String()
}

/* ITypeLit, Sig */

type ITypeLit struct {
	t  Type // Factor out embedded struct with STypeLit?  But constructor will need that struct?
	ss []Spec
}

var _ TDecl = ITypeLit{}

func (c ITypeLit) Ok(ds []Decl) {
	// TODO
}

func (c ITypeLit) GetType() Type {
	return c.t
}

func (c ITypeLit) GetName() Name {
	return Name(c.t)
}

func (c ITypeLit) Specs() []Spec {
	return c.ss
}

func (c ITypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(c.t.String())
	b.WriteString(" interface {")
	if len(c.ss) > 0 {
		b.WriteString(" ")
		b.WriteString(c.ss[0].String())
		for _, v := range c.ss[1:] {
			b.WriteString("; ")
			b.WriteString(v.String())
		}
		b.WriteString(" ")
	}
	b.WriteString("}")
	return b.String()
}

type Sig struct {
	m   Name
	pds []ParamDecl
	t   Type
}

func (s Sig) MethodName() Name          { return s.m }
func (s Sig) MethodParams() []ParamDecl { return s.pds }
func (s Sig) ReturnType() Type          { return s.t }

var _ Spec = Sig{}

// !!! Sig in FGR (also, Go spec) includes ~x, which naively breaks "impls"
func (g0 Sig) EqExceptVars(g Sig) bool {
	if len(g0.pds) != len(g.pds) {
		return false
	}
	for i := 0; i < len(g0.pds); i++ {
		if g0.pds[i].t != g.pds[i].t {
			return false
		}
	}
	return g0.m == g.m && g0.t == g.t
}

func (g Sig) GetSigs(_ []Decl) []Sig {
	return []Sig{g}
}

func (g Sig) String() string {
	var b strings.Builder
	b.WriteString(g.m)
	b.WriteString("(")
	writeParamDecls(&b, g.pds)
	b.WriteString(") ")
	b.WriteString(g.t.String())
	return b.String()
}

/* Helpers */

// Doesn't include "(...)" -- slightly more convenient for debug messages
func writeRepDecls(b *strings.Builder, rds []RepDecl) {
	if len(rds) > 0 {
		b.WriteString(rds[0].String())
		for _, v := range rds[1:] {
			b.WriteString("; " + v.String())
		}
	}
}

func writeFieldDecls(b *strings.Builder, fds []FieldDecl) {
	if len(fds) > 0 {
		b.WriteString(fds[0].String())
		for _, v := range fds[1:] {
			b.WriteString("; " + v.String())
		}
	}
}

func writeParamDecls(b *strings.Builder, pds []ParamDecl) {
	if len(pds) > 0 {
		b.WriteString(pds[0].String())
		for _, v := range pds[1:] {
			b.WriteString(", " + v.String())
		}
	}
}

/* Old */

//*/

// RH: Possibly refactor aspects of this and related as "Decl.Wf()" -- the parts of "Ok()" omitted from the paper
func isDistinctDecl(decl Decl, ds []Decl) bool {
	var count int
	for _, d := range ds {
		switch d := d.(type) {
		case TDecl:
			// checks that type-name is unique regardless of definition
			// RH: Refactor as a single global pass (use a temp map), or into a TDecl.Wf() -- done: currently integrated into FGRProgram.Ok for now (to avoid a second iteration)
			if td, ok := decl.(TDecl); ok && d.GetName() == td.GetName() {
				count++
			}
		case MDecl:
			// checks that (method-type, method-name) is unique
			// RH: CHECKME: this would allow (bad) "return overloading"? -- note, d.t is the method return type
			if md, ok := decl.(MDecl); ok && d.t.String() == md.t.String() && d.GetName() == md.GetName() {
				count++
			}
		}
	}
	return count == 1
}

//*/
