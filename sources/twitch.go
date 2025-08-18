package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-back/services"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TwitchCollector struct {
	reconnectDelay time.Duration
	conn           *websocket.Conn
	stop           chan struct{}
	ctx            context.Context
	sessionID      string
	eventChan      chan<- CollectorEvent
}

func NewTwitchCollector(ctx context.Context, ch chan<- CollectorEvent) *TwitchCollector {
	return &TwitchCollector{
		reconnectDelay: 5 * time.Second,
		stop:           make(chan struct{}),
		ctx:            ctx,
		eventChan:      ch,
	}
}

func (tc *TwitchCollector) GetCollectorType() string {
	return "Twitch"
}

func (tc *TwitchCollector) Start(ctx context.Context) error {
	if res, err := services.TwitchHasAuth(); err == nil && !res {
		return fmt.Errorf("no twitch auth")
	}
	err := services.TwitchAuthIfNot()
	if err != nil {
		return err
	}
	for {
		err := tc.connectAndRun()
		if err != nil {
			log.Printf("twitch: connection error: %v", err)
		}

		select {
		case <-ctx.Done():
			log.Println("twitch коллектор отключен (контекст)")
			return ctx.Err()
		case <-tc.stop:
			return nil
		case <-time.After(tc.reconnectDelay):
			// retry
		}
	}
}

func (tc *TwitchCollector) Stop() error {
	close(tc.stop)
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}

func (tc *TwitchCollector) connectAndRun() error {
	tc.setGUIState(Disonnected)
	conn, _, err := websocket.DefaultDialer.Dial("wss://eventsub.wss.twitch.tv/ws", nil)
	tc.setGUIState(Connecting)
	if err != nil {
		return err
	}
	defer conn.Close()
	tc.conn = conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return err // reconnect
		}

		var base socketAnswerBase
		if err := json.Unmarshal(msg, &base); err != nil {
			log.Println("twitch: unmarshal error:", err)
			continue
		}

		switch base.Metadata.Type {
		case "session_welcome":
			var hello socketSessionHello
			if err := json.Unmarshal(base.Payload, &hello); err == nil {
				tc.sessionID = hello.Session.Id
				log.Println("twitch: session id:", tc.sessionID)

				// подписка на события через services
				// if err := services.TwitchSubscribeChatMessages(tc.sessionID); err != nil {
				// 	log.Println("twitch: subscribe chat error:", err)
				// }
				if err := services.TwitchSubscribeRewardRedemption(tc.sessionID); err != nil {
					log.Println("twitch: subscribe rewards error:", err)
				}
				tc.setGUIState(Connected)
			}
		case "session_reconnect":
			tc.setGUIState(Disonnected)
			return fmt.Errorf("twitch: session reconnect")
		case "notification":
			tc.handleEvent(base.Payload)
		case "session_keepalive":
		default:
			log.Println("twitch: unknown type:", base.Metadata.Type)
		}
	}
}

func (tc *TwitchCollector) handleEvent(payload json.RawMessage) {
	log.Printf("TWITCH PAYLOAD: %s", payload)
	var msg struct {
		Subscription struct {
			Type string `json:"type"`
		} `json:"subscription"`
		Event json.RawMessage `json:"event"`
	}
	if err := json.Unmarshal(payload, &msg); err != nil {
		log.Println("twitch: event parse error:", err)
		return
	}

	switch msg.Subscription.Type {
	case "channel.chat.message":
		tc.handleChatMessage(msg.Event)
	case "channel.channel_points_custom_reward_redemption.add":
		tc.handleRewardRedemption(msg.Event)
	default:
		log.Println("twitch: unhandled event type:", msg.Subscription.Type)
	}
}

func (tc *TwitchCollector) handleChatMessage(event json.RawMessage) {

}

func (tc *TwitchCollector) handleRewardRedemption(event json.RawMessage) {
	log.Printf("json: %s", event)
	var parsedEvent rewardEvent
	if err := json.Unmarshal(event, &parsedEvent); err != nil {
		log.Printf("Error parsing event JSON: %v", err)
		return
	}

	if parsedEvent.Reward.ID != "" && parsedEvent.Reward.Title != "Растя-я-яжка" {
		return
	}

	roulette := RouletteEvent{
		Name:        parsedEvent.UserName,
		SpinsAmount: 1,
	}
	resultEvent, err := NewCollectorEvent("RouletteEvent", &roulette)
	if err != nil {
		log.Printf("Ошибка создания ивента крутки: %v", err)
		return
	}
	select {
	case tc.eventChan <- resultEvent:
	case <-tc.ctx.Done():
		return
	}
}

func (tc *TwitchCollector) setGUIState(state string) {
	runtime.EventsEmit(tc.ctx, "twitchConnectionUpdated", state)
}

// ==== Базовые структуры ====

type socketAnswerMetadata struct {
	Id        string `json:"message_id"`
	Type      string `json:"message_type"`
	Timestamp string `json:"message_timestamp"`
}

type socketAnswerBase struct {
	Metadata socketAnswerMetadata `json:"metadata"`
	Payload  json.RawMessage      `json:"payload"`
}

type socketSessionHello struct {
	Session struct {
		Id string `json:"id"`
	} `json:"session"`
}

// Message represents the top-level structure of the JSON.
type rewardMessage struct {
	Subscription rewardSubscription `json:"subscription"`
	Event        rewardEvent        `json:"event"`
}

// Subscription contains details about the subscription.
type rewardSubscription struct {
	ID        string          `json:"id"`
	Status    string          `json:"status"`
	Type      string          `json:"type"`
	Version   string          `json:"version"`
	Condition rewardCondition `json:"condition"`
	Transport rewardTransport `json:"transport"`
	CreatedAt string          `json:"created_at"`
	Cost      int             `json:"cost"`
}

// Condition holds the condition fields.
type rewardCondition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
	RewardID          string `json:"reward_id"`
}

// Transport holds the transport details.
type rewardTransport struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id"`
}

// Event represents the event data.
type rewardEvent struct {
	BroadcasterUserID    string       `json:"broadcaster_user_id"`
	BroadcasterUserLogin string       `json:"broadcaster_user_login"`
	BroadcasterUserName  string       `json:"broadcaster_user_name"`
	ID                   string       `json:"id"`
	UserID               string       `json:"user_id"`
	UserLogin            string       `json:"user_login"`
	UserName             string       `json:"user_name"`
	UserInput            string       `json:"user_input"`
	Status               string       `json:"status"`
	RedeemedAt           string       `json:"redeemed_at"`
	Reward               rewardReward `json:"reward"`
}

// Reward contains reward details.
type rewardReward struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Prompt string `json:"prompt"`
	Cost   int    `json:"cost"`
}
