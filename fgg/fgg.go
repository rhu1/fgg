package fgg

import "fmt"
import "reflect"
import "strings"

import "github.com/rhu1/fgg/base"

var _ = fmt.Errorf
var _ = reflect.Append

/* Export */

func NewTName(t Name, us []Type) TName { return TName{t, us} }

/* Name, Type, Type param, Type name -- !!! submission version, "Type name" overloaded */

type Name = base.Name // TODO: tidy up refactoring, due to introducing base

type Type interface {
	TSubs(subs map[TParam]Type) Type // TODO: factor out Subs type? -- N.B. map is TEnv
	Impls(ds []Decl, delta TEnv, u Type) bool
	Ok(ds []Decl, delta TEnv)
	Equals(u Type) bool
	String() string
	ToGoString() string
}

type TParam Name

var _ Type = TParam("")

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
func (a TParam) Impls(ds []Decl, delta TEnv, u Type) bool {
	if a1, ok := u.(TParam); ok {
		return a == a1
	} else {
		return bounds(delta, a).Impls(ds, delta, u)
	}
}

func (a TParam) Ok(ds []Decl, delta TEnv) {
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

// TODO: rename TNamed
type TName struct {
	t  Name
	us []Type
}

var _ Type = TName{}
var _ Spec = TName{}

// TODO: refactor
func (u0 TName) GetName() Name {
	return u0.t
}

// TODO: refactor
func (u0 TName) GetTArgs() []Type {
	return u0.us
}

func (u0 TName) TSubs(subs map[TParam]Type) Type {
	us := make([]Type, len(u0.us))
	for i := 0; i < len(us); i++ {
		us[i] = u0.us[i].TSubs(subs)
	}
	return TName{u0.t, us}
}

// u0 <: 1
func (u0 TName) Impls(ds []Decl, delta TEnv, u Type) bool {
	if isStructTName(ds, u) {
		return isStructTName(ds, u0) && u0.Equals(u) // Asks equality of nested TParam
	}
	if _, ok := u.(TParam); ok { // e.g., fgg_test.go, Test014
		panic("Type name does not implement open type param: found=" + u0.String() + ", expected=" + u.String())
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

func (u0 TName) Ok(ds []Decl, delta TEnv) {
	td := getTDecl(ds, u0.t)
	psi := td.GetTFormals()
	if len(psi.tfs) != len(u0.us) {
		var b strings.Builder
		b.WriteString("Arity mismatch between type formals and actuals: formals=")
		b.WriteString(psi.String())
		b.WriteString(" actuals=")
		writeTypes(&b, u0.us)
		b.WriteString("\n\t")
		b.WriteString(u0.String())
		panic(b.String())
	}
	subs := make(map[TParam]Type)
	for i := 0; i < len(psi.tfs); i++ {
		subs[psi.tfs[i].a] = u0.us[i]
	}
	for i := 0; i < len(psi.tfs); i++ {
		actual := psi.tfs[i].a.TSubs(subs) // CHECKME: submission version T-Named, subs applied to Delta?
		formal := psi.tfs[i].u.TSubs(subs)
		if !actual.Impls(ds, delta, formal) { // tfs[i].u is a \tau_I, checked by TDecl.Ok
			panic("Type actual must implement type formal: actual=" + actual.String() +
				" formal=" + formal.String())
		}
	}
	for _, v := range u0.us {
		v.Ok(ds, delta)
	}
}

// \tau_I is a Spec, but not \tau_S -- this aspect is currently "dynamically typed"
func (u TName) GetSigs(ds []Decl) []Sig {
	if !isInterfaceTName(ds, u) { // isStructType would be more efficient
		panic("Cannot use non-interface type as a Spec: " + u.String() +
			" is a " + reflect.TypeOf(u).String())
	}
	td := getTDecl(ds, u.t).(ITypeLit)
	var res []Sig
	for _, s := range td.ss {
		res = append(res, s.GetSigs(ds)...)
	}
	return res
}

func (u0 TName) Equals(u Type) bool {
	if _, ok := u.(TName); !ok {
		return false
	}
	u1 := u.(TName)
	if u0.t != u1.t || len(u0.us) != len(u1.us) {
		return false
	}
	for i := 0; i < len(u0.us); i++ {
		if !u0.us[i].Equals(u1.us[i]) { // Asks equality of nested TParam
			return false
		}
	}
	return true
}

func (u TName) String() string {
	var b strings.Builder
	b.WriteString(string(u.t))
	b.WriteString("(")
	writeTypes(&b, u.us)
	b.WriteString(")")
	return b.String()
}

func (u TName) ToGoString() string {
	var b strings.Builder
	b.WriteString("main.")
	b.WriteString(string(u.t))
	b.WriteString("(")
	writeToGoTypes(&b, u.us)
	b.WriteString(")")
	return b.String()
}

/* Context, Type context */

//type Env map[Variable]Type  // CHECKME: refactor?
type Env map[Name]Type

type TEnv map[TParam]Type

func (delta TEnv) String() string {
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

// TODO: tidy up refactoring, due to introducing base
type FGGNode = base.AstNode
type Decl = base.Decl

type TDecl interface {
	Decl
	GetTFormals() TFormals // TODO: rename? potential clash with, e.g., MDecl, can cause "false" interface satisfaction
}

type Spec interface {
	FGGNode
	GetSigs(ds []Decl) []Sig
}

type Expr interface {
	base.Expr // Using the same name "Expr", maybe rename this type to FGGExpr
	Subs(subs map[Variable]Expr) Expr
	TSubs(subs map[TParam]Type) Expr
	// Like gamma, delta is effectively immutable
	Typing(ds []Decl, delta TEnv, gamma Env, allowStupid bool) Type
	Eval(ds []Decl) (Expr, string)
}

/* Helpers */

// Based on FG version -- but currently no FGG equiv of isInterfaceType
func isStructType(ds []Decl, t Name) bool {
	for _, v := range ds {
		d, ok := v.(STypeLit)
		if ok && d.t == t {
			return true
		}
	}
	return false
}

// Check if u is a \tau_S -- implicitly must be a TName
func isStructTName(ds []Decl, u Type) bool {
	if u1, ok := u.(TName); ok {
		for _, v := range ds {
			d, ok := v.(STypeLit)
			if ok && d.t == u1.t {
				return true
			}
		}
	}
	return false
}

// TODO FIXME -- temp for visibility
func IsStructTName1(ds []Decl, u Type) bool {
	return isStructTName(ds, u)
}

// Check if u is a \tau_I -- N.B. looks for a *TName*, i.e., not a TParam
func isInterfaceTName(ds []Decl, u Type) bool {
	if u1, ok := u.(TName); ok {
		for _, v := range ds {
			d, ok := v.(ITypeLit)
			if ok && d.t == u1.t {
				return true
			}
		}
	}
	return false
}

// TODO: refactor
func IsInterfaceTName1(ds []Decl, u Type) bool {
	return isInterfaceTName(ds, u)
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
