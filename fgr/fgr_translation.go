package fgr

import (
	"fmt"
	"reflect"
	//"strings"

	"github.com/rhu1/fgg/fgg"
)

var _ = fmt.Errorf

/*HERE
- do assert->ifthenelse trans
- add meth-param RepDecls
- FGR eval
*/

/* FGGProgram */

func Translate(p fgg.FGGProgram) FGRProgram { // TODO FIXME: FGR -- TODO also can subsume existing FGG-FG trans?
	var ds_fgr []Decl

	// Add t_0 (etc.) to ds_fgr
	// TODO: factor out constants
	Any_0 := NewITypeLit(Type("Any_0"), []Spec{})
	Dummy_0 := NewSTypeLit(Type("Dummy_0"), []RepDecl{}, []FieldDecl{})
	ToAny_0 := NewSTypeLit(Type("ToAny_0"), []RepDecl{}, []FieldDecl{NewFieldDecl("any", Type("Any_0"))})
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
		adptr := NewSTypeLit(k, []RepDecl{}, fds)
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
			//e = NewAssert(e, toFgTypeFromBounds(delta, u))
			e = addGetValueCast(delta, e, u)
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
	tfs := s.GetTFormals().GetFormals()
	rds := make([]RepDecl, len(tfs))
	for i := 0; i < len(tfs); i++ {
		rds[i] = RepDecl{tfs[i].GetTParam(), Rep{tfs[i].GetType()}}
	}
	fds_fgg := s.GetFieldDecls()
	fds := make([]FieldDecl, len(fds_fgg))
	for i := 0; i < len(fds_fgg); i++ {
		fd := fds_fgg[i]
		fds[i] = NewFieldDecl(fd.GetName(), toFgTypeFromBounds(delta, fd.GetType()))
	}
	return NewSTypeLit(Type(s.GetName()), rds, fds)
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

	// Substituting TmpTParam's
	subs := make(map[Variable]Expr) // FIXME: Variable hack, actually subbing TmpTParam -- do as a separate disamb pass?
	subs[NewVariable(x_recv)] = NewVariable(x_recv)
	for _, pd := range pds_fgg {
		x := NewVariable(pd.GetName())
		subs[x] = x
	}
	tfs_recv := psi_recv.GetFormals()
	for _, tf := range tfs_recv {
		a := tf.GetTParam().String()
		subs[NewVariable(a)] = NewSelect(NewVariable(x_recv), a)
	}
	e = e.Subs(subs)

	return NewMDecl(recv, m, pds, t, e)
}

/* Expr */

// TODO: rename ds -> ds_fgg
// |e_FGG|_(\Delta; \Gamma) = e_FGR
func fgrTransExpr(ds []Decl, delta fgg.TEnv, gamma fgg.Env, e fgg.Expr,
	wrappers map[Type]adptrPair) Expr {
	switch e1 := e.(type) {
	case fgg.Variable:
		u := e1.Typing(ds, delta, gamma, false) // FIXME: should target type be `u` (possibly struct), or "base" (always i/face)?
		var res Expr
		res = NewVariable(e1.GetName())
		//if isInterfaceTName(ds, u) {
		if isFggITypeLit(ds, toFgTypeFromBounds(delta, u)) { // CHECKME FIXME: should check "base" type, not `u` instantiated FGG type?
			/*res = NewCall(res, Name("getValue"), []Expr{})
			res = NewAssert(res, toFgTypeFromBounds(delta, u))*/
			res = addGetValueCast(delta, res, u)
		}
		return res
	case fgg.StructLit:
		u := e1.GetTName()
		t := u.GetName()
		us := u.GetTArgs()
		es_fgg := e1.GetArgs()
		es := make([]Expr, (len(us) + len(es_fgg)))
		for i := 0; i < len(us); i++ {
			es[i] = mkRep(us[i])
		}
		fds_fgg := fgg.Fields1(ds, u)
		subs := make(map[fgg.TParam]fgg.Type)
		tfs := fgg.GetTDecl1(ds, t).GetTFormals().GetFormals()
		for i := 0; i < len(tfs); i++ {
			subs[tfs[i].GetTParam()] = //us[i]  // !!! Cf. ParamDecls in Call
				tfs[i].GetType()
		}
		for i := 0; i < len(es_fgg); i++ {
			u_i := fds_fgg[i].GetType().TSubs(subs)
			es[i+len(us)] = wrapExpr(ds, delta, gamma, es_fgg[i], u_i, wrappers)
		}
		return NewStructLit(Type(t), es)
	case fgg.Select:
		e_fgg := e1.GetExpr()
		f := e1.GetName()
		u := e1.Typing(ds, delta, gamma, false)
		var res Expr
		res = NewSelect(fgrTransExpr(ds, delta, gamma, e_fgg, wrappers), f)

		u_expr := fgg.Bounds1(delta, e1.GetExpr().Typing(ds, delta, gamma, false)).(fgg.TName)
		td := fgg.GetTDecl1(ds, u_expr.GetName()).(fgg.STypeLit)
		fds := td.GetFieldDecls() // Could use fields aux using a "dummy", cf. Call using methods
		var u_f fgg.Type = nil
		for _, fd := range fds {
			if fd.GetName() == f {
				u_f = fd.GetType()
			}
		}
		if u_f == nil {
			panic("Field not found in " + u_expr.String() + ": " + f)
		}
		delta1 := td.GetTFormals().ToTEnv()

		////if isInterfaceTName(ds, u) {
		//if isFggITypeLit(ds, toFgTypeFromBounds(delta, u)) {
		if isFggITypeLit(ds, toFgTypeFromBounds(delta1, u_f)) {
			// TODO FIXME: factor out with Variable
			/*res = NewCall(res, Name("getValue"), []Expr{})
			res = NewAssert(res, toFgTypeFromBounds(delta, u))*/
			res = addGetValueCast(delta, res, u)
		}
		return res
	case fgg.Call:
		e_recv_fgg := e1.GetRecv()
		m := e1.GetName()
		//targs := e1.GetTArgs()
		args_fgg := e1.GetArgs()

		u_recv := e_recv_fgg.Typing(ds, delta, gamma, false)
		/*g := fgg.Methods1(ds, fgg.Bounds1(delta, u_recv))[m]
		subs := make(map[fgg.TParam]fgg.Type)
		tfs := g.GetTFormals().GetFormals()
		for i := 0; i < len(tfs); i++ {
			subs[tfs[i].GetTParam()] = targs[i]
		}*/
		// !!! wrap target should be "raw" FGR decl, not FGG type -- don't want type arg instantiation, which may be t_S, we always want upper bound t_I(?)
		//t_recv := toFgTypeFromBounds(delta, u_recv)
		td := fgg.GetTDecl1(ds, fgg.Bounds1(delta, u_recv).(fgg.TName).GetName())
		tfs_recv := td.GetTFormals().GetFormals()
		//md := getMDecl(ds, t_recv, m)

		// TODO factor out -- cf. add-wrapper-meths part in Translate
		us := make([]fgg.Type, len(tfs_recv))
		for i := 0; i < len(tfs_recv); i++ {
			us[i] = tfs_recv[i].GetTParam()
		}
		dummy := fgg.NewTName(td.GetName(), us) // From the "base" type decl, not the instantiated type
		g := fgg.Methods1(ds, dummy)[m]

		delta1 := make(map[fgg.TParam]fgg.Type)
		for i := 0; i < len(tfs_recv); i++ {
			tf := tfs_recv[i]
			delta1[tf.GetTParam()] = tf.GetType()
		}
		//tfs := md.GetTFormals().GetFormals()
		tfs := g.GetTFormals().GetFormals()
		for i := 0; i < len(tfs); i++ {
			tf := tfs[i]
			delta1[tf.GetTParam()] = tf.GetType()
		}

		args := make([]Expr, len(args_fgg))
		//pds_fgg := md.GetParamDecls()
		pds_fgg := g.GetParamDecls()
		for i := 0; i < len(args_fgg); i++ {
			u_i := pds_fgg[i].GetType(). //TSubs(subs)
							TSubs(delta1) // Not toFgTypeFromBounds, need FGG Type target for wrap
			args[i] = wrapExpr(ds, delta, gamma, args_fgg[i], u_i, wrappers)
		}
		//u := e1.Typing(ds, delta, gamma, false)
		e_recv := fgrTransExpr(ds, delta, gamma, e_recv_fgg, wrappers)

		var res Expr
		res = NewCall(e_recv, m, args)

		u := g.GetType()
		////u_ret := g.u.TSubs(delta1)
		//u_ret := toFgTypeFromBounds(delta1, md.GetReturn())
		u_ret := toFgTypeFromBounds(delta1, u) // CHECKME: same as "direct" md.GetReturn().TSubs(delta1) ?
		//if isInterfaceTName(ds, u_ret) {
		//if _, ok := u_ret.(TParam); ok || isInterfaceTName(ds, u_ret) {
		//if isInterfaceTName(ds, u_ret) {
		if !isFggSTypeLit(ds, u_ret) {
			// TODO FIXME: factor out with Variable
			/*res = NewCall(res, Name("getValue"), []Expr{})
			//res = fg.NewAssert(res, fg.Type(erase(delta, u_ret))) // TODO FIXME: mkRep -- "FG" for now, not FGR
			res = NewAssert(res, u_ret)*/
			res = addGetValueCast(delta, res, u)
		}
		return res
	case fgg.Assert:
		u := e1.GetType()
		e2 := fgrTransExpr(ds, delta, gamma, e1.GetExpr(), wrappers)
		return IfThenElse{NewCall(e2, "getTypeRep", []Expr{}), mkRep(u), e2}
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

// Post: TypeTree or TmpTParam
func mkRep(u fgg.Type) Expr {
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
		panic("Unknown fgg.Type kind " + reflect.TypeOf(u).String() + ": " + u.String())
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

// ds are from FGG source (t is from toFgTypeFromBounds)
func getMDecl(ds []Decl, t Type, m Name) fgg.MDecl {
	for _, v := range ds {
		md, ok := v.(fgg.MDecl)
		if ok {
			fmt.Println("bbb:", v.String(), (md.GetRecvTypeName() == string(t)), (md.GetName() == m))
		}
		if ok && md.GetRecvTypeName() == string(t) && md.GetName() == m {
			return md
		}
	}
	panic("Method not found for type " + string(t) + ": " + m)
}

func addGetValueCast(delta fgg.TEnv, e Expr, u fgg.Type) Expr {
	e3 := NewCall(e, Name("getValue"), []Expr{})
	//e = NewAssert(e, toFgTypeFromBounds(delta, u)) // TODO FIXME: mkRep -- "FG" for now, not FGR
	e2 := mkRep(u)
	e1 := NewCall(e3, "getTypeRep", []Expr{})
	e = IfThenElse{e1, e2, e3}
	return e
}
