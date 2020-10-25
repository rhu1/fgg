package frontend

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/rhu1/fgg/internal/base"
	"github.com/rhu1/fgg/internal/fg"
	"github.com/rhu1/fgg/internal/fgg"
	"github.com/rhu1/fgg/internal/fgr"
	"github.com/rhu1/fgg/internal/parser"
)

var _ = os.Args
var _ = reflect.TypeOf
var _ = strconv.Itoa

const (
	EVAL_TO_VAL = -1 // Must be < 0
	NO_EVAL     = -2 // Must be < EVAL_TO_VAL
)

var (
	OblitEvalSteps int
	Verbose        bool
)

/* -- Interp */

type Interp interface {
	GetSource() base.Program // TODO: factor out decls with below
	GetProgram() base.Program
	SetProgram(p base.Program)
	Eval(steps int) base.Type
	VPrint(x string)
	VPrintln(x string)
}

func parse(verbose bool, a base.Adaptor, src string, strict bool) base.Program {
	VPrintln(verbose, "\nParsing AST:")
	prog := a.Parse(strict, src) // AST (Program root)
	VPrintln(verbose, prog.String())
	return prog
}

// N.B. currently FG panic comes out implicitly as an underlying run-time panic
// CHECKME: add explicit FG panics?
// If steps == EVAL_TO_VAL, then Eval to value
// Post: intrp.GetProgram() contains the Eval result; result type is returned
func Eval(intrp Interp, steps int) base.Type {
	if steps < NO_EVAL {
		panic("Invalid number of steps: " + strconv.Itoa(steps))
	}
	p_init := intrp.GetProgram()
	allowStupid := true
	t_init := p_init.Ok(allowStupid)
	intrp.VPrintln("\nEntering Eval loop:")
	intrp.VPrintln("Decls:")
	ds := p_init.GetDecls()
	for _, v := range ds {
		intrp.VPrintln("\t" + v.String() + ";")
	}
	intrp.VPrintln("Eval steps:")
	intrp.VPrintln(fmt.Sprintf("%6d: %8s %v", 0, "", p_init.GetMain())) // Initial prog OK already checked

	done := steps > EVAL_TO_VAL || // Ignore 'done' if num steps fixed (set true, for `||!done` below)
		p_init.GetMain().IsValue() // O/w evaluate until a val -- here, check if init expr is already a val
	var rule string
	p := p_init
	t := t_init

	for i := 1; i <= steps || !done; i++ {
		p, rule = p.Eval()
		intrp.SetProgram(p)
		intrp.VPrintln(fmt.Sprintf("%6d: %8s %v", i, "["+rule+"]", p.GetMain()))
		intrp.VPrint("Checking OK:") // TODO: maybe disable by default, enable by flag
		t = p.Ok(allowStupid)
		intrp.VPrintln(" " + t.String())
		if !t.Impls(ds, t_init) { // Check type preservation
			panic("Type not preserved by evaluation.")
		}
		if !done && p.GetMain().IsValue() { // N.B. IsValue, not CanEval -- bad asserts panics, like Go (but not actual FGG)
			done = true
		}
	}
	intrp.VPrintln(p.GetMain().String()) // Final result  // CHECKME: check prog.printf, for ToGoString?
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
	var a parser.FGAdaptor
	orig := parse(verbose, &a, src, strict).(fg.FGProgram)

	VPrintln(verbose, "\nChecking source program OK:")
	allowStupid := false
	orig.Ok(allowStupid)

	prog := fg.NewFGProgram(orig.GetDecls(), orig.GetMain().(fg.FGExpr), orig.IsPrintf())
	return &FGInterp{verboseHelper{verbose}, orig, prog}
}

func (intrp *FGInterp) GetSource() base.Program   { return intrp.orig }
func (intrp *FGInterp) GetProgram() base.Program  { return intrp.prog }
func (intrp *FGInterp) SetProgram(p base.Program) { intrp.prog = p.(fg.FGProgram) }

func (intrp *FGInterp) Eval(steps int) base.Type {
	return Eval(intrp, steps)
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
	return Eval(intrp, steps)
}

/* FGG */

type FGGInterp struct {
	verboseHelper
	orig fgg.FGGProgram
	prog fgg.FGGProgram
}

var _ Interp = &FGGInterp{}

func NewFGGInterp(verbose bool, src string, strict bool) *FGGInterp {
	var a parser.FGGAdaptor
	orig := parse(verbose, &a, src, strict).(fgg.FGGProgram)
	renamed := renameParams(orig)
	/*vPrintln(verbose, "\nRenamed method type params:")
	vPrintln(verbose, renamed.String())*/

	VPrintln(verbose, "\nChecking source program OK:")
	allowStupid := false
	renamed.Ok(allowStupid)

	prog := fgg.NewProgram(renamed.GetDecls(), renamed.GetMain().(fgg.FGGExpr),
		renamed.IsPrintf())
	return &FGGInterp{verboseHelper{verbose}, renamed, prog}
}

func (intrp *FGGInterp) GetSource() base.Program   { return intrp.orig }
func (intrp *FGGInterp) GetProgram() base.Program  { return intrp.prog }
func (intrp *FGGInterp) SetProgram(p base.Program) { intrp.prog = p.(fgg.FGGProgram) }

func (intrp *FGGInterp) Eval(steps int) base.Type {
	return Eval(intrp, steps)
}

// Pre: (monom == true || compile != "") => -fgg is set
// rename
func (intrp *FGGInterp) Monom(monom bool, compile string) {
	if !monom && compile == "" {
		return
	}

	p_fgg := intrp.GetSource().(fgg.FGGProgram)

	if ok, msg := fgg.IsMonomOK(intrp.orig); !ok { // nomonom
		//fmt.Println("\nCannot monomorphise (nomono detected):\n\t" + msg)
		panic("\nCannot monomorphise (nomono detected):\n\t" + msg)
	}

	p_mono := fgg.Monomorph(p_fgg)
	if monom {
		intrp.VPrintln("\nMonomorphising:")
		fmt.Println(p_mono.String())
	}
	if compile != "" {
		intrp.VPrintln("\nMonomorphising:")
		out := monomOutputHack(p_mono.String())
		if compile == "--" {
			fmt.Println(out)
		} else {
			intrp.VPrintln(out)
			intrp.VPrintln("Writing output to: " + compile)
			bs := []byte(out)
			err := ioutil.WriteFile(compile, bs, 0644)
			CheckErr(err)
		}
	}
}

// WIP
func (intrp *FGGInterp) Oblit(compile string) {
	if compile == "" {
		return
	}
	intrp.VPrintln("\nTranslating FGG to FG(R) using Obliteration: [Warning] WIP [Warning]")
	p_fgr := fgr.Obliterate(intrp.GetSource().(fgg.FGGProgram))
	out := p_fgr.String()
	// TODO: factor out with -monomc
	if compile == "--" {
		fmt.Println(out)
	} else {
		intrp.VPrintln(out)
		intrp.VPrintln("Writing output to: " + compile)
		bs := []byte(out)
		err := ioutil.WriteFile(compile, bs, 0644)
		CheckErr(err)
	}

	// cf. interp -- TODO: factor out with others
	p_fgr.Ok(false)
	if OblitEvalSteps > NO_EVAL {
		intrp.VPrint("\nEvaluating FGR:") // eval prints a leading "\n"
		intrp_fgr := NewFGRInterp(Verbose, p_fgr)
		intrp_fgr.Eval(OblitEvalSteps)
		fmt.Println(intrp_fgr.GetProgram().GetMain())
	}
}

/* Type param renaming */

func renameParams(p fgg.FGGProgram) fgg.FGGProgram {
	ds := make([]base.Decl, len(p.GetDecls()))
	for i, v := range p.GetDecls() {
		ds[i] = renameParamsDecl(v.(fgg.Decl))
	}
	return fgg.NewProgram(ds, p.GetMain().(fgg.FGGExpr), p.IsPrintf())
}

func renameParamsDecl(p fgg.Decl) fgg.Decl {
	switch d := p.(type) {
	case fgg.STypeLit:
		return d
	case fgg.ITypeLit:
		return renameParamsITypeLit(d)
	case fgg.MethDecl:
		return renameParamsMethDecl(d)
	default:
		panic("Unknown Decl type: " + reflect.TypeOf(d).String() +
			"\n\t" + d.String())
	}
}

func renameParamsITypeLit(c fgg.ITypeLit) fgg.ITypeLit {
	orig := c.GetSpecs()
	ss := make([]fgg.Spec, len(orig))
	for i, s1 := range orig {
		switch s := s1.(type) {
		case fgg.TNamed:
			ss[i] = s
		case fgg.Sig:
			subs := makeParamIndexSubs(s.Psi)
			ss[i] = s.TSubs(subs)
		default:
			panic("Unknown Spec type: " + reflect.TypeOf(s).String() +
				"\n\t" + s.String())
		}
	}
	return fgg.NewITypeLit(c.GetName(), c.Psi, ss)
}

func renameParamsMethDecl(m fgg.MethDecl) fgg.MethDecl {
	subs := makeParamIndexSubs(m.Psi_meth)
	tfs_orig := m.Psi_meth.GetTFormals()
	tfs := make([]fgg.TFormal, len(tfs_orig))
	for i, v := range tfs_orig {
		tfs[i] = fgg.NewTFormal(v.GetTParam().TSubs(subs).(fgg.TParam), v.GetUpperBound().TSubs(subs))
	}
	Psi_meth := fgg.NewBigPsi(tfs)
	pds_orig := m.GetParamDecls()
	pds := make([]fgg.ParamDecl, len(pds_orig))
	for i, v := range pds_orig {
		pds[i] = fgg.NewParamDecl(v.GetName(), v.GetType().TSubs(subs))
	}
	return fgg.NewMethDecl(m.GetRecvName(), m.GetRecvTypeName(), m.Psi_recv, m.GetName(), Psi_meth, pds, m.GetReturn().TSubs(subs), m.GetBody().TSubs(subs))
}

func makeParamIndexSubs(Psi fgg.BigPsi) fgg.Delta {
	subs := make(fgg.Delta)
	tfs := Psi.GetTFormals()
	for j := 0; j < len(tfs); j++ {
		subs[tfs[j].GetTParam()] = fgg.TParam("β" + strconv.Itoa(j+1)) // β for add meth params
	}
	return subs
}

/* Aux */

func monomOutputHack(out string) string {
	// TODO: refactor -- cf. fgg_monom, toMonomId
	out = strings.Replace(out, ",,", "ᐨ", -1) // U+1428 Canadian Aboriginal Syllabics Final Short Horizontal Stroke
	// U+035C Combining Double Breve Below -- CHECKME: doesn't work with ANTLR?
	out = strings.Replace(out, "<", "ᐸ", -1) // U+1438 Canadian Aboriginal Syllabics Pa
	out = strings.Replace(out, ">", "ᐳ", -1) // U+1433 Canadian Aboriginal Syllabics Po
	return out
}

/* Helpers */

type verboseHelper struct {
	verbose bool
}

func (vh verboseHelper) GetVerbose() bool {
	return vh.verbose
}

func (vh verboseHelper) VPrint(x string) {
	VPrint(vh.verbose, x)
}

func (vh verboseHelper) VPrintln(x string) {
	VPrintln(vh.verbose, x)
}

func VPrint(verbose bool, x string) {
	if verbose {
		fmt.Print(x)
	}
}

func VPrintln(verbose bool, x string) {
	if verbose {
		fmt.Println(x)
	}
}

// ECheckErr
func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}

func PrintResult(printf bool, p base.Program) {
	res := p.GetMain()
	if printf {
		fmt.Println(res.ToGoString(p.GetDecls()))
	} else {
		fmt.Println(res)
	}
}

func ReadSourceFile(arg0 string) string {
	b, err := ioutil.ReadFile(arg0)
	if err != nil {
		CheckErr(err)
	}
	return string(b)
}
