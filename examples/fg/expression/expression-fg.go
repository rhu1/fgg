//$ go run github.com/rhu1/fgg -eval=-1 -v fg/examples/expression/expression-fg.go

// An almost-solution to the Expression Problem in FG.

// 1. Define Eval() on Con and Plus
// 2. Define String() on Con and Plus
// 3. Define Eval() on Neg.
// 4. Define String() on Neg.
// Each of 1, 2, 3, 4 should be doable without
// altering the answers to the others.

// It is not quite a solution, because one wants
// to first define Eval() and only later define String(),
// and to do so *without* altering the Expr interface
// (the line in Expr marked "// 2").

package main;

// 1

// TODO
type Int interface { isInt() Int };
type One struct { };
func (x0 One) isInt() Int { return x0 };

type Expr interface {
  Eval() Int;
  String() string // 2  // TODO: WF
};

type Con struct {
  value Int
};

func (e Con) Eval() Int {
  return e.value
};

type Add struct {
  left Expr;
  right Expr
};

func (e Add) Eval() Int {
  return e.left.Eval()// + e.right.Eval()  // TODO
};

// 2

func (e Con) String() string {
  return fmt.Sprintf("%v", e.value)
};

func (e Add) String() string {
  return fmt.Sprintf("(%v+%v)", e.left.String(), e.right.String())
};

// 3

type Neg struct {
  expr Expr
};

func (e Neg) Eval() Int {
  //return - e.expr.Eval()
  return One{}  // TODO
};

// 4

func (e Neg) String() string {
  return fmt.Sprintf("-%v", e.expr.String())
};

func main() {
	_ = Add{Con{One{}}, Con{One{}}}.String()
}
