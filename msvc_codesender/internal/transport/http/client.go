package http_client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HTTPTransport struct {
	Client *http.Client
}

func NewHTTPTransport() *HTTPTransport {
	return &HTTPTransport{
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: false,
				MaxIdleConns:      10,
				IdleConnTimeout:   30 * time.Second,
			},
		},
	}
}

func (t *HTTPTransport) Post(apiUrl string, payload map[string]string) (string, error) {
	// мы не передаем контекст в Post() специально.
	// Контекст из main нужен был чтоб в process_message мы перестали принимать новые сообшения если пришла отмена контекста
	// но здесь мы должны дать время на довыполнение работы, поэтому создаем новый контекст не связанный с контекстом из main
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// jsonData, err := json.Marshal(payload)
	// if err != nil {
	// 	return fmt.Errorf("marshal error: %w", err)
	// }

	// Преобразуем в form data
	form := url.Values{}
	for k, v := range payload {
		form.Set(k, v)
	}

	// request debug
	fmt.Printf("REQUEST to %s:\nHeaders: Content-Type: application/x-www-form-urlencoded\nBody: %s\n", apiUrl, form.Encode())
	// end debug

	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body failed: %w", err)
	}

	return string(bodyBytes), nil

}
