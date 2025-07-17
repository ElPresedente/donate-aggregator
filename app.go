package main

import (
	"context"
	"go-back/logic"
	"go-back/sources"
	"log"
	"os"
	"github.com/joho/godotenv"
)

type App struct {
  ctx      context.Context
  logic    logic.Logic
  ws       *sources.WebSocketHub
}

func NewApp() *App {
  return &App{
    logic:	logic.NewLogicProcessor(),
	ws:		sources.NewWebSocketHub(),
  }
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.ws.Start()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %s", err)
	}

	// Создаём канал для событий
	eventCh := make(chan sources.DonationEvent, 100)

	//включение-выключение коллекторов по сигналу от фронта

	// Список коллекторов
	collectors := []sources.EventCollector{
		sources.NewDonattyCollector(os.Getenv("DONATTY_TOKEN"), os.Getenv("DONATTY_REF"), eventCh),
		sources.NewDonatePayCollector(os.Getenv("DONATPAY_TOKEN"), os.Getenv("DONATPAY_USERID"), eventCh),
	}

	// Запускаем все коллекторы
	for _, collector := range collectors {
		go func(c sources.EventCollector) {
			if err := c.Start(ctx); err != nil {
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
				a.logic.Process(donation)
				//runtime.EventsEmit(a.ctx, "donation", donation) -> в logic.Process
			}
		}
	}()
}