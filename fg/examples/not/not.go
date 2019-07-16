//$ go run github.com/rhu1/fgg -v -eval=4 fg/examples/not/not.go
// Cf.
//$ go run github.com/rhu1/fgg/fg/examples/not

package main;

/* Base decls: Any, Booleans, Functions */

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

/* Example code */

func main() {
	_ = not{}.Apply(TT{}).(Bool)
}
