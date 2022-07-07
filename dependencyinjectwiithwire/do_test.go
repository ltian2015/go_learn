package dependencyinjectwiithwire

import (
	"testing"
)

func TestNoDi(t *testing.T) {
	msg := NewMessage()
	greeter := NewGreeter(msg)
	event := NewEvent(greeter)
	event.Start()
}
func TestEventUsingDI(t *testing.T) {
	event := InitializeEvent1()
	event.Start()

}
func TestEventUsingDI2(t *testing.T) {
	event, err := InitializeEvent2(false, true)
	if err == nil {
		event.Start()
	} else {
		println(err.Error())
	}

}
