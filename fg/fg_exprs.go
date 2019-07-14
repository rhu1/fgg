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

func (v Variable) CanEval(ds []Decl) bool {
	return false
}

func (v Variable) Eval(ds []Decl) Expr {
	panic("Cannot evaluate free variable: " + v.id)
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

func (s StructLit) CanEval(ds []Decl) bool {
	for _, v := range s.es {
		if v.CanEval(ds) {
			return true
		}
	}
	return false
}

func (s StructLit) Eval(ds []Decl) Expr {
	if !s.CanEval(ds) {
		panic("Stuck: " + s.String())
	}
	done := false
	es := make([]Expr, len(s.es))
	for i := 0; i < len(s.es); i++ {
		v := s.es[i]
		if !done && v.CanEval(ds) {
			v = v.Eval(ds)
			done = true
		}
		es[i] = v
	}
	return StructLit{s.t, es}
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

type Select struct {
	e Expr
	f Name
}

func (s Select) Subs(m map[Variable]Expr) Expr {
	return Select{s.e.Subs(m), s.f}
}

func (s Select) CanEval(ds []Decl) bool {
	if _, ok := s.e.(StructLit); !ok {
		return false
	}
	v := s.e.(StructLit)
	td := getTDecl(ds, v.t)
	if t_S, ok := td.(STypeLit); !ok { // Unnecessary?
		return false
	} else {
		for i := 0; i < len(t_S.fds); i++ {
			if t_S.fds[i].f == s.f {
				return true
			}
		}
		return false
	}
}

func (s Select) Eval(ds []Decl) Expr {
	if !s.CanEval(ds) {
		panic("Stuck: " + s.String())
	}
	v := s.e.(StructLit)
	td := getTDecl(ds, v.t).(STypeLit)
	for i := 0; i < len(td.fds); i++ {
		if td.fds[i].f == s.f {
			return v.es[i]
		}
	}
	panic("Field not found: " + s.f + " in\n\t" + td.String())
}

func (s Select) Typing(ds []Decl, gamma Env) Type {
	t := s.e.Typing(ds, gamma)
	if !isStructType(ds, t) {
		panic("Illegal select on non-struct type expr: " + t)
	}
	td := getTDecl(ds, t).(STypeLit)
	for _, v := range td.fds {
		if v.f == s.f {
			return v.t
		}
	}
	panic("Field not found: " + s.f + " in\n\t" + td.String())
}

func (s Select) String() string {
	return s.e.String() + "." + s.f
}

/*
type Call struct {
	recv Expr
	m    Name
	args []Expr
}

func (c Call) Eval() Expr {

}

type Assert struct {
	e Expr
	t Name
}

func (a Assert) Eval() Expr {

}
*/
