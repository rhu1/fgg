//$ go run github.com/rhu1/fgg -v -eval=-1 fg/examples/oopsla20/fig1/functions.go

package main;

import "fmt";

/* Library: Bool, Nat */

type Eq interface {
	Equal(that Any) Bool
};
type Bool interface {
	Not() Bool;
	Equal(that Any) Bool;
	Cond(br Branches) Any
};
type Branches interface {
	IfTT() Any;
	IfFF() Any
};
type TT struct{};
type FF struct{};

func (this TT) Not() Bool { return FF{} };
func (this FF) Not() Bool { return TT{} };
func (this TT) Equal(that Any) Bool { return that.(Bool) };
func (this FF) Equal(that Any) Bool { return that.(Bool).Not() };
func (this TT) Cond(br Branches) Any { return br.IfTT() };
func (this FF) Cond(br Branches) Any { return br.IfFF() };

type Nat interface {
	Add(n Nat) Nat;
	Equal(n Any) Bool;
	equalZero() Bool;
	equalSucc(m Nat) Bool;
	Gt(n Nat) Bool
};

type Zero struct {};
func (m Zero) Add (n Nat) Nat { return n };
func (m Zero) Equal(n Any) Bool { return n.(Nat).equalZero() };
func (m Zero) equalZero() Bool { return TT{} };
func (m Zero) equalSucc(n Nat) Bool { return FF{} };
func (m Zero) Gt(b Nat) Bool { return FF{} };

type Succ struct {pred Nat};
func (m Succ) Add (n Nat) Nat { return Succ{m.pred.Add(n)} };
func (m Succ) Equal(n Any) Bool { return n.(Nat).equalSucc(m.pred) };
func (m Succ) equalZero() Bool { return FF{} };
func (m Succ) equalSucc(n Nat) Bool { return n.Equal(m.pred) };
func (m Succ) Gt(n Nat) Bool {
	return n.equalZero().Cond(SuccGtCond{m, n}).(Bool)
};

type SuccGtCond struct { m Succ; n Nat };
func (x0 SuccGtCond) IfTT() Any { return TT{} };
func (x0 SuccGtCond) IfFF() Any { return x0.m.pred.Gt(x0.n.(Succ).pred) };

type Const struct {};
func (d Const) _1() Nat { return Succ{Zero{}} };
func (d Const) _2() Nat { return Const{}._1().Add(Const{}._1()) };
func (d Const) _3() Nat { return Const{}._2().Add(Const{}._1()) };
func (d Const) _4() Nat { return Const{}._3().Add(Const{}._1()) };
func (d Const) _5() Nat { return Const{}._4().Add(Const{}._1()) };

type List interface {
	Map(f Function) List;
	Member(x Eq) Bool
};
type Nil struct {};
type Cons struct {
	head Any;
	tail List
};
func (xs Nil) Map(f Function) List { return Nil{} };
func (xs Cons) Map(f Function) List { return Cons{f.Apply(xs.head), xs.tail.Map(f)} };
type memberBr struct {
	xs Cons;
	x Eq
};
func (this memberBr) IfTT() Any { return TT{} };
func (this memberBr) IfFF() Any { return this.xs.tail.Member(this.x) };
func (xs Nil) Member(x Eq) Bool { return FF{} };
func (xs Cons) Member(x Eq) Bool { return x.Equal(xs.head).Cond(memberBr{xs,x}).(Bool) };

/* Example code */

type Any interface {};
type Function interface {
	Apply(x Any) Any
};
//type incr struct { n int };
type incr struct { n Nat };
func (this incr) Apply(x Any) Any {
	//return x.(int) + n
	return x.(Nat).Add(this.n)
};
type pos struct {};
func (this pos) Apply(x Any) Any {
	//return x.(int) > 0
	return x.(Nat).Gt(Zero{})
};

type compose struct {
	f Function;
	g Function
};
func (this compose) Apply(x Any) Any {
	return this.g.Apply(this.f.Apply(x))
};


func main() {
	/*var f Functiontion = compose{incr{-5}, pos{}}
	var b bool = f.Apply(3).(bool)*/
	_ = compose{incr{Const{}._5()} , pos{}}.Apply(Const{}._3()).(Bool)
}
