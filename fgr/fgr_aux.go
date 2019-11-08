package fgr

import "fmt"

var _ = fmt.Errorf

/* fields(t_S), methods(t), body(t_S, m), type(v_S) */

// Pre: t_S is a struct type
func fields(ds []Decl, t_S Type) []FieldDecl {
	s, ok := getTDecl(ds, t_S).(STypeLit)
	if !ok {
		panic("Not a struct type: " + t_S.String())
	}
	return s.fds
}

// Go has no overloading, meth names are a unique key
func methods(ds []Decl, t Type) map[Name]Sig {
	res := make(map[Name]Sig)
	if isStructType(ds, t) {
		for _, v := range ds { // Factor out getMDecl?
			md, ok := v.(MDecl)
			if ok && md.recv.t == t {
				res[md.m] = md.ToSig()
			}
		}
	} else if isInterfaceType(ds, t) {
		td := getTDecl(ds, t).(ITypeLit)
		for _, s := range td.ss {
			for _, v := range s.GetSigs(ds) { // CHECKME: can this cycle indefinitely? (cf. submission version, recursive "methods")
				res[v.m] = v
			}
		}
	} else if t != TRep { // !!! Rep // Perhaps redundant if all TDecl OK checked first
		panic("Unknown type: " + t.String())
	}
	return res
}

// Pre: t_S is a struct type
func body(ds []Decl, t_S Type, m Name) (Name, []Name, Expr) {
	for _, v := range ds {
		md, ok := v.(MDecl)
		if ok && md.recv.t == t_S && md.m == m {
			xs := make([]Name, len(md.pds))
			for i := 0; i < len(md.pds); i++ {
				xs[i] = md.pds[i].x
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
