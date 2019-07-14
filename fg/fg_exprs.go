/*
 * This file contains defs for "concrete" syntax w.r.t. exprs.
 * Base (abstract) types, interfaces, etc. are in fg.go.
 */

package fg

import "fmt"
import "strings"

type Variable struct {
	id Name
}

var _ Expr = Variable{}

func (v Variable) Subs(m map[Variable]Expr) Expr {
	res, ok := m[v]
	if !ok {
		panic("Unknown var: " + v.String())
	}
	return res
}

func (v Variable) Eval() Expr {
	panic(v.id)
}

func (v Variable) Typing(ds []Decl, gamma Env) Type {
	res, ok := gamma[v.id]
	if !ok {
		panic("Var not in env: " + v.String())
	}
	return res
}

func (v Variable) String() string {
	return v.id
}

type StructLit struct {
	t  Type
	es []Expr
}

var _ Expr = StructLit{}

func (s StructLit) Subs(m map[Variable]Expr) Expr {
	return s
}

func (s StructLit) Eval() Expr {
	return s
}

func (s StructLit) Typing(ds []Decl, gamma Env) Type {
	fs := fields(ds, s.t)
	if len(s.es) != len(fs) {
		tmp := ""
		if len(fs) > 0 {
			tmp = fs[0].String()
			for _, v := range fs[1:] {
				tmp = tmp + ", " + v.String()
			}
		}
		panic("Arity mismatch: found=" +
			strings.Join(strings.Split(fmt.Sprint(s.es), " "), ", ") +
			", expected=[" + tmp + "]" + "\n\t" + s.String())
	}
	for i := 0; i < len(s.es); i++ {
		t := s.es[i].Typing(ds, gamma)
		u := fs[i].t
		if !t.Impls(ds, u) {
			panic("Arg expr must impl field type: arg=" + t.String() + ", field=" +
				u.String() + "\n\t" + s.String())
		}
	}
	return s.t
}

func (s StructLit) String() string {
	var sb strings.Builder
	sb.WriteString(s.t.String())
	sb.WriteString("{")
	sb.WriteString(strings.Trim(strings.Join(strings.Split(fmt.Sprint(s.es), " "), ", "), "[]"))
	sb.WriteString("}")
	return sb.String()
}

/*
type Select struct {
	e Expr
	f Name
}

func (this Select) Eval() Expr {

}

type Call struct {
	recv Expr
	m    Name
	args []Expr
}

func (this Call) Eval() Expr {

}

type Assert struct {
	e Expr
	t Name
}

func (this assert) Eval() Expr {

}
*/
