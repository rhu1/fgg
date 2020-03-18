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
then t never appears in T (inter-procedurally)

*/