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

type AsciiCategorizer struct {
	EventDrivenModule
	filename string
	file     os.File
	scanner  *bufio.Scanner
}

func (a *AsciiCategorizer) HandleEvent(event Event) {
	handlers := map[string]func(string){
		"open":  a.OpenHandler,
		"read":  a.ReadHandler,
		"close": a.CloseHandler}
	handler := handlers[event.Name]
	handler(event.Arg)
}

func (a *AsciiCategorizer) OpenHandler(filename string) {
	f, err := os.Open(filename)
	check(err)

	a.scanner = bufio.NewScanner(f)
	a.AddEvent(Event{"read", ""})
}

func (a *AsciiCategorizer) ReadHandler(arg string) {
	a.scanner.Scan()
	fmt.Println(a.scanner.Text())
	if text := a.scanner.Scan(); text {
		fmt.Println(a.scanner.Text())
		a.AddEvent(Event{"read", ""})
	} else {
		a.AddEvent(Event{"close", ""})
	}
}

func (a *AsciiCategorizer) CloseHandler(arg string) {
	fmt.Println("closed")
	a.file.Close()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
