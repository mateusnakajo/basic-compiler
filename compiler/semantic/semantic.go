package lexer

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

// type Semantic struct {
// 	functions []Function
// }

// type Function struct {
// 	assembly    string
// 	returnvalue Expression
// 	name        string
// }

// type Expression struct {
// 	assembly string
// }

// func (f *Function) generateFunction() {
// 	f.assembly = fmt.Sprintf(`%v:
// 		pushq	%%rbp
// 		movq	%%rsp, %%rbp
// 		movl	$1, %%eax
// 		popq	%%rbp
// 		ret`, f.name)
// }

type Semantic struct {
	strings    []string
	DataFloat  map[string]float64
	DataArray  map[string][]float64
	accString  string
	accFloat   float64
	Expression string
}

func (s *Semantic) PopString() string {
	temp := s.accString
	s.accString = ""
	return temp
}

func (s *Semantic) PopFloat() float64 {
	temp := s.accFloat
	s.accFloat = 0
	return temp
}

func (s *Semantic) SaveString(value string) {
	s.accString = value
}

func (s *Semantic) SaveInt(value float64) {
	s.accFloat = value
}

func (s *Semantic) Evaluate() {
	fmt.Println("ANTES EVAL")
	fmt.Println(s.Expression)
	expression, _ := govaluate.NewEvaluableExpression(s.Expression)
	temp, _ := expression.Evaluate(nil)
	s.SaveInt(temp.(float64))
	s.Expression = ""
}

type ExpressionAssembly struct {
	assembly string
}

func (e ExpressionAssembly) getAssembly() string {
	return e.assembly
}

type AssemblyInterface interface {
	getAssembly() string
}
