package frontend

func FGmain(verbose bool, src string, strictParse bool, evalSteps int,
	printf bool) {
	intrp_fg := NewFGInterp(verbose, src, strictParse)
	if evalSteps > NO_EVAL {
		intrp_fg.Eval(evalSteps)
		PrintResult(printf, intrp_fg.GetProgram())
	}
}
