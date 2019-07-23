package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

/* For "strict" parsing, *parser* errors -- cf. F(G)GBailLexer */

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
	//panic(message) // + " " + e.GetMessage())
	fmt.Println(errors.New(message))
	os.Exit(1)
}

//ErrorStrategy.RecoverInline(Parser) Token
func (s *StrictErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	token := recognizer.GetCurrentToken()
	message := "[Parser] error at line " + strconv.Itoa(token.GetLine()) +
		", position " + strconv.Itoa(token.GetColumn()) + " right before " +
		s.GetTokenErrorDisplay(token)
	//panic(message) // + fmt.Sprintf("%v", antlr.NewInputMisMatchException(recognizer)))
	fmt.Println(errors.New(message))
	os.Exit(1)
	return nil
}
