package fgg

import "fmt"
import "reflect"

var _ = fmt.Errorf

/* bounds(delta, u), fields(u_S), methods(u), body(u_S, m), type(v_S) */

// CHECKME: does bounds ever need to work recursively? i.e., can an upperbound be a TParam?
// CHECKME: return type TName?
func bounds(delta TEnv, u Type) Type {
	if a, ok := u.(TParam); ok {
		if res, ok := delta[a]; ok {
			return res
		}
	}
	return u // CHECKME: submission version, includes when TParam 'a' not in delta, correct?
}

// Pre: len(s.psi.as) == len (u_S.typs), where s is the STypeLit decl for u_S.t
func fields(ds []Decl, u_S TName) []FieldDecl {
	s, ok := getTDecl(ds, u_S.t).(STypeLit)
	if !ok {
		panic("Not a struct type: " + u_S.String())
	}
	subs := make(map[TParam]Type)
	for i := 0; i < len(s.psi.tfs); i++ {
		subs[s.psi.tfs[i].a] = u_S.us[i]
	}
	fds := make([]FieldDecl, len(s.fds))
	for i := 0; i < len(s.fds); i++ {
		fds[i] = s.fds[i].Subs(subs)
	}
	return fds
}

// Go has no overloading, meth names are a unique key
func methods(ds []Decl, u Type) map[Name]Sig {
	res := make(map[Name]Sig)
	if isStructTName(ds, u) {
		for _, v := range ds {
			md, ok := v.(MDecl)
			if ok && isStructType(ds, md.t_recv) {
				//sd := md.recv.u.(TName)
				u_S := u.(TName)
				if md.t_recv == u_S.t {
					subs := make(map[TParam]Type)
					for i := 0; i < len(md.psi_recv.tfs); i++ {
						subs[md.psi_recv.tfs[i].a] = u_S.us[i]
					}
					/*for i := 0; i < len(md.psi.tfs); i++ { // CHECKME: because TParam.TSubs will panic o/w -- refactor?
						subs[md.psi.tfs[i].a] = md.psi.tfs[i].a
					}*/
					res[md.m] = md.ToSig().TSubs(subs)
				}
			}
		}
	} else if isInterfaceTName(ds, u) { // N.B. u is a TName, \tau_I (not a TParam)
		u_I := u.(TName)
		td := getTDecl(ds, u_I.t).(ITypeLit)
		subs := make(map[TParam]Type)
		for i := 0; i < len(td.psi.tfs); i++ {
			subs[td.psi.tfs[i].a] = u_I.us[i]
		}
		for _, s := range td.ss {
			/*for _, v := range s.GetSigs(ds) {
				res[v.m] = v
			}*/
			switch c := s.(type) {
			case Sig:
				res[c.m] = c.TSubs(subs)
			case TName: // Embedded u_I
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
func body(ds []Decl, u_S TName, m Name, targs []Type) (Name, []Name, Expr) {
	for _, v := range ds {
		md, ok := v.(MDecl)
		if ok && md.t_recv == u_S.t && md.m == m {
			xs := make([]Name, len(md.pds))
			for i := 0; i < len(md.pds); i++ {
				xs[i] = md.pds[i].x
			}
			subs := make(map[TParam]Type)
			for i := 0; i < len(md.psi_recv.tfs); i++ {
				subs[md.psi_recv.tfs[i].a] = u_S.us[i]
			}
			for i := 0; i < len(md.psi.tfs); i++ {
				subs[md.psi.tfs[i].a] = targs[i]
			}
			return md.x_recv, xs, md.e.TSubs(subs)
		}
	}
	panic("Method not found: " + u_S.String() + "." + m)
}

// Post: returns a struct type
func typ(ds []Decl, v StructLit) Type {
	u_S := v.u
	if !isStructTName(ds, u_S) {
		panic("Non struct type found in struct lit: " + v.String())
	}
	return u_S
}

func getTDecl(ds []Decl, t Name) TDecl {
	for _, v := range ds {
		td, ok := v.(TDecl)
		if ok && td.GetName() == t {
			return td
		}
	}
	panic("Type not found: " + t)
}
