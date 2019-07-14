package fg

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Pre: len(elems) > 1
// elems[:len(elems)-1] -- type/meth decls
// elems[len(elems)-1] -- "main" func body expression
func MakeFgProgram(elems ...string) string {
	if len(elems) == 0 {
		panic("Bad empty args: must supply at least body expression for \"main\"")
	}
	var b strings.Builder
	b.WriteString("package main;\n")
	for _, v := range elems[:len(elems)-1] {
		b.WriteString(v)
		b.WriteString(";\n")
	}
	b.WriteString("func main() { _ = " + elems[len(elems)-1] + "}")
	return b.String()
}

// For testing
func parseAndCheckOk(prog string) {
	var adptr FGAdaptor
	ast := adptr.Parse(true, prog)
	ast.Ok()
}

func parseAndOkGood(t *testing.T, elems ...string) {
	prog := MakeFgProgram(elems...)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: " + fmt.Sprintf("%v", r) + "\n" + prog)
		}
	}()
	parseAndCheckOk(prog)
}

// N.B. do not use to check for bad *syntax* -- see below, "[Parser]" panic check
func parseAndOkBad(t *testing.T, msg string, elems ...string) {
	prog := MakeFgProgram(elems...)
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
	parseAndCheckOk(prog)
}

// Cf. https://stackoverflow.com/questions/51683104/how-to-catch-minor-errors
type StrictErrorStrategy struct {
	antlr.DefaultErrorStrategy
}

var _ antlr.ErrorStrategy = &StrictErrorStrategy{}

func (s *StrictErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {
	token := recognizer.GetCurrentToken()
	message := "[Parser] error at line {0}, position {1} right before {2} " +
		strconv.Itoa(token.GetLine()) + ":" + strconv.Itoa(token.GetColumn()) + ": " +
		s.GetTokenErrorDisplay(token)
	panic(message) // + " " + e.GetMessage())
}

//ErrorStrategy.RecoverInline(Parser) Token
func (s *StrictErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	token := recognizer.GetCurrentToken()
	message := "[Parser] error at line {0}, position {1} right before {2} " +
		strconv.Itoa(token.GetLine()) + ":" + strconv.Itoa(token.GetColumn()) + ": " +
		s.GetTokenErrorDisplay(token)
	panic(message) // + fmt.Sprintf("%v", antlr.NewInputMisMatchException(recognizer)))
}

/*public class BailLexer : ...FGLexer...
{
    public BailLexer(ICharStream input) : base(input) { }

    public override void Recover(LexerNoViableAltException e)
    {
        string message = string.Format("lex error after token {0} at position {1}", _lasttoken.Text, e.StartIndex);
        BasicEnvironment.SyntaxError = message;
        BasicEnvironment.ErrorStartIndex = e.StartIndex;
        throw new ParseCanceledException(BasicEnvironment.SyntaxError);
    }
}*/
