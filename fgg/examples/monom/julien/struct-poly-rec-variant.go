// This is not monomorphisable 
package main;

type Any(type ) interface {};

type A(type ) struct {};

type BS(type a Any()) struct {};

func (x BS(type a Any())) m(type )() BS(BS(a)) {
	return BS(BS(a)){}
};


func main() { _ =  BS(A()){}.m()()}