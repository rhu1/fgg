package fg

import (
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/rhu1/fgg/parser/fg"
)

// Pre: len(elems) > 1
// Pre: elems[:len(elems)-1] -- type/meth decls; elems[len(elems)-1] -- "main" func body expression
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
	b.WriteString("func main() { _ = " + elems[len(elems)-1] + " }")
	return b.String()
}

/* For "strict" parsing, *parser* errors */

// Cf. https://stackoverflow.com/questions/51683104/how-to-catch-minor-errors
type StrictErrorStrategy struct {
	antlr.DefaultErrorStrategy
}

var _ antlr.ErrorStrategy = &StrictErrorStrategy{}

func (s *StrictErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {
	token := recognizer.GetCurrentToken()
	message := "[Parser] error at line " + strconv.Itoa(token.GetLine()) +
		", position " + strconv.Itoa(token.GetColumn()) + " right before " +
		s.GetTokenErrorDisplay(token)
	panic(message) // + " " + e.GetMessage())
}

//ErrorStrategy.RecoverInline(Parser) Token
func (s *StrictErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	token := recognizer.GetCurrentToken()
	message := "[Parser] error at line " + strconv.Itoa(token.GetLine()) +
		", position " + strconv.Itoa(token.GetColumn()) + " right before " +
		s.GetTokenErrorDisplay(token)
	panic(message) // + fmt.Sprintf("%v", antlr.NewInputMisMatchException(recognizer)))
}

/* For "strict" parsing, *lexer* errors */

type FGBailLexer struct {
	*parser.FGLexer
}

// FIXME: not working -- e.g., incr{1}, bad token
// Want to "override" *BaseLexer.Recover -- XXX that's not how Go works (because BaseLexer is a struct, not interface)
func (b *FGBailLexer) Recover(re antlr.RecognitionException) {
	message := "lex error after token " + re.GetOffendingToken().GetText() +
		" at position " + strconv.Itoa(re.GetOffendingToken().GetStart())
	panic(message)
}

/*public FGBailLexer(ICharStream input) : base(input) { }

public override void Recover(LexerNoViableAltException e)
{
	string message = string.Format("lex error after token {0} at position {1}", _lasttoken.Text, e.StartIndex);
	BasicEnvironment.SyntaxError = message;
	BasicEnvironment.ErrorStartIndex = e.StartIndex;
	throw new ParseCanceledException(BasicEnvironment.SyntaxError);
}*/
