package fgg

import (
	"fmt"
	"reflect"
)

var _ = fmt.Errorf

/**
 * Build Omega -- (morally) a map from ground FGG types to Sigs of (potential)
 * calls on that receiver.  N.B., calls are recorded only as seen for each
 * specific receiver type -- i.e., omega does not attempt to "respect"
 * subtyping (cf. "zigzagging" in fgg_monom).
 */

// Attempt to statically collect all ground types, and method instantiations
// called on those types, that may arise during execution
// Pre: isMonomorphisable -- TODO
func GetOmega(ds []Decl, e_main FGGExpr) Omega {
	var gamma GroundEnv
	ground := make(Omega)
	collectGroundTypesFromExpr(ds, gamma, e_main, ground)
	fixOmega(ds, gamma, ground)
	return ground
}

/* GroundMap, GroundEnv, GroundTypeAndSigs, GroundSig */

// Maps u_ground.String() -> GroundTypeAndSigs{u_ground, sigs}
type Omega map[string]GroundTypeAndSigs

// Pre: isGround(u_ground)
func toWKey(u_ground TNamed) string {
	return u_ground.String()
}

// Basically a Gamma for only TNamed
type GroundEnv map[Name]TNamed // Pre: forall TName, isGround

// A ground TNamed and the sigs of methods called on it as a receiver.
// sigs should include all potential such calls that may occur at run-time
type GroundTypeAndSigs struct {
	u_ground TNamed               // Pre: isGround(u_ground)
	sigs     map[string]GroundSig // string key is GroundSig.sig.String()
	// Morally, sigs is a map: fgg.Sig -> []Type -- all sigs on u_ground receiver, including empty add-meth-targs
}

func toGroundSigsKey(g Sig) string {
	return g.String()
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
// .. repeating until no "new" ground types encountered.
// Currently, very non-optimal.
// N.B. mutates `ground` -- encountered ground types collected into `ground`
func fixOmega(ds []Decl, gamma GroundEnv, omega Omega) {
	delta_empty := make(Delta)
	for again := true; again; {
		again = false

		for _, wv_upper := range omega {
			//fmt.Println("aaa: ", wv_upper)
			if !IsNamedIfaceType(ds, wv_upper.u_ground) || len(wv_upper.sigs) == 0 {
				continue
			}
			for _, wv_lower := range omega {

				//fmt.Println("bbb: ", wv_lower, wv_lower.u_ground.ImplsDelta(ds, delta_empty, wv_upper.u_ground))

				if //!IsStructType(ds, wv_S.u_ground) ||  // !!! Now include interfaces
				wv_lower.u_ground.Equals(wv_upper.u_ground) ||
					!wv_lower.u_ground.ImplsDelta(ds, delta_empty, wv_upper.u_ground) {
					continue
				}

				u_S := wv_lower.u_ground
				for _, g_I := range wv_upper.sigs {
					if len(g_I.targs) == 0 {
						continue
					}
					g_Ikey := toGroundSigsKey(g_I.sig)
					if _, ok := wv_lower.sigs[g_Ikey]; ok {
						continue
					}
					wv_lower.sigs[g_Ikey] = g_I

					// Very non-optimal, may revisit the same g_I/u_S pair many times
					gamma1, e_body := getGroundEnvAndBody(ds, g_I, u_S)
					omega1 := make(Omega)
					collectGroundTypesFromExpr(ds, gamma1, e_body, omega1)
					for _, wv_body := range omega1 {
						if _, ok := omega[toWKey(wv_body.u_ground)]; !ok {
							omega[toWKey(wv_body.u_ground)] = wv_body
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
	omega Omega) (res Type) {

	switch e1 := e.(type) {
	case Variable:
		res = gamma[e1.name]
	case StructLit:
		collectGroundTypesFromType(ds, e1.u_S, omega)
		for _, elem := range e1.elems {
			collectGroundTypesFromExpr(ds, gamma, elem, omega) // Discard return
		}
		res = e1.u_S
	case Select:
		u_S := collectGroundTypesFromExpr(ds, gamma, e1.e_S, omega).(TNamed) // Field types already collected via the structlit?
		// !!! we don't just visit e1.e_S, we also visit the type of e_S
		collectGroundTypesFromType(ds, u_S, omega)
		for _, fd := range fields(ds, u_S) {
			if fd.field == e1.field {
				res = fd.u.(TNamed)
				break
			}
		}
	case Call:
		u_recv := collectGroundTypesFromExpr(ds, gamma, e1.e_recv, omega)
		collectGroundTypesFromType(ds, u_recv, omega)
		for _, t_arg := range e1.t_args {
			collectGroundTypesFromType(ds, t_arg, omega)
		}
		for _, e_arg := range e1.args {
			collectGroundTypesFromExpr(ds, gamma, e_arg, omega) // Discard return
		}
		collectGroundTypesFromSigAndBody(ds, u_recv, e1, omega)

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
		collectGroundTypesFromType(ds, u, omega)
		collectGroundTypesFromExpr(ds, gamma, e1.e_I, omega)
		res = u
	case String: // CHECKME
		k := toWKey(STRING_TYPE)
		if _, ok := omega[k]; !ok {
			omega[k] = GroundTypeAndSigs{STRING_TYPE, make(map[string]GroundSig)}
		}
		res = STRING_TYPE
	case Sprintf:
		k := toWKey(STRING_TYPE)
		if _, ok := omega[k]; !ok {
			omega[k] = GroundTypeAndSigs{STRING_TYPE, make(map[string]GroundSig)}
		}
		for _, arg := range e1.args {
			collectGroundTypesFromExpr(ds, gamma, arg, omega) // Discard return
		}
		res = STRING_TYPE
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
	return res
}

// Collect ground types from a "standalone" type according to struct/interface,
// .. if u itself is ground.
// N.B. mutates `ground`
func collectGroundTypesFromType(ds []Decl, u Type, omega Omega) {

	if cast, ok := u.(TNamed); !ok || !isGround(cast) {
		return
	}
	u1 := u.(TNamed)
	if _, ok := omega[toWKey(u1)]; ok {
		return
	}

	groundsigs := make(map[string]GroundSig) // CHECKME: make GroundSigs type?
	omega[toWKey(u1)] = GroundTypeAndSigs{u1, groundsigs}

	if IsStructType(ds, u1) { // Struct case
		u_S := u1

		// Visit fields
		fds := fields(ds, u_S)
		for _, fd := range fds {
			u_f := fd.u.(TNamed)
			collectGroundTypesFromType(ds, u_f, omega)
		}

		// Visit meths
		gs := methods(ds, u_S)
		for _, g := range gs {
			collectGroudTypesInSig(ds, g, omega)

			// Visit body (if no add-meth-tparams)
			pds := g.GetParamDecls()
			if len(g.GetPsi().GetTFormals()) == 0 {
				x_recv, xs, e_body := body(ds, u_S, g.meth, []Type{})
				gamma := make(GroundEnv)
				gamma[x_recv] = u_S
				for i := 0; i < len(pds); i++ {
					gamma[xs[i]] = pds[i].GetType().(TNamed)
				}
				collectGroundTypesFromExpr(ds, gamma, e_body, omega)
			}
		}

		// CHECKME: check all super interfaces, and (recursively) visit all meths of sub-structs?
		// no: fixOmega does the zig-zagging

	} else { //if IsNamedIfaceType(ds, u1) { // Interface case -- cf. u1 is TNamed
		u_I := u1

		// Visit meths
		gs := methods(ds, u_I)
		for _, g := range gs {
			collectGroudTypesInSig(ds, g, omega)
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
				collectGroundTypesFromType(ds, u.TSubs(subs), omega)
			}
		}

		// Visit all meths of subtype structs? -- no: leave to fixOmega
	}
}

// Visit types in sig (for tparams, the upper bounds)
func collectGroudTypesInSig(ds []Decl, g Sig, omega Omega) {
	psi_meth := g.GetPsi()
	for _, v := range psi_meth.GetTFormals() {
		collectGroundTypesFromType(ds, v.GetUpperBound(), omega)
	}
	pds := g.GetParamDecls()
	for i := 0; i < len(pds); i++ {
		u_pd := pds[i].GetType()
		collectGroundTypesFromType(ds, u_pd, omega)
	}
	collectGroundTypesFromType(ds, g.u_ret, omega)
}

// Record sig (i.e., add-meth-targs) for u_recv, and collect ground types from
// .. sig, if receiver and add-meth-targs all ground.
// Also visit call target e_body, if u_recv is a struct.
// Pre: if u0 is ground, then already in `ground` (cf. collectGroundTypesFromExpr, Call case).
// Can proceed without a Delta when u0 is ground Delta, as we also have add-targs here.
func collectGroundTypesFromSigAndBody(ds []Decl, u_recv Type, c Call,
	omega Omega) {

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
	g = g.TSubs(subs) // CHECKME: keeping add-meth-params, used to create subs (e.g., getGroundEnvAndBody)
	//g = Sig{g.meth, Psi{[]TFormal{}}, g.pDecls, g.u_ret}
	gs := omega[toWKey(u_recv.(TNamed))].sigs
	if _, ok := gs[g.String()]; ok {
		return
	}
	gs[toGroundSigsKey(g)] = GroundSig{g, c.t_args} // Record sig for u_recv
	// N.B. recorded only for u_recv, and not, e.g., super interfaces -- cf. monomTDecl, ITypeLit case

	// If sig not already seen (checked above), use sig to collect from
	// .. tparams upper bounds, params and return
	collectGroudTypesInSig(ds, g, omega)

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
		collectGroundTypesFromExpr(ds, gamma1, e_body, omega)
	} else {
		// CHECKME: visit all possible bodies -- now subsumed by fixOmega?
	}
}

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
