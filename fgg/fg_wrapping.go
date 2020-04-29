package fgg

import (
	"fmt"
	"reflect"
	//"strings"

	"github.com/rhu1/fgg/fg"
	//"github.com/rhu1/fgg/fgg"
)

/**
 * !!! Deprecated !!!
 */

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

	for i := 0; i < len(p.decls); i++ {
		d := p.decls[i]
		switch d1 := d.(type) {
		case STypeLit:
			s := fgAdaptSTypeLit(d1)

			// Add getValue/getTypeRep to all (existing) t_S -- every t_S must implement t_0 -- TODO: factor out with fgWrappers
			//e_getv := fg.NewSelect(fg.NewVariable("x"), "value") // CHECKME: but t_S doesn't have value field, fgWrapper does?
			e_getv := fg.NewStructLit(fg.Type("Dummy_0"), []fg.FGExpr{})
			getv := fg.NewMDecl(fg.NewParamDecl("x", fg.Type(d1.t_name)), "getValue",
				[]fg.ParamDecl{}, fg.Type("Any_0"), e_getv)
			// gettr := ...TODO...

			ds = append(ds, s, getv)
		case ITypeLit:
			ds = append(ds, fgAdaptITypeLit(d1))
		case MDecl:
			delta := d1.PsiRecv.ToDelta()
			for _, v := range d1.PsiMeth.tFormals {
				delta[v.name] = v.u_I
			}
			gamma := make(Gamma)
			us := make([]Type, len(d1.PsiRecv.tFormals))
			for i := 0; i < len(us); i++ {
				us[i] = d1.PsiRecv.tFormals[i].name
			}
			gamma[d1.x_recv] = TNamed{d1.t_recv, us} // !!! also receiver
			for _, v := range d1.pDecls {
				gamma[v.name] = v.u
			}
			recv := fg.NewParamDecl(d1.x_recv, fg.Type(d1.t_recv))
			pds := make([]fg.ParamDecl, len(d1.pDecls))
			for i := 0; i < len(d1.pDecls); i++ {
				pd := d1.pDecls[i]
				pds[i] = fg.NewParamDecl(pd.name, fg.Type(fgErase(delta, pd.u)))
			}
			t := fg.Type(fgErase(delta, d1.u_ret)) // !!! tau_p typo
			u := bounds(delta, d1.u_ret)           // !!! cf. fgWrap, pre: u is a TName
			e := fgWrap(p.decls, delta, gamma, d1.e_body, u, fgWrappers)
			// ^TODO FIXME: subs ~alpha?
			md := fg.NewMDecl(recv, d1.name, pds, t, e)
			ds = append(ds, md)
		default:
			panic("Unexpected Decl type " + reflect.TypeOf(d).String() +
				": " + d.String())
		}
	}

	var delta Delta // Empty envs for main -- duplicated from FGGProgram.OK
	var gamma Gamma
	p.e_main.Typing(p.decls, delta, gamma, false) // Populates delta and gamma
	e := fgAdaptExpr(p.decls, delta, gamma, p.e_main, fgWrappers)

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

		c := GetTDecl(p.decls, string(v.super)).(ITypeLit)
		us := make([]Type, len(c.psi.tFormals))
		for i := 0; i < len(us); i++ {
			us[i] = c.psi.tFormals[i].name
		}
		dummy := TNamed{c.t_I, us}    // `us` are just the decl TParams, args not actually needed for `methods` or below
		gs := methods(p.decls, dummy) // !!! all meths of t_I target

		//for _, s := range c.ss {
		for _, g := range gs {
			delta := make(Delta)
			for _, v1 := range c.psi.tFormals {
				delta[v1.name] = v1.u_I
			}
			for _, v1 := range g.psi.tFormals {
				delta[v1.name] = v1.u_I
			}
			/*delta1 := make(TEnv)
			psi := GetTDecl(p.ds, string(v.sub)).GetTFormals()
			for _, v1 := range psi.tfs {
				delta1[v1.a] = v1.u
			}
			for _, v1 := range g.psi.tfs {
				delta1[v1.a] = v1.u
			}*/

			//if g, ok := s.(Sig); ok { // !!! need all meths in meth set (i.e., from embedded)
			pds := make([]fg.ParamDecl, len(g.pDecls))
			for i := 0; i < len(g.pDecls); i++ {
				pd := g.pDecls[i]
				pds[i] = fg.NewParamDecl(pd.name, fg.Type(fgErase(delta, pd.u)))
			}
			t := fg.Type(fgErase(delta, g.u_ret)) // !!! tau_p typo, and delta'?
			var e fg.FGExpr
			e = fg.NewStructLit("Dummy_0", []fg.FGExpr{})
			e = fg.NewStructLit("ToAny_0", []fg.FGExpr{e})
			e = fg.NewSelect(e, "any")
			e = fg.NewAssert(e, fg.Type(fgErase(delta, g.u_ret)))
			md := fg.NewMDecl(fg.NewParamDecl("x", k), g.meth,
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
	delta := s.Psi.ToDelta()
	fds := make([]fg.FieldDecl, len(s.fDecls)) // TODO FIXME: additional typerep fields
	for i := 0; i < len(s.fDecls); i++ {
		fd := s.fDecls[i]
		fds[i] = fg.NewFieldDecl(fd.field, fg.Type(fgErase(delta, fd.u)))
	}
	return fg.NewSTypeLit(fg.Type(s.t_name), fds)
}

func fgAdaptITypeLit(c ITypeLit) fg.ITypeLit {
	delta := c.psi.ToDelta()
	ss := make([]fg.Spec, len(c.specs)+1)
	ss[0] = fg.Type("t_0") // TODO: factor out
	for i := 1; i <= len(c.specs); i++ {
		s := c.specs[i-1]
		switch s1 := s.(type) {
		case TNamed:
			ss[i] = fg.Type(s1.t_name)
		case Sig:
			ss[i] = fgAdaptSig(delta, s1)
		default:
			panic("Unknown Spec type " + reflect.TypeOf(s).String() + ": " + s.String())
		}
	}
	return fg.NewITypeLit(fg.Type(c.t_I), ss)
}

func fgAdaptSig(delta Delta, g Sig) fg.Sig {
	delta1 := make(Delta)
	for k, v := range delta {
		delta1[k] = v
	}
	for _, v := range g.psi.tFormals {
		delta1[v.name] = v.u_I
	}
	pds := make([]fg.ParamDecl, len(g.pDecls))
	for i := 0; i < len(g.pDecls); i++ {
		pds[i] = fg.NewParamDecl(g.pDecls[i].name, fg.Type(fgErase(delta1, g.pDecls[i].u)))
	}
	t := fg.Type(fgErase(delta1, g.u_ret))
	return fg.NewSig(g.meth, pds, t)
}

/* Expr */

// |e_FGG|_(\Delta; \Gamma) = e_FGR
// TODO: rename
func fgAdaptExpr(ds []Decl, delta Delta, gamma Gamma, e FGGExpr, fgWrappers map[fg.Type]fgAdaptPair) fg.FGExpr {
	switch e1 := e.(type) {
	case Variable:
		u := e1.Typing(ds, delta, gamma, false)
		var res fg.FGExpr
		res = fg.NewVariable(e1.name)
		if IsNamedIfaceType(ds, u) {
			// x.getValue().((mkRep u))
			res = fg.NewCall(res, Name("getValue"), []fg.FGExpr{})
			res = fg.NewAssert(res, fg.Type(fgErase(delta, u))) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case StructLit:
		t := e1.u_S.t_name
		es := make([]fg.FGExpr, len(e1.elems)) // TODO FIXME: additional mkRep args
		fds := fields(ds, e1.u_S)
		subs := make(map[TParam]Type)
		psi := GetTDecl(ds, t).GetBigPsi()
		for i := 0; i < len(psi.tFormals); i++ {
			//subs[psi.tfs[i].a] = e1.u.us[i]
			// !!! type arg may be a TParam (when visiting MDecl), e.g., Box(a){...}
			subs[psi.tFormals[i].name] = bounds(delta, e1.u_S.u_args[i])
		}
		for i := 0; i < len(e1.elems); i++ {
			u_i := fds[i].u.TSubs(subs)
			es[i] = fgWrap(ds, delta, gamma, e1.elems[i], u_i, fgWrappers)
		}
		return fg.NewStructLit(fg.Type(t), es)
	case Select:
		u := e1.Typing(ds, delta, gamma, false)
		var res fg.FGExpr
		res = fg.NewSelect(fgAdaptExpr(ds, delta, gamma, e1.e_S, fgWrappers), e1.field)
		if IsNamedIfaceType(ds, u) {
			// TODO FIXME: factor out with Variable
			res = fg.NewCall(res, Name("getValue"), []fg.FGExpr{})
			res = fg.NewAssert(res, fg.Type(fgErase(delta, u))) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case Call:
		u_recv := e1.e_recv.Typing(ds, delta, gamma, false)
		g := methods(ds, bounds(delta, u_recv))[e1.meth]
		subs := make(map[TParam]Type)
		for i := 0; i < len(g.psi.tFormals); i++ {
			//subs[g.psi.tfs[i].a] = e1.targs[i]  // !!! cf. StructLit case
			subs[g.psi.tFormals[i].name] = bounds(delta, e1.t_args[i])
		}
		args := make([]fg.FGExpr, len(e1.args))
		for i := 0; i < len(e1.args); i++ {
			u_i := g.pDecls[i].u.TSubs(subs)
			args[i] = fgWrap(ds, delta, gamma, e1.args[i], u_i, fgWrappers)
		}
		//u := e1.Typing(ds, delta, gamma, false)
		e_recv := fgAdaptExpr(ds, delta, gamma, e1.e_recv, fgWrappers)
		var res fg.FGExpr
		res = fg.NewCall(e_recv, e1.meth, args)
		//if isInterfaceTName(ds, fgErase(delta, u)) {  // !!! fgErase returns fg.Type (cf. fgWrap isStructTName)

		//fmt.Println("aaa:", e1, u, bounds(delta, u))

		//u_ret := u.TSubs(subs) // Cf. bounds(delta, u) ?

		delta1 := make(map[TParam]Type)
		for i := 0; i < len(g.psi.tFormals); i++ {
			tf := g.psi.tFormals[i]
			delta1[tf.name] = tf.u_I
		}
		td := GetTDecl(ds, bounds(delta, u_recv).(TNamed).t_name)
		psi := td.GetBigPsi()
		for i := 0; i < len(psi.tFormals); i++ {
			tf := psi.tFormals[i]
			delta1[tf.name] = tf.u_I
		}

		//u_ret := g.u.TSubs(delta1)
		u_ret := fgErase(delta1, g.u_ret)

		//if isInterfaceTName(ds, u_ret) {
		//if _, ok := u_ret.(TParam); ok || isInterfaceTName(ds, u_ret) {
		//if isInterfaceTName(ds, u_ret) {
		if !isSTypeLit(ds, u_ret) {
			// TODO FIXME: factor out with Variable
			res = fg.NewCall(res, Name("getValue"), []fg.FGExpr{})
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
func fgErase(delta Delta, u Type) Name { //fg.Type {  // CHECKME: change return back to fg.Type?
	return bounds(delta, u).(TNamed).t_name
}

// Pre: u is a TName  // !!!
// Pre: type of e <: u
// `u` is "target type"
func fgWrap(ds []fg.Decl, delta Delta, gamma Gamma, e FGGExpr, u Type, fgWrappers map[fg.Type]fgAdaptPair) fg.FGExpr {
	/*t := fgErase(u, delta)
	if _, ok := fg.IsStructType(t)*/
	if IsStructType(ds, u) { // !!! differs slightly from def -- because there is no FG t_S decl (yet)?
		return fgAdaptExpr(ds, delta, gamma, e, fgWrappers)
	} else if IsNamedIfaceType(ds, u) {
		u1 := e.Typing(ds, delta, gamma, false)
		e1 := fgAdaptExpr(ds, delta, gamma, e, fgWrappers)
		return fgWrapper(delta, e1, u1, u, fgWrappers)
	} else {
		panic("Invalid fgWrap case: e=" + e.String() + ", u=" + u.String())
	}
}

// targ is a t_I
// TODO: rename, cf. fgWrap(ds, delta, gamma, e, u)
func fgWrapper(delta Delta, e fg.FGExpr, subj Type, targ Type, fgWrappers map[fg.Type]fgAdaptPair) fg.StructLit {
	t1 := fg.Type(fgErase(delta, subj))
	t_I := fg.Type(fgErase(delta, targ))
	adptr := fg.Type("Adptr_" + t1 + "_" + t_I) // TODO: factor out naming
	if _, ok := fgWrappers[adptr]; !ok {
		fgWrappers[adptr] = fgAdaptPair{t1, t_I}
	}
	return fg.NewStructLit(adptr, []fg.FGExpr{e})
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
