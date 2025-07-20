package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/r3labs/sse/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var api_donatty_uri string = "http://api-013.donatty.com"

var zone_offset int32 = -180

var ping_interval time.Duration = 50 * time.Second

// DonattyCollector реализует коллектор для Donatty

type token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireAt     string `json:"expireAt"`
}

type DonattyCollector struct {
	ctx            context.Context
	mainToken      string
	token          token
	ref            string
	reconnectDelay time.Duration
	client         *http.Client
	eventChan      chan<- DonationEvent
	stop           chan struct{}
	sseCancel      context.CancelFunc
}

func (dc *DonattyCollector) setGUIState(state string) {
	runtime.EventsEmit(dc.ctx, "donattyConnectionUpdated", state)
}

// NewDonattyCollector создаёт новый коллектор для Donatty
func NewDonattyCollector(ctx context.Context, token_str, ref string, ch chan<- DonationEvent) *DonattyCollector {
	return &DonattyCollector{
		ctx:            ctx,
		mainToken:      token_str,
		token:          token{},
		ref:            ref,
		reconnectDelay: 5 * time.Second,
		//pingInterval:   30 * time.Second,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		eventChan: ch,
		stop:      make(chan struct{}),
		sseCancel: nil,
	}
}

func (dc *DonattyCollector) GetCollectorType() string {
	return "Donatty"
}

// Start запускает коллектор
func (dc *DonattyCollector) Start(ctx context.Context) error {
	dc.getAccessToken()
	dc.stop = make(chan struct{})

	// Создаем контекст с возможностью отмены для SSE
	sseCtx, sseCancel := context.WithCancel(ctx)
	dc.sseCancel = sseCancel
	defer sseCancel() // Гарантируем вызов отмены при выходе из функции

	// ping секция
	go func() {
		lastPing := time.Now()
		for {
			select {
			case <-dc.stop:
				log.Println("Donatty коллектор отключен (ping горутина)")
				return
			case <-ctx.Done():
				log.Println("Donatty коллектор отключен (контекст отменен)")
				return
			default:
				if time.Since(lastPing) > ping_interval {
					req, err := http.NewRequest("POST", fmt.Sprintf("https://api.donatty.com/widgets/%s/ping", dc.ref), nil)
					if err != nil {
						log.Printf("❌ Ошибка создания PONG-запроса: %v", err)
						continue
					}
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dc.token.AccessToken))
					resp, err := dc.client.Do(req)
					if err != nil {
						log.Printf("❌ Ошибка при отправке PONG: %v", err)
					} else {
						log.Println("📡 Отправлен PONG Donatty")
						resp.Body.Close()
					}
					lastPing = time.Now()
				}
				time.Sleep(ping_interval)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Donatty коллектор отключен (контекст)")
			sseCancel() // Прерываем SSE-подключение
			return ctx.Err()
		case <-dc.stop:
			log.Println("Donatty коллектор отключен (stop канал)")
			sseCancel() // Прерываем SSE-подключение
			return nil
		default:
			sseUrl := fmt.Sprintf("%s/widgets/%s/sse?zoneOffset=%d&jwt=%s", api_donatty_uri, dc.ref, zone_offset, dc.token.AccessToken)
			sseClient := sse.NewClient(sseUrl)

			dc.setGUIState(Connecting)

			sseClient.OnConnect(func(c *sse.Client) {
				dc.setGUIState(Connected)
			})

			sseClient.OnDisconnect(func(c *sse.Client) {
				dc.setGUIState(Disonnected)
			})

			// Подписываемся с использованием контекста
			err := sseClient.SubscribeRawWithContext(sseCtx, func(msg *sse.Event) {
				// Логика обработки сообщений (без изменений)
				var outer struct {
					Action string          `json:"action"`
					Data   json.RawMessage `json:"data"`
				}
				if err := json.Unmarshal(msg.Data, &outer); err != nil {
					log.Printf("❌ Ошибка парсинга события: %v", err)
					return
				}

				if outer.Action != "DATA" {
					return
				}

				var wrapper struct {
					StreamEventType string  `json:"streamEventType"`
					StreamEventData string  `json:"streamEventData"`
					Subscriber      string  `json:"subscriber"`
					Message         string  `json:"message"`
					Amount          float64 `json:"amount"`
					Currency        string  `json:"currency"`
				}
				if err := json.Unmarshal(outer.Data, &wrapper); err != nil {
					log.Printf("❌ Ошибка парсинга wrapper: %v", err)
					return
				}

				var streamData struct {
					DisplayName string  `json:"displayName"`
					Amount      float64 `json:"amount"`
					Currency    string  `json:"currency"`
					Message     string  `json:"message"`
					Date        string  `json:"date"`
				}
				if wrapper.StreamEventData != "" {
					if err := json.Unmarshal([]byte(wrapper.StreamEventData), &streamData); err != nil {
						log.Printf("❌ Ошибка парсинга streamEventData: %v", err)
						return
					}
				}

				donation := DonationEvent{
					SourceID:  "donatty",
					User:      streamData.DisplayName,
					Amount:    wrapper.Amount,
					Message:   wrapper.Message,
					Timestamp: time.Now(),
				}

				if donation.Amount == 0 {
					donation.Amount = streamData.Amount
				}
				if donation.Message == "" {
					donation.Message = streamData.Message
				}
				if streamData.Date != "" {
					date, err := time.Parse(time.RFC3339, streamData.Date)
					if err != nil {
						log.Printf("❌ Ошибка парсинга даты: %v", err)
					} else {
						donation.Date = date
					}
				}
				fmt.Printf("\n🎁 Донат через DONATTY:\n")
				fmt.Printf("👤 От: %s\n", donation.User)
				fmt.Printf("💬 Сообщение: %s\n", donation.Message)
				fmt.Printf("💸 Сумма: %.2f\n", donation.Amount)
				fmt.Printf("📅 Дата: %s\n", donation.Date.Format("2006-01-02 15:04:05"))
				fmt.Printf("🕒 Время (локальное): %s\n", donation.Timestamp.Format("15:04:05"))
				fmt.Printf("----------------------------------------\n")

				select {
				case dc.eventChan <- donation:
				case <-ctx.Done():
					return
				}
			})
			if err != nil {
				dc.setGUIState(Connecting)
				log.Printf("❌ Ошибка подключения к Donatty: %v", err)
				log.Printf("🔁 Переподключение через %v...", dc.reconnectDelay)
				time.Sleep(dc.reconnectDelay)
			}
		}
	}
}

// Stop останавливает коллектор
func (dc *DonattyCollector) Stop() error {
	close(dc.stop)
	if dc.sseCancel != nil {
		dc.sseCancel()
	}
	dc.setGUIState(Disonnected)
	return nil
}

// getAccessToken получает access token для Donatty
func (dc *DonattyCollector) getAccessToken() error {
	url := fmt.Sprintf("%s/auth/tokens/%s", api_donatty_uri, dc.mainToken)
	resp, err := dc.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Response struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
			ExpireAt     string `json:"expireAt"`
		} `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("DONATTY TOKEN ERROR")
		return err
	}
	dc.token.AccessToken = result.Response.AccessToken
	dc.token.RefreshToken = result.Response.RefreshToken
	dc.token.ExpireAt = result.Response.ExpireAt
	return nil
}
