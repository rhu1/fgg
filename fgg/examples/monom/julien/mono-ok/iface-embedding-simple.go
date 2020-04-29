// go run github.com/rhu1/fgg -fgg -eval=-1 -monomc=tmp/test/fg/monom/julien/iface-embedding-simple.go fgg/examples/monom/julien/iface-embedding-simple.go
// go run github.com/rhu1/fgg -eval=-1 tmp/test/fg/monom/julien/iface-embedding-simple.go

package main;

type Any(type ) interface {};

type DummyFunc(type A Any(), B Any()) interface { apply(type )(a A) B };

type Func(type A Any(), B Any()) interface { DummyFunc(A,B) };



type Dummy(type ) struct{};

type D(type ) struct {};
type E(type ) struct {};

type DtoE(type ) struct {};
func (x0 DtoE(type )) apply(type )(d D()) E() { return E(){} };


func main() { _ =
	DtoE(){}.apply()(D(){})
}


