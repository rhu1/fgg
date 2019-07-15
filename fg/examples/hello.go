//rhu@HZHL4 MINGW64 ~/code/go/src/github.com/rhu1/fgg
//$ go run . -eval=10 example.go

// N.B. FG (or at least this implementation) is intended to be white-space *insensitive*.
// (E.g., this program could be written all on one line.)
// So the ';' are mandatory -- cannot replace by newlines (as in actual Go).
// (Newlines and whitespace may be freely added, though.)
package main;
type A struct{};
func (x0 A) m() A { return x0.m() };
func main() {
	_ = A{}.m()
}
