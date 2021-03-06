package semantic

import (
	"fmt"
	"strconv"

	"github.com/Knetic/govaluate"
	"github.com/mateusnakajo/basic-compiler/compiler"
)

type StackIdentifier struct {
	identifier []string
}

func (s *StackIdentifier) Add(i string) {
	s.identifier = append(s.identifier, i)
}

func (s *StackIdentifier) Pop() string {
	lastIndex := len(s.identifier) - 1
	last := s.identifier[lastIndex]
	s.identifier = s.identifier[0:lastIndex]
	return last
}

func (s *StackIdentifier) Top() string {
	return s.identifier[len(s.identifier)-1]
}

func (s StackIdentifier) IsEmpty() bool {
	return len(s.identifier) == 0
}

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
	ExpressionSaved string
	IfExpression    string
	identifiers     StackIdentifier
	arrayIdentifier ArrayIdentifier
	Rerun           bool
	forIdentifiers  []ForVariables
	varToAssign     string
	arrayToAssign   ArrayIdentifier
	identifierSaved string
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
		"saveIdentifier":      s.saveIdentifierHandler,
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
		"saveExpression":      s.saveExpressionHandler,
	}
	handler := handlers[event.Name]
	handler(event.Arg)
}

func (s *Semantic) addToExpHandler(v interface{}) {
	s.Expression += v.(string)
}

func (s *Semantic) addIdentifierToExpHandler(v interface{}) {
	s.Expression += fmt.Sprintf("%f", s.DataFloat[s.identifiers.Top()])
}

func (s *Semantic) addArrayIdentifierToExpHandler(v interface{}) { //FIXME
	s.Expression += fmt.Sprintf("%f", s.DataFloat[v.(string)])
}

func (s *Semantic) createNewAssignHandler(v interface{}) {
	if s.arrayToAssign != (ArrayIdentifier{}) {
		s.DataArray[s.arrayToAssign.name][s.arrayToAssign.index] = evaluate(s.Expression)
		s.arrayToAssign = ArrayIdentifier{}
	} else {
		s.DataFloat[s.varToAssign] = evaluate(s.Expression)
	}
	s.Expression = ""
}

func (s *Semantic) saveIdentifierHandler(identifier interface{}) {
	s.identifiers.Add(identifier.(string))
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
	for len(s.identifiers.identifier) > 1 {
		s.identifiers.Pop()
	}
	if _, ok := s.DataFloat[s.identifiers.Top()]; !ok {
		s.DataFloat[s.identifiers.Top()] = evaluate(s.Expression)
		s.forIdentifiers = append(s.forIdentifiers, ForVariables{s.identifiers.Pop(), 0, 1, v.(string)})
	}
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
	if s.identifierSaved != s.identifiers.Top() {
		s.identifiers.Pop()
	}
	s.arrayIdentifier = ArrayIdentifier{s.identifiers.Pop(), int(evaluate(s.Expression))}
	s.Expression = s.ExpressionSaved + fmt.Sprintf("%f", s.DataArray[s.arrayIdentifier.name][s.arrayIdentifier.index])
}

func (s *Semantic) defineArrayHandler(v interface{}) {
	s.arrayIdentifier = ArrayIdentifier{v.(string), 0}
}

func (s *Semantic) defineArrayIndexHandler(v interface{}) {
	s.arrayIdentifier.index, _ = strconv.Atoi(v.(string))
	s.DataArray[s.arrayIdentifier.name] = make([]float64, s.arrayIdentifier.index, s.arrayIdentifier.index)
	s.arrayIdentifier = ArrayIdentifier{}
}

func (s *Semantic) varToAssignHandler(v interface{}) {
	if s.arrayIdentifier != (ArrayIdentifier{}) {
		s.arrayToAssign = s.arrayIdentifier
		s.arrayIdentifier = ArrayIdentifier{}
	} else {
		s.varToAssign = s.identifiers.Pop()
	}
}

func (s *Semantic) beginExpressionHandler(v interface{}) {
	s.Expression = ""
}

func (s *Semantic) saveExpressionHandler(v interface{}) {
	s.ExpressionSaved = s.Expression
	s.identifierSaved = s.identifiers.Top()
}
