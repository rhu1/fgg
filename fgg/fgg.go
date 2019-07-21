package fgg

import "fmt"
import "reflect"
import "strings"

var _ = fmt.Errorf
var _ = reflect.Append

/* Name, Type, Type param, Type name -- !!! submission version, "Type name" overloaded */

type Name = string

type Type interface {
	Subs(subs map[TParam]Type) Type
	Impls(ds []Decl, delta TEnv, u Type) bool
	Ok(ds []Decl, delta TEnv)
	Equals(u Type) bool
	String() string
}

type TParam Name

var _ Type = TParam("")

func (a TParam) Subs(subs map[TParam]Type) Type {
	res, ok := subs[a]
	if !ok {
		panic("Unknown param: " + a.String())
	}
	return res
}

// u0 <: u
func (a TParam) Impls(ds []Decl, delta TEnv, u Type) bool {
	return a == u || a.Impls(ds, delta, bounds(delta, u))
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

type TName struct {
	t  Name
	us []Type
}

var _ Type = TName{}

func (u0 TName) Subs(subs map[TParam]Type) Type {
	us := make([]Type, len(u0.us))
	for i := 0; i < len(us); i++ {
		us[i] = u0.us[i].Subs(subs)
	}
	return TName{u0.t, us}
}

// u0 <: 1
func (u0 TName) Impls(ds []Decl, delta TEnv, u Type) bool {
	if isStructTName(ds, u) {
		return isStructTName(ds, u0) && u0.Equals(u) // Asks equality of nested TParam
	}

	gs := methods(ds, u)   // u is a t_I
	gs0 := methods(ds, u0) // t0 may be any
	for k, g := range gs {
		g0, ok := gs0[k]
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
		panic(b.String())
	}
	subs := make(map[TParam]Type)
	for i := 0; i < len(psi.tfs); i++ {
		subs[psi.tfs[i].a] = u0.us[i]
	}
	for i := 0; i < len(psi.tfs); i++ {
		actual := psi.tfs[i].a.Subs(subs) // CHECKME: submission version T-Named, subs applied to Delta?
		formal := psi.tfs[i].u.Subs(subs)
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

/* Context, Type context */

//type Env map[Variable]Type  // FIXME ?
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

type FGGNode interface {
	String() string
}

type Decl interface {
	FGGNode
	GetName() Name
	Ok(ds []Decl)
}

type TDecl interface {
	Decl
	GetTFormals() TFormals
}

type Spec interface {
	FGGNode
	GetSigs(ds []Decl) []Sig
}

type Expr interface {
	FGGNode
	Subs(subs map[Variable]Expr) Expr
	TSubs(subs map[TParam]Type) Expr
	Eval(ds []Decl) (Expr, string)
	// Like gamma, delta is effectively immutable
	Typing(ds []Decl, delta TEnv, gamma Env, allowStupid bool) Type
}

/* Helpers */

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

func writeTypes(b *strings.Builder, us []Type) {
	if len(us) > 0 {
		b.WriteString(us[0].String())
		for _, v := range us[1:] {
			b.WriteString(", " + v.String())
		}
	}
}
