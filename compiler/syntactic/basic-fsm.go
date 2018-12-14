package syntactic

import (
	"fmt"

	lexer "github.com/mateusnakajo/basic-compiler/compiler/lexer"
	semantic "github.com/mateusnakajo/basic-compiler/compiler/semantic"
)

type program struct {
	fsm
}

func NewProgram() program {

	program := program{}
	program.name = "program"
	state1 := State{
		name:    "1",
		isFinal: false}
	state1.next = func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
		b := NewBStatement()
		b.ConsumeToken(t, s, semantic)
		s.AddFSM(&b)
		return state1
	}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			b := NewBStatement()
			b.ConsumeToken(t, s, semantic)
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.Lexeme == "END" {
				return finalState //FIXME
			}
			assignFSM := NewAssign()
			assignFSM.ConsumeToken(t, s, semantic)
			if !assignFSM.InInvalidState() {
				s.AddFSM(&assignFSM)
				return finalState
			}

			predefFSM := NewPredef()
			predefFSM.ConsumeToken(t, s, semantic)
			if !predefFSM.InInvalidState() {
				s.AddFSM(&predefFSM)
				return finalState
			}

			readFSM := NewRead()
			readFSM.ConsumeToken(t, s, semantic)
			if !readFSM.InInvalidState() {
				s.AddFSM(&readFSM)
				return finalState
			}

			printFSM := NewPrint()
			printFSM.ConsumeToken(t, s, semantic)
			if !printFSM.InInvalidState() {
				s.AddFSM(&printFSM)
				return finalState
			}

			gotoFSM := NewGoto()
			gotoFSM.ConsumeToken(t, s, semantic)
			if !gotoFSM.InInvalidState() {
				s.AddFSM(&gotoFSM)
				return finalState
			}

			ifFSM := NewIf()
			ifFSM.ConsumeToken(t, s, semantic)
			if !ifFSM.InInvalidState() {
				s.AddFSM(&ifFSM)
				return finalState
			}

			forFSM := NewFor()
			forFSM.ConsumeToken(t, s, semantic)
			if !forFSM.InInvalidState() {
				s.AddFSM(&forFSM)
				return finalState
			}

			nextFSM := NewNext()
			nextFSM.ConsumeToken(t, s, semantic)
			if !nextFSM.InInvalidState() {
				s.AddFSM(&nextFSM)
				return finalState
			}

			dimFSM := NewDim()
			dimFSM.ConsumeToken(t, s, semantic)
			if !dimFSM.InInvalidState() {
				s.AddFSM(&dimFSM)
				return finalState
			}

			defFSM := NewDef()
			defFSM.ConsumeToken(t, s, semantic)
			if !defFSM.InInvalidState() {
				s.AddFSM(&defFSM)
				return finalState
			}

			gosubFSM := NewGosub()
			gosubFSM.ConsumeToken(t, s, semantic)
			if !gosubFSM.InInvalidState() {
				s.AddFSM(&gosubFSM)
				return finalState
			}

			returnFSM := NewReturn()
			returnFSM.ConsumeToken(t, s, semantic)
			if !returnFSM.InInvalidState() {
				s.AddFSM(&returnFSM)
				return finalState
			}

			return invalidState()
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Number {
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			semantic.Evaluate()
			semantic.DataFloat[semantic.PopString()] = semantic.PopFloat()
			return invalidState()
		},
		isFinal: true}
	state3 := State{
		name: "3",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			e := NewExp()
			e.ConsumeToken(t, s, semantic)
			if e.InInvalidState() {
				return invalidState()
			}
			s.AddFSM(&e)

			return state4
		},
		isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Equal {
				return state3
			}
			return invalidState()
		}}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			v := NewVar()
			v.ConsumeToken(t, s, semantic)
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.Lexeme == "LET" {
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Identifier {
				semantic.SaveString(t.Lexeme)
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			ebFSM := NewEB()
			ebFSM.ConsumeToken(t, s, semantic)
			if ebFSM.GetCurrent().name != invalidState().name {
				s.AddFSM(&ebFSM)
				return state1
			}
			return invalidState()
		}, isFinal: false}

	state1.next = func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
		if t.TokenType == lexer.Plus || t.TokenType == lexer.Minus || t.TokenType == lexer.Star || t.TokenType == lexer.Slash {
			semantic.Expression += t.Lexeme
			return state2
		}

		return invalidState()
	}

	state0 := State{
		name:    "0",
		isFinal: false}
	state0.next = func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
		if t.TokenType == lexer.Plus || t.TokenType == lexer.Minus {
			return state0
		}
		ebFSM := NewEB()
		ebFSM.ConsumeToken(t, s, semantic)
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state4 := State{
		name: "4",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			semantic.Expression += t.Lexeme
			if t.TokenType == lexer.RightParen {
				return state5
			}
			return invalidState()
		},
		isFinal: false}
	state3 := State{
		name: "3",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			semantic.Expression += t.Lexeme
			if t.TokenType == lexer.LeftParen {
				return state3
			}
			return invalidState()
		},
		isFinal: false}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			semantic.Expression += t.Lexeme
			if t.TokenType == lexer.Identifier {
				return state2
			}
			return invalidState()
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.Lexeme == "FN" {
				return state1
			}
			if t.TokenType == lexer.LeftParen {
				semantic.Expression += t.Lexeme
				return state3
			}
			if t.TokenType == lexer.Number {
				semantic.Expression += t.Lexeme
				return state5
			}
			if t.TokenType == lexer.Identifier {
				semantic.Expression += fmt.Sprintf("%f", semantic.DataFloat[t.Lexeme])
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.Lexeme == "SIN" { //FIXME
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true} //FIXME
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			varFSM := NewVar()
			varFSM.ConsumeToken(t, s, semantic)
			if !varFSM.InInvalidState() {
				s.AddFSM(&varFSM)
				return state2
			}
			return invalidState()
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.Lexeme == "READ" { //FIXME
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true} //FIXME
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			// varFSM := NewVar()
			// varFSM.ConsumeToken(t, s, semantic)
			// if !varFSM.InInvalidState() {
			// 	s.AddFSM(&varFSM)
			// 	return state2
			// } FIXME
			return state2
		},
		isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.Lexeme == "DATA" { //FIXME
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state3 := State{
		name: "3",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			pitemFSM := NewPitem()
			pitemFSM.ConsumeToken(t, s, semantic)
			if !pitemFSM.InInvalidState() {
				s.AddFSM(&pitemFSM)
				return state1
			}
			return state1
		},
		isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Comma {
				return state3
			}
			return invalidState()
		},
		isFinal: false} //FIXME
	state1.next = func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
		if t.TokenType == lexer.Comma {
			return state4
		}
		pitemFSM := NewPitem()
		pitemFSM.ConsumeToken(t, s, semantic)
		if !pitemFSM.InInvalidState() {
			s.AddFSM(&pitemFSM)
			return state2
		}
		return state2
	}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.Lexeme == "PRINT" { //FIXME
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
			if !expFSM.InInvalidState() {
				return state2
			}
			return invalidState()
		}, isFinal: true}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.String {
				return state1
			}
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Number {
				return state2
			}
			return invalidState()
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.GoTo {
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
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		},
		isFinal: true}
	state5 := State{
		name: "5",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Number {
				return state6
			}
			return invalidState()
		},
		isFinal: false}
	state4 := State{
		name: "4",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Then {
				return state5
			}
			return invalidState()
		}, isFinal: false}
	state3 := State{
		name: "3",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
			if !expFSM.InInvalidState() {
				s.AddFSM(&expFSM)
			}
			return state4
		}, isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Greater || t.TokenType == lexer.GreaterEqual ||
				t.TokenType == lexer.Different || t.TokenType == lexer.Less ||
				t.TokenType == lexer.LessEqual || t.TokenType == lexer.Equal {
				return state3
			}
			return invalidState()
		}, isFinal: false}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
			if !expFSM.InInvalidState() {
				s.AddFSM(&expFSM)
			}
			return state2
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.If {
				return state1
			}
			return invalidState()
		}, isFinal: false}
	ifFSM.initial = state0
	ifFSM.current = state0
	return ifFSM
}

type forFSM struct {
	fsm
}

func NewFor() forFSM {
	forFSM := forFSM{}
	forFSM.name = "for"
	state8 := State{
		name: "8",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		}, isFinal: true}
	state7 := State{
		name: "7",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
			if !expFSM.InInvalidState() {
				s.AddFSM(&expFSM)
			}
			return state8
		}, isFinal: false}
	state6 := State{
		name: "6",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Step {
				return state7
			}
			return invalidState()
		},
		isFinal: true}
	state5 := State{
		name: "5",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
			if !expFSM.InInvalidState() {
				s.AddFSM(&expFSM)
			}
			return state6
		}, isFinal: false}
	state4 := State{
		name: "4",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.To {
				return state5
			}
			return invalidState()
		}, isFinal: false}
	state3 := State{
		name: "3",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
			if !expFSM.InInvalidState() {
				s.AddFSM(&expFSM)
			}
			return state4
		}, isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Equal {
				return state3
			}
			return invalidState()
		}, isFinal: false}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Identifier {
				return state2
			}
			return invalidState()
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.For {
				return state1
			}
			return invalidState()
		}, isFinal: false}
	forFSM.initial = state0
	forFSM.current = state0
	return forFSM
}

type nextFSM struct {
	fsm
}

func NewNext() nextFSM {
	nextFSM := nextFSM{}
	nextFSM.name = "next"
	state2 := State{
		name: "2",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		}, isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Identifier {
				return state2
			}
			return invalidState()
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Next {
				return state1
			}
			return invalidState()
		}, isFinal: false}

	nextFSM.initial = state0
	nextFSM.current = state0
	return nextFSM
}

type dimFSM struct {
	fsm
}

func NewDim() dimFSM {
	dimFSM := dimFSM{}
	dimFSM.name = "dim"
	state1 := State{
		name:    "1",
		isFinal: false}
	state6 := State{
		name: "6",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Comma {
				return state1
			}
			return invalidState()
		}, isFinal: false}
	state4 := State{
		name:    "4",
		isFinal: false}
	state5 := State{
		name: "5",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Identifier {
				return state4
			}
			return invalidState()
		}, isFinal: false}
	state4.next = func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
		if t.TokenType == lexer.Comma {
			return state5
		}
		if t.TokenType == lexer.RightParen {
			return state6
		}
		return invalidState()
	}
	state3 := State{
		name: "3",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Number {
				return state4
			}
			return invalidState()
		}, isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.LeftParen {
				return state3
			}
			return invalidState()
		}, isFinal: false}
	state1.next = func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
		if t.TokenType == lexer.Identifier { //FIXME: tem que ser uma letra só
			return state2
		}
		return invalidState()
	}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Dim {
				return state1
			}
			return invalidState()
		}, isFinal: false}

	dimFSM.initial = state0
	dimFSM.current = state0
	return dimFSM
}

type defFSM struct {
	fsm
}

func NewDef() defFSM {
	defFSM := defFSM{}
	defFSM.name = "def"
	state7 := State{
		name: "7",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		}, isFinal: true}
	state6 := State{
		name: "6",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			expFSM := NewExp()
			expFSM.ConsumeToken(t, s, semantic)
			if !expFSM.InInvalidState() {
				s.AddFSM(&expFSM)
			}
			return state7
		}, isFinal: false}
	state5 := State{
		name: "5",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Equal {
				return state6
			}
			return invalidState()
		}, isFinal: false}
	state4 := State{
		name: "4",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.RightParen {
				return state5
			}
			return invalidState()
		}, isFinal: false}
	state3 := State{
		name: "3",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Identifier {
				return state4
			}
			return invalidState()
		}, isFinal: false}
	state2 := State{
		name: "2",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.LeftParen {
				return state3
			}
			return invalidState()
		}, isFinal: false}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Identifier { //fnf
				return state2
			}
			return invalidState()
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Def {
				return state1
			}
			return invalidState()
		}, isFinal: false}

	defFSM.initial = state0
	defFSM.current = state0
	return defFSM
}

type gosubFSM struct {
	fsm
}

func NewGosub() gosubFSM {
	gosubFSM := gosubFSM{}
	gosubFSM.name = "gosub"
	state2 := State{
		name: "2",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		}, isFinal: true}
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Number {
				return state2
			}
			return invalidState()
		}, isFinal: false}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Gosub {
				return state1
			}
			return invalidState()
		}, isFinal: false}

	gosubFSM.initial = state0
	gosubFSM.current = state0
	return gosubFSM
}

type returnFSM struct {
	fsm
}

func NewReturn() returnFSM {
	returnFSM := returnFSM{}
	returnFSM.name = "return"
	state1 := State{
		name: "1",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			return invalidState()
		}, isFinal: true}
	state0 := State{
		name: "0",
		next: func(f *fsm, t lexer.Token, s *Stack, semantic *semantic.Semantic) State {
			if t.TokenType == lexer.Return {
				return state1
			}
			return invalidState()
		}, isFinal: false}
	returnFSM.initial = state0
	returnFSM.current = state0
	return returnFSM
}
