package main

import (
	"context"
	"go-back/logic"
	"go-back/sources"
	"go-back/widget"
)

type App struct {
	ctx         context.Context
	logic       logic.Logic
	widgetHub   widget.WidgetsHub
	collManager *sources.CollectorManager
}

func NewApp() *App {
	a := App{
		logic:       logic.NewLogicProcessor(),
		widgetHub:   widget.NewWidgetsHub(),
		collManager: &sources.CollectorManager{},
	}
	a.widgetHub.LogicEventHandler = &a.logic
	a.logic.WidgetEventHandler = &a.widgetHub
	return &a
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.logic.AppCtx = ctx
	// Создаём канал для событий
	eventCh := make(chan sources.DonationEvent, 100)
	a.collManager = sources.NewCollectorManager(ctx, eventCh)

	go a.widgetHub.Start(":8080")

	// Обрабатываем события из канала
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case donation := <-eventCh:
				a.logic.Process(donation)
			}
		}
	}()
}
