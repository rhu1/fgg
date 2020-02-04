package fgg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rhu1/fgg/fg"
)

var _ = fmt.Errorf

/**
 * [WIP] Naive monomorph -- `isMonomorphisable` check not implemented yet
 */

// TODO: isMonomorphisable
/*func isMonomorphisable(p FGGProgram) bool {
	panic("[TODO]")
}*/

/* Monomoprh: FGGProgram -> FGProgram */

func Monomorph(p FGGProgram) fg.FGProgram {
	ds_fgg := p.GetDecls()
	omega := GetOmega(ds_fgg, p.GetMain().(FGGExpr)) // TODO: do "supertype closure" over omega (cf. collectSuperMethInstans)

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

// Pre: `wv` (an Omega map value) represents an instantiation of the `td` type
// TODO: decompose
func monomTDecl(ds []Decl, omega Omega, td TDecl,
	wv GroundTypeAndSigs) fg.TDecl {

	subs := make(map[TParam]Type) // Type is a TNamed
	psi := td.GetPsi()
	for i := 0; i < len(psi.tFormals); i++ {
		subs[psi.tFormals[i].name] = wv.u_ground.u_args[i]
	}
	switch d := td.(type) {
	case STypeLit:
		fds := make([]fg.FieldDecl, len(d.fDecls))
		for i := 0; i < len(d.fDecls); i++ {
			fd := d.fDecls[i]
			u_f := fd.u.TSubs(subs).(TNamed)      // "Inlined" substitution actions here -- cf. TDecl.TSubs
			if _, ok := omega[toWKey(u_f)]; !ok { // Cf. BuildWMap, extra loop over non-param TDecls, for those not seen o/w
				panic("Unknown type: " + u_f.String())
			}
			t_f_monom := toMonomId(omega[toWKey(u_f)].u_ground)
			fds[i] = fg.NewFieldDecl(fd.field, t_f_monom)
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
						pd := s.pDecls[i]
						u_p := pd.u.TSubs(subs).(TNamed)
						t_p_monom := toMonomId(omega[toWKey(u_p)].u_ground)
						pds[i] = fg.NewParamDecl(pd.name, t_p_monom)
					}
					u_ret := s.u_ret.TSubs(subs).(TNamed)
					t_ret_monom := toMonomId(omega[toWKey(u_ret)].u_ground)
					ss = append(ss, fg.NewSig(s.meth, pds, t_ret_monom))
				} else {
					// Instantiate sig for all calls of this method on this type.
					// Collect add-meth-targs for all meths called on wv.u_ground.
					mInstans := make(map[string][]Type) // Key is getTypeArgsHash([]Type)
					addMethInstans(wv, s.meth, mInstans)
					// CHECKME: if targs empty, methods "discarded" -- replace meth-params by bounds?
					for _, targs := range mInstans {
						subs1 := make(map[TParam]Type)
						for k1, v1 := range subs {
							subs1[k1] = v1
						}
						for i := 0; i < len(targs); i++ {
							subs1[s.psi.tFormals[i].name] = targs[i]
						}
						pds := make([]fg.ParamDecl, len(s.pDecls))
						for i := 0; i < len(s.pDecls); i++ {
							pd := s.pDecls[i]
							u_p := pd.u.TSubs(subs1).(TNamed)
							t_p_monom := toMonomId(omega[toWKey(u_p)].u_ground)
							pds[i] = fg.NewParamDecl(pd.name, t_p_monom)
						}
						u_ret := s.u_ret.TSubs(subs1).(TNamed)
						t_ret_monom := toMonomId(omega[toWKey(u_ret)].u_ground)
						g1 := fg.NewSig(getMonomMethName(omega, s.meth, targs), pds,
							t_ret_monom)
						ss = append(ss, g1)
					}
				}
			case TNamed: // Embedded
				u_I := s.TSubs(subs).(TNamed)
				t_monom := toMonomId(omega[toWKey(u_I)].u_ground)
				ss = append(ss, t_monom)
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

// Pre: `wv` (an Omega map value) represents an instantiation of `md.t_recv`
// N.B. `md.t_recv` is a t_S
// TODO: decompose
func monomMDecl(ds []Decl, omega Omega, md MDecl,
	wv GroundTypeAndSigs) (res []fg.MDecl) {

	subs := make(map[TParam]Type) // Type is a TNamed
	for i := 0; i < len(md.psi_recv.tFormals); i++ {
		subs[md.psi_recv.tFormals[i].name] = wv.u_ground.u_args[i]
	}
	recv := fg.NewParamDecl(md.x_recv, toMonomId(wv.u_ground))
	if len(md.psi_meth.tFormals) == 0 {
		pds := make([]fg.ParamDecl, len(md.pDecls))
		for i := 0; i < len(md.pDecls); i++ {
			pd := md.pDecls[i]
			u_p := pd.u.TSubs(subs).(TNamed) // "Inlined" substitution actions here -- cf. TDecl.TSubs
			t_p_monom := toMonomId(omega[toWKey(u_p)].u_ground)
			pds[i] = fg.NewParamDecl(pd.name, t_p_monom)
		}
		t_ret_monom := toMonomId(omega[toWKey(md.u_ret.TSubs(subs).(TNamed))].u_ground)
		e_monom := monomExpr(omega, md.e_body.TSubs(subs))
		res = append(res, fg.NewMDecl(recv, md.name, pds, t_ret_monom, e_monom))
	} else {
		// Instantiate method for all calls of md.name on any supertype.
		//mInstans := collectSuperMethInstans(ds, omega, md, wv) // reflexive
		mInstans := make(map[string][]Type) // CHECKME: should be sufficient given omega?
		addMethInstans(wv, md.name, mInstans)
		for _, targs := range mInstans {
			subs1 := make(map[TParam]Type)
			for k1, v1 := range subs {
				subs1[k1] = v1
			}
			for i := 0; i < len(targs); i++ {
				subs1[md.psi_meth.tFormals[i].name] = targs[i]
			}
			pds := make([]fg.ParamDecl, len(md.pDecls))
			for i := 0; i < len(md.pDecls); i++ {
				pd := md.pDecls[i]
				u_p := pd.u.TSubs(subs1).(TNamed)
				t_p_monom := toMonomId(omega[toWKey(u_p)].u_ground)
				pds[i] = fg.NewParamDecl(pd.name, t_p_monom)
			}
			u_ret := md.u_ret.TSubs(subs1).(TNamed)
			t_ret_monom := toMonomId(omega[toWKey(u_ret)].u_ground)
			recv := fg.NewParamDecl(md.x_recv, toMonomId(wv.u_ground))
			e_monom := monomExpr(omega, md.e_body.TSubs(subs1))
			m_monom := getMonomMethName(omega, md.name, targs)
			md1 := fg.NewMDecl(recv, m_monom, pds, t_ret_monom, e_monom)
			res = append(res, md1)
		}
	}
	return res
}

/*// Collect all instantations of calls to md on any supertype of wv.u_ground.
// - return is a map, so "duplicate" add-meth-param type instans are implicitly set-ified
// ^E.g., Calling m(A()) on some struct separately via two interfaces T1 and T2 where T2 <: T1
// Pre: `wv` (an Omega map value) represents an instantiation of `md.t_recv`
// N.B. `md.t_recv` is a t_S
func collectSuperMethInstans(ds []Decl, omega Omega, md MDecl,
	wv GroundTypeAndSigs) (mInstans map[string][]Type) {

	empty := make(Delta)
	mInstans = make(map[string][]Type)
	// Given m = md.m, forall u_I s.t. m in meths(u_I) && wv.u_ground <: u_I,
	// .. collect targs from all calls of m on u_I
	for _, wv1 := range omega {
		if wv.u_ground.ImplsDelta(ds, empty, wv1.u_ground) {
			gs := methods(ds, wv1.u_ground) // Includes embedded meths for i/face wv1.u_ground
			if _, ok := gs[md.name]; ok {
				addMethInstans(wv1, md.name, mInstans)
			}
		}
	}
	return mInstans
}*/

// Add instans of `m` in `wv` (an Omega map value) to `mInstans`
// (Only Adding instances with non-empty add-meth-targs, but that should simply depend on m's decl)
func addMethInstans(wv GroundTypeAndSigs, m Name, mInstans map[string][]Type) {
	for _, v := range wv.sigs {
		m1 := v.sig.GetMethod()
		if m1 == m && len(v.targs) > 0 {
			hash := getTypeArgsHash(v.targs)
			mInstans[hash] = v.targs
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

func monomExpr(omega Omega, e FGGExpr) fg.FGExpr {
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
		t_monom := toMonomId(omega[wk].u_ground)
		return fg.NewStructLit(t_monom, es)
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
		t_monom := toMonomId(omega[wk].u_ground)
		return fg.NewAssert(monomExpr(omega, e1.e_I), t_monom)
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

/* Helpers */

func toMonomId(u TNamed) fg.Type {
	res := u.String()
	res = strings.Replace(res, ",", ",,", -1) // TODO: refactor, cf. main.go, doMonom
	res = strings.Replace(res, "(", "<", -1)
	res = strings.Replace(res, ")", ">", -1)
	res = strings.Replace(res, " ", "", -1)
	return fg.Type(res)
}

// Pre: len(targs) > 0
func getMonomMethName(omega Omega, m Name, targs []Type) Name {
	first := toMonomId(omega[toWKey(targs[0].(TNamed))].u_ground)
	res := m + "<" + first.String()
	for _, v := range targs[1:] {
		next := toMonomId(omega[toWKey(v.(TNamed))].u_ground)
		res = res + "," + next.String()
	}
	res = res + ">"
	return Name(res)
}
