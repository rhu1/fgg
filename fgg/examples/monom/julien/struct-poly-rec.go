//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/fg/monom/julien/struct-prob.go -v fgg/examples/monom/julien/struct-poly-rec.go

// This should be monomorphisable!
package main;

type Any(type ) interface {};

type B(type a Any()) struct {val C(a)};

type C(type a Any()) struct {};

type A(type ) struct {};

func (x B(type a Any())) m(type )() B(C(a)) {
	return B(C(a)){C(C(a)){}}  // N.B. no recursion
};
/*func (x A(type )) m(type a Any())() C(B(a)) {
	return A(){}.m(a)()
};*/

func main() { _ =  B(A()){C(A()){}}.m()()}
//func main() { _ =  A(){}.m(A())() }
