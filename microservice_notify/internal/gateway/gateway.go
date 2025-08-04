package gateway

import (
	"context"
	"errors"
	"fmt"

	"oprosdom.ru/microservice_notify/internal/models"
	http_client "oprosdom.ru/microservice_notify/internal/transport/http"
)

type Gateway struct {
	Name      string
	URL       string
	Type      string
	Transport *http_client.HTTPTransport
	Config    map[string]string
}

func (g *Gateway) Send(ctx context.Context, msg models.SMSMessage) error {
	var payload interface{}

	switch g.Name {
	case "API1":
		payload = map[string]interface{}{
			"phone":   msg.PhoneNumber,
			"message": fmt.Sprintf("Ваш код: %d", msg.Code),
			"api_key": g.Config["api_key"],
		}
	case "API2":
		payload = map[string]interface{}{
			"to":    msg.PhoneNumber,
			"text":  fmt.Sprintf("Code: %d", msg.Code),
			"from":  "SMS_SERVICE",
			"token": g.Config["token"],
		}
	case "API3":
		payload = map[string]interface{}{
			"recipient": msg.PhoneNumber,
			"content":   fmt.Sprintf("Verification code: %d", msg.Code),
			"sender_id": "VERIFY",
			"auth_key":  g.Config["auth_key"],
		}
	default:
		return errors.New("unknown gateway")
	}

	return g.Transport.Post(ctx, g.URL, payload)
}
