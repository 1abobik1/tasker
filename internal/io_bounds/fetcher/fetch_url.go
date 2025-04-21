package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type FetchURLPayload struct {
	Input string `json:"input"`
}

type FetchURLProcessor struct{}

func NewFetchURLProcessor() *FetchURLProcessor {
	return &FetchURLProcessor{}
}

func (p *FetchURLProcessor) Process(ctx context.Context, payload []byte) ([]byte, error) {
	time.Sleep(time.Second * 8) // имитация долгой работы
    
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, string(payload), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return content, nil
}
