//$ go run github.com/rhu1/fgg -v -eval=7 fg/examples/booleans/booleans.go
// Cf.
//$ go run github.com/rhu1/fgg/fg/examples/booleans

package main;

/* Base decls: Any, Booleans */

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

/* Example code */

type exampleBr struct {
	x t;
	y t
};
func (this exampleBr) IfTT() Any {
	return this.x.m(this.y)
};
func (this exampleBr) IfFF() Any {
	return this.x
};

type t struct { };
func (x0 t) m(x1 t) t { return x1 };

type Ex struct {};
func (d Ex) example(b Bool, x t, y t) t {
	return b.Cond(exampleBr{x,y}).(t).m(t{})  // Main example, .(t)
};
func main() {
	_ = Ex{}.example(TT{}, t{}, t{})
}
