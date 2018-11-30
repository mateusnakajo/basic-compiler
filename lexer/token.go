package lexer

type TokenType int

const (
	// delimitors
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot

	// operators
	Equal
	Minus
	Plus
	Star
	Slash
	E10

	// comparators
	Greater
	Less
	Different
	GreaterEqual
	EqualEqual
	LessEqual

	//keywords
	Let
	Fn
	Read
	Data
	Print
	GoTo
	If
	For
	To
	Step
	Next
	Dim
	Def
	Gosub
	Return
	Rem
	Then

	// literals
	Identifier
	String
	Number

	// control
	EndOfLine

	Keyword
)

type Token struct {
	tokenType TokenType
	lexeme    string
}
