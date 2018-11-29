package lexer

type State struct {
	name    string
	next    func(*fsm, Token, *Stack) State
	isFinal bool
}

func invalidState() State {
	return State{name: "INVALID", isFinal: true}
}

type fsmInterface interface {
	ConsumeToken(Token, *Stack)
	GetCurrent() State
	GetName() string
}

type fsm struct {
	initial State
	current State
	name    string
}

func (f *fsm) ConsumeToken(token Token, s *Stack) {
	//fmt.Println("Antes", f.GetName(), f.GetCurrent().name, token)
	f.current = f.current.next(f, token, s)
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
