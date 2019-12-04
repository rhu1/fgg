/*
 * This file contains defs for "concrete" syntax w.r.t. programs and decls.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fg

import "fmt"
import "reflect"
import "strings"

import "github.com/rhu1/fgg/base"

/* "Exported" constructors for fgg (monom) */

func NewFGProgram(ds []Decl, e FGExpr, printf bool) FGProgram {
	return FGProgram{ds, e, printf}
}

func NewSTypeLit(t Type, fds []FieldDecl) STypeLit { return STypeLit{t, fds} }
func NewITypeLit(t Type, ss []Spec) ITypeLit       { return ITypeLit{t, ss} }

func NewMDecl(recv ParamDecl, m Name, pds []ParamDecl, t Type, e FGExpr) MDecl {
	return MDecl{recv, m, pds, t, e}
}

func NewFieldDecl(f Name, t Type) FieldDecl      { return FieldDecl{f, t} }
func NewParamDecl(x Name, t Type) ParamDecl      { return ParamDecl{x, t} } // For fgg_monom.MakeWMap
func NewSig(m Name, pds []ParamDecl, t Type) Sig { return Sig{m, pds, t} }  // For fgg_monom.MakeWMap

/* Program */

var _ base.Program = FGProgram{}
var _ FGNode = FGProgram{}

type FGProgram struct {
	decls  []Decl
	e_main FGExpr
	printf bool // false = "original" `_ = e_main` syntax; true = import-fmt/printf syntax
}

// From base.Program
func (p FGProgram) GetDecls() []Decl   { return p.decls } // Return a copy?
func (p FGProgram) GetMain() base.Expr { return p.e_main }

// From base.Program
func (p FGProgram) Ok(allowStupid bool) {
	if !allowStupid { // Hack, to print the following only for "top-level" programs (not during Eval)
		fmt.Println("[Warning] Type/method decl OK not fully checked yet " +
			"(e.g., distinct field/param names, etc.)")
	}
	tds := make(map[Type]TDecl)
	mds := make(map[string]MDecl) // Hack, string = string(md.recv.t) + "." + md.GetName()
	for _, v := range p.decls {
		switch d := v.(type) {
		case TDecl:
			d.Ok(p.decls) // Currently empty -- TODO: check, e.g., unique field names -- cf., above [Warning]
			// N.B. checks also omitted from submission version
			t := Type(d.GetName())
			if _, ok := tds[t]; ok {
				panic("Multiple declarations of type name: " + string(t) + "\n\t" +
					d.String())
			}
			tds[t] = d
		case MDecl:
			d.Ok(p.decls)
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
	p.e_main.Typing(p.decls, gamma, allowStupid)
}

// CHECKME: resulting FGProgram is not parsed from source, OK? -- cf. Expr.Eval
// But doesn't affect FGPprogam.Ok() (i.e., Expr.Typing)
// From base.Program
func (p FGProgram) Eval() (base.Program, string) {
	e, rule := p.e_main.Eval(p.decls)
	return FGProgram{p.decls, e.(FGExpr), p.printf}, rule
}

func (p FGProgram) String() string {
	var b strings.Builder
	b.WriteString("package main;\n")
	if p.printf {
		b.WriteString("import \"fmt\";\n")
	}
	for _, v := range p.decls {
		b.WriteString(v.String())
		b.WriteString(";\n")
	}
	b.WriteString("func main() { ")
	if p.printf {
		b.WriteString("fmt.Printf(\"%#v\", ")
		b.WriteString(p.e_main.String())
		b.WriteString(")")
	} else {
		b.WriteString("_ = ")
		b.WriteString(p.e_main.String())
	}
	b.WriteString(" }")
	return b.String()
}

/* STypeLit, FieldDecl */

var _ TDecl = STypeLit{}

type STypeLit struct {
	t_S    Type
	fDecls []FieldDecl
}

func (s STypeLit) GetType() Type              { return s.t_S }
func (s STypeLit) GetFieldDecls() []FieldDecl { return s.fDecls }

// From Decl
func (s STypeLit) GetName() Name { return Name(s.t_S) }

// From Decl
func (s STypeLit) Ok(ds []Decl) {
	// TODO
}

func (s STypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(s.t_S.String())
	b.WriteString(" struct {")
	if len(s.fDecls) > 0 {
		b.WriteString(" ")
		writeFieldDecls(&b, s.fDecls)
		b.WriteString(" ")
	}
	b.WriteString("}")
	return b.String()
}

var _ FGNode = FieldDecl{}

type FieldDecl struct {
	name Name
	t    Type
}

func (f FieldDecl) GetType() Type { return f.t }

// From Decl
func (f FieldDecl) GetName() Name { return f.name }

func (fd FieldDecl) String() string {
	return fd.name + " " + fd.t.String()
	var b strings.Builder
	b.WriteString(fd.name)
	b.WriteString(" ")
	b.WriteString(fd.t.String())
	return b.String()
}

/* MDecl, ParamDecl */

var _ Decl = MDecl{}

type MDecl struct {
	recv   ParamDecl
	name   Name // Not embedding Sig because Sig doesn't take xs
	pDecls []ParamDecl
	t_ret  Type // Return
	e_body FGExpr
}

func (md MDecl) GetReceiver() ParamDecl     { return md.recv }
func (md MDecl) GetParamDecls() []ParamDecl { return md.pDecls } // Returns non-receiver params
func (md MDecl) GetReturn() Type            { return md.t_ret }
func (md MDecl) GetBody() FGExpr            { return md.e_body }

// From Decl
func (md MDecl) GetName() Name { return md.name }

func (md MDecl) Ok(ds []Decl) {
	if !isStructType(ds, md.recv.t) {
		panic("Receiver must be a struct type: not " + md.recv.t.String() +
			"\n\t" + md.String())
	}
	env := Env{md.recv.name: md.recv.t}
	for _, v := range md.pDecls {
		env[v.name] = v.t
	}
	allowStupid := false
	t := md.e_body.Typing(ds, env, allowStupid)
	if !t.Impls(ds, md.t_ret) {
		panic("Method body type must implement declared return type: found=" +
			t.String() + ", expected=" + md.t_ret.String() + "\n\t" + md.String())
	}
}

func (md MDecl) ToSig() Sig {
	return Sig{md.name, md.pDecls, md.t_ret}
}

func (md MDecl) String() string {
	var b strings.Builder
	b.WriteString("func (")
	b.WriteString(md.recv.String())
	b.WriteString(") ")
	b.WriteString(md.name)
	b.WriteString("(")
	writeParamDecls(&b, md.pDecls)
	b.WriteString(") ")
	b.WriteString(md.t_ret.String())
	b.WriteString(" { return ")
	b.WriteString(md.e_body.String())
	b.WriteString(" }")
	return b.String()
}

var _ FGNode = ParamDecl{}

// Cf. FieldDecl
type ParamDecl struct {
	name Name // CHECKME: Variable? (also Env key)
	t    Type
}

func (pd ParamDecl) GetType() Type { return pd.t }

// From Decl
func (pd ParamDecl) GetName() Name { return pd.name }

func (pd ParamDecl) String() string {
	var b strings.Builder
	b.WriteString(pd.name)
	b.WriteString(" ")
	b.WriteString(pd.t.String())
	return b.String()
}

/* ITypeLit, Sig */

var _ TDecl = ITypeLit{}

type ITypeLit struct {
	t_I   Type // Factor out embedded struct with STypeLit?  But constructor will need that struct?
	specs []Spec
}

func (c ITypeLit) GetType() Type    { return c.t_I }
func (c ITypeLit) GetSpecs() []Spec { return c.specs }

// From Decl
func (c ITypeLit) GetName() Name { return Name(c.t_I) }

// From Decl
func (c ITypeLit) Ok(ds []Decl) {
	// TODO
}

func (c ITypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(c.t_I.String())
	b.WriteString(" interface {")
	if len(c.specs) > 0 {
		b.WriteString(" ")
		b.WriteString(c.specs[0].String())
		for _, v := range c.specs[1:] {
			b.WriteString("; ")
			b.WriteString(v.String())
		}
		b.WriteString(" ")
	}
	b.WriteString("}")
	return b.String()
}

var _ Spec = Sig{}

type Sig struct {
	name   Name
	pDecls []ParamDecl
	t_ret  Type
}

func (s Sig) GetName() Name              { return s.name }
func (s Sig) GetParamDecls() []ParamDecl { return s.pDecls }
func (s Sig) GetReturn() Type            { return s.t_ret }

// !!! Sig in FG (also, Go spec) includes ~x, which naively breaks "impls"
func (g0 Sig) EqExceptVars(g Sig) bool {
	if len(g0.pDecls) != len(g.pDecls) {
		return false
	}
	for i := 0; i < len(g0.pDecls); i++ {
		if g0.pDecls[i].t != g.pDecls[i].t {
			return false
		}
	}
	return g0.name == g.name && g0.t_ret == g.t_ret
}

// From Spec
func (g Sig) GetSigs(_ []Decl) []Sig {
	return []Sig{g}
}

func (g Sig) String() string {
	var b strings.Builder
	b.WriteString(g.name)
	b.WriteString("(")
	writeParamDecls(&b, g.pDecls)
	b.WriteString(") ")
	b.WriteString(g.t_ret.String())
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

/* Old */

//*/

// RH: Possibly refactor aspects of this and related as "Decl.Wf()" -- the parts of "Ok()" omitted from the paper
func isDistinctDecl(decl Decl, ds []Decl) bool {
	var count int
	for _, d := range ds {
		switch d := d.(type) {
		case TDecl:
			// checks that type-name is unique regardless of definition
			// RH: Refactor as a single global pass (use a temp map), or into a TDecl.Wf() -- done: currently integrated into FGProgram.Ok for now (to avoid a second iteration)
			if td, ok := decl.(TDecl); ok && d.GetName() == td.GetName() {
				count++
			}
		case MDecl:
			// checks that (method-type, method-name) is unique
			// RH: CHECKME: this would allow (bad) "return overloading"? -- note, d.t is the method return type
			if md, ok := decl.(MDecl); ok && d.t_ret.String() == md.t_ret.String() && d.GetName() == md.GetName() {
				count++
			}
		}
	}
	return count == 1
}

//*/
