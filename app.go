package main

import (
	"context"
	"fmt"
	"go-back/sources"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Создаём канал для событий
	eventCh := make(chan sources.DonationEvent, 100)

	// Список коллекторов
	collectors := []sources.EventCollector{
		sources.NewDonattyCollector(os.Getenv("DONATTY_TOKEN"), os.Getenv("DONATTY_REF")),
		sources.NewDonatePayCollector(os.Getenv("DONATPAY_TOKEN"), os.Getenv("DONATPAY_USERID")),
	}

	// Запускаем все коллекторы
	for _, collector := range collectors {
		go func(c sources.EventCollector) {
			if err := c.Start(ctx, eventCh); err != nil {
				log.Printf("❌ Ошибка коллектора: %v", err)
			}
		}(collector)
	}

	// Обрабатываем события из канала
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Останавливаем все коллекторы при завершении
				for _, collector := range collectors {
					if err := collector.Stop(); err != nil {
						log.Printf("❌ Ошибка остановки коллектора: %v", err)
					}
				}
				return
			case donation := <-eventCh:
				// Отправка события в фронтенд (для будущего GUI)
				runtime.EventsEmit(a.ctx, "donation", donation)
			}
		}
	}()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
