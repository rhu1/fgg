//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/test/fg/monom/misc/mono-ok/poly-rec-iface.go fgg/examples/monom/misc/mono-ok/poly-rec-iface.go
//$ go run github.com/rhu1/fgg -eval=10 -v tmp/test/fg/monom/misc/mono-ok/poly-rec-iface.go

// This is monomorphisable !
package main;

import "fmt";

type Any(type ) interface {};

type Dummy(type ) struct {};

func (x Dummy(type )) toAny(type )(y Any()) Any() {
	return y
};

type I(type a Any()) interface { m(type b Any())() I(I(b))};

type A(type ) struct {};

func (x A(type )) m(type b Any())() I(I(b)) {
	//return Dummy(){}.toAny()(A(){}).(I(a)).m(a)()
	return Dummy(){}.toAny()(A(){}).(I(b)).m(b)()
};

func main() {
	//_ =
	fmt.Printf("%#v",
		A(){}.m(A())()
	)
}
