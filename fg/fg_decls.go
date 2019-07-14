/*
 * This file contains defs for "concrete" syntax w.r.t. programs and decls.
 * Base (abstract) types, interfaces, etc. are in fg.go.
 */

package fg

import "reflect"
import "strings"

type FGProgram struct {
	ds []Decl
	e  Expr
}

var _ FGNode = FGProgram{}

func (p FGProgram) Ok() {
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
	p.e.Typing(p.ds, gamma)
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

type MDecl struct {
	recv ParamDecl
	m    Name
	ps   []ParamDecl
	t    Type
	e    Expr
}

var _ Decl = MDecl{}

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
	t := m.e.Typing(ds, env)
	if !t.Impls(ds, m.t) {
		panic("Method body type must implement declared return type: found=" +
			t.String() + ", expected=" + m.t.String() + "\n\t" + m.String())
	}
}

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
	b.WriteString(m.t.String())
	b.WriteString(" { return ")
	b.WriteString(m.e.String())
	b.WriteString(" }")
	return b.String()
}

// Cf. FieldDecl
type ParamDecl struct {
	x Name
	t Type
}

var _ FGNode = ParamDecl{}

func (p ParamDecl) String() string {
	return p.x + " " + p.t.String()
}

type STypeLit struct { // TODO: rename STypeLit
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
	t Type
}

var _ FGNode = FieldDecl{}

func (fd FieldDecl) String() string {
	return fd.f + " " + fd.t.String()
}

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
