//$ go run github.com/rhu1/fgg -eval=-1 -v fg/examples/oopsla20/fig2/equality.go

package main;

import "fmt";

/* Library: Bool, Int */

type Bool interface {
	Not() Bool;
	Equal(that Eq) Bool;
	Cond(br Branches) Any;
	And(x Bool) Bool
};
type Branches interface {
	IfTT() Any;
	IfFF() Any
};
type TT struct{};
func (this TT) Not() Bool { return FF{} };
func (this TT) Equal(that Eq) Bool { return that.(Bool) };
func (this TT) Cond(br Branches) Any { return br.IfTT() };
func (this TT) And(x Bool) Bool { return x };

type FF struct{};
func (this FF) Not() Bool { return TT{} };
func (this FF) Equal(that Eq) Bool { return that.(Bool).Not() };
func (this FF) Cond(br Branches) Any { return br.IfFF() };
func (this FF) And(x Bool) Bool { return this };

type Int interface {
	Inc() Int;
	Dec() Int;
	Add(x Int) Int;
	Gt(x Int) Bool;
	IsNeg() Bool;
	IsZero() Bool;
	Eq(x Int) Bool
};

type Zero struct {};
func (x0 Zero) Inc() Int { return Pos{x0} };
func (x0 Zero) Dec() Int { return Neg{x0} };
func (x0 Zero) Add(x Int) Int { return x };
func (x0 Zero) Gt(x Int) Bool { return x.IsNeg() };
func (x0 Zero) IsNeg() Bool { return FF{} };
func (x0 Zero) IsZero() Bool { return TT{} };
func (x0 Zero) Eq(x Int) Bool { return x.IsZero() };

type Pos struct { dec Int };
func (x0 Pos) Inc() Int { return Pos{x0} };
func (x0 Pos) Dec() Int { return x0.dec };
func (x0 Pos) Add(x Int) Int { return x0.dec.Add(x.Inc()) };
func (x0 Pos) Gt(x Int) Bool { return x0.dec.Gt(x.Dec()) };
func (x0 Pos) IsNeg() Bool { return FF{} };
func (x0 Pos) IsZero() Bool { return FF{} };
func (x0 Pos) Eq(x Int) Bool { return x0.dec.Eq(x) };

type Neg struct { inc Int };
func (x0 Neg) Inc() Int { return x0.inc };
func (x0 Neg) Dec() Int { return Neg{x0} };
func (x0 Neg) Add(x Int) Int { return x0.inc.Add(x.Dec()) };
func (x0 Neg) Gt(x Int) Bool { return x0.inc.Gt(x.Inc()) };
func (x0 Neg) IsNeg() Bool { return TT{} };
func (x0 Neg) IsZero() Bool { return FF{} };
func (x0 Neg) Eq(x Int) Bool { return x0.inc.Eq(x) };

type Ints struct {};
func (d Ints) _1() Int { return Pos{Zero{}} };
func (d Ints) _2() Int { return d._1().Add(d._1()) };
func (d Ints) _3() Int { return d._2().Add(d._1()) };
//func (d Ints) _4() Int { return d._3().Add(d._1()) };
//func (d Ints) _5() Int { return d._4().Add(d._1()) };
func (d Ints) __1() Int { return Neg{Zero{}} };
func (d Ints) __2() Int { return d.__1().Add(d.__1()) };
func (d Ints) __3() Int { return d.__2().Add(d.__1()) };
func (d Ints) __4() Int { return d.__3().Add(d.__1()) };
func (d Ints) __5() Int { return d.__4().Add(d.__1()) };


/* Example code */

type Any interface {};
type Eq interface {
	Equal(that Eq) Bool
};
func (this Zero) Equal(that Eq) Bool {
	return this.Eq(that.(Int))
};
func (this Pos) Equal(that Eq) Bool {
	return this.Eq(that.(Int))
};
func (this Neg) Equal(that Eq) Bool {
	return this.Eq(that.(Int))
};
type Pair struct {
	left Eq;
	right Eq
};
func (this Pair) Equal(that Eq) Bool {
	return this.left.Equal(that.(Pair).left).And(this.right.Equal(that.(Pair).right))
};
func main() {
	/*var i, j Int = 1, 2
	var p Pair = Pair{i, j}
	var _ bool = p.Equal(p) // true*/
	fmt.Printf("%#v", Pair{Ints{}._1().(Pos), Ints{}._2().(Pos)}.
			Equal(Pair{Ints{}._1().(Pos), Ints{}._2().(Pos)}))
}
