package fgg

import (
	"fmt"
	"reflect"
	//"strings"

	"github.com/rhu1/fgg/fg"
	//"github.com/rhu1/fgg/fgg"
)

var _ = fmt.Errorf

/* FGGProgram */

type adaptPair struct {
	sub   fg.Type
	super fg.Type // The "target" type, a t_I
}

func Translate(p FGGProgram) fg.FGProgram { // TODO FIXME: FGR -- TODO also can subsume existing FGG-FG trans?
	var ds []Decl

	// Add t_0 (etc.) to ds
	// TODO: factor out constants
	Any_0 := fg.NewITypeLit(fg.Type("Any_0"), []fg.Spec{})
	Dummy_0 := fg.NewSTypeLit(fg.Type("Dummy_0"), []fg.FieldDecl{})
	ToAny_0 := fg.NewSTypeLit(fg.Type("ToAny_0"), []fg.FieldDecl{fg.NewFieldDecl("any", fg.Type("Any_0"))})
	ds = append(ds, Any_0, Dummy_0, ToAny_0)

	getValue := fg.NewSig("getValue", []fg.ParamDecl{}, fg.Type("Any_0")) // TODO: rename "unwrap"?
	//getTypeRep := fg.NewSig("getTypeRep", []fg.ParamDecl{}, fg.Type("...TODO..."))
	ss_0 := []fg.Spec{getValue}
	t_0 := fg.NewITypeLit(fg.Type("t_0"), ss_0) // TODO FIXME? Go doesn't allow "overlapping" interfaces
	ds = append(ds, t_0)

	wrappers := make(map[fg.Type]adaptPair) // Populated by visiting MDecl and main Expr

	for i := 0; i < len(p.ds); i++ {
		d := p.ds[i]
		switch d1 := d.(type) {
		case STypeLit:
			s := translateSTypeLit(d1)

			// Add getValue/getTypeRep to all (existing) t_S -- every t_S must implement t_0 -- TODO: factor out with wrappers
			//e_getv := fg.NewSelect(fg.NewVariable("x"), "value") // CHECKME: but t_S doesn't have value field, wrapper does?
			e_getv := fg.NewStructLit(fg.Type("Dummy_0"), []fg.Expr{})
			getv := fg.NewMDecl(fg.NewParamDecl("x", fg.Type(d1.t)), "getValue",
				[]fg.ParamDecl{}, fg.Type("Any_0"), e_getv)
			// gettr := ...TODO...

			ds = append(ds, s, getv)
		case ITypeLit:
			ds = append(ds, translateITypeLit(d1))
		case MDecl:
			delta := d1.psi_recv.ToTEnv()
			for _, v := range d1.psi.tfs {
				delta[v.a] = v.u
			}
			gamma := make(Env)
			us := make([]Type, len(d1.psi_recv.tfs))
			for i := 0; i < len(us); i++ {
				us[i] = d1.psi_recv.tfs[i].a
			}
			gamma[d1.x_recv] = TName{d1.t_recv, us}
			for _, v := range d1.pds {
				gamma[v.x] = v.u
			}
			recv := fg.NewParamDecl(d1.x_recv, fg.Type(d1.t_recv))
			pds := make([]fg.ParamDecl, len(d1.pds))
			for i := 0; i < len(d1.pds); i++ {
				pd := d1.pds[i]
				pds[i] = fg.NewParamDecl(pd.x, fg.Type(erase(delta, pd.u)))
			}
			t := fg.Type(erase(delta, d1.u))
			e := wrap(p.ds, delta, gamma, d1.e, d1.u, wrappers) // TODO FIXME: subs ~alpha?
			md := fg.NewMDecl(recv, d1.m, pds, t, e)
			ds = append(ds, md)
		default:
			panic("Unexpected Decl type " + reflect.TypeOf(d).String() +
				": " + d.String())
		}
	}

	var delta TEnv // Empty envs for main -- duplicated from FGGProgram.OK
	var gamma Env
	p.e.Typing(p.ds, delta, gamma, false) // Populates delta and gamma
	e := translateExpr(p.ds, delta, gamma, p.e, wrappers)

	// Add wrappers, wrapper meths -- also getValue/getTypeRep (TODO: factor out with above)
	// Needs to follow translateExpr, for wrappers to be populated
	for k, v := range wrappers {
		fds := []fg.FieldDecl{fg.NewFieldDecl("value", v.sub)} // TODO: factor out
		adptr := fg.NewSTypeLit(k, fds)

		// TODO: factor out with STypeLits
		e_getv := fg.NewSelect(fg.NewVariable("x"), "value") // CHECKME: but t_S doesn't have value field, wrapper does?
		getv := fg.NewMDecl(fg.NewParamDecl("x", fg.Type(k)), "getValue",
			[]fg.ParamDecl{}, fg.Type("Any_0"), e_getv)
		// gettr := ...TODO...

		c := getTDecl(p.ds, string(v.super)).(ITypeLit)
		us := make([]Type, len(c.psi.tfs))
		for i := 0; i < len(us); i++ {
			us[i] = c.psi.tfs[i].a
		}
		dummy := TName{c.t, us}    // `us` are just the decl TParams, args not actually needed for `methods` or below
		gs := methods(p.ds, dummy) //map[Name]Sig

		//for _, s := range c.ss {
		for _, g := range gs {
			delta := make(TEnv)
			for _, v1 := range c.psi.tfs {
				delta[v1.a] = v1.u
			}
			for _, v1 := range g.psi.tfs {
				delta[v1.a] = v1.u
			}
			/*delta1 := make(TEnv)
			psi := getTDecl(p.ds, string(v.sub)).GetTFormals()
			for _, v1 := range psi.tfs {
				delta1[v1.a] = v1.u
			}
			for _, v1 := range g.psi.tfs {
				delta1[v1.a] = v1.u
			}*/

			//if g, ok := s.(Sig); ok { // TODO FIXME: need all meths in meth set (i.e., from embedded)
			pds := make([]fg.ParamDecl, len(g.pds))
			for i := 0; i < len(g.pds); i++ {
				pd := g.pds[i]
				pds[i] = fg.NewParamDecl(pd.x, fg.Type(erase(delta, pd.u)))
			}
			t := fg.Type(erase(delta, g.u))
			var e fg.Expr
			e = fg.NewStructLit("Dummy_0", []fg.Expr{})
			e = fg.NewStructLit("ToAny_0", []fg.Expr{e})
			e = fg.NewSelect(e, "any")
			e = fg.NewAssert(e, fg.Type(erase(delta, g.u)))
			md := fg.NewMDecl(fg.NewParamDecl("x", k), g.m,
				pds, t, e)
			ds = append(ds, md)
			//}
		}

		ds = append(ds, adptr, getv)
	}

	return fg.NewFGProgram(ds, e)
}

/* TDecl */

func translateSTypeLit(s STypeLit) fg.STypeLit {
	delta := s.psi.ToTEnv()
	fds := make([]fg.FieldDecl, len(s.fds)) // TODO FIXME: additional typerep fields
	for i := 0; i < len(s.fds); i++ {
		fd := s.fds[i]
		fds[i] = fg.NewFieldDecl(fd.f, fg.Type(erase(delta, fd.u)))
	}
	return fg.NewSTypeLit(fg.Type(s.t), fds)
}

func translateITypeLit(c ITypeLit) fg.ITypeLit {
	delta := c.psi.ToTEnv()
	ss := make([]fg.Spec, len(c.ss)+1)
	ss[0] = fg.Type("t_0") // TODO: factor out
	for i := 1; i <= len(c.ss); i++ {
		s := c.ss[i-1]
		switch s1 := s.(type) {
		case TName:
			ss[i] = fg.Type(s1.t)
		case Sig:
			ss[i] = translateSig(delta, s1)
		default:
			panic("Unknown Spec type " + reflect.TypeOf(s).String() + ": " + s.String())
		}
	}
	return fg.NewITypeLit(fg.Type(c.t), ss)
}

func translateSig(delta TEnv, g Sig) fg.Sig {
	delta1 := make(TEnv)
	for k, v := range delta {
		delta1[k] = v
	}
	for _, v := range g.psi.tfs {
		delta1[v.a] = v.u
	}
	pds := make([]fg.ParamDecl, len(g.pds))
	for i := 0; i < len(g.pds); i++ {
		pds[i] = fg.NewParamDecl(g.pds[i].x, fg.Type(erase(delta1, g.pds[i].u)))
	}
	t := fg.Type(erase(delta1, g.u))
	return fg.NewSig(g.m, pds, t)
}

/* Expr */

// |e_FGG|_(\Delta; \Gamma) = e_FGR
// TODO: rename
func translateExpr(ds []Decl, delta TEnv, gamma Env, e Expr, wrappers map[fg.Type]adaptPair) fg.Expr {
	switch e1 := e.(type) {
	case Variable:
		u := e1.Typing(ds, delta, gamma, false)
		var res fg.Expr
		res = fg.NewVariable(e1.id)
		if isInterfaceTName(ds, u) {
			// x.getValue().((mkRep u))
			res = fg.NewCall(res, Name("getValue"), []fg.Expr{})
			res = fg.NewAssert(res, fg.Type(erase(delta, u))) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case StructLit:
		t := e1.u.t
		es := make([]fg.Expr, len(e1.es)) // TODO FIXME: additional mkRep args
		fds := fields(ds, e1.u)
		subs := make(map[TParam]Type)
		psi := getTDecl(ds, t).GetTFormals()
		for i := 0; i < len(psi.tfs); i++ {
			subs[psi.tfs[i].a] = e1.u.us[i]
		}
		for i := 0; i < len(e1.es); i++ {
			u_i := fds[i].u
			es[i] = wrap(ds, delta, gamma, e1.es[i], u_i.TSubs(subs), wrappers)
		}
		return fg.NewStructLit(fg.Type(t), es)
	case Select:
		u := e1.Typing(ds, delta, gamma, false)
		var res fg.Expr
		res = fg.NewSelect(translateExpr(ds, delta, gamma, e1.e, wrappers), e1.f)
		if isInterfaceTName(ds, u) {
			// TODO FIXME: factor out with Variable
			res = fg.NewCall(res, Name("getValue"), []fg.Expr{})
			res = fg.NewAssert(res, fg.Type(erase(delta, u))) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case Call:
		u_recv := e1.e.Typing(ds, delta, gamma, false)
		g := methods(ds, bounds(delta, u_recv))[e1.m] // map[Name]Sig
		subs := make(map[TParam]Type)
		for i := 0; i < len(g.psi.tfs); i++ {
			subs[g.psi.tfs[i].a] = e1.targs[i]
		}
		args := make([]fg.Expr, len(e1.args))
		for i := 0; i < len(e1.args); i++ {
			args[i] = wrap(ds, delta, gamma, e1.args[i], g.pds[i].u.TSubs(subs), wrappers)
		}
		u := e1.Typing(ds, delta, gamma, false)
		e_recv := translateExpr(ds, delta, gamma, e1.e, wrappers)
		var res fg.Expr
		res = fg.NewCall(e_recv, e1.m, args)
		if isInterfaceTName(ds, u) {
			// TODO FIXME: factor out with Variable
			res = fg.NewCall(res, Name("getValue"), []fg.Expr{})
			res = fg.NewAssert(res, fg.Type(erase(delta, u))) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case Assert:
		panic("TODO " + reflect.TypeOf(e).String() + ": " + e.String())
	default:
		panic("Unknown Expr type " + reflect.TypeOf(e).String() + ": " + e.String())
	}
}

/* Aux */

// |\tau|_\Delta = t
func erase(delta TEnv, u Type) Name { //fg.Type {  // CHECKME: change return back to fg.Type?
	return bounds(delta, u).(TName).t
}

// Pre: type of e <: u
// `u` is "target type"
func wrap(ds []fg.Decl, delta TEnv, gamma Env, e Expr, u Type, wrappers map[fg.Type]adaptPair) fg.Expr {
	/*t := erase(u, delta)
	if _, ok := fg.isStructType(t)*/
	if isStructTName(ds, u) { // N.B. differs slightly from def -- because there is no FG t_S decl (yet)?
		return translateExpr(ds, delta, gamma, e, wrappers)
	} else if isInterfaceTName(ds, u) {
		u1 := e.Typing(ds, delta, gamma, false)
		e1 := translateExpr(ds, delta, gamma, e, wrappers)
		return wrapper(delta, e1, u1, u, wrappers)
	} else {
		panic("Invalid wrap case: e=" + e.String() + ", u=" + u.String())
	}
}

// targ is a t_I
// TODO: rename, cf. wrap(ds, delta, gamma, e, u)
func wrapper(delta TEnv, e fg.Expr, subj Type, targ Type, wrappers map[fg.Type]adaptPair) fg.StructLit {
	t1 := fg.Type(erase(delta, subj))
	t_I := fg.Type(erase(delta, targ))
	adptr := fg.Type("Adptr_" + t1 + "_" + t_I) // TODO: factor out naming
	if _, ok := wrappers[adptr]; !ok {
		wrappers[adptr] = adaptPair{t1, t_I}
	}
	return fg.NewStructLit(adptr, []fg.Expr{e})
}
