package syntactic

import (
	lexer "github.com/mateusnakajo/basic-compiler/compiler/lexer"
	semantic "github.com/mateusnakajo/basic-compiler/compiler/semantic"
)

type State struct {
	name    string
	next    func(*fsm, lexer.Token, *Stack, *semantic.Semantic) State
	isFinal bool
}

func invalidState() State {
	return State{name: "INVALID", isFinal: true}
}

type fsmInterface interface {
	ConsumeToken(lexer.Token, *Stack, *semantic.Semantic)
	GetCurrent() State
	GetName() string
	InInvalidState() bool
}

type fsm struct {
	initial  State
	current  State
	name     string
	assembly semantic.AssemblyInterface
}

func (f *fsm) ConsumeToken(token lexer.Token, s *Stack, semantic *semantic.Semantic) {
	//fmt.Println("Antes", f.GetName(), f.GetCurrent().name, token)
	f.current = f.current.next(f, token, s, semantic)
	//fmt.Println("Depois", f.GetName(), f.GetCurrent().name, token)
}

func (f fsm) GetCurrent() State {
	return f.current
}

func (f fsm) GetName() string {
	return f.name
}

func (f fsm) InInvalidState() bool {
	return f.current.name == "INVALID"
}

type program struct {
	fsm
}
