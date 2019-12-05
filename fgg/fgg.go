package fgg

import "fmt"
import "reflect"
import "strings"

import "github.com/rhu1/fgg/base"

var _ = fmt.Errorf
var _ = reflect.Append

/* Export */

func NewTName(t Name, us []Type) TNamed       { return TNamed{t, us} }
func IsStructType(ds []Decl, u Type) bool     { return isStructType(ds, u) }
func IsNamedIfaceType(ds []Decl, u Type) bool { return isNamedIfaceType(ds, u) }

/* Aliases from base */
// TODO: refactor?

type Name = base.Name
type FGGNode = base.AstNode
type Decl = base.Decl

/* Name, Type, Type param, Type name -- !!! submission version, "Type name" overloaded */

// Name: see Aliases (at top)

type Type interface {
	TSubs(subs map[TParam]Type) Type // N.B. map is Delta -- TODO: factor out Subs type?
	Impls(ds []Decl, delta Delta, u Type) bool
	Ok(ds []Decl, delta Delta)
	Equals(u Type) bool
	String() string
	ToGoString() string
}

var _ Type = TParam("")

type TParam Name

func (a TParam) TSubs(subs map[TParam]Type) Type {
	res, ok := subs[a]
	if !ok {
		//panic("Unknown param: " + a.String())
		return a // CHECKME: ok? -- see TSubs in methods aux, w.r.t. meth-tparams that aren't in the subs map
		// Cf. Variable.Subs?
	}
	return res
}

// u0 <: u
func (a TParam) Impls(ds []Decl, delta Delta, u Type) bool {
	if a1, ok := u.(TParam); ok {
		return a == a1
	} else {
		return bounds(delta, a).Impls(ds, delta, u)
	}
}

func (a TParam) Ok(ds []Decl, delta Delta) {
	if _, ok := delta[a]; !ok {
		panic("Unknown type param: " + a.String())
	}
}

func (a TParam) Equals(u Type) bool {
	if b, ok := u.(TParam); ok {
		return a == b
	}
	return false
}

func (a TParam) String() string {
	return string(a)
}

func (a TParam) ToGoString() string {
	return string(a)
}

var _ Type = TNamed{}
var _ Spec = TNamed{}

// Convention: t=type name (t), u=FGG type (tau)
type TNamed struct {
	t_name Name
	u_args []Type
}

func (u0 TNamed) GetName() Name    { return u0.t_name }
func (u0 TNamed) GetTArgs() []Type { return u0.u_args }

func (u0 TNamed) TSubs(subs map[TParam]Type) Type {
	us := make([]Type, len(u0.u_args))
	for i := 0; i < len(us); i++ {
		us[i] = u0.u_args[i].TSubs(subs)
	}
	return TNamed{u0.t_name, us}
}

// u0 <: u
func (u0 TNamed) Impls(ds []Decl, delta Delta, u Type) bool {
	if isStructType(ds, u) {
		return isStructType(ds, u0) && u0.Equals(u) // Asks equality of nested TParam
	}
	if _, ok := u.(TParam); ok { // e.g., fgg_test.go, Test014
		panic("Type name does not implement open type param: found=" +
			u0.String() + ", expected=" + u.String())
	}

	gs := methods(ds, u)   // u is a t_I
	gs0 := methods(ds, u0) // t0 may be any
	for k, g := range gs {
		g0, ok := gs0[k]
		if ok {
		}
		if !ok || !g.EqExceptTParamsAndVars(g0) {
			return false
		}
	}
	return true
}

func (u0 TNamed) Ok(ds []Decl, delta Delta) {
	td := GetTDecl(ds, u0.t_name)
	psi := td.GetPsi()
	if len(psi.tFormals) != len(u0.u_args) {
		var b strings.Builder
		b.WriteString("Arity mismatch between type formals and actuals: formals=")
		b.WriteString(psi.String())
		b.WriteString(" actuals=")
		writeTypes(&b, u0.u_args)
		b.WriteString("\n\t")
		b.WriteString(u0.String())
		panic(b.String())
	}
	subs := make(map[TParam]Type)
	for i := 0; i < len(psi.tFormals); i++ {
		subs[psi.tFormals[i].name] = u0.u_args[i]
	}
	for i := 0; i < len(psi.tFormals); i++ {
		actual := psi.tFormals[i].name.TSubs(subs) // CHECKME: submission version T-Named, subs applied to Delta?
		formal := psi.tFormals[i].u_I.TSubs(subs)
		if !actual.Impls(ds, delta, formal) { // tfs[i].u is a \tau_I, checked by TDecl.Ok
			panic("Type actual must implement type formal: actual=" +
				actual.String() + " formal=" + formal.String())
		}
	}
	for _, v := range u0.u_args {
		v.Ok(ds, delta)
	}
}

// \tau_I is a Spec, but not \tau_S -- this aspect is currently "dynamically typed"
// From Spec
func (u TNamed) GetSigs(ds []Decl) []Sig {
	if !isNamedIfaceType(ds, u) { // isStructType would be more efficient
		panic("Cannot use non-interface type as a Spec: " + u.String() +
			" is a " + reflect.TypeOf(u).String())
	}
	td := GetTDecl(ds, u.t_name).(ITypeLit)
	var res []Sig
	for _, s := range td.specs {
		res = append(res, s.GetSigs(ds)...)
	}
	return res
}

func (u0 TNamed) Equals(u Type) bool {
	if _, ok := u.(TNamed); !ok {
		return false
	}
	u1 := u.(TNamed)
	if u0.t_name != u1.t_name || len(u0.u_args) != len(u1.u_args) {
		return false
	}
	for i := 0; i < len(u0.u_args); i++ {
		if !u0.u_args[i].Equals(u1.u_args[i]) { // Asks equality of nested TParam
			return false
		}
	}
	return true
}

func (u TNamed) String() string {
	var b strings.Builder
	b.WriteString(string(u.t_name))
	b.WriteString("(")
	writeTypes(&b, u.u_args)
	b.WriteString(")")
	return b.String()
}

func (u TNamed) ToGoString() string {
	var b strings.Builder
	b.WriteString("main.")
	b.WriteString(string(u.t_name))
	b.WriteString("(")
	writeToGoTypes(&b, u.u_args)
	b.WriteString(")")
	return b.String()
}

/* Context, Type context */

//type Gamma map[Variable]Type  // CHECKME: refactor?
type Gamma map[Name]Type
type Delta map[TParam]Type

func (delta Delta) String() string {
	res := "["
	first := true
	for k, v := range delta {
		if first {
			first = false
		} else {
			res = res + ", "
		}
		res = k.String() + ":" + v.String()
	}
	return res + "]"
}

/* AST base intefaces: FGGNode, Decl, TDecl, Spec, Expr */

// FGGNode, Decl: see Aliases (at top)

type TDecl interface {
	Decl
	GetPsi() pDecls // TODO: rename? potential clash with, e.g., MDecl, can cause "false" interface satisfaction
}

type Spec interface {
	FGGNode
	GetSigs(ds []Decl) []Sig
}

type FGGExpr interface {
	base.Expr
	Subs(subs map[Variable]FGGExpr) FGGExpr
	TSubs(subs map[TParam]Type) FGGExpr
	// gamma and delta should be treated immutably
	Typing(ds []Decl, delta Delta, gamma Gamma, allowStupid bool) Type
	Eval(ds []Decl) (FGGExpr, string)
}

/* Helpers */

// Based on FG version -- but currently no FGG equiv of isInterfaceType
// Helpful for MDecl.t_recv
func isStructName(ds []Decl, t Name) bool {
	for _, v := range ds {
		d, ok := v.(STypeLit)
		if ok && d.t_name == t {
			return true
		}
	}
	return false
}

// Check if u is a \tau_S -- implicitly must be a TNamed
func isStructType(ds []Decl, u Type) bool {
	if u1, ok := u.(TNamed); ok {
		for _, v := range ds {
			d, ok := v.(STypeLit)
			if ok && d.t_name == u1.t_name {
				return true
			}
		}
	}
	return false
}

// Check if u is a \tau_I -- N.B. looks for a *TNamed*, i.e., not a TParam
func isNamedIfaceType(ds []Decl, u Type) bool {
	if u1, ok := u.(TNamed); ok {
		for _, v := range ds {
			d, ok := v.(ITypeLit)
			if ok && d.t_I == u1.t_name {
				return true
			}
		}
	}
	return false
}

func writeTypes(b *strings.Builder, us []Type) {
	if len(us) > 0 {
		b.WriteString(us[0].String())
		for _, v := range us[1:] {
			b.WriteString(", " + v.String())
		}
	}
}

func writeToGoTypes(b *strings.Builder, us []Type) {
	if len(us) > 0 {
		b.WriteString(us[0].ToGoString())
		for _, v := range us[1:] {
			b.WriteString(", " + v.ToGoString())
		}
	}
}
