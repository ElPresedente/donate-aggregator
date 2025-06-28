package logic

import (
	"go-back/sources"
)

type LogicEvent struct {
	event_type string
	message    string
}

type Logic struct {
	// internal data
	//Тут короче нужен такой же канал для ивентов от логики - они также будут оформленны в одной структуре для передачи на фронт
	roulette *Roulette
}

func (a *Logic) NewLogicProcessor() {
	//да, я всё делаю на объектах, отъебитесь
	a.roulette = &Roulette{}
}

func (a *Logic) Process(donate sources.DonationEvent) {

}

func (a *Logic) EraseRouletteQueue() {
	a.roulette.queue = make([]sources.DonationEvent, 0)
}
