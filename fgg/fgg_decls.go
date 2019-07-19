/*
 * This file contains defs for "concrete" syntax w.r.t. programs and decls.
 * Base ("abstract") types, interfaces, etc. are in fgg.go.
 */

package fgg

import "fmt"
import "reflect"
import "strings"

var _ = fmt.Errorf
var _ = reflect.Append

/* Program */

type FGGProgram struct {
	ds []Decl
	e  Expr
}

var _ FGGNode = FGGProgram{}

/*func (p FGGProgram) Ok(allowStupid bool) {
	if !allowStupid { // Hack, to print only "top-level" programs (not during Eval)
		fmt.Println("[Warning] Type decl OK not checked yet (e.g., distinct type/field/method names, etc.)")
	}
	for _, v := range p.ds {
		switch c := v.(type) {
		case TDecl:
			// TODO: e.g., unique type names, unique field names, unique method names
			// N.B. omitted from submission version
		case MDecl:
			c.Ok(p.ds)
		default:
			panic("Unknown decl: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}
	var gamma Env
	p.e.Typing(p.ds, gamma, allowStupid)
}

// CHECKME: resulting FGGProgram is not parsed from source, OK? -- cf. Expr.Eval
// But doesn't affect FGPprogam.Ok() (i.e., Expr.Typing)
func (p FGGProgram) Eval() (FGGProgram, string) {
	e, rule := p.e.Eval(p.ds)
	return FGGProgram{p.ds, e}, rule
}

func (p FGGProgram) GetDecls() []Decl {
	return p.ds // Returns a copy?
}

func (p FGGProgram) GetExpr() Expr {
	return p.e
}*/

func (p FGGProgram) String() string {
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

/* Type formals */

// Pre: len(as) == len(us)
// Wrapper for []TFormal (cf. e.g., FieldDecl), only because of "(type ...)" syntax
type TFormals struct {
	tfs []TFormal
}

func (psi TFormals) String() string {
	var b strings.Builder
	b.WriteString("(type ") // Includes "(...)" -- cf. e.g., writeFieldDecls
	if len(psi.tfs) > 0 {
		b.WriteString(psi.tfs[0].String())
		for _, v := range psi.tfs[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
	b.WriteString(")")
	return b.String()
}

type TFormal struct {
	a TParam
	u Type
	// CHECKME: submission version, upper bound \tau_I is only "of the form t_I(~\tau)"? -- i.e., not \alpha?
	// ^If so, then can refine to TName
}

func (tf TFormal) String() string {
	return string(tf.a) + " " + tf.u.String()
}

/* STypeLit, FieldDecl */

type STypeLit struct {
	t   Name
	psi TFormals
	fds []FieldDecl
}

var _ TDecl = STypeLit{}

/*func (s STypeLit) GetType() Type {
	return s.t
}*/

func (s STypeLit) GetName() Name {
	return s.t
}

func (s STypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(string(s.t))
	b.WriteString(s.psi.String())
	b.WriteString(" struct {")
	if len(s.fds) > 0 {
		b.WriteString(" ")
		writeFieldDecls(&b, s.fds)
		b.WriteString(" ")
	}
	b.WriteString("}")
	return b.String()
}

type FieldDecl struct {
	f Name
	u Type
}

var _ FGGNode = FieldDecl{}

func (fd FieldDecl) Subs(subs map[TParam]Type) FieldDecl {
	return fd // FIXME TODO
}

func (fd FieldDecl) String() string {
	return fd.f + " " + fd.u.String()
}

/* MDecl, ParamDecl */

type MDecl struct {
	//recv ParamDecl
	x_recv   Name
	t_recv   Name // N.B. t_S
	psi_recv TFormals
	m        Name // Refactor to embed Sig?
	psi      TFormals
	pds      []ParamDecl
	u        Type
	e        Expr
}

var _ Decl = MDecl{}

func (md MDecl) ToSig() Sig {
	return Sig{md.m, md.psi, md.pds, md.u}
}

/*
func (m MDecl) Ok(ds []Decl) {
	if !isStructType(ds, m.recv.t) {
		panic("Receiver must be a struct type: not " + m.recv.t.String() +
			"\n\t" + m.String())
	}
	env := make(map[Name]Type)
	env[m.recv.x] = m.recv.t
	for _, v := range m.ps {
		env[v.x] = v.t
	}
	allowStupid := false
	t := m.e.Typing(ds, env, allowStupid)
	if !t.Impls(ds, m.t) {
		panic("Method body type must implement declared return type: found=" +
			t.String() + ", expected=" + m.t.String() + "\n\t" + m.String())
	}
}*/

func (md MDecl) GetName() Name {
	return md.m
}

func (md MDecl) String() string {
	var b strings.Builder
	b.WriteString("func (")
	//b.WriteString(md.recv.String())
	b.WriteString(md.x_recv)
	b.WriteString(" ")
	b.WriteString(md.t_recv)
	b.WriteString(md.psi_recv.String())
	b.WriteString(") ")
	b.WriteString(md.m)
	b.WriteString("(")
	writeParamDecls(&b, md.pds)
	b.WriteString(") ")
	b.WriteString(md.u.String())
	b.WriteString(" { return ")
	b.WriteString(md.e.String())
	b.WriteString(" }")
	return b.String()
}

// Cf. FieldDecl
type ParamDecl struct {
	x Name
	u Type
}

var _ FGGNode = ParamDecl{}

func (pd ParamDecl) String() string {
	return pd.x + " " + pd.u.String()
}

/*/* ITypeLit, Sig * /

type ITypeLit struct {
	t  Type // Factor out embedded struct with STypeLit?  But constructor will need that struct?
	ss []Spec
}

var _ TDecl = ITypeLit{}

func (r ITypeLit) GetType() Type {
	return r.t
}

func (r ITypeLit) GetName() Name {
	return Name(r.t)
}

func (r ITypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(r.t.String())
	b.WriteString(" interface {")
	if len(r.ss) > 0 {
		b.WriteString(" ")
		b.WriteString(r.ss[0].String())
		for _, v := range r.ss[1:] {
			b.WriteString("; ")
			b.WriteString(v.String())
		}
		b.WriteString(" ")
	}
	b.WriteString("}")
	return b.String()
}
*/

type Sig struct {
	m   Name
	psi TFormals
	pds []ParamDecl
	t   Type
}

var _ Spec = Sig{}

func (g Sig) Subs(subs map[TParam]Type) Sig {
	tfs := make([]TFormal, len(g.psi.tfs))
	for i := 0; i < len(g.psi.tfs); i++ {
		tf := g.psi.tfs[i]
		tfs[i] = TFormal{tf.a, tf.u.Subs(subs)}
	}
	ps := make([]ParamDecl, len(g.pds))
	for i := 0; i < len(ps); i++ {
		pd := g.pds[i]
		ps[i] = ParamDecl{pd.x, pd.u.Subs(subs)}
	}
	t := g.t.Subs(subs)
	return Sig{g.m, TFormals{tfs}, ps, t}
}

// !!! Sig in FGG includes ~a and ~x, which naively breaks "impls"
func (g0 Sig) EqExceptTParamsAndVars(g Sig) bool {
	if len(g0.psi.tfs) != len(g.psi.tfs) || len(g0.pds) != len(g.pds) {
		return false
	}
	for i := 0; i < len(g0.psi.tfs); i++ {
		if g0.psi.tfs[i].u != g.psi.tfs[i].u {
			return false
		}
	}
	for i := 0; i < len(g0.pds); i++ {
		if g0.pds[i].u != g.pds[i].u {
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
	b.WriteString(g.psi.String())
	b.WriteString("(")
	writeParamDecls(&b, g.pds)
	b.WriteString(") ")
	b.WriteString(g.t.String())
	return b.String()
}

/* Helpers */

// Doesn't include "(...)" -- slightly more convenient for debug messages
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
