package lexer

func NewProgram() program {

	program := program{}
	program.name = "program"
	state1 := State{
		name:    "1",
		isFinal: false}
	state1.next = func(f *fsm, t Token, s *Stack) State {
		b := NewBStatement()
		b.ConsumeToken(t, s)
		s.AddFSM(&b)
		return state1
	}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			b := NewBStatement()
			b.ConsumeToken(t, s)
			s.AddFSM(&b)
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
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.lexeme == "END" {
				return finalState //FIXME
			}
			assignFSM := NewAssign()
			assignFSM.ConsumeToken(t, s)
			if !assignFSM.InInvalidState() {
				s.AddFSM(&assignFSM)
				return finalState
			}

			nextFSM := NewPredef()
			nextFSM.ConsumeToken(t, s)
			if !nextFSM.InInvalidState() {
				s.AddFSM(&nextFSM)
				return finalState
			}

			readFSM := NewRead()
			readFSM.ConsumeToken(t, s)
			if !readFSM.InInvalidState() {
				s.AddFSM(&readFSM)
				return finalState
			}

			printFSM := NewPrint()
			printFSM.ConsumeToken(t, s)
			if !printFSM.InInvalidState() {
				s.AddFSM(&printFSM)
				return finalState
			}

			gotoFSM := NewGoto()
			gotoFSM.ConsumeToken(t, s)
			if !gotoFSM.InInvalidState() {
				s.AddFSM(&gotoFSM)
				return finalState
			}

			ifFSM := NewIf()
			ifFSM.ConsumeToken(t, s)
			if !ifFSM.InInvalidState() {
				s.AddFSM(&ifFSM)
				return finalState
			}

			return invalidState()
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == Number {
				return state1
			}
			return invalidState()
		},
		isFinal: false}
	bstatement.initial = state0
	bstatement.current = state0
	return bstatement
}

type assignFSM struct {
	fsm
}

func NewAssign() assignFSM {
	assignFSM := assignFSM{}
	assignFSM.name = "assign"
	state4 := State{
		name: "4",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state3 := State{
		name: "3",
		next: func(f *fsm, t Token, s *Stack) State {
			e := NewExp()
			e.ConsumeToken(t, s)
			if e.InInvalidState() {
				return invalidState()
			}
			s.AddFSM(&e)
			return state4
		},
		isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == Equal {
				return state3
			}
			return invalidState()
		}}
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			v := NewVar()
			v.ConsumeToken(t, s)
			if !v.InInvalidState() {
				s.AddFSM(&v)
			} else {
				return invalidState()
			}
			return state2
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
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
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
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
	expFSM.name = "exp"

	state1 := State{
		name:    "1",
		isFinal: true}
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			ebFSM := NewEB()
			ebFSM.ConsumeToken(t, s)
			if ebFSM.GetCurrent().name != invalidState().name {
				s.AddFSM(&ebFSM)
				return state1
			}
			return invalidState()
		}, isFinal: false}

	state1.next = func(f *fsm, t Token, s *Stack) State {
		if t.tokenType == Plus || t.tokenType == Minus || t.tokenType == Star || t.tokenType == Slash {
			return state2
		}
		return invalidState()
	}

	state0 := State{
		name:    "0",
		isFinal: false}
	state0.next = func(f *fsm, t Token, s *Stack) State {
		if t.tokenType == Plus || t.tokenType == Minus {
			return state0
		}
		ebFSM := NewEB()
		ebFSM.ConsumeToken(t, s)
		if ebFSM.GetCurrent().name != invalidState().name {
			s.AddFSM(&ebFSM)
			return state1
		}
		return invalidState()
	}
	expFSM.initial = state0
	expFSM.current = state0
	return expFSM
}

type ebFSM struct {
	fsm
}

func NewEB() ebFSM {
	ebFSM := ebFSM{}
	ebFSM.name = "eb"
	state5 := State{
		name: "5",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state4 := State{
		name: "4",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == RightParen {
				return state5
			}
			return invalidState()
		},
		isFinal: false}
	state3 := State{
		name: "3",
		next: func(f *fsm, t Token, s *Stack) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s)
			if expFSM.GetCurrent().name != invalidState().name {
				s.AddFSM(&expFSM)
			} else {
				return invalidState()
			}
			return state4
		},
		isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == LeftParen {
				return state3
			}
			return invalidState()
		},
		isFinal: false}
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == Identifier {
				return state2
			}
			return invalidState()
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.lexeme == "FN" {
				return state1
			}
			if t.tokenType == LeftParen {
				return state3
			}
			if t.tokenType == Number || t.tokenType == Identifier {
				return state5
			}
			return invalidState()
		},
		isFinal: false}
	ebFSM.initial = state0
	ebFSM.current = state0
	return ebFSM
}

type predefFSM struct {
	fsm
}

func NewPredef() predefFSM {
	predefFSM := predefFSM{}
	predefFSM.name = "predef"
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.lexeme == "SIN" { //FIXME
				return state1
			}
			return invalidState()
		},
		isFinal: false}
	predefFSM.initial = state0
	predefFSM.current = state0
	return predefFSM
}

type readFSM struct {
	fsm
}

func NewRead() readFSM {
	readFSM := readFSM{}
	readFSM.name = "read"
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true} //FIXME
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			varFSM := NewVar()
			varFSM.ConsumeToken(t, s)
			if !varFSM.InInvalidState() {
				s.AddFSM(&varFSM)
				return state2
			}
			return invalidState()
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.lexeme == "READ" { //FIXME
				return state1
			}
			return invalidState()
		},
		isFinal: false}
	readFSM.initial = state0
	readFSM.current = state0
	return readFSM
}

type dataFSM struct {
	fsm
}

func NewData() dataFSM {
	dataFSM := dataFSM{}
	dataFSM.name = "data"
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true} //FIXME
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			// varFSM := NewVar()
			// varFSM.ConsumeToken(t, s)
			// if !varFSM.InInvalidState() {
			// 	s.AddFSM(&varFSM)
			// 	return state2
			// } FIXME
			return state2
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.lexeme == "DATA" { //FIXME
				return state1
			}
			return invalidState()
		},
		isFinal: false}
	dataFSM.initial = state0
	dataFSM.current = state0
	return dataFSM
}

type printFSM struct {
	fsm
}

func NewPrint() printFSM {
	printFSM := printFSM{}
	printFSM.name = "print"
	state4 := State{
		name: "4",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state3 := State{
		name: "3",
		next: func(f *fsm, t Token, s *Stack) State {
			pitemFSM := NewPitem()
			pitemFSM.ConsumeToken(t, s)
			if !pitemFSM.InInvalidState() {
				s.AddFSM(&pitemFSM)
				return state1
			}
			return state1
		},
		isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == Comma {
				return state3
			}
			return invalidState()
		},
		isFinal: false} //FIXME
	state1.next = func(f *fsm, t Token, s *Stack) State {
		if t.tokenType == Comma {
			return state4
		}
		pitemFSM := NewPitem()
		pitemFSM.ConsumeToken(t, s)
		if !pitemFSM.InInvalidState() {
			s.AddFSM(&pitemFSM)
			return state2
		}
		return state2
	}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.lexeme == "PRINT" { //FIXME
				return state1
			}
			return invalidState()
		},
		isFinal: false}
	printFSM.initial = state0
	printFSM.current = state0
	return printFSM
}

type pitemFSM struct {
	fsm
}

func NewPitem() pitemFSM {
	pitemFSM := pitemFSM{}
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s)
			if !expFSM.InInvalidState() {
				return state2
			}
			return invalidState()
		}, isFinal: true}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == String {
				return state1
			}
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s)
			if !expFSM.InInvalidState() {
				return state2
			}
			return invalidState()
		}, isFinal: false}
	pitemFSM.initial = state0
	pitemFSM.current = state0
	return pitemFSM
}

type gotoFSM struct {
	fsm
}

func NewGoto() gotoFSM {
	gotoFSM := gotoFSM{}
	gotoFSM.name = "goto"
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == Number {
				return state2
			}
			return invalidState()
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == GoTo {
				return state1
			}
			return invalidState()
		}, isFinal: false}
	gotoFSM.initial = state0
	gotoFSM.current = state0
	return gotoFSM
}

type ifFSM struct {
	fsm
}

func NewIf() ifFSM {
	ifFSM := ifFSM{}
	ifFSM.name = "if"
	state6 := State{
		name: "6",
		next: func(f *fsm, t Token, s *Stack) State {
			return invalidState()
		},
		isFinal: true}
	state5 := State{
		name: "5",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == Number {
				return state6
			}
			return invalidState()
		},
		isFinal: false}
	state4 := State{
		name: "4",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == Then {
				return state5
			}
			return invalidState()
		}, isFinal: false}
	state3 := State{
		name: "3",
		next: func(f *fsm, t Token, s *Stack) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s)
			if !expFSM.InInvalidState() {
				s.AddFSM(&expFSM)
			}
			return state4
		}, isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == Greater || t.tokenType == GreaterEqual ||
				t.tokenType == Different || t.tokenType == Less ||
				t.tokenType == LessEqual || t.tokenType == Equal {
				return state3
			}
			return invalidState()
		}, isFinal: false}
	state1 := State{
		name: "1",
		next: func(f *fsm, t Token, s *Stack) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s)
			if !expFSM.InInvalidState() {
				s.AddFSM(&expFSM)
			}
			return state2
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t Token, s *Stack) State {
			if t.tokenType == If {
				return state1
			}
			return invalidState()
		}, isFinal: false}
	ifFSM.initial = state0
	ifFSM.current = state0
	return ifFSM
}
