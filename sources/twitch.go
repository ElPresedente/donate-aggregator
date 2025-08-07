package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"go-back/services"
	_ "go-back/services"
)

type ConnectionStatus = int

const (
	INIT      = 0
	CONNECTED = 1
)

type TwitchCollector struct {
	status ConnectionStatus
}

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

func NewTwitchCollector() *TwitchCollector {
	return &TwitchCollector{
		status: INIT,
	}
}

func (dc *TwitchCollector) GetCollectorType() string {
	return "Twitch"
}

func (dc *TwitchCollector) Start(ctx context.Context) error {
	if res, err := services.TwitchHasAuth(); err == nil || !res {
		return fmt.Errorf("twitch не авторизован")
	}

	return nil
}
