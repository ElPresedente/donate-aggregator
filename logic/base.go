package logic

import(
	"go-back/sources"
)

// ВЫНЕСТИ В logic/roulette.go (мне тут неудобно прост)

type RouletteSectorItem struct{ //может потом еще что то понадобится (отдельные шансы для элемента внутри????)
	name string 
}

type RouletteSector struct{
	name        string
	probability int //оно в интерфейсе инт, пусть и в базе и тут будет инт, на месте поделим
	items       []RouletteSectorItem //надеюсь массив так пишется)
}

type RouletteSettings struct{
	sectors []RouletteSector
	
}

type Roulette struct{
	sum int
	roll_price int
	timeout time.Time

	queue []sources.DonationEvent
}

type LogicEvent struct{
	event_type string
	message string	
}

type Logic struct{
	// internal data
	//Тут короче нужен такой же канал для ивентов от логики - они также будут оформленны в одной структуре для передачи на фронт
	roulette *Roulette
}

func (a *Logic) NewLogicProcessor(){
	//да, я всё делаю на объектах, отъебитесь
	a.roulette = Roulette{}
}

func (a *Logic) process (donate sources.DonationEvent){
	
}

func (a *Logic) erace_roulette_queue(){
	// roulette.queue.erase() хз как)))))))))
}
