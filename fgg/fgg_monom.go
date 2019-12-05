package fgg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rhu1/fgg/fg"
)

var _ = fmt.Errorf

/* Naive monomorph -- !!WIP!! */

// CHECKME: -monom skips any meth with add-param that isn't instantiated?  is that sound? (any previously failing cast now permitted by this "relaxed interface constraint"? i.e., increased duck typing potential)

type ClosedEnv map[Name]TNamed // Pre: forall TName, isClosed

// func isMonomorphisable(p FGGProgram) bool { ... }

// TODO: reformat (e.g., "<...>") to make an actual FG program
func Monomorph(p FGGProgram) fg.FGProgram {
	//var gamma ClosedEnv // CHECKME: nil map -- so never used?
	omega := make(WMap)
	//MakeWMap(p.GetDecls(), gamma, p.GetExpr().(Expr), omega) // Populates omega

	//fmt.Println("1111:\n", omega)
	var gamma1 ClosedEnv
	ground := make(map[string]Ground)
	//collectGroundFggTypes(p.GetDecls(), gamma1, p.GetExpr().(Expr), ground)
	fix(p.GetDecls(), gamma1, p.GetMain().(FGGExpr), ground)
	MakeWMap2(p.GetDecls(), ground, omega)
	//fmt.Println("2222:\n", ground)

	var ds []Decl
	for _, v := range p.decls {
		switch d := v.(type) {
		case TDecl:
			t := d.GetName()
			for k1, v1 := range omega { // CHECKME: "prunes" unused types -- OK?
				if k1.t == t {
					ds = append(ds, monomTDecl(p.decls, omega, d, v1))
				}
			}
		case MDecl:
			for k1, v1 := range omega { // CHECKME: "prunes" unused types -- OK?
				if k1.t == d.t_recv {
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

// Pre: `wv` represents an instantiation of the `td` type  // TODO: refactor, decompose
func monomTDecl(ds []Decl, omega WMap, td TDecl, wv WVal) fg.TDecl {
	subs := make(map[TParam]Type) // Type is a TName
	psi := td.GetPsi()
	for i := 0; i < len(psi.tFormals); i++ {
		subs[psi.tFormals[i].name] = wv.u.u_args[i]
	}
	switch d := td.(type) {
	case STypeLit:
		fds := make([]fg.FieldDecl, len(d.fDecls))
		for i := 0; i < len(d.fDecls); i++ {
			tmp := d.fDecls[i]
			u := tmp.u.TSubs(subs).(TNamed)     // "Inlined" substitution actions here -- cf. TDecl.TSubs
			if _, ok := omega[toWKey(u)]; !ok { // Cf. MakeWMap2, extra loop over non-param TDecls, for those non seen o/w
				panic("Unknown type: " + u.String())
			}
			fds[i] = fg.NewFieldDecl(tmp.field, omega[toWKey(u)].id)
		}
		return fg.NewSTypeLit(wv.id, fds)
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
						pds[i] = fg.NewParamDecl(tmp.name, omega[toWKey(u_p)].id)
					}
					u := s.u_ret.TSubs(subs).(TNamed)
					ss = append(ss, fg.NewSig(s.meth, pds, omega[toWKey(u)].id))
				} else {
					// forall u_S s.t. u_S <: wv.u, collect m.targs for all wv.m and mono(u_S.m)
					// ^Correction: forall u, not only u_S, i.e., including interface type receivers
					// (Cf. map.fgg, Bool().Cond(Bool())(...))
					gs := methods(ds, wv.u)
					empty := make(Delta)
					targs := make(map[string][]Type)
					for _, v := range omega {
						if /*IsStructType(ds, v.u.t) &&*/ v.u.Impls(ds, empty, wv.u) { // N.B. now adding reflexively
							// Collect meth instans from *all* subtypes, i.e., including calls on interface receivers
							for _, v1 := range gs {
								addMethInstans(v, v1.meth, targs)
							}
							/*for _, v1 := range v.gs {
								// TODO: factor out with addMethInstans? but here, adding all m's, not filtering
								m1 := getOrigMethName(v1.g.GetMethName())
								if _, ok := gs[m1]; ok && len(v1.targs) > 0 { // len check redundant?
									hash := "" // Use WriteTypes?
									for _, v2 := range v1.targs {
										hash = hash + v2.String()
									}
									targs[hash] = v1.targs
								}
							}*/
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
							pds[i] = fg.NewParamDecl(tmp.name, omega[toWKey(u_p)].id)
						}
						u := s.u_ret.TSubs(subs1).(TNamed)
						g1 := fg.NewSig(getMonomMethName(omega, s.meth, v), pds,
							omega[toWKey(u)].id)
						ss = append(ss, g1)
					}
				}
			case TNamed:
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
	for i := 0; i < len(md.psi_recv.tFormals); i++ {
		subs[md.psi_recv.tFormals[i].name] = wv.u.u_args[i]
	}
	recv := fg.NewParamDecl(md.x_recv, wv.id)
	if len(md.psi_meth.tFormals) == 0 {
		pds := make([]fg.ParamDecl, len(md.pDecls))
		for i := 0; i < len(md.pDecls); i++ {
			tmp := md.pDecls[i]
			u := tmp.u.TSubs(subs).(TNamed) // "Inlined" substitution actions here -- cf. TDecl.TSubs
			pds[i] = fg.NewParamDecl(tmp.name, omega[toWKey(u)].id)
		}
		t := omega[toWKey(md.u_ret.TSubs(subs).(TNamed))].id
		e := monomExpr(omega, md.e_body.TSubs(subs))
		res = append(res, fg.NewMDecl(recv, md.name, pds, t, e))
	} else {
		targs := collectZigZagMethInstans(ds, omega, md, wv)
		if len(targs) == 0 { // Means no u_I, if len(wv.gs)>0 -- targs doesn't (yet) include wv.gs
			addMethInstans(wv, md.name, targs)
		}
		for _, v := range targs { // CHECKME: factor out with ITypeLit?
			subs1 := make(map[TParam]Type)
			for k1, v1 := range subs {
				subs1[k1] = v1
			}
			for i := 0; i < len(v); i++ {
				subs1[md.psi_meth.tFormals[i].name] = v[i]
			}
			recv := fg.NewParamDecl(md.x_recv, wv.id)
			pds := make([]fg.ParamDecl, len(md.pDecls))
			for i := 0; i < len(md.pDecls); i++ {
				tmp := md.pDecls[i]
				u_p := tmp.u.TSubs(subs1).(TNamed)
				pds[i] = fg.NewParamDecl(tmp.name, omega[toWKey(u_p)].id)
			}
			u := md.u_ret.TSubs(subs1).(TNamed)
			e := monomExpr(omega, md.e_body.TSubs(subs1))
			md1 := fg.NewMDecl(recv, getMonomMethName(omega, md.name, v), pds,
				omega[toWKey(u)].id, e)
			res = append(res, md1)
		}
	}
	return res
}

// N.B. return is empty, i.e., does not include wv.gs, if no u_I
// N.B. return is a map, so "duplicate" add-meth-param type instans are implicitly setify-ed
// ^E.g., Calling m(A()) on some struct separately via two interfaces T1 and T2 where T2 <: T1
func collectZigZagMethInstans(ds []Decl, omega WMap, md MDecl, wv WVal) map[string][]Type {
	empty := make(Delta)
	targs := make(map[string][]Type)
	// Given m = md.m, forall u_I s.t. m in meths(u_I) && wv.u <: u_I, ..
	// ..forall u_S s.t. u_S <: u_I, collect targs for all mono(u_S.m)
	// ^Correction: forall u, not only u_S
	for _, v := range omega {
		if IsNamedIfaceType(ds, v.u) && wv.u.Impls(ds, empty, v.u) {
			gs := methods(ds, v.u)
			if _, ok := gs[md.name]; ok {
				addMethInstans(v, md.name, targs)
				for _, v1 := range omega {
					if /*isStructTName(ds, v1.u) &&*/ v1.u.Impls(ds, empty, v.u) {
						addMethInstans(v1, md.name, targs)
					}
				}
			}
		}
	}
	return targs
}

// Add meth instans from `wv`, filtered by `m`, to `targs`
func addMethInstans(wv WVal, m Name, targs map[string][]Type) {
	for _, v := range wv.gs {
		m1 := getOrigMethName(v.g.GetName())
		if m1 == m && len(v.targs) > 0 {
			hash := "" // Use WriteTypes?
			for _, v1 := range v.targs {
				hash = hash + v1.String()
			}
			targs[hash] = v.targs
		}
	}
}

func monomExpr(omega WMap, e FGGExpr) fg.FGExpr {
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
			//panic("Unknown type: " + e1.u.String())
		}
		return fg.NewStructLit(omega[wk].id, es)
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
		return fg.NewAssert(monomExpr(omega, e1.expr),
			omega[wk].id)
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
func toWKey(u TNamed) WKey {
	hash := ""
	if len(u.u_args) > 0 {
		hash = u.u_args[0].String()
		for _, v := range u.u_args[1:] {
			hash = hash + ",," + v.String()
		}
	}
	return WKey{u.t_name, hash}
}

type WVal struct {
	u  TNamed             // Pre: isClosed(u)
	id fg.Type            // Monomorph identifier
	gs map[string]MonoSig // Only records methods with "additional params" // HACK: string key is MonoSig.g.String()
}

func (wv WVal) GetTName() TNamed {
	return wv.u
}

func (wv WVal) GetMonomId() fg.Type {
	return wv.id
}

func toMonomId(u TNamed) fg.Type {
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
	u     TNamed // The "actual" return type -- cf. "declared" return type
	// "Actual" return type means the (static) type of body 'e' of the source..
	// ..method under targs (and the TEnv of the parent TDecl instance)
}

// @Deprecated
// CHECKME: "whole-program" approach starting from the "main" Expr means monom
// may somtimes work in the presence of irregular types and polymorphic
// recursion? -- the cost is, of course, to give up separate compilation
//
// N.B. mutates omega -- i.e., omega is populated with the results
// Pre: `e` is typeable under an empty TEnv, i.e., does not feature any TParams
func MakeWMap(ds []Decl, gamma ClosedEnv, e FGGExpr, omega WMap) (res Type) {
	var todo []TNamed // Pre: forall u, isClosed(u)
	// Usage contract: if addTypeToWMap true, then append `u` to `todo`

	switch e1 := e.(type) {
	case Variable:
		res = gamma[e1.name]
	case StructLit:
		if addTypeToWMap(e1.u_S, omega) { // CHECKME: do recursively on e1.u.us?
			todo = append(todo, e1.u_S) // Cannot refactor inside addTypeToWMap
		}
		fds := fields(ds, e1.u_S)
		for _, fd := range fds {
			u := fd.u.(TNamed)
			if addTypeToWMap(u, omega) {
				todo = append(todo, u)
			}
		}
		for _, v := range e1.elems {
			MakeWMap(ds, gamma, v, omega) // Discard return -- recursive recording done inside, here we want the decl type fd.u as done above
		}
		res = e1.u_S
	case Select:
		u_S := MakeWMap(ds, gamma, e1.e_S, omega).(TNamed)
		for _, v := range fields(ds, u_S) {
			if v.field == e1.field {
				res = v.u.(TNamed)
				break
			}
		}
	case Call:
		u0 := MakeWMap(ds, gamma, e1.e_recv, omega)
		for _, v := range e1.args {
			MakeWMap(ds, gamma, v, omega) // Discard return
		}
		g := methods(ds, u0)[e1.meth]
		res = g.u_ret // May be a TParam, e.g., `Cond(type a Any())(br Branches(a)) a` (map.fgg) -- then below is skipped
		if u0_closed, ok := u0.(TNamed); ok && isClosed(u0_closed) &&
			len(e1.t_args) > 0 {
			isC := true
			for _, v := range e1.t_args {
				if u, ok := v.(TNamed); !ok || !isClosed(u) { // CHECKME: do recursively on targs?
					isC = false
					break
				}
			}
			if isC {
				subs := make(map[TParam]Type)
				for i := 0; i < len(g.psi.tFormals); i++ {
					subs[g.psi.tFormals[i].name] = e1.t_args[i]
				}
				g_subs := g.TSubs(subs)
				hash := g_subs.String()
				mds := omega[toWKey(u0_closed)].gs // Pre: MakeWMap above ensures u0 in omega
				if tmp, ok := mds[hash]; ok {
					res = tmp.u
				} else {
					m := getMonomMethName(omega, g_subs.meth, e1.t_args)
					var pds []fg.ParamDecl
					for _, v := range g_subs.pDecls {
						pds = append(pds, fg.NewParamDecl(v.name, toMonomId(v.u.(TNamed))))
					}
					res = g_subs.u_ret.(TNamed)
					mds[hash] = MonoSig{fg.NewSig(m, pds, toMonomId(res.(TNamed))),
						e1.t_args, res.(TNamed)}
					_, todo1 := visitSig(ds, u0_closed, g_subs, e1.t_args, omega)
					todo = append(todo, todo1...)
				}
			}
		}
	case Assert:
		u := e1.u_cast.(TNamed)
		if addTypeToWMap(u, omega) {
			todo = append(todo, u)
		}
		MakeWMap(ds, gamma, e1.expr, omega)
		res = e1.u_cast.(TNamed) // Factor out
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
			if len(v.psi.tFormals) > 0 { // Mutually exclusive with Call counterpart
				continue
			}
			_, todo1 := visitSig(ds, u, v, empty, omega)
			todo = append(todo, todo1...)
		}
	}

	return res
}

// N.B. mutates omega -- adds WKey, WVal pair (if `u` closed)
// @return `true` if type added, `false` o/w
func addTypeToWMap(u TNamed, omega WMap) bool {
	if !isClosed(u) { // CHECKME: necessary?
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
// Pre: isClosed(u0), g is ground
func visitSig(ds []Decl, u0 TNamed, g Sig, targs []Type, omega WMap) (res TNamed,
	todo []TNamed) {
	u_ret := g.u_ret.(TNamed) // Closed, since u0 closed and no meth-params
	wk1 := toWKey(u_ret)
	if _, ok := omega[wk1]; !ok {
		omega[wk1] = WVal{u_ret, toMonomId(u_ret), make(map[string]MonoSig)}
		todo = append(todo, u_ret)
	}
	for _, v := range g.pDecls {
		u_p := v.u.(TNamed)
		wk2 := toWKey(u_p)
		if _, ok := omega[wk2]; !ok {
			omega[wk2] = WVal{u_p, toMonomId(u_p), make(map[string]MonoSig)}
			todo = append(todo, u_ret)
		}
	}
	if IsStructType(ds, u0) { // CHECKME: for interface types, visit all possible methods?  Or visiting all struct types already enough?
		x0, xs, e := body(ds, u0, g.meth, targs)
		gamma1 := ClosedEnv{x0: u0}
		for i := 0; i < len(xs); i++ {
			gamma1[xs[i]] = g.pDecls[i].u.(TNamed)
		}
		res = MakeWMap(ds, gamma1, e, omega).(TNamed) // isClosed(u0), g is ground
	} else {
		res = u_ret
	}
	return res, todo
}

/* Helpers */

func isClosed(u TNamed) bool {
	for _, v := range u.u_args {
		if u1, ok := v.(TNamed); !ok {
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
	res := m + "<" + string(omega[toWKey(targs[0].(TNamed))].id)
	for _, v := range targs[1:] {
		res = res + "," + string(omega[toWKey(v.(TNamed))].id)
	}
	res = res + ">"
	return Name(res)
}

func getOrigMethName(m Name) Name { // Hack
	return m[:strings.Index(m, "<")]
}

/****************/

func MakeWMap2(ds []Decl, ground map[string]Ground, omega WMap) {
	for _, v := range ground {
		wk := toWKey(v.u)
		gs := make(map[string]MonoSig)
		omega[wk] = WVal{v.u, toMonomId(v.u), gs}
		/*}

		for _, v := range ground {*/
		for _, pair := range v.gs {
			if len(pair.targs) == 0 {
				continue
			}
			hash := pair.g.String()
			pds := pair.g.GetParamDecls()
			pds_fg := make([]fg.ParamDecl, len(pds))
			for i := 0; i < len(pds); i++ {
				pd := pds[i]
				pds_fg[i] = fg.NewParamDecl(pd.name, toMonomId(pd.u.(TNamed)))
			}
			ret := pair.g.u_ret.(TNamed)
			m := getMonomMethName(omega, pair.g.meth, pair.targs)
			//gs := omega[toWKey(v.u)].gs
			gs[hash] = MonoSig{fg.NewSig(m, pds_fg, toMonomId(ret)), pair.targs, ret}
		}
	}
}

/*for _, d1 := range ds {
	switch d := d1.(type) {
	case STypeLit:
		if len(d.GetTFormals().GetFormals()) == 0 {
			u := TName{d.t, []Type{}}
			wk := toWKey(u)
			if _, ok := omega[wk]; ok {
				continue
			}
			omega[wk] = WVal{u, toMonomId(u), make(map[string]MonoSig)}
		}
	case ITypeLit:
		// CHECKME: meth decls?
	}
}*/

// Cf. WVal
type Ground struct {
	u  TNamed                // Pre: isClosed(u)
	gs map[string]GroundPair // // HACK: string key is Sig.String
}

type GroundPair struct {
	g     Sig
	targs []Type
}

func fix(ds []Decl, gamma ClosedEnv, e FGGExpr, ground map[string]Ground) {
	empty := make(Delta)

	again := true
	for again {
		again = false
		collectGroundFggTypes(ds, gamma, e, ground)

		for _, v := range ground {
			if IsNamedIfaceType(ds, v.u) {

				if len(v.gs) == 0 {
					continue
				}

				for _, v1 := range ground {
					if IsStructType(ds, v1.u) {

						if !v1.u.Impls(ds, empty, v.u) {
							continue
						}

						u_S := v1.u

						for _, gp := range v.gs {

							if len(gp.targs) == 0 {
								continue
							}

							subs := make(map[TParam]Type)
							td := GetTDecl(ds, u_S.GetName())
							targs := u_S.GetTArgs()
							tfs := td.GetPsi().GetTFormals()
							for i := 0; i < len(targs); i++ {
								subs[tfs[i].name] = targs[i]
							}
							tfs_c := gp.g.GetPsi().GetTFormals()
							for i := 0; i < len(tfs_c); i++ {
								subs[tfs_c[i].name] = gp.targs[i]
							}

							var pds []ParamDecl = nil
							for _, d := range ds {
								if md, ok := d.(MDecl); ok {
									if md.t_recv == v1.u.t_name && md.name == gp.g.meth {
										pds = md.pDecls
										break
									}
								}
							}
							if pds == nil {
								panic("Method not found on " + v1.u.String() + ": " + gp.g.meth)
							}

							x0, xs, e := body(ds, u_S, gp.g.meth, gp.targs)
							gamma1 := make(ClosedEnv)
							gamma1[x0] = u_S
							for i := 0; i < len(xs); i++ { // xs = ys in pds
								gamma1[xs[i]] = pds[i].GetType().TSubs(subs).(TNamed)
							}

							ground1 := make(map[string]Ground)
							collectGroundFggTypes(ds, gamma1, e, ground1)
							for _, gp2 := range ground1 {
								if _, ok := ground[gp2.u.String()]; !ok {
									ground[gp2.u.String()] = gp2
									again = true
								}
							}
						}
					}
				}
			}
		}
	}
}

// gamma needed when we're visiting an e of a "standalone" meth decl (via collectGroundFggType)
// CHECKME: Post: res already collected?
func collectGroundFggTypes(ds []Decl, gamma ClosedEnv, e FGGExpr, ground map[string]Ground) (res Type) {
	//fmt.Println("444:", e)
	switch e1 := e.(type) {
	case Variable:
		res = gamma[e1.name]
	case StructLit:
		collectGroundFggType(ds, e1.u_S, ground)
		for _, v := range e1.elems {
			collectGroundFggTypes(ds, gamma, v, ground) // Discard return
		}
		res = e1.u_S
	case Select:
		u_S := collectGroundFggTypes(ds, gamma, e1.e_S, ground).(TNamed) // Field types already collected via the structlit?
		collectGroundFggType(ds, u_S, ground)
		for _, v := range fields(ds, u_S) {
			if v.field == e1.field {
				res = v.u.(TNamed)
				break
			}
		}
	case Call:
		u0 := collectGroundFggTypes(ds, gamma, e1.e_recv, ground)
		for _, v := range e1.t_args {
			collectGroundFggType(ds, v, ground)
		}
		for _, v := range e1.args {
			collectGroundFggTypes(ds, gamma, v, ground) // Discard return
		}
		collectGroundFggCall(ds, u0, e1, ground)

		gamma1 := make(Gamma)
		for k, v := range gamma {
			gamma1[k] = v
		}
		res = e1.Typing(ds, make(Delta), gamma1, true) // CHECKME: typing vs. sig? -- CHECKME: currently this typing mixed with res
		//res = g.u // May be a TParam, e.g., `Cond(type a Any())(br Branches(a)) a` (map.fgg)
	case Assert:
		u := e1.u_cast.(TNamed) // CHECKME: guaranteed?
		collectGroundFggType(ds, u, ground)
		collectGroundFggTypes(ds, gamma, e1.expr, ground)
		res = u
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
	return res
}

func collectGroundFggType(ds []Decl, u Type, ground map[string]Ground) {
	//fmt.Println("555:", u)
	if _, ok := ground[u.String()]; ok {
		return
	}
	if cast, ok := u.(TNamed); !ok || !isClosed(cast) {
		return
	}

	u1 := u.(TNamed)

	gs := make(map[string]GroundPair)
	ground[u.String()] = Ground{u1, gs}
	if IsStructType(ds, u1) {
		u_S := u1

		fds := fields(ds, u_S)
		for _, fd := range fds {
			u := fd.u.(TNamed)
			collectGroundFggType(ds, u, ground)
		}

		// visit meths
		gs := methods(ds, u_S)
		for _, g := range gs {
			// visit types in sig
			pds := g.GetParamDecls()
			for i := 0; i < len(pds); i++ {
				u_pd := pds[i].GetType()
				collectGroundFggType(ds, u_pd, ground)
			}
			collectGroundFggType(ds, g.u_ret, ground)

			// visit body
			if len(g.GetPsi().GetTFormals()) == 0 {
				x_recv, xs, e := body(ds, u_S, g.meth, []Type{})
				gamma := make(ClosedEnv)
				gamma[x_recv] = u_S
				for i := 0; i < len(pds); i++ {
					gamma[xs[i]] = pds[i].GetType().(TNamed)
				}
				collectGroundFggTypes(ds, gamma, e, ground)
			}
		}

		// check all super interfaces, and visit all meths of sub structs (recursively) -- no
	} else { // interface
		u_I := u1

		// visit meths
		gs := methods(ds, u_I)
		for _, g := range gs {
			// visit types in sig // TODO: duplicated from above, factor out
			pds := g.GetParamDecls()
			for i := 0; i < len(pds); i++ {
				u_pd := pds[i].GetType()
				collectGroundFggType(ds, u_pd, ground)
			}
			collectGroundFggType(ds, g.u_ret, ground)
		}

		// visit embedded
		td := GetTDecl(ds, u_I.t_name).(ITypeLit)
		tfs := td.GetPsi().GetTFormals()
		subs := make(map[TParam]Type)
		for i := 0; i < len(u_I.u_args); i++ {
			subs[tfs[i].name] = u_I.u_args[i]
		}
		for _, s := range td.specs {
			if u, ok := s.(TNamed); ok {
				collectGroundFggType(ds, u.TSubs(subs), ground)
			}
		}

		// visit all meths of sub structs -- no
	}
}

// Pre: if u0 is ground, then already in `ground` -- no XXX
func collectGroundFggCall(ds []Decl, u0 Type, c Call, ground map[string]Ground) {
	//fmt.Println("666:", u0, c)
	if cast, ok := u0.(TNamed); !ok || !isClosed(cast) {
		return
	}
	for _, v := range c.t_args {
		if cast, ok := v.(TNamed); !ok || !isClosed(cast) {
			return
		}
	}
	if _, ok := ground[u0.String()]; !ok {
		collectGroundFggType(ds, u0, ground)
	}

	g := methods(ds, u0)[c.meth]
	//if len(c.targs) > 0 {
	subs := make(map[TParam]Type)
	for i := 0; i < len(g.psi.tFormals); i++ {
		subs[g.psi.tFormals[i].name] = c.t_args[i]
	}
	g = g.TSubs(subs)
	gs := ground[u0.String()].gs
	if _, ok := gs[g.String()]; ok {
		return
	}
	/*targs := make([]Type, len(c.targs)) // CHECKME: unnecessary to copy?
	copy(targs, c.targs)*/
	gs[g.String()] = GroundPair{g, c.t_args}
	//}
	pds := g.GetParamDecls()
	for _, v := range pds {
		collectGroundFggType(ds, v.GetType(), ground)
	}
	psi_g := g.GetPsi()
	for _, v := range psi_g.GetTFormals() {
		collectGroundFggType(ds, v.GetUpperBound(), ground)
	}
	collectGroundFggType(ds, g.GetType(), ground)

	if IsStructType(ds, u0) {
		u_S := u0.(TNamed)

		subs := make(map[TParam]Type)
		td := GetTDecl(ds, u_S.GetName())
		targs := u_S.GetTArgs()
		tfs := td.GetPsi().GetTFormals()
		for i := 0; i < len(targs); i++ {
			subs[tfs[i].name] = targs[i]
		}
		tfs_c := psi_g.GetTFormals()
		for i := 0; i < len(tfs_c); i++ {
			subs[tfs_c[i].name] = c.t_args[i]
		}

		x0, xs, e := body(ds, u_S, c.meth, c.t_args)
		gamma1 := make(ClosedEnv)
		gamma1[x0] = u_S
		for i := 0; i < len(xs); i++ { // xs = ys in pds
			gamma1[xs[i]] = pds[i].GetType().TSubs(subs).(TNamed) // Param names in g should be same as actual MDecl
		}
		collectGroundFggTypes(ds, gamma1, e, ground)
	} else {
		// visit all possible bodies
	}
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
