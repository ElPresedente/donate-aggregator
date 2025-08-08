package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-back/services"

	"github.com/gorilla/websocket"
)

type TwitchCollector struct {
	reconnectDelay time.Duration
	conn           *websocket.Conn
	stop           chan struct{}
	ctx            context.Context
	sessionID      string
}

func NewTwitchCollector(ctx context.Context) *TwitchCollector {
	return &TwitchCollector{
		reconnectDelay: 5 * time.Second,
		stop:           make(chan struct{}),
		ctx:            ctx,
	}
}

func (tc *TwitchCollector) Start() {
	for {
		err := tc.connectAndRun()
		if err != nil {
			log.Printf("twitch: connection error: %v", err)
		}

		select {
		case <-tc.stop:
			return
		case <-time.After(tc.reconnectDelay):
			// retry
		}
	}
}

func (tc *TwitchCollector) Stop() {
	close(tc.stop)
	if tc.conn != nil {
		tc.conn.Close()
	}
}

func (tc *TwitchCollector) connectAndRun() error {
	conn, _, err := websocket.DefaultDialer.Dial("wss://eventsub.wss.twitch.tv/ws", nil)
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
				if err := services.TwitchSubscribeChatMessages(tc.sessionID); err != nil {
					log.Println("twitch: subscribe chat error:", err)
				}
				if err := services.TwitchSubscribeRewardRedemption(tc.sessionID); err != nil {
					log.Println("twitch: subscribe rewards error:", err)
				}
			}
		case "session_reconnect":
			return fmt.Errorf("twitch: session reconnect")
		case "notification":
			tc.handleEvent(base.Payload)
		default:
			log.Println("twitch: unknown type:", base.Metadata.Type)
		}
	}
}

func (tc *TwitchCollector) handleEvent(payload json.RawMessage) {
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
