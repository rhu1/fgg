//$ go run github.com/rhu1/fgg -eval=-1 -v fg/examples/oopsla20/fig2/lists.go

package main;

import "fmt";

/* Library: Bool, Int */

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
	Equal(x Any) Bool;
	EqualZero() Bool;
	EqualNonZero(x Int) Bool;
	Add(x Int) Int;
	Gt(x Int) Bool;
	IsNeg() Bool
};

type Zero struct {};
func (x0 Zero) Inc() Int { return Pos{x0} };
func (x0 Zero) Dec() Int { return Neg{x0} };
func (x0 Zero) Equal(x Any) Bool { return x.(Int).EqualZero() };
func (x0 Zero) EqualZero() Bool { return TT{} };
func (x0 Zero) EqualNonZero(x Int) Bool { return FF{} };
func (x0 Zero) Add(x Int) Int { return x };
func (x0 Zero) Gt(x Int) Bool { return x.IsNeg() };
func (x0 Zero) IsNeg() Bool { return FF{} };

type Pos struct { dec Int };
func (x0 Pos) Inc() Int { return Pos{x0} };
func (x0 Pos) Dec() Int { return x0.dec };
func (x0 Pos) Equal(x Any) Bool { return x0.EqualNonZero(x.(Int)) };
func (x0 Pos) EqualZero() Bool { return FF{} };
func (x0 Pos) EqualNonZero(x Int) Bool { return x.Equal(x0.dec) };
func (x0 Pos) Add(x Int) Int { return x0.dec.Add(x.Inc()) };
func (x0 Pos) Gt(x Int) Bool { return x0.dec.Gt(x.Dec()) };
func (x0 Pos) IsNeg() Bool { return FF{} };

type Neg struct { inc Int };
func (x0 Neg) Inc() Int { return x0.inc };
func (x0 Neg) Dec() Int { return Neg{x0} };
func (x0 Neg) Equal(x Any) Bool { return x0.EqualNonZero(x.(Int)) };
func (x0 Neg) EqualZero() Bool { return FF{} };
func (x0 Neg) EqualNonZero(x Int) Bool { return x.Equal(x0.inc) };
func (x0 Neg) Add(x Int) Int { return x0.inc.Add(x.Dec()) };
func (x0 Neg) Gt(x Int) Bool { return x0.inc.Gt(x.Inc()) };
func (x0 Neg) IsNeg() Bool { return TT{} };

type Ints struct {};
func (d Ints) _1() Int { return Pos{Zero{}} };
func (d Ints) _2() Int { return d._1().Add(d._1()) };
func (d Ints) _3() Int { return d._2().Add(d._1()) };
func (d Ints) _4() Int { return d._3().Add(d._1()) };
func (d Ints) _5() Int { return d._4().Add(d._1()) };
func (d Ints) _6() Int { return d._5().Add(d._1()) };
func (d Ints) __1() Int { return Neg{Zero{}} };
func (d Ints) __2() Int { return d.__1().Add(d.__1()) };
func (d Ints) __3() Int { return d.__2().Add(d.__1()) };
func (d Ints) __4() Int { return d.__3().Add(d.__1()) };
func (d Ints) __5() Int { return d.__4().Add(d.__1()) };


/* Prev. example */

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


/* Example code */

type Eq interface {
	Equal(that Eq) Bool
};
/*func (this int) Equal(that Eq) bool {  // "already" implemented
	return this == that.(int)
};
func (this bool) Equal(that Eq) bool {
	return this == that.(bool)
};*/

type List interface {
	Map(f Function) List
};
type Nil struct {};
type Cons struct {
	head Any;
	tail List
};
func (xs Nil) Map(f Function) List {
	return Nil{}
};
func (xs Cons) Map(f Function) List {
	return Cons{f.Apply(xs.head), xs.tail.Map(f)}
};


func main() {
	/*var xs List = Cons{3, Cons{6, Nil{}}}
	var ys List = xs.Map(incr{-5})
	var zs List = ys.Map(pos)*/
	_ = Cons{Ints{}._3(), Cons{Ints{}._6(), Nil{}}}.
				Map(incr{Ints{}.__5()}).
				Map(pos{})  // !!!
}
