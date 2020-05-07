
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


type tSC(type a Any()) struct {};
func (x tSC(type a Any())) m(type b Any())(y tI(a)) A(){
	return y.m(Box(b))(y)
};


func (x Dummy) method() { tSA(A).m(tSA) }

func main() { _ =
	tSA(A()){}.m(A())(tSA(A()){})

}

/*
for i++
 for each Di in ov(D)
	  G^i(Di) and check whether you reached a fixpoint (terminate that Di)
	   or it breaks occurs check then abort
finish and all Di are stable, or you have aborted
*/