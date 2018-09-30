package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	Source string
	Tokens []Token
}

func (l Lexer) ScanTokens() []Token {
	lines := strings.Split(l.Source, "\n")
	for _, line := range lines {
		for _, char := range line {
			fmt.Print(string(char))
		}
	}
	return nil
}
