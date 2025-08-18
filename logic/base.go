package logic

import (
	"context"
	"encoding/json"
	"go-back/database"
	"go-back/l2db"
	"go-back/l2wbridge"
	"go-back/sources"
	"log"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Logic struct {
	roulette           Roulette
	WidgetEventHandler l2wbridge.L2WHandler
	AppCtx             context.Context
}

func NewLogicProcessor() Logic {
	return Logic{
		roulette:           NewRouletteProcessor(),
		WidgetEventHandler: nil,
		AppCtx:             nil,
	}
}

func (l *Logic) LogicEventHandler(request string, data string) {
	switch request {
	case "spins-done":
		l.roulette.isWorking = false
		l.roulette.rouletteLoop(l)
	case "rouletteConnected":
		log.Println("Roulette widget connected")
		l.roulette.isWorking = false
		l.roulette.rouletteLoop(l)
		runtime.EventsEmit(l.AppCtx, "rouletteConnectionUpdated", "connected")
	case "rouletteDisconnected":
		log.Println("Roulette widget disconnected")
		runtime.EventsEmit(l.AppCtx, "rouletteConnectionUpdated", "disconnected")
	}
}

func (l *Logic) ManualRouletteSpin() {
	var spin sources.RouletteEvent
	spin.Name = "Пользователь"
	spin.SpinsAmount = 1
	l.roulette.ProcessSpin(&spin, l)
}

func (l *Logic) Process(event sources.CollectorEvent) {
	switch event.EventType.GetTypeName() {
	case "DonationEvent":
		l.roulette.ProcessDonate(event.Event.(*DonateEvent), l)
	case "RouletteSpinEvent":
		l.roulette.ProcessSpin(event.Event.(*RouletteEvent), l)
	}

}

func (l *Logic) DispatchLogicEvent(le LogicEvent) {
	switch le.name {
	case RouletteSpin:
		for key, val := range le.data.(l2db.ResponseData).Spins {
			log.Printf("Прокрут рулетки №%d. Результат: категория:%s сектор:%s", key+1, val.WinnerCategory, val.WinnerSector)
		}

		database.LogDB.InsertSpins(le.data.(l2db.ResponseData))

		l.sendToWidgetSpinData(le.data.(l2db.ResponseData))
		l.sendToFrontSpinData(le.data.(l2db.ResponseData))
	case RouletteBalanceUpdate:
		runtime.EventsEmit(l.AppCtx, "currentAmountUpdate", le.data.(float64))
		//	l.WidgetEventHandler.WidgetEventHandler("set-roulette-balance", le.data.(float64)) //ченить такое будет
	case RouletteDonateQueueLengthUpdate:
		runtime.EventsEmit(l.AppCtx, "donateQueueLengthUpdate", le.data.(int))

	}
}

func (l *Logic) sendToWidgetSpinData(spinData l2db.ResponseData) {
	var data struct {
		Request string          `json:"request"`
		Spins   []l2db.SpinData `json:"spins"`
	}
	data.Request = "enqueue-spins"
	data.Spins = spinData.Spins
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Json encoding failed %v", err)
	}

	l.WidgetEventHandler.WidgetEventHandler("enqueue-spins", string(jsonData))
}

func (l *Logic) sendToFrontSpinData(spinData l2db.ResponseData) {
	var data struct {
		User  string          `json:"user"`
		Time  string          `json:"time"`
		Spins []l2db.SpinData `json:"spins"`
	}
	data.User = spinData.User
	data.Time = time.Now().Format("02.01 15:04")
	data.Spins = spinData.Spins
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Json encoding failed %v", err)
	}

	runtime.EventsEmit(l.AppCtx, "logUpdated", string(jsonData))
}

func (l *Logic) ReloadRoulette() {
	l.roulette.Reload(l)
	l.roulette.rouletteLoop(l)
	l.WidgetEventHandler.WidgetEventHandler("reloadRoulette", "")

}

func (l *Logic) EraseRouletteQueue() {
	l.roulette.queue = make([]sources.DonationEvent, 0)
}
