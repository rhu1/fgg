//$ go run github.com/rhu1/fgg -fgg -eval=-1 -v fgg/examples/monom/julien/Parameterised-Map.go

package main;


type Any(type ) interface {};

type Int(type ) struct {};
type Bool(type ) struct {};


type Box(type a Any()) interface {
	get(type )() a
};

type ABox(type a Any()) struct{ val a };

func (x ABox(type a Any())) get(type )() a {return x.val};

type Func(type a Any(), b Any()) interface {
	apply(type )(x a) b
};


type ABox(type a Any()) struct{ val a };

type BadBox(type a Any()) struct { 
		val Box(Box(Box(a)))
		};

type Dummy(type ) struct{};

func (x Dummy(type )) toAny(type )(y Any()) Any() {
	return y.(Any())

};


func main() { _ =
	Dummy(){}.toAny()(ABox(Bool()){Bool(){}}).(BadBox(Int()))

}
