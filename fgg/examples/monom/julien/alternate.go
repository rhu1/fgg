package main;

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

func main() { _ =  A(){}.m1(A())(A(){}) }

