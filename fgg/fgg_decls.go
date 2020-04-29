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
}

var _ base.Program = FGGProgram{}
var _ FGGNode = FGGProgram{}

func (p FGGProgram) GetDecls() []Decl   { return p.decls } // Return a copy?
func (p FGGProgram) GetMain() base.Expr { return p.e_main }
func (p FGGProgram) IsPrintf() bool     { return p.printf } // HACK

func (p FGGProgram) Ok(allowStupid bool) base.Type {
	if !allowStupid { // Hack, to print only "top-level" programs (not during Eval) -- cf. verbose
		/*fmt.Println("[Warning] Type lit OK (\"T ok\") not fully implemented yet " +
		"(e.g., distinct type/field/method names, etc.)")*/
	}
	for _, v := range p.decls {
		switch d := v.(type) {
		case TDecl:
			d.Ok(p.decls)
		case MDecl:
			d.Ok(p.decls)
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

var _ TDecl = STypeLit{}

func (s STypeLit) GetName() Name              { return s.t_name }
func (s STypeLit) GetBigPsi() BigPsi          { return s.Psi }
func (s STypeLit) GetFieldDecls() []FieldDecl { return s.fDecls }

func (s STypeLit) Ok(ds []Decl) {
	TDeclOk(ds, s)
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

/* MDecl, ParamDecl */

type MDecl struct {
	x_recv  Name // CHECKME: better to be Variable?  (etc. for other such Names)
	t_recv  Name // N.B. t_S
	PsiRecv BigPsi
	// N.B. receiver elements "decomposed" because Psi (not TNamed, cf. fg.MDecl uses ParamDecl)
	name    Name // Refactor to embed Sig?
	PsiMeth BigPsi
	pDecls  []ParamDecl
	u_ret   Type // Return
	e_body  FGGExpr
}

var _ Decl = MDecl{}

func (md MDecl) GetRecvName() Name          { return md.x_recv }
func (md MDecl) GetRecvTypeName() Name      { return md.t_recv }
func (md MDecl) GetRecvPsi() BigPsi         { return md.PsiRecv }
func (md MDecl) GetName() Name              { return md.name }
func (md MDecl) GetMDeclPsi() BigPsi        { return md.PsiMeth } // MDecl in name to prevent false capture by TDecl interface
func (md MDecl) GetParamDecls() []ParamDecl { return md.pDecls }
func (md MDecl) GetReturn() Type            { return md.u_ret }
func (md MDecl) GetBody() FGGExpr           { return md.e_body }

func (md MDecl) Ok(ds []Decl) {
	if !isStructName(ds, md.t_recv) {
		panic("Receiver must be a struct type: not " + md.t_recv +
			"\n\t" + md.String())
	}
	md.PsiRecv.Ok(ds)
	md.PsiMeth.Ok(ds)

	td := getTDecl(ds, md.t_recv)
	tfs_td := td.GetBigPsi().tFormals
	if len(tfs_td) != len(md.PsiRecv.tFormals) {
		panic("Receiver parameter arity mismatch:\n\tmdecl=" + md.t_recv +
			md.PsiRecv.String() + ", tdecl=" + td.GetName() + td.GetBigPsi().String())
	}
	for i := 0; i < len(tfs_td); i++ {
		subs_md := makeParamIndexSubs(md.PsiRecv)
		subs_td := makeParamIndexSubs(td.GetBigPsi())
		if !md.PsiRecv.tFormals[i].u_I.TSubs(subs_md).
			Impls(ds, tfs_td[i].u_I.TSubs(subs_td)) {
			panic("Receiver parameter upperbound not a subtype of type decl upperbound:" +
				"\n\tmdecl=" + md.PsiRecv.tFormals[i].String() + ", tdecl=" +
				tfs_td[i].String())
		}
	}

	delta := md.PsiRecv.ToDelta()
	for _, v := range md.PsiRecv.tFormals {
		v.u_I.Ok(ds, delta)
	}

	delta1 := md.PsiMeth.ToDelta()
	for k, v := range delta {
		delta1[k] = v
	}
	for _, v := range md.PsiMeth.tFormals {
		v.u_I.Ok(ds, delta1)
	}

	as := make([]Type, len(md.PsiRecv.tFormals)) // !!! submission version, x:t_S(a) => x:t_S(~a)
	for i := 0; i < len(md.PsiRecv.tFormals); i++ {
		as[i] = md.PsiRecv.tFormals[i].name
	}
	gamma := Gamma{md.x_recv: TNamed{md.t_recv, as}} // CHECKME: can we give the bounds directly here instead of 'as'?
	for _, v := range md.pDecls {
		gamma[v.name] = v.u
	}
	allowStupid := false
	u := md.e_body.Typing(ds, delta1, gamma, allowStupid)
	if !u.ImplsDelta(ds, delta1, md.u_ret) {
		panic("Method body type must implement declared return type: found=" +
			u.String() + ", expected=" + md.u_ret.String() + "\n\t" + md.String())
	}
}

func (md MDecl) ToSig() Sig {
	return Sig{md.name, md.PsiMeth, md.pDecls, md.u_ret}
}

func (md MDecl) String() string {
	var b strings.Builder
	b.WriteString("func (")
	//b.WriteString(md.recv.String())
	b.WriteString(md.x_recv)
	b.WriteString(" ")
	b.WriteString(md.t_recv)
	b.WriteString(md.PsiRecv.String())
	b.WriteString(") ")
	b.WriteString(md.name)
	b.WriteString(md.PsiMeth.String())
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
	psi   BigPsi // TODO: rename Psi
	specs []Spec
}

var _ TDecl = ITypeLit{}

func (c ITypeLit) GetName() Name     { return c.t_I }
func (c ITypeLit) GetBigPsi() BigPsi { return c.psi }
func (c ITypeLit) GetSpecs() []Spec  { return c.specs }

func (c ITypeLit) Ok(ds []Decl) {
	TDeclOk(ds, c)
	for _, v := range c.specs {
		// TODO: check Sigs OK?  e.g., "type IA(type ) interface { m1(type )() Any };" while missing Any
		if g, ok := v.(Sig); ok {
			g.Ok(ds)
		}
	}
	// In general, also missing checks for, e.g., unique type/field/method names -- cf., TDeclOk
}

func (c ITypeLit) String() string {
	var b strings.Builder
	b.WriteString("type ")
	b.WriteString(string(c.t_I))
	b.WriteString(c.psi.String())
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
	psi    BigPsi // Add-meth-tparams  // TODO: rename Psi
	pDecls []ParamDecl
	u_ret  Type
}

var _ Spec = Sig{}

func (g Sig) GetMethod() Name            { return g.meth }
func (g Sig) GetPsi() BigPsi             { return g.psi }
func (g Sig) GetParamDecls() []ParamDecl { return g.pDecls }
func (g Sig) GetReturn() Type            { return g.u_ret }

func (g Sig) TSubs(subs map[TParam]Type) Sig {
	tfs := make([]TFormal, len(g.psi.tFormals))
	for i := 0; i < len(g.psi.tFormals); i++ {
		tf := g.psi.tFormals[i]
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

func (g Sig) Ok(ds []Decl) {
	g.psi.Ok(ds)
	// TODO: check distinct param names, etc. -- N.B. interface may not be *used* (so may not be checked else where)
}

func (g Sig) GetSigs(_ []Decl) []Sig {
	return []Sig{g}
}

func (g Sig) String() string {
	var b strings.Builder
	b.WriteString(g.meth)
	b.WriteString(g.psi.String())
	b.WriteString("(")
	writeParamDecls(&b, g.pDecls)
	b.WriteString(") ")
	b.WriteString(g.u_ret.String())
	return b.String()
}

/* Aux, helpers */

func TDeclOk(ds []Decl, td TDecl) {
	psi := td.GetBigPsi()
	psi.Ok(ds)
	delta := psi.ToDelta()
	for _, v := range psi.tFormals {
		u_I, _ := v.u_I.(TNamed) // \tau_I, already checked by psi.Ok
		u_I.Ok(ds, delta)        // !!! Submission version T-Type, t_i => t_I
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
