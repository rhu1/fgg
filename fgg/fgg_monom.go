package fgg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rhu1/fgg/fg"
)

var _ = fmt.Errorf

/**
 * [WIP] Naive monomorph -- `isMonomophisable` check not implemented yet
 */

// TODO: isMonomorphisable
// func isMonomorphisable(p FGGProgram) bool { ... }

/* Monomoprh: FGGProgram -> FGProgram */

func Monomorph(p FGGProgram) fg.FGProgram {
	ds_fgg := p.GetDecls()
	omega := GetOmega(ds_fgg, p.GetMain().(FGGExpr))

	var ds_monom []Decl
	for _, v := range p.decls {
		switch d := v.(type) {
		case TDecl:
			t := d.GetName()
			for _, wv := range omega { // CHECKME: "prunes" unused types, OK?
				if wv.u_ground.t_name == t {
					td_monom := monomTDecl(p.decls, omega, d, wv)
					ds_monom = append(ds_monom, td_monom)
				}
			}
		case MDecl:
			for _, wv := range omega { // CHECKME: "prunes" unused types, OK?
				if wv.u_ground.t_name == d.t_recv {
					mds_monom := monomMDecl(p.decls, omega, d, wv)
					for _, v := range mds_monom {
						ds_monom = append(ds_monom, v)
					}
				}
			}
		default:
			panic("Unknown Decl kind: " + reflect.TypeOf(d).String() +
				"\n\t" + d.String())
		}
	}

	e_monom := monomExpr(omega, p.e_main)
	return fg.NewFGProgram(ds_monom, e_monom, p.printf)
}

/* Monom TDecl */

// Pre: `wv` (an "omega" map value) represents an instantiation of the `td` type
// TODO: refactor, decompose
func monomTDecl(ds []Decl, omega GroundMap, td TDecl,
	wv GroundTypeAndSigs) fg.TDecl {
	subs := make(map[TParam]Type) // Type is a TName
	psi := td.GetPsi()
	for i := 0; i < len(psi.tFormals); i++ {
		subs[psi.tFormals[i].name] = wv.u_ground.u_args[i]
	}
	switch d := td.(type) {
	case STypeLit:
		fds := make([]fg.FieldDecl, len(d.fDecls))
		for i := 0; i < len(d.fDecls); i++ {
			tmp := d.fDecls[i]
			u := tmp.u.TSubs(subs).(TNamed)     // "Inlined" substitution actions here -- cf. TDecl.TSubs
			if _, ok := omega[toWKey(u)]; !ok { // Cf. BuildWMap, extra loop over non-param TDecls, for those non seen o/w
				panic("Unknown type: " + u.String())
			}
			fds[i] = fg.NewFieldDecl(tmp.field, toMonomId(omega[toWKey(u)].u_ground))
		}
		return fg.NewSTypeLit(toMonomId(wv.u_ground), fds)
	case ITypeLit:
		var ss []fg.Spec
		for _, v := range d.specs {
			switch s := v.(type) {
			case Sig:
				if len(s.psi.tFormals) == 0 {
					pds := make([]fg.ParamDecl, len(s.pDecls))
					for i := 0; i < len(s.pDecls); i++ {
						tmp := s.pDecls[i]
						u_p := tmp.u.TSubs(subs).(TNamed)
						pds[i] = fg.NewParamDecl(tmp.name, toMonomId(omega[toWKey(u_p)].u_ground))
					}
					u := s.u_ret.TSubs(subs).(TNamed)
					ss = append(ss, fg.NewSig(s.meth, pds, toMonomId(omega[toWKey(u)].u_ground)))
				} else {
					// forall u s.t. u <: wv.u, collect add-meth-targs for all meths called on u
					gs := methods(ds, wv.u_ground)
					delta_empty := make(Delta)
					targs := make(map[string][]Type) // Key is getTypeArgsHash([]Type)
					for _, wv1 := range omega {
						if wv1.u_ground.Impls(ds, delta_empty, wv.u_ground) {
							// Collect meth instans from *all* subtypes
							// (including calls on i/face receivers -- cf. map.fgg, Bool().Cond(Bool())(...))
							// Includes reflexive
							for _, v1 := range gs {
								addMethInstans(wv1, v1.meth, targs)
							}
						}
					}
					// CHECKME: if targs empty, methods "discarded" -- replace meth-params by bounds?
					for _, v := range targs { // CHECKME: factor out with MDecl?
						subs1 := make(map[TParam]Type)
						for k1, v1 := range subs {
							subs1[k1] = v1
						}
						for i := 0; i < len(v); i++ {
							subs1[s.psi.tFormals[i].name] = v[i]
						}
						pds := make([]fg.ParamDecl, len(s.pDecls))
						for i := 0; i < len(s.pDecls); i++ {
							tmp := s.pDecls[i]
							u_p := tmp.u.TSubs(subs1).(TNamed)
							pds[i] = fg.NewParamDecl(tmp.name, toMonomId(omega[toWKey(u_p)].u_ground))
						}
						u := s.u_ret.TSubs(subs1).(TNamed)
						g1 := fg.NewSig(getMonomMethName(omega, s.meth, v), pds,
							toMonomId(omega[toWKey(u)].u_ground))
						ss = append(ss, g1)
					}
				}
			case TNamed:
				ss = append(ss, toMonomId(omega[toWKey(s)].u_ground))
			default:
				panic("Unknown Spec kind: " + reflect.TypeOf(v).String() +
					"\n\t" + v.String())
			}
		}
		return fg.NewITypeLit(toMonomId(wv.u_ground), ss)
	default:
		panic("Unknown TDecl kind: " + reflect.TypeOf(d).String() +
			"\n\t" + d.String())
	}
}

/* Monom MDecl */

// Pre: `wval` represents an instantiation of `md.t_recv`  // TODO: decompose
func monomMDecl(ds []Decl, omega GroundMap, md MDecl,
	wval GroundTypeAndSigs) (res []fg.MDecl) {
	subs := make(map[TParam]Type) // Type is a TName
	for i := 0; i < len(md.psi_recv.tFormals); i++ {
		subs[md.psi_recv.tFormals[i].name] = wval.u_ground.u_args[i]
	}
	recv := fg.NewParamDecl(md.x_recv, toMonomId(wval.u_ground))
	if len(md.psi_meth.tFormals) == 0 {
		pds := make([]fg.ParamDecl, len(md.pDecls))
		for i := 0; i < len(md.pDecls); i++ {
			tmp := md.pDecls[i]
			u := tmp.u.TSubs(subs).(TNamed) // "Inlined" substitution actions here -- cf. TDecl.TSubs
			pds[i] = fg.NewParamDecl(tmp.name, toMonomId(omega[toWKey(u)].u_ground))
		}
		t := toMonomId(omega[toWKey(md.u_ret.TSubs(subs).(TNamed))].u_ground)
		e := monomExpr(omega, md.e_body.TSubs(subs))
		res = append(res, fg.NewMDecl(recv, md.name, pds, t, e))
	} else {
		targs := collectZigZagMethInstans(ds, omega, md, wval) // CHECKME: maybe not needed? (w.r.t. revised fgg_omega)
		if len(targs) == 0 {
			// ^Means no u_I, if len(wv.gs)>0 -- targs doesn't (yet) include wv.gs
			addMethInstans(wval, md.name, targs)
		}
		for _, v := range targs { // CHECKME: factor out with ITypeLit?
			subs1 := make(map[TParam]Type)
			for k1, v1 := range subs {
				subs1[k1] = v1
			}
			for i := 0; i < len(v); i++ {
				subs1[md.psi_meth.tFormals[i].name] = v[i]
			}
			recv := fg.NewParamDecl(md.x_recv, toMonomId(wval.u_ground))
			pds := make([]fg.ParamDecl, len(md.pDecls))
			for i := 0; i < len(md.pDecls); i++ {
				tmp := md.pDecls[i]
				u_p := tmp.u.TSubs(subs1).(TNamed)
				pds[i] = fg.NewParamDecl(tmp.name, toMonomId(omega[toWKey(u_p)].u_ground))
			}
			u := md.u_ret.TSubs(subs1).(TNamed)
			e := monomExpr(omega, md.e_body.TSubs(subs1))
			md1 := fg.NewMDecl(recv, getMonomMethName(omega, md.name, v), pds,
				toMonomId(omega[toWKey(u)].u_ground), e)
			res = append(res, md1)
		}
	}
	return res
}

// CHECKME: is this still needed now?  (given revised fgg_omega?)
//
// N.B. return is empty, i.e., does not include wv.gs, if no u_I
// N.B. return is a map, so "duplicate" add-meth-param type instans are implicitly setify-ed
// ^E.g., Calling m(A()) on some struct separately via two interfaces T1 and T2 where T2 <: T1
func collectZigZagMethInstans(ds []Decl, omega GroundMap, md MDecl,
	wval GroundTypeAndSigs) map[string][]Type {
	empty := make(Delta)
	targs := make(map[string][]Type)
	// Given m = md.m, forall u_I s.t. m in meths(u_I) && wv.u <: u_I, ..
	// ..forall u_S s.t. u_S <: u_I, collect targs for all mono(u_S.m)
	// ^Correction: forall u, not only u_S
	for _, v := range omega {
		if IsNamedIfaceType(ds, v.u_ground) && wval.u_ground.Impls(ds, empty, v.u_ground) {
			gs := methods(ds, v.u_ground)
			if _, ok := gs[md.name]; ok {
				addMethInstans(v, md.name, targs)
				for _, v1 := range omega {
					if /*isStructTName(ds, v1.u) &&*/ v1.u_ground.Impls(ds, empty, v.u_ground) {
						addMethInstans(v1, md.name, targs)
					}
				}
			}
		}
	}
	return targs
}

// Add instans of `m` in `wv` (an "omega" map value) to `targs`
// (Adding instances with non-empty add-meth-targs, but that should simply depend on m's decl)
func addMethInstans(wv GroundTypeAndSigs, m Name, targs map[string][]Type) {
	for _, v := range wv.sigs {
		m1 := v.sig.GetMethod()
		if m1 == m && len(v.targs) > 0 {
			hash := getTypeArgsHash(v.targs)
			targs[hash] = v.targs
		}
	}
}

func getTypeArgsHash(us []Type) string {
	hash := "" // Use WriteTypes?
	for _, v1 := range us {
		hash = hash + v1.String()
	}
	return hash
}

/* Monom FGGExprs */

func monomExpr(omega GroundMap, e FGGExpr) fg.FGExpr {
	switch e1 := e.(type) {
	case Variable:
		return fg.NewVariable(e1.name)
	case StructLit:
		es := make([]fg.FGExpr, len(e1.elems))
		for i := 0; i < len(e1.elems); i++ {
			es[i] = monomExpr(omega, e1.elems[i])
		}
		wk := toWKey(e1.u_S)
		if _, ok := omega[wk]; !ok {
			panic("Unknown type: " + e1.u_S.String())
		}
		return fg.NewStructLit(toMonomId(omega[wk].u_ground), es)
	case Select:
		return fg.NewSelect(monomExpr(omega, e1.e_S), e1.field)
	case Call:
		e2 := monomExpr(omega, e1.e_recv)
		var m Name
		if len(e1.t_args) == 0 {
			m = e1.meth
		} else {
			m = getMonomMethName(omega, e1.meth, e1.t_args)
		}
		es := make([]fg.FGExpr, len(e1.args))
		for i := 0; i < len(e1.args); i++ {
			es[i] = monomExpr(omega, e1.args[i])
		}
		return fg.NewCall(e2, m, es)
	case Assert:
		wk := toWKey(e1.u_cast.(TNamed))
		if _, ok := omega[wk]; !ok {
			panic("Unknown type: " + e1.u_cast.String())
		}
		return fg.NewAssert(monomExpr(omega, e1.e_I),
			toMonomId(omega[wk].u_ground))
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

/* Helpers */

func toMonomId(u TNamed) fg.Type {
	res := u.String()
	res = strings.Replace(res, ",", ",,", -1)
	res = strings.Replace(res, "(", "<", -1)
	res = strings.Replace(res, ")", ">", -1)
	res = strings.Replace(res, " ", "", -1)
	return fg.Type(res)
}

// Pre: len(targs) > 0
//func getMonomMethName(omega WMap, m Name, targs []Type) Name {
func getMonomMethName(omega GroundMap, m Name, targs []Type) Name {
	res := m + "<" + toMonomId(omega[toWKey(targs[0].(TNamed))].u_ground).String()
	for _, v := range targs[1:] {
		res = res + "," + toMonomId(omega[toWKey(v.(TNamed))].u_ground).String()
	}
	res = res + ">"
	return Name(res)
}
