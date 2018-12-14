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
	dataFloat  map[string]float64
	dataArray  map[string][]float64
	accString  string
	accFloat   float64
	expression string
}

func (s *Semantic) popString() string {
	temp := s.accString
	s.accString = ""
	return temp
}

func (s *Semantic) popFloat() float64 {
	temp := s.accFloat
	s.accFloat = 0
	return temp
}

func (s *Semantic) saveString(value string) {
	s.accString = value
}

func (s *Semantic) saveInt(value float64) {
	s.accFloat = value
}

func (s *Semantic) evaluate() {
	fmt.Println("ANTES EVAL")
	fmt.Println(s.expression)
	expression, _ := govaluate.NewEvaluableExpression(s.expression)
	temp, _ := expression.Evaluate(nil)
	s.saveInt(temp.(float64))
	s.expression = ""
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
