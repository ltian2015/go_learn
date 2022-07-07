package dependencyinjectwiithwire

import (
	"errors"
)

type Message string
type Greeter struct {
	Grumpy bool
	Msg    Message
}
type Event struct {
	Greeter Greeter
}

func (e Event) Start() {
	println(e.Greeter.Msg)
}

func NewMessage() Message {
	return Message("Hello here!")
}

func NewMessage2(grumpy bool) Message {
	if grumpy {
		return Message("go away!")
	}
	return Message("Hello here!")
}

func NewGreeter(msg Message) Greeter {
	return Greeter{Msg: msg}
}
func NewGreeter2(msg Message, grumpy bool) Greeter {

	return Greeter{Msg: msg, Grumpy: grumpy}

}
func NewEvent(greeter Greeter) Event {
	return Event{Greeter: greeter}
}
func NewEvent2(greeter Greeter) (event Event, err error) {
	if greeter.Grumpy {
		err = errors.New("could not create event: event greeter is grumpy")
		return
	}
	event = Event{greeter}
	return
}
