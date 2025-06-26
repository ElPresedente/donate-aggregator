package sources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/golang-jwt/jwt/v5"
)

// DonatePayCollector реализует коллектор для DonatePay
type DonatePayCollector struct {
	accessToken    string
	userID         string
	reconnectDelay time.Duration
	client         *centrifuge.Client
	stop           chan struct{}
}

// NewDonatePayCollector создаёт новый коллектор для DonatePay
func NewDonatePayCollector(accessToken, userID string) *DonatePayCollector {
	return &DonatePayCollector{
		accessToken:    accessToken,
		userID:         userID,
		reconnectDelay: 5 * time.Second,
		stop:           make(chan struct{}),
	}
}

type ConnetionEventhandler struct{}

func (c ConnetionEventhandler) OnConnect(client *centrifuge.Client, event centrifuge.ConnectEvent) {
	log.Println("✅ Подключено к DonatePay Centrifugo")
}

func (c ConnetionEventhandler) OnError(client *centrifuge.Client, event centrifuge.ErrorEvent) {
	log.Printf("❌ Ошибка Centrifugo: %v", event)
}

func (c ConnetionEventhandler) OnDisconnect(client *centrifuge.Client, event centrifuge.DisconnectEvent) {
	log.Printf("🔌 Отключено от Centrifugo: %v", event)
}

func (c ConnetionEventhandler) OnSubscribeError(sub *centrifuge.Subscription, event centrifuge.SubscribeErrorEvent) {
	log.Printf("❌ Subscribing on channel %s - %s", sub.Channel(), event.Error)
}

func (c ConnetionEventhandler) OnSubscribeSuccess(sub *centrifuge.Subscription, event centrifuge.SubscribeSuccessEvent) {
	log.Printf("Subscribed on channel %s, (was Recovered: %v, Resubscribed: %v)", sub.Channel(), event.Recovered, event.Resubscribed)
}

func (c ConnetionEventhandler) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	log.Printf("Unsubscribed from channel %s", sub.Channel())
}

func (c ConnetionEventhandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	log.Printf("CENTRIFUGO JSON from channel %s\n %s", sub.Channel(), string(e.Data))
}

func (c ConnetionEventhandler) OnPrivateSub(client *centrifuge.Client, e centrifuge.PrivateSubEvent) (string, error) {
	log.Printf("PrivateSub channel - %s; ClientID - %s", e.Channel, e.ClientID)

	claims := jwt.MapClaims{
		"client":  e.ClientID,
		"channel": e.Channel,
		"exp":     time.Now().Add(time.Hour * 48).Unix(), // по желанию: срок действия
	}

	// Секрет
	secret := []byte("secret")

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подпись токена
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Fatalf("Ошибка создания токена: %v", err)
		return "", err
	}
	return tokenString, err
}

// Start запускает коллектор
func (dc *DonatePayCollector) Start(ctx context.Context, ch chan<- DonationEvent) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-dc.stop:
			return nil
		default:
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
			config.Name = "js"
			client := centrifuge.NewJsonClient(
				"wss://centrifugo.donatepay.ru/connection/websocket",
				config,
			)
			client.SetToken(token)
			dc.client = client

			// Обработка событий Centrifugo

			handler := ConnetionEventhandler{}

			client.OnConnect(handler)
			client.OnError(handler)
			client.OnDisconnect(handler)
			client.OnPrivateSub(handler)

			// Подписка на канал //:= fmt.Sprintf("notifications#%s", dc.userID)
			events := "events:events#1393285"
			notifications := "notifications#1393285"
			settings := "settings#1393285"
			widgets_AB := "widgets:AlertBox#1393285"
			widgets_LE := "widgets:LastEvents#1393285"

			sub_e, err := client.NewSubscription(events)
			if err != nil {
				log.Fatalln(err)
			}

			sub_e.OnSubscribeError(handler)
			sub_e.OnSubscribeSuccess(handler)
			sub_e.OnUnsubscribe(handler)

			sub_e.OnPublish(handler)

			sub_n, err := client.NewSubscription(notifications)
			if err != nil {
				log.Fatalln(err)
			}

			sub_n.OnSubscribeError(handler)
			sub_n.OnSubscribeSuccess(handler)
			sub_n.OnUnsubscribe(handler)

			sub_n.OnPublish(handler)

			sub_s, err := client.NewSubscription(settings)
			if err != nil {
				log.Fatalln(err)
			}

			sub_s.OnSubscribeError(handler)
			sub_s.OnSubscribeSuccess(handler)
			sub_s.OnUnsubscribe(handler)

			sub_s.OnPublish(handler)

			sub_wab, err := client.NewSubscription(widgets_AB)
			if err != nil {
				log.Fatalln(err)
			}

			sub_wab.OnSubscribeError(handler)
			sub_wab.OnSubscribeSuccess(handler)
			sub_wab.OnUnsubscribe(handler)

			sub_wab.OnPublish(handler)

			sub_wle, err := client.NewSubscription(widgets_LE)
			if err != nil {
				log.Fatalln(err)
			}

			sub_wle.OnSubscribeError(handler)
			sub_wle.OnSubscribeSuccess(handler)
			sub_wle.OnUnsubscribe(handler)

			sub_wle.OnPublish(handler)

			// Подключение
			if err := client.Connect(); err != nil {
				log.Printf("❌ Ошибка подключения к Centrifugo: %v", err)
				client.Close()
				time.Sleep(dc.reconnectDelay)
				continue
			} else {
				log.Printf("!!!!!! Подключено к Centrifugo")

			}

			sub_e.Subscribe()
			sub_n.Subscribe()
			sub_s.Subscribe()
			sub_wab.Subscribe()
			sub_wle.Subscribe()

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
	if dc.client != nil {
		dc.client.Close()
	}
	return nil
}

// getConnectionToken получает токен подключения к Centrifugo
func (dc *DonatePayCollector) getConnectionToken() (string, error) {
	url := "https://donatepay.ru/api/v2/socket/token"
	payload, _ := json.Marshal(map[string]string{"access_token": dc.accessToken})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("!!!!!!1 Ошибка получения DONATEPAY токена")
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("!!!!!!2 Ошибка получения DONATEPAY токена")
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("!!!!!!3 Ошибка получения DONATEPAY токена")
		return "", err
	}
	fmt.Printf("!!!!!TOKEN %s", result.Token)

	return result.Token, nil
}
