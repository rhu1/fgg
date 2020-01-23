package main

import (
	"fmt"
	//"os"
	"reflect"
	"strconv"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
	"github.com/rhu1/fgg/fgr"
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

/*const (
	EVAL_TO_VAL = -1 // Must be < 0
	NO_EVAL     = -2 // Must be < EVAL_TO_VAL
)*/

type Interp interface {
	GetProgram() base.Program
	SetProgram(p base.Program)
	Eval(steps int) base.Program
}

func parseaux(verbose bool, a base.Adaptor, src string, strict bool) base.Program {
	vPrintln("\nParsing AST:")
	prog := a.Parse(strict, src) // AST (Program root)
	vPrintln(prog.String())

	vPrintln("\nChecking source program OK:")
	allowStupid := false
	prog.Ok(allowStupid)

	return prog
}

// N.B. currently FG panic comes out implicitly as an underlying run-time panic
// CHECKME: add explicit FG panics?
// If steps == EVAL_TO_VAL, then eval to value
func evalaux(intrp Interp, steps int) base.Program {
	if steps < NO_EVAL {
		panic("Invalid number of steps: " + strconv.Itoa(steps))
	}
	p_init := intrp.GetProgram()
	allowStupid := true
	vPrintln("\nEntering Eval loop:")
	vPrintln("Decls:")
	for _, v := range p_init.GetDecls() {
		vPrintln("\t" + v.String() + ";")
	}
	vPrintln("Eval steps:")
	vPrintln(fmt.Sprintf("%6d: %8s %v", 0, "", p_init.GetMain())) // Initial prog OK already checked

	done := steps > EVAL_TO_VAL || // Ignore 'done' if num steps fixed (set true, for `||!done` below)
		p_init.GetMain().IsValue() // O/w evaluate until a val -- here, check if init expr is already a val
	var rule string
	p := intrp.GetProgram() // Convenient for re-assign to p inside loop
	for i := 1; i <= steps || !done; i++ {
		p, rule = p.Eval()
		intrp.SetProgram(p)
		vPrintln(fmt.Sprintf("%6d: %8s %v", i, "["+rule+"]", p.GetMain()))
		vPrintln("Checking OK:") // TODO: maybe disable by default, enable by flag
		// TODO FIXME: check actual type preservation of e_main (not just typeability)
		p.Ok(allowStupid)
		if !done && p.GetMain().IsValue() {
			done = true
		}
	}
	p_res := intrp.GetProgram()
	vPrintln(p_res.GetMain().ToGoString()) // Final result
	return p_res
}

/* FG */

type FGInterp struct {
	verbose bool
	prog    fg.FGProgram
}

var _ Interp = &FGInterp{}

func NewFGInterp(verbose bool, src string, strict bool) *FGInterp {
	var a fg.FGAdaptor
	prog := parseaux(verbose, &a, src, strict).(fg.FGProgram)
	return &FGInterp{verbose, prog}
}

func (intrp *FGInterp) GetProgram() base.Program {
	return intrp.prog
}

func (intrp *FGInterp) SetProgram(p base.Program) {
	intrp.prog = p.(fg.FGProgram)
}

func (intrp *FGInterp) Eval(steps int) base.Program {
	return evalaux(intrp, steps)
}

/* FGR */

type FGRInterp struct {
	verbose bool
	prog    fgr.FGRProgram
}

var _ Interp = &FGRInterp{}

func NewFGRInterp(verbose bool, p fgr.FGRProgram) *FGRInterp {
	return &FGRInterp{verbose, p}
}

func (intrp *FGRInterp) GetProgram() base.Program {
	return intrp.prog
}

func (intrp *FGRInterp) SetProgram(p base.Program) {
	intrp.prog = p.(fgr.FGRProgram)
}

func (intrp *FGRInterp) Eval(steps int) base.Program {
	return evalaux(intrp, steps)
}

/* FGG */

type FGGInterp struct {
	verbose bool
	prog    fgg.FGGProgram
}

var _ Interp = &FGGInterp{}

func NewFGGInterp(verbose bool, src string, strict bool) *FGGInterp {
	var a fgg.FGGAdaptor
	prog := parseaux(verbose, &a, src, strict).(fgg.FGGProgram)
	return &FGGInterp{verbose, prog}
}

func (intrp *FGGInterp) GetProgram() base.Program {
	return intrp.prog
}

func (intrp *FGGInterp) SetProgram(p base.Program) {
	intrp.prog = p.(fgg.FGGProgram)
}

func (intrp *FGGInterp) Eval(steps int) base.Program {
	return evalaux(intrp, steps)
}

/* Helpers */

/*func vPrint(x string) {
	if verbose {
		fmt.Print(x)
	}
}

func vPrintln(x string) {
	if verbose {
		fmt.Println(x)
	}
}*/
