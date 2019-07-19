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
type TFormals struct {
	as []TParam
	us []Type
}

func (psi TFormals) String() string {
	var b strings.Builder
	if len(psi.as) > 0 {
		b.WriteString(psi.as[0].String() + " " + psi.us[0].String())
		for i := 1; i < len(psi.as); i++ {
			b.WriteString(", " + psi.as[i].String() + " " + psi.us[i].String())
		}
	}
	return b.String()
}

/* MDecl, ParamDecl */

type MDecl struct {
	recv ParamDecl
	m    Name // Refactor to embed Sig?
	psi  TFormals
	ps   []ParamDecl
	u    Type
	e    Expr
}

var _ Decl = MDecl{}

/*
func (m MDecl) ToSig() Sig {
	return Sig{m.m, m.ps, m.t}
}

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

func (m MDecl) GetName() Name {
	return m.m
}

func (m MDecl) String() string {
	var b strings.Builder
	b.WriteString("func (")
	b.WriteString(m.recv.String())
	b.WriteString(") ")
	b.WriteString(m.m)
	b.WriteString("(")
	if len(m.ps) > 0 {
		b.WriteString(m.ps[0].String())
		for _, v := range m.ps[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
	b.WriteString(") ")
	b.WriteString(m.u.String())
	b.WriteString(" { return ")
	b.WriteString(m.e.String())
	b.WriteString(" }")
	return b.String()
}

// Cf. FieldDecl
type ParamDecl struct {
	x Name
	u Type
}

var _ FGGNode = ParamDecl{}

func (p ParamDecl) String() string {
	return p.x + " " + p.u.String()
}

/* STypeLit, FieldDecl */

type STypeLit struct {
	t   Type
	psi TFormals
	fds []FieldDecl
}

var _ TDecl = STypeLit{}

func (s STypeLit) GetType() Type {
	return s.t
}

func (s STypeLit) GetName() Name {
	return Name(s.t.(TName).t)
}

func (s STypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(s.t.String())
	b.WriteString(" struct {")
	if len(s.fds) > 0 {
		b.WriteString(" ")
		b.WriteString(s.fds[0].String())
		for _, v := range s.fds[1:] {
			b.WriteString("; ")
			b.WriteString(v.String())
		}
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

/*type Sig struct {
	m  Name
	//TODO: TFormals
	ps []ParamDecl
	t  Type
}

var _ Spec = Sig{}

// !!! Sig in FG (also, Go spec) includes ~x, which breaks "impls"
func (s0 Sig) EqExceptVars(s Sig) bool {
	if len(s0.ps) != len(s.ps) {
		return false
	}
	for i := 0; i < len(s0.ps); i++ {
		if s0.ps[i].t != s.ps[i].t {
			return false
		}
	}
	return s0.m == s.m && s0.t == s.t
}

func (s Sig) GetSigs(_ []Decl) []Sig {
	return []Sig{s}
}

func (s Sig) String() string {
	var b strings.Builder
	b.WriteString(s.m)
	b.WriteString("(")
	if len(s.ps) > 0 {
		b.WriteString(s.ps[0].String())
		for _, v := range s.ps[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
	b.WriteString(") ")
	b.WriteString(s.t.String())
	return b.String()
}*/
