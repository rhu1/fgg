//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/test/fg/monom/julien/mono-ok/i-closure-bad.go fgg/examples/monom/julien/mono-ok/i-closure-bad.go
//$ go run github.com/rhu1/fgg -eval=-1 -v tmp/test/fg/monom/julien/mono-ok/i-closure-bad.go

// This is monomorphisable ! -- rename "-bad"
package main;

type Any(type ) interface {};

type Dummy(type ) struct {};

type Pair(type ) struct {fst Any(); snd Any()};

// we need IB <: IA here
func (x Dummy(type )) useInterfaceA(type )(y IA()) IA() {
	return y.m1(S())() // m1() returns IB
};

func (x Dummy(type )) useInterfaceB(type )(y IB()) Any() {
	return y.m2(S())()
};

type IA(type ) interface { m1(type a Any())() IB()};
type IB(type ) interface {
	m1(type a Any())() IB();   // IB.m1() never occurs => need I-closure()
	m2(type a Any())() IB() }; // IB <: IA

type S(type ) struct {};
func (x S(type )) m1(type a Any())() IB() {return S(){}};
func (x S(type )) m2(type a Any())() IB() {return S(){}};

func main() { _ =
	Pair(){
		Dummy(){}.useInterfaceA()(S(){})
		,
		Dummy(){}.useInterfaceB()(S(){})
		}
	}
