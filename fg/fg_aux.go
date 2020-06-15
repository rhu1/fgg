package fg

import "fmt"

var _ = fmt.Errorf

/* fields(t_S), methods(t), body(t_S, m) */

// Pre: t_S is a struct type
func fields(ds []Decl, t_S Type) []FieldDecl {
	s, ok := getTDecl(ds, t_S).(STypeLit)
	if !ok {
		panic("Not a struct type: " + t_S.String())
	}
	return s.fDecls
}

// Go has no overloading, meth names are a unique key
func methods(ds []Decl, t Type) map[Name]Sig {
	res := make(map[Name]Sig)
	if isStructType(ds, t) {
		for _, v := range ds { // Factor out getMDecl?
			md, ok := v.(MethDecl)
			if ok && md.recv.t == t {
				res[md.name] = md.ToSig()
			}
		}
	} else if isInterfaceType(ds, t) {
		td := getTDecl(ds, t).(ITypeLit)
		for _, s := range td.specs {
			for _, v := range s.GetSigs(ds) { // cycles? (cf. submission version, recursive "methods")
				res[v.meth] = v
			}
		}
	} else { // Perhaps redundant if all TDecl OK checked first
		panic("Unknown type: " + t.String())
	}
	return res
}

// Pre: t_S is a struct type
func body(ds []Decl, t_S Type, m Name) (Name, []Name, FGExpr) {
	for _, v := range ds {
		md, ok := v.(MethDecl)
		if ok && md.recv.t == t_S && md.name == m {
			xs := make([]Name, len(md.pDecls))
			for i := 0; i < len(md.pDecls); i++ {
				xs[i] = md.pDecls[i].name
			}
			return md.recv.name, xs, md.e_body
		}
	}
	panic("Method not found: " + t_S.String() + "." + m)
}

/* Additional */

func getTDecl(ds []Decl, t Type) TDecl {
	for _, v := range ds {
		td, ok := v.(TDecl)
		if ok && td.GetType() == t {
			return td
		}
	}
	panic("Type not found: " + t)
}
