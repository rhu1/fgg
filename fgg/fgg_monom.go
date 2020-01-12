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

/* Omega -- implemented as WMap, MonomTypeAndSigs, MonomSig */

type WMap map[WKey]MonomTypeAndSigs

//type MonomMap map[WKey]fg.Type

// Hack, because TNamed cannot be used as map key directly
type WKey struct {
	t_name Name   // Just the `t` of the WVal.u_closed
	hash   string // "Hash" of the WVal.u_closed -- a closed TName
}

// MonomTypeAndSigs = closed FGG type, it's monom-name, and monom's sigs with add-meth-targs on this receiver
// TODO: integrate with GroundTypeAndSigs
type MonomTypeAndSigs struct {
	u_ground TNamed               // Pre: isGround(u_ground) -- the source (ground) FGG type
	t_monom  fg.Type              // Monom'd (i.e., FG) type name
	sigs     map[string]GroundSig // *FG* sigs on t_monom receiver -- HACK: key is MonomSig.sig.String()
	// ^Also tracks the FGG add-meth-args that gave the sigs
}

/* Monomoprh: FGGProgram -> FGProgram */

func Monomorph(p FGGProgram) fg.FGProgram {
	omega := BuildWMap(p)

	var ds []Decl
	for _, v := range p.decls {
		switch d := v.(type) {
		case TDecl:
			t := d.GetName()
			for _, v1 := range omega { // CHECKME: "prunes" unused types -- OK?
				//if k1.t_name == t {
				if v1.u_ground.t_name == t {
					ds = append(ds, monomTDecl(p.decls, omega, d, v1))
				}
			}
		case MDecl:
			for _, v1 := range omega { // CHECKME: "prunes" unused types -- OK?
				//if k1.t_name == d.t_recv {
				if v1.u_ground.t_name == d.t_recv {
					//ds = append(ds, monomMDecl(omega, d, v1)...)  // Not allowed
					for _, v := range monomMDecl(p.decls, omega, d, v1) {
						ds = append(ds, v)
					}
				}
			}
		default:
			panic("Unknown Decl kind: " + reflect.TypeOf(d).String() +
				"\n\t" + d.String())
		}
	}
	e := monomExpr(omega, p.e_main)
	return fg.NewFGProgram(ds, e, p.printf)
}

/* Monom TDecl */

// Pre: `wval` represents an instantiation of the `td` type  // TODO: refactor, decompose
func monomTDecl(ds []Decl, omega GroundMap, td TDecl,
	wval GroundTypeAndSigs) fg.TDecl {
	subs := make(map[TParam]Type) // Type is a TName
	psi := td.GetPsi()
	for i := 0; i < len(psi.tFormals); i++ {
		subs[psi.tFormals[i].name] = wval.u_ground.u_args[i]
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
		return fg.NewSTypeLit(toMonomId(wval.u_ground), fds)
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
					// forall u_S s.t. u_S <: wv.u, collect m.targs for all wv.m and mono(u_S.m)
					// ^Correction: forall u, not only u_S, i.e., including interface type receivers
					// (Cf. map.fgg, Bool().Cond(Bool())(...))
					gs := methods(ds, wval.u_ground)
					empty := make(Delta)
					targs := make(map[string][]Type)
					for _, v := range omega {
						if /*IsStructType(ds, v.u.t) &&*/ v.u_ground.Impls(ds, empty, wval.u_ground) { // N.B. now adding reflexively
							// Collect meth instans from *all* subtypes, i.e., including calls on interface receivers
							for _, v1 := range gs {
								addMethInstans(v, v1.meth, targs)
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
		return fg.NewITypeLit(toMonomId(wval.u_ground), ss)
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

// Add meth instans from `wval`, filtered by `m`, to `targs`
func addMethInstans(wval GroundTypeAndSigs, m Name, targs map[string][]Type) {
	for _, v := range wval.sigs {
		//m1 := getOrigMethName(v.sig.GetMethod())
		m1 := v.sig.GetMethod()
		if m1 == m && len(v.targs) > 0 {
			hash := "" // Use WriteTypes?
			for _, v1 := range v.targs {
				hash = hash + v1.String()
			}
			targs[hash] = v.targs
		}
	}
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

// Pre: isClosed(u)
//func toWKey(u TNamed) WKey {
func toWKey(u TNamed) string {
	/*hash := ""
	if len(u.u_args) > 0 {
		hash = u.u_args[0].String()
		for _, v := range u.u_args[1:] {
			hash = hash + ",," + v.String()
		}
	}
	return WKey{u.t_name, hash}*/
	return u.String()
}

func toMonomId(u TNamed) fg.Type {
	res := u.String()
	res = strings.Replace(res, ",", ",,", -1)
	res = strings.Replace(res, "(", "<", -1)
	res = strings.Replace(res, ")", ">", -1)
	res = strings.Replace(res, " ", "", -1)
	return fg.Type(res)
}

func isGround(u TNamed) bool {
	for _, v := range u.u_args {
		if u1, ok := v.(TNamed); !ok {
			return false
		} else if !isGround(u1) {
			return false
		}
	}
	return true
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

// Hack
func getOrigMethName(m Name) Name {
	return m[:strings.Index(m, "<")]
}

/* Convert (FGG) GroundTypeAndSigs to MonomTypeAndSigs -- temp workaround */

// TODO: refactor below
// Temp. workaround code to interface older "apply-omega" (this file) and
// .. newer "build-omega" (fgg_omega).
//func BuildWMap(p FGGProgram) WMap {
func BuildWMap(p FGGProgram) GroundMap {
	ds := p.GetDecls()
	//omega := make(WMap)

	ground := GetOmega(ds, p.GetMain().(FGGExpr))
	return ground
	/*for _, v := range ground {
		wk := toWKey(v.u_ground)
		//gs := make(map[string]MonomSig)
		gs := make(map[string]GroundSig)
		omega[wk] = MonomTypeAndSigs{v.u_ground, toMonomId(v.u_ground), gs}

		for _, pair := range v.sigs {
			if len(pair.targs) == 0 {
				continue
			}
			hash := pair.sig.String()
			pds := pair.sig.GetParamDecls()
			pds_fg := make([]fg.ParamDecl, len(pds))
			for i := 0; i < len(pds); i++ {
				pd := pds[i]
				pds_fg[i] = fg.NewParamDecl(pd.name, toMonomId(pd.u.(TNamed)))
			}
			//ret := pair.sig.u_ret.(TNamed)
			//m := getMonomMethName(omega, pair.sig.meth, pair.targs)
			//gs[hash] = MonomSig{fg.NewSig(m, pds_fg, toMonomId(ret)), pair.targs, ret}
			gs[hash] = pair
		}
	}

	return omega*/
}
