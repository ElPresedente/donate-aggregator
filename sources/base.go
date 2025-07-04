package sources

import (
	"context"
	"time"
)

// DonationEvent представляет нормализованные данные о донате
type DonationEvent struct {
	SourceID string  `json:"sourceId"`
	User     string  `json:"user"`
	Amount   float64 `json:"amount"`
	//Currency   string    `json:"currency"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Date      time.Time `json:"date"`
}

// EventCollector интерфейс для всех источников событий
type EventCollector interface {
	Start(ctx context.Context) error
	Stop() error
}
