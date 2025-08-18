package logic

import (
	"go-back/database"
	"go-back/l2db"
	"go-back/sources"
	"log"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type DonateEvent = sources.DonationEvent
type RouletteEvent = sources.RouletteEvent

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

type Roulette struct {
	actualAmount      float64
	rollPrice         int
	rollPriceIncrease int
	lastRoll          time.Time
	timeout           time.Duration
	lastDonate        DonateEvent

	queue       []DonateEvent
	eventsQueue []RouletteEvent
	stop        chan struct{}
	mu          sync.Mutex
	settings    RouletteSettings
	isWorking   bool
}

func NewRouletteProcessor() Roulette {
	return Roulette{
		rollPrice:         100,
		rollPriceIncrease: 0,
		timeout:           1 * time.Second,
		stop:              make(chan struct{}),
		settings:          RouletteSettings{},
		isWorking:         true,
	}
}

func (r *Roulette) UpdateDataFromDB() {
	dbRollPrice, err := database.WidgetDB.GetRouletteSettingValue("rollPrice")
	if err != nil {
		log.Printf("❌ Ошибка получения стоимости прокрута: %s", err)
	}
	dbRollIncreasePrice, err := database.WidgetDB.GetRouletteSettingValue("rollPriceIncrease")

	if err != nil {
		log.Printf("❌ Ошибка получения размера увеличения стоимости прокрута: %s", err)
	}

	dbRollPriceInt, err := strconv.Atoi(dbRollPrice)
	if err != nil {
		log.Printf("❌ Ошибка приведения типа стоимости прокрута: %s", err)
	}
	dbRollIncreasePriceInt, err := strconv.Atoi(dbRollIncreasePrice)

	if err != nil {
		log.Printf("❌ Ошибка приведения типа размера увеличения стоимости прокрута: %s", err)
	}

	if dbRollPriceInt > 0 {
		r.rollPrice = dbRollPriceInt
	}

	r.rollPriceIncrease = dbRollIncreasePriceInt

	categories, err := database.WidgetDB.GetRouletteCategorys()

	if err != nil {
		log.Printf("❌ Ошибка получения категорий: %s", err)
	}

	for _, category := range categories {
		newCategory := RouletteCategory{}
		newCategory.name = category.Name
		newCategory.probability = int(category.Percentage * 100)

		sectors, err := database.WidgetDB.GetSectorsByCategoryID(category.ID)
		if err != nil {
			log.Printf("❌ Ошибка получения секторов: %s", err)
		}

		for _, sector := range sectors {
			newSector := RouletteCategorySector{}
			newSector.name = sector.Name

			newCategory.sectors = append(newCategory.sectors, newSector)

		}

		r.settings.categories = append(r.settings.categories, newCategory)
	}
}

func (r *Roulette) rouletteLoop(logic *Logic) {
	r.UpdateDataFromDB()

	if len(r.eventsQueue) > 0 {
		lastEvent := r.eventsQueue[0]
		r.DequeueEvent(logic)
		responses := l2db.ResponseData{
			User: lastEvent.Name,
		}
		for range lastEvent.SpinsAmount {
			spinResult := r.generateSpin()
			responses.Spins = append(responses.Spins, spinResult)
		}
		logic.DispatchLogicEvent(LogicEvent{
			name: RouletteSpin,
			data: responses,
		})
		r.isWorking = true
		return
	}

	for len(r.queue) > 0 {
		r.lastDonate = r.queue[0]
		r.actualAmount += r.lastDonate.Amount
		logic.DispatchLogicEvent(LogicEvent{
			name: RouletteBalanceUpdate,
			data: r.actualAmount,
		})
		r.DequeueDonate(logic)
		if r.actualAmount >= float64(r.rollPrice) {
			responses := l2db.ResponseData{
				User: r.lastDonate.User,
			}

			baseRollPrice := float64(r.rollPrice)

			for r.actualAmount >= float64(r.rollPrice) {
				spinResult := r.generateSpin()
				responses.Spins = append(responses.Spins, spinResult)
				r.actualAmount -= float64(r.rollPrice)
				//r.rollPrice каждый раз берется заного из БД, поэтому её можно "портить" (если я всё правильно понял)
				r.rollPrice += r.rollPriceIncrease
			}

			if r.actualAmount >= baseRollPrice {
				spinResult := r.generateSpin()
				responses.Spins = append(responses.Spins, spinResult)
				r.actualAmount = math.Mod(r.actualAmount, baseRollPrice)
			}

			logic.DispatchLogicEvent(LogicEvent{
				name: RouletteSpin,
				data: responses,
			})
			logic.DispatchLogicEvent(LogicEvent{
				name: RouletteBalanceUpdate,
				data: r.actualAmount,
			})
			r.isWorking = true
			return
		}
	}
}

func (r *Roulette) ProcessSpin(spins *RouletteEvent, logic *Logic) {
	r.UpdateDataFromDB()

	r.EnqueueEvent(spins, logic)

	if r.isWorking {
		return
	}

	r.rouletteLoop(logic)
}

func (r *Roulette) generateSpin() l2db.SpinData {
	winnerCategory := chooseCategory(r.settings.categories)
	winnerSector := chooseCategorySector(winnerCategory.sectors)

	spinResult := l2db.SpinData{
		WinnerCategory: winnerCategory.name,
		WinnerSector:   winnerSector.name,
	}

	return spinResult
}

func (r *Roulette) Reload(logic *Logic) {
	r.actualAmount = 0
	r.isWorking = false
	r.queue = []DonateEvent{}
	r.eventsQueue = []RouletteEvent{}
	logic.DispatchLogicEvent(LogicEvent{
		name: RouletteBalanceUpdate,
		data: r.actualAmount,
	})
	logic.DispatchLogicEvent(LogicEvent{
		name: RouletteDonateQueueLengthUpdate,
		data: len(r.queue),
	})
}

func (r *Roulette) ProcessDonate(event *DonateEvent, logic *Logic) {
	r.EnqueueDonate(event, logic)

	if r.isWorking {
		return
	}

	r.rouletteLoop(logic)
}

func (r *Roulette) EnqueueDonate(event *DonateEvent, logic *Logic) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.queue = append(r.queue, *event)
	logic.DispatchLogicEvent(LogicEvent{
		name: RouletteDonateQueueLengthUpdate,
		data: len(r.queue),
	})
}

func (r *Roulette) DequeueDonate(logic *Logic) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.queue = r.queue[1:]
	logic.DispatchLogicEvent(LogicEvent{
		name: RouletteDonateQueueLengthUpdate,
		data: len(r.queue),
	})
}

func (r *Roulette) EnqueueEvent(event *RouletteEvent, logic *Logic) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.eventsQueue = append(r.eventsQueue, *event)
}

func (r *Roulette) DequeueEvent(logic *Logic) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.eventsQueue = r.eventsQueue[1:]
}

func chooseCategory(categories []RouletteCategory) RouletteCategory {
	total := 0
	for _, s := range categories {
		total += s.probability
	}

	// Добавить проверку на total > 0 и придумать, что в этом случае делать
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
