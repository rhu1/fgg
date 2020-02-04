package main;

type Any(type ) interface {};


type Clothing(type a Any()) interface {
	Wash(type )() a ;
	Wear(type b Any())(y b) a
};

type Shirt(type ) struct{};

func (x Shirt(type )) Wash(type )() Shirt() {return Shirt(){}};
func (x Shirt(type )) Wear(type b Any())(y b) Shirt() {return Shirt(){}};

type Tyre(type a Any()) interface {
	Inflate(type )() a ;
	Wear(type b Any())(x b) b
};

type Bridgestone(type ) struct{};

func (x Bridgestone(type )) Inflate(type )() Bridgestone() {return Bridgestone(){}};
func (x Bridgestone(type )) Wear(type b Any())(y b) Bridgestone() {return Bridgestone(){}};


type Human(type ) struct {};
type Road(type ) struct {};


type Pair(type a Any(), b Any() ) struct { 
	fst a;
	snd b
};


func main() { _ =

	Pair(Shirt(), Bridgestone()){Shirt(){}.Wear(Human())(Human(){}), Bridgestone(){}.Wear(Road())(Road(){})}

}


