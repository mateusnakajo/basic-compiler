package lexer

func NewProgram() program {

	program := program{}
	program.name = "program"
	state3 := State{
		name:    "3",
		isFinal: true}
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token) State {
			if t.lexeme == "END" {
				return state3
			}
			return invalidState()
		},
		isFinal: false}
	state1 := State{
		name:    "1",
		isFinal: false}
	state1.next = func(f *fsm, t Token) State {
		if t.tokenType == Number {
			return state2
		}
		b := NewBStatement()
		b.ConsumeToken(t)
		f.SetChildren(&b)
		return state1
	}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token) State {
			b := NewBStatement()
			b.ConsumeToken(t)
			f.SetChildren(&b)
			return state1
		},
		isFinal: false}
	program.initial = state0
	program.current = state0
	return program
}

type bstatement struct {
	fsm
}

func NewBStatement() bstatement {
	bstatement := bstatement{}
	bstatement.name = "bstatement"
	finalState := State{
		name:    "2",
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token) State {
			assignFSM := NewAssign()
			assignFSM.ConsumeToken(t)
			if assignFSM.InInvalidState() {
				return invalidState()
			}
			f.SetChildren(&assignFSM)
			return finalState
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token) State {
			if t.tokenType == Number {
				return state1
			}
			return invalidState()
		},
		isFinal: false}
	bstatement.initial = state0
	bstatement.current = state0
	bstatement.children = nil
	return bstatement
}

type assignFSM struct {
	fsm
}

func NewAssign() assignFSM {
	assignFSM := assignFSM{}
	assignFSM.name = "assign"
	state4 := State{
		name:    "4",
		isFinal: true}
	state3 := State{
		name: "3",
		next: func(f *fsm, t Token) State {
			e := NewExp()
			e.ConsumeToken(t)
			if e.InInvalidState() {
				return invalidState()
			}
			f.SetChildren(&e)
			return state4
		}}
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token) State {
			if t.tokenType == Equal {
				return state3
			}
			return invalidState()
		}}
	state1 := State{
		name: "asdasddasd",
		next: func(f *fsm, t Token) State {
			v := NewVar()
			v.ConsumeToken(t)
			if v.GetCurrent().name != invalidState().name {
				f.SetChildren(&v)
			} else {
				return invalidState()
			}
			return state2
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token) State {
			if t.lexeme == "LET" {
				return state1
			}
			return invalidState()
		},
		isFinal: false}
	assignFSM.initial = state0
	assignFSM.current = state0
	return assignFSM
}

type varFSM struct {
	fsm
}

func NewVar() varFSM { //FIXME
	varFSM := varFSM{}
	varFSM.name = "var"
	state1 := State{
		name:    "1",
		isFinal: true}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token) State {
			if t.tokenType == Identifier {
				return state1
			}
			return invalidState()
		},
		isFinal: false}
	varFSM.initial = state0
	varFSM.current = state0
	return varFSM
}

type expFSM struct {
	fsm
}

func NewExp() expFSM {
	expFSM := expFSM{}
	return expFSM
}
