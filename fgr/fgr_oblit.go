package fgr

import (
	"fmt"
	"reflect"
	//"strings"

	"github.com/rhu1/fgg/fgg"
)

var _ = fmt.Errorf
var _ = reflect.Append

/* [WIP] Obliteration: obliterate param/return types, erase field types, rep-ify type-args */

// See fgr.go for RepType (FggType) const
const GET_REP = "getRep"
const HAS_REP = "HasRep"

func Obliterate(p_fgg fgg.FGGProgram) FGRProgram { // CHECKME can also subsume existing FGG-FG trans?
	ds_fgg := p_fgg.GetDecls()

	e_fgg := p_fgg.GetMain().(fgg.FGGExpr)
	var delta fgg.Delta
	var gamma fgg.Gamma
	e_fgr := oblitExpr(ds_fgg, delta, gamma, e_fgg)

	// Translate Decls
	ds_fgr := make([]Decl, 1)                                    // There will also be an additional getRep MDecl for each t_S
	ss_HasRep := []Spec{NewSig(GET_REP, []ParamDecl{}, RepType)} // !!! Rep type name -- TODO: factor out constants
	ds_fgr[0] = NewITypeLit(Type(HAS_REP), ss_HasRep)            // TODO: factor out constant
	for i := 0; i < len(ds_fgg); i++ {
		d_fgg := ds_fgg[i]
		switch d := d_fgg.(type) {
		case fgg.STypeLit:
			recv_getRep := NewParamDecl("x0", Type(d.GetName())) // TODO: factor out constant
			t_S := d.GetName()
			tfs := d.GetPsi().GetTFormals()
			es := make([]FGRExpr, len(tfs))
			for i := 0; i < len(es); i++ {
				es[i] = NewSelect(NewVariable("x0"), tfs[i].GetTParam().String())
			}
			e_getRep := TRep{Name(t_S), es} // TODO: New constructor
			getRep := NewMDecl(recv_getRep, GET_REP /*[]RepDecl{},*/, []ParamDecl{},
				RepType, e_getRep) // TODO: factor out constants
			ds_fgr = append(ds_fgr, oblitSTypeLit(d), getRep)
		case fgg.ITypeLit:
			ds_fgr = append(ds_fgr, oblitITypeLit(d))
		case fgg.MDecl:
			ds_fgr = append(ds_fgr, oblitMDecl(ds_fgg, d))
		default:
			panic("Unexpected Decl type " + reflect.TypeOf(d).String() + ": " +
				d.String())
		}
	}

	return NewFGRProgram(ds_fgr, e_fgr)
}

/* Obliterate STypeLit, ITypeLit, Sig */

func oblitSTypeLit(s fgg.STypeLit) STypeLit {
	t := Type(s.GetName())
	psi := s.GetPsi()
	tfs := psi.GetTFormals()
	fds_fgg := s.GetFieldDecls()
	fds_fgr := make([]FieldDecl, len(tfs)+len(fds_fgg))
	for i := 0; i < len(tfs); i++ {
		fds_fgr[i] = NewFieldDecl(tfs[i].GetTParam().String(), RepType)
	}
	delta := psi.ToDelta()
	for i := 0; i < len(fds_fgg); i++ {
		fd_fgg := fds_fgg[i]
		erased := toFgrTypeFromBounds(delta, fd_fgg.GetType())
		fds_fgr[len(tfs)+i] = NewFieldDecl(fd_fgg.GetName(), erased)
	}
	return NewSTypeLit(t /*rds,*/, fds_fgr)
}

func oblitITypeLit(c fgg.ITypeLit) ITypeLit {
	t := Type(c.GetName())
	ss_fgg := c.GetSpecs()
	ss_fgr := make([]Spec, 1+len(ss_fgg))
	ss_fgr[0] = Type(HAS_REP) // TODO: add HasRep to decls -- and factor out constant
	for i := 0; i < len(ss_fgg); i++ {
		s_fgg := ss_fgg[i]
		switch s := s_fgg.(type) {
		case fgg.Type:
			panic("[TODO]: " + s.String()) // !!!
		case fgg.Sig:
			ss_fgr[i+1] = oblitSig(s)
		}
	}
	return NewITypeLit(t, ss_fgr)
}

func oblitSig(g_fgg fgg.Sig) Sig {
	m := g_fgg.GetMethod()
	tfs := g_fgg.GetPsi().GetTFormals()
	pds_fgg := g_fgg.GetParamDecls()
	pds_fgr := make([]ParamDecl, len(tfs)+len(pds_fgg))
	for i := 0; i < len(tfs); i++ {
		tf := tfs[i]
		pds_fgr[i] = NewParamDecl(tf.GetTParam().String(),
			RepType)
	}
	for i := 0; i < len(pds_fgg); i++ {
		pd_fgg := pds_fgg[i]
		pds_fgr[len(tfs)+i] = NewParamDecl(pd_fgg.GetName(), Type(HAS_REP)) // TODO: factor out constant
	}
	return NewSig(m, pds_fgr, Type(HAS_REP)) // TODO: constant
}

/* Obliterate MDecl */

func oblitMDecl(ds_fgg []Decl, d fgg.MDecl) MDecl {
	x_recv := d.GetRecvName()
	t_recv := Type(d.GetRecvTypeName())
	recv_fgr := NewParamDecl(x_recv, t_recv)
	m := d.GetName()
	tfs := d.GetMDeclPsi().GetTFormals()
	recv_tfs := d.GetRecvPsi().GetTFormals()
	pds_fgg := d.GetParamDecls()
	pds_fgr := make([]ParamDecl, len(tfs)+len(pds_fgg)) // Cf. TStructLit
	for i := 0; i < len(tfs); i++ {
		tf := tfs[i]
		pds_fgr[i] = NewParamDecl(tf.GetTParam().String(), RepType)
	}
	t_fgr := Type(HAS_REP)
	delta := d.GetRecvPsi().ToDelta()
	for i := 0; i < len(tfs); i++ {
		tf := tfs[i]
		a := tf.GetTParam()
		delta[a] = tf.GetUpperBound() // CHECKME: bounds on GetType?
	}
	subs := make(map[Variable]FGRExpr)
	v_recv := NewVariable(x_recv)
	subs[v_recv] = v_recv // CHECKME: needed o/w Variable.Subs panics -- refactor?
	for i := 0; i < len(recv_tfs); i++ {
		recv_tf := recv_tfs[i]
		a := recv_tf.GetTParam()
		subs[NewVariable(a.String())] = NewSelect(v_recv, a.String())
	}
	for i := 0; i < len(pds_fgg); i++ {
		pd := pds_fgg[i]
		x := pd.GetName()
		pds_fgr[len(tfs)+i] = NewParamDecl(x, Type(HAS_REP))
		v := NewVariable(x)
		u := pd.GetType()
		//if _, ok := u.(fgg.TParam); ok || fgg.IsInterfaceTName1(ds_fgg, u) { // !!! cf. y := y.(erase(\sigma)) -- no: allowStupid
		subs[v] = NewAssert(v, toFgrTypeFromBounds(delta, u))
	}
	gamma := make(fgg.Gamma)
	tfs_recv := d.GetRecvPsi().GetTFormals()
	us_fgg := make([]fgg.Type, len(tfs_recv))
	for i := 0; i < len(tfs_recv); i++ {
		us_fgg[i] = tfs_recv[i].GetTParam()
	}
	gamma[x_recv] = fgg.NewTName(t_recv.String(), us_fgg)
	for i := 0; i < len(pds_fgg); i++ {
		pd := pds_fgg[i]
		gamma[pd.GetName()] = pd.GetType()
	}
	e_fgr := oblitExpr(ds_fgg, delta, gamma, d.GetBody())
	e_fgr = e_fgr.Subs(subs)
	return NewMDecl(recv_fgr, m /*rds,*/, pds_fgr, t_fgr, e_fgr)
}

/* Obliterate Expr */

func oblitExpr(ds_fgg []Decl, delta fgg.Delta, gamma fgg.Gamma, e_fgg fgg.FGGExpr) FGRExpr {
	switch e := e_fgg.(type) {
	case fgg.Variable:
		return NewVariable(e.GetName())
	case fgg.StructLit:
		u := e.GetNamedType()
		t := Type(u.GetName())
		us := u.GetTArgs()
		es_fgg := e.GetElems()
		es_fgr := make([]FGRExpr, len(us)+len(es_fgg))
		for i := 0; i < len(us); i++ {
			es_fgr[i] = mkRep_oblit(us[i])
		}
		for i := 0; i < len(es_fgg); i++ {
			es_fgr[len(us)+i] = oblitExpr(ds_fgg, delta, gamma, es_fgg[i]) // !!!
		}
		return NewStructLit(t, es_fgr)
	case fgg.Select:
		e_fgg := e.GetExpr() // Shadows original e_fgg
		e_fgr := oblitExpr(ds_fgg, delta, gamma, e_fgg)
		f := e.GetField()
		u := e_fgg.Typing(ds_fgg, delta, gamma, true).(fgg.TNamed)
		fds_fgg := fgg.Fields(ds_fgg, u)
		var u_f fgg.Type = nil
		for _, fd_fgg := range fds_fgg {
			if fd_fgg.GetName() == f {
				u_f = fd_fgg.GetType()
				break
			}
		}
		if u_f == nil {
			panic("Field not found in " + u.String() + ": " + f)
		}
		var res FGRExpr
		res = NewSelect(e_fgr, f)
		du := dtype(ds_fgg, delta, gamma, e)
		//if !fgg.IsStrucrName1(ds_fgg, u_f) { // !!! don't add cast when field type is a struct
		if !fgg.IsStructType(ds_fgg, du) { // if the FGR field decl type is (erased) non-struct, in general need to cast the select result to the (erasure of the) expected FGG type
			res = NewAssert(res, toFgrTypeFromBounds(delta, u_f))
		}
		return res
	case fgg.Call:
		e_fgg := e.GetRecv() // Shadows original e_fgg
		e_fgr := oblitExpr(ds_fgg, delta, gamma, e_fgg)
		m := e.GetMethod()
		targs := e.GetTArgs()
		es_fgg := e.GetArgs()
		es_fgr := make([]FGRExpr, len(targs)+len(es_fgg))
		for i := 0; i < len(targs); i++ {
			es_fgr[i] = mkRep_oblit(targs[i])
		}
		for i := 0; i < len(es_fgg); i++ {
			es_fgr[len(targs)+i] = oblitExpr(ds_fgg, delta, gamma, es_fgg[i]) // !!!
		}

		u_recv := e_fgg.Typing(ds_fgg, delta, gamma, true)
		g := fgg.Methods(ds_fgg, fgg.Bounds(delta, u_recv))[m]
		tsubs := make(map[fgg.TParam]fgg.Type)
		tfs := g.GetPsi().GetTFormals()
		for i := 0; i < len(targs); i++ {
			tsubs[tfs[i].GetTParam()] = targs[i]
		}
		t_ret := toFgrTypeFromBounds(delta, g.GetType().TSubs(tsubs))

		var res FGRExpr
		res = NewCall(e_fgr, m, es_fgr)
		res = NewAssert(res, t_ret)
		return res
	case fgg.Assert:
		x := oblitExpr(ds_fgg, delta, gamma, e.GetExpr())
		e1 := NewCall(x, GET_REP, []FGRExpr{})
		u := e.GetType()
		e3 := NewAssert(x, toFgrTypeFromBounds(delta, u))
		p_fgg := fgg.NewProgram(ds_fgg, fgg.NewVariable(fgg.Name("dummy")), false)
		return IfThenElse{e1, mkRep_oblit(u), e3, p_fgg.String()} // TODO: New constructor
	default:
		panic("Unknown FGG Expr type: " + e_fgg.String())
	}
}

/* Aux */

// i.e., "erase" -- cf. oblit
func toFgrTypeFromBounds(delta fgg.Delta, u fgg.Type) Type {
	return Type(fgg.Bounds(delta, u).(fgg.TNamed).GetName())
}

// TODO: check where dtype should be used in wrapper translation -- and add unit tests (when return type is type param, don't want the FGG type arg, which may be struct; want the FGR target decl type as the wrapper target)
func dtype(ds []Decl, delta fgg.Delta, gamma fgg.Gamma, d fgg.FGGExpr) fgg.Type {
	switch e := d.(type) {
	case fgg.Variable:
		return gamma[e.GetName()]
	case fgg.StructLit:
		t_S := e.GetNamedType().GetName()
		td := fgg.GetTDecl(ds, t_S).(fgg.STypeLit)
		tfs := td.GetPsi().GetTFormals()
		us := make([]fgg.Type, len(tfs))
		for i := 0; i < len(us); i++ {
			us[i] = fgg.TParam(tfs[i].GetTParam().String())
		}
		return fgg.NewTName(t_S, us)
	case fgg.Select:
		u := dtype(ds, delta, gamma, e.GetExpr()).(fgg.TNamed)
		fds := fgg.Fields(ds, u)
		f := e.GetField()
		for _, fd := range fds {
			if fd.GetName() == f {
				return fd.GetType()
			}
		}
		panic("Field " + f + "not found in: " + u.String())
	case fgg.Call:
		u := fgg.Bounds(delta, dtype(ds, delta, gamma, e.GetRecv()))
		g := fgg.Methods(ds, u)[e.GetMethod()]
		return g.GetType()
	case fgg.Assert:
		return e.GetType()
	default:
		panic("Unknown FGG expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

// Post: TRep or TmpTParam
func mkRep_oblit(u fgg.Type) FGRExpr { // Duplicated from fgr_translation
	switch u1 := u.(type) {
	case fgg.TParam:
		return TmpTParam{u1.String()}
	case fgg.TNamed:
		us := u1.GetTArgs()
		es := make([]FGRExpr, len(us))
		for i := 0; i < len(us); i++ {
			es[i] = mkRep(us[i])
		}
		return TRep{u1.GetName(), es}
	default:
		panic("Unknown fgg.Type kind " + reflect.TypeOf(u).String() +
			": " + u.String())
	}
}
