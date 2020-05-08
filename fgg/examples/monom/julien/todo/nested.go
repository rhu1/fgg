//$ go run github.com/rhu1/fgg -fgg -monomc=-- -v fgg/examples/monom/julien/nested.go

package main;

type Any(type ) interface {};

type Int(type ) struct {};

type Box(type a Any()) struct { cell a};


type NestedCons(type a Any()) struct {
	val a;
	tail Box(a)
};



type Arg(type a Any()) struct {};


// Badly typed (correct) -- cf. nested-fix
func (x Arg(type a Any())) mkNesting(type )(y a) Box(a) {
	return NestedCons(a){
		y,
		//Arg(Box(y)){}.mkNesting()(Box(y){y})  // FIXME: ImplsDelta blows up
		 Arg(Box(a)){}.mkNesting()(Box(a){y})
		 }.tail
};

func main() { _ =  Arg(){}.mkNesting(Int())(Int(){}) }

