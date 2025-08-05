package sources

import "fmt"

type EventType struct {
	name string
}

var baseEventTypeNames = map[string]struct{}{
	"DonationEvent": {},
}

func NewEventType(typeName string) (EventType, error) {
	if _, ok := baseEventTypeNames[typeName]; !ok {
		return EventType{}, fmt.Errorf("Ошибка создания типа ивента. Не существует типа %s", typeName)
	}

	return EventType{name: typeName}, nil
}

func (et *EventType) GetTypeName() string {
	return et.name
}
