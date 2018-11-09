package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mateusnakajo/basic-compiler/lexer"
)

func main() {
	args := os.Args[1:]
	switch {
	case len(args) > 1:
		fmt.Println("Usage: basic [script]")
	case len(args) == 1:
		lexer.RunLexer(readFile(args[0]))
	case len(args) == 0:
		fmt.Println(">>")
	}
}

func readFile(filename string) string {
	program, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(program)
}
