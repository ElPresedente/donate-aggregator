package main

import (
	"context"
	"encoding/json"
	"go-back/database"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) FrontendDispatcher(endpoint string, argJSON string) {
	//log.Printf("🛰 Вызов FrontendDispatcher: %s, argJSON: %s", endpoint, argJSON)

	switch endpoint {
	// Получение предметов по ID группы
	case "getItemsByGroupId":
		getItemsByGroupId(a.ctx, argJSON)

	case "itemsToSave":
		itemsToSave(a.ctx, argJSON)

	// Получение всех групп и их итемов
	case "getGroups":
		getGroups(a.ctx)

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
	case "getNumLogs":
		getNumLogs(a.ctx, argJSON)
	case "newStream":
		newStream(a)
	default:
		log.Printf("⚠️ Неизвестный endpoint: %s", endpoint)
	}
}

func newStream(a *App) {
	database.LogDB.ClearDatabase()
	runtime.EventsEmit(a.ctx, "logNumData", map[string]any{})
}

func getNumLogs(ctx context.Context, data string) {
	var num int
	if err := json.Unmarshal([]byte(data), &num); err != nil {
		log.Println("❌ Ошибка парсинга JSON:", err)
		return
	}
	items, err := database.LogDB.GetLastNLogs(num)
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
	runtime.EventsEmit(ctx, "logNumData", formattedItems)
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

func getItemsByGroupId(ctx context.Context, data string) {
	var payload struct {
		GroupID int `json:"group_id"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Println("❌ Ошибка парсинга JSON:", err)
		return
	}
	items, err := database.RouletteDB.GetItemsByGroupID(payload.GroupID)
	if err != nil {
		log.Println("❌ Ошибка при получении предметов:", err)
		return
	}

	var formattedItems []map[string]interface{}
	for _, item := range items {
		formattedItems = append(formattedItems, map[string]interface{}{
			"id":     item.ID,
			"data":   item.Name,
			"status": nil,
		})
	}
	runtime.EventsEmit(ctx, "itemsByGroupIdData", formattedItems)
}

func itemsToSave(ctx context.Context, data string) {
	var payload struct {
		GroupID int `json:"id"` //Если потом произойдет логичный ренейм в групID, то тут тоже поменять
		Items   []struct {
			ID     int     `json:"id"`
			Data   string  `json:"data"`
			Status *string `json:"status"` // может быть null
		} `json:"items"`
	}

	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Println("❌ Ошибка парсинга JSON itemsToSave:", err)
		return
	}
	log.Println(payload)

	for _, item := range payload.Items {
		switch {
		case item.Status == nil:
			continue

		case *item.Status == "add":
			err := database.RouletteDB.AddItem(payload.GroupID, item.Data)
			if err != nil {
				log.Printf("❌ Ошибка добавления: %v", err)
			}

		case *item.Status == "edit":
			err := database.RouletteDB.UpdateItem(item.ID, item.Data)
			if err != nil {
				log.Printf("❌ Ошибка обновления: %v", err)
			}

		case *item.Status == "delete":
			err := database.RouletteDB.DeleteItem(item.ID)
			if err != nil {
				log.Printf("❌ Ошибка удаления: %v", err)
			}

		default:
			log.Printf("⚠️ Неизвестный статус '%v' для элемента ID %d", *item.Status, item.ID)
		}
	}

	items, err := database.RouletteDB.GetItemsByGroupID(payload.GroupID)
	if err != nil {
		log.Println("❌ Ошибка при повторном получении предметов:", err)
		return
	}

	var formattedItems []map[string]interface{}
	for _, item := range items {
		formattedItems = append(formattedItems, map[string]interface{}{
			"id":     item.ID,
			"data":   item.Name,
			"status": nil,
		})
	}
	runtime.EventsEmit(ctx, "itemsByGroupIdData", formattedItems)
}

func getGroups(ctx context.Context) {
	groups, err := database.RouletteDB.GetRouletteGroups()
	if err != nil {
		log.Println("❌ Ошибка при получении предметов:", err)
		return
	}
	result := make([]map[string]interface{}, 0)

	for _, group := range groups {
		items, err := database.RouletteDB.GetItemsByGroupID(group.ID)
		if err != nil {
			log.Printf("❌ Ошибка при получении предметов для группы %d: %s", group.ID, err)
			continue // пропускаем группу, если что-то пошло не так
		}

		itemNames := make([]string, 0, len(items))
		for _, item := range items {
			itemNames = append(itemNames, item.Name)
		}

		groupData := map[string]interface{}{
			"title":      group.Name,
			"items":      itemNames,
			"percentage": group.Percentage,
			"color":      group.Color,
		}
		result = append(result, groupData)
		log.Println("✅ Группы:", result)
		// Отправляем на фронт
		runtime.EventsEmit(ctx, "groupsData", result)
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
