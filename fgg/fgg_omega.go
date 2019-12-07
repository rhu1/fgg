package fgg

import (
	"fmt"
	"reflect"
)

var _ = fmt.Errorf

/* GroundEnv */

// Basically Gamma for TNamed only.
type GroundEnv map[Name]TNamed // Pre: forall TName, isClosed

// Cf. WVal
type Ground struct {
	u  TNamed                 // Pre: isClosed(u)
	gs map[string]GroundEntry // Morally, Sig->[]Type -- HACK: string key is Sig.String
}

// The actual map entry, because g cannot be used as map key directly
type GroundEntry struct {
	g     Sig
	targs []Type
}

// Collects ground types in `ground`.
func fixOmega(ds []Decl, gamma GroundEnv, e FGGExpr, ground map[string]Ground) {
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
							gamma1 := make(GroundEnv)
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

// gamma needed when we're visiting an `e` of a "standalone" meth decl (via collectGroundFggType)
// CHECKME: Post: res already collected?
func collectGroundFggTypes(ds []Decl, gamma GroundEnv, e FGGExpr, ground map[string]Ground) (res Type) {
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
		// !!! we don't just visit e1.e_S, we also visit the type of e_S
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
		// !!! CHECKME: "actual" vs. "declared -- declared is highest, most exhaustive? -- or need both?
		res = e1.Typing(ds, make(Delta), gamma1, true) // CHECKME: typing vs. sig? -- CHECKME: currently this typing mixed with res
		//res = g.u // May be a TParam, e.g., `Cond(type a Any())(br Branches(a)) a` (map.fgg)
	case Assert:
		u := e1.u_cast.(TNamed) // CHECKME: guaranteed?
		collectGroundFggType(ds, u, ground)
		collectGroundFggTypes(ds, gamma, e1.e_I, ground)
		res = u
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
	return res
}

func collectGroundFggType(ds []Decl, u Type, ground map[string]Ground) {
	if _, ok := ground[u.String()]; ok {
		return
	}
	if cast, ok := u.(TNamed); !ok || !isGround(cast) {
		return
	}

	u1 := u.(TNamed)

	gs := make(map[string]GroundEntry)
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
				gamma := make(GroundEnv)
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
// can proceed when u0 is ground without a Delta as we also have add type args here
func collectGroundFggCall(ds []Decl, u0 Type, c Call, ground map[string]Ground) {
	if cast, ok := u0.(TNamed); !ok || !isGround(cast) {
		return
	}
	for _, v := range c.t_args {
		if cast, ok := v.(TNamed); !ok || !isGround(cast) {
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
	gs[g.String()] = GroundEntry{g, c.t_args}
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
		gamma1 := make(GroundEnv)
		gamma1[x0] = u_S
		for i := 0; i < len(xs); i++ { // xs = ys in pds
			gamma1[xs[i]] = pds[i].GetType().TSubs(subs).(TNamed) // Param names in g should be same as actual MDecl
		}
		collectGroundFggTypes(ds, gamma1, e, ground)
	} else {
		// visit all possible bodies
	}
}
