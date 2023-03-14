package httpreceiver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/goto/salt/log"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/httpclient"
	"github.com/goto/siren/pkg/retry"
	"github.com/goto/siren/plugins/receivers/base"
	"github.com/mitchellh/mapstructure"
)

type PluginService struct {
	base.UnimplementedService
	httpClient *httpclient.Client
	retrier    retry.Runner
	logger     log.Logger
}

func NewPluginService(logger log.Logger, cfg AppConfig, opts ...ServiceOption) *PluginService {
	s := &PluginService{}

	for _, opt := range opts {
		opt(s)
	}

	s.logger = logger

	if s.httpClient == nil {
		s.httpClient = httpclient.New(cfg.HTTPClient)
	}

	return s
}

func (s *PluginService) PreHookDBTransformConfigs(ctx context.Context, receiverConfigMap map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(receiverConfigMap, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, errors.ErrInvalid.WithMsgf(err.Error())
	}

	return receiverConfig.AsMap(), nil
}

func (s *PluginService) PreHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to slack notification config: %w", err)
	}

	if err := notificationConfig.Validate(); err != nil {
		return nil, err
	}

	return notificationConfig.AsMap(), nil
}

func (s *PluginService) Send(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return false, err
	}

	bodyBytes, err := json.Marshal(notificationMessage.Details)
	if err != nil {
		return false, err
	}

	if err := s.Notify(ctx, notificationConfig.URL, bodyBytes); err != nil {
		if errors.As(err, new(retry.RetryableError)) {
			return true, err
		} else {
			return false, err
		}
	}

	return false, nil
}

func (s *PluginService) Notify(ctx context.Context, apiURL string, body []byte) error {
	if s.retrier != nil {
		if err := s.retrier.Run(ctx, func(ctx context.Context) error {
			return s.notify(ctx, apiURL, body)
		}); err != nil {
			return err
		}
	}
	return s.notify(ctx, apiURL, body)
}

func (s *PluginService) notify(ctx context.Context, apiURL string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request body: %w", err)
	}

	resp, err := s.httpClient.HTTP().Do(req)
	if err != nil {
		return retry.RetryableError{Err: fmt.Errorf("failure in http call: %w", err)}
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
		s.logger.Info("httpreceiver call success", "url", apiURL, "response", string(bodyBytes))
	}

	return nil
}
