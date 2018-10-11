package lexer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/pkg/errors"
)

type Lexer struct {
	Source      string
	Tokens      []Token
	Start       int
	Current     int
	line        int
	lines       []string
	eventsQueue EventDrivenModule
}

type CategorizedChar struct {
	Char byte
	Type string
}

func MergeErrors(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(cs))

	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func WaitForPipeline(errs ...<-chan error) error {
	errc := MergeErrors(errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func lineListSource(ctx context.Context, program string) (
	<-chan string, <-chan error, error) {
	lines := strings.Split(program, "\n")
	if len(lines) == 0 {
		return nil, nil, errors.Errorf("no lines provided")
	}
	out := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for lineIndex, line := range lines {
			if line == "" {
				errc <- errors.Errorf("line %v is empty", lineIndex+1)
				return
			}
			select {
			case out <- line:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil
}

func charListSource(ctx context.Context, in <-chan string) (<-chan CategorizedChar, <-chan error, error) {
	errc := make(chan error, 1)
	out := make(chan CategorizedChar, 1)
	go func() {
		defer close(errc)
		defer close(out)
		for n := range in {
			for i := 0; i < len(n); i++ {
				select {
				case out <- categorizeChar(n[i]):
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out, errc, nil
}

func tokenListSource(ctx context.Context, in <-chan CategorizedChar) (<-chan Token, <-chan error, error) {
	errc := make(chan error, 1)
	out := make(chan Token, 1)
	go func() {
		defer close(errc)
		defer close(out)
		acc := ""
		token := Token{}
		for n := range in {
			c := rune(n.Char)
			if c == ' ' && acc == "" {
				continue
			}
			token = convertCharToToken(c) //ex: (
			if (Token{}) == token {
				if c == ' ' {
					//ex: DIM
					token = convertStringToToken(acc)
					emitToken(ctx, token, out)
					acc = ""
				} else {
					acc += string(c)
				}
			} else if acc != "" {
				//ex: F(i+1)
				token1 := Token{tokenType: Identifier, lexeme: acc}
				acc = ""
				emitToken(ctx, token1, out)
				emitToken(ctx, token, out)
			}
		}
	}()
	//10 DIM F(20)

	return out, errc, nil
}

func convertCharToToken(lexeme rune) Token {
	token := Token{}
	switch lexeme {
	case '(':
		token = Token{tokenType: LeftParen, lexeme: string(lexeme)}
	case ')':
		token = Token{tokenType: RightParen, lexeme: string(lexeme)}
	}
	return token
}

func convertStringToToken(lexeme string) Token {
	if _, err := strconv.Atoi(lexeme); err == nil {
		return Token{tokenType: Number, lexeme: lexeme}
	} else if unicode.IsLetter(rune(lexeme[0])) {
		return Token{tokenType: Identifier, lexeme: lexeme} //FIXME
	}
	return Token{}
}

func emitToken(ctx context.Context, token Token, out chan<- Token) {
	select {
	case out <- token:
	case <-ctx.Done():
		return
	}
}

func categorizeChar(in byte) CategorizedChar {
	categorizedChar := CategorizedChar{Char: in}
	switch char := rune(in); {
	case unicode.IsLetter(char):
		categorizedChar.Type = "Letter"
	case unicode.IsDigit(char):
		categorizedChar.Type = "Digit"
	case char == ' ':
		categorizedChar.Type = "Delimiter"
	case char == ':' || char == '=' || char == '+' || char == ';':
		categorizedChar.Type = "Special"
	}

	return categorizedChar
}

func RunLexer(filename string) error {
	program := readFile(filename)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var errcList []<-chan error
	linec, errc, err := lineListSource(ctx, program)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	charc, errc, err := charListSource(ctx, linec)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)

	tokenc, errc, err := tokenListSource(ctx, charc)
	for t := range tokenc {
		fmt.Println(t.lexeme)
	}
	return WaitForPipeline(errcList...)

}

func readFile(filename string) string {
	program, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(program)
}
