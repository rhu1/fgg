package fgg

import (
	"fmt"
	"reflect"
)

var _ = fmt.Errorf

/* GroundEnv, GroundTypeAndSigs, GroundEntry */

// Basically Gamma for TNamed only.
type GroundEnv map[Name]TNamed // Pre: forall TName, isGround

// Cf. MonomTypeAndSigs
type GroundTypeAndSigs struct {
	u_ground TNamed               // Pre: isGround(u_ground)
	sigs     map[string]GroundSig // Morally, Sig->[]Type -- HACK: string key is Sig.String
	// ^(a) FGG sigs; (b) all sigs on u_ground receiver, including empty add-meth-targs (cf. fgg_monom.go)
}

// The actual map entry, because sig cannot be used as map key directly
type GroundSig struct {
	sig   Sig // May only need meth name given receiver type, but Sig is convenient(?)
	targs []Type
}

/* Attempt to reach a closure on ground types */

// N.B. mutates `ground` -- encountered ground types collected into `ground`.
func fixOmega(ds []Decl, gamma GroundEnv, e FGGExpr,
	ground map[string]GroundTypeAndSigs) {
	collectGroundTypesFromExpr(ds, gamma, e, ground)

	empty := make(Delta)
	again := true
	for again {
		again = false

		for _, v_I := range ground {
			if !IsNamedIfaceType(ds, v_I.u_ground) || len(v_I.sigs) == 0 {
				continue
			}
			for _, v_S := range ground {
				if !IsStructType(ds, v_S.u_ground) ||
					!v_S.u_ground.Impls(ds, empty, v_I.u_ground) {
					continue
				}

				u_S := v_S.u_ground
				for _, ge := range v_I.sigs {
					if len(ge.targs) == 0 {
						continue
					}

					subs := make(map[TParam]Type)
					td := GetTDecl(ds, u_S.GetName())
					targs := u_S.GetTArgs()
					tfs := td.GetPsi().GetTFormals()
					for i := 0; i < len(targs); i++ {
						subs[tfs[i].name] = targs[i]
					}
					tfs_c := ge.sig.GetPsi().GetTFormals()
					for i := 0; i < len(tfs_c); i++ {
						subs[tfs_c[i].name] = ge.targs[i]
					}

					var pds []ParamDecl = nil
					for _, d := range ds {
						if md, ok := d.(MDecl); ok {
							if md.t_recv == v_S.u_ground.t_name && md.name == ge.sig.meth {
								pds = md.pDecls
								break
							}
						}
					}
					if pds == nil {
						panic("Method not found on " + v_S.u_ground.String() + ": " +
							ge.sig.meth)
					}

					x0, xs, e := body(ds, u_S, ge.sig.meth, ge.targs)
					gamma1 := make(GroundEnv)
					gamma1[x0] = u_S
					for i := 0; i < len(xs); i++ { // xs = ys in pds
						gamma1[xs[i]] = pds[i].GetType().TSubs(subs).(TNamed)
					}

					ground1 := make(map[string]GroundTypeAndSigs)
					collectGroundTypesFromExpr(ds, gamma1, e, ground1)
					for _, ge1 := range ground1 {
						if _, ok := ground[ge1.u_ground.String()]; !ok {
							ground[ge1.u_ground.String()] = ge1
							again = true
						}
					}
				}
			}
		}
	}
}

// gamma needed when we're visiting an `e` of a "standalone" meth decl (via collectGroundFggType)
// CHECKME: Post: res already collected?
func collectGroundTypesFromExpr(ds []Decl, gamma GroundEnv, e FGGExpr,
	ground map[string]GroundTypeAndSigs) (res Type) {
	switch e1 := e.(type) {
	case Variable:
		res = gamma[e1.name]
	case StructLit:
		collectGroundTypesFromType(ds, e1.u_S, ground)
		for _, v := range e1.elems {
			collectGroundTypesFromExpr(ds, gamma, v, ground) // Discard return
		}
		res = e1.u_S
	case Select:
		u_S := collectGroundTypesFromExpr(ds, gamma, e1.e_S, ground).(TNamed) // Field types already collected via the structlit?
		// !!! we don't just visit e1.e_S, we also visit the type of e_S
		collectGroundTypesFromType(ds, u_S, ground)
		for _, v := range fields(ds, u_S) {
			if v.field == e1.field {
				res = v.u.(TNamed)
				break
			}
		}
	case Call:
		u0 := collectGroundTypesFromExpr(ds, gamma, e1.e_recv, ground)
		for _, v := range e1.t_args {
			collectGroundTypesFromType(ds, v, ground)
		}
		for _, v := range e1.args {
			collectGroundTypesFromExpr(ds, gamma, v, ground) // Discard return
		}
		collectGroundTypesByVisitingCall(ds, u0, e1, ground)

		gamma1 := make(Gamma)
		for k, v := range gamma {
			gamma1[k] = v
		}
		// !!! CHECKME: "actual" vs. "declared -- declared is highest, most exhaustive? -- or need both?
		res = e1.Typing(ds, make(Delta), gamma1, true) // CHECKME: typing vs. sig? -- CHECKME: currently this typing mixed with res
		//res = g.u // May be a TParam, e.g., `Cond(type a Any())(br Branches(a)) a` (map.fgg)
	case Assert:
		u := e1.u_cast.(TNamed) // CHECKME: guaranteed?
		collectGroundTypesFromType(ds, u, ground)
		collectGroundTypesFromExpr(ds, gamma, e1.e_I, ground)
		res = u
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
	return res
}

func collectGroundTypesFromType(ds []Decl, u Type, ground map[string]GroundTypeAndSigs) {
	if _, ok := ground[u.String()]; ok {
		return
	}
	if cast, ok := u.(TNamed); !ok || !isGround(cast) {
		return
	}

	u1 := u.(TNamed)
	gs := make(map[string]GroundSig)
	ground[u.String()] = GroundTypeAndSigs{u1, gs}
	if IsStructType(ds, u1) {
		u_S := u1

		fds := fields(ds, u_S)
		for _, fd := range fds {
			u := fd.u.(TNamed)
			collectGroundTypesFromType(ds, u, ground)
		}

		// Visit meths
		gs := methods(ds, u_S)
		for _, g := range gs {
			// Visit types in sig
			pds := g.GetParamDecls()
			for i := 0; i < len(pds); i++ {
				u_pd := pds[i].GetType()
				collectGroundTypesFromType(ds, u_pd, ground)
			}
			collectGroundTypesFromType(ds, g.u_ret, ground)

			// Visit body
			if len(g.GetPsi().GetTFormals()) == 0 {
				x_recv, xs, e := body(ds, u_S, g.meth, []Type{})
				gamma := make(GroundEnv)
				gamma[x_recv] = u_S
				for i := 0; i < len(pds); i++ {
					gamma[xs[i]] = pds[i].GetType().(TNamed)
				}
				collectGroundTypesFromExpr(ds, gamma, e, ground)
			}
		}

		// CHECKME: check all super interfaces, and visit all meths of sub structs (recursively)? -- no

	} else { // Interface
		u_I := u1

		// Visit meths
		gs := methods(ds, u_I)
		for _, g := range gs {
			// Visit types in sig // TODO: duplicated from above, factor out
			pds := g.GetParamDecls()
			for i := 0; i < len(pds); i++ {
				u_pd := pds[i].GetType()
				collectGroundTypesFromType(ds, u_pd, ground)
			}
			collectGroundTypesFromType(ds, g.u_ret, ground)
		}

		// Visit embedded
		td := GetTDecl(ds, u_I.t_name).(ITypeLit)
		tfs := td.GetPsi().GetTFormals()
		subs := make(map[TParam]Type)
		for i := 0; i < len(u_I.u_args); i++ {
			subs[tfs[i].name] = u_I.u_args[i]
		}
		for _, s := range td.specs {
			if u, ok := s.(TNamed); ok {
				collectGroundTypesFromType(ds, u.TSubs(subs), ground)
			}
		}

		// Visit all meths of subtype structs? -- no
	}
}

// Pre: if u0 is ground, then already in `ground` -- no XXX
// can proceed when u0 is ground without a Delta as we also have add type args here
func collectGroundTypesByVisitingCall(ds []Decl, u0 Type, c Call,
	ground map[string]GroundTypeAndSigs) {

	if cast, ok := u0.(TNamed); !ok || !isGround(cast) {
		return
	}
	for _, v := range c.t_args {
		if cast, ok := v.(TNamed); !ok || !isGround(cast) {
			return
		}
	}
	if _, ok := ground[u0.String()]; !ok {
		collectGroundTypesFromType(ds, u0, ground)
	}

	g := methods(ds, u0)[c.meth]
	//if len(c.targs) > 0 {
	subs := make(map[TParam]Type)
	for i := 0; i < len(g.psi.tFormals); i++ {
		subs[g.psi.tFormals[i].name] = c.t_args[i]
	}
	g = g.TSubs(subs)
	gs := ground[u0.String()].sigs
	if _, ok := gs[g.String()]; ok {
		return
	}
	gs[g.String()] = GroundSig{g, c.t_args}
	//}
	pds := g.GetParamDecls()
	for _, v := range pds {
		collectGroundTypesFromType(ds, v.GetType(), ground)
	}
	psi_g := g.GetPsi()
	for _, v := range psi_g.GetTFormals() {
		collectGroundTypesFromType(ds, v.GetUpperBound(), ground)
	}
	collectGroundTypesFromType(ds, g.GetType(), ground)

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
		gamma1 := make(GroundEnv)
		gamma1[x0] = u_S
		for i := 0; i < len(xs); i++ { // xs = ys in pds
			gamma1[xs[i]] = pds[i].GetType().TSubs(subs).(TNamed) // Param names in g should be same as actual MDecl
		}
		collectGroundTypesFromExpr(ds, gamma1, e, ground)
	} else {
		// CHECME: visit all possible bodies -- now subsumed by fixOmega?
	}
}
