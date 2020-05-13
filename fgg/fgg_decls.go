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

/* Public constructors */

func NewProgram(ds []Decl, e FGGExpr, printf bool) FGGProgram {
	return FGGProgram{ds, e, printf}
}

/* Program */

type FGGProgram struct {
	decls  []Decl
	e_main FGGExpr
	printf bool // false = "original" `_ = e_main` syntax; true = import-fmt/printf syntax
	// N.B. coincidentally "behaves" like an actual printf because interpreter prints out final eval result
}

var _ base.Program = FGGProgram{}
var _ FGGNode = FGGProgram{}

func (p FGGProgram) GetDecls() []Decl   { return p.decls } // Return a copy?
func (p FGGProgram) GetMain() base.Expr { return p.e_main }
func (p FGGProgram) IsPrintf() bool     { return p.printf } // HACK

func (p FGGProgram) Ok(allowStupid bool) base.Type {
	tds := make(map[string]TypeDecl) // Type name
	mds := make(map[string]MethDecl) // Hack, string = md.recv.t + "." + md.name
	for _, v := range p.decls {
		switch d := v.(type) {
		case TypeDecl:
			d.Ok(p.decls)
			t := d.GetName()
			if _, ok := tds[t]; ok {
				panic("Multiple declarations of type name: " + t + "\n\t" +
					d.String())
			}
			tds[t] = d
		case MethDecl:
			d.Ok(p.decls)
			hash := string(d.t_recv) + "." + d.name
			if _, ok := mds[hash]; ok {
				panic("Multiple declarations for receiver " + string(d.t_recv) +
					" of the method name: " + d.name + "\n\t" + d.String())
			}
			mds[hash] = d
		default:
			panic("Unknown decl: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}
	// Empty envs for main
	var delta Delta
	var gamma Gamma
	return p.e_main.Typing(p.decls, delta, gamma, allowStupid)
}

func (p FGGProgram) Eval() (base.Program, string) {
	e, rule := p.e_main.Eval(p.decls)
	return FGGProgram{p.decls, e.(FGGExpr), p.printf}, rule
}

func (p FGGProgram) String() string {
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
	t_name Name
	Psi    BigPsi
	fDecls []FieldDecl
}

var _ TypeDecl = STypeLit{}

func (s STypeLit) GetName() Name              { return s.t_name }
func (s STypeLit) GetBigPsi() BigPsi          { return s.Psi }
func (s STypeLit) GetFieldDecls() []FieldDecl { return s.fDecls }

func (s STypeLit) Ok(ds []Decl) {
	s.Psi.Ok(ds, BigPsi{})
	seen := make(map[Name]FieldDecl)
	delta := s.Psi.ToDelta()
	for _, v := range s.fDecls {
		if _, ok := seen[v.field]; ok {
			panic("Duplicate field name: " + v.field + "\n\t" + s.String())
		}
		v.u.Ok(ds, delta)
	}
}

func (s STypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(string(s.t_name))
	b.WriteString(s.Psi.String())
	b.WriteString(" struct {")
	if len(s.fDecls) > 0 {
		b.WriteString(" ")
		writeFieldDecls(&b, s.fDecls)
		b.WriteString(" ")
	}
	b.WriteString("}")
	return b.String()
}

type FieldDecl struct {
	field Name
	u     Type // u=tau
}

var _ FGGNode = FieldDecl{}

func (fd FieldDecl) GetName() Name { return fd.field }
func (fd FieldDecl) GetType() Type { return fd.u }

func (fd FieldDecl) Subs(subs map[TParam]Type) FieldDecl {
	return FieldDecl{fd.field, fd.u.TSubs(subs)}
}

func (fd FieldDecl) String() string {
	return fd.field + " " + fd.u.String()
}

/* MethDecl, ParamDecl */

type MethDecl struct {
	x_recv   Name // CHECKME: better to be Variable?  (etc. for other such Names)
	t_recv   Name // N.B. t_S
	Psi_recv BigPsi
	// N.B. receiver elements "decomposed" because Psi (not TNamed, cf. fg.MDecl uses ParamDecl)
	name     Name // Refactor to embed Sig?
	Psi_meth BigPsi
	pDecls   []ParamDecl
	u_ret    Type // Return
	e_body   FGGExpr
}

var _ Decl = MethDecl{}

func (md MethDecl) GetRecvName() Name          { return md.x_recv }
func (md MethDecl) GetRecvTypeName() Name      { return md.t_recv }
func (md MethDecl) GetRecvPsi() BigPsi         { return md.Psi_recv }
func (md MethDecl) GetName() Name              { return md.name }
func (md MethDecl) GetMDeclPsi() BigPsi        { return md.Psi_meth } // MDecl in name to prevent false capture by TDecl interface
func (md MethDecl) GetParamDecls() []ParamDecl { return md.pDecls }
func (md MethDecl) GetReturn() Type            { return md.u_ret }
func (md MethDecl) GetBody() FGGExpr           { return md.e_body }

func (md MethDecl) Ok(ds []Decl) {
	if !isStructName(ds, md.t_recv) {
		panic("Receiver must be a struct type: not " + md.t_recv +
			"\n\t" + md.String())
	}
	md.Psi_recv.Ok(ds, BigPsi{}) // !!! premise ok missing
	md.Psi_meth.Ok(ds, md.Psi_recv)
	delta := md.Psi_recv.ToDelta()
	for _, v := range md.Psi_meth.tFormals {
		delta[v.name] = v.u_I
	}

	td := getTDecl(ds, md.t_recv)
	tfs_td := td.GetBigPsi().tFormals
	if len(tfs_td) != len(md.Psi_recv.tFormals) {
		panic("Receiver type parameter arity mismatch:\n\tmdecl=" + md.t_recv +
			md.Psi_recv.String() + ", tdecl=" + td.GetName() +
			"\n\t" + td.GetBigPsi().String())
	}
	for i := 0; i < len(tfs_td); i++ {
		subs_md := makeParamIndexSubs(md.Psi_recv)
		subs_td := makeParamIndexSubs(td.GetBigPsi())
		if !md.Psi_recv.tFormals[i].u_I.TSubs(subs_md). // Canonicalised
								Impls(ds, tfs_td[i].u_I.TSubs(subs_td)) {
			panic("Receiver parameter upperbound not a subtype of type decl upperbound:" +
				"\n\tmdecl=" + md.Psi_recv.tFormals[i].String() + ", tdecl=" +
				tfs_td[i].String())
		}
	}

	as := md.Psi_recv.Hat()                          // !!! submission version, x:t_S(a) => x:t_S(~a)
	gamma := Gamma{md.x_recv: TNamed{md.t_recv, as}} // CHECKME: can we give the bounds directly here instead of 'as'?
	seen := make(map[Name]Name)
	seen[md.x_recv] = md.x_recv
	for _, v := range md.pDecls {
		if _, ok := seen[v.name]; ok {
			panic("Duplicate receiver/param name: " + v.name + "\n\t" + md.String())
		}
		seen[v.name] = v.name
		v.u.Ok(ds, delta)
		gamma[v.name] = v.u
	}
	md.u_ret.Ok(ds, delta)
	allowStupid := false
	u := md.e_body.Typing(ds, delta, gamma, allowStupid)
	if !u.ImplsDelta(ds, delta, md.u_ret) {
		panic("Method body type must implement declared return type: found=" +
			u.String() + ", expected=" + md.u_ret.String() + "\n\t" + md.String())
	}
}

func (md MethDecl) ToSig() Sig {
	return Sig{md.name, md.Psi_meth, md.pDecls, md.u_ret}
}

func (md MethDecl) String() string {
	var b strings.Builder
	b.WriteString("func (")
	//b.WriteString(md.recv.String())
	b.WriteString(md.x_recv)
	b.WriteString(" ")
	b.WriteString(md.t_recv)
	b.WriteString(md.Psi_recv.String())
	b.WriteString(") ")
	b.WriteString(md.name)
	b.WriteString(md.Psi_meth.String())
	b.WriteString("(")
	writeParamDecls(&b, md.pDecls)
	b.WriteString(") ")
	b.WriteString(md.u_ret.String())
	b.WriteString(" { return ")
	b.WriteString(md.e_body.String())
	b.WriteString(" }")
	return b.String()
}

// Cf. FieldDecl
type ParamDecl struct {
	name Name // CHECKME: Variable?
	u    Type
}

var _ FGGNode = ParamDecl{}

func (pd ParamDecl) GetName() Name { return pd.name }
func (pd ParamDecl) GetType() Type { return pd.u }

func (pd ParamDecl) String() string {
	return pd.name + " " + pd.u.String()
}

/* ITypeLit, Sig */

type ITypeLit struct {
	t_I   Name
	Psi   BigPsi
	specs []Spec
}

var _ TypeDecl = ITypeLit{}

func (c ITypeLit) GetName() Name     { return c.t_I }
func (c ITypeLit) GetBigPsi() BigPsi { return c.Psi }
func (c ITypeLit) GetSpecs() []Spec  { return c.specs }

func (c ITypeLit) Ok(ds []Decl) {
	c.Psi.Ok(ds, BigPsi{})
	seen_g := make(map[Name]Sig)    // !!! unique(~S) more flexible
	seen_u := make(map[string]Type) // key is u.String()
	for _, v := range c.specs {
		switch s := v.(type) {
		case Sig:
			if _, ok := seen_g[s.meth]; ok {
				panic("Multiple sigs with name: " + s.meth + "\n\t" + c.String())
			}
			seen_g[s.meth] = s
			s.Ok(ds, c.Psi)
		case TNamed:
			k := s.String()
			if _, ok := seen_u[k]; ok {
				panic("Repeat embedding of type: " + k + "\n\t" + c.String())
			}
			seen_u[k] = s
			if !IsNamedIfaceType(ds, s) { // CHECKME: allow embed type param?
				panic("Embedded type must be a named interface, not: " + k + "\n\t" + c.String())
			}
			s.Ok(ds, c.Psi.ToDelta())
		default:
			panic("Unknown Spec kind: " + reflect.TypeOf(v).String() + "\n\t" +
				c.String())
		}
	}
}

func (c ITypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(string(c.t_I))
	b.WriteString(c.Psi.String())
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
	Psi    BigPsi // Add-meth-tparams
	pDecls []ParamDecl
	u_ret  Type
}

var _ Spec = Sig{}

func (g Sig) GetMethod() Name            { return g.meth }
func (g Sig) GetPsi() BigPsi             { return g.Psi }
func (g Sig) GetParamDecls() []ParamDecl { return g.pDecls }
func (g Sig) GetReturn() Type            { return g.u_ret }

func (g Sig) TSubs(subs map[TParam]Type) Sig {
	tfs := make([]TFormal, len(g.Psi.tFormals))
	for i := 0; i < len(g.Psi.tFormals); i++ {
		tf := g.Psi.tFormals[i]
		tfs[i] = TFormal{tf.name, tf.u_I.TSubs(subs)}
	}
	ps := make([]ParamDecl, len(g.pDecls))
	for i := 0; i < len(ps); i++ {
		pd := g.pDecls[i]
		ps[i] = ParamDecl{pd.name, pd.u.TSubs(subs)}
	}
	u := g.u_ret.TSubs(subs)
	return Sig{g.meth, BigPsi{tfs}, ps, u}
}

func (g Sig) Ok(ds []Decl, env BigPsi) {
	env.Ok(ds, BigPsi{})
	g.Psi.Ok(ds, env)
	delta := env.ToDelta()
	for _, v := range g.Psi.tFormals {
		delta[v.name] = v.u_I
	}
	seen := make(map[Name]ParamDecl)
	for _, v := range g.pDecls {
		if _, ok := seen[v.name]; ok {
			panic("Duplicate variable name " + v.name + ":\n\t" + g.String())
		}
		seen[v.name] = v
		v.u.Ok(ds, delta)
	}
	g.u_ret.Ok(ds, delta)
}

func (g Sig) GetSigs(_ []Decl) []Sig {
	return []Sig{g}
}

func (g Sig) String() string {
	var b strings.Builder
	b.WriteString(g.meth)
	b.WriteString(g.Psi.String())
	b.WriteString("(")
	writeParamDecls(&b, g.pDecls)
	b.WriteString(") ")
	b.WriteString(g.u_ret.String())
	return b.String()
}

/* Aux, helpers */

/*func BigPsiOk(ds []Decl, env BigPsi, Psi BigPsi) {
	Psi.Ok(ds)
	delta := Psi.ToDelta()
	for _, v := range Psi.tFormals {
		u_I, _ := v.u_I.(TNamed) // \tau_I, already checked by psi.Ok
		u_I.Ok(ds, delta)        // !!! Submission version T-Type, t_i => t_I
	}
}*/

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
