package fgr

import (
	"fmt"
	"reflect"
	//"strings"

	"github.com/rhu1/fgg/fgg"
)

var _ = fmt.Errorf
var _ = reflect.Append

/*HERE
- getTypeReps
- add meth-param RepDecls
- FGR eval
*/

/* FGGProgram */

func Obliterate(p_fgg fgg.FGGProgram) FGRProgram { // CHECKME can also subsume existing FGG-FG trans?

	ds_fgg := p_fgg.GetDecls()

	e_fgg := p_fgg.GetExpr().(fgg.Expr)
	var delta fgg.TEnv
	var gamma fgg.Env
	e_fgr := oblitExpr(ds_fgg, delta, gamma, e_fgg)

	var ds_fgr []Decl

	// Translate Decls (and collect wrappers from MDecls)
	/*ds_fgg := p.GetDecls()
	for i := 0; i < len(ds_fgg); i++ {
		d := ds_fgg[i]
		switch d1 := d.(type) {
		case fgg.STypeLit:
			s := fgrTransSTypeLit(d1)

			// Add getValue/getTypeRep to all (existing) t_S -- every t_S must implement t_0 -- TODO: factor out with wrappers
			//e_getv := fg.NewSelect(fg.NewVariable("x"), "value") // CHECKME: but t_S doesn't have value field, wrapper does?
			e_getv := NewStructLit(Type("Dummy_0"), []Expr{})
			t := Type(d1.GetName())
			getv := NewMDecl(NewParamDecl("x", t), "getValue",
				[]RepDecl{}, // FIXME
				[]ParamDecl{}, Type("Any_0"), e_getv)
			// gettr := ...TODO FIXME...

			ds_fgr = append(ds_fgr, s, getv)
		case fgg.ITypeLit:
			ds_fgr = append(ds_fgr, fgrTransITypeLit(d1))
		case fgg.MDecl:
			ds_fgr = append(ds_fgr, fgrTransMDecl(ds_fgg, d1, wrappers))
		default:
			panic("Unexpected Decl type " + reflect.TypeOf(d).String() +
				": " + d.String())
		}
	}*/

	return NewFGRProgram(ds_fgr, e_fgr)
}

func oblitExpr(ds_fgg []Decl, delta fgg.TEnv, gamma fgg.Env,
	e_fgg fgg.Expr) Expr {
	switch e := e_fgg.(type) {
	case fgg.Variable:
		return NewVariable(e.GetName())
	case fgg.StructLit:
		u := e.GetTName()
		t := Type(u.GetName())
		us := u.GetTArgs()
		es_fgg := e.GetArgs()
		es_fgr := make([]Expr, len(us)+len(es_fgg))
		for i := 0; i < len(us); i++ {
			es_fgr[i] = oblitMkRep(us[i])
		}
		for i := 0; i < len(es_fgg); i++ {
			es_fgr[len(us)+i] = oblitExpr(ds_fgg, delta, gamma, es_fgg[i])
		}
		return NewStructLit(t, es_fgr)
	case fgg.Select:
		e_fgr := oblitExpr(ds_fgg, delta, gamma, e.GetExpr())
		u := e_fgg.Typing(ds_fgg, delta, gamma, true).(fgg.TName)
		fds_fgg := fgg.Fields1(ds_fgg, u)
		f := e.GetName()
		var u_f fgg.Type = nil
		for _, fd_fgg := range fds_fgg {
			if fd_fgg.GetName() == f {
				u_f = fd_fgg.GetType()
			}
		}
		if u_f == nil {
			panic("Field not found in " + u.String() + ": " + f)
		}
		var res Expr
		res = NewSelect(e_fgr, f)
		res = NewAssert(res, toFgrTypeFromBounds(delta, u_f))
		return res
	case fgg.Call:
		panic("[TODO]: " + e.String())
	case fgg.Assert:
		x := oblitExpr(ds_fgg, delta, gamma, e.GetExpr())
		e1 := NewCall(x, "getTypeRep", []Expr{})
		u := e.GetType()
		e3 := NewAssert(x, toFgrTypeFromBounds(delta, u))
		return IfThenElse{e1, mkRep(u), e3} // TODO: New constructor
	default:
		panic("Unknown FGG Expr type: " + e_fgg.String())
	}
}

/* Helper */

func toFgrTypeFromBounds(delta fgg.TEnv, u fgg.Type) Type {
	return Type(fgg.Bounds1(delta, u).(fgg.TName).GetName())
}

// Post: TypeTree or TmpTParam
func oblitMkRep(u fgg.Type) Expr {
	switch u1 := u.(type) {
	case fgg.TParam:
		return TmpTParam{u1.String()}
	case fgg.TName:
		us := u1.GetTArgs()
		es := make([]Expr, len(us))
		for i := 0; i < len(us); i++ {
			es[i] = mkRep(us[i])
		}
		return TypeTree{Type(u1.GetName()), es}
	default:
		panic("Unknown fgg.Type kind " + reflect.TypeOf(u).String() + ": " +
			u.String())
	}
}
