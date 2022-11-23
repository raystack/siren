package pagerduty

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/odpf/siren/pkg/httpclient"
	"github.com/odpf/siren/pkg/retry"
)

const (
	defaultPagerdutyHost = "https://events.pagerduty.com"
)

type eventsV1HTTPResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	IncidentKey string `json:"incident_key"`
}

type ClientOption func(*Client)

// ClientWithHTTPClient assigns custom client when creating a pagerduty client
func ClientWithHTTPClient(cli *httpclient.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = cli
	}
}

// ClientWithRetrier wraps client call with retrier
func ClientWithRetrier(runner retry.Runner) ClientOption {
	return func(c *Client) {
		c.retrier = runner
	}
}

type Client struct {
	cfg        AppConfig
	httpClient *httpclient.Client
	retrier    retry.Runner
	// httpClientTracer *telemetry.HTTPClientSpan
}

func NewClient(cfg AppConfig, opts ...ClientOption) *Client {
	c := &Client{
		cfg: cfg,
	}

	for _, opt := range opts {
		opt(c)
	}

	if cfg.APIHost == "" {
		c.cfg.APIHost = defaultPagerdutyHost
	}

	if c.httpClient == nil {
		c.httpClient = httpclient.New(cfg.HTTPClient)
	}

	return c
}

func (c *Client) NotifyV1(ctx context.Context, msg MessageV1) error {
	if c.retrier != nil {
		if err := c.retrier.Run(ctx, func(ctx context.Context) error {
			return c.notifyV1(ctx, msg)
		}); err != nil {
			return err
		}
	}
	return c.notifyV1(ctx, msg)
}

func (c *Client) notifyV1(ctx context.Context, message MessageV1) error {
	// TODO need to sanitize body first?
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.APIHost+"/generic/2010-04-15/create_event.json", bytes.NewReader(messageJSON))
	if err != nil {
		return fmt.Errorf("failed to create request body: %w", err)
	}

	resp, err := c.httpClient.HTTP().Do(req)
	if err != nil {
		return fmt.Errorf("failure in http call: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 || resp.StatusCode >= 500 {
		return retry.RetryableError{Err: errors.New(http.StatusText(resp.StatusCode))}
	}

	if resp.StatusCode >= 300 {
		return errors.New(http.StatusText(resp.StatusCode))
	} else {
		// Status code 2xx only
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		apiResponse := eventsV1HTTPResponse{}
		if err = json.Unmarshal(bodyBytes, &apiResponse); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		if apiResponse.Status != "success" {
			return fmt.Errorf("something wrong when sending pagerduty event: %v", apiResponse)
		}
	}

	return nil
}
