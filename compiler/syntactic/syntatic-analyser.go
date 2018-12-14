package syntactic

import (
	"fmt"

	compiler "github.com/mateusnakajo/basic-compiler/compiler"
	lexer "github.com/mateusnakajo/basic-compiler/compiler/lexer"
	semantic "github.com/mateusnakajo/basic-compiler/compiler/semantic"
)

type syntaticAnalyser struct {
	compiler.EventDrivenModule
	program  fsmInterface
	fsmStack Stack
	semantic semantic.Semantic
}

func NewSyntaticAnalyser() syntaticAnalyser {
	syntaticAnalyser := syntaticAnalyser{}
	program := NewProgram()
	syntaticAnalyser.semantic = semantic.Semantic{}
	syntaticAnalyser.semantic.DataFloat = make(map[string]float64)
	syntaticAnalyser.fsmStack.AddFSM(&program)
	return syntaticAnalyser
}

func (s *syntaticAnalyser) HandleEvent(event compiler.Event) {
	handlers := map[string]func(lexer.Token){
		"consumeToken": s.ConsumeToken}
	handler := handlers[event.Name]
	handler(event.Arg.(lexer.Token))
}

func (s *syntaticAnalyser) ConsumeToken(token lexer.Token) {

	if true {
		fmt.Println("\n")
		fmt.Println(token)
		s.fsmStack.PrintStack()
	}

	if !s.fsmStack.TopFSM().InInvalidState() {
		s.fsmStack.TopFSM().ConsumeToken(token, &s.fsmStack, &s.semantic)
	}
	for s.fsmStack.TopFSM().InInvalidState() {
		s.fsmStack.PopFSM()
		s.fsmStack.TopFSM().ConsumeToken(token, &s.fsmStack, &s.semantic)
	}
}

type Stack struct {
	fsm []fsmInterface
}

func (s *Stack) AddFSM(fsm fsmInterface) {
	s.fsm = append(s.fsm, fsm)
}

func (s *Stack) PopFSM() fsmInterface {
	lastIndex := len(s.fsm) - 1
	lastFSM := s.fsm[lastIndex]
	s.fsm = s.fsm[0:lastIndex]
	return lastFSM
}

func (s *Stack) TopFSM() fsmInterface {
	return s.fsm[len(s.fsm)-1]
}

func (s Stack) IsEmpty() bool {
	return len(s.fsm) == 0
}

func (s Stack) PrintStack() {
	for i := range s.fsm {
		fmt.Println("FSM:", s.fsm[i].GetName(), "Estado:", s.fsm[i].GetCurrent())
	}
}
