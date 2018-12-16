package syntactic

import (
	compiler "github.com/mateusnakajo/basic-compiler/compiler"
	lexer "github.com/mateusnakajo/basic-compiler/compiler/lexer"
)

type State struct {
	name    string
	next    func(*fsm, lexer.Token, *Stack, *string, func(compiler.Event)) State
	isFinal bool
}

func invalidState() State {
	return State{name: "INVALID", isFinal: true}
}

type fsmInterface interface {
	ConsumeToken(lexer.Token, *Stack, *string, func(compiler.Event))
	GetCurrent() State
	GetName() string
	InInvalidState() bool
}

type fsm struct {
	initial State
	current State
	name    string
}

func (f *fsm) ConsumeToken(
	token lexer.Token,
	s *Stack,
	numberOfNewLine *string,
	external func(compiler.Event)) {
	f.current = f.current.next(f, token, s, numberOfNewLine, external)
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
