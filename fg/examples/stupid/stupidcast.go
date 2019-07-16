//$ go run github.com/rhu1/fgg -eval=-1 -v fg/examples/stupid/stupid.go

package main; type Any interface{}; type ToAny struct { any Any }; type A struct {}; type B struct{}; func main() { _ = ToAny{A{}}.any.(B) }
