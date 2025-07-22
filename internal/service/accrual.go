package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

const (
	defaultAccrualRequestTimeout = 5 * time.Second
	defaultAccrualRetryInterval  = 1 * time.Second
	maxAccrualRetries            = 3
)

type AccrualClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAccrualClient(baseURL string) *AccrualClient {
	return &AccrualClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultAccrualRequestTimeout,
		},
	}
}

func (c *AccrualClient) GetOrderAccrual(
	ctx context.Context,
	orderNumber string,
) (*models.OrderAccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)

	var resp *models.OrderAccrualResponse
	var lastErr error

	for i := 0; i < maxAccrualRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		httpResp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			time.Sleep(defaultAccrualRetryInterval)
			continue
		}
		defer httpResp.Body.Close()

		switch httpResp.StatusCode {
		case http.StatusOK:
			if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
				return nil, fmt.Errorf("failed to decode response: %w", err)
			}
			return resp, nil
		case http.StatusNoContent:
			return nil, nil
		case http.StatusTooManyRequests:
			retryAfter := defaultAccrualRetryInterval
			if retryHeader := httpResp.Header.Get("Retry-After"); retryHeader != "" {
				if seconds, err := time.ParseDuration(retryHeader + "s"); err == nil {
					retryAfter = seconds
				}
			}
			time.Sleep(retryAfter)
			continue
		default:
			return nil, fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
		}
	}

	return nil, lastErr
}
