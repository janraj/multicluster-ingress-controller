package events

import (
	"errors"
)

type Event struct {
	eventName            string
	eventWatchStopSwitch chan struct{}
}

const (
	Endpoints  string = "endpoints"
	Services   string = "services"
	Pods       string = "pods"
	Secrets    string = "secrets"
	Ingresses  string = "ingresses"
	Namespaces string = "namespaces"
)

var (
	allowedEventNames = [...]string{Ingresses, Endpoints, Services, Pods, Secrets, Namespaces}
)

func (ev *Event) GetEventName() string {
	return ev.eventName
}

func (ev *Event) GetEventStopChannel() chan struct{} {
	return ev.eventWatchStopSwitch
}
func New(eventName string) (*Event, error) {
	ev := Event{eventName: eventName, eventWatchStopSwitch: make(chan struct{})}
	err := ev.validate()
	if err != nil {
		return nil, err
	}
	return &ev, nil
}
func (ev *Event) validate() error {
	for _, v := range allowedEventNames {
		if ev.eventName == v {
			return nil
		}
	}
	return errors.New("Invalid EventName: " + ev.eventName)

}

func (ev *Event) stopEventWatch() {
	close(ev.eventWatchStopSwitch)
}

func (ev *Event) DeleteAll() {
	ev.stopEventWatch()
}
