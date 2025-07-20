package sources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var api_donatepay_uri string = "https://donatepay.ru/api/v2"

// DonatePayCollector реализует коллектор для DonatePay
type DonatePayCollector struct {
	ctx            context.Context
	accessToken    string
	userID         string
	reconnectDelay time.Duration
	client         *centrifuge.Client
	eventChan      chan<- DonationEvent
	stop           chan struct{}
}

func (dc *DonatePayCollector) setGUIState(state string) {
	runtime.EventsEmit(dc.ctx, "donatepayConnectionUpdated", state)
}

// NewDonatePayCollector создаёт новый коллектор для DonatePay
func NewDonatePayCollector(ctx context.Context, accessToken, userID string, ch chan<- DonationEvent) *DonatePayCollector {
	return &DonatePayCollector{
		ctx:            ctx,
		accessToken:    accessToken,
		userID:         userID,
		reconnectDelay: 5 * time.Second,
		eventChan:      ch,
		stop:           make(chan struct{}),
	}
}

// ConnetionEventHandler обрабатывает события Centrifugo
type ConnetionEventHandler struct {
	dc *DonatePayCollector
}

// OnConnect вызывается при успешном подключении
func (c ConnetionEventHandler) OnConnect(client *centrifuge.Client, event centrifuge.ConnectEvent) {
	c.dc.setGUIState(Connected)
	log.Println("✅ Подключено к DonatePay Centrifugo")
}

// OnError вызывается при ошибке
func (c ConnetionEventHandler) OnError(client *centrifuge.Client, event centrifuge.ErrorEvent) {
	log.Printf("❌ Ошибка Centrifugo: %v", event)
}

// OnDisconnect вызывается при отключении
func (c ConnetionEventHandler) OnDisconnect(client *centrifuge.Client, event centrifuge.DisconnectEvent) {
	c.dc.setGUIState(Disonnected)
	log.Printf("🔌 Отключено от Centrifugo: %v", event)
}

// OnSubscribeSuccess вызывается при успешной подписке
func (c ConnetionEventHandler) OnSubscribeSuccess(sub *centrifuge.Subscription, event centrifuge.SubscribeSuccessEvent) {
	log.Printf("✅ Подписка на канал %s успешна (Recovered: %v, Resubscribed: %v)", sub.Channel(), event.Recovered, event.Resubscribed)
}

// OnSubscribeError вызывается при ошибке подписки
func (c ConnetionEventHandler) OnSubscribeError(sub *centrifuge.Subscription, event centrifuge.SubscribeErrorEvent) {
	c.dc.setGUIState(Disonnected)
	log.Printf("❌ Ошибка подписки на канал %s: %s", sub.Channel(), event.Error)
}

// OnUnsubscribe вызывается при отписке
func (c ConnetionEventHandler) OnUnsubscribe(sub *centrifuge.Subscription, event centrifuge.UnsubscribeEvent) {
	log.Printf("🔌 Отписка от канала %s", sub.Channel())
}

// PublishHandler обрабатывает сообщения о донатах
type PublishHandler struct {
	ctx context.Context
	ch  chan<- DonationEvent
}

// OnPublish обрабатывает публикации в канале
func (h PublishHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	// Логируем сырой JSON для отладки
	log.Printf("📩 CENTRIFUGO JSON from channel %s: %s", sub.Channel(), string(e.Data))

	// Парсинг сообщения
	var msg struct {
		Type         string `json:"type"`
		Notification struct {
			Type   string `json:"type"`
			UserID int    `json:"user_id"`
			Vars   string `json:"vars"` // Изменено на string
		} `json:"notification"`
	}
	if err := json.Unmarshal(e.Data, &msg); err != nil {
		log.Printf("❌ Ошибка парсинга сообщения: %v", err)
		return
	}
	if msg.Type != "event" || msg.Notification.Type != "donation" {
		log.Printf("ℹ️ Пропущено сообщение с type=%s, notification.type=%s", msg.Type, msg.Notification.Type)

		return // Игнорируем сообщения, не связанные с донатами
	}

	// Проверка vars на пустоту
	if msg.Notification.Vars == "" {
		log.Printf("❌ Поле vars пустое")
		return
	}

	// Парсинг vars
	var vars struct {
		Name      string  `json:"name"`
		Comment   string  `json:"comment"`
		Sum       float64 `json:"sum"`
		Currency  string  `json:"currency"` // Предполагается, если отсутствует
		Target    string  `json:"target"`
		VideoLink string  `json:"video_link"`
		VideoID   string  `json:"video_id"`
	}
	if err := json.Unmarshal([]byte(msg.Notification.Vars), &vars); err != nil {
		log.Printf("❌ Ошибка парсинга vars: %v (vars: %s)", err, msg.Notification.Vars)
		return
	}

	// Нормализация данных
	donation := DonationEvent{
		SourceID: "donatepay",
		User:     vars.Name, //fmt.Sprintf("donatepay-%d", msg.Notification.UserID),
		Amount:   vars.Sum,
		//Currency:   vars.Currency, //мб пригодится потом
		Message:   vars.Comment,
		Timestamp: time.Now(),
		Date:      time.Now(), // DonatePay не предоставляет дату
	}

	if donation.User == "" {
		//если нет нормального имени, будет временное (надеюсь нет)
		donation.User = fmt.Sprintf("donatepay-%d", msg.Notification.UserID)
	}

	//if donation.Currency == "" {
	//	donation.Currency = "RUB" // Предполагаем RUB
	//}

	// Декодирование Unicode-экранированных символов
	if decodedComment, err := decodeUnicode(vars.Comment); err == nil {
		donation.Message = decodedComment
	} else {
		log.Printf("❌ Ошибка декодирования сообщения: %v", err)
	}

	// Вывод в консоль
	fmt.Printf("\n🎁 Донат через DONATEPAY:\n")
	fmt.Printf("👤 От: %s\n", donation.User)
	fmt.Printf("💬 Сообщение: %s\n", donation.Message)
	fmt.Printf("💸 Сумма: %.2f\n", donation.Amount /*, donation.Currency*/)
	fmt.Printf("📅 Дата: %s\n", donation.Date.Format("2006-01-02 15:04:05"))
	fmt.Printf("🕒 Время (локальное): %s\n", donation.Timestamp.Format("15:04:05"))
	fmt.Printf("----------------------------------------\n")

	// Отправка события в канал
	select {
	case h.ch <- donation:
	case <-h.ctx.Done():
		return
	}
}

func (dc *DonatePayCollector) GetCollectorType() string {
	return "DonatePay"
}

// Start запускает коллектор
func (dc *DonatePayCollector) Start(ctx context.Context) error {
	dc.stop = make(chan struct{})
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-dc.stop:
			return nil
		default:
			dc.setGUIState(Connecting)
			log.Println("🔑 Получение токена подключения DonatePay...")
			token, err := dc.getConnectionToken()
			if err != nil {
				log.Printf("❌ Ошибка получения токена: %v", err)
				log.Printf("🔁 Переподключение через %v...", dc.reconnectDelay)
				time.Sleep(dc.reconnectDelay)
				continue
			}

			log.Println("🔌 Настройка DonatePay Centrifugo...")
			config := centrifuge.DefaultConfig()
			config.Name = "go"
			client := centrifuge.NewJsonClient("wss://centrifugo.donatepay.ru/connection/websocket", config)
			client.SetToken(token)
			dc.client = client

			// Обработка событий Centrifugo
			handler := ConnetionEventHandler{dc}
			client.OnConnect(handler)
			client.OnError(handler)
			client.OnDisconnect(handler)

			// Подписка на канал widgets:LastEvents#<userID>
			channel := fmt.Sprintf("widgets:LastEvents#%s", dc.userID)
			sub, err := client.NewSubscription(channel)
			if err != nil {
				log.Printf("❌ Ошибка создания подписки на канал %s: %v", channel, err)
				client.Close()
				time.Sleep(dc.reconnectDelay)
				continue
			}

			// Обработка событий подписки
			sub.OnSubscribeSuccess(handler)
			sub.OnSubscribeError(handler)
			sub.OnUnsubscribe(handler)

			// Обработка сообщений о донатах
			sub.OnPublish(PublishHandler{ctx: ctx, ch: dc.eventChan})

			// Подключение и подписка
			if err := client.Connect(); err != nil {
				log.Printf("❌ Ошибка подключения к Centrifugo: %v", err)
				client.Close()
				time.Sleep(dc.reconnectDelay)
				continue
			}
			if err := sub.Subscribe(); err != nil {
				log.Printf("❌ Ошибка подписки на канал %s: %v", channel, err)
				client.Close()
				time.Sleep(dc.reconnectDelay)
				continue
			}

			// Ожидание завершения
			select {
			case <-ctx.Done():
				client.Close()
				return ctx.Err()
			case <-dc.stop:
				client.Close()
				return nil
			}
		}
	}
}

// Stop останавливает коллектор
func (dc *DonatePayCollector) Stop() error {
	close(dc.stop)
	dc.setGUIState(Disonnected)
	if dc.client != nil {
		dc.client.Close()
	}
	return nil
}

// getConnectionToken получает токен подключения к Centrifugo
func (dc *DonatePayCollector) getConnectionToken() (string, error) {
	url := fmt.Sprintf("%s/socket/token", api_donatepay_uri)
	payload, _ := json.Marshal(map[string]string{"access_token": dc.accessToken})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса для токена: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса для токена: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа токена: %v", err)
	}

	if result.Token == "" {
		return "", fmt.Errorf("получен пустой токен")
	}
	return result.Token, nil
}

// decodeUnicode декодирует Unicode-экранированные символы
func decodeUnicode(s string) (string, error) {
	var result strings.Builder
	var tempRune rune
	var err error
	for i := 0; i < len(s); {
		if i+1 < len(s) && s[i] == '\\' && s[i+1] == 'u' {
			if i+5 < len(s) {
				var code uint32
				_, err = fmt.Sscanf(s[i+2:i+6], "%04x", &code)
				if err != nil {
					return "", fmt.Errorf("ошибка декодирования Unicode: %v", err)
				}
				tempRune = rune(code)
				result.WriteRune(tempRune)
				i += 6
			} else {
				return "", fmt.Errorf("неполная Unicode-последовательность")
			}
		} else {
			tempRune, _, err = readRune(s, i)
			if err != nil {
				return "", fmt.Errorf("ошибка чтения руны: %v", err)
			}
			result.WriteRune(tempRune)
			i += runeLen(tempRune)
		}
	}
	return result.String(), nil
}

// readRune читает одну руну из строки
func readRune(s string, i int) (rune, int, error) {
	if i >= len(s) {
		return 0, 0, fmt.Errorf("индекс вне диапазона")
	}
	for _, r := range s[i:] {
		return r, len(string(r)), nil
	}
	return 0, 0, fmt.Errorf("пустая строка")
}

// runeLen возвращает длину руны в байтах
func runeLen(r rune) int {
	return len(string(r))
}
