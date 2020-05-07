package fgg

import (
	"fmt"
	"reflect"
	"strconv"
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

/* */

// All m (MethInstan.meth) belong to the same t (MethInstan.u_recv.t_name)
type Mu map[string]MethInstan // Cf. Omega1, toKey_Wm

var empty_S = fg.Type("Empty")

/* Monomorph: FGGProgram -> FGProgram */

func Monomorph(p FGGProgram) fg.FGProgram {
	ds_fgg := p.GetDecls()
	omega := GetOmega1(ds_fgg, p.GetMain().(FGGExpr))
	return ApplyOmega1(p, omega)
}

func ApplyOmega1(p FGGProgram, omega Omega1) fg.FGProgram {
	var ds_monom []Decl
	for _, v := range p.decls {
		switch d := v.(type) {
		case TypeDecl:
			tds_monom := monomTDecl1(p.decls, omega, d)
			for _, v := range tds_monom {
				ds_monom = append(ds_monom, v)
			}
		case MethDecl:
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
	ds_monom = append(ds_monom, fg.NewSTypeLit(empty_S, []fg.FieldDecl{}))
	return fg.NewFGProgram(ds_monom, e_monom, p.printf)
}

func monomTDecl1(ds []Decl, omega Omega1, td TypeDecl) []fg.TDecl {
	var res []fg.TDecl
	for _, u := range omega.us {
		t := td.GetName()
		if u.t_name == t {
			eta := MakeEta(td.GetBigPsi(), u.u_args)
			mu := make(Mu)
			for k, m := range omega.ms {
				if m.u_recv.t_name == t &&
					SmallPsi(m.u_recv.GetTArgs()).Equals(SmallPsi(u.u_args)) { // TODO: fix conversions
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
	pds_empty := []fg.ParamDecl{}
	subs := make(Delta) // TODO: refactor -- because of Sig.TSubs
	for k, v := range eta {
		subs[k] = v
	}
	for _, v := range c.specs {
		switch s := v.(type) {
		case Sig: // !!! M contains Psi
			for _, m := range mu {
				if m.meth != s.meth {
					continue
				}
				theta := MakeEta(s.Psi, m.psi)
				for k, v := range eta {
					theta[k] = v
				}
				g_monom := monomSig1(s, m, theta) // !!! small psi
				ss = append(ss, g_monom)
			}
			hash := fg.NewSig(toHashSig(s.TSubs(subs)), pds_empty, empty_S)
			ss = append(ss, hash)
		case TNamed: // Embedded
			u_I := s.SubsEta(eta)
			emb_monom := toMonomId(u_I)
			ss = append(ss, emb_monom)
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

func monomMDecl1(omega Omega1, md MethDecl) []fg.MethDecl {
	var res []fg.MethDecl
	for _, m := range omega.ms {
		if !(m.u_recv.t_name == md.t_recv && m.meth == md.name) {
			continue
		}
		theta := MakeEta(md.Psi_recv, m.u_recv.u_args)
		for i := 0; i < len(md.Psi_meth.tFormals); i++ {
			theta[md.Psi_meth.tFormals[i].name] = m.psi[i].(TNamed)
		}
		recv_monom := fg.NewParamDecl(md.x_recv, toMonomId(m.u_recv))                  // !!! t_S(phi) already ground receiver
		g_monom := monomSig1(Sig{md.name, md.Psi_meth, md.pDecls, md.u_ret}, m, theta) // !!! small psi
		e_monom := monomExpr1(md.e_body, theta)
		md_monom := fg.NewMDecl(recv_monom, g_monom.GetMethod(), g_monom.GetParamDecls(), g_monom.GetReturn(), e_monom)
		res = append(res, md_monom)
	}
	pds_empty := []fg.ParamDecl{}
	empty := fg.Type("Empty")
	e_empty := fg.NewStructLit(empty, []fg.FGExpr{})
	for _, u := range omega.us {
		recv_monom := fg.NewParamDecl(md.x_recv, toMonomId(u)) // !!! t_S(phi) already ground receiver
		if u.t_name != md.t_recv {
			continue
		}
		eta := MakeEta(md.Psi_recv, u.u_args)
		subs := make(Delta) // TODO: refactor -- because of Sig.TSubs
		for k, v := range eta {
			subs[k] = v
		}
		g := md.ToSig().TSubs(subs)
		hash := fg.NewMDecl(recv_monom, toHashSig(g), pds_empty, empty, e_empty)
		res = append(res, hash)
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
		u_cast := e.u_cast.SubsEta(eta) // "Inlined" substitution actions here -- cf. M-Type
		t_monom := toMonomId(u_cast)
		return fg.NewAssert(e_monom, t_monom)
	case String: // CHECKME
		return fg.NewString(e.val)
	case Sprintf:
		args := make([]fg.FGExpr, len(e.args))
		for i := 0; i < len(e.args); i++ {
			args[i] = monomExpr1(e.args[i], eta)
		}
		return fg.NewSprintf(e.format, args)
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

/*// Pre: len(targs) > 0
func getMonomMethName(omega Omega, m Name, targs []Type) Name {
	first := toMonomId(omega[toWKey(targs[0].(TNamed))].u_ground)
	res := m + "<" + first.String()
	for _, v := range targs[1:] {
		next := toMonomId(omega[toWKey(v.(TNamed))].u_ground)
		res = res + "," + next.String()
	}
	res = res + ">"
	return Name(res)
}*/

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

func toHashSig(g Sig) string {
	subs := make(Delta)
	for i := 0; i < len(g.Psi.tFormals); i++ {
		subs[g.Psi.tFormals[i].name] = TParam("α" + strconv.Itoa(i+1))
	}
	g1 := g.TSubs(subs)
	var b strings.Builder
	b.WriteString(g.meth)
	b.WriteString("_")
	for _, v := range g1.Psi.tFormals {
		b.WriteString("_")
		b.WriteString(v.name.String())
		b.WriteString("_")
		b.WriteString(v.u_I.String())
	}
	b.WriteString("_")
	for _, v := range g1.pDecls { // Discard formal param names
		b.WriteString("_")
		b.WriteString(v.u.String())
	}
	b.WriteString("_")
	b.WriteString(g1.u_ret.String())
	res := b.String()
	res = strings.Replace(res, "(", "_", -1) // TODO
	res = strings.Replace(res, ")", "_", -1)
	res = strings.Replace(res, ",", "_", -1)
	res = strings.Replace(res, " ", "", -1)
	return res
}

/*







































 */

/* Deprecated -- Simplistic isMonom check:
   no typeparam nested in a named type in typeargs of StructLit/Call exprs */

/*
func IsMonomable(p FGGProgram) (FGGExpr, bool) {
	for _, v := range p.decls {
		switch d := v.(type) {
		case STypeLit:
		case ITypeLit:
		case MethDecl:
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

func isMonomableMDecl(d MethDecl) (FGGExpr, bool) {
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
*/

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
