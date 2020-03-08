package fgg

import (
	"fmt"
	"reflect"
	"strings"
)

var _ = fmt.Errorf
var _ = reflect.Append
var _ = strings.Compare

type RecvMethPair struct {
	//t_recv   Name   // Cf. MDecl
	//psi_recv string // HACK: string is psi.String()
	u_recv string
	m      Name
}

var ddd []RecvMethPair = make([]RecvMethPair, 0)

func Foo(ds []Decl) {
	graph := make(map[RecvMethPair]map[RecvMethPair]bool)
	for _, v := range ds {
		switch d := v.(type) {
		case STypeLit:
		case ITypeLit:
		case MDecl:
			delta := d.GetMDeclPsi().ToDelta()
			tfs := d.GetMDeclPsi().GetTFormals()
			u_args := make([]Type, len(tfs))
			for i := 0; i < len(tfs); i++ {
				u_args[i] = tfs[i].GetUpperBound()
			}
			u_recv := TNamed{d.t_recv, u_args}
			gamma := make(Gamma)
			gamma[d.x_recv] = u_recv
			for _, v := range d.GetParamDecls() {
				gamma[v.GetName()] = v.GetType()
			}
			ctxt := RecvMethPair{u_recv.String(), d.name}
			ddd = append(ddd, ctxt)
			bar(ds, delta, gamma, ctxt, d.e_body, graph)
		default:
			panic("Unknown Decl kind: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}

	war(graph)
	fmt.Println("1111: ", graph)
}

// N.B. mutates graph
func bar(ds []Decl, delta Delta, gamma Gamma, ctxt RecvMethPair, e FGGExpr,
	graph map[RecvMethPair]map[RecvMethPair]bool) {

	switch e1 := e.(type) {
	case Variable:
	case StructLit:
		for _, elem := range e1.elems {
			bar(ds, delta, gamma, ctxt, elem, graph)
		}
	case Select:
	case Call:
		bar(ds, delta, gamma, ctxt, e1.e_recv, graph)
		for _, arg := range e1.args {
			bar(ds, delta, gamma, ctxt, arg, graph)
		}
		//g := methods(u_recv)[e1.meth]  // Want u_recv from Typing...
		var psi Psi
		for _, v := range ds {
			if v1, ok := v.(MDecl); ok && v1.name == e1.meth {
				psi = v1.GetMDeclPsi()
				break
			}
		}
		delta1 := psi.ToDelta()
		for k, v := range delta {
			delta1[k] = v
		}
		u_recv := e1.e_recv.Typing(ds, delta1, gamma, true) // CHECKME: TParam possible? or already bounds
		tmp := graph[ctxt]
		if tmp == nil {
			tmp = make(map[RecvMethPair]bool)
			graph[ctxt] = tmp
		}
		if isStructType(ds, u_recv) {
			tmp[RecvMethPair{u_recv.String(), e1.meth}] = true
		} else {
			u_I := u_recv
			for _, v := range ds {
				switch d := v.(type) {
				case STypeLit:
					tfs := d.GetPsi().GetTFormals()
					u_args := make([]Type, len(tfs))
					for i := 0; i < len(tfs); i++ {
						u_args[i] = tfs[i].GetUpperBound()
					}
					u_S := TNamed{d.t_name, u_args}
					if u_S.ImplsDelta(ds, delta1, u_I) {
						tmp[RecvMethPair{u_S.String(), e1.meth}] = true
					}
				case ITypeLit:
				case MDecl:
				default:
					panic("Unknown Decl kind: " + reflect.TypeOf(e).String() + "\n\t" +
						e.String())
				}
			}
		}
	case Assert:
		bar(ds, delta, gamma, ctxt, e1.e_I, graph)
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

/* Aux */

// Mutates graph
func war(graph map[RecvMethPair]map[RecvMethPair]bool) {
	for k := 0; k < len(ddd); k++ {
		for i := 0; i < len(ddd); i++ {
			for j := 0; j < len(ddd); j++ {
				tmp := graph[ddd[i]]
				if tmp == nil {
					tmp = make(map[RecvMethPair]bool)
					graph[ddd[i]] = tmp
				}
				if !tmp[ddd[j]] {
					tmp2 := graph[ddd[i]]
					tmp3 := graph[ddd[k]]
					if tmp2 != nil && tmp3 != nil {
						tmp[ddd[j]] = tmp2[ddd[k]] && tmp3[ddd[j]]
					}
				}
			}
		}
	}
}
