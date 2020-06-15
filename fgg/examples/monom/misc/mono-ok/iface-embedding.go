// go run github.com/rhu1/fgg -fgg -eval=-1 -monomc=tmp/test/fg/monom/misc/iface-embedding.go fgg/examples/monom/misc/iface-embedding.go
// go run github.com/rhu1/fgg -eval=-1 tmp/test/fg/monom/misc/iface-embedding.go

package main;

import "fmt";

type Any(type ) interface {};

type DummyFunc(type X Any(), Y Any()) interface { apply(type )(a X) Y };

type Func(type X Any(), Y Any()) interface { DummyFunc(X,Y) };

type Box(type X Any()) interface {
	Map(type Y Any())(f Func(X,Y)) Box(Y)
};

type ABox(type X Any()) struct{
	value X
};


func (a ABox(type X Any())) Map(type Y Any())(f Func(X,Y)) Box(Y) {
	return ABox(Y){f.apply()(a.value)}
};

type Dummy(type ) struct{};

type D(type ) struct {};
type E(type ) struct {};

type DtoE(type ) struct {};
func (x0 DtoE(type )) apply(type )(d D()) E() { return E(){} };

func (x Dummy(type )) takeBox(type )(b Box(D())) Any() {
	return b.Map(E())(DtoE(){})  // Map<E>     // m(type a tau) ---> t\dagger
};

func main() {
	//_ =
	fmt.Printf("%#v",
		Dummy(){}.takeBox()(ABox(D()){D(){}}) // ABox<D>
	)
}


