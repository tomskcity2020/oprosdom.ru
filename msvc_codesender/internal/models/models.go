package models

import (
	"time"
)

type Config struct {
	WorkerInterval time.Duration
	MaxJitterMs    int
	Gateways       []GatewayConfig
}

type GatewayConfig struct {
	Name string
	URL  string
	Type string // "regular" или "premium"
	Auth map[string]string
}

type MsgFromRepo struct {
	Id    int
	Phone string
	Code  int
	Retry int
}
