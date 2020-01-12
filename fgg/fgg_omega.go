package fgg

import (
	"fmt"
	"reflect"
)

var _ = fmt.Errorf

/* Build Omega -- (morally) a map from ground FGG types to Sigs of (potential) calls on that receiver */

// Attempt to statically collect all
func GetOmega(ds []Decl, e_main FGGExpr) GroundMap {
	var gamma GroundEnv
	ground := make(GroundMap)
	collectGroundTypesFromExpr(ds, gamma, e_main, ground)
	fixOmega(ds, gamma, ground)
	return ground
}

/* GroundMap, GroundEnv, GroundTypeAndSigs, GroundSig */

// Maps u_ground.String() -> GroundTypeAndSigs{u_ground, sigs}
type GroundMap map[string]GroundTypeAndSigs

// Basically a Gamma for only TNamed
type GroundEnv map[Name]TNamed // Pre: forall TName, isGround

// A ground TNamed and the sigs of methods called on it as a receiver.
// sigs should include all potential such calls that may occur at run-time
type GroundTypeAndSigs struct {
	u_ground TNamed               // Pre: isGround(u_ground)
	sigs     map[string]GroundSig // Morally, Sig->[]Type -- HACK: string key is Sig.String
	// ^(i) FGG sigs; (ii) all sigs on u_ground receiver, including empty add-meth-targs
}

// The actual GroundTypeAndSigs.sigs map entry: Sig -> add-meth-targs
// i.e., the add-meth-targs that gives this Sig instance (param/return types).
// (Because Sig cannot be used as map key directly.)
type GroundSig struct {
	sig   Sig // CHECKME: may only need meth name (given receiver type), but Sig is convenient?
	targs []Type
}

/* fixOmega */

// Attempt to form a closure on encountered ground types.
// Iterate over `ground` using add-meth-targs recorded on i/face receivers to
// .. visit all possible method bodies of implementing struct types --
// .. repeating until no "new" ground types encoutered.
// Currently, very non-optimal.
// N.B. mutates `ground` -- encountered ground types collected into `ground`
func fixOmega(ds []Decl, gamma GroundEnv, ground GroundMap) {
	delta_empty := make(Delta)
	for again := true; again; {
		again = false

		for _, v_I := range ground {
			if !IsNamedIfaceType(ds, v_I.u_ground) || len(v_I.sigs) == 0 {
				continue
			}
			for _, v_S := range ground {
				if !IsStructType(ds, v_S.u_ground) ||
					!v_S.u_ground.Impls(ds, delta_empty, v_I.u_ground) {
					continue
				}

				u_S := v_S.u_ground
				for _, g_I := range v_I.sigs {
					if len(g_I.targs) == 0 { // CHECKME: dropping this skip obsoletes monom zigzag?
						continue
					}

					// Very non-optimal, may revisit the same g_I/u_S pair many times
					gamma1, e_body := getGroundEnvAndBody(ds, g_I, u_S)
					ground1 := make(map[string]GroundTypeAndSigs)
					collectGroundTypesFromExpr(ds, gamma1, e_body, ground1)
					for _, v_body := range ground1 {
						if _, ok := ground[v_body.u_ground.String()]; !ok {
							ground[v_body.u_ground.String()] = v_body
							again = true
						}
					}
				}
			}
		}
	}
}

// Get the Gamma and e_body for visiting the target meth of g_I on receiver u_S
func getGroundEnvAndBody(ds []Decl, g_I GroundSig, u_S TNamed) (
	GroundEnv, FGGExpr) {

	subs := make(map[TParam]Type)
	td_S := GetTDecl(ds, u_S.GetName())
	targs_recv := u_S.GetTArgs()
	tfs_recv := td_S.GetPsi().GetTFormals()
	for i := 0; i < len(targs_recv); i++ {
		subs[tfs_recv[i].name] = targs_recv[i]
	}
	tfs_meth := g_I.sig.GetPsi().GetTFormals()
	for i := 0; i < len(tfs_meth); i++ {
		subs[tfs_meth[i].name] = g_I.targs[i]
	}

	var pds []ParamDecl = nil
	for _, d := range ds {
		if md, ok := d.(MDecl); ok {
			if md.t_recv == u_S.t_name && md.name == g_I.sig.meth {
				pds = md.pDecls
				break
			}
		}
	}
	if pds == nil {
		panic("Method not found on " + u_S.String() + ": " + g_I.sig.meth)
	}

	x0, xs, e := body(ds, u_S, g_I.sig.meth, g_I.targs)
	gamma1 := make(GroundEnv)
	gamma1[x0] = u_S
	for i := 0; i < len(xs); i++ { // xs = ys in pds
		gamma1[xs[i]] = pds[i].GetType().TSubs(subs).(TNamed)
	}
	return gamma1, e
}

/* collectGroundTypesFromExpr, collectGroundTypesFromType, collectGroundTypesFromSigAndBody */

// Collect ground types from an Expr according to the Expr kind.
// gamma is needed when visiting an `e` of a "standalone" MDecl (via collectGroundTypesFromType)
// CHECKME: Post: res already collected?
// N.B. mutates `ground`
func collectGroundTypesFromExpr(ds []Decl, gamma GroundEnv, e FGGExpr,
	ground GroundMap) (res Type) {

	switch e1 := e.(type) {
	case Variable:
		res = gamma[e1.name]
	case StructLit:
		collectGroundTypesFromType(ds, e1.u_S, ground)
		for _, elem := range e1.elems {
			collectGroundTypesFromExpr(ds, gamma, elem, ground) // Discard return
		}
		res = e1.u_S
	case Select:
		u_S := collectGroundTypesFromExpr(ds, gamma, e1.e_S, ground).(TNamed) // Field types already collected via the structlit?
		// !!! we don't just visit e1.e_S, we also visit the type of e_S
		collectGroundTypesFromType(ds, u_S, ground)
		for _, fd := range fields(ds, u_S) {
			if fd.field == e1.field {
				res = fd.u.(TNamed)
				break
			}
		}
	case Call:
		u_recv := collectGroundTypesFromExpr(ds, gamma, e1.e_recv, ground)
		collectGroundTypesFromType(ds, u_recv, ground)
		for _, t_arg := range e1.t_args {
			collectGroundTypesFromType(ds, t_arg, ground)
		}
		for _, e_arg := range e1.args {
			collectGroundTypesFromExpr(ds, gamma, e_arg, ground) // Discard return
		}
		collectGroundTypesFromSigAndBody(ds, u_recv, e1, ground)

		gamma1 := make(Gamma)
		for k, v := range gamma {
			gamma1[k] = v
		}
		// !!! CHECKME: "actual" vs. "declared" (currently, actual)
		// .. declared is "higher", most exhaustive? (but declared collected via FromType) -- or do both?
		//res = g.u // May be a TParam, e.g., `Cond(type a Any())(br Branches(a)) a` (map.fgg)
		res = e1.Typing(ds, make(Delta), gamma1, true) // CHECKME: typing vs. sig? -- CHECKME: currently this typing mixed with res
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

// Collect ground types from a "standalone" type according to struct/interface,
// .. if u itself is ground.
// N.B. mutates `ground`
func collectGroundTypesFromType(ds []Decl, u Type, ground GroundMap) {

	if _, ok := ground[u.String()]; ok {
		return
	}
	if cast, ok := u.(TNamed); !ok || !isGround(cast) {
		return
	}

	u1 := u.(TNamed)
	gs := make(map[string]GroundSig) // CHECKME: make GroundSigs type?
	ground[u1.String()] = GroundTypeAndSigs{u1, gs}

	if IsStructType(ds, u1) { // Struct case
		u_S := u1

		// Visit fields
		fds := fields(ds, u_S)
		for _, fd := range fds {
			u_f := fd.u.(TNamed)
			collectGroundTypesFromType(ds, u_f, ground)
		}

		// Visit meths
		gs := methods(ds, u_S)
		for _, g := range gs {
			collectGroudTypesInSig(ds, g, ground)

			// Visit body (if no add-meth-tparams)
			pds := g.GetParamDecls()
			if len(g.GetPsi().GetTFormals()) == 0 {
				x_recv, xs, e_body := body(ds, u_S, g.meth, []Type{})
				gamma := make(GroundEnv)
				gamma[x_recv] = u_S
				for i := 0; i < len(pds); i++ {
					gamma[xs[i]] = pds[i].GetType().(TNamed)
				}
				collectGroundTypesFromExpr(ds, gamma, e_body, ground)
			}
		}

		// CHECKME: check all super interfaces, and (recursively) visit all meths of sub-structs?
		// no: fixOmega does the zig-zagging

	} else { //if IsNamedIfaceType(ds, u1) { // Interface case -- cf. u1 is TNamed
		u_I := u1

		// Visit meths
		gs := methods(ds, u_I)
		for _, g := range gs {
			collectGroudTypesInSig(ds, g, ground)
		}

		// Visit embedded
		td_I := GetTDecl(ds, u_I.t_name).(ITypeLit)
		tfs_I := td_I.GetPsi().GetTFormals()
		subs := make(map[TParam]Type)
		for i := 0; i < len(u_I.u_args); i++ {
			subs[tfs_I[i].name] = u_I.u_args[i]
		}
		for _, s := range td_I.specs {
			if u, ok := s.(TNamed); ok {
				collectGroundTypesFromType(ds, u.TSubs(subs), ground)
			}
		}

		// Visit all meths of subtype structs? -- no: leave to fixOmega
	}
}

// Visit types in sig (for tparams, the upper bounds)
func collectGroudTypesInSig(ds []Decl, g Sig, ground GroundMap) {
	psi_meth := g.GetPsi()
	for _, v := range psi_meth.GetTFormals() {
		collectGroundTypesFromType(ds, v.GetUpperBound(), ground)
	}
	pds := g.GetParamDecls()
	for i := 0; i < len(pds); i++ {
		u_pd := pds[i].GetType()
		collectGroundTypesFromType(ds, u_pd, ground)
	}
	collectGroundTypesFromType(ds, g.u_ret, ground)
}

// Record sig (i.e., add-meth-targs) for u_recv, and collect ground types from
// .. sig, if receiver and add-meth-targs all ground.
// Also visit call target e_body, if u_recv is a struct.
// Pre: if u0 is ground, then already in `ground` (cf. collectGroundTypesFromExpr, Call case).
// Can proceed without a Delta when u0 is ground Delta, as we also have add-targs here.
func collectGroundTypesFromSigAndBody(ds []Decl, u_recv Type, c Call,
	ground GroundMap) {

	// Receiver/add-meth-targs must be ground for the remainder
	if cast, ok := u_recv.(TNamed); !ok || !isGround(cast) {
		return
	}
	for _, t_arg := range c.t_args {
		if cast, ok := t_arg.(TNamed); !ok || !isGround(cast) {
			return
		}
	}

	g := methods(ds, u_recv)[c.meth]

	// If u_recv and add-meth-targs ground, add GroundSig for u_recv to GroundMap
	// .. if sig not already seen
	subs := make(map[TParam]Type)
	for i := 0; i < len(g.psi.tFormals); i++ {
		subs[g.psi.tFormals[i].name] = c.t_args[i]
	}
	g = g.TSubs(subs)
	gs := ground[u_recv.String()].sigs
	if _, ok := gs[g.String()]; ok {
		return
	}
	gs[g.String()] = GroundSig{g, c.t_args}

	// If sig not already seen (checked above), use sig to collect from
	// .. tparams upper bounds, params and return
	collectGroudTypesInSig(ds, g, ground)

	// If u_recv is ground struct and add-meth-targs ground, visit call target e_body
	if IsStructType(ds, u_recv) {
		u_S := u_recv.(TNamed)

		subs := make(map[TParam]Type)
		td_S := GetTDecl(ds, u_S.GetName())
		targs_S := u_S.GetTArgs()
		tfs_S := td_S.GetPsi().GetTFormals()
		for i := 0; i < len(targs_S); i++ {
			subs[tfs_S[i].name] = targs_S[i]
		}
		tfs_meth := g.GetPsi().GetTFormals()
		for i := 0; i < len(tfs_meth); i++ {
			subs[tfs_meth[i].name] = c.t_args[i]
		}

		pds := g.GetParamDecls()
		x0, xs, e_body := body(ds, u_S, c.meth, c.t_args)
		gamma1 := make(GroundEnv)
		gamma1[x0] = u_S
		for i := 0; i < len(xs); i++ { // xs = ys in pds
			gamma1[xs[i]] = pds[i].GetType().TSubs(subs).(TNamed) // Param names in g should be same as actual MDecl
		}
		collectGroundTypesFromExpr(ds, gamma1, e_body, ground)
	} else {
		// CHECKME: visit all possible bodies -- now subsumed by fixOmega?
	}
}
