package fgg

import (
	"fmt"
	"reflect"
)

var _ = fmt.Errorf

/* GroundEnv */

// Basically a Gamma for only TNamed
type GroundEnv map[Name]TNamed // Pre: forall TName, isGround

/**
 * Build Omega -- (morally) a map from ground FGG types to Sigs of (potential)
 * calls on that receiver.  N.B., calls are recorded only as seen for each
 * specific receiver type -- i.e., omega does not attempt to "respect"
 * subtyping (cf. "zigzagging" in fgg_monom).
 */

// Attempt to statically collect all ground types, and method instantiations
// called on those types, that may arise during execution
// Pre: isMonomorphisable -- TODO
/*func GetOmega(ds []Decl, e_main FGGExpr) Omega {
	omega := make(Omega)
	var gamma GroundEnv
	collectGroundTypesFromExpr(ds, gamma, e_main, omega, true)
	fixOmega(ds, gamma, omega)
	return omega
}*/

// Pre: isMonomorphisable -- TODO
func GetOmega1(ds []Decl, e_main FGGExpr) Omega1 {
	omega1 := Omega1{make(map[string]TNamed), make(map[string]MethInstan)}
	collectExpr(ds, make(GroundEnv), e_main, omega1)
	fixOmega1(ds, omega1)
	//omega1.Println()
	return omega1
}

/* Omega, GroundTypeAndSigs, GroundSig, GroundEnv */

type Omega1 struct {
	us map[string]TNamed // Pre: all TNamed are isGround
	//ms map[string]GroundTypeAndSigs // Maps u_ground.String() -> GroundTypeAndSigs{u_ground, sigs}
	ms map[string]MethInstan
}

func (w Omega1) Println() {
	fmt.Println("=== Type instances:")
	for _, v := range w.us {
		fmt.Println(v)
	}
	fmt.Println("--- Method instances:")
	for _, v := range w.ms {
		fmt.Println(v.u_recv, v.meth, v.psi)
	}
	fmt.Println("===")
}

type MethInstan struct {
	u_recv TNamed // Pre: isGround
	meth   Name
	psi    SmallPsi // Pre: all isGround
}

// Pre: isGround(u_ground)
func toKey_Wt(u_ground TNamed) string {
	return u_ground.String()
}

// Pre: isGround(x.u_ground)
func toKey_Wm(x MethInstan) string {
	return x.u_recv.String() + "_" + x.meth + "_" + x.psi.String()
}

/* fixOmega */

func fixOmega1(ds []Decl, omega Omega1) {
	for auxG(ds, omega) {
	}
}

/* Expressions */

// gamma used to type Call receiver
func collectExpr(ds []Decl, gamma GroundEnv, e FGGExpr, omega Omega1) {

	switch e1 := e.(type) {
	case Variable:
		return
	case StructLit:
		for _, elem := range e1.elems {
			collectExpr(ds, gamma, elem, omega)
		}
		omega.us[toKey_Wt(e1.u_S)] = e1.u_S
	case Select:
		collectExpr(ds, gamma, e1.e_S, omega)
	case Call:
		collectExpr(ds, gamma, e1.e_recv, omega)
		for _, e_arg := range e1.args {
			collectExpr(ds, gamma, e_arg, omega)
		}
		gamma1 := make(Gamma)
		for k, v := range gamma {
			gamma1[k] = v
		}
		u_recv := e1.e_recv.Typing(ds, make(Delta), gamma1, false).(TNamed)
		omega.us[toKey_Wt(u_recv)] = u_recv
		m := MethInstan{u_recv, e1.meth, e1.GetTArgs()} // CHECKME: why add u_recv separately?
		omega.ms[toKey_Wm(m)] = m
	case Assert:
		collectExpr(ds, gamma, e1.e_I, omega)
		u := e1.u_cast.(TNamed)
		omega.us[toKey_Wt(u)] = u
	case String: // CHECKME
		k := toKey_Wt(STRING_TYPE)
		if _, ok := omega.us[k]; !ok {
			omega.us[k] = STRING_TYPE
		}
	case Sprintf:
		k := toKey_Wt(STRING_TYPE)
		if _, ok := omega.us[k]; !ok {
			omega.us[k] = STRING_TYPE
		}
		for _, arg := range e1.args {
			collectExpr(ds, gamma, arg, omega) // Discard return
		}
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

/* Aux */

// N.B. no closure over types occurring in bounds, or *interface decl* method sigs
// CHECKME: clone omega?
func auxG(ds []Decl, omega Omega1) bool {
	res := false
	res = auxF(ds, omega) || res
	res = auxI(ds, omega) || res
	res = auxM(ds, omega) || res
	res = auxS(ds, make(Delta), omega) || res
	// I/face embeddings
	res = auxE1(ds, omega) || res
	res = auxE2(ds, omega) || res
	//res = auxP(ds, omega) || res
	return res
}

func auxF(ds []Decl, omega Omega1) bool {
	res := false
	tmp := make(map[string]TNamed)
	for _, u := range omega.us {
		if !isStructType(ds, u) || u.Equals(STRING_TYPE) { // CHECKME
			continue
		}
		for _, u_f := range Fields(ds, u) {
			cast := u_f.u.(TNamed)
			tmp[toKey_Wt(cast)] = cast
		}
	}
	for k, v := range tmp {
		if _, ok := omega.us[k]; !ok {
			omega.us[k] = v
			res = true
		}
	}
	return res
}

func auxI(ds []Decl, omega Omega1) bool {
	res := false
	tmp := make(map[string]MethInstan)
	for _, m := range omega.ms {
		if !IsNamedIfaceType(ds, m.u_recv) {
			continue
		}
		for _, m1 := range omega.ms {
			if !IsNamedIfaceType(ds, m1.u_recv) {
				continue
			}
			if m1.u_recv.Impls(ds, m.u_recv) {
				mm := MethInstan{m1.u_recv, m.meth, m.psi}
				tmp[toKey_Wm(mm)] = mm
			}
		}
	}
	for k, v := range tmp {
		if _, ok := omega.ms[k]; !ok {
			omega.ms[k] = v
			res = true
		}
	}
	return res
}

func auxM(ds []Decl, omega Omega1) bool {
	res := false
	tmp := make(map[string]TNamed)
	for _, m := range omega.ms {
		gs := methods(ds, m.u_recv)
		for _, g := range gs { // Should be only g s.t. g.meth == m.meth
			if g.meth != m.meth {
				continue
			}
			eta := MakeEta(g.Psi, m.psi)
			//fmt.Println("333:", m.u_recv, ";", m.meth)
			for _, pd := range g.pDecls {
				//fmt.Println("444:", pd.name, pd.u)
				u_pd := pd.u.SubsEta(eta) // HERE: need receiver subs also? cf. map.fgg "type b Eq(b)" -- methods should be ok?
				tmp[toKey_Wt(u_pd)] = u_pd
			}
			u_ret := g.u_ret.SubsEta(eta)
			tmp[toKey_Wt(u_ret)] = u_ret
		}
	}
	for k, v := range tmp {
		if _, ok := omega.us[k]; !ok {
			omega.us[k] = v
			res = true
		}
	}
	return res
}

func auxS(ds []Decl, delta Delta, omega Omega1) bool {
	res := false
	tmp := make(map[string]MethInstan)
	for _, m := range omega.ms {
		for _, u := range omega.us {
			if !isStructType(ds, u) || !u.ImplsDelta(ds, delta, m.u_recv) {
				continue
			}
			x0, xs, e := body(ds, u, m.meth, m.psi)
			gamma := make(GroundEnv)
			gamma[x0.name] = x0.u.(TNamed)
			for _, pd := range xs {
				gamma[pd.name] = pd.u.(TNamed)
			}
			m1 := MethInstan{u, m.meth, m.psi}
			k := toKey_Wm(m1)
			//if _, ok := omega.ms[k]; !ok { // No: initial collectExpr already adds to omega.ms
			tmp[k] = m1
			collectExpr(ds, gamma, e, omega)
			//}
		}
	}
	for k, v := range tmp {
		if _, ok := omega.ms[k]; !ok {
			omega.ms[k] = v
			res = true
		}
	}
	return res
}

// Add embedded types
func auxE1(ds []Decl, omega Omega1) bool {
	res := false
	tmp := make(map[string]TNamed)
	for _, u := range omega.us {
		if !isNamedIfaceType(ds, u) {
			continue
		}
		td_I := getTDecl(ds, u.t_name).(ITypeLit)
		eta := MakeEta(td_I.Psi, u.u_args)
		for _, s := range td_I.specs {
			if u_emb, ok := s.(TNamed); ok {
				u_sub := u_emb.SubsEta(eta)
				tmp[toKey_Wt(u_sub)] = u_sub
			}
		}
	}
	for k, v := range tmp {
		if _, ok := omega.us[k]; !ok {
			omega.us[k] = v
			res = true
		}
	}
	return res
}

// Propagate method instances up to embedded supertypes
func auxE2(ds []Decl, omega Omega1) bool {
	res := false
	tmp := make(map[string]MethInstan)
	for _, m := range omega.ms {
		if !isNamedIfaceType(ds, m.u_recv) {
			continue
		}
		td_I := getTDecl(ds, m.u_recv.t_name).(ITypeLit)
		eta := MakeEta(td_I.Psi, m.u_recv.u_args)
		for _, s := range td_I.specs {
			if u_emb, ok := s.(TNamed); ok {
				u_sub := u_emb.SubsEta(eta)
				gs := methods(ds, u_sub)
				for _, g := range gs {
					if m.meth == g.meth {
						m_emb := MethInstan{u_sub, m.meth, m.psi}
						tmp[toKey_Wm(m_emb)] = m_emb
					}
				}
			}
		}
	}
	for k, v := range tmp {
		if _, ok := omega.ms[k]; !ok {
			omega.ms[k] = v
			res = true
		}
	}
	return res
}

/*func auxP(ds []Decl, omega Omega1) bool {
	res := false
	tmp := make(map[string]MethInstan)
	for _, u := range omega.us {
		if !isNamedIfaceType(ds, u) {
			continue
		}
		gs := methods(ds, u)
		for _, g := range gs {
			psi := make(SmallPsi, len(g.psi.tFormals))
			for i := 0; i < len(psi); i++ {
				psi[i] = g.psi.tFormals[i].u_I
			}
			m := MethInstan{u, g.meth, psi}
			tmp[toKey_Wm(m)] = m
			fmt.Println("222:", u, ";", m)
		}
	}
	for k, v := range tmp {
		if _, ok := omega.ms[k]; !ok {
			omega.ms[k] = v
			res = true
		}
	}
	return res
}*/

/* Helpers */

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
