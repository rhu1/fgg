//$ go run github.com/rhu1/fgg -fgg -monomc=tmp/test/fg/oopsla20/fig19/dispatcher.go fgg/examples/oopsla20/fig19/dispatcher.go
//$ go run github.com/rhu1/fgg -eval=-1 -v tmp/test/fg/oopsla20/fig19/dispatcher.go

package main;

import "fmt";

type Any(type ) interface {};

type Int(type ) struct {};

type Event(type ) interface{
	Process(type b Any())(y b) Int()
};

type UIEvent(type ) struct {};

func (x UIEvent(type )) Process(type b Any())(y b) Int() {return Int(){}};

type Dispatcher(type ) struct{};

func (x Dispatcher(type )) Dispatch(type )(y Event()) Int() {
	return y.Process(Int())(Int(){}) // IA : MyFunction<Int>
};


func main() {
	//_ =
	fmt.Printf("%#v",
		Dispatcher(){}.Dispatch()(UIEvent(){})
	)
}


