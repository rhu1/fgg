// This is not monomorphisable ?
package main;

type Any(type ) interface {};


type I(type a Any()) interface { method(type )() I(I(a))};

type A(type ) struct {};

func (x A(type )) m(type a Any())() I(a) {
	return A(){}.m(a)()
};



func main() { _ =  A(){}.m(A())() }