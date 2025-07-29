package logic

type EventName string

// вместо enum для LogicResponse.name
const (
	RouletteSpin                    EventName = "roulette-spin-result"
	RouletteBalanceUpdate           EventName = "roulette-update-remaining-amount"
	RouletteDonateQueueLengthUpdate EventName = "roulette-update-queue-length"
)

type LogicEvent struct {
	name EventName
	data any
}

// ?? enum для event type (пока только 1 значение - roulette spin)
// Я так понимаю, LogicEvent при добавлении канала больше не нужен будто,
// но пока оставил на подумать

// type LogicEvent struct {
// 	event_type string
// 	message    string //json
// }
