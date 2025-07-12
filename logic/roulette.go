package logic

import (
	"go-back/sources"
	"math/rand"
	"sync"
	"time"
)

type DonateEvent = sources.DonationEvent

type RouletteCategorySector struct {
	// probability int // (отдельные шансы для элемента внутри????)
	name string
}

type RouletteCategory struct {
	name        string
	probability int //оно в интерфейсе инт, пусть и в базе и тут будет инт, на месте поделим
	sectors     []RouletteCategorySector
}

type RouletteSettings struct {
	categories []RouletteCategory
}

type ResponseData struct {
	User  string
	Spins []SpinData
}

type SpinData struct {
	winnerCategory string
	winnerSector   string
}

type Roulette struct {
	actualAmount float64
	rollPrice    int
	lastRoll     time.Time
	timeout      time.Duration
	lastDonate   DonateEvent

	queue     []DonateEvent
	stop      chan struct{}
	mu        sync.Mutex
	settings  RouletteSettings
	isWorking bool
}

func NewRouletteProcessor() Roulette {
	return Roulette{
		rollPrice: 100,
		timeout:   1 * time.Second,
		stop:      make(chan struct{}),
		settings: RouletteSettings{
			categories: []RouletteCategory{
				{
					name:        "Виски",
					probability: 20,
					sectors: []RouletteCategorySector{
						{name: "1"},
						{name: "2"},
						{name: "3"},
					},
				},
				{
					name:        "Виски",
					probability: 20,
					sectors: []RouletteCategorySector{
						{name: "4"},
						{name: "5"},
						{name: "6"},
					},
				},
				{
					name:        "Виски",
					probability: 20,
					sectors: []RouletteCategorySector{
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
	for len(r.queue) > 0 {
		r.lastDonate = r.queue[0]
		r.actualAmount += r.lastDonate.Amount
		r.DequeueDonate()

		if r.actualAmount >= float64(r.rollPrice) {
			responses := ResponseData{
				User: r.lastDonate.User,
			}
			for r.actualAmount >= float64(r.rollPrice) {
				winnerCategory := chooseCategory(r.settings.categories)
				winnerSector := chooseCategorySector(winnerCategory.sectors)

				spinResult := SpinData{
					winnerCategory: winnerCategory.name,
					winnerSector:   winnerSector.name,
				}
				responses.Spins = append(responses.Spins, spinResult)
				r.actualAmount -= float64(r.rollPrice)
			}
			logic.DispatchLogicEvent(LogicEvent{
				name: RouletteSpin,
				data: responses,
			})
			r.isWorking = true
			return
		}
	}
}

func (r *Roulette) Process(event *DonateEvent, logic *Logic) {
	r.EnqueueDonate(event)
	r.UpdateDataFromDB()

	if r.isWorking {
		return
	}

	r.rouletteLoop(logic)
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

func chooseCategory(categories []RouletteCategory) RouletteCategory {
	total := 0
	for _, s := range categories {
		total += s.probability
	}

	r := rand.Intn(total)
	sum := 0
	for _, s := range categories {
		sum += s.probability
		if r < sum {
			return s
		}
	}

	// Это на всякий случай — никогда не должно сработать
	return categories[len(categories)-1]
}

func chooseCategorySector(sectors []RouletteCategorySector) RouletteCategorySector {
	// Если нужны будут разные шансы на sector
	// total := 0
	// for _, i := range sectors {
	// 	total += i.probability
	// }

	// r := rand.Intn(total)
	// sum := 0
	// for _, i := range sectors {
	// 	sum += i.probability
	// 	if r < sum {
	// 		return i
	// 	}
	// }

	return sectors[rand.Intn(len(sectors))]
}
