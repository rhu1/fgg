// This is not monomorphisable but not well-formed (recursive struct)
package main;

type Any(type ) interface {};

type A(type ) struct {};

type B(type a Any()) struct {val C(a)};

type C(type a Any()) struct {val B(B(a))}; // no monom-related prob with val B(a) (but ruled out by real Go)

func (x A(type )) m(type a Any())() B(a) {
	return A(){}.m(a)()
};

func main() { _ =  A(){}.m(A())() }