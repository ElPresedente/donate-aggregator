package logic

import (
	"go-back/sources"
)

//?? enum для event type (пока только 1 значение - roulette spin)

type LogicEvent struct {
	event_type string
	message    string //json
}

type Logic struct {
	// internal data
	//Тут короче нужен такой же канал для ивентов от логики - они также будут оформленны в одной структуре для передачи на фронт
	roulette Roulette
	//channel для приёма ивентов от процессоров (пока только рулетка)
}

func NewLogicProcessor() Logic {
	//да, я всё делаю на объектах, отъебитесь
	return Logic{
		roulette: NewRouletteProcessor(),
	}
}

func (l *Logic) Process(donate sources.DonationEvent) {

	l.roulette.Process(&donate)

	//db.SaveLog( donate )

	//проверить на сообщения от процессоров
	//dispatch инветов
}

func (l *Logic) DispatchLogicEvent(le LogicEvent) {
	//do smthn
	switch le.event_type {
	case "roulette-spin":
		//front.emitEvent(...)
		//передать виджету по websocket результат прокрутки
		//case ...
	}
}

func (l *Logic) EraseRouletteQueue() {
	l.roulette.queue = make([]sources.DonationEvent, 0)
}
