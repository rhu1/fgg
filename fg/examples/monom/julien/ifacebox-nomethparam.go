package main;
type Any interface {};
type BoxE interface { Make() BoxE };
type ABoxE struct {};
func (a ABoxE) Make() BoxE { return ABoxE{} };
type E struct { val D };
type D struct { val E };
type Dummy struct {};
func (x Dummy) doSomething(y BoxE) BoxE { return y.Make() };
func (x Dummy) makeBoxE() BoxE { return ABoxE{} };
func main() { _ = Dummy{}.doSomething(Dummy{}.makeBoxE()) }