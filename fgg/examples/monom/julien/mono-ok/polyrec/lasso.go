// Should monomorphise (m3() is not called, m2() is not recursive and m1() is not poly-rec)

package main;

type Any(type ) interface {};

type Box(type a Any()) struct {};

type A(type ) struct {};

type B(type ) struct {};

func (x A(type )) m1(type a Any())() A(){
	return A(){}.m1(a)()
};

func (x B(type )) m2(type a Any())() A(){
	return A(){}.m1(Box(a))()
};

func (x A(type )) m3(type a Any())() A(){
	return A(){}.m3(Box(a))()
};

func main() { _ =  B(){}.m2(A())() }

