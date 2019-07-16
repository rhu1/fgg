//$ go run github.com/rhu1/fgg -v -eval=10 fg/examples/hello/example.go
// Cf.
//$ go run github.com/rhu1/fgg/fg/examples/hello

// N.B. FG (or at least this implementation) is intended to be white-space *insensitive*.
// (E.g., this program could be written all on one line.)
// So the ';' are mandatory -- cannot replace by newlines (as in actual Go).
// (Newlines and other whitespace may be freely added, though.)
package main;
type A struct{};
func (x0 A) m() A { return x0.m() };
func main() {
	_ = A{}.m()
}
