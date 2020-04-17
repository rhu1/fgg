package fgg

import "fmt"
import "reflect"

var _ = fmt.Errorf

/* Export */

func Bounds(delta Delta, u Type) Type          { return bounds(delta, u) }
func Fields(ds []Decl, u_S TNamed) []FieldDecl { return fields(ds, u_S) }
func Methods(ds []Decl, u Type) map[Name]Sig   { return methods(ds, u) }
func GetTDecl(ds []Decl, t Name) TDecl         { return getTDecl(ds, t) }

/* bounds(delta, u), fields(u_S), methods(u), body(u_S, m) */

// CHECKME: does Bounds ever need to work recursively? i.e., can an upperbound be a TParam?
// CHECKME: return type TName?
func bounds(delta Delta, u Type) Type {
	if a, ok := u.(TParam); ok {
		if res, ok := delta[a]; ok {
			return res
		}
	}
	return u // CHECKME: submission version, includes when TParam 'a' not in delta, correct?
}

// Pre: len(s.psi.as) == len (u_S.typs), where s is the STypeLit decl for u_S.t
func fields(ds []Decl, u_S TNamed) []FieldDecl {
	s, ok := getTDecl(ds, u_S.t_name).(STypeLit)
	if !ok {
		panic("Not a struct type: " + u_S.String())
	}
	subs := make(map[TParam]Type)
	for i := 0; i < len(s.Psi.tFormals); i++ {
		subs[s.Psi.tFormals[i].name] = u_S.u_args[i]
	}
	fds := make([]FieldDecl, len(s.fDecls))
	for i := 0; i < len(s.fDecls); i++ {
		fds[i] = s.fDecls[i].Subs(subs)
	}
	return fds
}

// Go has no overloading, meth names are a unique key
func methods(ds []Decl, u Type) map[Name]Sig {
	res := make(map[Name]Sig)
	if IsStructType(ds, u) {
		for _, v := range ds {
			md, ok := v.(MDecl)
			if ok && isStructName(ds, md.t_recv) {
				//sd := md.recv.u.(TName)
				u_S := u.(TNamed)
				if md.t_recv == u_S.t_name {
					subs := make(map[TParam]Type)
					for i := 0; i < len(md.PsiRecv.tFormals); i++ {
						subs[md.PsiRecv.tFormals[i].name] = u_S.u_args[i]
					}
					/*for i := 0; i < len(md.psi.tfs); i++ { // CHECKME: because TParam.TSubs will panic o/w -- refactor?
						subs[md.psi.tfs[i].a] = md.psi.tfs[i].a
					}*/
					res[md.name] = md.ToSig().TSubs(subs)
				}
			}
		}
	} else if IsNamedIfaceType(ds, u) { // N.B. u is a TName, \tau_I (not a TParam)
		u_I := u.(TNamed)
		td := getTDecl(ds, u_I.t_name).(ITypeLit)
		subs := make(map[TParam]Type)
		for i := 0; i < len(td.psi.tFormals); i++ {
			subs[td.psi.tFormals[i].name] = u_I.u_args[i]
		}
		for _, s := range td.specs {
			/*for _, v := range s.GetSigs(ds) {
				res[v.m] = v
			}*/
			switch c := s.(type) {
			case Sig:
				res[c.meth] = c.TSubs(subs)
			case TNamed: // Embedded u_I
				for k, v := range methods(ds, c.TSubs(subs)) { // CHECKME: can this cycle indefinitely? (cf. submission version)
					res[k] = v
				}
			default:
				panic("Unknown Spec kind: " + reflect.TypeOf(s).String())
			}
		}
	} else { // Perhaps redundant if all TDecl OK checked first
		panic("Unknown type: " + u.String())
	}
	return res
}

// Pre: t_S is a struct type
// Submission version, m(~\rho) informal notation
func body(ds []Decl, u_S TNamed, m Name, targs []Type) (Name, []Name, FGGExpr) {
	for _, v := range ds {
		md, ok := v.(MDecl)
		if ok && md.t_recv == u_S.t_name && md.name == m {
			xs := make([]Name, len(md.pDecls))
			for i := 0; i < len(md.pDecls); i++ {
				xs[i] = md.pDecls[i].name
			}
			subs := make(map[TParam]Type)
			for i := 0; i < len(md.PsiRecv.tFormals); i++ {
				subs[md.PsiRecv.tFormals[i].name] = u_S.u_args[i]
			}
			for i := 0; i < len(md.PsiMeth.tFormals); i++ {
				subs[md.PsiMeth.tFormals[i].name] = targs[i]
			}
			return md.x_recv, xs, md.e_body.TSubs(subs)
		}
	}
	panic("Method not found: " + u_S.String() + "." + m)
}

/* Additional */

func getTDecl(ds []Decl, t Name) TDecl {
	for _, v := range ds {
		td, ok := v.(TDecl)
		if ok && td.GetName() == t {
			return td
		}
	}
	panic("Type not found: " + t)
}
