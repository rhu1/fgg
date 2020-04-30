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
	u_recv string
	m      Name
}

var meths []RecvMethPair = make([]RecvMethPair, 0) // TODO refactor

// TODO: rename
func Foo(ds []Decl) {
	graph := make(map[RecvMethPair]map[RecvMethPair]([][]Type))
	bools := make(map[RecvMethPair]map[RecvMethPair]bool)
	recvargs := make(map[RecvMethPair]map[RecvMethPair]([][]Type))
	for _, v := range ds {
		switch d := v.(type) {
		case STypeLit:
		case ITypeLit:
		case MethDecl:
			delta := d.GetMDeclPsi().ToDelta()
			tfs := d.GetRecvPsi().GetTFormals()
			u_args := make([]Type, len(tfs))
			for i := 0; i < len(tfs); i++ {
				//u_args[i] = tfs[i].GetUpperBound()
				u_args[i] = tfs[i].GetTParam()
				delta[tfs[i].GetTParam()] = tfs[i].GetUpperBound()
			}
			u_recv := TNamed{d.t_recv, u_args}
			gamma := make(Gamma)
			gamma[d.x_recv] = u_recv
			for _, v := range d.GetParamDecls() {
				gamma[v.GetName()] = v.GetType()
			}
			//ctxt := RecvMethPair{u_recv.TSubs(delta).String(), d.name}
			ctxt := RecvMethPair{u_recv.t_name, d.name}
			meths = append(meths, ctxt)
			bar(ds, delta, gamma, ctxt, d.e_body, graph, bools, recvargs)
		default:
			panic("Unknown Decl kind: " + reflect.TypeOf(v).String() + "\n\t" +
				v.String())
		}
	}

	////war(bools)
	//fmt.Println("1111a: ", graph, "\n")
	//fmt.Println("1111b: ", recvargs)
	//fmt.Println("2222: ", bools)

	findCycles(bools)
	//fmt.Println("3333: ", cycles)

	for _, v := range cycles {
		for i := 0; i < len(v); i++ {
			var next RecvMethPair
			if i == len(v)-1 {
				next = v[0]
			} else {
				next = v[i+1]
			}
			tmp := graph[v[i]]
			if tmp != nil {
				tmp2 := tmp[next]
				if tmp2 != nil {
					for _, t_args := range tmp2 {
						for _, u := range t_args {
							if u1, ok := u.(TNamed); ok {
								for _, x := range u1.u_args {
									if isOrContainsTParam(x) { // CHECKME: basically the naive syntactic restriction, OK?
										panic("Not monomorphisable, potential polymorphic recursion: " +
											fmt.Sprintf("%v", v))
									}
								}
							}
						}
					}
				}
			}

			rtmp := recvargs[v[i]]
			if rtmp != nil {
				rtmp2 := rtmp[next]
				if rtmp2 != nil {
					for _, t_args := range rtmp2 {
						for _, u := range t_args {
							if u1, ok := u.(TNamed); ok {
								for _, x := range u1.u_args {
									if isOrContainsTParam(x) { // CHECKME: basically the naive syntactic restriction, OK?
										panic("Not monomorphisable, potential polymorphic recursion: " +
											fmt.Sprintf("%v", v))
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

// N.B. mutates graph
func bar(ds []Decl, delta Delta, gamma Gamma, ctxt RecvMethPair, e FGGExpr,
	graph map[RecvMethPair]map[RecvMethPair]([][]Type),
	bools map[RecvMethPair]map[RecvMethPair]bool,
	recvargs map[RecvMethPair]map[RecvMethPair]([][]Type)) {
	switch e1 := e.(type) {
	case Variable:
	case StructLit:
		for _, elem := range e1.elems {
			bar(ds, delta, gamma, ctxt, elem, graph, bools, recvargs)
		}
	case Select:
		bar(ds, delta, gamma, ctxt, e1.e_S, graph, bools, recvargs)
	case Call:
		bar(ds, delta, gamma, ctxt, e1.e_recv, graph, bools, recvargs)
		for _, arg := range e1.args {
			bar(ds, delta, gamma, ctxt, arg, graph, bools, recvargs)
		}
		//g := methods(u_recv)[e1.meth]  // Want u_recv from Typing...
		/*var psi Psi
		for _, v := range ds {
			if v1, ok := v.(MDecl); ok && v1.name == e1.meth {
				psi = v1.GetMDeclPsi()
				break
			}
		}
		delta1 := psi.ToDelta()
		for k, v := range delta {
			delta1[k] = v
		}*/
		delta1 := delta // TODO refactor
		u_recv := e1.e_recv.Typing(ds, delta1, gamma, true)

		/*if _, ok := u_recv.(TParam); ok { // E.g., compose, x.Equal()(xs.head), x is `a`
			u_recv = delta[u_recv.(TParam)]
		}*/

		tmp := graph[ctxt]
		btmp := bools[ctxt]
		rtmp := recvargs[ctxt]
		if tmp == nil {
			tmp = make(map[RecvMethPair]([][]Type))
			graph[ctxt] = tmp
			btmp = make(map[RecvMethPair]bool)
			bools[ctxt] = btmp
			rtmp = make(map[RecvMethPair]([][]Type))
			recvargs[ctxt] = rtmp
		}
		if isStructType(ds, u_recv) {
			//key := RecvMethPair{u_recv.TSubs(delta1).String(), e1.meth}
			key := RecvMethPair{u_recv.TSubs(delta).(TNamed).t_name, e1.meth}
			tmp2 := tmp[key]
			if tmp2 == nil {
				tmp2 = make([][]Type, 0)
			}
			tmp2 = append(tmp2, e1.t_args)
			tmp[key] = tmp2
			btmp[key] = true
			if y, ok := u_recv.(TNamed); ok { // CHECKME: how about TParam?
				rtmp2 := rtmp[key]
				if rtmp2 == nil {
					rtmp2 = make([][]Type, 0)
				}
				rtmp2 = append(rtmp2, y.u_args)
				rtmp[key] = rtmp2
			}
		} else {
			u_I := u_recv // Or type param
			for _, v := range ds {
				switch d := v.(type) {
				case STypeLit:
					tfs := d.GetBigPsi().GetTFormals()
					u_args := make([]Type, len(tfs))
					for i := 0; i < len(tfs); i++ {
						u_args[i] = tfs[i].GetUpperBound()
					}
					u_S := TNamed{d.t_name, u_args}
					if p, ok := u_I.(TParam); (ok && u_S.ImplsDelta(ds, delta1, delta1[p])) || // CHECKME: delta1[p] ?
						(!ok && u_S.ImplsDelta(ds, delta1, u_I)) {
						//key := RecvMethPair{u_S.TSubs(delta1).String(), e1.meth}
						key := RecvMethPair{u_S.t_name, e1.meth}
						// TODO factor out below with above
						tmp2 := tmp[key]
						if tmp2 == nil {
							tmp2 = make([][]Type, 0)
						}
						tmp2 = append(tmp2, e1.t_args)
						tmp[key] = tmp2
						btmp[key] = true
						if y, ok := u_recv.(TNamed); ok { // CHECKME: how about TParam?
							rtmp2 := rtmp[key]
							if rtmp2 == nil {
								rtmp2 = make([][]Type, 0)
							}
							rtmp2 = append(rtmp2, y.u_args)
							rtmp[key] = rtmp2
						}
					}
				case ITypeLit:
				case MethDecl:
				default:
					panic("Unknown Decl kind: " + reflect.TypeOf(e).String() + "\n\t" +
						e.String())
				}
			}
		}
	case Assert:
		bar(ds, delta, gamma, ctxt, e1.e_I, graph, bools, recvargs)
	default:
		panic("Unknown Expr kind: " + reflect.TypeOf(e).String() + "\n\t" +
			e.String())
	}
}

/* Aux */

// Currently redundant
// Mutates graph
func war(graph map[RecvMethPair]map[RecvMethPair]bool) {
	for k := 0; k < len(meths); k++ {
		for i := 0; i < len(meths); i++ {
			for j := 0; j < len(meths); j++ {
				tmp := graph[meths[i]]
				if tmp == nil {
					/*tmp = make(map[RecvMethPair]bool)
					graph[meths[i]] = tmp*/
					return
				}
				if !tmp[meths[j]] {
					tmp2 := graph[meths[i]]
					tmp3 := graph[meths[k]]
					if tmp2 != nil && tmp3 != nil {
						tmp[meths[j]] = tmp2[meths[k]] && tmp3[meths[j]]
					}
				}
			}
		}
	}
}

var cycles [][]RecvMethPair

func findCycles(bools map[RecvMethPair]map[RecvMethPair]bool) {
	for _, v := range meths {
		stack := []RecvMethPair{v}
		aux(bools, stack)
	}
}

// DFS
func aux(bools map[RecvMethPair]map[RecvMethPair]bool, stack []RecvMethPair) {
	tmp := bools[stack[len(stack)-1]]
	if tmp == nil {
		return
	}
	for i := 0; i < len(meths); i++ {
		m := meths[i]
		if tmp[m] {
			if stack[0] == m {
				//stack1 := append(stack, m)
				cycles = append(cycles, stack)
				return
			}
			for j := 1; j < len(stack); j++ {
				if stack[j] == m {
					return
				}
			}
			stack1 := append(stack, m)
			aux(bools, stack1)
		}
	}
}
