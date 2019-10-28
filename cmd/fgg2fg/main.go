package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "not enough arguments (expected FGG file path)")
		os.Exit(1)
	}
	b, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	fggAdaptor := new(fgg.FGGAdaptor)
	fggProg := fggAdaptor.Parse(false, string(b))

	obliterate(fggProg.(fgg.FGGProgram))
}

func obliterate(prog fgg.FGGProgram) fg.FGProgram {
	return obliProg(prog)
}

// type GetRep interface{ getRep() Rep }
var _ = fg.NewITypeLit(fg.Type("GetRep"),
	[]fg.Spec{fg.NewSig("getRep", []fg.ParamDecl{}, fg.Type("Rep"))},
)

// obliProg implements ||P||
func obliProg(p fgg.FGGProgram) fg.FGProgram {
	ds, e := p.GetDecls(), p.GetExpr()

	var decls []fg.Decl
	for _, d := range ds {
		decls = append(decls, obliDecl(d))
	}

	return fg.NewFGProgram(decls, obliExpr(e))
}

// obliDecl implements ||D||
func obliDecl(d fgg.Decl) fg.Decl {
	switch d := d.(type) {
	case fgg.TDecl: // Type declaration

		switch d := d.(type) {
		case fgg.STypeLit: // type t(type α τ...) struct{fρ...}
			return obliStructDecl(d)
		case fgg.ITypeLit: // type t(type α τ...) interface {S...}
			return obliIfaceDecl(d)
		}

	case fgg.MDecl: // Method declaration
		return obliMDecl(d)
	}

	log.Fatalf("%v is not a valid declaration", d)
	return nil
}

// obliStructDecl translates struct type declaration.
func obliStructDecl(decl fgg.STypeLit) fg.Decl {

	ts := decl.GetName()
	ατs := decl.GetTFormals().Get()
	fρs := decl.GetFields()

	var fds []fg.FieldDecl
	for _, ατ := range ατs {
		fds = append(fds, fg.NewFieldDecl(ατ.Name(), fg.Type("Rep")))
	}
	for _, fρ := range fρs {
		fds = append(fds, fg.NewFieldDecl(fρ.GetName(), convertFGGTypeToFGType(fρ.GetType())))
	}

	return fg.NewSTypeLit(fg.Type(ts), fds)
}

// obliIfaceDecl translates interface type declaration.
func obliIfaceDecl(decl fgg.ITypeLit) fg.Decl {

	ti := decl.GetName()
	// TODO(nickng): ||S||
	var methods []fg.Spec

	return fg.NewITypeLit(fg.Type(ti), methods)
}

// obliMDecl translates method declaration.
func obliMDecl(decl fgg.MDecl) fg.Decl {

	xtατ := decl.Receiver()
	m := decl.GetName()
	βρs := decl.TFormals()
	yσs := decl.ParamDecls()
	e := decl.Body()

	xt := receiverToParamDecl(xtατ)
	var params []fg.ParamDecl
	for _, βρ := range βρs {
		params = append(params, typeParamToParamDecl(βρ))
	}
	for _, yσ := range yσs {
		params = append(params, paramToGetRep(yσ))
	}
	d := obliExpr(e)

	// func (x t) m(β Reps; y GetRep) GetRep { return d }
	return fg.NewMDecl(xt, m, params, fg.Type("GetRep"), d)
}

func convertFGGTypeToFGType(t fgg.Type) fg.Type {
	switch t := t.(type) {
	case fgg.TParam:
		return fg.Type(t)
	case fgg.TName:
		// TODO(nickng): what do we do with parameters?
		return fg.Type(t.Name())
	}

	log.Fatalf("%v is not a valid type", t)
	return fg.Type("")
}

// receiverToParamDecl converts fgg (x t(type α τ)) → fg (x t)
// removing the type parameter component.
func receiverToParamDecl(xtατ fgg.ParamDecl) fg.ParamDecl {
	switch t := xtατ.Type().(type) {
	case fgg.TParam: // Just name
		return fg.NewParamDecl(xtατ.Name(), fg.Type(t))
	case fgg.TName: // Ignore param part
		return fg.NewParamDecl(xtατ.Name(), fg.Type(t.Name()))
	}

	log.Fatalf("%v is not a valid receiver", xtατ)
	return fg.ParamDecl{}
}

// typeParamToParamDecl converts fgg (type β ρ) → fg (β Rep)
func typeParamToParamDecl(typeParam fgg.TFormal) fg.ParamDecl {
	return fg.NewParamDecl(typeParam.Name(), fg.Type("Rep"))
}

// paramToGetRep converts fgg (y σ) → fg (y GetRep)
func paramToGetRep(param fgg.ParamDecl) fg.ParamDecl {
	return fg.NewParamDecl(param.Name(), fg.Type("GetRep"))
}

func obliExpr(e base.Expr) fg.Expr {
	return nil // TODO
}
