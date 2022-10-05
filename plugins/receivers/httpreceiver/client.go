package httpreceiver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/pkg/httpclient"
	"github.com/odpf/siren/pkg/retry"
)

type ClientOption func(*Client)

// ClientWithHTTPClient assigns custom client when creating a httpreceiver client
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
	logger     log.Logger
	httpClient *httpclient.Client
	retrier    retry.Runner
}

func NewClient(logger log.Logger, cfg AppConfig, opts ...ClientOption) *Client {
	c := &Client{
		cfg:    cfg,
		logger: logger,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = httpclient.New(httpclient.Config{})
	}

	return c
}

func (c *Client) Notify(ctx context.Context, apiURL string, body []byte) error {
	if c.retrier != nil {
		if err := c.retrier.Run(ctx, func(ctx context.Context) error {
			return c.notify(ctx, apiURL, body)
		}); err != nil {
			return err
		}
	}
	return c.notify(ctx, apiURL, body)
}

func (c *Client) notify(ctx context.Context, apiURL string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
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
		c.logger.Info("httpreceiver call success", "url", apiURL, "response", string(bodyBytes))
	}

	return nil
}
