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
	timeout    time.Time

	queue      []DonateEvent
	stop       chan struct{}
}

func (a *Roulette) rouletteLoop(){
	lastRoll := time.Now()
	for{
		if a.queue.empty(){
			return
		}
		select{
			case <- a.stop:
				return
			default:
				if time.Since( lastRoll ) >= timeout{
					//прокрутить рандомно число от 1 до 100 
					//динамически присвоить каждому сектору промежуток чисел
					//получить сектор, получить варианты сектора
					//выбрать случайно один из вариантов
					//выдать событие рулетки

					//и всё это я не могу сделать без инета потому что не знаю как пользоваться этим языком
					
				} else{
					time.Sleep( timeout /* - time.Since( lastRoll )*/)
				} 
		}
	}
}

func (a *Roulette) process(){
	if !queue.empty(){
		return
	}
	go rouletteLoop();
}

func (a *Roulette) EnqueueDonate( event *DonateEvent ){
	//queue.emplace_back( event ) не знаю как правильно)
	process();	
}
