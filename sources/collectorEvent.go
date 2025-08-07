package sources

import "fmt"

type CollectorEvent struct {
	EventType EventType
	Event     any
}

func NewCollectorEvent(eventTypeName string, eventRef any) (CollectorEvent, error) {
	eventType, err := NewEventType(eventTypeName)
	if err != nil {
		return CollectorEvent{}, fmt.Errorf("ошибка создания доната: %v", err)
	}

	event := CollectorEvent{
		EventType: eventType,
		Event:     eventRef,
	}

	return event, nil
}
