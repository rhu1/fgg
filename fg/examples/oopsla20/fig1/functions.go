//$ go run github.com/rhu1/fgg -v -eval=-1 fg/examples/oopsla20/fig1/functions.go

package main;

import "fmt";

/* Library: Bool, Int */

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
func (this TT) Not() Bool { return FF{} };
func (this TT) Equal(that Any) Bool { return that.(Bool) };
func (this TT) Cond(br Branches) Any { return br.IfTT() };

type FF struct{};
func (this FF) Not() Bool { return TT{} };
func (this FF) Equal(that Any) Bool { return that.(Bool).Not() };
func (this FF) Cond(br Branches) Any { return br.IfFF() };

type Int interface {
	Inc() Int;
	Dec() Int;
	Eq(x Any) Bool;
	EqZero() Bool;
	EqNonZero(x Int) Bool;
	Add(x Int) Int;
	Gt(x Int) Bool;
	IsNeg() Bool
};

type Zero struct {};
func (x0 Zero) Inc() Int { return Pos{x0} };
func (x0 Zero) Dec() Int { return Neg{x0} };
func (x0 Zero) Eq(x Any) Bool { return x.(Int).EqZero() };
func (x0 Zero) EqZero() Bool { return TT{} };
func (x0 Zero) EqNonZero(x Int) Bool { return FF{} };
func (x0 Zero) Add(x Int) Int { return x };
func (x0 Zero) Gt(x Int) Bool { return x.IsNeg() };
func (x0 Zero) IsNeg() Bool { return FF{} };

type Pos struct { dec Int };
func (x0 Pos) Inc() Int { return Pos{x0} };
func (x0 Pos) Dec() Int { return x0.dec };
func (x0 Pos) Eq(x Any) Bool { return x0.EqNonZero(x.(Int)) };
func (x0 Pos) EqZero() Bool { return FF{} };
func (x0 Pos) EqNonZero(x Int) Bool { return x.Eq(x0.dec) };
func (x0 Pos) Add(x Int) Int { return x0.dec.Add(x.Inc()) };
func (x0 Pos) Gt(x Int) Bool { return x0.dec.Gt(x.Dec()) };
func (x0 Pos) IsNeg() Bool { return FF{} };

type Neg struct { inc Int };
func (x0 Neg) Inc() Int { return x0.inc };
func (x0 Neg) Dec() Int { return Neg{x0} };
func (x0 Neg) Eq(x Any) Bool { return x0.EqNonZero(x.(Int)) };
func (x0 Neg) EqZero() Bool { return FF{} };
func (x0 Neg) EqNonZero(x Int) Bool { return x.Eq(x0.inc) };
func (x0 Neg) Add(x Int) Int { return x0.inc.Add(x.Dec()) };
func (x0 Neg) Gt(x Int) Bool { return x0.inc.Gt(x.Inc()) };
func (x0 Neg) IsNeg() Bool { return TT{} };

type Ints struct {};
func (d Ints) _1() Int { return Pos{Zero{}} };
func (d Ints) _2() Int { return Ints{}._1().Add(Ints{}._1()) };
func (d Ints) _3() Int { return Ints{}._2().Add(Ints{}._1()) };
//func (d Ints) _4() Int { return Ints{}._3().Add(Ints{}._1()) };
//func (d Ints) _5() Int { return Ints{}._4().Add(Ints{}._1()) };
func (d Ints) __1() Int { return Neg{Zero{}} };
func (d Ints) __2() Int { return Ints{}.__1().Add(Ints{}.__1()) };
func (d Ints) __3() Int { return Ints{}.__2().Add(Ints{}.__1()) };
func (d Ints) __4() Int { return Ints{}.__3().Add(Ints{}.__1()) };
func (d Ints) __5() Int { return Ints{}.__4().Add(Ints{}.__1()) };

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
type incr struct { n Int };
func (this incr) Apply(x Any) Any {
	//return x.(int) + n
	return x.(Int).Add(this.n)
};
type pos struct {};
func (this pos) Apply(x Any) Any {
	//return x.(int) > 0
	return x.(Int).Gt(Zero{})
};

type compose struct {
	f Function;
	g Function
};
func (this compose) Apply(x Any) Any {
	return this.g.Apply(this.f.Apply(x))
};


func main() {
	/*var f Function = compose{incr{-5}, pos{}}
	var b bool = f.Apply(3).(bool)*/
	_ = compose{incr{Ints{}.__5()} , pos{}}.Apply(Ints{}._3()).(Bool)
}
