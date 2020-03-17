// This is not monomorphisable 
package main;

type Any(type ) interface {};

type A(type ) struct {};

type B(type a Any()) struct {val C(a)};

type C(type a Any()) struct {val B(B(a))};

func (x A(type )) m(type a Any())() B(a) {
	return A(){}.m(a)()
};


func main() { _ =  A(){}.m(A())() }