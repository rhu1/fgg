package fg

func fields(ds []Decl, t_S Type) []FieldDecl {
	for _, v := range ds {
		s, ok := v.(TStruct)
		if ok && s.t == t_S {
			return s.fds
		}
	}
	panic("Unknown type: " + t_S.String())
}

// Go has no overloading, meth names are a unique key
func methods(ds []Decl, t Type) map[Name]MDecl {
	res := make(map[Name]MDecl)
	if isStructType(ds, t) {
		for _, v := range ds {
			m, ok := v.(MDecl)
			if ok && m.t == t {
				res[m.m] = m
			}
		}
	} else if isInterfaceType(ds, t) {
		panic("[TODO] interface types: " + t.String())
	} else { // Perhaps redundant if all TDecl OK checked first
		panic("Unknown type: " + t.String())
	}
	return res
}

func body(ds []Decl, t_S Type, m Name) Sig {
	for _, v := range ds {
		md, ok := v.(MDecl)
		if ok && md.t == t_S && md.m == m {
			return md.ToSig()
		}
	}
	panic("Method not found: " + t_S.String() + "." + m)
}

func isStructType(ds []Decl, t Type) bool {
	for _, v := range ds {
		d, ok := v.(TStruct)
		if ok && d.t == t {
			return true
		}
	}
	return false
}

func isInterfaceType(ds []Decl, t Type) bool {
	return !isStructType(ds, t) // FIXME: could be neither
}
