package fgg

import "fmt"

var _ = fmt.Errorf

/* bounds(delta, u), fields(u_S), methods(u), body(u_S, m), type(v_S) */

func bounds(delta TEnv, u Type) Type {
	if a, ok := u.(TParam); ok {
		if res, ok := delta[a]; ok {
			return res
		}
	}
	return u // CHECKME: includes when TParam 'a' not in delta, correct?
}

/*func SCheckErr(msg string) {
	if e != nil {
		panic(msg)
	}
}*/

// Pre: len(s.psi.as) == len (u_S.typs), where s is the STypeLit decl for u_S.t
func fields(ds []Decl, u_S TName) []FieldDecl {
	s, ok := getTDecl(ds, u_S.t).(STypeLit)
	if !ok {
		panic("Not a struct type: " + u_S.String())
	}
	subs := make(map[TParam]Type) // TODO FIXME: use s.psi to do subs
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
	if isStructType(ds, u) {
		for _, v := range ds {
			md, ok := v.(MDecl)
			if ok && isStructType(ds, TName{md.t_recv, []Type{}}) { // FIXME HACK: TName
				//sd := md.recv.u.(TName)
				u1 := u.(TName)
				if md.t_recv == u1.t {
					subs := make(map[TParam]Type)
					for i := 0; i < len(md.psi_recv.tfs); i++ {
						subs[md.psi_recv.tfs[i].a] = u1.us[i]
					}
					res[md.m] = md.ToSig().Subs(subs)
				}
			}
		}
	} else if isInterfaceType(ds, u) {
		/*td := getTDecl(ds, t).(ITypeLit)
		for _, s := range td.ss {
			for _, v := range s.GetSigs(ds) {
				res[v.m] = v
			}
		}*/
		panic("[TODO]: ")
	} else { // Perhaps redundant if all TDecl OK checked first
		panic("Unknown type: " + u.String())
	}
	return res
}

/*
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
*/

func getTDecl(ds []Decl, t Name) TDecl {
	for _, v := range ds {
		td, ok := v.(TDecl)
		if ok && td.GetName() == t {
			return td
		}
	}
	panic("Type not found: " + t)
}
