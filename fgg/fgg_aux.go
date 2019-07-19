package fgg

import "fmt"

var _ = fmt.Errorf

/* bounds(delta, u), fields(u_S), methods(u), body(u_S, m), type(v_S) */

func bounds(delta TEnv, u Type) Type {
	if a, ok := u.(TParam); ok {
		return delta[a]
	}
	return u
}

// Pre: u_S is a TName
func fields(ds []Decl, u_S Type) []FieldDecl {
	for _, v := range ds {
		s, ok := v.(STypeLit)
		if ok && s.t == u_S {
			subs := make(map[TParam]Type) // TODO FIXME: use s.psi to do subs
			fds := make([]FieldDecl, len(s.fds))
			for i := 0; i < len(s.fds); i++ {
				fds[i] = s.fds[i].Subs(subs)
			}
			return fds
		}
	}
	panic("Not a struct type: " + u_S.String())
}

/*
// Go has no overloading, meth names are a unique key
func methods(ds []Decl, t Type) map[Name]Sig {
	res := make(map[Name]Sig)
	if isStructType(ds, t) {
		for _, v := range ds {
			md, ok := v.(MDecl)
			if ok && md.recv.t == t {
				res[md.m] = md.ToSig()
			}
		}
	} else if isInterfaceType(ds, t) {
		td := getTDecl(ds, t).(ITypeLit)
		for _, s := range td.ss {
			for _, v := range s.GetSigs(ds) {
				res[v.m] = v
			}
		}
	} else { // Perhaps redundant if all TDecl OK checked first
		panic("Unknown type: " + t.String())
	}
	return res
}

// Pre: t_S is a struct type
func body(ds []Decl, t_S Type, m Name) (Name, []Name, Expr) {
	for _, v := range ds {
		md, ok := v.(MDecl)
		if ok && md.recv.t == t_S && md.m == m {
			xs := make([]Name, len(md.ps))
			for i := 0; i < len(md.ps); i++ {
				xs[i] = md.ps[i].x
			}
			return md.recv.x, xs, md.e
		}
	}
	panic("Method not found: " + t_S.String() + "." + m)
}

// Post: returns a struct type
func typ(ds []Decl, s StructLit) Type {
	t_S := s.t
	if !isStructType(ds, t_S) {
		panic("Non struct type found in struct lit: " + s.String())
	}
	return t_S
}

func getTDecl(ds []Decl, t Type) TDecl {
	for _, v := range ds {
		td, ok := v.(TDecl)
		if ok && td.GetType() == t {
			return td
		}
	}
	panic("Type not found: " + t)
}
*/
