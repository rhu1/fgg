package fgg

import "fmt"
import "reflect"
import "strconv"
import "strings"

import "github.com/rhu1/fgg/base"

var _ = fmt.Errorf
var _ = reflect.Append
var _ = strconv.AppendBool

/* Export */

func NewTName(t Name, us []Type) TNamed       { return TNamed{t, us} }
func IsStructType(ds []Decl, u Type) bool     { return isStructType(ds, u) }
func IsNamedIfaceType(ds []Decl, u Type) bool { return isNamedIfaceType(ds, u) }

/* Constants */

var STRING_TYPE TNamed = TNamed{"string", []Type{}}

/* Aliases from base */
// TODO: refactor?

type Name = base.Name
type FGGNode = base.AstNode
type Decl = base.Decl

/* Name, Type, Type param, Type name -- !!! submission version, "Type name" overloaded */

// Name: see Aliases (at top)

type Type interface {
	base.Type
	ImplsDelta(ds []Decl, delta Delta, u Type) bool
	TSubs(subs map[TParam]Type) Type // N.B. map is Delta -- factor out a Subs type?
	SubsEta(eta Eta) TNamed
	SubsEta2(eta Eta2) Type
	Ok(ds []Decl, delta Delta)
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

func (a TParam) SubsEta(eta Eta) TNamed {
	res, ok := eta[a]
	if !ok {
		panic("Shouldn't get here: " + a)
	}
	return res
}

func (a TParam) SubsEta2(eta Eta2) Type {
	res, ok := eta[a]
	if !ok {
		return a
	}
	return res
}

// u0 <: u
func (a TParam) ImplsDelta(ds []Decl, delta Delta, u Type) bool {
	if a1, ok := u.(TParam); ok {
		return a == a1
	} else {
		return bounds(delta, a).ImplsDelta(ds, delta, u)
	}
}

// Cf. base.Type
func (a TParam) Impls(ds []Decl, u base.Type) bool {
	if _, ok := u.(Type); !ok {
		panic("Expected FGG type, not " + reflect.TypeOf(u).String() +
			":\n\t" + u.String())
	}
	return a.ImplsDelta(ds, make(Delta), u.(Type))
}

func (a TParam) Ok(ds []Decl, delta Delta) {
	if _, ok := delta[a]; !ok {
		panic("Unknown type param: " + a.String())
	}
}

func (a TParam) Equals(u base.Type) bool {
	if _, ok := u.(Type); !ok {
		panic("Expected FGG type, not " + reflect.TypeOf(u).String() +
			":\n\t" + u.String())
	}
	if b, ok := u.(TParam); ok {
		return a == b // TODO FIXME: need alpha -- pre- index-ify TNamed (and Impls)?  // TODO: add tests
	}
	return false
}

func (a TParam) String() string {
	return string(a)
}

func (a TParam) ToGoString() string {
	return string(a)
}

// Convention: t=type name (t), u=FGG type (tau)
type TNamed struct {
	t_name Name
	u_args []Type // SmallPsi
}

var _ Type = TNamed{}
var _ Spec = TNamed{}

func (u0 TNamed) GetName() Name    { return u0.t_name }
func (u0 TNamed) GetTArgs() []Type { return u0.u_args } // SmallPsi

func (u0 TNamed) TSubs(subs map[TParam]Type) Type {
	us := make([]Type, len(u0.u_args))
	for i := 0; i < len(us); i++ {
		us[i] = u0.u_args[i].TSubs(subs)
	}
	return TNamed{u0.t_name, us}
}

func (u0 TNamed) SubsEta(eta Eta) TNamed {
	//fmt.Println("555:", u0, eta)
	us := make([]Type, len(u0.u_args))
	for i := 0; i < len(us); i++ {
		us[i] = u0.u_args[i].SubsEta(eta)
	}
	return TNamed{u0.t_name, us}
}

func (u0 TNamed) SubsEta2(eta Eta2) Type {
	//fmt.Println("555:", u0, eta)
	us := make([]Type, len(u0.u_args))
	for i := 0; i < len(us); i++ {
		us[i] = u0.u_args[i].SubsEta2(eta)
	}
	return TNamed{u0.t_name, us}
}

// u0 <: u
// delta unused here (cf. TParam.ImplsDelta)
func (u0 TNamed) ImplsDelta(ds []Decl, delta Delta, u Type) bool {
	u_fgg := u.(Type)
	if isStructType(ds, u_fgg) {
		return isStructType(ds, u0) && u0.Equals(u_fgg) // Asks equality of nested TParam
	}
	if _, ok := u.(TParam); ok { // e.g., fgg_test.go, Test014
		panic("Type name does not implement open type param: found=" +
			u0.String() + ", expected=" + u_fgg.String())
	}

	gs := methods(ds, u_fgg) // u is a t_I
	gs0 := methods(ds, u0)   // t0 may be any
	for k, g := range gs {
		g0, ok := gs0[k]
		if !ok || !sigAlphaEquals(g0, g) {
			return false
		}
	}
	return true
}

// !!! Sig in FGG includes ~a and ~x, which naively breaks "impls"
func sigAlphaEquals(g0 Sig, g Sig) bool {
	if len(g0.Psi.tFormals) != len(g.Psi.tFormals) || len(g0.pDecls) != len(g.pDecls) {
		return false
	}
	subs0 := makeParamIndexSubs(g0.Psi)
	subs := makeParamIndexSubs(g.Psi)
	for i := 0; i < len(g0.Psi.tFormals); i++ {
		if !g0.Psi.tFormals[i].u_I.TSubs(subs0).
			Equals(g.Psi.tFormals[i].u_I.TSubs(subs)) {
			return false
		}
	}
	for i := 0; i < len(g0.pDecls); i++ {
		if !g0.pDecls[i].u.TSubs(subs0).Equals(g.pDecls[i].u.TSubs(subs)) {
			return false
		}
	}
	return g0.meth == g.meth && g0.u_ret.Equals(g.u_ret)
}

// CHECKME: Used by sigAlphaEquals, and MDecl.OK (for covariant receiver bounds)
func makeParamIndexSubs(Psi BigPsi) Delta {
	subs := make(Delta)
	for j := 0; j < len(Psi.tFormals); j++ {
		//subs[Psi.tFormals[j].name] = Psi.tFormals[j].name
		subs[Psi.tFormals[j].name] = TParam("α" + strconv.Itoa(j+1))
	}
	return subs
}

// Cf. base.Type
func (u0 TNamed) Impls(ds []Decl, u base.Type) bool {
	if _, ok := u.(Type); !ok {
		panic("Expected FGG type, not " + reflect.TypeOf(u).String() +
			":\n\t" + u.String())
	}
	return u0.ImplsDelta(ds, make(Delta), u.(Type))
}

func (u0 TNamed) Ok(ds []Decl, delta Delta) {
	td := GetTDecl(ds, u0.t_name)
	psi := td.GetBigPsi()
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
		actual := psi.tFormals[i].name.TSubs(subs)
		// CHECKME: submission T-Named, subs applied to Delta? -- already applied, Delta is coming from the subs context
		formal := psi.tFormals[i].u_I.TSubs(subs)
		if !actual.ImplsDelta(ds, delta, formal) { // tfs[i].u is a \tau_I, checked by TDecl.Ok
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

func (u0 TNamed) Equals(u base.Type) bool {
	if _, ok := u.(Type); !ok {
		panic("Expected FGG type, not " + reflect.TypeOf(u).String() +
			":\n\t" + u.String())
	}
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

/* Type formals and actuals */

// Pre: len(as) == len(us)
// Wrapper for []TFormal (cf. e.g., FieldDecl), only because of "(type ...)" syntax
// Also ranged over by big phi
type BigPsi struct {
	tFormals []TFormal
}

func (Psi BigPsi) GetTFormals() []TFormal { return Psi.tFormals }

func (Psi BigPsi) Ok(ds []Decl) {
	for _, v := range Psi.tFormals {
		u, ok := v.u_I.(TNamed)
		if !ok {
			panic("Upper bound must be of the form \"t_I(type ...)\": not " +
				v.u_I.String())
		}
		if !IsNamedIfaceType(ds, u) { // CHECKME: subsumes above TName check (looks for \tau_S)
			panic("Upper bound must be an interface type: not " + u.String())
		}
	}
}

func (Psi BigPsi) ToDelta() Delta {
	delta := make(map[TParam]Type)
	for _, v := range Psi.tFormals {
		delta[v.name] = v.u_I
	}
	return delta
}

// The ordered value set of ToDelta
func (Psi BigPsi) Hat() SmallPsi {
	res := make(SmallPsi, len(Psi.tFormals))
	for i, v := range Psi.tFormals {
		res[i] = v.u_I
	}
	return res
}

func (Psi BigPsi) String() string {
	var b strings.Builder
	b.WriteString("(type ") // Includes "(...)" -- cf. e.g., writeFieldDecls
	if len(Psi.tFormals) > 0 {
		b.WriteString(Psi.tFormals[0].String())
		for _, v := range Psi.tFormals[1:] {
			b.WriteString(", ")
			b.WriteString(v.String())
		}
	}
	b.WriteString(")")
	return b.String()
}

type TFormal struct {
	name TParam
	u_I  Type
	// CHECKME: submission version, upper bound \tau_I is only "of the form t_I(~\tau)"? -- i.e., not \alpha?
	// ^If so, then can refine to TNamed
}

func (tf TFormal) GetTParam() TParam   { return tf.name }
func (tf TFormal) GetUpperBound() Type { return tf.u_I }

func (tf TFormal) String() string {
	return string(tf.name) + " " + tf.u_I.String()
}

// Type actuals
// Also ranged over by small phi
type SmallPsi []Type // CHECKME: Currently only used in omega/monom, maybe deprecate?

func (x0 SmallPsi) TSubs(subs map[TParam]Type) SmallPsi {
	res := make(SmallPsi, len(x0))
	for i, v := range x0 {
		res[i] = v.TSubs(subs)
	}
	return res
}

func (x0 SmallPsi) String() string {
	var b strings.Builder
	for _, v := range x0 {
		b.WriteString(v.String())
	}
	return b.String()
}

func (x0 SmallPsi) Equals(x SmallPsi) bool {
	if len(x0) != len(x) {
		return false
	}
	for i := 0; i < len(x0); i++ {
		if !x0[i].Equals(x[i]) {
			return false
		}
	}
	return true
}

/* Context, Type context, Substitutions */

//type Gamma map[Variable]Type  // CHECKME: refactor?
type Gamma map[Name]Type
type Delta map[TParam]Type // Type intended to be an upper bound
type Eta map[TParam]TNamed // TNamed intended to be a ground

type Eta2 map[TParam]Type

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

// Pre: len(psi) == len(Psi.GetTFormals()); psi all ground
func MakeEta(Psi BigPsi, psi SmallPsi) Eta {
	eta := make(Eta)
	tfs := Psi.tFormals
	for i := 0; i < len(tfs); i++ {
		eta[tfs[i].name] = psi[i].(TNamed)
	}
	return eta
}

func MakeEta2(Psi BigPsi, psi SmallPsi) Eta2 {
	eta2 := make(Eta2)
	tfs := Psi.tFormals
	for i := 0; i < len(tfs); i++ {
		eta2[tfs[i].name] = psi[i]
	}
	return eta2
}

/* AST base intefaces: FGGNode, Decl, TypeDecl, Spec, Expr */

// FGGNode, Decl: see Aliases (at top)

type TypeDecl interface {
	Decl
	GetBigPsi() BigPsi // TODO: rename? potential clash with, e.g., MDecl, can cause "false" interface satisfaction
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
		if u1.Equals(STRING_TYPE) { // TODO CHECKME
			return true
		}
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
