//$ go run github.com/rhu1/fgg -fgg -v -monomc=-- fgg/examples/monom/julien/mono-ko/nested-fix.go

package main;
type Any(type ) interface {};
type Int(type ) struct {};
type Box(type a Any()) struct { cell a };

type Arg(type a Any()) struct {};

func (x Arg(type a Any())) mkNesting(type )(y a) a {
		return Arg(Box(a)){}.mkNesting()(Box(a){y}).cell
	};

func main() { _ = Arg(Int()){}.mkNesting()(Int(){}) }
