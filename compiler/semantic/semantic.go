package semantic

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/mateusnakajo/basic-compiler/compiler"
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
	compiler.EventDrivenModule
	strings    []string
	DataFloat  map[string]float64
	DataArray  map[string][]float64
	accString  string
	accFloat   float64
	Expression string
	identifier string
}

func NewSemantic() Semantic {
	semantic := Semantic{}
	semantic.DataFloat = make(map[string]float64)

	return semantic
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

func evaluate(expression string) float64 {
	eval, _ := govaluate.NewEvaluableExpression(expression)
	temp, _ := eval.Evaluate(nil)
	return temp.(float64)
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

func (s *Semantic) HandleEvent(event compiler.Event) {
	handlers := map[string]func(string){
		"addToExp":           s.addToExpHandler,
		"addIdentifierToExp": s.addIdentifierToExpHandler,
		"createNewAssign":    s.createNewAssignHandler,
		"saveIdentifier":     s.saveIdentifier,
		"print":              s.printHandler,
	}
	handler := handlers[event.Name]
	handler(event.Arg.(string))
}

func (s *Semantic) addToExpHandler(v string) {
	s.Expression += v
}

func (s *Semantic) addIdentifierToExpHandler(v string) {
	s.Expression += fmt.Sprintf("%f", s.DataFloat[v])
}

func (s *Semantic) createNewAssignHandler(v string) {
	s.DataFloat[s.identifier] = evaluate(s.Expression)
	fmt.Println(s.DataFloat)
	s.identifier = ""
	s.Expression = ""
}

func (s *Semantic) saveIdentifier(identifier string) {
	s.identifier = identifier
}

func (s *Semantic) printHandler(v string) {
	fmt.Println(s.Expression)
	s.Expression = ""
}
