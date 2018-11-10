package lexer

import (
	"bufio"
	"fmt"
	"os"
)

type EventDrivenModule struct {
	Events      []Event
	Handlers    map[string]func(string)
	AddExternal func(Event)
}

func (e *EventDrivenModule) AddEvent(newEvent Event) {
	e.Events = append(e.Events, newEvent)
}

func (e *EventDrivenModule) PopEvent() Event {
	firstEvent := e.Events[0]
	e.Events = e.Events[1:]
	return firstEvent
}

func (e *EventDrivenModule) IsEmpty() bool {
	return len(e.Events) == 0
}

type Event struct {
	Name string
	Arg  string
}

type FileReader struct {
	EventDrivenModule
	filename string
	file     os.File
	scanner  *bufio.Scanner
}

func (a *FileReader) HandleEvent(event Event) {
	handlers := map[string]func(string){
		"open":  a.OpenHandler,
		"read":  a.ReadHandler,
		"close": a.CloseHandler}
	handler := handlers[event.Name]
	handler(event.Arg)
}

func (a *FileReader) OpenHandler(filename string) {
	f, err := os.Open(filename)
	check(err)

	a.scanner = bufio.NewScanner(f)
	a.AddEvent(Event{"read", ""})
}

func (a *FileReader) ReadHandler(arg string) {
	a.scanner.Scan()
	if text := a.scanner.Scan(); text {
		a.AddExternal(Event{"categorizeLineHandler", a.scanner.Text()})
		a.AddEvent(Event{"read", ""})
	} else {
		a.AddEvent(Event{"close", ""})
	}
}

func (a *FileReader) CloseHandler(arg string) {
	fmt.Println("closed")
	a.file.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type AsciiCategorizer struct {
	EventDrivenModule
}

func (a *AsciiCategorizer) HandleEvent(event Event) {
	handlers := map[string]func(string){
		"categorizeLineHandler": a.CategorizeLineHandler}
	handler := handlers[event.Name]
	handler(event.Arg)
}

func (a *AsciiCategorizer) CategorizeLineHandler(line string) {
	fmt.Println(line)
}
