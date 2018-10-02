package lexer

import (
	"io/ioutil"
	"log"
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

func Start(filename string) {
	l := Lexer{}
	l.eventsQueue = EventDrivenModule{Events: []Event{Event{l.readFile, filename}}}
	l.eventsQueue.AddEvent(Event{l.firstLine, nil})
	for !l.eventsQueue.IsEmpty() {
		l.executeHandler()
	}

	// for !eventsQueue.IsEmpty() {
	// 	firstEvent := eventsQueue.PopEvent()
	// 	firstEvent.Handler()
	// }
}

func (l *Lexer) executeHandler() {
	firstEvent := l.eventsQueue.PopEvent()
	firstEvent.Handler(firstEvent.Args)
}

func (l *Lexer) readFile(arg interface{}) {
	filename := arg.(string)
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	l.Source = string(dat)
}

func (l *Lexer) firstLine(arg interface{}) {
	l.line = 0
	l.eventsQueue.AddEvent(Event{l.readLine, nil})
}

func (l *Lexer) readLine(arg interface{}) {
	line := "TODO"
	for _, c := range line {
		print(c)
	}
	l.line++
	l.eventsQueue.AddEvent(Event{l.readLine, nil})
}

func (l Lexer) filtroAsc(arg interface{}) {

}

// func

// func (l Lexer) ScanTokens() []Token {
// 	e := EventDrivenModule{Events: []Event{Event{l.scanToken}}}
// 	f := e.PopEvent()
// 	fmt.Printf("%v", f)

// 	for !l.isAtEnd() {
// 		l.Start = l.Current
// 		l.scanToken()
// 	}

// 	return nil
// }

// func readFile(fileName string) interface{} {
// 	dat, err := ioutil.ReadFile(fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return string(dat)
// }

// func (l Lexer) readLine() interface{} {
// 	return nil
// }

// func nextLine() {

// }

// func (l Lexer) isAtEnd() bool {
// 	return l.Current >= len(l.Source)
// }

// func (l *Lexer) scanToken() {
// 	fmt.Print("oi")
// 	c := l.advance()
// 	switch c {
// 	case "(":
// 		fmt.Print("(\n")
// 	}
// }

// func (l *Lexer) advance() string {
// 	l.Current++
// 	return string(l.Source[l.Current-1])
// }
