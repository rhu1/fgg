package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/rhu1/fgg/base"
	"github.com/rhu1/fgg/fg"
	"github.com/rhu1/fgg/fgg"
	"github.com/rhu1/fgg/fgr"
)

var _ = reflect.TypeOf
var _ = strconv.Itoa

const (
	EVAL_TO_VAL = -1 // Must be < 0
	NO_EVAL     = -2 // Must be < EVAL_TO_VAL
)

/* !!! WIP -- refactoring whole main package !!! */

/* -- Interp */

type Interp interface {
	GetSource() base.Program // TODO: factor out decls with below
	GetProgram() base.Program
	SetProgram(p base.Program)
	Eval(steps int) base.Type
	vPrint(x string)
	vPrintln(x string)
}

func parse(verbose bool, a base.Adaptor, src string, strict bool) base.Program {
	vPrintln(verbose, "\nParsing AST:")
	prog := a.Parse(strict, src) // AST (Program root)
	vPrintln(verbose, prog.String())

	vPrintln(verbose, "\nChecking source program OK:")
	allowStupid := false
	prog.Ok(allowStupid)

	return prog
}

// N.B. currently FG panic comes out implicitly as an underlying run-time panic
// CHECKME: add explicit FG panics?
// If steps == EVAL_TO_VAL, then eval to value
// Post: intrp.GetProgram() contains the eval result; result type is returned
func eval(intrp Interp, steps int) base.Type {
	if steps < NO_EVAL {
		panic("Invalid number of steps: " + strconv.Itoa(steps))
	}
	p_init := intrp.GetProgram()
	allowStupid := true
	t_init := p_init.Ok(allowStupid)
	intrp.vPrintln("\nEntering Eval loop:")
	intrp.vPrintln("Decls:")
	ds := p_init.GetDecls()
	for _, v := range ds {
		intrp.vPrintln("\t" + v.String() + ";")
	}
	intrp.vPrintln("Eval steps:")
	intrp.vPrintln(fmt.Sprintf("%6d: %8s %v", 0, "", p_init.GetMain())) // Initial prog OK already checked

	done := steps > EVAL_TO_VAL || // Ignore 'done' if num steps fixed (set true, for `||!done` below)
		p_init.GetMain().IsValue() // O/w evaluate until a val -- here, check if init expr is already a val
	var rule string
	p := p_init
	t := t_init
	for i := 1; i <= steps || !done; i++ {
		p, rule = p.Eval()
		intrp.SetProgram(p)
		intrp.vPrintln(fmt.Sprintf("%6d: %8s %v", i, "["+rule+"]", p.GetMain()))
		intrp.vPrint("Checking OK:") // TODO: maybe disable by default, enable by flag
		t = p.Ok(allowStupid)
		intrp.vPrintln(" " + t.String())
		if !t.Impls(ds, t_init) { // Check type preservation
			panic("Type not preserved by evaluation.")
		}
		if !done && p.GetMain().IsValue() {
			done = true
		}
	}
	intrp.vPrintln(p.GetMain().ToGoString()) // Final result
	//return p_res
	return t
}

/* FG */

type FGInterp struct {
	verboseHelper
	orig fg.FGProgram
	prog fg.FGProgram
}

var _ Interp = &FGInterp{}

func NewFGInterp(verbose bool, src string, strict bool) *FGInterp {
	var a fg.FGAdaptor
	orig := parse(verbose, &a, src, strict).(fg.FGProgram)
	prog := fg.NewFGProgram(orig.GetDecls(), orig.GetMain().(fg.FGExpr), orig.IsPrintf())
	return &FGInterp{verboseHelper{verbose}, orig, prog}
}

func (intrp *FGInterp) GetSource() base.Program   { return intrp.orig }
func (intrp *FGInterp) GetProgram() base.Program  { return intrp.prog }
func (intrp *FGInterp) SetProgram(p base.Program) { intrp.prog = p.(fg.FGProgram) }

func (intrp *FGInterp) Eval(steps int) base.Type {
	return eval(intrp, steps)
}

/* FGR */

type FGRInterp struct {
	verboseHelper
	prog fgr.FGRProgram
}

var _ Interp = &FGRInterp{}

func NewFGRInterp(verbose bool, p fgr.FGRProgram) *FGRInterp {

	return &FGRInterp{verboseHelper{verbose}, p}
}

func (intrp *FGRInterp) GetSource() base.Program   { panic("TODO") }
func (intrp *FGRInterp) GetProgram() base.Program  { return intrp.prog }
func (intrp *FGRInterp) SetProgram(p base.Program) { intrp.prog = p.(fgr.FGRProgram) }

func (intrp *FGRInterp) Eval(steps int) base.Type {
	return eval(intrp, steps)
}

/* FGG */

type FGGInterp struct {
	verboseHelper
	orig fgg.FGGProgram
	prog fgg.FGGProgram
}

var _ Interp = &FGGInterp{}

func NewFGGInterp(verbose bool, src string, strict bool) *FGGInterp {
	var a fgg.FGGAdaptor
	orig := parse(verbose, &a, src, strict).(fgg.FGGProgram)
	prog := fgg.NewProgram(orig.GetDecls(), orig.GetMain().(fgg.FGGExpr), orig.IsPrintf())
	return &FGGInterp{verboseHelper{verbose}, orig, prog}
}

func (intrp *FGGInterp) GetSource() base.Program   { return intrp.orig }
func (intrp *FGGInterp) GetProgram() base.Program  { return intrp.prog }
func (intrp *FGGInterp) SetProgram(p base.Program) { intrp.prog = p.(fgg.FGGProgram) }

func (intrp *FGGInterp) Eval(steps int) base.Type {
	return eval(intrp, steps)
}

// Pre: (monom == true || compile != "") => -fgg is set
// TODO: rename
func (intrp *FGGInterp) Monom(monom bool, compile string) {
	if !monom && compile == "" {
		return
	}

	p_fgg := intrp.GetSource().(fgg.FGGProgram)
	/*if e, ok := fgg.IsMonomable(p_fgg); !ok {
		panic("\nNot monomorphisable according to \"type param under named type\"" +
			" restriction.\n\t" + e.String())
	}*/

	fgg.Foo(intrp.orig.GetDecls())

	p_mono := fgg.Monomorph(p_fgg)
	if monom {
		intrp.vPrintln("\nMonomorphising, formal notation: [Warning] WIP [Warning]")
		fmt.Println(p_mono.String())
	}
	if compile != "" {
		intrp.vPrintln("\nMonomorphising, FG output: [Warning] WIP [Warning]")
		out := monomOutputHack(p_mono.String())
		if compile == "--" {
			fmt.Println(out)
		} else {
			intrp.vPrintln(out)
			intrp.vPrintln("Writing output to: " + compile)
			bs := []byte(out)
			err := ioutil.WriteFile(compile, bs, 0644)
			checkErr(err)
		}
	}
}

func (intrp *FGGInterp) Oblit(compile string) {
	if compile == "" {
		return
	}
	intrp.vPrintln("\nTranslating FGG to FG(R) using Obliteration: [Warning] WIP [Warning]")
	p_fgr := fgr.Obliterate(intrp.GetSource().(fgg.FGGProgram))
	out := p_fgr.String()
	// TODO: factor out with -monomc
	if compile == "--" {
		fmt.Println(out)
	} else {
		intrp.vPrintln(out)
		intrp.vPrintln("Writing output to: " + compile)
		bs := []byte(out)
		err := ioutil.WriteFile(compile, bs, 0644)
		checkErr(err)
	}

	// cf. interp -- TODO: factor out with others
	p_fgr.Ok(false)
	if oblitEvalSteps > NO_EVAL {
		intrp.vPrint("\nEvaluating FGR:") // eval prints a leading "\n"
		intrp_fgr := NewFGRInterp(verbose, p_fgr)
		intrp_fgr.Eval(oblitEvalSteps)
		fmt.Println(intrp_fgr.GetProgram().GetMain())
	}
}

/* Aux */

func monomOutputHack(out string) string {
	out = strings.Replace(out, ",,", "ᐨ", -1) // TODO: refactor -- cf. fgg_monom, toMonomId
	out = strings.Replace(out, "<", "ᐸ", -1)
	out = strings.Replace(out, ">", "ᐳ", -1)
	return out
}

/* Helpers */

type verboseHelper struct {
	verbose bool
}

func (vh verboseHelper) GetVerbose() bool {
	return vh.verbose
}

func (vh verboseHelper) vPrint(x string) {
	vPrint(vh.verbose, x)
}

func (vh verboseHelper) vPrintln(x string) {
	vPrintln(vh.verbose, x)
}

func vPrint(verbose bool, x string) {
	if verbose {
		fmt.Print(x)
	}
}

func vPrintln(verbose bool, x string) {
	if verbose {
		fmt.Println(x)
	}
}
