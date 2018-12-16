package syntactic

import (
	"fmt"

	compiler "github.com/mateusnakajo/basic-compiler/compiler"
	lexer "github.com/mateusnakajo/basic-compiler/compiler/lexer"
)

type syntaticAnalyser struct {
	compiler.EventDrivenModule
	program         fsmInterface
	fsmStack        Stack
	IndexOfLine     map[string]int
	numberOfNewLine string
	numberOfTokens  int
}

func NewSyntaticAnalyser() syntaticAnalyser {
	syntaticAnalyser := syntaticAnalyser{}
	syntaticAnalyser.IndexOfLine = make(map[string]int)
	syntaticAnalyser.numberOfTokens = 0
	syntaticAnalyser.numberOfNewLine = ""
	program := NewProgram()
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
	if _, ok := s.IndexOfLine[s.numberOfNewLine]; !ok {
		s.IndexOfLine[s.numberOfNewLine] = s.numberOfTokens - 1
	}

	s.numberOfTokens++
	if false {
		fmt.Println("\n")
		fmt.Println(token)
		s.fsmStack.PrintStack()
	}
	if !s.fsmStack.TopFSM().InInvalidState() {
		s.fsmStack.TopFSM().ConsumeToken(token, &s.fsmStack, &s.numberOfNewLine, s.AddExternal)
	}
	for s.fsmStack.TopFSM().InInvalidState() {
		s.fsmStack.PopFSM()
		s.fsmStack.TopFSM().ConsumeToken(token, &s.fsmStack, &s.numberOfNewLine, s.AddExternal)
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
		fmt.Println("FSM:", s.fsm[i], "Estado:", s.fsm[i].GetCurrent())
	}
}
