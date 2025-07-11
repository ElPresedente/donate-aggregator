package logic

import (
	"go-back/sources"
)

// ?? enum для event type (пока только 1 значение - roulette spin)
// Я так понимаю, LogicEvent при добавлении канала больше не нужен будто,
// но пока оставил на подумать
// type LogicEvent struct {
// 	event_type string
// 	message    string //json
// }

type Logic struct {
	// internal data
	//Тут короче нужен такой же канал для ивентов от логики - они также будут оформленны в одной структуре для передачи на фронт
	roulette           Roulette
	CommunicateChannel chan LogicResponse
	stop               chan struct{}
}

func NewLogicProcessor() Logic {
	//да, я всё делаю на объектах, отъебитесь
	return Logic{
		roulette:           NewRouletteProcessor(),
		CommunicateChannel: make(chan LogicResponse),
	}
}

func (l *Logic) Process(donate sources.DonationEvent) {

	l.roulette.Process(&donate, l.CommunicateChannel)

	//db.SaveLog( donate )

	//проверить на сообщения от процессоров
	//dispatch инветов
}

func (l *Logic) DispatchLogicEvent( /*le LogicEvent*/ ) {
	for {
		select {
		case response := <-l.CommunicateChannel:
			switch response.name {
			case RouletteSpin:
				//front.emitEvent(...)
				//передать виджету по websocket результат прокрутки

			}
			// case
		case <-l.stop:
			return
		}
	}
}

func (l *Logic) EraseRouletteQueue() {
	l.roulette.queue = make([]sources.DonationEvent, 0)
}
