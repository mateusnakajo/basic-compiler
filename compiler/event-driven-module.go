package compiler

type EventDrivenModule struct {
	Events      []Event
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

func (e *EventDrivenModule) LookAhead() Event {
	return e.Events[0]
}

func (e *EventDrivenModule) IsEmpty() bool {
	return len(e.Events) == 0
}

type Event struct {
	Name string
	Arg  interface{}
}
