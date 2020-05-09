//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/test/fg/monom/julien/mono-ok/i-closure.go -v fgg/examples/monom/julien/mono-ok/i-closure.go
//$ go run github.com/rhu1/fgg -eval=-1 -v tmp/test/fg/monom/julien/mono-ok/i-closure.go

// This is monomorphisable !
package main;

import "fmt";

type Any(type ) interface {};

type Dummy(type ) struct {};

func (x Dummy(type )) useInterface(type )(y IA()) Any() {
	return y.m1(S())()
};

// // Adding this fixes the panic
// func (x Dummy(type )) useInterfaceB(type )(y IB()) Any() {
// 	return y.m1(S())()
// };

type IA(type ) interface { m1(type a Any())() S() };
type IB(type ) interface { m1(type a Any())() S();  m2(type a Any())() S()};

type S(type ) struct {};
func (x S(type )) m1(type a Any())() S() {return S(){}};
func (x S(type )) m2(type a Any())() S() {return S(){}};

func main() {
	//_ =
	fmt.Printf("%#v",
		Dummy(){}.useInterface()(S(){}).(IB())
	)
}
