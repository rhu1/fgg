//$ go run github.com/rhu1/fgg -v -eval=13 fg/examples/popl20/map/map.go
// Cf.
//$ go run github.com/rhu1/fgg/fg/examples/popl20/map

package main;

import "fmt";

/* Base decls: Any, Booleans, Functions, Lists */

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

/* Functions */

type Func interface {
	Apply(x Any) Any
};
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

func main() {
	// Submission version was missing a "}"
	/*_ =  Cons{TT{}, Cons{FF{}, Nil{}}}.Map(not{}). // Main example
			(Cons).head.(Bool).Not() // Extra, assertion necessary*/
	fmt.Printf("%#v", Cons{TT{}, Cons{FF{}, Nil{}}}.Map(not{}). // Main example
			(Cons).head.(Bool).Not())
}
