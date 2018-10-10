package main

import (
	"fmt"
	"os"

	"github.com/mateusnakajo/basic-compiler/lexer"
)

func main() {
	args := os.Args[1:]
	switch {
	case len(args) > 1:
		fmt.Println("Usage: basic [script]")
	case len(args) == 1:
		lexer.RunLexer(args[0])
	case len(args) == 0:
		fmt.Println(">>")
	}
}
