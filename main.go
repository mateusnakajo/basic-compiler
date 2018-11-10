package main

import (
	"io/ioutil"
	"log"

	"github.com/mateusnakajo/basic-compiler/lexer"
)

func main() { /*
		args := os.Args[1:]
		switch {
		case len(args) > 1:
			fmt.Println("Usage: basic [script]")
		case len(args) == 1:
			lexer.RunLexer(readFile(args[0]))
		case len(args) == 0:
			fmt.Println(">>")
		}*/
	f := lexer.FileReader{}
	a := lexer.AsciiCategorizer{}

	f.AddEvent(lexer.Event{"open", "sample-program/quicksort.bas"})
	f.AddExternal = a.AddEvent
	for !f.IsEmpty() {
		event := f.PopEvent()
		f.HandleEvent(event)
	}
	for !a.IsEmpty() {
		event := a.PopEvent()
		a.HandleEvent(event)
	}
}

func readFile(filename string) string {
	program, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(program)
}
