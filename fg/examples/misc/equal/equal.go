//$ go run github.com/rhu1/fgg -v -eval=3 fg/examples/misc/equal/equal.go
// Cf.
// Cf.
//$ go run github.com/rhu1/fgg/fg/examples/misc/equals

package main;

import "fmt";

type I1 interface { Equal(that Any) Bool };
type I2 interface { Equal(n Any) Bool };

type T struct {};
func (t T) Equal(foo Any) Bool { return Bool{} };
type Any interface {};
type ToAny struct { any Any };
type Bool struct {};  // Just for the purposes of this example

func main() {
	//_ = ToAny{T{}}.any.(I1).(I2)
	fmt.Printf("%#v", ToAny{T{}}.any.(I1).(I2))
}
