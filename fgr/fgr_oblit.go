package fgr

import (
	"fmt"
	"reflect"
	//"strings"

	"github.com/rhu1/fgg/fgg"
)

var _ = fmt.Errorf
var _ = reflect.Append

/* FGGProgram */

func Obliterate(p_fgg fgg.FGGProgram) FGRProgram { // CHECKME can also subsume existing FGG-FG trans?
	ds_fgg := p_fgg.GetDecls()

	e_fgg := p_fgg.GetExpr().(fgg.Expr)
	var delta fgg.TEnv
	var gamma fgg.Env
	e_fgr := oblitExpr(ds_fgg, delta, gamma, e_fgg)

	// Translate Decls
	ds_fgr := make([]Decl, 1)                                         // There will also be an additional getRep MDecl for each t_S
	ss_GetRep := []Spec{NewSig("getRep", []ParamDecl{}, Type("Rep"))} // !!! Rep type name -- TODO: factor out constants
	ds_fgr[0] = NewITypeLit(Type("GetRep"), ss_GetRep)                // TODO: factor out constant
	for i := 0; i < len(ds_fgg); i++ {
		d_fgg := ds_fgg[i]
		switch d := d_fgg.(type) {
		case fgg.STypeLit:
			recv_getRep := NewParamDecl("x", Type(d.GetName())) // TODO: factor out constant
			t_S := d.GetName()
			tfs := d.GetTFormals().GetFormals()
			es := make([]Expr, len(tfs))
			for i := 0; i < len(es); i++ {
				es[i] = NewSelect(NewVariable("x"), tfs[i].GetTParam().String())
			}
			e_getRep := TypeTree{Type(t_S), es} // TODO: New constructor
			getRep := NewMDecl(recv_getRep, "getRep", []RepDecl{}, []ParamDecl{},
				Type("Rep"), e_getRep) // TODO: factor out constants
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

func oblitSTypeLit(s fgg.STypeLit) STypeLit {
	t := Type(s.GetName())
	tfs := s.GetTFormals().GetFormals()
	rds := make([]RepDecl, len(tfs))
	for i := 0; i < len(rds); i++ {
		tf := tfs[i]
		rds[i] = RepDecl{tf.GetTParam(), Rep{tf.GetType()}} // TODO: make `New` constructor
	}
	fds_fgg := s.GetFieldDecls()
	fds_fgr := make([]FieldDecl, len(fds_fgg))
	for i := 0; i < len(fds_fgg); i++ {
		fd_fgg := fds_fgg[i]
		fds_fgr[i] = NewFieldDecl(fd_fgg.GetName(), Type("GetRep")) // TODO: factor out constant
	}
	return NewSTypeLit(t, rds, fds_fgr)
}

func oblitITypeLit(c fgg.ITypeLit) ITypeLit {
	t := Type(c.GetName())
	ss_fgg := c.GetSpecs()
	ss_fgr := make([]Spec, 1+len(ss_fgg))
	ss_fgr[0] = Type("GetRep") // TODO: add GetRep to decls -- and factor out constant
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
	m := g_fgg.GetName()
	tfs := g_fgg.GetTFormals().GetFormals()
	pds_fgg := g_fgg.GetParamDecls()
	pds_fgr := make([]ParamDecl, len(tfs)+len(pds_fgg))
	for i := 0; i < len(tfs); i++ {
		tf := tfs[i]
		pds_fgr[i] = NewParamDecl(tf.GetTParam().String(),
			//Rep{tf.GetType()}) // TODO: !!! Rep `New` constructor
			Type("Rep")) // !!! TODO: factor out constant
	}
	for i := 0; i < len(pds_fgg); i++ {
		pd_fgg := pds_fgg[i]
		pds_fgr[i] = NewParamDecl(pd_fgg.GetName(), Type("GetRep")) // TODO: factor out constant
	}
	return NewSig(m, pds_fgr, Type("GetRep")) // TODO: constant
	// FIXME: need RepDecl in Sig?
}

func oblitMDecl(ds_fgg []Decl, d fgg.MDecl) MDecl {
	x_recv := d.GetRecvName()
	t_recv := Type(d.GetRecvTypeName())
	recv_fgr := NewParamDecl(x_recv, t_recv)
	m := d.GetName()
	tfs := d.GetTFormals().GetFormals()
	rds := make([]RepDecl, len(tfs))
	for i := 0; i < len(tfs); i++ {
		tf := tfs[i]
		rds[i] = RepDecl{tf.GetTParam(), Rep{tf.GetType()}} // TODO: `New` constructors
	}
	pds_fgg := d.GetParamDecls()
	pds_fgr := make([]ParamDecl, len(pds_fgg))
	for i := 0; i < len(pds_fgg); i++ {
		pd := pds_fgg[i]
		pds_fgr[i] = NewParamDecl(pd.GetName(), Type("GetRep"))
	}
	t_fgr := Type("GetRep")
	delta := d.GetRecvTFormals().ToTEnv()
	for i := 0; i < len(tfs); i++ {
		tf := tfs[i]
		delta[tf.GetTParam()] = tf.GetType() // CHECKME: bounds on GetType?
	}
	gamma := make(fgg.Env)
	tfs_recv := d.GetRecvTFormals().GetFormals()
	us_fgg := make([]fgg.Type, len(tfs_recv))
	for i := 0; i < len(tfs_recv); i++ {
		us_fgg[i] = tfs_recv[i].GetTParam()
	}
	gamma[x_recv] = fgg.NewTName(string(t_recv), us_fgg)
	for i := 0; i < len(pds_fgg); i++ {
		pd := pds_fgg[i]
		gamma[pd.GetName()] = pd.GetType()
	}
	e_fgr := oblitExpr(ds_fgg, delta, gamma, d.GetExpr())
	return NewMDecl(recv_fgr, m, rds, pds_fgr, t_fgr, e_fgr)
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
			es_fgr[len(us)+i] = oblitExpr(ds_fgg, delta, gamma, es_fgg[i]) // !!!
		}
		return NewStructLit(t, es_fgr)
	case fgg.Select:
		e_fgr := oblitExpr(ds_fgg, delta, gamma, e.GetExpr())
		u := e_fgg.Typing(ds_fgg, delta, gamma, true).(fgg.TName)
		fds_fgg := fgg.Fields1(ds_fgg, u) // !!! CHECKME: bounds on u
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
		e_fgg := e.GetRecv()
		e_fgr := oblitExpr(ds_fgg, delta, gamma, e_fgg)
		m := e.GetName()
		targs := e.GetTArgs()
		es_fgg := e.GetArgs()
		es_fgr := make([]Expr, len(targs)+len(es_fgg))
		for i := 0; i < len(targs); i++ {
			es_fgr[i] = oblitMkRep(targs[i])
		}
		for i := 0; i < len(es_fgg); i++ {
			es_fgr[len(targs)+i] = oblitExpr(ds_fgg, delta, gamma, es_fgg[i]) // !!!
		}

		u_recv := e_fgg.Typing(ds_fgg, delta, gamma, true)
		g := fgg.Methods1(ds_fgg, fgg.Bounds1(delta, u_recv))[m]
		tsubs := make(map[fgg.TParam]fgg.Type)
		tfs := g.GetTFormals().GetFormals()
		for i := 0; i < len(targs); i++ {
			tsubs[tfs[i].GetTParam()] = targs[i]
		}
		t_ret := toFgrTypeFromBounds(delta, g.GetType().TSubs(tsubs))

		var res Expr
		res = NewCall(e_fgr, m, es_fgr)
		res = NewAssert(res, t_ret)
		return res
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
		return makeTypeTree(u1)
	default:
		panic("Unknown fgg.Type kind " + reflect.TypeOf(u).String() + ": " +
			u.String())
	}
}

func makeTypeTree(u1 fgg.TName) TypeTree {
	us := u1.GetTArgs()
	es := make([]Expr, len(us))
	for i := 0; i < len(us); i++ {
		es[i] = mkRep(us[i])
	}
	return TypeTree{Type(u1.GetName()), es}
}