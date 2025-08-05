package gateway

import (
	"errors"
	"fmt"

	"oprosdom.ru/msvc_codesender/internal/models"
	http_client "oprosdom.ru/msvc_codesender/internal/transport/http"
)

type Gateway struct {
	Name      string
	URL       string
	Type      string
	Transport *http_client.HTTPTransport
	Config    map[string]string
}

func (g *Gateway) Send(msg models.MsgFromRepo) error {
	var payload map[string]string

	switch g.Name {
	case "Zima1reg":
		payload = map[string]string{
			"phone":  msg.Phone,
			"msg":    fmt.Sprintf("OprosDom.ru Никому не сообщайте: %d", msg.Code),
			"device": g.Config["device"],
			"token":  g.Config["token"],
		}
	case "Zima2reg":
		payload = map[string]string{
			"phone":  msg.Phone,
			"msg":    fmt.Sprintf("OprosDom.ru Никому не сообщайте: %d", msg.Code),
			"device": g.Config["device"],
			"token":  g.Config["token"],
		}
	case "Zima3prem":
		payload = map[string]string{
			"phone":  msg.Phone,
			"msg":    fmt.Sprintf("OprosDom.ru Никому не сообщайте: %d", msg.Code),
			"device": g.Config["device"],
			"token":  g.Config["token"],
		}
	default:
		return errors.New("unknown gateway")
	}

	return g.Transport.Post(g.URL, payload)
}
