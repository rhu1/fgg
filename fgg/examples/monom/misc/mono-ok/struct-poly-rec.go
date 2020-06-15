//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/fg/monom/misc/mono-ok/struct-prob.go -v fgg/examples/monom/misc/mono-ok/struct-poly-rec.go
//
//$ go run github.com/rhu1/fgg -eval=-1 tmp/fg/monom/misc/mono-ok/struct-poly-rec.go

// This should be monomorphisable!
package main;

import "fmt";



type Any(type ) interface {};

type A(type ) struct {};

type B(type a Any()) struct {val C(a)};
//type B(type a Any()) struct {val C(C(a))};

type C(type a Any()) struct {};
//type C(type a Any()) struct {val B(a)};

func (x B(type a Any())) m(type )() B(C(a)) {  // Recurisve type nesting
	return B(C(a)){C(C(a)){}}  // N.B. but no actual recursion
};
/*
func (x A(type )) m(type a Any())() C(B(a)) {
	return A(){}.m(a)()
};
*/

/*func (x A(type )) m(type a Any())() C(B(a)) {
	return A(){}.m(a)()
};*/

func main() {
	//_ =
	fmt.Printf("%#v",
		B(A()){C(A()){}}.m()()
	)
}
//func main() { _ =  A(){}.m(A())() }


/*

for all:
type t(type ) struct T,
then t never appears in fields*(T) (transitively / inter-procedurally for struct)

*/















/*type Box(type a Any()) struct {val IBox(a)};
type IBox(type a Any()) interface { box() IBox(IBox(a)); unbox() a};

type Cons(a Any()) struct { head a; tail List(a)};
type Nil(a Any()) struct { };*/
