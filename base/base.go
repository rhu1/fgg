package base

import "fmt"
import "strings"
import "testing"

type Adaptor interface {
	Parse(strictParse bool, input string) Program
}

type Name = string // Type alias (cf. definition)

type AstNode interface {
	String() string
}

type Decl interface {
	AstNode
	GetName() Name
	Ok(ds []Decl)
}

type Program interface {
	AstNode
	GetDecls() []Decl
	GetExpr() Expr
	Ok(allowStupid bool)
	Eval() (Program, string) // Eval one step; string is the name of the (innermost) applied rule
}

type Expr interface {
	AstNode
	IsValue() bool
}

/* Test harness functions */

func parseAndCheckOk(a Adaptor, prog string) Program {
	ast := a.Parse(true, prog)
	allowStupid := false
	ast.Ok(allowStupid)
	return ast
}

func ParseAndOkGood(t *testing.T, a Adaptor, prog string) Program {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: " + fmt.Sprintf("%v", r) + "\n" +
				prog)
		}
	}()
	return parseAndCheckOk(a, prog)
}

// N.B. do not use to check for bad *syntax* -- see the "[Parser]" panic check
func ParseAndOkBad(t *testing.T, msg string, a Adaptor, prog string) Program {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but none occurred: " + msg + "\n" +
				prog)
		} else {
			rec := fmt.Sprintf("%v", r)
			if strings.HasPrefix(rec, "[Parser]") {
				t.Errorf("Unexpected panic: " + rec + "\n" + prog)
			}
			// TODO FIXME: check panic more specifically
		}
	}()
	return parseAndCheckOk(a, prog)
}

// Pre: parseAndOkGood
func EvalAndOkGood(t *testing.T, p Program, steps int) Program {
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
func EvalAndOkBad(t *testing.T, p Program, msg string, steps int) Program {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but none occurred: " + msg + "\n" +
				p.String())
		} else {
			// [Parser] panic should be already checked by parseAndOkGood
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
