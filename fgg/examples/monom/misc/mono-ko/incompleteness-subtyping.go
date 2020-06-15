//$ go run github.com/rhu1/fgg -fgg -monomc=-- -v fgg/examples/monom/misc/mono-ko/incompleteness-subtyping.go

// This is not monomorphisable !
// panic: Not monomorphisable, potential polymorphic recursion: [{SA m1}]

package main;

type Any(type ) interface {};

type Dummy(type ) struct {};

type Box(type a Any()) struct {val a};

type Pair(type ) struct {fst Any(); snd Any()};

func (x Dummy(type )) useInterfaceA(type )(y IA()) S() {
	return y.m1(S())()
};

type IA(type ) interface { m1(type a Any())() S()};

type S(type ) struct {};

type SA(type ) struct {};
func (x SA(type )) m1(type a Any())() S() {return SA(){}.m1(Box(a))()};

type SB(type ) struct {};
func (x SB(type )) m1(type a Any())() S() {return S(){}};

func main() { _ =
		Dummy(){}.useInterfaceA()(SB(){})
	}
