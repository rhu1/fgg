//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/test/fg/monom/misc/mono-ok/Parameterised-Map.go -v fgg/examples/monom/misc/mono-ok/Parameterised-Map.go
//$ go run github.com/rhu1/fgg -eval=-1 -v tmp/test/fg/monom/misc/mono-ok/Parameterised-Map.go

package main;

import "fmt";

type Any(type ) interface {};

type Int(type ) struct {};
type Bool(type ) struct {};

type Func(type a Any(), b Any()) interface {
	apply(type )(x a) b
};

type Bool2Int(type ) struct {};
type ParamBox(type a Any()) struct {v1 a};
func (x ParamBox(type a Any())) map(type b Any())(f Func(a,b)) Box(b) {return ParamBox(b){f.apply()(x.v1)}};


func (x Bool2Int(type )) apply(type )(y Bool()) Int() {return Int(){} };

type Box(type a Any()) interface{
	map(type b Any())(f Func(a,b)) Box(b)

};

type IntBox(type ) struct {v1 Int()}; // IntBox <:
func (x IntBox(type )) map(type b Any())(f Func(Int(),b)) Box(b) {return ParamBox(b){f.apply()(x.v1)}};


type BoolBox(type ) struct {v1 Bool()}; // BoolBox <: IA(Bool())
func (x BoolBox(type )) map(type b Any())(f Func(Bool(),b)) Box(b) {return ParamBox(b){f.apply()(x.v1)}};


type Dummy(type ) struct{};

func (x Dummy(type )) CallFunctionBool(type )(y Box(Bool())) Box(Int()) {
	return y.map(Int())(Bool2Int(){})

};


func main() {
	//_ =
	fmt.Printf("%#v",
		Dummy(){}.CallFunctionBool()(BoolBox(){Bool(){}})
	)
}
