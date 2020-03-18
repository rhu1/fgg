// This is monomorphisable !
package main;

type Any(type ) interface {};

type Dummy(type ) struct {};

type Pair(type ) struct {fst Any(); snd Any()};

func (x Dummy(type )) useInterfaceA(type )(y IA()) IA() {
	return y.m1(S())()
};

func (x Dummy(type )) useInterfaceB(type )(y IB()) Any() {
	return y.m2(S())()
};

type IA(type ) interface { m1(type a Any())() IB()};
type IB(type ) interface { m1(type a Any())() IB();  m2(type a Any())() IB() }; // IB <: IA

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