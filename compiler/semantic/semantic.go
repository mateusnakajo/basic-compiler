package semantic

import (
	"fmt"
	"strconv"

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

type ForVariables struct {
	identifier string
	limit      float64
	step       float64
	line       string
}

type ArrayIdentifier struct {
	name  string
	index int
}

type Semantic struct {
	TokenEvents    []compiler.Event
	NewTokenEvents []compiler.Event

	IndexOfLine map[string]int
	compiler.EventDrivenModule
	strings   []string
	DataFloat map[string]float64
	DataArray map[string][]float64

	accString       string
	accFloat        float64
	Expression      string
	IfExpression    string
	identifier      string
	arrayIdentifier ArrayIdentifier
	Rerun           bool
	forIdentifiers  []ForVariables
	varToAssign     string
	arrayToAssign   ArrayIdentifier
}

func NewSemantic() Semantic {
	semantic := Semantic{}
	semantic.Rerun = false
	semantic.DataFloat = make(map[string]float64)
	semantic.DataArray = make(map[string][]float64)
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

func evaluateBoolean(expression string) bool {
	eval, _ := govaluate.NewEvaluableExpression(expression)
	temp, _ := eval.Evaluate(nil)
	return temp.(bool)
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
	handlers := map[string]func(interface{}){
		"addToExp":            s.addToExpHandler,
		"addIdentifierToExp":  s.addIdentifierToExpHandler,
		"createNewAssign":     s.createNewAssignHandler,
		"saveIdentifier":      s.saveIdentifier,
		"print":               s.printHandler,
		"goto":                s.gotoHandler,
		"ifExp":               s.ifExpHandler,
		"evaluateIf":          s.evaluateIfHandler,
		"ifComparator":        s.ifComparatorHandler,
		"forAssign":           s.forAssignHandler,
		"forLimit":            s.forLimitHandler,
		"stepFor":             s.stepForHandler,
		"endFor":              s.endForHandler,
		"saveArrayIdentifier": s.saveArrayIdentifierHandler,
		"defineArray":         s.defineArrayHandler,
		"defineArrayIndex":    s.defineArrayIndexHandler,
		"varToAssign":         s.varToAssignHandler,
		"beginExpression":     s.beginExpressionHandler,
	}
	handler := handlers[event.Name]
	handler(event.Arg)
}

func (s *Semantic) addToExpHandler(v interface{}) {
	s.Expression += v.(string)
}

func (s *Semantic) addIdentifierToExpHandler(v interface{}) {
	s.Expression += fmt.Sprintf("%f", s.DataFloat[s.identifier])
}

func (s *Semantic) addArrayIdentifierToExpHandler(v interface{}) { //FIXME
	s.Expression += fmt.Sprintf("%f", s.DataFloat[v.(string)])
}

func (s *Semantic) createNewAssignHandler(v interface{}) {
	s.DataFloat[s.varToAssign] = evaluate(s.Expression)

	// if s.arrayIdentifier != (ArrayIdentifier{}) {
	// 	s.DataArray[s.arrayIdentifier.name][s.arrayIdentifier.index] = evaluate(s.Expression)
	// 	s.arrayIdentifier = ArrayIdentifier{}
	// } else {
	// 	s.DataFloat[s.identifier] = evaluate(s.Expression)
	// }
	s.Expression = ""
}

func (s *Semantic) saveIdentifier(identifier interface{}) {
	s.identifier = identifier.(string)
}

func (s *Semantic) printHandler(v interface{}) {
	fmt.Println(s.Expression)
	s.Expression = ""
}

func (s *Semantic) gotoHandler(v interface{}) {
	start := s.IndexOfLine[v.(string)]
	s.NewTokenEvents = s.TokenEvents[start:len(s.TokenEvents)]
	s.Rerun = true
	for !s.IsEmpty() {
		_ = s.PopEvent()
	}
}

func (s *Semantic) ifExpHandler(v interface{}) {
	s.IfExpression += s.Expression
	s.Expression = ""
	fmt.Println("DEBUG:", s.IfExpression)
}

func (s *Semantic) evaluateIfHandler(v interface{}) {
	evalIf := evaluateBoolean(s.IfExpression)
	s.IfExpression = ""
	if evalIf {
		s.gotoHandler(v)
	}
}

func (s *Semantic) ifComparatorHandler(v interface{}) {
	s.IfExpression += v.(string)
}

func (s *Semantic) forAssignHandler(v interface{}) {
	if _, ok := s.DataFloat[s.identifier]; !ok {
		s.DataFloat[s.identifier] = evaluate(s.Expression)
		s.forIdentifiers = append(s.forIdentifiers, ForVariables{s.identifier, 0, 0, v.(string)})
	}
	s.identifier = ""
	s.Expression = ""
}

func (s *Semantic) forLimitHandler(v interface{}) {
	s.forIdentifiers[len(s.forIdentifiers)-1].limit = evaluate(s.Expression)
	s.Expression = ""
}

func (s *Semantic) stepForHandler(v interface{}) {
	s.forIdentifiers[len(s.forIdentifiers)-1].step = evaluate(s.Expression)
	s.Expression = ""
}

func (s *Semantic) endForHandler(v interface{}) {
	forIdentifier := s.forIdentifiers[len(s.forIdentifiers)-1].identifier
	forStep := s.forIdentifiers[len(s.forIdentifiers)-1].step
	forLimit := s.forIdentifiers[len(s.forIdentifiers)-1].limit
	forLine := s.forIdentifiers[len(s.forIdentifiers)-1].line

	s.DataFloat[forIdentifier] += forStep
	if s.DataFloat[forIdentifier] <= forLimit {
		s.gotoHandler(forLine)
	} else {
		delete(s.DataFloat, forIdentifier)
		s.forIdentifiers = s.forIdentifiers[:len(s.forIdentifiers)-1]
	}
}

func (s *Semantic) saveArrayIdentifierHandler(identifier interface{}) {
	s.arrayIdentifier = ArrayIdentifier{s.identifier, int(evaluate(s.Expression))}
	s.identifier = ""
	s.Expression = ""
}

func (s *Semantic) defineArrayHandler(v interface{}) {
	s.arrayIdentifier = ArrayIdentifier{v.(string), 0}
}

func (s *Semantic) defineArrayIndexHandler(v interface{}) {
	s.arrayIdentifier.index, _ = strconv.Atoi(v.(string))
	s.DataArray[s.arrayIdentifier.name] = make([]float64, s.arrayIdentifier.index, s.arrayIdentifier.index)
}

func (s *Semantic) varToAssignHandler(v interface{}) {
	if s.arrayIdentifier != (ArrayIdentifier{}) {
		s.arrayToAssign = s.arrayIdentifier
		s.arrayIdentifier = ArrayIdentifier{}
	} else {
		s.varToAssign = s.identifier
		s.identifier = ""
	}
}

func (s *Semantic) beginExpressionHandler(v interface{}) {
	s.Expression = ""
	s.identifier = ""
}

//TODO: zerar exp na começar exp
