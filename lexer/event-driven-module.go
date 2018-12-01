package lexer

import (
	"bufio"
	"os"
	"strconv"
	"unicode"
)

type EventDrivenModule struct {
	Events      []Event
	AddExternal func(Event)
}

func (e *EventDrivenModule) AddEvent(newEvent Event) {
	e.Events = append(e.Events, newEvent)
}

func (e *EventDrivenModule) PopEvent() Event {
	firstEvent := e.Events[0]
	e.Events = e.Events[1:]
	return firstEvent
}

func (e *EventDrivenModule) LookAhead() Event {
	return e.Events[0]
}

func (e *EventDrivenModule) IsEmpty() bool {
	return len(e.Events) == 0
}

type Event struct {
	Name string
	Arg  interface{}
}

type FileReader struct {
	EventDrivenModule
	filename string
	file     os.File
	scanner  *bufio.Scanner
}

func (a *FileReader) HandleEvent(event Event) {
	handlers := map[string]func(string){
		"open":  a.OpenHandler,
		"read":  a.ReadHandler,
		"close": a.CloseHandler}
	handler := handlers[event.Name]
	handler(event.Arg.(string))
}

func (a *FileReader) OpenHandler(filename string) {
	f, err := os.Open(filename)
	check(err)

	a.scanner = bufio.NewScanner(f)
	a.AddEvent(Event{"read", ""})
}

func (a *FileReader) ReadHandler(arg string) {
	if text := a.scanner.Scan(); text {
		a.AddEvent(Event{"read", ""})
		a.AddExternal(Event{"categorizeLineHandler", a.scanner.Text()})
	} else {
		a.AddEvent(Event{"close", ""})
	}
}

func (a *FileReader) CloseHandler(arg string) {
	a.file.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type AsciiCategorizer struct {
	EventDrivenModule
}

func (a *AsciiCategorizer) HandleEvent(event Event) {
	handlers := map[string]func(string){
		"categorizeLineHandler": a.CategorizeLineHandler}
	handler := handlers[event.Name]
	handler(event.Arg.(string))
}

func (a *AsciiCategorizer) CategorizeLineHandler(line string) {
	for i := 0; i < len(line); i++ {
		a.AddExternal(Event{"tokenizeString", categorizeChar(line[i])})
	}
	a.AddExternal(Event{"tokenizeString", categorizeChar('\n')})
}

func categorizeChar(in byte) CategorizedChar {
	categorizedChar := CategorizedChar{Char: rune(in)}
	switch char := rune(in); {
	case unicode.IsLetter(char):
		categorizedChar.Type = "Letter"
	case unicode.IsDigit(char):
		categorizedChar.Type = "Digit"
	case char == ' ' || char == '\n':
		categorizedChar.Type = "Delimiter"
	default:
		categorizedChar.Type = "Special"
	}

	return categorizedChar
}

type TokenCategorizer struct {
	EventDrivenModule
	acc string
}

func (t *TokenCategorizer) HandleEvent(event Event) {
	handlers := map[string]func(CategorizedChar){
		"tokenizeString": t.tokenizeString}
	handler := handlers[event.Name]
	handler(event.Arg.(CategorizedChar))
}

func (t *TokenCategorizer) tokenizeString(char CategorizedChar) {
	//fmt.Println(string(char.Type), string(char.Char))
	if char.Type == "Letter" || char.Type == "Digit" {
		t.acc += string(char.Char)
	} else if char.Type == "Delimiter" {
		token := t.createTokenFromAcc()
		if token != (Token{}) {
			t.AddExternal(Event{"consumeToken", token})
		}
		t.acc = ""
	} else if char.Type == "Special" {
		token1 := t.createTokenFromAcc()
		if token1 != (Token{}) {
			t.AddExternal(Event{"consumeToken", token1})
		}
		token2 := t.createSpecialToken(string(char.Char))
		if token2 != (Token{}) {
			t.AddExternal(Event{"consumeToken", token2})
		}
		t.acc = ""
	}
}

func (t *TokenCategorizer) createTokenFromAcc() (token Token) {
	token = Token{}
	if t.acc != "" {
		if _, err := strconv.Atoi(t.acc); err == nil {
			token = Token{tokenType: Number, lexeme: t.acc}
		} else {
			switch t.acc {
			case "GOTO":
				token = Token{tokenType: GoTo, lexeme: "GOTO"}
			case "IF":
				token = Token{tokenType: If, lexeme: "IF"}
			case "THEN":
				token = Token{tokenType: Then, lexeme: "THEN"}
			case "FOR":
				token = Token{tokenType: For, lexeme: "FOR"}
			case "TO":
				token = Token{tokenType: To, lexeme: "TO"}
			case "STEP":
				token = Token{tokenType: Step, lexeme: "STEP"}
			case "NEXT":
				token = Token{tokenType: Next, lexeme: "NEXT"}
			case "DIM":
				token = Token{tokenType: Dim, lexeme: "DIM"}
			case "DEF":
				token = Token{tokenType: Def, lexeme: "DEF"}
			case "GOSUB":
				token = Token{tokenType: Gosub, lexeme: "GOSUB"}
			case "RETURN":
				token = Token{tokenType: Return, lexeme: "RETURN"}
			default:
				token = Token{tokenType: Identifier, lexeme: t.acc}
			}
		}
	}
	return token
}

func (t *TokenCategorizer) createSpecialToken(special string) (token Token) {
	token = t.createSpecialTokenWithMoreThenOneChar(special)
	if token == (Token{}) {
		switch special {
		case "(":
			token = Token{tokenType: LeftParen, lexeme: string(special)}
		case ")":
			token = Token{tokenType: RightParen, lexeme: string(special)}
		case "=":
			token = Token{tokenType: Equal, lexeme: string(special)}
		case "\n":
			token = Token{tokenType: EndOfLine, lexeme: string(special)}
		case ",":
			token = Token{tokenType: Comma, lexeme: string(special)}
		case "-":
			token = Token{tokenType: Minus, lexeme: string(special)}
		case "+":
			token = Token{tokenType: Plus, lexeme: string(special)}
		case "*":
			token = Token{tokenType: Star, lexeme: string(special)}
		case "/":
			token = Token{tokenType: Slash, lexeme: string(special)}
		case ">":
			token = Token{tokenType: Greater, lexeme: string(special)}
		case "<":
			token = Token{tokenType: Less, lexeme: string(special)}
		}
	}
	return token
}

func (t *TokenCategorizer) createSpecialTokenWithMoreThenOneChar(special string) (token Token) {
	token = Token{}
	if special == ">" && t.LookAhead().Arg.(CategorizedChar).Char == '=' {
		token = Token{tokenType: GreaterEqual, lexeme: ">="}
	} //
	if token != (Token{}) {
		t.PopEvent()
	}
	return token
}
