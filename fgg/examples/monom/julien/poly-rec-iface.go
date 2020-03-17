// This is monomorphisable !
package main;

type Any(type ) interface {};

type Dummy(type ) struct {};

func (x Dummy(type )) toAny(type )(y Any()) Any() {
	return y
};

type I(type a Any()) interface { m(type b Any())() I(I(b))};

type A(type ) struct {};

func (x A(type )) m(type b Any())() I(I(b)) {
	return Dummy(){}.toAny()(A(){}).(I(a)).m(a)()
};

func main() { _ =  A(){}.m(A())() }