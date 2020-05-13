//$ go run github.com/rhu1/fgg -v -eval=-1 fg/examples/oopsla20/fig1/functions.go

package main;

import "fmt";

/* Library: Bool, Int */

type Bool interface {
	Not() Bool;
	Equal(that Eq) Bool;
	Cond(br Branches) Any
};
type Branches interface {
	IfTT() Any;
	IfFF() Any
};
type TT struct{};
func (this TT) Not() Bool { return FF{} };
func (this TT) Equal(that Eq) Bool { return that.(Bool) };
func (this TT) Cond(br Branches) Any { return br.IfTT() };

type FF struct{};
func (this FF) Not() Bool { return TT{} };
func (this FF) Equal(that Eq) Bool { return that.(Bool).Not() };
func (this FF) Cond(br Branches) Any { return br.IfFF() };

type Int interface {
	Inc() Int;
	Dec() Int;
	Equal(x Eq) Bool;
	EqualZero() Bool;
	EqualNonZero(x Int) Bool;
	Add(x Int) Int;
	Gt(x Int) Bool;
	IsNeg() Bool
};

type Zero struct {};
func (x0 Zero) Inc() Int { return Pos{x0} };
func (x0 Zero) Dec() Int { return Neg{x0} };
func (x0 Zero) Equal(x Eq) Bool { return x.(Int).EqualZero() };
func (x0 Zero) EqualZero() Bool { return TT{} };
func (x0 Zero) EqualNonZero(x Int) Bool { return FF{} };
func (x0 Zero) Add(x Int) Int { return x };
func (x0 Zero) Gt(x Int) Bool { return x.IsNeg() };
func (x0 Zero) IsNeg() Bool { return FF{} };

type Pos struct { dec Int };
func (x0 Pos) Inc() Int { return Pos{x0} };
func (x0 Pos) Dec() Int { return x0.dec };
func (x0 Pos) Equal(x Eq) Bool { return x0.EqualNonZero(x.(Int)) };
func (x0 Pos) EqualZero() Bool { return FF{} };
func (x0 Pos) EqualNonZero(x Int) Bool { return x.Equal(x0.dec) };
func (x0 Pos) Add(x Int) Int { return x0.dec.Add(x.Inc()) };
func (x0 Pos) Gt(x Int) Bool { return x0.dec.Gt(x.Dec()) };
func (x0 Pos) IsNeg() Bool { return FF{} };

type Neg struct { inc Int };
func (x0 Neg) Inc() Int { return x0.inc };
func (x0 Neg) Dec() Int { return Neg{x0} };
func (x0 Neg) Equal(x Eq) Bool { return x0.EqualNonZero(x.(Int)) };
func (x0 Neg) EqualZero() Bool { return FF{} };
func (x0 Neg) EqualNonZero(x Int) Bool { return x.Equal(x0.inc) };
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


/* Later example */

type Eq interface {
	Equal(that Eq) Bool
};


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
