package sources

import (
	"context"
)

type TwitchCollector struct {
}

func NewTwitchCollector() *TwitchCollector {
	return &TwitchCollector{}
}

func (dc *TwitchCollector) GetCollectorType() string {
	return "Twitch"
}

func (dc *TwitchCollector) Start(ctx context.Context) error {

	return nil
}
