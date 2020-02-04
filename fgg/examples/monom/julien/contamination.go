package main;

type Any(type ) interface {};
type Bool(type ) struct {};

type Clothing(type a Any()) interface {
	Wash(type )() Bool() ;
	Wear(type b Any())(y b) Bool()
};

type Shirt(type ) struct{};

func (x Shirt(type )) Wash(type )() Bool() {return Bool(){}};
func (x Shirt(type )) Wear(type b Any())(y b) Bool() {return Bool(){}};

type Tyre(type a Any()) interface {
	Inflate(type )() Bool() ;
	Wear(type b Any())(x b) Bool()
};

type Bridgestone(type ) struct{};

func (x Bridgestone(type )) Inflate(type )() Bool() {return Bool(){}};
func (x Bridgestone(type )) Wear(type b Any())(y b) Bool() {return Bool(){}};


type Human(type ) struct {};
type Road(type ) struct {};


type Pair(type a Any(), b Any() ) struct { 
	fst a;
	snd b
};


type Dummy(type ) struct {};
func (x Dummy(type )) makePair(type )(c Clothing(Any()), t Tyre(Any())) Pair(Bool(), Bool()) {
	return 	Pair(Bool(), Bool()){c.Wear(Human())(Human(){}), t.Wear(Road())(Road(){})}
};


func main() { _ =

	Dummy(){}.makePair()(Shirt(){}, Bridgestone(){})
}


