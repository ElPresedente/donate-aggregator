package logic

import (
	"go-back/sources"
	"log"
	"time"
)

type DonateEvent = sources.DonationEvent

type RouletteSectorItem struct { //может потом еще что то понадобится (отдельные шансы для элемента внутри????)
	name string
}

type RouletteSector struct {
	name        string
	probability int                  //оно в интерфейсе инт, пусть и в базе и тут будет инт, на месте поделим
	items       []RouletteSectorItem //надеюсь массив так пишется)
}

type RouletteSettings struct {
	sectors []RouletteSector
}

type Roulette struct {
	sum        int
	roll_price int
	timeout    time.Duration

	queue []DonateEvent
	stop  chan struct{}
}

func (a *Roulette) rouletteLoop() {
	lastRoll := time.Now()
	for {
		if len(a.queue) == 0 {
			return
		}
		select {
		case <-a.stop:
			return
		default:
			if time.Since(lastRoll) >= a.timeout {
				//прокрутить рандомно число от 1 до 100
				//динамически присвоить каждому сектору промежуток чисел
				//получить сектор, получить варианты сектора
				//выбрать случайно один из вариантов
				//выдать событие рулетки

				//и всё это я не могу сделать без инета потому что не знаю как пользоваться этим языком

			} else {
				time.Sleep(a.timeout /* - time.Since( lastRoll )*/)
			}
		}
	}
}

// ТУТ ТЫКАЕТСЯ ЯРЧЕ ДЛЯ ПРОВЕРКИ ДОСТУПА ИЗ JS-------------------------------------------------------------------------
func (a *Roulette) roll() {
	log.Printf("Нажатие кнопки Крутить")
}

//ТУТ ТЫКАЕТСЯ ЯРЧЕ ДЛЯ ПРОВЕРКИ ДОСТУПА ИЗ JS-------------------------------------------------------------------------

func (a *Roulette) process() {
	if !(len(a.queue) == 0) {
		return
	}
	go RouletteLoop()
}

func (a *Roulette) EnqueueDonate(event *DonateEvent) {
	//queue.emplace_back( event ) не знаю как правильно)
	a.process()
}

func RouletteLoop() {

}
