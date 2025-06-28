package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/r3labs/sse/v2"
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
	mainToken      string
	token          token
	ref            string
	reconnectDelay time.Duration
	client         *http.Client
	stop           chan struct{}
}

// NewDonattyCollector создаёт новый коллектор для Donatty
func NewDonattyCollector(token_str, ref string) *DonattyCollector {
	return &DonattyCollector{
		mainToken:      token_str,
		token:          token{},
		ref:            ref,
		reconnectDelay: 5 * time.Second,
		//pingInterval:   30 * time.Second,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		stop: make(chan struct{}),
	}
}

// Start запускает коллектор
func (dc *DonattyCollector) Start(ctx context.Context, ch chan<- DonationEvent) error {
	dc.getAccessToken()

	//ping секция
	go func() {
		lastPing := time.Now()
		for {
			select {
			case <-dc.stop:
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
				} else {
					time.Sleep(ping_interval)
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-dc.stop:
			return nil
		default:
			//check is token expired
			//if expired - getAccessToken()
			sseUrl := fmt.Sprintf("%s/widgets/%s/sse?zoneOffset=%d&jwt=%s", api_donatty_uri, dc.ref, zone_offset, dc.token.AccessToken)
			sseClient := sse.NewClient(sseUrl)

			err := sseClient.SubscribeRaw(func(msg *sse.Event) {
				var outer struct {
					Action string          `json:"action"`
					Data   json.RawMessage `json:"data"`
				}
				if err := json.Unmarshal(msg.Data, &outer); err != nil {
					log.Printf("❌ Ошибка парсинга события: %v", err)
					return
				}

				if outer.Action != "DATA" {
					//возможно где то тут надо отслеживать пришел ли пинг или нет, и на это реагировать
					return
				}

				//log.Printf("!!!! SSE EVENT %s, %s, %t", outer.Action, outer.Data, outer.Action != "DATA")

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
					SourceID: "donatty",
					User:     streamData.DisplayName,
					Amount:   wrapper.Amount,
					//Currency:   wrapper.Currency,
					Message:   wrapper.Message,
					Timestamp: time.Now(),
				}

				if donation.Amount == 0 {
					donation.Amount = streamData.Amount
				}
				// if donation.Currency == "" {
				// 	donation.Currency = streamData.Currency
				// }
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
				fmt.Printf("💸 Сумма: %.2f\n", donation.Amount /*, donation.Currency*/)
				fmt.Printf("📅 Дата: %s\n", donation.Date.Format("2006-01-02 15:04:05"))
				fmt.Printf("🕒 Время (локальное): %s\n", donation.Timestamp.Format("15:04:05"))
				fmt.Printf("----------------------------------------\n")

				// Отправка события в канал
				select {
				case ch <- donation:
				case <-ctx.Done():
					return
				}
			})
			if err != nil {
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
