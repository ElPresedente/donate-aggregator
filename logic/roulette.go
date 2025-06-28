package logic

import (
	"go-back/sources"
	"time"
)

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

	queue []sources.DonationEvent
}
