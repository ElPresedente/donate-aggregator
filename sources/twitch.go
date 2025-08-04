package sources

import (
	"context"
	_ "embed"
	"strings"
)

//go:embed twitch_client_credentials.txt
var secrets string

var (
	clientID     string
	clientSecret string
)

type TwitchCollector struct {
}

func NewTwitchCollector() *TwitchCollector {
	lines := strings.Split(strings.TrimSpace(secrets), "\n")

	// Проверяем, что в файле ровно 2 строки
	if len(lines) != 2 {
		panic("Ошибка: файл twitch_client_credentials.txt должен содержать ровно 2 строки")
	}
	clientID = lines[0]
	clientSecret = lines[1]
	return &TwitchCollector{}
}

func (dc *TwitchCollector) GetCollectorType() string {
	return "Twitch"
}

func (dc *TwitchCollector) Start(ctx context.Context) error {

	return nil
}
