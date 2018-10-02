package lexer

type EventDrivenModule struct {
	Events []Event
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
	Handler func(interface{})
	Args    interface{}
}
