//$ go run github.com/rhu1/fgg -fgg -monomc=-- -v fgg/examples/monom/misc/mono-ko/monom-imp.go

package main;


type Any(type ) interface {};

type A(type ) struct {};

type Box(type a Any()) struct {};

type tI(type a Any()) interface {
	m(type b Any())(y tI(a)) A()
};

type tSA(type a Any()) struct {};

func (x tSA(type a Any())) m(type b Any())(y tI(a)) A(){
	return y.m(Box(b))(y)
};

type tSB(type a Any()) struct {};

func (x tSB(type a Any())) m(type b Any())(y tI(a)) A(){
	return A(){}
};

func main() { _ =
	tSA(A()){}.m(A())(tSA(A()){})

}
