package lexer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
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

	for c := range charc {
		fmt.Println(string(c.Type))
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
