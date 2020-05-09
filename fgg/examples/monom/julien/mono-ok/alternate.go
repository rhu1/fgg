//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/test/fg/monom/julien/mono-ok/alternate.go fgg/examples/monom/julien/mono-ok/alternate.go
//$ go run github.com/rhu1/fgg -eval=10 -v tmp/test/fg/monom/julien/mono-ok/alternate.go

// Should monomorphise

package main;

import "fmt";

type Any(type ) interface {};

type Box(type a Any()) interface { unbox(type )() a};

type SBox(type a Any()) struct {val a};

func (x SBox(type a Any())) unbox(type )() a {return x.val};

type A(type ) struct {};

func (x A(type )) m1(type a Any())(y a) A(){
	return A(){}.m2(a, Box(a))(SBox(a){y})
};

func (x A(type )) m2(type a Any(), b Box(a))(y Box(a)) A(){
	return A(){}.m1(a)(y.unbox()())
};

func main() {
	//_ =
	fmt.Printf("%#v",
		A(){}.m1(A())(A(){})
	)
}

