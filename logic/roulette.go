package logic

import (
	"go-back/sources"
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
	lastRoll   time.Time
	timeout    time.Duration

	queue []DonateEvent
	stop  chan struct{}
}

// func (r *Roulette) rouletteLoop() {
// 	for {
// 		if len(r.queue) == 0 {
// 			return
// 		}
// 		select {
// 		case <-r.stop:
// 			return
// 		default:
// 			if time.Since(r.lastRoll) >= r.timeout {
// 				//прокрутить рандомно число от 1 до 100
// 				//динамически присвоить каждому сектору промежуток чисел
// 				//получить сектор, получить варианты сектора
// 				//выбрать случайно один из вариантов
// 				//выдать событие рулетки

// 				//и всё это я не могу сделать без инета потому что не знаю как пользоваться этим языком

// 			} else {
// 				time.Sleep(r.timeout /* - time.Since( lastRoll )*/)
// 			}
// 		}
// 	}
// }

// короче потом буду дальше думать, пока не работает вообще
func (r *Roulette) process() {
	if !(len(r.queue) == 0) {
		return
	}

	var currentSum float64

	for _, val := range r.queue {
		currentSum += val.Amount
	}

	if currentSum >= float64(r.roll_price) {
		// жёсткая прокрутка
		// go r.roll()

	}

	// go r.rouletteLoop()
}

func (r *Roulette) EnqueueDonate(event *DonateEvent) {
	//queue.emplace_back( event ) не знаю как правильно)
	r.queue = append(r.queue, *event)

	r.process()
}
