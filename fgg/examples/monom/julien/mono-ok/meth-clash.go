//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/test/fg/monom/julien/mono-ok/meth-clash.go fgg/examples/monom/julien/mono-ok/meth-clash.go
//$ go run github.com/rhu1/fgg -eval=7 -v tmp/test/fg/monom/julien/mono-ok/meth-clash.go

// Proposed (correct) monomed: https://play.golang.org/p/9nT1MahKCvE
package main;

type Any(type ) interface {};

type Top(type ) interface {};

type Dummy(type ) struct {};

func (x Dummy(type )) toAny(type )(y Any()) Any() {
	return y.(Any())
};

type Triple(type ) struct{ fst Any(); snd Any(); thd Any()};

type Foo(type ) interface { m(type a Top())() Dummy() };

// S </: Foo
type S(type ) struct {};
func (x S(type )) m(type a Any())() Dummy() {
	return Dummy(){}
};


// T <: Foo
type T(type ) struct {};
func (x T(type )) m(type a Top())() Dummy() {
	return Dummy(){}
};

func main() { _ =
	// Here both S.m() and Foo.m()  -which dont have matching signatures-
	// are instantiated to the same type (S).
	// which leads to mono'd FGG to accept a cast that fails in FGG -- edit: now fine
	Triple(){
		S(){}.m(S())(),
		//
		Dummy(){}.toAny()(T(){}).(Foo()).m(S())(),  // type-assert succeeds
		//
		Dummy(){}.toAny()(S(){}).(Foo()) // type-assert fails
	}
	}
