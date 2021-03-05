package frontend

func FGGmain(verbose bool, src string, strictParse bool, evalSteps int,
	printf bool, monom bool, monomc string, oblitc string) {
	intrp_fgg := NewFGGInterp(verbose, src, strictParse)
	if evalSteps > NO_EVAL {
		intrp_fgg.Eval(evalSteps)
		PrintResult(printf, intrp_fgg.GetProgram())
	}

	// TODO: refactor (cf. Frontend, Interp)
	intrp_fgg.Monom(monom, monomc)
	intrp_fgg.Oblit(oblitc)
	////doWrappers(prog, wrapperc)
}
