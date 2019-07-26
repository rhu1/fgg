package fgg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rhu1/fgg/fg"
)

var _ = fmt.Errorf

/* Naive monomorph -- !!WIP!! */

type ClosedEnv map[Name]TName // Pre: forall TName, isClosed

// func isMonomorphisable(p FGGProgram) bool { ... }

// TODO: reformat (e.g., "<...>") to make an actual FG program
func Monomorph(p FGGProgram) fg.FGProgram {
	var gamma ClosedEnv
	omega := make(WMap)
	MakeWMap(p.GetDecls(), gamma, p.GetExpr().(Expr), omega) // Populates omega

	var ds []Decl
	for _, v := range p.ds {
		switch d := v.(type) {
		case TDecl:
			t := d.GetName()
			for k1, v1 := range omega { // CHECKME: "prunes" unused types -- OK?
				if k1.t == t {
					ds = append(ds, monomTDecl(p.ds, omega, d, v1))
				}
			}
		case MDecl:
			for k1, v1 := range omega { // CHECKME: "prunes" unused types -- OK?
				if k1.t == d.t_recv {
					//ds = append(ds, monomMDecl(omega, d, v1)...)  // Not allowed
					for _, v := range monomMDecl(p.ds, omega, d, v1) {
						ds = append(ds, v)
					}
				}
			}
		default:
			panic("Unknown Decl kind: " + reflect.TypeOf(d).String() +
				"\n\t" + d.String())
		}
	}
	e := monomExpr(omega, p.e)
	return fg.NewFGProgram(ds, e)
}

// Pre: `wv` represents an instantiation of the `td` type  // TODO: refactor, decompose
func monomTDecl(ds []Decl, omega WMap, td TDecl, wv WVal) fg.TDecl {
	subs := make(map[TParam]Type) // Type is a TName
	psi := td.GetTFormals()
	for i := 0; i < len(psi.tfs); i++ {
		subs[psi.tfs[i].a] = wv.u.us[i]
	}
	switch d := td.(type) {
	case STypeLit:
		fds := make([]fg.FieldDecl, len(d.fds))
		for i := 0; i < len(d.fds); i++ {
			tmp := d.fds[i]
			u := tmp.u.TSubs(subs).(TName) // "Inlined" substitution actions here -- cf. TDecl.TSubs
			fds[i] = fg.NewFieldDecl(tmp.f, omega[toWKey(u)].id)
		}
		return fg.NewSTypeLit(wv.id, fds)
	case ITypeLit:
		var ss []fg.Spec
		for _, v := range d.ss {
			switch s := v.(type) {
			case Sig:
				if len(s.psi.tfs) == 0 {
					pds := make([]fg.ParamDecl, len(s.pds))
					for i := 0; i < len(s.pds); i++ {
						tmp := s.pds[i]
						u_p := tmp.u.TSubs(subs).(TName)
						pds[i] = fg.NewParamDecl(tmp.x, omega[toWKey(u_p)].id)
					}
					u := s.u.TSubs(subs).(TName)
					ss = append(ss, fg.NewSig(s.m, pds, omega[toWKey(u)].id))
				} else {
					/*for _, v := range wv.gs {  // Subsumed by below?
						if v.g.GetMethName() == s.m {
							ss = append(ss, v.g)
						}
					}*/
					// forall u_S s.t. u_S <: wv.u, collect m.targs for all wv.m and mono(u_S.m)
					gs := methods(ds, wv.u)
					empty := make(TEnv)
					targs := make(map[string][]Type)
					for _, v := range omega {
						if isStructType(ds, v.u.t) && v.u.Impls(ds, empty, wv.u) {
							for _, v1 := range v.gs {
								m1 := getOrigMethName(v1.g.GetMethName())
								if _, ok := gs[m1]; ok && len(v1.targs) > 0 { // Redundant?
									hash := "" // TODO: factor out  // Use WriteTypes?
									for _, v2 := range v1.targs {
										hash = hash + v2.String()
									}
									targs[hash] = v1.targs
								}
							}
						}
					}
					// CHECKME: if targs empty, methods "discarded" -- replace meth-params by bounds?
					for _, v := range targs {
						subs1 := make(map[TParam]Type)
						for k1, v1 := range subs {
							subs1[k1] = v1
						}
						for i := 0; i < len(v); i++ {
							subs1[s.psi.tfs[i].a] = v[i]
						}
						pds := make([]fg.ParamDecl, len(s.pds))
						for i := 0; i < len(s.pds); i++ {
							tmp := s.pds[i]
							u_p := tmp.u.TSubs(subs1).(TName)
							pds[i] = fg.NewParamDecl(tmp.x, omega[toWKey(u_p)].id)
						}
						u := s.u.TSubs(subs1).(TName)
						g1 := fg.NewSig(getMonomMethName(omega, s.m, v), pds,
							omega[toWKey(u)].id)
						ss = append(ss, g1)
					}
				}
			case TName:
				ss = append(ss, omega[toWKey(s)].id)
			default:
				panic("Unknown Spec kind: " + reflect.TypeOf(v).String() +
					"\n\t" + v.String())
			}
		}
		return fg.NewITypeLit(wv.id, ss)
	default:
		panic("Unknown TDecl kind: " + reflect.TypeOf(d).String() +
			"\n\t" + d.String())
	}
}

// Pre: `wv` represents an instantiation of `md.t_recv`  // TODO: refactor, decompose
func monomMDecl(ds []Decl, omega WMap, md MDecl, wv WVal) (res []fg.MDecl) {
	subs := make(map[TParam]Type) // Type is a TName
	for i := 0; i < len(md.psi_recv.tfs); i++ {
		subs[md.psi_recv.tfs[i].a] = wv.u.us[i]
	}
	recv := fg.NewParamDecl(md.x_recv, wv.id)
	if len(md.psi.tfs) == 0 {
		pds := make([]fg.ParamDecl, len(md.pds))
		for i := 0; i < len(md.pds); i++ {
			tmp := md.pds[i]
			u := tmp.u.TSubs(subs).(TName) // "Inlined" substitution actions here -- cf. TDecl.TSubs
			pds[i] = fg.NewParamDecl(tmp.x, omega[toWKey(u)].id)
		}
		t := omega[toWKey(md.u.TSubs(subs).(TName))].id
		e := monomExpr(omega, md.e.TSubs(subs))
		res = append(res, fg.NewMDecl(recv, md.m, pds, t, e))
	} else {
		empty := make(TEnv)
		targs := make(map[string][]Type)
		// Given m = md.m, forall u_I s.t. m in meths(u_I) && wv.u <: u_I, ..
		// ..forall u_S s.t. u_S <: u_I, collect targs for all mono(u_S.m)
		for _, v := range omega {
			if isInterfaceTName(ds, v.u) && wv.u.Impls(ds, empty, v.u) {
				gs := methods(ds, v.u)
				if _, ok := gs[md.m]; ok {
					for _, v1 := range omega {
						if isStructTName(ds, v1.u) && v1.u.Impls(ds, empty, v.u) {
							for _, v2 := range v1.gs {
								m2 := getOrigMethName(v2.g.GetMethName())
								if m2 == md.m && len(v2.targs) > 0 {
									hash := "" // TODO: factor out
									for _, v3 := range v2.targs {
										hash = hash + v3.String()
									}
									targs[hash] = v2.targs
								}
							}
						}
					}
				}
			}
		}
		if len(targs) == 0 { // Means no u_I, if len(wv.gs)>0 -- targs doesn't include wv.gs
			for _, v := range wv.gs {
				if getOrigMethName(v.g.GetMethName()) == md.m {
					if len(v.targs) > 0 {
						hash := "" // TODO: factor out
						for _, v1 := range v.targs {
							hash = hash + v1.String()
						}
						targs[hash] = v.targs
					}
				}
			}
		}
		for _, v := range targs {
			subs1 := make(map[TParam]Type)
			for k1, v1 := range subs {
				subs1[k1] = v1
			}
			for i := 0; i < len(v); i++ {
				subs1[md.psi.tfs[i].a] = v[i]
			}
			recv := fg.NewParamDecl(md.x_recv, wv.id)
			pds := make([]fg.ParamDecl, len(md.pds))
			for i := 0; i < len(md.pds); i++ {
				tmp := md.pds[i]
				u_p := tmp.u.TSubs(subs1).(TName)
				pds[i] = fg.NewParamDecl(tmp.x, omega[toWKey(u_p)].id)
			}
			u := md.u.TSubs(subs1).(TName)
			e := monomExpr(omega, md.e.TSubs(subs1))
			md1 := fg.NewMDecl(recv, getMonomMethName(omega, md.m, v), pds,
				omega[toWKey(u)].id, e)
			res = append(res, md1)
		}
	}
	return res
}

func monomExpr(omega WMap, e Expr) fg.Expr {
	switch e1 := e.(type) {
	case Variable:
		return fg.NewVariable(e1.id)
	case StructLit:
		es := make([]fg.Expr, len(e1.es))
		for i := 0; i < len(e1.es); i++ {
			es[i] = monomExpr(omega, e1.es[i])
		}
		return fg.NewStructLit(omega[toWKey(e1.u)].id, es)
	case Select:
		return fg.NewSelect(monomExpr(omega, e1.e), e1.f)
	case Call:
		e2 := monomExpr(omega, e1.e)
		var m Name
		if len(e1.targs) == 0 {
			m = e1.m
		} else {
			m = getMonomMethName(omega, e1.m, e1.targs)
		}
		es := make([]fg.Expr, len(e1.args))
		for i := 0; i < len(e1.args); i++ {
			es[i] = monomExpr(omega, e1.args[i])
		}
		return fg.NewCall(e2, m, es)
	case Assert:
		return fg.NewAssert(monomExpr(omega, e1.e),
			omega[toWKey(e1.u.(TName))].id)
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

type WMap map[WKey]WVal
type WEnv map[TParam]Type // Pre: Type is closed

type WKey struct {
	t    Name
	hash string // Hack, represents a closed TName
}

// Pre: isClosed(u)
func toWKey(u TName) WKey {
	hash := ""
	if len(u.us) > 0 {
		hash = u.us[0].String()
		for _, v := range u.us[1:] {
			hash = hash + ",," + v.String()
		}
	}
	return WKey{u.t, hash}
}

type WVal struct {
	u  TName              // Pre: isClosed(u)
	id fg.Type            // Monomorph identifier
	gs map[string]MonoSig // Only records methods with "additional params" // HACK: string key is MonoSig.g.String()
}

func (wv WVal) GetTName() TName {
	return wv.u
}

func (wv WVal) GetMonomId() fg.Type {
	return wv.id
}

func toMonomId(u TName) fg.Type {
	res := u.String()
	res = strings.Replace(res, ",", ",,", -1)
	res = strings.Replace(res, "(", "<", -1)
	res = strings.Replace(res, ")", ">", -1)
	res = strings.Replace(res, " ", "", -1)
	return fg.Type(res)
}

type MonoSig struct {
	g     fg.Sig
	targs []Type // "Additional method type actuals" that give 'g'
	u     TName  // The "actual" return type -- cf. "declared" return type
	// "Actual" return type means the (static) type of body 'e' of the source..
	// ..method under targs (and the TEnv of the parent TDecl instance)
}

// CHECKME: "whole-program" approach starting from the "main" Expr means monom
// may still work in the presence of irregular types and polymorphic
// recursion? -- the cost is, of course, to give up separate compilation
//
// N.B. mutates omega -- i.e., omega is populated with the results
// Pre: 'e' is typeable under an empty TEnv, i.e., does not feature any TParams
func MakeWMap(ds []Decl, gamma ClosedEnv, e Expr, omega WMap) (res TName) {
	var todo []TName // Pre: forall u, isClosed(u)
	// Usage contract: if addTypeToWMap true, then append 'u' to 'todo'

	switch e1 := e.(type) {
	case Variable:
		res = gamma[e1.id]
	case StructLit:
		if addTypeToWMap(e1.u, omega) { // CHECKME: do recursively on e1.u.us?
			todo = append(todo, e1.u) // Cannot refactor inside addTypeToWMap
		}
		for _, v := range e1.es {
			MakeWMap(ds, gamma, v, omega) // Discard return
		}
		res = e1.u
	case Select:
		u_S := MakeWMap(ds, gamma, e1.e, omega)
		for _, v := range fields(ds, u_S) {
			if v.f == e1.f {
				res = v.u.(TName)
				break
			}
		}
	case Call:
		u0 := MakeWMap(ds, gamma, e1.e, omega)
		for _, v := range e1.args {
			MakeWMap(ds, gamma, v, omega) // Discard return
		}
		if isClosed(u0) && len(e1.targs) > 0 {
			isC := true
			for _, v := range e1.targs {
				if u, ok := v.(TName); !ok || !isClosed(u) { // CHECKME: do recursively on targs?
					isC = false
					break
				}
			}
			if isC {
				g := methods(ds, u0)[e1.m]
				subs := make(map[TParam]Type)
				for i := 0; i < len(g.psi.tfs); i++ {
					subs[g.psi.tfs[i].a] = e1.targs[i]
				}
				g = g.TSubs(subs)
				hash := g.String()
				mds := omega[toWKey(u0)].gs // Pre: MakeWMap above ensures u0 in omega
				if tmp, ok := mds[hash]; ok {
					res = tmp.u
				} else {
					m := getMonomMethName(omega, g.m, e1.targs)
					var pds []fg.ParamDecl
					for _, v := range g.pds {
						pds = append(pds, fg.NewParamDecl(v.x, toMonomId(v.u.(TName))))
					}
					mds[hash] = MonoSig{fg.NewSig(m, pds, toMonomId(g.u.(TName))),
						e1.targs, res}
					_, todo1 := visitSig(ds, u0, g, e1.targs, omega)
					todo = append(todo, todo1...)
				}
			}
		}
	case Assert:
		u := e1.u.(TName)
		if addTypeToWMap(u, omega) {
			todo = append(todo, u)
		}
		MakeWMap(ds, gamma, e1.e, omega)
		res = e1.u.(TName) // Factor out
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}

	var empty []Type
	for len(todo) > 0 {
		u := todo[0]
		todo = todo[1:]
		gs := methods(ds, u)
		for _, v := range gs {
			if len(v.psi.tfs) > 0 { // Mutually exclusive with Call counterpart
				continue
			}
			_, todo1 := visitSig(ds, u, v, empty, omega)
			todo = append(todo, todo1...)
		}
	}

	return res
}

// N.B. mutates omega -- adds WKey, WVal pair (if 'u' closed)
func addTypeToWMap(u TName, omega WMap) bool {
	if !isClosed(u) {
		return false
	}
	wk := toWKey(u)
	if _, ok := omega[wk]; ok {
		return false
	}
	omega[wk] = WVal{u, toMonomId(u), make(map[string]MonoSig)}
	return true
}

// N.B. mutates omega -- i.e., omega is populated with the results
func visitSig(ds []Decl, u0 TName, g Sig, targs []Type, omega WMap) (res TName,
	todo []TName) {
	u_ret := g.u.(TName) // Closed, since u0 closed and no meth-params
	wk1 := toWKey(u_ret)
	if _, ok := omega[wk1]; !ok {
		omega[wk1] = WVal{u_ret, toMonomId(u_ret), make(map[string]MonoSig)}
		todo = append(todo, u_ret)
	}
	for _, v := range g.pds {
		u_p := v.u.(TName)
		wk2 := toWKey(u_p)
		if _, ok := omega[wk2]; !ok {
			omega[wk2] = WVal{u_p, toMonomId(u_p), make(map[string]MonoSig)}
			todo = append(todo, u_ret)
		}
	}
	if isStructTName(ds, u0) { // CHECKME: for interface types, visit all possible methods?  Or visiting all struct types already enough?
		x0, xs, e := body(ds, u0, g.m, targs)
		gamma1 := ClosedEnv{x0: u0}
		for i := 0; i < len(xs); i++ {
			gamma1[xs[i]] = g.pds[i].u.(TName)
		}
		res = MakeWMap(ds, gamma1, e, omega)
	} else {
		res = g.u.(TName)
	}
	return res, todo
}

/* Helpers */

func isClosed(u TName) bool {
	for _, v := range u.us {
		if u1, ok := v.(TName); !ok {
			return false
		} else {
			if !isClosed(u1) {
				return false
			}
		}
	}
	return true
}

// Pre: len(targs) > 0
func getMonomMethName(omega WMap, m Name, targs []Type) Name {
	res := m + "<" + string(omega[toWKey(targs[0].(TName))].id)
	for _, v := range targs[1:] {
		res = res + "," + string(omega[toWKey(v.(TName))].id)
	}
	res = res + ">"
	return Name(res)
}

func getOrigMethName(m Name) Name { // Hack
	return m[:strings.Index(m, "<")]
}

/*
// ...debug print WMap
func ... {
	var gamma fgg.ClosedEnv
	omega := make(fgg.WMap)
	fgg.MakeWMap(prog.GetDecls(), gamma, prog.GetExpr().(fgg.Expr), omega)
	for _, v := range omega {
		vPrintln(v.GetTName().String() + " |-> " + string(v.GetMonomId()))
		gs := fgg.GetParameterisedSigs(v)
		if len(gs) > 0 {
			vPrintln("Instantiations of parameterised methods: (i.e., those that had \"additional method params\")")
			for _, g := range gs {
				vPrintln("\t" + g.String())
			}
		}
	}
}

func GetParameterisedSigs(wv WVal) []fg.Sig {
	var res []fg.Sig
	for _, v := range wv.gs {
		res = append(res, v.g)
	}
	return res
}
*/

/* IGNORE */

/*
	// TODO: factor out with Call case
	u_ret := v.u.(TName) // Closed, since u closed and no meth-params
	key1 := toWKey(u_ret)
	if _, ok := omega[key1]; !ok {
		omega[key1] = WVal{u_ret, toMonomId(u_ret), make(map[string]MonoSig)}
		todo = append(todo, u_ret)
	}
	for i := 0; i < len(v.pds); i++ {
		u_p := v.pds[i].u.(TName)
		key2 := toWKey(u_p)
		if _, ok := omega[key2]; !ok {
			omega[key2] = WVal{u_p, toMonomId(u_p), make(map[string]MonoSig)}
			todo = append(todo, u_ret)
		}
	}
	if isStructTName(ds, u) {
		x0, xs, e := body(ds, u, v.m, empty)
		gamma1 := make(Env)
		gamma1[x0] = u
		for i := 0; i < len(xs); i++ {
			gamma1[xs[i]] = v.pds[i].u.(TName)
		}
		MakeWMap(ds, gamma1, e, omega)
	}
*/
