package main

import (
	"context"
	"encoding/json"
	"go-back/database"
	"go-back/services"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) FrontendDispatcher(endpoint string, argJSON string) {
	//log.Printf("🛰 Вызов FrontendDispatcher: %s, argJSON: %s", endpoint, argJSON)

	switch endpoint {
	// Получение предметов по ID группы
	case "getSectorsByCategoryId":
		getSectorsByCategoryId(a.ctx, argJSON)

	case "sectorsToSave":
		sectorsToSave(a.ctx, argJSON)

	// Получение всех групп и их итемов
	case "getSectors":
		getSectors(a.ctx)

	case "getSettings":
		getSettings(a.ctx)

	case "updateSettings":
		updateSettings(a, argJSON)

	case "startCollector":
		startCollector(argJSON, a)

	case "startAllCollector":
		startAllCollector(a)

	case "stopAllCollector":
		stopAllCollector(a)

	case "reconnectAllCollector":
		reconnectAllCollector(a)

	case "reconnectDonatty":
		reconnectDonatty(a)
	case "reconnectDonatepay":
		reconnectDonatepay(a)
	case "reloadRoulette":
		reloadRoulette(a)
	case "manualRouletteSpin":
		manualRouletteSpin(a)
	case "getLogs":
		getLogs(a.ctx)
	case "newStream":
		newStream(a)
	case "twitchLoginProcedure":
		twitchLoginProcedure()
	default:
		log.Printf("⚠️ Неизвестный endpoint: %s", endpoint)
	}
}

func twitchLoginProcedure() {
	services.TwitchNewToken()
}

func newStream(a *App) {
	database.LogDB.ClearDatabase()
	a.logic.ReloadRoulette()
	runtime.EventsEmit(a.ctx, "logData", map[string]any{})
}

func getLogs(ctx context.Context) {
	items, err := database.LogDB.GetLogs()
	if err != nil {
		log.Println("❌ Ошибка при получении предметов:", err)
		return
	}
	var formattedItems []map[string]interface{}
	for _, item := range items {
		formattedItems = append(formattedItems, map[string]interface{}{
			"time":  item.Time,
			"user":  item.User,
			"value": item.Item,
		})
	}
	runtime.EventsEmit(ctx, "logData", formattedItems)
}

func manualRouletteSpin(a *App) {
	a.logic.ManualRouletteSpin()
}

func reconnectDonatty(a *App) {
	a.collManager.StopCollector("Donatty")
	a.collManager.StartCollector("Donatty")
}

func reconnectDonatepay(a *App) {
	if a.collManager.IsCollectorActive("DonatePay") {
		a.collManager.StopCollector("DonatePay")
	}
	a.collManager.StartCollector("DonatePay")
}

func reloadRoulette(a *App) {
	a.logic.ReloadRoulette()
}

func getSectorsByCategoryId(ctx context.Context, data string) {
	var payload struct {
		CategoryID int `json:"category_id"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Println("❌ Ошибка парсинга JSON:", err)
		return
	}
	sectors, err := database.WidgetDB.GetSectorsByCategoryID(payload.CategoryID)
	if err != nil {
		log.Println("❌ Ошибка при получении предметов:", err)
		return
	}

	var formattedSectors []map[string]interface{}
	for _, sector := range sectors {
		formattedSectors = append(formattedSectors, map[string]interface{}{
			"id":     sector.ID,
			"data":   sector.Name,
			"status": nil,
		})
	}
	runtime.EventsEmit(ctx, "SectorsByCategoryIdData", formattedSectors)
}

func sectorsToSave(ctx context.Context, data string) {
	var payload struct {
		CategoryID int `json:"id"`
		Sectors    []struct {
			ID     int     `json:"id"`
			Data   string  `json:"data"`
			Status *string `json:"status"` // может быть null
		} `json:"sectors"`
	}

	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Println("❌ Ошибка парсинга JSON sectorsToSave:", err)
		return
	}
	log.Println(payload)

	for _, sector := range payload.Sectors {
		switch {
		case sector.Status == nil:
			continue

		case *sector.Status == "add":
			err := database.WidgetDB.AddSector(payload.CategoryID, sector.Data)
			if err != nil {
				log.Printf("❌ Ошибка добавления: %v", err)
			}

		case *sector.Status == "edit":
			err := database.WidgetDB.UpdateSector(sector.ID, sector.Data)
			if err != nil {
				log.Printf("❌ Ошибка обновления: %v", err)
			}

		case *sector.Status == "delete":
			err := database.WidgetDB.DeleteSector(sector.ID)
			if err != nil {
				log.Printf("❌ Ошибка удаления: %v", err)
			}

		default:
			log.Printf("⚠️ Неизвестный статус '%v' для элемента ID %d", *sector.Status, sector.ID)
		}
	}

	sectors, err := database.WidgetDB.GetSectorsByCategoryID(payload.CategoryID)
	if err != nil {
		log.Println("❌ Ошибка при повторном получении предметов:", err)
		return
	}

	var formattedSectors []map[string]interface{}
	for _, sector := range sectors {
		formattedSectors = append(formattedSectors, map[string]interface{}{
			"id":     sector.ID,
			"data":   sector.Name,
			"status": nil,
		})
	}
	runtime.EventsEmit(ctx, "sectorsByCategoryIdData", formattedSectors)
}

func getSectors(ctx context.Context) {
	sectors, err := database.WidgetDB.GetRouletteCategorys()
	if err != nil {
		log.Println("❌ Ошибка при получении предметов:", err)
		return
	}
	result := make([]map[string]interface{}, 0)

	for _, sector := range sectors {
		sectors, err := database.WidgetDB.GetSectorsByCategoryID(sector.ID)
		if err != nil {
			log.Printf("❌ Ошибка при получении предметов для группы %d: %s", sector.ID, err)
			continue // пропускаем группу, если что-то пошло не так
		}

		sectorNames := make([]string, 0, len(sectors))
		for _, sector := range sectors {
			sectorNames = append(sectorNames, sector.Name)
		}

		sectorData := map[string]interface{}{
			"title":      sector.Name,
			"sectors":    sectorNames,
			"percentage": sector.Percentage,
			"color":      sector.Color,
		}
		result = append(result, sectorData)
		log.Println("✅ Группы:", result)
		// Отправляем на фронт
		runtime.EventsEmit(ctx, "sectorsData", result)
	}
}

func getSettings(ctx context.Context) {
	settings, err := database.CredentialsDB.GetAllENVValues()
	if err != nil {
		log.Println("❌ Ошибка при получении настроек:", err)
		return
	}
	result := make([]map[string]interface{}, 0)

	for _, setting := range settings {
		settingsData := map[string]interface{}{
			"name":  setting.Name,
			"value": setting.Value,
		}
		result = append(result, settingsData)
	}
	runtime.EventsEmit(ctx, "SettingsData", result)
}

func updateSettings(a *App, data string) {
	var payload struct {
		Settings []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"settings"`
	}

	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Println("❌ Ошибка парсинга JSON updateSettings:", err)
		return
	}
	for _, setting := range payload.Settings {
		exists, err := database.CredentialsDB.CheckENVExists(setting.Name)
		if err != nil {
			log.Printf("❌ Ошибка проверки существования настройки '%s': %v", setting.Name, err)
			continue
		}

		if exists {
			database.CredentialsDB.UpdateENVValue(setting.Name, setting.Value)
		} else {
			database.CredentialsDB.InsertENVValue(setting.Name, setting.Value)
		}
	}
	if a.collManager.IsActive() {
		reconnectAllCollector(a)
	}

	toastData := map[string]interface{}{
		"message": "Данные сохранены",
		"type":    "success",
	}
	runtime.EventsEmit(a.ctx, "toastExec", toastData)
}

func startCollector(data string, a *App) {
	var payload struct {
		Collector string `json:"collectorName"`
	}

	// парсим JSON-строку
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Println("❌ Ошибка парсинга JSON startCollector:", err)
		return
	}

	collectorName := payload.Collector
	log.Println("🔁 Получен запрос на переключение коллектора:", collectorName)

	a.collManager.StartCollector(collectorName)

}

func startAllCollector(a *App) {

	log.Println("🔁 Получен запрос на включение всех коллекторов:")

	a.collManager.StartAllCollector()

}

func stopAllCollector(a *App) {

	log.Println("🔁 Получен запрос на выключение всех коллекторов:")

	a.collManager.StopAllCollector()

}

func reconnectAllCollector(a *App) {

	log.Println("🔁 Получен запрос на переподключение всех коллекторов:")

	a.collManager.StopAllCollector()
	a.collManager.StartAllCollector()

}
