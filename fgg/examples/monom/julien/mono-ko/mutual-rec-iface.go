//$ go run github.com/rhu1/fgg -fgg -v -monom=-- fgg/examples/monom/julien/mutual-poly-rec.go

// This is not monomorphisable
// Not monomorphisable, potential polymorphic recursion: [{A m1} {B m2}]

package main;

type Any(type ) interface {};

type IBox(type a Any()) interface {};

/*
IBox<A>
IBox<IBox<A>>
IBox<IBox<IBox<IBox<A>>>>
*/

type A(type ) struct {};

type B(type ) struct {};

func (x A(type )) m1(type a Any())() A(){
	return B(){}.m2(IBox(a))()
};

func (x B(type )) m2(type a Any())() A(){
	return A(){}.m1(a)()
};

func main() { _ =  A(){}.m1(A())() }

