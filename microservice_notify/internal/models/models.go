package models

import (
	"time"
)

type Config struct {
	WorkerInterval time.Duration
	GatewayTimeout time.Duration
	MaxJitterMs    int
	Gateways       []GatewayConfig
}

type GatewayConfig struct {
	Name string
	URL  string
	Type string // "regular" или "premium"
	Auth map[string]string
}

type SMSMessage struct {
	ID          int
	PhoneNumber string
	Message     string
	Code        int
	Retry       int
	Status      string
	Gateway     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
