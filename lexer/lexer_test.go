package lexer

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tokens, _ := RunLexer("10 DIM F(20)")
	expected := []Token{
		Token{Number, "10"},
		Token{Keyword, "DIM"},
		Token{Identifier, "F"},
		Token{LeftParen, "("},
		Token{Number, "20"},
		Token{RightParen, ")"},
		Token{EndOfLine, "\n"}}
	for i := range tokens {
		if tokens[i] != expected[i] {
			t.Errorf("Test failed: {%v} expected, {%v} received", expected[i], tokens[i])
		}
	}
}
