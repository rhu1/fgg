/*
 * This file contains defs for "concrete" syntax w.r.t. programs and decls.
 * Base ("abstract") types, interfaces, etc. are in fgg.go.
 */

package fgg

import "fmt"
import "reflect"
import "strings"

import "github.com/rhu1/fgg/base"

var _ = fmt.Errorf
var _ = reflect.Append

/* Program */

type FGGProgram struct {
	ds []Decl
	e  Expr
}

var _ base.Program = FGGProgram{}
var _ FGGNode = FGGProgram{}

func (p FGGProgram) Ok(allowStupid bool) {
	if !allowStupid { // Hack, to print only "top-level" programs (not during Eval)
		fmt.Println("[Warning] Type lit OK (\"T ok\") not fully implemented yet " +
			"(e.g., distinct type/field/method names, etc.)")
	}
	for _, v := range p.ds {
		switch d := v.(type) {
		case TDecl:
			d.Ok(p.ds)
		case MDecl:
			d.Ok(p.ds)
		default:
			panic("Unknown decl: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}
	// Empty envs for main
	var delta TEnv
	var gamma Env
	p.e.Typing(p.ds, delta, gamma, allowStupid)
}

func (p FGGProgram) Eval() (base.Program, string) {
	e, rule := p.e.Eval(p.ds)
	return FGGProgram{p.ds, e.(Expr)}, rule
}

func (p FGGProgram) GetDecls() []Decl {
	return p.ds // Returns a copy?
}

func (p FGGProgram) GetExpr() base.Expr {
	return p.e
}

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

func (psi TFormals) Ok(ds []Decl) {
	for _, v := range psi.tfs {
		u, ok := v.u.(TName)
		if !ok {
			panic("Upper bound must be of the form \"t_I(type ...)\": not " + v.u.String())
		}
		if !isInterfaceTName(ds, u) { // CHECKME: subsumes above TName check (looks for \tau_S)
			panic("Upper bound must be an interface type: not " + u.String())
		}
	}
}

func (psi TFormals) ToTEnv() TEnv {
	delta := make(map[TParam]Type)
	for _, v := range psi.tfs {
		delta[v.a] = v.u
	}
	return delta
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

func (s STypeLit) Ok(ds []Decl) {
	TDeclOk(ds, s)
}

func (s STypeLit) GetName() Name {
	return s.t
}

func (s STypeLit) GetTFormals() TFormals {
	return s.psi
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
	return FieldDecl{fd.f, fd.u.TSubs(subs)}
}

func (fd FieldDecl) String() string {
	return fd.f + " " + fd.u.String()
}

/* MDecl, ParamDecl */

type MDecl struct {
	x_recv   Name // CHECKME: better to be Variable?  (etc. for other such Names)
	t_recv   Name // N.B. t_S
	psi_recv TFormals
	// N.B. receiver elements "decomposed" because TFormals (not TName, cf. fg.MDecl uses ParamDecl)
	m   Name // Refactor to embed Sig?
	psi TFormals
	pds []ParamDecl
	u   Type // Return
	e   Expr
}

var _ Decl = MDecl{}

func (md MDecl) ToSig() Sig {
	return Sig{md.m, md.psi, md.pds, md.u}
}

func (md MDecl) Ok(ds []Decl) {
	if !isStructType(ds, md.t_recv) {
		panic("Receiver must be a struct type: not " + md.t_recv +
			"\n\t" + md.String())
	}
	md.psi_recv.Ok(ds)
	md.psi.Ok(ds)

	delta := md.psi_recv.ToTEnv()
	for _, v := range md.psi_recv.tfs {
		v.u.Ok(ds, delta)
	}

	delta1 := md.psi.ToTEnv()
	for k, v := range delta {
		delta1[k] = v
	}
	for _, v := range md.psi.tfs {
		v.u.Ok(ds, delta1)
	}

	as := make([]Type, len(md.psi_recv.tfs)) // !!! submission version, x:t_S(a) => x:t_S(~a)
	for i := 0; i < len(md.psi_recv.tfs); i++ {
		as[i] = md.psi_recv.tfs[i].a
	}
	gamma := Env{md.x_recv: TName{md.t_recv, as}} // CHECKME: can we give the bounds directly here instead of 'as'?
	for _, v := range md.pds {
		gamma[v.x] = v.u
	}
	allowStupid := false
	u := md.e.Typing(ds, delta1, gamma, allowStupid)
	if !u.Impls(ds, delta1, md.u) {
		panic("Method body type must implement declared return type: found=" +
			u.String() + ", expected=" + md.u.String() + "\n\t" + md.String())
	}
}

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
	b.WriteString(md.psi.String())
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
	x Name // CHECKME: Variable?
	u Type
}

var _ FGGNode = ParamDecl{}

func (pd ParamDecl) String() string {
	return pd.x + " " + pd.u.String()
}

/* ITypeLit, Sig */

type ITypeLit struct {
	t   Name
	psi TFormals
	ss  []Spec
}

var _ TDecl = ITypeLit{}

func (c ITypeLit) Ok(ds []Decl) {
	TDeclOk(ds, c)
	for _, v := range c.ss {
		// TODO: check Sigs OK?  e.g., "type IA(type ) interface { m1(type )() Any };" while missing Any
		if g, ok := v.(Sig); ok {
			g.Ok(ds)
		}
	}
	// In general, also missing checks for, e.g., unique type/field/method names -- cf., TDeclOk
}

func (c ITypeLit) GetName() Name {
	return c.t
}

func (c ITypeLit) GetTFormals() TFormals {
	return c.psi
}

func (c ITypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(string(c.t))
	b.WriteString(c.psi.String())
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
	psi TFormals
	pds []ParamDecl
	u   Type
}

var _ Spec = Sig{}

// TODO: rename TSubs
func (g Sig) TSubs(subs map[TParam]Type) Sig {
	tfs := make([]TFormal, len(g.psi.tfs))
	for i := 0; i < len(g.psi.tfs); i++ {
		tf := g.psi.tfs[i]
		tfs[i] = TFormal{tf.a, tf.u.TSubs(subs)}
	}
	ps := make([]ParamDecl, len(g.pds))
	for i := 0; i < len(ps); i++ {
		pd := g.pds[i]
		ps[i] = ParamDecl{pd.x, pd.u.TSubs(subs)}
	}
	u := g.u.TSubs(subs)
	return Sig{g.m, TFormals{tfs}, ps, u}
}

func (g Sig) Ok(ds []Decl) {
	g.psi.Ok(ds)
	// TODO: check distinct param names, etc. -- N.B. interface may not be *used* (so may not be checked else where)
}

func (g Sig) GetSigs(_ []Decl) []Sig {
	return []Sig{g}
}

// !!! Sig in FGG includes ~a and ~x, which naively breaks "impls"
func (g0 Sig) EqExceptTParamsAndVars(g Sig) bool {
	if len(g0.psi.tfs) != len(g.psi.tfs) || len(g0.pds) != len(g.pds) {
		return false
	}
	for i := 0; i < len(g0.psi.tfs); i++ {
		if !g0.psi.tfs[i].u.Equals(g.psi.tfs[i].u) {
			return false
		}
	}
	for i := 0; i < len(g0.pds); i++ {
		if !g0.pds[i].u.Equals(g.pds[i].u) {
			return false
		}
	}
	return g0.m == g.m && g0.u.Equals(g.u)
}

func (g Sig) String() string {
	var b strings.Builder
	b.WriteString(g.m)
	b.WriteString(g.psi.String())
	b.WriteString("(")
	writeParamDecls(&b, g.pds)
	b.WriteString(") ")
	b.WriteString(g.u.String())
	return b.String()
}

/* Aux, helpers */

func TDeclOk(ds []Decl, td TDecl) {
	psi := td.GetTFormals()
	psi.Ok(ds)
	delta := psi.ToTEnv()
	for _, v := range psi.tfs {
		u, _ := v.u.(TName) // \tau_I, checked by psi.Ok
		u.Ok(ds, delta)     // !!! Submission version T-Type, t_i => t_I
	}
	// TODO: Check, e.g., unique type/field/method names -- cf., FGGProgram.OK [Warning]
}

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
