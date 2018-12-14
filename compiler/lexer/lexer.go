package lexer

import (
	"bufio"
	"os"
	"strconv"
	"unicode"

	"github.com/mateusnakajo/basic-compiler/compiler"
)

type CategorizedChar struct {
	Char rune
	Type string
}

type FileReader struct {
	compiler.EventDrivenModule
	filename string
	file     os.File
	scanner  *bufio.Scanner
}

func (a *FileReader) HandleEvent(event compiler.Event) {
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
	a.AddEvent(compiler.Event{"read", ""})
}

func (a *FileReader) ReadHandler(arg string) {
	if text := a.scanner.Scan(); text {
		a.AddEvent(compiler.Event{"read", ""})
		a.AddExternal(compiler.Event{"categorizeLineHandler", a.scanner.Text()})
	} else {
		a.AddEvent(compiler.Event{"close", ""})
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
	compiler.EventDrivenModule
}

func (a *AsciiCategorizer) HandleEvent(event compiler.Event) {
	handlers := map[string]func(string){
		"categorizeLineHandler": a.CategorizeLineHandler}
	handler := handlers[event.Name]
	handler(event.Arg.(string))
}

func (a *AsciiCategorizer) CategorizeLineHandler(line string) {
	for i := 0; i < len(line); i++ {
		a.AddExternal(compiler.Event{"tokenizeString", categorizeChar(line[i])})
	}
	a.AddExternal(compiler.Event{"tokenizeString", categorizeChar('\n')})
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
	compiler.EventDrivenModule
	acc string
}

func (t *TokenCategorizer) HandleEvent(event compiler.Event) {
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
			t.AddExternal(compiler.Event{"consumeToken", token})
		}
		t.acc = ""
	} else if char.Type == "Special" {
		token1 := t.createTokenFromAcc()
		if token1 != (Token{}) {
			t.AddExternal(compiler.Event{"consumeToken", token1})
		}
		token2 := t.createSpecialToken(string(char.Char))
		if token2 != (Token{}) {
			t.AddExternal(compiler.Event{"consumeToken", token2})
		}
		t.acc = ""
	}
}

func (t *TokenCategorizer) createTokenFromAcc() (token Token) {
	token = Token{}
	if t.acc != "" {
		if _, err := strconv.Atoi(t.acc); err == nil {
			token = Token{TokenType: Number, Lexeme: t.acc}
		} else {
			switch t.acc {
			case "GOTO":
				token = Token{TokenType: GoTo, Lexeme: "GOTO"}
			case "IF":
				token = Token{TokenType: If, Lexeme: "IF"}
			case "THEN":
				token = Token{TokenType: Then, Lexeme: "THEN"}
			case "FOR":
				token = Token{TokenType: For, Lexeme: "FOR"}
			case "TO":
				token = Token{TokenType: To, Lexeme: "TO"}
			case "STEP":
				token = Token{TokenType: Step, Lexeme: "STEP"}
			case "NEXT":
				token = Token{TokenType: Next, Lexeme: "NEXT"}
			case "DIM":
				token = Token{TokenType: Dim, Lexeme: "DIM"}
			case "DEF":
				token = Token{TokenType: Def, Lexeme: "DEF"}
			case "GOSUB":
				token = Token{TokenType: Gosub, Lexeme: "GOSUB"}
			case "RETURN":
				token = Token{TokenType: Return, Lexeme: "RETURN"}
			default:
				token = Token{TokenType: Identifier, Lexeme: t.acc}
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
			token = Token{TokenType: LeftParen, Lexeme: string(special)}
		case ")":
			token = Token{TokenType: RightParen, Lexeme: string(special)}
		case "=":
			token = Token{TokenType: Equal, Lexeme: string(special)}
		case "\n":
			token = Token{TokenType: EndOfLine, Lexeme: string(special)}
		case ",":
			token = Token{TokenType: Comma, Lexeme: string(special)}
		case "-":
			token = Token{TokenType: Minus, Lexeme: string(special)}
		case "+":
			token = Token{TokenType: Plus, Lexeme: string(special)}
		case "*":
			token = Token{TokenType: Star, Lexeme: string(special)}
		case "/":
			token = Token{TokenType: Slash, Lexeme: string(special)}
		case ">":
			token = Token{TokenType: Greater, Lexeme: string(special)}
		case "<":
			token = Token{TokenType: Less, Lexeme: string(special)}
		}
	}
	return token
}

func (t *TokenCategorizer) createSpecialTokenWithMoreThenOneChar(special string) (token Token) {
	token = Token{}
	if special == ">" && t.LookAhead().Arg.(CategorizedChar).Char == '=' {
		token = Token{TokenType: GreaterEqual, Lexeme: ">="}
	} //
	if token != (Token{}) {
		t.PopEvent()
	}
	return token
}
