package fgg

import "fmt"
import "reflect"

var _ = fmt.Errorf

/* Export */

func Bounds(delta Delta, u Type) Type          { return bounds(delta, u) }
func Fields(ds []Decl, u_S TNamed) []FieldDecl { return fields(ds, u_S) }
func Methods(ds []Decl, u Type) map[Name]Sig   { return methods(ds, u) }
func GetTDecl(ds []Decl, t Name) TypeDecl      { return getTDecl(ds, t) }

/* bounds(delta, u), fields(u_S), methods(u), body(u_S, m) */

// return type TName?
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
	subs := make(map[TParam]Type) // Cf. MakeEta
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
	return methodsDelta(ds, make(Delta), u)
}

func methodsDelta(ds []Decl, delta Delta, u Type) map[Name]Sig {
	res := make(map[Name]Sig)
	if IsStructType(ds, u) {
		for _, v := range ds {
			md, ok := v.(MethDecl)
			if ok && isStructName(ds, md.t_recv) {
				//sd := md.recv.u.(TName)
				u_S := u.(TNamed)
				if md.t_recv == u_S.t_name {
					/*subs := make(map[TParam]Type)                    // Cf. MakeEta
					for i := 0; i < len(md.Psi_recv.tFormals); i++ { // TODO: md.Psi_recv.ToDelta
						subs[md.Psi_recv.tFormals[i].name] = u_S.u_args[i]
					}
					//for i := 0; i < len(md.psi.tfs); i++ { // CHECKME: because TParam.TSubs will panic o/w -- refactor?
					//	subs[md.psi.tfs[i].a] = md.psi.tfs[i].a
					//}
					res[md.name] = md.ToSig().TSubs(subs)*/
					if ok, eta := MakeEtaDelta(ds, delta, md.Psi_recv, u_S.u_args); ok {
						res[md.name] = md.ToSig().TSubs(eta)
					}
				}
			}
		}
	} else if IsNamedIfaceType(ds, u) { // N.B. u is a TName, \tau_I (not a TParam)
		u_I := u.(TNamed)
		td := getTDecl(ds, u_I.t_name).(ITypeLit)
		subs := make(map[TParam]Type) // Cf. MakeEta
		for i := 0; i < len(td.Psi.tFormals); i++ {
			subs[td.Psi.tFormals[i].name] = u_I.u_args[i]
		}
		for _, s := range td.specs {
			/*for _, v := range s.GetSigs(ds) {
				res[v.m] = v
			}*/
			switch s1 := s.(type) {
			case Sig:
				res[s1.meth] = s1.TSubs(subs)
			case TNamed: // Embedded u_I
				for k, v := range methods(ds, s1.TSubs(subs)) { // cycles? (cf. submission version)
					res[k] = v
				}
			default:
				panic("Unknown Spec kind: " + reflect.TypeOf(s).String())
			}
		}
	} else if cast, ok := u.(TParam); ok {
		upper, ok := delta[cast]
		if !ok {
			panic("Unknown type: " + u.String())
		}
		//return methodsDelta(ds, delta, bounds(delta, cast)) // !!! delegate to bounds
		return methodsDelta(ds, delta, upper)
	} else {
		panic("Unknown type: " + u.String()) // Perhaps redundant if all TDecl OK checked first
	}
	return res
}

// Pre: t_S is a struct type
// Submission version, m(~\rho) informal notation
//func body(ds []Decl, u_S TNamed, m Name, targs []Type) (Name, []Name, FGGExpr) {
func body(ds []Decl, u_S TNamed, m Name, targs []Type) (ParamDecl, []ParamDecl, FGGExpr) {
	for _, v := range ds {
		md, ok := v.(MethDecl)
		if ok && md.t_recv == u_S.t_name && md.name == m {
			subs := make(map[TParam]Type) // Cf. MakeEta
			for i := 0; i < len(md.Psi_recv.tFormals); i++ {
				subs[md.Psi_recv.tFormals[i].name] = u_S.u_args[i]
			}
			for i := 0; i < len(md.Psi_meth.tFormals); i++ {
				subs[md.Psi_meth.tFormals[i].name] = targs[i]
			}
			recv := ParamDecl{md.x_recv, u_S}
			pds := make([]ParamDecl, len(md.pDecls))
			for i := 0; i < len(md.pDecls); i++ {
				tmp := md.pDecls[i]
				pds[i] = ParamDecl{tmp.name, tmp.u.TSubs(subs)}
			}
			//return md.x_recv, xs, md.e_body.TSubs(subs)
			return recv, pds, md.e_body.TSubs(subs)
		}
	}
	panic("Method not found: " + u_S.String() + "." + m)
}

/* Additional */

func getTDecl(ds []Decl, t Name) TypeDecl {
	for _, v := range ds {
		td, ok := v.(TypeDecl)
		if ok && td.GetName() == t {
			return td
		}
	}
	panic("Type not found: " + t)
}
