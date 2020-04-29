package fgg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rhu1/fgg/fg"
)

var _ = fmt.Errorf

/**
 * Monomorph
 */

/* Export */

func ToMonomId(u TNamed) fg.Type {
	return toMonomId(u)
}

/* Monomorph: FGGProgram -> FGProgram */

func Monomorph(p FGGProgram) fg.FGProgram {
	ds_fgg := p.GetDecls()

	//fmt.Println("xxxx:")

	//omega := GetOmega(ds_fgg, p.GetMain().(FGGExpr)) // TODO: do "supertype closure" over omega (cf. collectSuperMethInstans)
	omega := GetOmega1(ds_fgg, p.GetMain().(FGGExpr))
	return ApplyOmega1(p, omega)
}

func ApplyOmega1(p FGGProgram, omega Omega1) fg.FGProgram {
	var ds_monom []Decl
	for _, v := range p.decls {
		switch d := v.(type) {
		case TDecl:
			tds_monom := monomTDecl1(p.decls, omega, d)
			for _, v := range tds_monom {
				ds_monom = append(ds_monom, v)
			}
		case MDecl:
			mds_monom := monomMDecl1(omega, d)
			for _, v := range mds_monom {
				ds_monom = append(ds_monom, v)
			}
		default:
			panic("Unknown Decl kind: " + reflect.TypeOf(d).String() +
				"\n\t" + d.String())
		}
	}
	e_monom := monomExpr1(p.e_main, make(Eta))
	return fg.NewFGProgram(ds_monom, e_monom, p.printf)
}

// All m (MethInstan.meth) belong to the same t (MethInstan.u_recv.t_name)
type Mu map[string]MethInstan // Cf. Omega1, toKey_Wm

func monomTDecl1(ds []Decl, omega Omega1, td TDecl) []fg.TDecl {
	var res []fg.TDecl
	for _, u := range omega.us {
		t := td.GetName()
		if u.t_name == t {
			eta := MakeEta(td.GetBigPsi(), u.u_args)
			mu := make(Mu)
			for k, m := range omega.ms {
				if m.u_recv.t_name == t { //&& m.psi.String() == SmallPsi(u.u_args).String() {
					//fmt.Println("aaa:", m.u_recv.t_name, m.psi, t, u.u_args, m.psi.Equals(u.u_args))
					//fmt.Println("bbb:", m.psi.String(), SmallPsi(u.u_args).String())
				}
				if m.u_recv.t_name == t && SmallPsi(m.u_recv.GetTArgs()).Equals(SmallPsi(u.u_args)) { // TODO: fix conversions
					mu[k] = m
				}
			}
			t_monom := toMonomId(u)
			switch cast := td.(type) {
			case STypeLit:
				res = append(res, monomSTypeLit1(t_monom, cast, eta))
			case ITypeLit:
				res = append(res, monomITypeLit1(t_monom, cast, eta, mu))
			default:
				panic("Unknown TDecl kind: " + reflect.TypeOf(td).String() +
					"\n\t" + td.String())
			}
		}
	}
	return res
}

func monomSTypeLit1(t_monom fg.Type, s STypeLit, eta Eta) fg.STypeLit {
	fds := make([]fg.FieldDecl, len(s.fDecls))
	for i := 0; i < len(s.fDecls); i++ {
		fd := s.fDecls[i]
		u_f := fd.u.SubsEta(eta) // "Inlined" substitution actions here -- cf. M-Type
		f_monom := toMonomId(u_f)
		fds[i] = fg.NewFieldDecl(fd.field, f_monom)
	}
	return fg.NewSTypeLit(t_monom, fds)
}

func monomITypeLit1(t_monom fg.Type, c ITypeLit, eta Eta, mu Mu) fg.ITypeLit {
	var ss []fg.Spec
	for _, v := range c.specs {
		switch s := v.(type) {
		case Sig: // !!! M contains Psi
			for _, m := range mu {
				if m.meth != s.meth {
					continue
				}
				theta := MakeEta(s.psi, m.psi)
				for k, v := range eta {
					theta[k] = v
				}
				g_monom := monomSig1(s, m, theta) // !!! small psi
				ss = append(ss, g_monom)
			}
		case TNamed: // Embedded
			u_I := s.SubsEta(eta)
			t_monom := toMonomId(u_I)
			ss = append(ss, t_monom)
		default:
			panic("Unknown Spec kind: " + reflect.TypeOf(v).String() +
				"\n\t" + v.String())
		}
	}
	return fg.NewITypeLit(t_monom, ss)
}

func monomSig1(g Sig, m MethInstan, eta Eta) fg.Sig {
	//getMonomMethName(omega Omega, m Name, targs []Type) Name {
	m_monom := toMonomMethName1(m.meth, m.psi, eta) // !!! small psi
	pds_monom := make([]fg.ParamDecl, len(g.pDecls))
	for i := 0; i < len(pds_monom); i++ {
		pd := g.pDecls[i]
		t_monom := toMonomId(pd.u.SubsEta(eta)) // Cf. M-Type
		pds_monom[i] = fg.NewParamDecl(pd.name, t_monom)
	}
	ret_monom := toMonomId(g.u_ret.SubsEta(eta)) // Cf. M-Type
	return fg.NewSig(m_monom, pds_monom, ret_monom)
}

func monomMDecl1(omega Omega1, md MDecl) []fg.MDecl {
	var res []fg.MDecl
	for _, m := range omega.ms {
		if !(m.u_recv.t_name == md.t_recv && m.meth == md.name) {
			continue
		}
		theta := MakeEta(md.PsiRecv, m.u_recv.u_args)
		for i := 0; i < len(md.PsiMeth.tFormals); i++ {
			theta[md.PsiMeth.tFormals[i].name] = m.psi[i].(TNamed)
		}
		recv_monom := fg.NewParamDecl(md.x_recv, toMonomId(m.u_recv))                 // !!! t_S(phi) already ground receiver
		g_monom := monomSig1(Sig{md.name, md.PsiMeth, md.pDecls, md.u_ret}, m, theta) // !!! small psi
		e_monom := monomExpr1(md.e_body, theta)
		md_monom := fg.NewMDecl(recv_monom, g_monom.GetMethod(), g_monom.GetParamDecls(), g_monom.GetReturn(), e_monom)
		res = append(res, md_monom)
	}
	return res
}

func monomExpr1(e1 FGGExpr, eta Eta) fg.FGExpr {
	switch e := e1.(type) {
	case Variable:
		return fg.NewVariable(e.name)
	case StructLit:
		es_monom := make([]fg.FGExpr, len(e.elems))
		for i := 0; i < len(e.elems); i++ {
			es_monom[i] = monomExpr1(e.elems[i], eta)
		}
		t_monom := toMonomId(e.u_S.SubsEta(eta))
		return fg.NewStructLit(t_monom, es_monom)
	case Select:
		return fg.NewSelect(monomExpr1(e.e_S, eta), e.field)
	case Call:
		e_monom := monomExpr1(e.e_recv, eta)
		var m_monom Name
		/*if len(e.t_args) == 0 {  // Cf. toMonomMethName1
			m_monom = e.meth
		} else {*/
		m_monom = toMonomMethName1(e.meth, e.t_args, eta)
		//}
		es_monom := make([]fg.FGExpr, len(e.args))
		for i := 0; i < len(e.args); i++ {
			es_monom[i] = monomExpr1(e.args[i], eta)
		}
		return fg.NewCall(e_monom, m_monom, es_monom)
	case Assert:
		e_monom := monomExpr1(e.e_I, eta)
		t_monom := toMonomId(e.u_cast.(TNamed))
		return fg.NewAssert(e_monom, t_monom)
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e1).String() + "\n\t" +
			e1.String())
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

// !!! CHECKME: psi should already be gorunded, eta unnecessary?
func toMonomMethName1(m Name, psi SmallPsi, eta Eta) Name {
	if len(psi) == 0 {
		return m + "<>"
	}
	first := toMonomId(psi[0].SubsEta(eta))
	res := m + "<" + first.String()
	for _, v := range psi[1:] {
		next := toMonomId(v.SubsEta(eta))
		res = res + "," + next.String()
	}
	res = res + ">"
	return Name(res)
}

// returns true iff u is a TParam or contains a TParam
func isOrContainsTParam(u Type) bool {
	if _, ok := u.(TParam); ok {
		return true
	}
	u1 := u.(TNamed)
	for _, v := range u1.u_args {
		if isOrContainsTParam(v) {
			return true
		}
	}
	return false
}

/* Simplistic conservative isMonom check:
   no typeparam nested in a named type in typeargs of StructLit/Call exprs */

func IsMonomable(p FGGProgram) (FGGExpr, bool) {
	for _, v := range p.decls {
		switch d := v.(type) {
		case STypeLit:
		case ITypeLit:
		case MDecl:
			if e, ok := isMonomableMDecl(d); !ok {
				return e, false
			}
		default:
			panic("Unknown Decl kind: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}
	return isMonomableExpr(p.e_main)
}

func isMonomableMDecl(d MDecl) (FGGExpr, bool) {
	return isMonomableExpr(d.e_body)
}

// Post: if bool is true, Expr is the offender; o/w disregard Expr
func isMonomableExpr(e FGGExpr) (FGGExpr, bool) {
	switch e1 := e.(type) {
	case Variable:
		return e1, true
	case StructLit:
		for _, v := range e1.u_S.u_args {
			if u1, ok := v.(TNamed); ok {
				if isOrContainsTParam(u1) {
					return e1, false
				}
			}
		}
		for _, v := range e1.elems {
			if e2, ok := isMonomableExpr(v); !ok {
				return e2, false
			}
		}
		return e1, true
	case Select:
		return isMonomableExpr(e1.e_S)
	case Call:
		for _, v := range e1.t_args {
			if u1, ok := v.(TNamed); ok {
				if isOrContainsTParam(u1) {
					return e1, false
				}
			}
		}
		if e2, ok := isMonomableExpr(e1.e_recv); !ok {
			return e2, false
		}
		for _, v := range e1.args {
			if e2, ok := isMonomableExpr(v); !ok {
				return e2, false
			}
		}
		return e1, true
	case Assert:
		if u1, ok := e1.u_cast.(TNamed); ok {
			if isOrContainsTParam(u1) {
				return e1, false
			}
		}
		return isMonomableExpr(e1.e_I)
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

/*





























 */

func ApplyOmega(p FGGProgram, omega Omega) fg.FGProgram {
	//fmt.Println("yyyy:")

	var ds_monom []Decl
	for _, v := range p.decls {
		switch d := v.(type) {
		case TDecl:
			t := d.GetName()
			for _, wv := range omega { // CHECKME: "prunes" unused types, OK?
				if wv.u_ground.t_name == t {
					//fmt.Println("yyyy:", t)
					td_monom := monomTDecl(p.decls, omega, d, wv)
					ds_monom = append(ds_monom, td_monom)
				}
			}
		case MDecl:
			for _, wv := range omega { // CHECKME: "prunes" unused types, OK?
				if wv.u_ground.t_name == d.t_recv {
					//fmt.Println("zzzz:", d.name, d.t_recv, wv.u_ground, wv.sigs)
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
	psi := td.GetBigPsi()
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
	for i := 0; i < len(md.PsiRecv.tFormals); i++ {
		subs[md.PsiRecv.tFormals[i].name] = wv.u_ground.u_args[i]
	}
	if len(md.PsiMeth.tFormals) == 0 {
		pds := make([]fg.ParamDecl, len(md.pDecls))
		for i := 0; i < len(md.pDecls); i++ {
			pd := md.pDecls[i]
			u_p := pd.u.TSubs(subs).(TNamed) // "Inlined" substitution actions here -- cf. TDecl.TSubs
			t_p_monom := toMonomId(omega[toWKey(u_p)].u_ground)
			pds[i] = fg.NewParamDecl(pd.name, t_p_monom)
		}
		u_ret := md.u_ret.TSubs(subs).(TNamed)
		t_ret_monom := toMonomId(omega[toWKey(u_ret)].u_ground)
		recv := fg.NewParamDecl(md.x_recv, toMonomId(wv.u_ground))
		e_monom := monomExpr(omega, md.e_body.TSubs(subs))
		md1 := fg.NewMDecl(recv, md.name, pds, t_ret_monom, e_monom)
		res = append(res, md1)
	} else {
		// Instantiate method for all calls of md.name (on any supertype receiver, as per omega)
		//mInstans := collectSuperMethInstans(ds, omega, md, wv) // reflexive
		mInstans := make(map[string][]Type) // CHECKME: should be sufficient given omega?
		addMethInstans(wv, md.name, mInstans)
		for _, targs := range mInstans {
			subs1 := make(map[TParam]Type)
			for k1, v1 := range subs {
				subs1[k1] = v1
			}
			for i := 0; i < len(targs); i++ {
				subs1[md.PsiMeth.tFormals[i].name] = targs[i]
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
