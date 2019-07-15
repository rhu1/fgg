//$ go run github.com/rhu1/fgg -eval=3 fg/examples/equal/equal.go

package main;

type I1 interface { Equal(that Any) Bool };
type I2 interface { Equal(n Any) Bool };

type T struct {};
func (t T) Equal(foo Any) Bool { return Bool{} };
type Any interface {};
type ToAny struct { any Any };
type Bool struct {};  // Just for the purposes of this example

func main() {
	_ = ToAny{T{}}.any.(I1).(I2)
}
