package lexer

import (
	"fmt"
)

type Lexer struct {
	Source  string
	Tokens  []Token
	Start   int
	Current int
}

func (l Lexer) ScanTokens() []Token {
	for !l.isAtEnd() {
		l.Start = l.Current
		l.scanToken()
	}

	return nil
}

func (l Lexer) isAtEnd() bool {
	return l.Current >= len(l.Source)
}

func (l *Lexer) scanToken() {
	c := l.advance()
	switch c {
	case "(":
		fmt.Print("(\n")
	}
}

func (l *Lexer) advance() string {
	l.Current++
	return string(l.Source[l.Current-1])
}
