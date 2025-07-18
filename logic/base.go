package logic

import (
	"go-back/sources"
	"log"
)

type Logic struct {
	// internal data
	//Тут короче нужен такой же канал для ивентов от логики - они также будут оформленны в одной структуре для передачи на фронт
	roulette Roulette
	stop     chan struct{}
}

func NewLogicProcessor() Logic {
	//да, я всё делаю на объектах, отъебитесь
	return Logic{
		roulette: NewRouletteProcessor(),
	}
}

func (l *Logic) Process(donate sources.DonationEvent) {

	l.roulette.Process(&donate, l)

	//db.SaveLog( donate )

	//проверить на сообщения от процессоров
	//dispatch инветов
}

func (l *Logic) DispatchLogicEvent(le LogicEvent) {
	switch le.name {
	case RouletteSpin:
		for key, val := range le.data.(ResponseData).Spins {
			log.Printf("Прокрут рулетки №%d. Результат: категория:%s сектор:%s", key+1, val.winnerCategory, val.winnerSector)
		}

		//front.emitEvent(...)
		//передать виджету по websocket результат прокрутки
		// объёкт поля для отправки в json
		// 1. (сектор)
		// 2. item

		// case
	}

}

func (l *Logic) EraseRouletteQueue() {
	l.roulette.queue = make([]sources.DonationEvent, 0)
}
