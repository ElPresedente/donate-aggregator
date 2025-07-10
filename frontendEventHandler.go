package main

import (
	"encoding/json"
	"go-back/database"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) FrontendDispatcher(endpoint string, argJSON string) {
	log.Printf("🛰 Вызов FrontendDispatcher: %s, argJSON: %s", endpoint, argJSON)

	switch endpoint {
	// Получение предметов по ID группы
	case "getItemsByGroupId":
		{ //to logic class
			var payload struct {
				GroupID int `json:"group_id"`
			}
			if err := json.Unmarshal([]byte(argJSON), &payload); err != nil {
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
			runtime.EventsEmit(a.ctx, "itemsByGroupIdData", formattedItems)
		}
	case "itemsToSave":
		{ //to logic class
			var payload struct {
				GroupID int `json:"id"` //Если потом произойдет логичный ренейм в групID, то тут тоже поменять
				Items   []struct {
					ID     int     `json:"id"`
					Data   string  `json:"data"`
					Status *string `json:"status"` // может быть null
				} `json:"items"`
			}

			if err := json.Unmarshal([]byte(argJSON), &payload); err != nil {
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
			runtime.EventsEmit(a.ctx, "itemsByGroupIdData", formattedItems)
		}

	// Получение всех групп и их итемов
	case "getGroups":
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
			runtime.EventsEmit(a.ctx, "groupsData", result)
		}

	case "getSettings":
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
		runtime.EventsEmit(a.ctx, "SettingsData", result)

	case "updateSettings":
		var payload struct {
			Settings []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"settings"`
		}

		if err := json.Unmarshal([]byte(argJSON), &payload); err != nil {
			log.Println("❌ Ошибка парсинга JSON updateSettings:", err)
			return
		}
		log.Println("Ы:", payload)
		for _, setting := range payload.Settings {
			exists, err := database.CredentialsDB.CheckENVExists(setting.Name)
			if err != nil {
				log.Printf("❌ Ошибка проверки существования настройки '%s': %v", setting.Name, err)
				continue
			}

			if exists {
				err = database.CredentialsDB.UpdateENVValue(setting.Name, setting.Value)
				if err != nil {
					log.Printf("❌ Ошибка обновления настройки '%s': %v", setting.Name, err)
				}
			} else {
				database.CredentialsDB.InsertENVValue(setting.Name, setting.Value)
			}
		}

	default:
		log.Printf("⚠️ Неизвестный endpoint: %s", endpoint)
	}
}
