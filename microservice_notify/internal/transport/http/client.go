package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPTransport struct {
	Client *http.Client
}

func NewHTTPTransport(timeout time.Duration) *HTTPTransport {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &HTTPTransport{
		Client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DisableKeepAlives: false,
				MaxIdleConns:      10,
				IdleConnTimeout:   30 * time.Second,
			},
		},
	}
}

func (t *HTTPTransport) Post(ctx context.Context, url string, payload interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, t.Client.Timeout)
	defer cancel()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	return nil
}
