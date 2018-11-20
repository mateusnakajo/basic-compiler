package lexer

type State struct {
	name    string
	next    func(*fsm, Token) State
	isFinal bool
}

type fsmInterface interface {
	ConsumeToken(Token)
	GetCurrent() State
	GetChildren() fsmInterface
	GetTest() int
}

type fsm struct {
	initial  State
	current  State
	children fsmInterface
	test     int
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

func (f fsm) GetTest() int {
	return f.test
}

type program struct {
	fsm
}

func NewProgram() program {
	program := program{}
	nextState := State{
		name:    "1",
		isFinal: true}
	initialState := State{
		name: "0",
		next: func(f *fsm, t Token) State {
			b := NewBStatement()
			b.ConsumeToken(t)
			b.test = 23
			f.children = &b
			return nextState
		},
		isFinal: false}
	program.initial = initialState
	program.current = initialState
	program.test = 22
	return program
}

type bstatement struct {
	fsm
}

func NewBStatement() bstatement {
	bstatement := bstatement{}
	nextState := State{
		name:    "1",
		isFinal: true}
	initialState := State{
		name: "0",
		next: func(f *fsm, t Token) State {
			if t.lexeme == "1" {
				return nextState
			}
			return State{}
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
	for ; leaf.GetChildren() != nil; leaf = leaf.GetChildren() {
	}
	leaf.ConsumeToken(token)
	if leaf.GetCurrent().isFinal {
		leaf = nil
	}
}
