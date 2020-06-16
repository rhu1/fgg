package main;


type Any(type ) interface {};

type Int(type ) struct {};
type Bool(type ) struct {};

type Func(type a Any(), b Any()) interface {
	apply(type )(x a) b
};

type Bool2Int(type ) struct {};

func (x Bool2Int(type )) apply(type )(x Bool()) Int() {return Int(){} };

type Box(type a Any()) interface{
	map(type b Any())(f Func(a,b)) Box(b)

};

type IntBox(type ) struct {v Int()}; // IntBox <: 
func (x IntBox(type )) map(type b Any())(f Func(Int(),b)) Box(b) {return BoolBox{f.apply()(x.v)}}; 


type BoolBox(type ) struct {v Bool()}; // BoolBox <: IA(Bool())
func (x BoolBox(type )) map(type b Any())(f Func(Bool(),b)) Box(b) {return IntBox{f.apply()(x.v)}}; 



type Dummy(type ) struct{};

func (x Dummy(type )) CallFunctionBool(type )(y Box(Bool())) Bool() {
	return y.map(Int())(Bool2Int(){}) 

};


func main() { _ =
	Dummy(){}.CallFunctionBool()(BoolBox(){})  

}


