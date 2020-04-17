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

	//fmt.Println("vvvv:")

	collectGroundTypesFromExpr(ds, gamma, e_main, ground, true)

	//fmt.Println("wwww:")

	fixOmega(ds, gamma, ground)
	return ground
}

/* Omega, GroundTypeAndSigs, GroundSig, GroundEnv */

type Omega1 struct {
	us map[string]TNamed // Pre: all TNamed are isGround
	//ms map[string]GroundTypeAndSigs // Maps u_ground.String() -> GroundTypeAndSigs{u_ground, sigs}
	ms map[string]MethInstan
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

// Attempt to form a closure on encountered ground types.
// Iterate over `ground` using add-meth-targs recorded on i/face receivers to
// .. visit all possible method bodies of implementing struct types --
// .. repeating until no "new" ground types encountered.
// Currently, very non-optimal.
// N.B. mutates `omega` -- encountered ground types collected into `ground`
/*func fixOmega(ds []Decl, gamma GroundEnv, omega Omega) {
	delta_empty := make(Delta)
	for again := true; again; {
		again = false

		//fmt.Println("000: ", omega, "\n")

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
					g_Ikey := toGroundSigsKey(g_I)
					if _, ok := wv_lower.sigs[g_Ikey]; ok {
						continue
					}
					wv_lower.sigs[g_Ikey] = g_I

					// Very non-optimal, may revisit the same g_I/u_S pair many times
					if IsStructType(ds, u_S) {
						gamma1, e_body := getGroundEnvAndBody(ds, g_I, u_S)
						omega1 := make(Omega)
						collectGroundTypesFromExpr(ds, gamma1, e_body, omega1, true)
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
}*/

/* Expressions */

// gamma used to type Call receiver
func collectExpr(ds []Decl, gamma GroundEnv, e FGGExpr,
	omega Omega1) {

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
		u_recv := e1.Typing(ds, make(Delta), gamma1, false).(TNamed)
		omega.us[toKey_Wt(u_recv)] = u_recv
		m := MethInstan{u_recv, e1.meth, e1.GetTArgs()} // CHECKME: why add u_recv separately?
		omega.ms[toKey_Wm(m)] = m
	case Assert:
		collectExpr(ds, gamma, e1.e_I, omega)
		u := e1.u_cast.(TNamed)
		omega.us[toKey_Wt(u)] = u
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

/* Aux */

func auxG(ds []Decl, omega Omega1) {
	auxF(ds, omega)
	auxI(ds, omega)
	auxM(ds, omega)
	auxS(ds, make(Delta), omega)
	auxP(ds, omega)
}

func auxF(ds []Decl, omega Omega1) {
	tmp := make(map[string]TNamed)
	for _, u := range omega.us {
		if !isStructType(ds, u) {
			continue
		}
		for _, u_f := range Fields(ds, u) {
			cast := u_f.u.(TNamed)
			tmp[toKey_Wt(cast)] = cast
		}
	}
	for k, v := range tmp {
		omega.us[k] = v
	}
}

func auxI(ds []Decl, omega Omega1) {
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
		omega.ms[k] = v
	}
}

func auxM(ds []Decl, omega Omega1) {
	tmp := make(map[string]TNamed)
	for _, m := range omega.ms {
		gs := methods(ds, m.u_recv)
		for _, g := range gs { // Should be singleton
			eta := MakeEta(g.psi, m.psi)
			for _, pd := range g.pDecls {
				u_pd := pd.u.SubsEta(eta)
				tmp[toKey_Wt(u_pd)] = u_pd
			}
			u_ret := g.u_ret.SubsEta(eta)
			tmp[toKey_Wt(u_ret)] = u_ret
		}
	}
	for k, v := range tmp {
		omega.us[k] = v
	}
}

func auxS(ds []Decl, delta Delta, omega Omega1) {
	tmp := make(map[string]MethInstan)
	for _, m := range omega.ms {
		for _, u := range omega.us {
			if !isStructType(ds, u) || !u.ImplsDelta(ds, delta, m.u_recv) {
				continue
			}
			x0, _, e := body(ds, u, m.meth, m.psi)
			gamma := make(GroundEnv)
			gamma[x0] = u
			// HERE gamma value param types
			collectExpr(ds, gamma, e, omega)
			m1 := MethInstan{u, m.meth, m.psi}
			tmp[toKey_Wm(m1)] = m1
		}
	}
	for k, v := range tmp {
		omega.ms[k] = v
	}
}

func auxP(ds []Decl, omega Omega1) {
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
		}
	}
	for k, v := range tmp {
		omega.ms[k] = v
	}
}

/*

















































 */

/* Omega, GroundTypeAndSigs, GroundSig, GroundEnv */

// Maps u_ground.String() -> GroundTypeAndSigs{u_ground, sigs}
type Omega map[string]GroundTypeAndSigs

// Pre: isGround(u_ground)
func toWKey(u_ground TNamed) string {
	return u_ground.String()
}

// A ground TNamed and the sigs of methods called on it as a receiver.
// sigs should include all potential such calls that may occur at run-time
type GroundTypeAndSigs struct {
	u_ground TNamed               // Pre: isGround(u_ground)
	sigs     map[string]GroundSig // string key is toGroundSigsKey, i.e., GroundSig.sig.String()
	// Morally, sigs is a map: fgg.Sig -> []Type -- all sigs on u_ground receiver, including empty add-meth-targs
}

func toGroundSigsKey(g GroundSig) string {
	return g.sig.String()
}

// The actual GroundTypeAndSigs.sigs map entry: Sig -> add-meth-targs
// i.e., the add-meth-targs that gives this Sig instance (param/return types).
// (Because Sig cannot be used as map key directly.)
type GroundSig struct {
	sig   Sig // CHECKME: may only need meth name (given receiver type), but Sig is convenient?
	targs []Type
}

// Basically a Gamma for only ground TNamed -- cf. Eta (TParam, not Name)
type GroundEnv map[Name]TNamed // Pre: forall TName, isGround

/* fixOmega */

// Attempt to form a closure on encountered ground types.
// Iterate over `ground` using add-meth-targs recorded on i/face receivers to
// .. visit all possible method bodies of implementing struct types --
// .. repeating until no "new" ground types encountered.
// Currently, very non-optimal.
// N.B. mutates `omega` -- encountered ground types collected into `ground`
func fixOmega(ds []Decl, gamma GroundEnv, omega Omega) {
	delta_empty := make(Delta)
	for again := true; again; {
		again = false

		//fmt.Println("000: ", omega, "\n")

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
					g_Ikey := toGroundSigsKey(g_I)
					if _, ok := wv_lower.sigs[g_Ikey]; ok {
						continue
					}
					wv_lower.sigs[g_Ikey] = g_I

					// Very non-optimal, may revisit the same g_I/u_S pair many times
					if IsStructType(ds, u_S) {
						gamma1, e_body := getGroundEnvAndBody(ds, g_I, u_S)
						omega1 := make(Omega)
						collectGroundTypesFromExpr(ds, gamma1, e_body, omega1, true)
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
// N.B. mutates `omega`
func collectGroundTypesFromExpr(ds []Decl, gamma GroundEnv, e FGGExpr,
	omega Omega, rec bool) (res Type) {

	//fmt.Println("2222:", e, "\n")

	switch e1 := e.(type) {
	case Variable:
		res = gamma[e1.name]
	case StructLit:
		collectGroundTypesFromType(ds, e1.u_S, omega, rec)
		for _, elem := range e1.elems {
			collectGroundTypesFromExpr(ds, gamma, elem, omega, rec) // Discard return
		}
		res = e1.u_S
	case Select:
		u_S := collectGroundTypesFromExpr(ds, gamma, e1.e_S, omega, rec).(TNamed) // Field types already collected via the structlit?
		// !!! we don't just visit e1.e_S, we also visit the type of e_S
		collectGroundTypesFromType(ds, u_S, omega, rec)
		for _, fd := range fields(ds, u_S) {
			if fd.field == e1.field {
				res = fd.u.(TNamed)
				break
			}
		}
	case Call:

		//fmt.Println("^cccc:", e1)

		u_recv := collectGroundTypesFromExpr(ds, gamma, e1.e_recv, omega, rec)
		collectGroundTypesFromType(ds, u_recv, omega, rec)
		for _, t_arg := range e1.t_args {
			collectGroundTypesFromType(ds, t_arg, omega, rec)
		}
		for _, e_arg := range e1.args {
			collectGroundTypesFromExpr(ds, gamma, e_arg, omega, rec) // Discard return
		}
		collectGroundTypesFromSigAndBody(ds, u_recv, e1, omega, rec)

		gamma1 := make(Gamma)
		for k, v := range gamma {
			gamma1[k] = v
		}
		// !!! CHECKME: "actual" vs. "declared" (currently, actual)
		// .. declared is "higher", most exhaustive? (but declared collected via FromType) -- or do both?
		//res = g.u // May be a TParam, e.g., `Cond(type a Any())(br Branches(a)) a` (map.fgg)
		res = e1.Typing(ds, make(Delta), gamma1, rec) // CHECKME: typing vs. sig? -- CHECKME: currently this typing mixed with res
	case Assert:
		u := e1.u_cast.(TNamed) // CHECKME: guaranteed?
		collectGroundTypesFromType(ds, u, omega, rec)
		collectGroundTypesFromExpr(ds, gamma, e1.e_I, omega, rec)
		res = u
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
	return res
}

var x int = 0

// Collect ground types from a "standalone" type according to struct/interface,
// .. if u itself is ground.
// N.B. mutates `omega`
func collectGroundTypesFromType(ds []Decl, u Type, omega Omega,
	rec bool) { // HACK FIXME

	x++
	if x > 100 {
		//panic("foo")
	}

	if cast, ok := u.(TNamed); !ok || !isGround(cast) {
		return
	}
	u1 := u.(TNamed)
	if _, ok := omega[toWKey(u1)]; ok {
		return
	}

	groundsigs := make(map[string]GroundSig) // CHECKME: make GroundSigs type?
	omega[toWKey(u1)] = GroundTypeAndSigs{u1, groundsigs}

	//fmt.Println("3333:", u1, IsStructType(ds, u1))

	if !rec {
		//return
	}

	if IsStructType(ds, u1) { // Struct case
		u_S := u1

		// Visit fields
		fds := fields(ds, u_S)
		for _, fd := range fds {
			u_f := fd.u.(TNamed)
			collectGroundTypesFromType(ds, u_f, omega, false)
		}
		//fmt.Println("4444:")

		// Visit meths
		gs := methods(ds, u_S)
		for _, g := range gs {
			collectGroudTypesInSig(ds, g, omega, false)

			// Visit body (if no add-meth-tparams)
			pds := g.GetParamDecls()
			if len(g.GetPsi().GetTFormals()) == 0 {
				x_recv, xs, e_body := body(ds, u_S, g.meth, []Type{})
				gamma := make(GroundEnv)
				gamma[x_recv] = u_S
				for i := 0; i < len(pds); i++ {
					gamma[xs[i]] = pds[i].GetType().(TNamed)
				}
				collectGroundTypesFromExpr(ds, gamma, e_body, omega, false)
				var _ = e_body
			}
		}
		//fmt.Println("5555:")

		// CHECKME: check all super interfaces, and (recursively) visit all meths of sub-structs?
		// no: fixOmega does the zig-zagging

	} else { //if IsNamedIfaceType(ds, u1) { // Interface case -- cf. u1 is TNamed
		u_I := u1

		// Visit meths
		gs := methods(ds, u_I)
		for _, g := range gs {
			collectGroudTypesInSig(ds, g, omega, false)
			var _ = g
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
				collectGroundTypesFromType(ds, u.TSubs(subs), omega, false)
			}
		}

		// Visit all meths of subtype structs? -- no: leave to fixOmega
	}
}

// Visit types in sig (for tparams, the upper bounds)
func collectGroudTypesInSig(ds []Decl, g Sig, omega Omega, rec bool) {
	psi_meth := g.GetPsi()
	for _, v := range psi_meth.GetTFormals() {
		collectGroundTypesFromType(ds, v.GetUpperBound(), omega, rec)
	}
	pds := g.GetParamDecls()
	for i := 0; i < len(pds); i++ {
		u_pd := pds[i].GetType()
		collectGroundTypesFromType(ds, u_pd, omega, rec)
	}
	collectGroundTypesFromType(ds, g.u_ret, omega, rec)
}

// Record sig (i.e., add-meth-targs) for u_recv, and collect ground types from
// .. sig, if receiver and add-meth-targs all ground.
// Also visit call target e_body, if u_recv is a struct.
// Pre: if u0 is ground, then already in `ground` (cf. collectGroundTypesFromExpr, Call case).
// Can proceed without a Delta when u0 is ground Delta, as we also have add-targs here.
func collectGroundTypesFromSigAndBody(ds []Decl, u_recv Type, c Call,
	omega Omega, rec bool) {

	//fmt.Println("^dddd:", c, omega[toWKey(u_recv.(TNamed))].sigs)

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
	//fmt.Println("^eeee:", u_recv, g, ",,", subs)
	if _, ok := gs[g.String()]; ok {
		return
	}
	//fmt.Println("^ffff:")
	tmp := GroundSig{g, c.t_args}
	gs[toGroundSigsKey(tmp)] = tmp // Record sig for u_recv
	// N.B. recorded only for u_recv, and not, e.g., super interfaces -- cf. monomTDecl, ITypeLit case

	// If sig not already seen (checked above), use sig to collect from
	// .. tparams upper bounds, params and return
	collectGroudTypesInSig(ds, g, omega, rec)

	// If u_recv is ground struct and add-meth-targs ground, visit call target e_body
	if IsStructType(ds, u_recv) {
		u_S := u_recv.(TNamed)

		//fmt.Println("^aaaa:", u_S, c)

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

		//fmt.Println("^bbbb:", e_body)

		gamma1 := make(GroundEnv)
		gamma1[x0] = u_S
		for i := 0; i < len(xs); i++ { // xs = ys in pds
			gamma1[xs[i]] = pds[i].GetType().TSubs(subs).(TNamed) // Param names in g should be same as actual MDecl
		}
		collectGroundTypesFromExpr(ds, gamma1, e_body, omega, rec)
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
