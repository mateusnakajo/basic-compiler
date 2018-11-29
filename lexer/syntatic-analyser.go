package lexer

import "fmt"

type syntaticAnalyser struct {
	EventDrivenModule
	program  fsmInterface
	fsmStack Stack
}

func NewSyntaticAnalyser() syntaticAnalyser {
	syntaticAnalyser := syntaticAnalyser{}
	program := NewProgram()
	syntaticAnalyser.fsmStack.AddFSM(&program)
	return syntaticAnalyser
}

func (s *syntaticAnalyser) HandleEvent(event Event) {
	handlers := map[string]func(Token){
		"consumeToken": s.ConsumeToken}
	handler := handlers[event.Name]
	handler(event.Arg.(Token))
}

func (s *syntaticAnalyser) ConsumeToken(token Token) {
	fmt.Println(token)
	// leaf := s.fsmStack.TopFSM()
	// leaf.ConsumeToken(token, &s.fsmStack)
	// fmt.Println(s.fsmStack.TopFSM())
	// for !s.fsmStack.IsEmpty() && s.fsmStack.TopFSM().GetCurrent().isFinal { //FIXME
	// 	s.fsmStack.PopFSM()
	// 	fmt.Println(s.fsmStack)
	// }
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
