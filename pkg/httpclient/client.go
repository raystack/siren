package httpclient

import (
	"net/http"
	"time"

	"github.com/raystack/siren/pkg/telemetry"
)

type ClientOpt func(*Client)

func WithHTTPClient(hc *http.Client) ClientOpt {
	return func(c *Client) {
		c.httpClient = hc
	}
}

type Client struct {
	cfg        Config
	httpClient *http.Client
}

func New(cfg Config, opts ...ClientOpt) *Client {
	c := &Client{}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		if c.cfg.MaxConnsPerHost != 0 {
			transport.MaxConnsPerHost = c.cfg.MaxConnsPerHost
		}
		if c.cfg.MaxIdleConns != 0 {
			transport.MaxIdleConns = c.cfg.MaxIdleConns
		}
		if c.cfg.MaxIdleConnsPerHost != 0 {
			transport.MaxIdleConnsPerHost = c.cfg.MaxIdleConnsPerHost
		}
		if c.cfg.IdleConnTimeoutMS != 0 {
			transport.IdleConnTimeout = time.Duration(c.cfg.IdleConnTimeoutMS)
		}

		c.httpClient = &http.Client{
			Transport: &telemetry.Transport{
				Base: transport,
			},
		}

		if c.cfg.TimeoutMS != 0 {
			c.httpClient.Timeout = time.Duration(c.cfg.TimeoutMS)
		}
	}

	return c
}

func (c *Client) HTTP() *http.Client {
	return c.httpClient
}
