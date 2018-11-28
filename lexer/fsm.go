package lexer

type State struct {
	name    string
	next    func(*fsm, Token) State
	isFinal bool
}

func invalidState() State {
	return State{name: "INVALID", isFinal: true}
}

type fsmInterface interface {
	ConsumeToken(Token)
	GetCurrent() State
	GetChildren() fsmInterface
	GetName() string
	SetChildren(fsmInterface)
}

type fsm struct {
	initial  State
	current  State
	children fsmInterface
	name     string
}

func (f *fsm) ConsumeToken(token Token) {
	//fmt.Println("Antes", f.GetName(), f.GetCurrent().name, token)
	f.current = f.current.next(f, token)
	//fmt.Println("Depois", f.GetName(), f.GetCurrent().name, token)
}

func (f fsm) GetChildren() fsmInterface {
	return f.children
}

func (f *fsm) SetChildren(fsm fsmInterface) {
	f.children = fsm
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
