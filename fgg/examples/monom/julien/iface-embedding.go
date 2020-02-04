// go run github.com/rhu1/fgg -fgg -eval=-1 -monomc=tmp/test/fg/monom/julien/iface-embedding.go fgg/examples/monom/julien/iface-embedding.go
// go run github.com/rhu1/fgg -eval=-1 tmp/test/fg/monom/julien/iface-embedding.go

package main;

type Any(type ) interface {};

type DummyFunc(type A Any(), B Any()) interface { apply(type )(a A) B };

type Func(type A Any(), B Any()) interface { DummyFunc(A,B) };

type Box(type A Any()) interface {
	Map(type B Any())(f Func(A,B)) Box(B)
};

type ABox(type A Any()) struct{
	value A
};


func (a ABox(type A Any())) Map(type B Any())(f Func(A,B)) Box(B) {
	return ABox(B){f.apply()(a.value)}
};

type Dummy(type ) struct{};

type D(type ) struct {};
type E(type ) struct {};

type DtoE(type ) struct {};
func (x0 DtoE(type )) apply(type )(d D()) E() { return E(){} };

func (x Dummy(type )) takeBox(type )(b Box(D())) Any() {
	return b.Map(E())(DtoE(){})  // Map<E>     // m(type a tau) ---> t\dagger
};

func main() { _ =
	Dummy(){}.takeBox()(ABox(D()){D(){}}) // ABox<D>
}


