// Will not monomorphise
package main;

type Any(type ) interface {};

type Box(type a Any()) interface { unbox(type )() a};

type SBox(type a Any()) struct {val a};

func (x SBox(type a Any())) unbox(type )() a {return x.val};

type A(type ) struct {};

func (x A(type )) m1(type a Any(), b Any())() A(){
	return A(){}.m1(Box(b), a)()
};


func main() { _ =  A(){}.m1(A(),A())() }

