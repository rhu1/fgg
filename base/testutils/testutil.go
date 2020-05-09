package testutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rhu1/fgg/base"
)

const PARSER_PANIC_PREFIX = "[Parser] "

/* Test harness functions */

func parseAndCheckOk(a base.Adaptor, src string) base.Program {
	ast := a.Parse(true, src)
	allowStupid := false
	ast.Ok(allowStupid)
	return ast
}

func ParseAndOkGood(t *testing.T, a base.Adaptor, src string) base.Program {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: " + fmt.Sprintf("%v", r) + "\n" +
				src)
		}
	}()
	return parseAndCheckOk(a, src)
}

// N.B. do not use to check for bad *syntax* -- see the PARSER_PANIC_PREFIX panic check
func ParseAndOkBad(t *testing.T, msg string, a base.Adaptor, src string) base.Program {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but none occurred: " + msg + "\n" +
				src)
		} else {
			rec := fmt.Sprintf("%v", r)
			if strings.HasPrefix(rec, PARSER_PANIC_PREFIX) {
				t.Errorf("Unexpected panic: " + rec + "\n" + src)
			}
			// TODO FIXME: check panic more specifically
		}
	}()
	return parseAndCheckOk(a, src)
}

// Pre: parseAndOkGood
func EvalAndOkGood(t *testing.T, p base.Program, steps int) base.Program {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: " + fmt.Sprintf("%v", r) + "\n" +
				p.String())
		}
	}()
	allowStupid := true
	for i := 0; i < steps; i++ {
		p, _ = p.Eval() // CHECKME: check rule names as part of test?
		p.Ok(allowStupid)
	}
	return p
}

// Pre: parseAndOkGood
func EvalAndOkBad(t *testing.T, p base.Program, msg string, steps int) base.Program {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but none occurred: " + msg + "\n" +
				p.String())
		} else {
			// PARSER_PANIC_PREFIX panic should be already checked by parseAndOkGood
			// TODO FIXME: check panic more specifically
		}
	}()
	allowStupid := true
	for i := 0; i < steps; i++ {
		p, _ = p.Eval()
		p.Ok(allowStupid)
	}
	return p
}
