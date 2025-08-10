package gateway

import (
	"context"
	"errors"
	"fmt"

	"oprosdom.ru/msvc_codesender/internal/models"
	"oprosdom.ru/msvc_codesender/internal/repo"
	http_client "oprosdom.ru/msvc_codesender/internal/transport/http"
)

type Gateway struct {
	Name      string
	URL       string
	Type      string
	Transport *http_client.HTTPTransport
	Repo      repo.NoSqlInterface
	Config    map[string]string
}

func (g *Gateway) Send(ctx context.Context, msg models.MsgFromRepo) error {
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

	gatewayRepsonse, err := g.Transport.Post(g.URL, payload)
	if err!=nil {
		return fmt.Errorf("gateway http transport failed: %w", err)
	}

	g.Repo.LogResponse(ctx, g.Name, gatewayRepsonse)

	return nil
}
