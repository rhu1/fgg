/*
 * This file contains defs for "concrete" syntax w.r.t. programs and decls.
 * Base ("abstract") types, interfaces, etc. are in fg.go.
 */

package fg

import "fmt"
import "reflect"
import "strings"

import "github.com/rhu1/fgg/base"

var _ = fmt.Errorf

/* "Exported" constructors (e.g., for fgg_monom)*/

func NewFGProgram(ds []Decl, e FGExpr, printf bool) FGProgram {
	return FGProgram{ds, e, printf}
}

func NewSTypeLit(t Type, fds []FieldDecl) STypeLit { return STypeLit{t, fds} }
func NewITypeLit(t Type, ss []Spec) ITypeLit       { return ITypeLit{t, ss} }
func NewMDecl(recv ParamDecl, m Name, pds []ParamDecl, t Type, e FGExpr) MethDecl {
	return MethDecl{recv, m, pds, t, e}
}
func NewFieldDecl(f Name, t Type) FieldDecl      { return FieldDecl{f, t} }
func NewParamDecl(x Name, t Type) ParamDecl      { return ParamDecl{x, t} } // For fgg_monom.MakeWMap
func NewSig(m Name, pds []ParamDecl, t Type) Sig { return Sig{m, pds, t} }  // For fgg_monom.MakeWMap

/* Program */

type FGProgram struct {
	decls  []Decl
	e_main FGExpr
	printf bool // false = "original" `_ = e_main` syntax; true = import-fmt/printf syntax
	// N.B. coincidentally "behaves" like an actual printf because interpreter prints out final eval result
}

var _ base.Program = FGProgram{}
var _ FGNode = FGProgram{}

// From base.Program
func (p FGProgram) GetDecls() []Decl   { return p.decls } // Return a copy?
func (p FGProgram) GetMain() base.Expr { return p.e_main }
func (p FGProgram) IsPrintf() bool     { return p.printf } // HACK

// From base.Program
func (p FGProgram) Ok(allowStupid bool) base.Type {
	tds := make(map[string]TDecl)    // Type name
	mds := make(map[string]MethDecl) // Hack, string = string(md.recv.t) + "." + md.name
	for _, v := range p.decls {
		switch d := v.(type) {
		case TDecl:
			d.Ok(p.decls) // Currently empty -- TODO: check, e.g., unique field names -- cf., above [Warning]
			// N.B. checks also omitted from submission version
			t := d.GetName()
			if _, ok := tds[t]; ok {
				panic("Multiple declarations of type name: " + t + "\n\t" +
					d.String())
			}
			tds[t] = d
		case MethDecl:
			d.Ok(p.decls)
			hash := string(d.recv.t) + "." + d.name
			if _, ok := mds[hash]; ok {
				panic("Multiple declarations for receiver " + string(d.recv.t) +
					" of the method name: " + d.name + "\n\t" + d.String())
			}
			mds[hash] = d
		default:
			panic("Unknown decl: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}
	var gamma Gamma // Empty env for main
	return p.e_main.Typing(p.decls, gamma, allowStupid)
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

type STypeLit struct {
	t_S    Type
	fDecls []FieldDecl
}

var _ TDecl = STypeLit{}

func (s STypeLit) GetType() Type              { return s.t_S }
func (s STypeLit) GetFieldDecls() []FieldDecl { return s.fDecls }

// From Decl
func (s STypeLit) GetName() Name { return Name(s.t_S) }

// From Decl
func (s STypeLit) Ok(ds []Decl) {
	fs := make(map[Name]FieldDecl)
	for _, v := range s.fDecls {
		if _, ok := fs[v.name]; ok {
			panic("Multiple fields with name: " + v.name + "\n\t" + s.String())
		}
		fs[v.name] = v
		if !isTypeOk(ds, v.t) {
			panic("Field " + v.name + " has an unknown type: " + string(v.t) +
				"\n\t" + s.String())
		}
	}
	if isRecursiveFieldType(ds, make(map[Type]Type), s.t_S) {
		panic("Invalid recursive struct type:\n\t" + s.String())
	}
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

// Rename FDecl?
type FieldDecl struct {
	name Name
	t    Type
}

var _ FGNode = FieldDecl{}

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

/* MethDecl, ParamDecl */

type MethDecl struct {
	recv   ParamDecl
	name   Name // Not embedding Sig because Sig doesn't take xs
	pDecls []ParamDecl
	t_ret  Type // Return
	e_body FGExpr
}

var _ Decl = MethDecl{}

func (md MethDecl) GetReceiver() ParamDecl     { return md.recv }
func (md MethDecl) GetName() Name              { return md.name }   // From Decl
func (md MethDecl) GetParamDecls() []ParamDecl { return md.pDecls } // Returns non-receiver params
func (md MethDecl) GetReturn() Type            { return md.t_ret }
func (md MethDecl) GetBody() FGExpr            { return md.e_body }

func (md MethDecl) Ok(ds []Decl) {
	if !isStructType(ds, md.recv.t) {
		panic("Receiver must be a struct type: not " + md.recv.t.String() +
			"\n\t" + md.String())
	}
	env := Gamma{md.recv.name: md.recv.t}
	for _, v := range md.pDecls {
		if !isTypeOk(ds, v.t) {
			panic("Parameter " + v.name + " has an unknown type: " + string(v.t) +
				"\n\t" + md.String())
		}
		if _, ok := env[v.name]; ok {
			panic("Multiple receiver/parameters with name " + v.name + "\n\t" +
				md.String())
		}
		env[v.name] = v.t
	}
	if !isTypeOk(ds, md.t_ret) {
		panic("Unknown return type: " + string(md.t_ret) + "\n\t" + md.String())
	}
	allowStupid := false
	t := md.e_body.Typing(ds, env, allowStupid)
	if !t.Impls(ds, md.t_ret) {
		panic("Method body type must implement declared return type: found=" +
			t.String() + ", expected=" + md.t_ret.String() + "\n\t" + md.String())
	}
}

func (md MethDecl) ToSig() Sig {
	return Sig{md.name, md.pDecls, md.t_ret}
}

func (md MethDecl) String() string {
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

// Cf. FieldDecl  // Rename PDecl?
type ParamDecl struct {
	name Name // CHECKME: Variable? (also Env key)
	t    Type
}

var _ FGNode = ParamDecl{}

func (pd ParamDecl) GetName() Name { return pd.name } // From Decl
func (pd ParamDecl) GetType() Type { return pd.t }

func (pd ParamDecl) String() string {
	var b strings.Builder
	b.WriteString(pd.name)
	b.WriteString(" ")
	b.WriteString(pd.t.String())
	return b.String()
}

/* ITypeLit, Sig */

type ITypeLit struct {
	t_I   Type // Factor out embedded struct with STypeLit?  But constructor will need that struct?
	specs []Spec
}

var _ TDecl = ITypeLit{}

func (c ITypeLit) GetType() Type    { return c.t_I }
func (c ITypeLit) GetSpecs() []Spec { return c.specs }

// From Decl
func (c ITypeLit) GetName() Name { return Name(c.t_I) }

// From Decl
func (c ITypeLit) Ok(ds []Decl) {
	seen := make(map[Name]Sig)
	for _, v := range c.specs {
		switch s := v.(type) {
		case Sig:
			if _, ok := seen[s.meth]; ok {
				panic("Multiple sigs with name: " + s.meth + "\n\t" + c.String())
			}
			seen[s.meth] = s
		case Type:
			if !isInterfaceType(ds, s) {
				panic("Embedded type must be an interface, not: " + string(s) +
					"\n\t" + s.String())
			}
			if isRecursiveInterfaceEmbedding(ds, make(map[Type]Type), s) {
				panic("Invalid recursive interface embedding type:\n\t" + c.String())
			}
		}
	}
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

type Sig struct {
	meth   Name
	pDecls []ParamDecl
	t_ret  Type
}

var _ Spec = Sig{}

func (g Sig) GetMethod() Name            { return g.meth }
func (g Sig) GetParamDecls() []ParamDecl { return g.pDecls }
func (g Sig) GetReturn() Type            { return g.t_ret }

func (g0 Sig) Ok(ds []Decl) {
	seen := make(map[Type]ParamDecl)
	for _, v := range g0.pDecls {
		if !isTypeOk(ds, v.t) {
			panic("Parameter " + v.name + " has an unknown type: " + string(v.t) +
				"\n\t" + g0.String())
		}
		if _, ok := seen[v.t]; ok {
			panic("Multiple parameters with same name: " + v.name +
				"\n\t" + g0.String())
		}
	}
	if !isTypeOk(ds, g0.t_ret) {
		panic("Unknown return type: " + string(g0.t_ret) +
			"\n\t" + g0.String())
	}
}

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
	return g0.meth == g.meth && g0.t_ret == g.t_ret
}

// From Spec
func (g Sig) GetSigs(_ []Decl) []Sig {
	return []Sig{g}
}

func (g Sig) String() string {
	var b strings.Builder
	b.WriteString(g.meth)
	b.WriteString("(")
	writeParamDecls(&b, g.pDecls)
	b.WriteString(") ")
	b.WriteString(g.t_ret.String())
	return b.String()
}

/* Helpers */

// N.B. returns bool, not implicit panic like other Ok's
func isTypeOk(ds []Decl, t Type) bool { // Cf. isStructType, etc.
	if t == STRING_TYPE {
		return true
	}
	for _, v := range ds {
		if d, ok := v.(STypeLit); ok && d.t_S == t {
			return true
		} else if d, ok := v.(ITypeLit); ok && d.t_I == t {
			return true
		}
	}
	return false
}

// Pre: isStruct(ds, t_S)
func isRecursiveFieldType(ds []Decl, seen map[Type]Type, t_S Type) bool {
	if _, ok := seen[t_S]; ok {
		return true
	}
	for _, v := range fields(ds, t_S) {
		if !isStructType(ds, v.t) {
			continue
		}
		seen1 := make(map[Type]Type)
		for k, v := range seen {
			seen1[k] = v
		}
		seen1[t_S] = t_S
		if isRecursiveFieldType(ds, seen1, v.t) {
			return true
		}
	}
	return false
}

// Pre: isNamedIfaceType(ds, t_I), t_I OK already checked
func isRecursiveInterfaceEmbedding(ds []Decl, seen map[Type]Type, t_I Type) bool {
	if _, ok := seen[t_I]; ok {
		return true
	}
	td := getTDecl(ds, t_I).(ITypeLit)
	for _, v := range td.specs {
		emb, ok := v.(Type)
		if !ok {
			continue
		}
		seen1 := make(map[Type]Type)
		for k, v := range seen {
			seen1[k] = v
		}
		seen1[t_I] = t_I
		if isRecursiveInterfaceEmbedding(ds, seen1, emb) {
			return true
		}
	}
	return false
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
		case MethDecl:
			// checks that (method-type, method-name) is unique
			// RH: CHECKME: this would allow (bad) "return overloading"? -- note, d.t is the method return type
			if md, ok := decl.(MethDecl); ok && d.t_ret.String() == md.t_ret.String() && d.GetName() == md.GetName() {
				count++
			}
		}
	}
	return count == 1
}

//*/
