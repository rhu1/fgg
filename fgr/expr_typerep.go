package fgr

import (
	"log"

	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
)

// TypeRep is runtime type representation t↓(e)
type TypeRep struct {
	t         fg.Type
	paramReps []*TypeRep
}

func (tr TypeRep) IsValue() bool  { return false }
func (tr TypeRep) String() string { return "_TypeRep_" } // TODO(nickng)

// Reify(e) recursively traverses the TypeRep e
// and reconstructs the corresponding FGG type.
func Reify(e *TypeRep) fgg.Type {
	switch len(e.paramReps) {
	case 0: // Rₜ
		return fgg.NewType(fg.Name(e.t))

	default:
		var params []fgg.Type
		for _, pr := range e.paramReps {
			params = append(params, Reify(pr))
		}
		return fgg.NewType(fg.Name(e.t), params...)
	}
}

// getTypeRep(x tₛ) is a convenient shortcut to
//
// func (x tₛ) getTypeRep() TypeRep
func getTypeRep(x fg.Type) TypeRep {
	return TypeRep{}
}

// ToTypeRep converts FGG type fggType to its type representation.
func ToTypeRep(tσ fgg.Type) *TypeRep {
	switch t := tσ.(type) {
	case fgg.TName: // Type
		return paramTypeRep(t)
	case fgg.TParam: // Type parameter
		return groundTypeRep(t)
	}

	log.Fatalf("unknown FGG type: %v (%T)", tσ, tσ)
	return nil
}

// paramTypeRep converts t to TypeRep if type has type parameter.
func paramTypeRep(t fgg.TName) *TypeRep {
	rep := &TypeRep{t: fg.Type(t.Name())}
	for _, tParam := range t.Params() {
		rep.paramReps = append(rep.paramReps, ToTypeRep(tParam))
	}
	return rep
}

// groundTypeRep converts t to TypeRep if type has no type parameter.
func groundTypeRep(t fgg.TParam) *TypeRep {
	return &TypeRep{t: fg.Type(t)}
}
