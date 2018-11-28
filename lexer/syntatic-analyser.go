package lexer

import "fmt"

type syntaticAnalyser struct {
	EventDrivenModule
	program fsmInterface
}

func NewSyntaticAnalyser() syntaticAnalyser {
	syntaticAnalyser := syntaticAnalyser{}
	program := NewProgram()
	syntaticAnalyser.program = &program
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
	leaf := s.program
	for {
		if leaf.GetChildren() == nil {
			break
		}
		if leaf.GetChildren().GetCurrent().isFinal {
			leaf.SetChildren(nil)
		} else {
			leaf = leaf.GetChildren()
		}
	}
	leaf.ConsumeToken(token)
}
