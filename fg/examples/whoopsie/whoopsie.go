//$ go run github.com/rhu1/fgg fg/examples/whoopsie/whoopsie.go

// TODO FIXME: error not checked yet
package main; type Bad struct { whoopsie Bad }; type A struct { }; func main() { _ = A{} }
