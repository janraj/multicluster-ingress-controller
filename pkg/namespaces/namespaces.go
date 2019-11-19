package namespaces

import (
	"errors"
	"k8s.io/api/core/v1"
	"multicluster-ingress-controller/pkg/events"
	"multicluster-ingress-controller/pkg/swaggerInterface"
	"strconv"
)

type Namespace struct {
	namespaceName string
	eventList     []events.Event
}

func New(namespaceName string) *Namespace {
	return &Namespace{namespaceName: namespaceName}
}

func (nS *Namespace) AddEvents(cS swaggerInterface.CreateClientServer) error {
	var evList []events.Event
	for _, v := range cS.WatchEvents {
		ev, err := events.New(v)
		if err != nil {
			return err
		}
		evList = append(evList, *ev)
	}

	nS.eventList = evList
	return nil
}

func (nS *Namespace) DeleteAll() {
	for _, e := range nS.eventList {
		e.DeleteAll()
	}
	nS.eventList = nil
}

func (nS *Namespace) GetNamespaceName() string {
	return nS.namespaceName
}

func (nS *Namespace) EventsLength() int {
	return len(nS.eventList)
}
func (nS *Namespace) GetNamespaceString() string {
	return GetNamespaceString(nS.namespaceName)
}

func GetNamespaceString(namespace string) string {
	if namespace == v1.NamespaceAll {
		return "allNamespace"
	}
	return namespace
}

func (nS *Namespace) GetEvent(index int) (*events.Event, error) {
	length := len(nS.eventList)
	if index > length {
		return nil, errors.New("Namespaces: Out of bounds Index: " + strconv.Itoa(index) + "total length: " + strconv.Itoa(length))
	}
	return &nS.eventList[index], nil
}
