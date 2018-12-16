package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mateusnakajo/basic-compiler/compiler"
	lexer "github.com/mateusnakajo/basic-compiler/compiler/lexer"
	"github.com/mateusnakajo/basic-compiler/compiler/semantic"
	syntactic "github.com/mateusnakajo/basic-compiler/compiler/syntactic"
)

func main() {
	//args := os.Args[1:]

	f := lexer.FileReader{}
	a := lexer.AsciiCategorizer{}
	t := lexer.TokenCategorizer{}
	s := syntactic.NewSyntaticAnalyser()
	semantic := semantic.NewSemantic()
	f.AddEvent(compiler.Event{Name: "open", Arg: "sample-program/test-array.bas"})
	//f.AddEvent(compiler.Event{Name: "open", Arg: args[0]})
	f.AddExternal = a.AddEvent
	a.AddExternal = t.AddEvent
	t.AddExternal = s.AddEvent
	s.AddExternal = semantic.AddEvent
	semantic.AddExternal = s.AddEvent

	for !f.IsEmpty() {
		event := f.PopEvent()
		f.HandleEvent(event)
	}
	for !a.IsEmpty() {
		event := a.PopEvent()
		a.HandleEvent(event)
	}
	for !t.IsEmpty() {
		event := t.PopEvent()
		t.HandleEvent(event)
	}
	semantic.TokenEvents = s.Events
	for !s.IsEmpty() {
		event := s.PopEvent()
		s.HandleEvent(event)
	}
	semantic.IndexOfLine = s.IndexOfLine
	for !semantic.IsEmpty() {
		event := semantic.PopEvent()
		semantic.HandleEvent(event)
		fmt.Print(event)
	}
	for semantic.Rerun {
		semantic.Rerun = false
		s.Events = semantic.NewTokenEvents

		for !s.IsEmpty() {
			event := s.PopEvent()
			s.HandleEvent(event)
		}
		for !semantic.IsEmpty() {
			event := semantic.PopEvent()
			semantic.HandleEvent(event)
			//fmt.Print(event)
		}
	}
	// fmt.Println(semantic.TokenEvents)

	//s := lexer.NewSyntaticAnalyser()
	//s.HandleEvent(lexer.Event{"consumeToken", lexer.Token{}})
}

func readFile(filename string) string {
	program, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(program)
}
