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

type fgAdaptPair struct {
	sub   fg.Type
	super fg.Type // The "target" type, a t_I
}

func FgAdptrTranslate(p FGGProgram) fg.FGProgram { // TODO FIXME: FGR -- TODO also can subsume existing FGG-FG trans?
	var ds []Decl

	// Add t_0 (etc.) to ds
	// TODO: factor out constants
	Any_0 := fg.NewITypeLit(fg.Type("Any_0"), []fg.Spec{})
	Dummy_0 := fg.NewSTypeLit(fg.Type("Dummy_0"), []fg.FieldDecl{})
	ToAny_0 := fg.NewSTypeLit(fg.Type("ToAny_0"), []fg.FieldDecl{fg.NewFieldDecl("any", fg.Type("Any_0"))})
	ds = append(ds, Any_0, Dummy_0, ToAny_0)

	getValue := fg.NewSig("getValue", []fg.ParamDecl{}, fg.Type("Any_0")) // TODO: rename "unfgWrap"?
	//getTypeRep := fg.NewSig("getTypeRep", []fg.ParamDecl{}, fg.Type("...TODO..."))
	ss_0 := []fg.Spec{getValue}
	t_0 := fg.NewITypeLit(fg.Type("t_0"), ss_0) // TODO FIXME? Go doesn't allow "overlapping" interfaces
	ds = append(ds, t_0)

	fgWrappers := make(map[fg.Type]fgAdaptPair) // Populated by visiting MDecl and main Expr

	for i := 0; i < len(p.ds); i++ {
		d := p.ds[i]
		switch d1 := d.(type) {
		case STypeLit:
			s := fgAdaptSTypeLit(d1)

			// Add getValue/getTypeRep to all (existing) t_S -- every t_S must implement t_0 -- TODO: factor out with fgWrappers
			//e_getv := fg.NewSelect(fg.NewVariable("x"), "value") // CHECKME: but t_S doesn't have value field, fgWrapper does?
			e_getv := fg.NewStructLit(fg.Type("Dummy_0"), []fg.Expr{})
			getv := fg.NewMDecl(fg.NewParamDecl("x", fg.Type(d1.t)), "getValue",
				[]fg.ParamDecl{}, fg.Type("Any_0"), e_getv)
			// gettr := ...TODO...

			ds = append(ds, s, getv)
		case ITypeLit:
			ds = append(ds, fgAdaptITypeLit(d1))
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
			gamma[d1.x_recv] = TName{d1.t_recv, us} // !!! also receiver
			for _, v := range d1.pds {
				gamma[v.x] = v.u
			}
			recv := fg.NewParamDecl(d1.x_recv, fg.Type(d1.t_recv))
			pds := make([]fg.ParamDecl, len(d1.pds))
			for i := 0; i < len(d1.pds); i++ {
				pd := d1.pds[i]
				pds[i] = fg.NewParamDecl(pd.x, fg.Type(fgErase(delta, pd.u)))
			}
			t := fg.Type(fgErase(delta, d1.u)) // !!! tau_p typo
			u := bounds(delta, d1.u)           // !!! cf. fgWrap, pre: u is a TName
			e := fgWrap(p.ds, delta, gamma, d1.e, u, fgWrappers)
			// ^TODO FIXME: subs ~alpha?
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
	e := fgAdaptExpr(p.ds, delta, gamma, p.e, fgWrappers)

	// Add fgWrappers, fgWrapper meths -- also getValue/getTypeRep (TODO: factor out with above)
	// Needs to follow fgAdaptExpr, for fgWrappers to be populated
	for k, v := range fgWrappers {
		fds := []fg.FieldDecl{fg.NewFieldDecl("value", v.sub)} // TODO: factor out
		adptr := fg.NewSTypeLit(k, fds)

		// TODO: factor out with STypeLits
		e_getv := fg.NewSelect(fg.NewVariable("x"), "value") // CHECKME: but t_S doesn't have value field, fgWrapper does?
		getv := fg.NewMDecl(fg.NewParamDecl("x", fg.Type(k)), "getValue",
			[]fg.ParamDecl{}, fg.Type("Any_0"), e_getv)
		// gettr := ...TODO...

		c := getTDecl(p.ds, string(v.super)).(ITypeLit)
		us := make([]Type, len(c.psi.tfs))
		for i := 0; i < len(us); i++ {
			us[i] = c.psi.tfs[i].a
		}
		dummy := TName{c.t, us}    // `us` are just the decl TParams, args not actually needed for `methods` or below
		gs := methods(p.ds, dummy) // !!! all meths of t_I target

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

			//if g, ok := s.(Sig); ok { // !!! need all meths in meth set (i.e., from embedded)
			pds := make([]fg.ParamDecl, len(g.pds))
			for i := 0; i < len(g.pds); i++ {
				pd := g.pds[i]
				pds[i] = fg.NewParamDecl(pd.x, fg.Type(fgErase(delta, pd.u)))
			}
			t := fg.Type(fgErase(delta, g.u)) // !!! tau_p typo, and delta'?
			var e fg.Expr
			e = fg.NewStructLit("Dummy_0", []fg.Expr{})
			e = fg.NewStructLit("ToAny_0", []fg.Expr{e})
			e = fg.NewSelect(e, "any")
			e = fg.NewAssert(e, fg.Type(fgErase(delta, g.u)))
			md := fg.NewMDecl(fg.NewParamDecl("x", k), g.m,
				pds, t, e)
			ds = append(ds, md)
			//}
		}

		ds = append(ds, adptr, getv)
	}

	return fg.NewFGProgram(ds, e, false) // FIXME: printf
}

/* TDecl */

func fgAdaptSTypeLit(s STypeLit) fg.STypeLit {
	delta := s.psi.ToTEnv()
	fds := make([]fg.FieldDecl, len(s.fds)) // TODO FIXME: additional typerep fields
	for i := 0; i < len(s.fds); i++ {
		fd := s.fds[i]
		fds[i] = fg.NewFieldDecl(fd.f, fg.Type(fgErase(delta, fd.u)))
	}
	return fg.NewSTypeLit(fg.Type(s.t), fds)
}

func fgAdaptITypeLit(c ITypeLit) fg.ITypeLit {
	delta := c.psi.ToTEnv()
	ss := make([]fg.Spec, len(c.ss)+1)
	ss[0] = fg.Type("t_0") // TODO: factor out
	for i := 1; i <= len(c.ss); i++ {
		s := c.ss[i-1]
		switch s1 := s.(type) {
		case TName:
			ss[i] = fg.Type(s1.t)
		case Sig:
			ss[i] = fgAdaptSig(delta, s1)
		default:
			panic("Unknown Spec type " + reflect.TypeOf(s).String() + ": " + s.String())
		}
	}
	return fg.NewITypeLit(fg.Type(c.t), ss)
}

func fgAdaptSig(delta TEnv, g Sig) fg.Sig {
	delta1 := make(TEnv)
	for k, v := range delta {
		delta1[k] = v
	}
	for _, v := range g.psi.tfs {
		delta1[v.a] = v.u
	}
	pds := make([]fg.ParamDecl, len(g.pds))
	for i := 0; i < len(g.pds); i++ {
		pds[i] = fg.NewParamDecl(g.pds[i].x, fg.Type(fgErase(delta1, g.pds[i].u)))
	}
	t := fg.Type(fgErase(delta1, g.u))
	return fg.NewSig(g.m, pds, t)
}

/* Expr */

// |e_FGG|_(\Delta; \Gamma) = e_FGR
// TODO: rename
func fgAdaptExpr(ds []Decl, delta TEnv, gamma Env, e Expr, fgWrappers map[fg.Type]fgAdaptPair) fg.Expr {
	switch e1 := e.(type) {
	case Variable:
		u := e1.Typing(ds, delta, gamma, false)
		var res fg.Expr
		res = fg.NewVariable(e1.id)
		if isInterfaceTName(ds, u) {
			// x.getValue().((mkRep u))
			res = fg.NewCall(res, Name("getValue"), []fg.Expr{})
			res = fg.NewAssert(res, fg.Type(fgErase(delta, u))) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case StructLit:
		t := e1.u.t
		es := make([]fg.Expr, len(e1.es)) // TODO FIXME: additional mkRep args
		fds := fields(ds, e1.u)
		subs := make(map[TParam]Type)
		psi := getTDecl(ds, t).GetTFormals()
		for i := 0; i < len(psi.tfs); i++ {
			//subs[psi.tfs[i].a] = e1.u.us[i]
			// !!! type arg may be a TParam (when visiting MDecl), e.g., Box(a){...}
			subs[psi.tfs[i].a] = bounds(delta, e1.u.us[i])
		}
		for i := 0; i < len(e1.es); i++ {
			u_i := fds[i].u.TSubs(subs)
			es[i] = fgWrap(ds, delta, gamma, e1.es[i], u_i, fgWrappers)
		}
		return fg.NewStructLit(fg.Type(t), es)
	case Select:
		u := e1.Typing(ds, delta, gamma, false)
		var res fg.Expr
		res = fg.NewSelect(fgAdaptExpr(ds, delta, gamma, e1.e, fgWrappers), e1.f)
		if isInterfaceTName(ds, u) {
			// TODO FIXME: factor out with Variable
			res = fg.NewCall(res, Name("getValue"), []fg.Expr{})
			res = fg.NewAssert(res, fg.Type(fgErase(delta, u))) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case Call:
		u_recv := e1.e.Typing(ds, delta, gamma, false)
		g := methods(ds, bounds(delta, u_recv))[e1.m]
		subs := make(map[TParam]Type)
		for i := 0; i < len(g.psi.tfs); i++ {
			//subs[g.psi.tfs[i].a] = e1.targs[i]  // !!! cf. StructLit case
			subs[g.psi.tfs[i].a] = bounds(delta, e1.targs[i])
		}
		args := make([]fg.Expr, len(e1.args))
		for i := 0; i < len(e1.args); i++ {
			u_i := g.pds[i].u.TSubs(subs)
			args[i] = fgWrap(ds, delta, gamma, e1.args[i], u_i, fgWrappers)
		}
		//u := e1.Typing(ds, delta, gamma, false)
		e_recv := fgAdaptExpr(ds, delta, gamma, e1.e, fgWrappers)
		var res fg.Expr
		res = fg.NewCall(e_recv, e1.m, args)
		//if isInterfaceTName(ds, fgErase(delta, u)) {  // !!! fgErase returns fg.Type (cf. fgWrap isStructTName)

		//fmt.Println("aaa:", e1, u, bounds(delta, u))

		//u_ret := u.TSubs(subs) // Cf. bounds(delta, u) ?

		delta1 := make(map[TParam]Type)
		for i := 0; i < len(g.psi.tfs); i++ {
			tf := g.psi.tfs[i]
			delta1[tf.a] = tf.u
		}
		td := getTDecl(ds, bounds(delta, u_recv).(TName).t)
		psi := td.GetTFormals()
		for i := 0; i < len(psi.tfs); i++ {
			tf := psi.tfs[i]
			delta1[tf.a] = tf.u
		}

		//u_ret := g.u.TSubs(delta1)
		u_ret := fgErase(delta1, g.u)

		//if isInterfaceTName(ds, u_ret) {
		//if _, ok := u_ret.(TParam); ok || isInterfaceTName(ds, u_ret) {
		//if isInterfaceTName(ds, u_ret) {
		if !isSTypeLit(ds, u_ret) {
			// TODO FIXME: factor out with Variable
			res = fg.NewCall(res, Name("getValue"), []fg.Expr{})
			//res = fg.NewAssert(res, fg.Type(fgErase(delta, u_ret))) // TODO FIXME: mkRep -- "FG" for now, not FGR
			res = fg.NewAssert(res, fg.Type(u_ret))
		}
		return res
	case Assert:
		// Need actual FGR
		panic("TODO " + reflect.TypeOf(e).String() + ": " + e.String())
	default:
		panic("Unknown Expr type " + reflect.TypeOf(e).String() + ": " + e.String())
	}
}

/* Aux */

// |\tau|_\Delta = t
func fgErase(delta TEnv, u Type) Name { //fg.Type {  // CHECKME: change return back to fg.Type?
	return bounds(delta, u).(TName).t
}

// Pre: u is a TName  // !!!
// Pre: type of e <: u
// `u` is "target type"
func fgWrap(ds []fg.Decl, delta TEnv, gamma Env, e Expr, u Type, fgWrappers map[fg.Type]fgAdaptPair) fg.Expr {
	/*t := fgErase(u, delta)
	if _, ok := fg.isStructType(t)*/
	if isStructTName(ds, u) { // !!! differs slightly from def -- because there is no FG t_S decl (yet)?
		return fgAdaptExpr(ds, delta, gamma, e, fgWrappers)
	} else if isInterfaceTName(ds, u) {
		u1 := e.Typing(ds, delta, gamma, false)
		e1 := fgAdaptExpr(ds, delta, gamma, e, fgWrappers)
		return fgWrapper(delta, e1, u1, u, fgWrappers)
	} else {
		panic("Invalid fgWrap case: e=" + e.String() + ", u=" + u.String())
	}
}

// targ is a t_I
// TODO: rename, cf. fgWrap(ds, delta, gamma, e, u)
func fgWrapper(delta TEnv, e fg.Expr, subj Type, targ Type, fgWrappers map[fg.Type]fgAdaptPair) fg.StructLit {
	t1 := fg.Type(fgErase(delta, subj))
	t_I := fg.Type(fgErase(delta, targ))
	adptr := fg.Type("Adptr_" + t1 + "_" + t_I) // TODO: factor out naming
	if _, ok := fgWrappers[adptr]; !ok {
		fgWrappers[adptr] = fgAdaptPair{t1, t_I}
	}
	return fg.NewStructLit(adptr, []fg.Expr{e})
}

/* Helper */

func isSTypeLit(ds []Decl, t Name) bool {
	for _, v := range ds {
		if _, ok := v.(STypeLit); ok && v.GetName() == string(t) {
			return true
		}
	}
	return false
}
