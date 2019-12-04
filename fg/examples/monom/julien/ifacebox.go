package main;
type Any interface {};
type FuncDE interface { apply(a D) E };
type BoxE interface {};
type BoxD interface { MapE(f FuncDE) BoxE };
type ABoxD struct { value D };
type ABoxE struct { value E };
func (a ABoxD) MapE(f FuncDE) BoxE { return ABoxE{f.apply(a.value)} };
type Dummy struct {};
type D struct {};
type E struct {};
type DtoE struct {};
func (x0 DtoE) apply(d D) E { return E{} };
func (x Dummy) takeBox(b BoxD) Any { return b.MapE(DtoE{}) };
func main() { _ = Dummy{}.takeBox(ABoxD{D{}}) }