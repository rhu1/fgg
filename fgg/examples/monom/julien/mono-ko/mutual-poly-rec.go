//$ go run github.com/rhu1/fgg -fgg -v -monomc=-- fgg/examples/monom/julien/mono-ko/mutual-poly-rec.go

package main;

type Any(type ) interface {};

type Box(type a Any()) struct {};

type A(type ) struct {};

type B(type ) struct {};

func (x A(type )) m1(type a Any())() A(){
	return B(){}.m2(Box(a))()
};

func (x B(type )) m2(type a Any())() A(){
	return A(){}.m1(a)()
};

func main() { _ =  A(){}.m1(A())() }

