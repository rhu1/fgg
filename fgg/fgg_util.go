package fgg

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/parser/fgg"
)

var _ = fmt.Errorf

// Pre: len(elems) > 1
// Pre: elems[:len(elems)-1] -- type/meth decls; elems[len(elems)-1] -- "main" func body expression
func MakeFggProgram(elems ...string) string {
	if len(elems) == 0 {
		panic("Bad empty args: must supply at least body expression for \"main\"")
	}
	var b strings.Builder
	b.WriteString("package main;\n")
	for _, v := range elems[:len(elems)-1] {
		b.WriteString(v)
		b.WriteString(";\n")
	}
	b.WriteString("func main() { _ = " + elems[len(elems)-1] + " }")
	return b.String()
}

/* Naive monomorph -- !!WIP!! */

//func monomorph(p FGGProgram) fg.FGProgram { /* TODO */ }

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

// CHECKME: "whole-program" approach starting from the "main" Expr means monom
// may still work in the presence of irregular types and polymorphic
// recursion? -- the cost is, of course, to give up separate compilation
//
// N.B. mutates omega -- i.e., omega is populated with the results
func MakeWMap(ds []Decl, gamma Env, e Expr, omega WMap) TName {
	var todo []TName // Pre: all TName are closed
	isTodo := func(u Type) bool {
		// N.B. mutates omega -- adds WKey, WVal pair (if 'u' closed)
		u1, ok := u.(TName) // Redundant?
		if !ok || !isClosed(u1) {
			return false
		}
		key := toWKey(u1)
		if _, ok := omega[key]; ok {
			return false
		}
		omega[key] = WVal{u1, toMonomId(u1), make(map[string]SigRetPair)}
		return true
	} // Usage contract: if return true, then append 'u' to 'todo'

	var res TName
	switch e1 := e.(type) {
	case Variable:
		res = gamma[e1.id].(TName)
	case StructLit:
		if isTodo(e1.u) { // CHECKME: do isTodo recursively on e1.u.us?
			todo = append(todo, e1.u) // Cannot refactor inside isTodo
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
		if //isStructTName(ds, u0) && // TODO FIXME: for interface types, visit all possible methods?  Or visiting all struct types already enough?
		// TODO: factor out with below for todo
		isClosed(u0) && len(e1.targs) > 0 {
			isC := true
			for _, v := range e1.targs {
				if u, ok := v.(TName); !ok || !isClosed(u) { // CHECKME: do isTodo recursively on targs?
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
				key := g.String()
				mds := omega[toWKey(u0)].mds // Pre: MakeWMap above ensures u0 in omega
				if tmp, ok := mds[key]; ok {
					res = tmp.u
				} else {
					u_ret := g.u.(TName) // Closed, since u closed and no meth-params
					key1 := toWKey(u_ret)
					if _, ok := omega[key1]; !ok {
						omega[key1] = WVal{u_ret, toMonomId(u_ret), make(map[string]SigRetPair)}
						todo = append(todo, u_ret)
					}
					for _, v := range g.pds {
						u_p := v.u.(TName)
						key2 := toWKey(u_p)
						if _, ok := omega[key2]; !ok {
							omega[key2] = WVal{u_p, toMonomId(u_p), make(map[string]SigRetPair)}
							todo = append(todo, u_ret)
						}
					}
					if isStructTName(ds, u0) {
						x0, xs, e2 := body(ds, u0, e1.m, e1.targs)
						gamma1 := make(Env)
						gamma1[x0] = u0
						for i := 0; i < len(xs); i++ {
							gamma1[xs[i]] = g.pds[i].u
						}
						res = MakeWMap(ds, gamma1, e2, omega)
						var pds []fg.ParamDecl
						for _, v := range g.pds {
							pds = append(pds, fg.NewParamDecl(v.x, toMonomId(v.u.(TName))))
						}
						mds[key] = SigRetPair{fg.NewSig(g.m, pds, toMonomId(g.u.(TName))), res}
					} else {
						res = g.u.(TName)
					}
				}
			}
		}
	case Assert:
		if isTodo(e1.u) {
			todo = append(todo, e1.u.(TName))
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
		if isInterfaceTName(ds, u) {
			continue
		}
		gs := methods(ds, u)
		for _, v := range gs {
			// TODO: factor out with Call case
			if len(v.psi.tfs) > 0 {
				continue
			}
			u_ret := v.u.(TName) // Closed, since u closed and no meth-params
			key := toWKey(u_ret)
			if _, ok := omega[key]; !ok {
				omega[key] = WVal{u_ret, toMonomId(u_ret), make(map[string]SigRetPair)}
				todo = append(todo, u_ret)
			}
			x0, xs, e := body(ds, u, v.m, empty)
			gamma1 := make(Env)
			gamma1[x0] = u
			for i := 0; i < len(xs); i++ {
				u_p := v.pds[i].u.(TName)
				key := toWKey(u_p)
				if _, ok := omega[key]; !ok {
					omega[key] = WVal{u_p, toMonomId(u_p), make(map[string]SigRetPair)}
					todo = append(todo, u_ret)
				}
				gamma1[xs[i]] = u_p
			}
			MakeWMap(ds, gamma1, e, omega)
		}
	}

	return res
}

/* Monomorph helper */

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

/* For "strict" parsing, *lexer* errors */

type FGGBailLexer struct {
	*parser.FGGLexer
}

// FIXME: not working -- e.g., incr{1}, bad token
// Want to "override" *BaseLexer.Recover -- XXX that's not how Go works (because BaseLexer is a struct, not interface)
func (b *FGGBailLexer) Recover(re antlr.RecognitionException) {
	message := "lex error after token " + re.GetOffendingToken().GetText() +
		" at position " + strconv.Itoa(re.GetOffendingToken().GetStart())
	panic(message)
}
