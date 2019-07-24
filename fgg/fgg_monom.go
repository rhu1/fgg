package fgg

import (
	"reflect"
	"strings"

	"github.com/rhu1/fgg/fg"
)

/* Naive monomorph -- !!WIP!! */

//func Monomorph(p FGGProgram) fg.FGProgram { /* TODO */ }

type WMap map[WKey]WVal
type WEnv map[TParam]Type // Pre: Type is closed

type WKey struct {
	t    Name
	hash string // Hack
}

type WVal struct {
	u   TName                 // Pre: isClosed(u)
	id  fg.Type               // Monomorph identifier
	mds map[string]SigRetPair // HACK: string key is Sig.String()
}

func (wv WVal) GetTName() TName {
	return wv.u
}

func (wv WVal) GetMonomId() fg.Type {
	return wv.id
}

func (wv WVal) GetParameterisedSigs() []fg.Sig {
	var res []fg.Sig
	for _, v := range wv.mds {
		res = append(res, v.g)
	}
	return res
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

func toMonomId(u TName) fg.Type {
	res := u.String()
	res = strings.Replace(res, "(", "<", -1)
	res = strings.Replace(res, ")", ">", -1)
	res = strings.Replace(res, " ", "", -1)
	return fg.Type(res)
}

type SigRetPair struct {
	g fg.Sig
	u TName // The "actual" return type -- cf. "declared" return type, g.u
}

type ClosedEnv map[Name]TName // Pre: forall TName, isClosed

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
				mds := omega[toWKey(u0)].mds // Pre: MakeWMap above ensures u0 in omega
				if tmp, ok := mds[hash]; ok {
					res = tmp.u
				} else {
					_, todo1 := visitSig(ds, u0, g, e1.targs, omega)
					todo = append(todo, todo1...)
					var pds []fg.ParamDecl
					for _, v := range g.pds {
						pds = append(pds, fg.NewParamDecl(v.x, toMonomId(v.u.(TName))))
					}
					mds[hash] = SigRetPair{fg.NewSig(g.m, pds, toMonomId(g.u.(TName))),
						res}
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
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String())
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
	omega[wk] = WVal{u, toMonomId(u), make(map[string]SigRetPair)}
	return true
}

// N.B. mutates omega -- i.e., omega is populated with the results
func visitSig(ds []Decl, u0 TName, g Sig, targs []Type, omega WMap) (res TName,
	todo []TName) {
	u_ret := g.u.(TName) // Closed, since u0 closed and no meth-params
	wk1 := toWKey(u_ret)
	if _, ok := omega[wk1]; !ok {
		omega[wk1] = WVal{u_ret, toMonomId(u_ret), make(map[string]SigRetPair)}
		todo = append(todo, u_ret)
	}
	for _, v := range g.pds {
		u_p := v.u.(TName)
		wk2 := toWKey(u_p)
		if _, ok := omega[wk2]; !ok {
			omega[wk2] = WVal{u_p, toMonomId(u_p), make(map[string]SigRetPair)}
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

/* Helper */

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

/* IGNORE */

/*
	// TODO: factor out with Call case
	u_ret := v.u.(TName) // Closed, since u closed and no meth-params
	key1 := toWKey(u_ret)
	if _, ok := omega[key1]; !ok {
		omega[key1] = WVal{u_ret, toMonomId(u_ret), make(map[string]SigRetPair)}
		todo = append(todo, u_ret)
	}
	for i := 0; i < len(v.pds); i++ {
		u_p := v.pds[i].u.(TName)
		key2 := toWKey(u_p)
		if _, ok := omega[key2]; !ok {
			omega[key2] = WVal{u_p, toMonomId(u_p), make(map[string]SigRetPair)}
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
