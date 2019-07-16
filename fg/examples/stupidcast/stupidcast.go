//$ go run github.com/rhu1/fgg -eval=-1 -v fg/examples/stupidcast/stupidcast.go
// Cf.
//$ go run github.com/rhu1/fgg/fg/examples/stupidcast

package main; type Any interface{}; type ToAny struct { any Any }; type A struct {}; type B struct{}; func main() { _ = ToAny{A{}}.any.(B) }
