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
		program := readFile(args[0])
		l := lexer.Lexer{Source: program}
		l.ScanTokens()
	case len(args) == 0:
		fmt.Println(">>")
	}
}

func readFile(fileName string) string {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}
