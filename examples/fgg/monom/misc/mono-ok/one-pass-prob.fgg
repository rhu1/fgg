//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/test/fg/monom/misc/mono-ok/one-pass-prob.go fgg/examples/monom/misc/mono-ok/one-pass-prob.go
//$ go run github.com/rhu1/fgg -eval=-1 -v tmp/test/fg/monom/misc/mono-ok/one-pass-prob.go

package main;

import "fmt";

type Any(type ) interface {};

type Int(type ) struct {};

type Pair(type a Any(), b Any() ) struct {
	fst a;
	snd b
};

type IA(type ) interface{
	MyFunction(type b Any())(y b) Int()
};

type SA(type ) struct {};

// NB: SA.MyFunction() is only called via interface,
// can we find the instantiation of Pair(Int,Int) in one pass?
func (x SA(type )) MyFunction(type b Any())(y b) Int() {return Pair(b,Int()){y, Int(){}}.snd};


type Dummy(type ) struct{};

func (x Dummy(type )) CallFunction(type )(y IA()) Int() {
	return y.MyFunction(Int())(Int(){}) // IA : MyFunction<Int>
};


func main() {
	//_ =
	fmt.Printf("%#v",
		Dummy(){}.CallFunction()(SA(){})
	)
}


