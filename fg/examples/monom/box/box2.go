package main;
type Any interface {};
type D struct {};
func (x0 D) badA(x1 A) Any { return D{}.badA(x1) };
type A struct {};
func main() { _ = D{}.badA(A{}) }