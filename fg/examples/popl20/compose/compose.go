//$ go run github.com/rhu1/fgg -v -eval=46 fg/examples/popl20/compose/compose.go
// Cf.
//$ go run github.com/rhu1/fgg/fg/examples/popl20/compose

package main;

/* Base decls: Any, Booleans, Nautrals, Functions, Lists */

type Any interface {};

/* Booleans */

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

/* Naturals */

type Nat interface {
	Add(n Nat) Nat;
	Equal(n Any) Bool;
	equalZero() Bool;
	equalSucc(m Nat) Bool
};
type Zero struct {};
type Succ struct {pred Nat};
func (m Zero) Add (n Nat) Nat { return n };
func (m Succ) Add (n Nat) Nat { return Succ{m.pred.Add(n)} };
func (m Zero) Equal(n Any) Bool { return n.(Nat).equalZero() };
func (m Succ) Equal(n Any) Bool { return n.(Nat).equalSucc(m.pred) };
func (n Zero) equalZero() Bool { return TT{} };
func (n Succ) equalZero() Bool { return FF{} };
func (n Zero) equalSucc(m Nat) Bool { return FF{} };
func (n Succ) equalSucc(m Nat) Bool { return m.Equal(n.pred) };

/* Functions */

type Func interface {
	Apply(x Any) Any
};
type incr struct {
	n Nat
};
func (this incr) Apply(x Any) Any { return x.(Nat).Add(this.n) };
type not struct {};
func (this not) Apply(x Any) Any { return x.(Bool).Not() };
type compose struct {
	f Func;
	g Func
};
func (this compose) Apply(x Any) Any { return this.g.Apply(this.f.Apply(x)) };

/* Lists */

type List interface {
	Map(f Func) List;
	Member(x Eq) Bool
};
type Nil struct {};
type Cons struct {
	head Any;
	tail List
};
func (xs Nil) Map(f Func) List { return Nil{} };
func (xs Cons) Map(f Func) List { return Cons{f.Apply(xs.head), xs.tail.Map(f)} };
type memberBr struct {
	xs Cons;
	x Eq
};
func (this memberBr) IfTT() Any { return TT{} };
func (this memberBr) IfFF() Any { return this.xs.tail.Member(this.x) };
func (xs Nil) Member(x Eq) Bool { return FF{} };
func (xs Cons) Member(x Eq) Bool { return x.Equal(xs.head).Cond(memberBr{xs,x}).(Bool) };

/* Example code */

type D struct {};
func (d D) _1() Nat { return Succ{Zero{}} };
func (d D) _2() Nat { return D{}._1().Add(D{}._1()) };
func (d D) _3() Nat { return D{}._2().Add(D{}._1()) };

func main() {
	// Submission version: compose{incr{1},incr{2}}.Apply(3).(Nat)
	//_ = compose{incr{Succ{Zero{}}},incr{Succ{Succ{Zero{}}}}}.Apply(Succ{Succ{Succ{Zero{}}}}).(Nat) // -eval=26
	_ = compose{incr{D{}._1()}, incr{D{}._2()}}.Apply(D{}._3()).(Nat) // -eval=46

	// Also: _ = incr{2}.Apply(3).(Nat)
}
