/*
 * This file contains defs for "concrete" syntax w.r.t. programs and decls.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fg

import "fmt"
import "reflect"
import "strings"

/* Program */

type FGProgram struct {
	ds []Decl
	e  Expr
}

var _ FGNode = FGProgram{}

func (p FGProgram) Ok(allowStupid bool) {
	if !allowStupid { // Hack, to print the following only for "top-level" programs (not during Eval)
		fmt.Println("[Warning] Type decl OK not checked yet " +
			"(e.g., distinct type/field/method names, etc.)")
	}
	for _, v := range p.ds {
		switch d := v.(type) {
		case TDecl:
			// TODO: Check, e.g., unique type/field/method names -- cf., above [Warning]
			// N.B. omitted from submission version
			// (call isDistinctDecl(d, p.ds))
		case MDecl:
			d.Ok(p.ds)
		default:
			panic("Unknown decl: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}
	var gamma Env // Empty env for main
	p.e.Typing(p.ds, gamma, allowStupid)
}

// Possibly refactor aspects of this and related as "Decl.Wf()" -- the parts of "Ok()" omitted from the paper
func isDistinctDecl(decl Decl, ds []Decl) bool {
	var count int
	for _, d := range ds {
		switch d := d.(type) {
		case TDecl:
			// checks that type-name is unique regardless of definition  // Refactor as a single global pass (use a temp map), or into a TDecl.Wf()
			if td, ok := decl.(TDecl); ok && d.GetName() == td.GetName() {
				count++
			}
		case MDecl:
			// checks that (method-type, method-name) is unique  // RH: CHECKME: this would allow (bad) "return overloading"?
			if md, ok := decl.(MDecl); ok && d.t.String() == md.t.String() && d.GetName() == md.GetName() {
				count++
			}
		}
	}
	return count == 1
}

// CHECKME: resulting FGProgram is not parsed from source, OK? -- cf. Expr.Eval
// But doesn't affect FGPprogam.Ok() (i.e., Expr.Typing)
func (p FGProgram) Eval() (FGProgram, string) {
	e, rule := p.e.Eval(p.ds)
	return FGProgram{p.ds, e}, rule
}

func (p FGProgram) GetDecls() []Decl {
	return p.ds // Returns a copy?
}

func (p FGProgram) GetExpr() Expr {
	return p.e
}

func (p FGProgram) String() string {
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

/* STypeLit, FieldDecl */

type STypeLit struct {
	t   Type
	fds []FieldDecl
}

var _ TDecl = STypeLit{}

func (s STypeLit) GetType() Type {
	return s.t
}

func (s STypeLit) GetName() Name {
	return Name(s.t)
}

func (s STypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(s.t.String())
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
	t Type
}

var _ FGNode = FieldDecl{}

func (fd FieldDecl) String() string {
	return fd.f + " " + fd.t.String()
}

/* MDecl, ParamDecl */

type MDecl struct {
	recv ParamDecl
	m    Name // Not embedding Sig because Sig doesn't take xs
	pds  []ParamDecl
	t    Type // Return
	e    Expr
}

var _ Decl = MDecl{}

func (md MDecl) ToSig() Sig {
	return Sig{md.m, md.pds, md.t}
}

func (md MDecl) Ok(ds []Decl) {
	if !isStructType(ds, md.recv.t) {
		panic("Receiver must be a struct type: not " + md.recv.t.String() +
			"\n\t" + md.String())
	}
	env := make(map[Name]Type)
	env[md.recv.x] = md.recv.t
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

var _ FGNode = ParamDecl{}

func (pd ParamDecl) String() string {
	return pd.x + " " + pd.t.String()
}

/* ITypeLit, Sig */

type ITypeLit struct {
	t  Type // Factor out embedded struct with STypeLit?  But constructor will need that struct?
	ss []Spec
}

var _ TDecl = ITypeLit{}

func (c ITypeLit) GetType() Type {
	return c.t
}

func (c ITypeLit) GetName() Name {
	return Name(c.t)
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

var _ Spec = Sig{}

// !!! Sig in FG (also, Go spec) includes ~x, which naively breaks "impls"
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
