package logic

import (
	"go-back/sources"
	"math/rand"
	"sync"
	"time"
)

type DonateEvent = sources.DonationEvent

type RouletteSectorItem struct {
	// probability int // (отдельные шансы для элемента внутри????)
	name string
}

type RouletteSector struct {
	name        string
	probability int //оно в интерфейсе инт, пусть и в базе и тут будет инт, на месте поделим
	items       []RouletteSectorItem
}

type RouletteSettings struct {
	sectors []RouletteSector
}

type Roulette struct {
	actualAmount float64
	rollPrice    int
	lastRoll     time.Time
	timeout      time.Duration

	queue    []DonateEvent
	stop     chan struct{}
	mu       sync.Mutex
	settings RouletteSettings
}

func NewRouletteProcessor() Roulette {
	return Roulette{
		rollPrice: 100,
		timeout:   1 * time.Second,
		stop:      make(chan struct{}),
		settings: RouletteSettings{
			sectors: []RouletteSector{
				{
					name:        "Виски",
					probability: 20,
					items: []RouletteSectorItem{
						{name: "1"},
						{name: "2"},
						{name: "3"},
					},
				},
				{
					name:        "Виски",
					probability: 20,
					items: []RouletteSectorItem{
						{name: "4"},
						{name: "5"},
						{name: "6"},
					},
				},
				{
					name:        "Виски",
					probability: 20,
					items: []RouletteSectorItem{
						{name: "7"},
						{name: "8"},
						{name: "9"},
					},
				},
			},
		},
	}
}

func (r *Roulette) UpdateDataFromDB() {
	// подсасываем настройки из бд
}

func (r *Roulette) rouletteLoop(logic *Logic) {
	ticker := time.NewTicker(r.timeout)
	defer ticker.Stop()

	for {
		if len(r.queue) == 0 {
			return
		}
		select {
		case <-r.stop:
			return
		case <-ticker.C:
			if time.Since(r.lastRoll) >= r.timeout {
				r.lastRoll = time.Now()

				r.actualAmount += r.queue[0].Amount
				r.DequeueDonate()

				if r.actualAmount >= float64(r.rollPrice) {
					winnerSector := chooseSector(r.settings.sectors)
					winnerItem := chooseSectorItem(winnerSector.items)

					logic.DispatchLogicEvent(LogicEvent{
						name: RouletteSpin,
						data: []any{
							winnerSector,
							winnerItem,
						},
					})
					r.actualAmount -= float64(r.rollPrice)
				}
			}
		}
	}
}

func (r *Roulette) Process(event *DonateEvent, logic *Logic) {
	r.EnqueueDonate(event)
	r.UpdateDataFromDB()

	if len(r.queue) > 1 {
		return
	}

	go r.rouletteLoop(logic)
}

func (r *Roulette) EnqueueDonate(event *DonateEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.queue = append(r.queue, *event)
}

func (r *Roulette) DequeueDonate() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.queue = r.queue[1:]
}

func chooseSector(sectors []RouletteSector) RouletteSector {
	total := 0
	for _, s := range sectors {
		total += s.probability
	}

	r := rand.Intn(total)
	sum := 0
	for _, s := range sectors {
		sum += s.probability
		if r < sum {
			return s
		}
	}

	// Это на всякий случай — никогда не должно сработать
	return sectors[len(sectors)-1]
}

func chooseSectorItem(items []RouletteSectorItem) RouletteSectorItem {
	// Если нужны будут разные шансы на item
	// total := 0
	// for _, i := range items {
	// 	total += i.probability
	// }

	// r := rand.Intn(total)
	// sum := 0
	// for _, i := range items {
	// 	sum += i.probability
	// 	if r < sum {
	// 		return i
	// 	}
	// }

	return items[rand.Intn(len(items))]
}
