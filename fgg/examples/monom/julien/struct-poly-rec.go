// This should be monomorphisable!
package main;

type Any(type ) interface {};

type B(type a Any()) struct {val C(a)};

type C(type a Any()) struct {};

type A(type ) struct {};

func (x B(type a Any())) m(type )() B(C(a)) {
	return B(C(a)){C(C(a)){}}  // N.B. no recursion
};


func main() { _ =  B(A()){C(A()){}}.m()()}
