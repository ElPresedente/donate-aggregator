package logic

import (
	"encoding/json"
	"go-back/l2wbridge"
	"go-back/sources"
	"log"
)

type Logic struct {
	roulette           Roulette
	WidgetEventHandler l2wbridge.L2WHandler
}

func NewLogicProcessor() Logic {
	return Logic{
		roulette:           NewRouletteProcessor(),
		WidgetEventHandler: nil,
	}
}

func (l *Logic) LogicEventHandler(request string, data string) {
	switch request {
	case "spins-done":
		l.roulette.isWorking = false
		l.roulette.rouletteLoop(l)
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
			log.Printf("Прокрут рулетки №%d. Результат: категория:%s сектор:%s", key+1, val.WinnerCategory, val.WinnerSector)
		}

		var data struct {
			Request string     `json:"request"`
			Spins   []SpinData `json:"spins"`
		}
		data.Request = "enqueue-spins"
		data.Spins = le.data.(ResponseData).Spins
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Json encoding failed %v", err)
		}
		l.WidgetEventHandler.WidgetEventHandler("enqueue-spins", string(jsonData))
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
