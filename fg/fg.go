package fg

import "reflect"

import "github.com/rhu1/fgg/base"

/* Aliases from base */
// TODO: refactor?

type Name = base.Name
type FGNode = base.AstNode
type Decl = base.Decl

/* Constants */

var STRING_TYPE Type = Type("string")

/* Name, Context, Type */

// Name: see Aliases (at top)

type Gamma map[Name]Type // TODO: should be Variable rather than Name? -- though Variable is an Expr

type Type Name // Type definition (cf. alias)

var _ base.Type = Type("")
var _ Spec = Type("")

// Pre: t0, t are known types
// t0 <: t
func (t0 Type) Impls(ds []Decl, t base.Type) bool {
	if _, ok := t.(Type); !ok {
		panic("Expected FGR type, not " + reflect.TypeOf(t).String() +
			":\n\t" + t.String())
	}
	t_fg := t.(Type)
	if isStructType(ds, t_fg) {
		return isStructType(ds, t0) && t0 == t_fg
	}

	gs := methods(ds, t_fg) // t is a t_I
	gs0 := methods(ds, t0)  // t0 may be any
	for k, g := range gs {
		g0, ok := gs0[k]
		if !ok || !g.EqExceptVars(g0) {
			return false
		}
	}
	return true
}

// t_I is a Spec, but not t_S -- this aspect is currently "dynamically typed"
// From Spec
func (t Type) GetSigs(ds []Decl) []Sig {
	if !isInterfaceType(ds, t) { // isStructType would be more efficient
		panic("Cannot use non-interface type as a Spec: " + t.String() +
			" is a " + reflect.TypeOf(t).String())
	}
	td := getTDecl(ds, t).(ITypeLit)
	var res []Sig
	for _, s := range td.specs {
		res = append(res, s.GetSigs(ds)...)
	}
	return res
}

func (t0 Type) Equals(t base.Type) bool {
	if _, ok := t.(Type); !ok {
		panic("Expected FGR type, not " + reflect.TypeOf(t).String() +
			":\n\t" + t.String())
	}
	return t0 == t.(Type)
}

func (t Type) String() string {
	return string(t)
}

/* AST base intefaces: FGNode, Decl, TDecl, Spec, Expr */

// FGNode, Decl: see Aliases (at top)

type TDecl interface {
	Decl
	GetType() Type // In FG, GetType() == Type(GetName())
}

// A Sig or a Type (specifically a t_I -- bad t_S usage raises a run-time error, cf. Type.GetSigs)
type Spec interface {
	FGNode
	GetSigs(ds []Decl) []Sig
}

type FGExpr interface {
	base.Expr
	Subs(subs map[Variable]FGExpr) FGExpr

	// N.B. gamma should be treated immutably (and ds, of course)
	// (No typing rule modifies gamma, except the T-Func bootstrap)
	Typing(ds []Decl, gamma Gamma, allowStupid bool) Type

	// string is the type name of the "actually evaluated" expr (within the eval context)
	// CHECKME: resulting Exprs are not "parsed" from source, OK?
	Eval(ds []Decl) (FGExpr, string)

	//IsPanic() bool  // TODO "explicit" FG panic -- cf. underlying runtime panic
}

/* Helpers */

func isStructType(ds []Decl, t Type) bool {
	if t == STRING_TYPE { // TODO CHECKME
		return true
	}
	for _, v := range ds {
		d, ok := v.(STypeLit)
		if ok && d.t_S == t {
			return true
		}
	}
	return false
}

func isInterfaceType(ds []Decl, t Type) bool {
	for _, v := range ds {
		d, ok := v.(ITypeLit)
		if ok && d.t_I == t {
			return true
		}
	}
	return false
}
