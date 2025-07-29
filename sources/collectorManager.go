package sources

import (
	"context"
	"fmt"
	"go-back/database"
	"log"
	"strings"
	"sync"
)

var CollectorType []string = []string{"Donatty", "DonatePay"}

type CollectorManager struct {
	mu         sync.Mutex
	ctx        context.Context
	collectors map[string]*managedCollector
	eventCh    chan DonationEvent
}

type managedCollector struct {
	collector EventCollector
	cancel    context.CancelFunc
}

func NewCollectorManager(ctx context.Context, eventCh chan DonationEvent) *CollectorManager {
	return &CollectorManager{
		ctx:        ctx,
		collectors: make(map[string]*managedCollector),
		eventCh:    eventCh,
	}
}

func (m *CollectorManager) NewManagedCollector(ctx context.Context, cancel context.CancelFunc, name string) error {

	var collector EventCollector

	switch name {
	case "Donatty":
		//GetENVValue
		token, _ := database.CredentialsDB.GetENVValue("donattyToken")
		url, _ := database.CredentialsDB.GetENVValue("donattyUrl")
		collector = NewDonattyCollector(m.ctx, token, url, m.eventCh)
	case "DonatePay":
		token, _ := database.CredentialsDB.GetENVValue("donatpayToken")
		userId, _ := database.CredentialsDB.GetENVValue("donatpayUserId")
		collector = NewDonatePayCollector(m.ctx, token, userId, m.eventCh)
	default:
		return fmt.Errorf("❌ Ошибка создания коллектора. Коллектор с именем %s не найден", name)
	}

	m.collectors[name] = &managedCollector{
		collector: collector,
		cancel:    cancel,
	}
	return nil
}

func (m *CollectorManager) StartCollector(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx, cancel := context.WithCancel(m.ctx)

	if _, exists := m.collectors[name]; exists {
		cancel()
		return fmt.Errorf("❌ Ошибка запуска коллектора. Коллектор с именем %s уже работает", name)
	}

	err := m.NewManagedCollector(ctx, cancel, name)

	if err != nil {
		return err
	}

	// Запускаем в отдельной горутине
	go func() {
		if err := m.collectors[name].collector.Start(ctx); err != nil && ctx.Err() == nil {
			log.Printf("❌ Ошибка в коллекторе %s: %v", name, err)
		}
	}()

	log.Printf("✅ Коллектор %s запущен", name)

	return nil
}

func (m *CollectorManager) StartAllCollector() error {
	failedCollectors := []string{}
	for _, name := range CollectorType {
		err := m.StartCollector(name)
		if err != nil {
			failedCollectors = append(failedCollectors, name)
		}
	}

	if len(failedCollectors) > 0 {
		return fmt.Errorf("❌ Ошибка запуска коллекторов: %s", strings.Join(failedCollectors, ", "))
	}
	return nil
}

func (m *CollectorManager) StopCollector(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	mc, exists := m.collectors[name]
	if !exists {
		return fmt.Errorf("❌ Ошибка остановки. Коллектор %s не запущен", name)
	}

	// Сигнал отмены
	mc.cancel()

	// остановка логики коллектора
	if err := mc.collector.Stop(); err != nil {
		log.Printf("❌ Ошибка остановки %s: %v", name, err)
	}

	delete(m.collectors, name)
	return nil
}

func (m *CollectorManager) StopAllCollector() error {
	failedCollectors := []string{}
	for _, name := range CollectorType {
		err := m.StopCollector(name)
		if err != nil {
			failedCollectors = append(failedCollectors, name)
		}
	}

	if len(failedCollectors) > 0 {
		return fmt.Errorf("❌ Ошибка остановки коллекторов: %s", strings.Join(failedCollectors, ", "))
	}
	return nil
}

func (m *CollectorManager) IsCollectorActive(name string) bool {
	_, exists := m.collectors[name]
	return exists
}
