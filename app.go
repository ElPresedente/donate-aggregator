package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-back/database"
	"go-back/logic"
	"go-back/sources"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"net/http"

	"github.com/gorilla/websocket"
)

//ЯРЧЕ: Мужики, сори, не ебу куда это пихнуть, сами разберётесь---------------------------------------------------------------------------------------------------------------

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// Объект для апгрейда HTTP соединения до WebSocket
func (a *App) StartWebSocketServer() {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Ошибка апгрейда:", err)
			return
		}
		a.clientsMu.Lock()
		a.clients[conn] = true
		a.clientsMu.Unlock()

		defer func() {
			a.clientsMu.Lock()
			delete(a.clients, conn)
			a.clientsMu.Unlock()
			conn.Close()
		}()

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	})

	go http.ListenAndServe(":8080", nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Апгрейд соединения
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка апгрейда:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Новое подключение WebSocket")

	for {
		// Чтение сообщения от клиента
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Ошибка чтения:", err)
			break
		}

		fmt.Printf("Получено сообщение: %s\n", message)

		// Ответ клиенту
		err = conn.WriteMessage(messageType, []byte("Принято: "+string(message)))
		if err != nil {
			fmt.Println("Ошибка записи:", err)
			break
		}
	}
}

func (a *App) SendMessageFromFrontend(msg string) {

	a.clientsMu.Lock()
	defer a.clientsMu.Unlock()

	payload := Message{
		Type:    "chat",
		Payload: msg,
	}
	data, _ := json.Marshal(payload)

	for conn := range a.clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Println("Ошибка при отправке:", err)
			conn.Close()
			delete(a.clients, conn)
		}
	}
}

//ЯРЧЕ: Мужики, сори, не ебу куда это пихнуть, сами разберётесь---------------------------------------------------------------------------------------------------------------

// App struct
type App struct {
	ctx context.Context
	//ЯРЧЕ: ЭТО ТОЖЕ МОЁ---------------------------------------------------------------------------------------------------------------
	clients   map[*websocket.Conn]bool
	clientsMu sync.Mutex
	//ЯРЧЕ: ЭТО ТОЖЕ МОЁ---------------------------------------------------------------------------------------------------------------
}

// NewApp creates a new App application struct
func NewApp() *App {
	//ЯРЧЕ: НЕ ЕБУ ЧО ЭТО ДЕЛАЕТ, НО ЭТО ТОЖЕ МОЁ---------------------------------------------------------------------------------------------------------------
	return &App{
		clients: make(map[*websocket.Conn]bool),
	}
	//ЯРЧЕ: НЕ ЕБУ ЧО ЭТО ДЕЛАЕТ, НО ЭТО ТОЖЕ МОЁ---------------------------------------------------------------------------------------------------------------
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.StartWebSocketServer()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %s", err)
	}

	// Создаём канал для событий
	eventCh := make(chan sources.DonationEvent, 100)

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
				// Отправка события в фронтенд (для будущего GUI) -------------------

				//ПО СУТИ ВОТ ТУТ МЫ БУДЕМ ЮЗАТЬ FrontendDispatcher
				runtime.EventsEmit(a.ctx, "donation", donation)
			}
		}
	}()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) FrontendDispatcher(endpoint string, argJSON string)  {
	log.Printf("🛰 Вызов FrontendDispatcher: %s, argJSON: %s", endpoint, argJSON)
	
	switch endpoint {
	case "test":
		logic.NotifyDBChange(a.ctx, []string{"попа", "жопа", "попа"})
	
	// Получение предметов по ID группы
	case "getItemsByGroup":
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
		runtime.EventsEmit(a.ctx, "groupItems", items)

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
	}
	log.Println("✅ Группы:", result)
	// Отправляем на фронт
		runtime.EventsEmit(a.ctx, "groupsData", result)
		//runtime.EventsEmit(a.ctx, "groupItems", items)

	// Добавление нового предмета в группу
	case "addItemToGroup":
		var payload struct {
			GroupID  int    `json:"group_id"`
			ItemName string `json:"item_name"`
		}
		if err := json.Unmarshal([]byte(argJSON), &payload); err != nil {
			log.Println("❌ Ошибка парсинга JSON:", err)
			return
		}
		err := database.RouletteDB.AddItemToGroup(payload.GroupID, payload.ItemName)
		if err != nil {
			log.Println("❌ Ошибка добавления предмета:", err)
			return
		}
		runtime.EventsEmit(a.ctx, "itemAdded", map[string]interface{}{
			"group_id": payload.GroupID,
			"name":     payload.ItemName,
		})
	case "getGroupById":
		var payload struct {
			GroupID int `json:"group_id"`
		}
		if err := json.Unmarshal([]byte(argJSON), &payload); err != nil {
			log.Println("⚠️ Ошибка парсинга JSON:", err)
			return
		}
		groupData, err := database.RouletteDB.GetGroupWithItemsByID(payload.GroupID)
		if err != nil {
			log.Printf("❌ Ошибка получения группы: %s", err)
			return
		}
		jsonData, err := json.Marshal(groupData)
		if err != nil {
			log.Printf("❌ Ошибка маршалинга JSON: %s", err)
			return
		}
		runtime.EventsEmit(a.ctx, "db_updated", string(jsonData))

	default:
		log.Printf("⚠️ Неизвестный endpoint: %s", endpoint)
	}
}
