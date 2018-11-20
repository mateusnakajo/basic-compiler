package lexer

import (
	"fmt"
)

type State struct {
	name    string
	next    func(*fsm, Token) State
	isFinal bool
}

type fsmInterface interface {
	ConsumeToken(Token)
	GetCurrent() State
	GetChildren() fsmInterface
	GetName() string
}

type fsm struct {
	initial  State
	current  State
	children fsmInterface
	name     string
}

func (f *fsm) ConsumeToken(token Token) {
	f.current = f.current.next(f, token)
}

func (f fsm) GetChildren() fsmInterface {
	return f.children
}

func (f fsm) GetCurrent() State {
	return f.current
}

func (f fsm) GetName() string {
	return f.name
}

type program struct {
	fsm
}

func NewProgram() program {
	program := program{}
	program.name = "program"
	nextState := State{
		name:    "1",
		isFinal: true}
	initialState := State{
		name: "0",
		next: func(f *fsm, t Token) State {
			b := NewBStatement()
			b.ConsumeToken(t)
			f.children = &b
			return nextState
		},
		isFinal: false}
	program.initial = initialState
	program.current = initialState
	return program
}

type bstatement struct {
	fsm
}

func NewBStatement() bstatement {
	bstatement := bstatement{}
	bstatement.name = "bstatement"
	nextState := State{
		name:    "1",
		isFinal: true}
	initialState := State{
		name: "0",
		next: func(f *fsm, t Token) State {
			if t.tokenType == Number {
				return nextState
			}
			return nextState // FIXME
		},
		isFinal: false}
	bstatement.initial = initialState
	bstatement.current = initialState
	bstatement.children = nil
	return bstatement
}

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
	leaf := s.program
	token = Token{lexeme: "1"}
	for {
		if leaf.GetChildren() == nil {
			break
		}
		leaf = leaf.GetChildren()
	}
	fmt.Println(leaf.GetName())
	if !leaf.GetCurrent().isFinal {
		leaf.ConsumeToken(token)
	}
}
