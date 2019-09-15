package fgr

import (
	"fmt"
	"reflect"
	//"strings"

	"github.com/rhu1/fgg/fgg"
)

var _ = fmt.Errorf

/* FGGProgram */

func Translate(p fgg.FGGProgram) FGRProgram { // TODO FIXME: FGR -- TODO also can subsume existing FGG-FG trans?
	var ds_fgr []Decl

	// Add t_0 (etc.) to ds_fgr
	// TODO: factor out constants
	Any_0 := NewITypeLit(Type("Any_0"), []Spec{})
	Dummy_0 := NewSTypeLit(Type("Dummy_0"), []FieldDecl{})
	ToAny_0 := NewSTypeLit(Type("ToAny_0"), []FieldDecl{NewFieldDecl("any", Type("Any_0"))})
	getValue := NewSig("getValue", []ParamDecl{}, Type("Any_0")) // TODO: rename "unwrap"?
	//getTypeRep := fg.NewSig("getTypeRep", []fg.ParamDecl{}, fg.Type("...TODO..."))
	ss_0 := []Spec{getValue}
	t_0 := NewITypeLit(Type("t_0"), ss_0) // TODO FIXME? Go doesn't allow "overlapping" interfaces
	ds_fgr = append(ds_fgr, Any_0, Dummy_0, ToAny_0, t_0)

	wrappers := make(map[Type]adptrPair) // Populated by visiting MDecl and main Expr

	// Translate Decls (and collect wrappers from MDecls)
	ds_fgg := p.GetDecls()
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
				[]ParamDecl{}, Type("Any_0"), e_getv)
			// gettr := ...TODO...

			ds_fgr = append(ds_fgr, s, getv)
		case fgg.ITypeLit:
			ds_fgr = append(ds_fgr, fgrTransITypeLit(d1))
		case fgg.MDecl:
			ds_fgr = append(ds_fgr, fgrTransMDecl(ds_fgg, d1, wrappers))
		default:
			panic("Unexpected Decl type " + reflect.TypeOf(d).String() +
				": " + d.String())
		}
	}

	// Translate main Expr (and collect wrappers)
	e_fgg := p.GetExpr().(fgg.Expr)
	var delta fgg.TEnv // Empty envs for main -- duplicated from FGGProgram.OK
	var gamma fgg.Env
	e_fgg.Typing(ds_fgg, delta, gamma, false) // Populates delta and gamma
	e := fgrTransExpr(ds_fgg, delta, gamma, e_fgg, wrappers)

	// Add wrappers: Adptr types, getValue/getTypeRep (TODO: factor out with above) and duck-typing meths
	// Pre: fgrTransExpr, for wrappers to be populated
	for k, v := range wrappers {
		// Add Adptr and getValue/getTypeRep
		fds := []FieldDecl{NewFieldDecl("value", v.sub)} // TODO: factor out
		adptr := NewSTypeLit(k, fds)
		// TODO: factor out with STypeLits
		e_getv := NewSelect(NewVariable("x"), "value") // CHECKME: but t_S doesn't have value field, wrapper does?
		getv := NewMDecl(NewParamDecl("x", Type(k)), "getValue",
			[]ParamDecl{}, Type("Any_0"), e_getv)
		// gettr := ...TODO...
		ds_fgr = append(ds_fgr, adptr, getv)

		// Add duck-typing meths
		c := fgg.GetTDecl1(ds_fgg, string(v.super)).(fgg.ITypeLit)
		tfs := c.GetTFormals().GetFormals()
		us := make([]fgg.Type, len(tfs))
		for i := 0; i < len(us); i++ {
			us[i] = tfs[i].GetTParam()
		}
		dummy := fgg.NewTName(c.GetName(), us) // `us` are just the decl TParams, args not actually needed for `methods` or below
		gs := fgg.Methods1(ds_fgg, dummy)      // !!! all meths of t_I target
		//for _, s := range c.ss {
		for _, g := range gs {
			delta := make(fgg.TEnv)
			for _, v1 := range tfs {
				delta[v1.GetTParam()] = v1.GetType()
			}
			tfs := g.GetTFormals().GetFormals()
			for _, v1 := range tfs {
				delta[v1.GetTParam()] = v1.GetType()
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
			pds_fgg := g.GetParamDecls()
			pds := make([]ParamDecl, len(pds_fgg))
			for i := 0; i < len(pds_fgg); i++ {
				pd := pds_fgg[i]
				pds[i] = NewParamDecl(pd.GetName(), toFgTypeFromBounds(delta, pd.GetType()))
			}
			u := g.GetType()
			t := toFgTypeFromBounds(delta, u) // !!! tau_p typo, and delta'?
			var e Expr
			e = NewStructLit("Dummy_0", []Expr{})
			e = NewStructLit("ToAny_0", []Expr{e})
			e = NewSelect(e, "any")
			e = NewAssert(e, toFgTypeFromBounds(delta, u))
			md := NewMDecl(NewParamDecl("x", k), g.GetName(),
				pds, t, e)
			ds_fgr = append(ds_fgr, md)
			//}
		}
	}

	return NewFGRProgram(ds_fgr, e)
}

/* TDecl */

func fgrTransSTypeLit(s fgg.STypeLit) STypeLit {
	delta := s.GetTFormals().ToTEnv()
	fds_fgg := s.GetFieldDecls()
	fds := make([]FieldDecl, len(fds_fgg)) // TODO FIXME: additional typerep fields
	for i := 0; i < len(fds_fgg); i++ {
		fd := fds_fgg[i]
		fds[i] = NewFieldDecl(fd.GetName(), toFgTypeFromBounds(delta, fd.GetType()))
	}
	return NewSTypeLit(Type(s.GetName()), fds)
}

func fgrTransITypeLit(c fgg.ITypeLit) ITypeLit {
	delta := c.GetTFormals().ToTEnv()
	ss_fgg := c.GetSpecs()
	ss := make([]Spec, len(ss_fgg)+1)
	ss[0] = Type("t_0") // TODO: factor out
	for i := 1; i <= len(ss_fgg); i++ {
		s := ss_fgg[i-1]
		switch s1 := s.(type) {
		case fgg.TName:
			ss[i] = Type(s1.GetName())
		case fgg.Sig:
			ss[i] = fgrTransSig(delta, s1)
		default:
			panic("Unknown Spec type " + reflect.TypeOf(s).String() + ": " + s.String())
		}
	}
	return NewITypeLit(Type(c.GetName()), ss)
}

func fgrTransSig(delta fgg.TEnv, g fgg.Sig) Sig {
	delta1 := make(fgg.TEnv)
	for k, v := range delta {
		delta1[k] = v
	}
	for _, v := range g.GetTFormals().GetFormals() {
		delta1[v.GetTParam()] = v.GetType()
	}
	pds_fgg := g.GetParamDecls()
	pds := make([]ParamDecl, len(pds_fgg))
	for i := 0; i < len(pds_fgg); i++ {
		pds[i] = NewParamDecl(pds_fgg[i].GetName(), toFgTypeFromBounds(delta1, pds_fgg[i].GetType()))
	}
	t := toFgTypeFromBounds(delta1, g.GetType())
	return NewSig(g.GetName(), pds, t)
}

/* MDecl */

func fgrTransMDecl(ds []Decl, d1 fgg.MDecl, wrappers map[Type]adptrPair) MDecl {
	x_recv := d1.GetRecvName()
	t_recv := d1.GetRecvTypeName()
	psi_recv := d1.GetRecvTFormals()
	m := d1.GetName()
	psi := d1.GetTFormals()
	pds_fgg := d1.GetParamDecls()
	u := d1.GetReturn()
	e_fgg := d1.GetExpr()

	delta := psi_recv.ToTEnv()
	tfs := psi.GetFormals()
	for _, v := range tfs {
		delta[v.GetTParam()] = v.GetType()
	}
	gamma := make(fgg.Env)
	us := make([]fgg.Type, len(tfs))
	for i := 0; i < len(us); i++ {
		us[i] = tfs[i].GetTParam()
	}
	gamma[x_recv] = fgg.NewTName(t_recv, us) // !!! also receiver
	for _, v := range pds_fgg {
		gamma[v.GetName()] = v.GetType()
	}
	recv := NewParamDecl(x_recv, Type(t_recv))
	pds := make([]ParamDecl, len(pds_fgg))
	for i := 0; i < len(pds_fgg); i++ {
		pd := pds_fgg[i]
		pds[i] = NewParamDecl(pd.GetName(), toFgTypeFromBounds(delta, pd.GetType()))
	}
	t := toFgTypeFromBounds(delta, u)                   // !!! tau_p typo
	e := wrapExpr(ds, delta, gamma, e_fgg, u, wrappers) // TODO FIXME: subs ~alpha?
	return NewMDecl(recv, m, pds, t, e)
}

/* Expr */

// TODO: rename ds -> ds_fgg
// |e_FGG|_(\Delta; \Gamma) = e_FGR
func fgrTransExpr(ds []Decl, delta fgg.TEnv, gamma fgg.Env, e fgg.Expr,
	wrappers map[Type]adptrPair) Expr {
	switch e1 := e.(type) {
	case fgg.Variable:
		u := e1.Typing(ds, delta, gamma, false)
		var res Expr
		res = NewVariable(e1.GetName())
		//if isInterfaceTName(ds, u) {
		if isFggITypeLit(ds, toFgTypeFromBounds(delta, u)) {
			// x.getValue().((mkRep u))
			res = NewCall(res, Name("getValue"), []Expr{})
			res = NewAssert(res, toFgTypeFromBounds(delta, u)) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case fgg.StructLit:
		u := e1.GetTName()
		t := u.GetName()
		us := u.GetTArgs()
		es_fgg := e1.GetArgs()
		es := make([]Expr, len(es_fgg)) // TODO FIXME: additional mkRep args
		fds_fgg := fgg.Fields1(ds, u)
		subs := make(map[fgg.TParam]fgg.Type)
		tfs := fgg.GetTDecl1(ds, t).GetTFormals().GetFormals()
		for i := 0; i < len(tfs); i++ {
			subs[tfs[i].GetTParam()] = us[i]
		}
		for i := 0; i < len(es_fgg); i++ {
			u_i := fds_fgg[i].GetType().TSubs(subs)
			es[i] = wrapExpr(ds, delta, gamma, es_fgg[i], u_i, wrappers)
		}
		return NewStructLit(Type(t), es)
	case fgg.Select:
		e_fgg := e1.GetExpr()
		f := e1.GetName()
		u := e1.Typing(ds, delta, gamma, false)
		var res Expr
		res = NewSelect(fgrTransExpr(ds, delta, gamma, e_fgg, wrappers), f)
		//if isInterfaceTName(ds, u) {
		if isFggITypeLit(ds, toFgTypeFromBounds(delta, u)) {
			// TODO FIXME: factor out with Variable
			res = NewCall(res, Name("getValue"), []Expr{})
			res = NewAssert(res, toFgTypeFromBounds(delta, u)) // TODO FIXME: mkRep -- "FG" for now, not FGR
		}
		return res
	case fgg.Call:
		e_recv_fgg := e1.GetRecv()
		m := e1.GetName()
		targs := e1.GetTArgs()
		args_fgg := e1.GetArgs()
		u_recv := e_recv_fgg.Typing(ds, delta, gamma, false)
		g := fgg.Methods1(ds, fgg.Bounds1(delta, u_recv))[m]
		subs := make(map[fgg.TParam]fgg.Type)
		tfs := g.GetTFormals().GetFormals()
		for i := 0; i < len(tfs); i++ {
			subs[tfs[i].GetTParam()] = targs[i]
		}
		args := make([]Expr, len(args_fgg))
		pds_fgg := g.GetParamDecls()
		for i := 0; i < len(args_fgg); i++ {
			u_i := pds_fgg[i].GetType().TSubs(subs)
			args[i] = wrapExpr(ds, delta, gamma, args_fgg[i], u_i, wrappers)
		}
		//u := e1.Typing(ds, delta, gamma, false)
		e_recv := fgrTransExpr(ds, delta, gamma, e_recv_fgg, wrappers)
		var res Expr
		res = NewCall(e_recv, m, args)
		//if isInterfaceTName(ds, erase(delta, u)) {  // !!! erase returns fg.Type (cf. wrap isStructTName)

		//fmt.Println("aaa:", e1, u, bounds(delta, u))

		//u_ret := u.TSubs(subs) // Cf. bounds(delta, u) ?

		delta1 := make(map[fgg.TParam]fgg.Type)
		for i := 0; i < len(tfs); i++ {
			tf := tfs[i]
			delta1[tf.GetTParam()] = tf.GetType()
		}
		td := fgg.GetTDecl1(ds, fgg.Bounds1(delta, u_recv).(fgg.TName).GetName())
		tfs_recv := td.GetTFormals().GetFormals()
		for i := 0; i < len(tfs_recv); i++ {
			tf := tfs_recv[i]
			delta1[tf.GetTParam()] = tf.GetType()
		}

		//u_ret := g.u.TSubs(delta1)
		u_ret := toFgTypeFromBounds(delta1, g.GetType())

		//if isInterfaceTName(ds, u_ret) {
		//if _, ok := u_ret.(TParam); ok || isInterfaceTName(ds, u_ret) {
		//if isInterfaceTName(ds, u_ret) {
		if !isFggSTypeLit(ds, u_ret) {
			// TODO FIXME: factor out with Variable
			res = NewCall(res, Name("getValue"), []Expr{})
			//res = fg.NewAssert(res, fg.Type(erase(delta, u_ret))) // TODO FIXME: mkRep -- "FG" for now, not FGR
			res = NewAssert(res, u_ret)
		}
		return res
	case fgg.Assert:
		// Need actual FGR
		panic("TODO " + reflect.TypeOf(e).String() + ": " + e.String())
	default:
		panic("Unknown Expr type " + reflect.TypeOf(e).String() + ": " + e.String())
	}
}

/* Aux */

// TODO: rename toFgrType...
// |\tau|_\Delta = t
// Basically, type name from bounds
func toFgTypeFromBounds(delta fgg.TEnv, u fgg.Type) Type {
	return Type(fgg.Bounds1(delta, u).(fgg.TName).GetName())
}

type adptrPair struct {
	sub   Type
	super Type // The "target" type, a t_I
}

// Pre: type of e <: u
// u is "target type"
func wrapExpr(ds []Decl, delta fgg.TEnv, gamma fgg.Env, e fgg.Expr, u fgg.Type,
	wrappers map[Type]adptrPair) Expr {
	t := toFgTypeFromBounds(delta, u)

	fmt.Println("aaa:", u, t, isFggSTypeLit(ds, t), isFggITypeLit(ds, t))

	if isFggSTypeLit(ds, t) {
		return fgrTransExpr(ds, delta, gamma, e, wrappers)
	} else if isFggITypeLit(ds, t) {
		u1 := e.Typing(ds, delta, gamma, false)
		e1 := fgrTransExpr(ds, delta, gamma, e, wrappers)
		return makeAdptr(delta, e1, u1, u, wrappers)
	} else {
		panic("Invalid wrap case: e=" + e.String() + ", u=" + u.String() + ", t=" + t.String())
	}
}

/* Helper */

// targ is a t_I
// TODO: rename, cf. wrap(ds, delta, gamma, e, u)
func makeAdptr(delta fgg.TEnv, e Expr, subj fgg.Type, targ fgg.Type,
	wrappers map[Type]adptrPair) StructLit {
	t1 := Type(toFgTypeFromBounds(delta, subj))
	t_I := Type(toFgTypeFromBounds(delta, targ))
	adptr := Type("Adptr_" + t1 + "_" + t_I) // TODO: factor out naming
	if _, ok := wrappers[adptr]; !ok {
		wrappers[adptr] = adptrPair{t1, t_I}
	}
	return NewStructLit(adptr, []Expr{e})
}

// ds are from FGG source (t is from toFgTypeFromBounds)
func isFggSTypeLit(ds []Decl, t Type) bool {
	for _, v := range ds {
		if _, ok := v.(fgg.STypeLit); ok && v.GetName() == string(t) {
			return true
		}
	}
	return false
}

// ds are from FGG source (t is from toFgTypeFromBounds)
func isFggITypeLit(ds []Decl, t Type) bool {
	for _, v := range ds {
		if _, ok := v.(fgg.ITypeLit); ok && (v.GetName() == string(t)) {
			return true
		}
	}
	return false
}
