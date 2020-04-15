// This is not monomorphisable because it is not well-formed (recursive struct)
package main;

type Any(type ) interface {};

type A(type ) struct {};

type B(type a Any()) struct {val C(C(a))};

type C(type a Any()) struct {val B(a)}; 

func (x A(type )) m(type a Any())() C(B(a)) {
	return A(){}.m(a)()
};

func main() { _ =  A(){}.m(A())() }


/*

for all: 
type t(type ) struct T, 
then t never appears in fields*(T) (transitively / inter-procedurally for struct)

*/

type Box(type a Any()) struct {val IBox(a)};
type IBox(type a Any()) interface { box() IBox(IBox(a)); unbox() a};

type Cons(a Any()) struct { head a; tail List(a)};
type Nil(a Any()) struct { };