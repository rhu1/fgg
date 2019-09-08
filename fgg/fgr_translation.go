package fgg

import (
	//"fmt"
	"reflect"
	//"strings"

	"github.com/rhu1/fgg/fg"
	//"github.com/rhu1/fgg/fgg"
)

// |\tau|_\Delta = t
func erase(delta TEnv, u Type) Name { //fg.Type {
	return bounds(delta, u).(TName).t
}

// |e_FGG|_(\Delta; \Gamma) = e_FGR
func translate(ds []Decl, delta TEnv, gamma Env, e Expr) fg.Expr {
	switch e1 := e.(type) {
	case Variable:
		u := e1.Typing(ds, delta, gamma, false)
		if isStructTName(ds, u) {
			return fg.NewVariable(e1.id)
		} else { // "interface" case
			// x.getValue().((mkRep u))
			getVal := fg.NewCall(fg.NewVariable(e1.id), Name("getValue"), []fg.Expr{})
			return fg.NewAssert(getVal, fg.Type(erase(delta, u))) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
	case StructLit:
		t := e1.u.t
		es := make([]fg.Expr, len(e1.es))
		fds := fields(ds, e1.u)
		subs := make(map[TParam]Type)
		psi := getTDecl(ds, t).GetTFormals()
		for i := 0; i < len(psi.tfs); i++ {
			subs[psi.tfs[i].a] = e1.u.us[i]
		}
		for i := 0; i < len(e1.es); i++ {
			u_i := fds[i].u
			es[i] = wrap(ds, delta, gamma, e1.es[i], u_i.TSubs(subs))
		}
		return fg.NewStructLit(fg.Type(t), es)
	default:
		panic("TODO " + reflect.TypeOf(e).String() + ": " + e.String())
	}
}

// Pre: type of e <: u
// `u` is "target type"
func wrap(ds []fg.Decl, delta TEnv, gamma Env, e Expr, u Type) fg.Expr {
	/*t := erase(u, delta)
	if _, ok := fg.isStructType(t)*/
	if isStructTName(ds, u) { // N.B. differs slightly from def -- because there is no FG t_S decl (yet)?
		return translate(ds, delta, gamma, e)
	} else if isInterfaceTName(ds, u) {
		targ := erase(delta, u)
		u1 := e.Typing(ds, delta, gamma, false)
		subj := erase(delta, u1)
		e1 := translate(ds, delta, gamma, e)
		return wrapper(targ, subj, e1)
	} else {
		panic("Invalid wrap case: e=" + e.String() + ", u=" + u.String())
	}
}

func wrapper(targ Name, subj Name, e fg.Expr) fg.StructLit {
	return fg.NewStructLit(fg.Type("Adptr_"+targ+"_"+subj), []fg.Expr{e}) // TODO: factor out naming
}
