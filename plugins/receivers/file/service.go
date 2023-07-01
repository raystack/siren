package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/raystack/siren/core/notification"
	"github.com/raystack/siren/pkg/errors"
	"github.com/raystack/siren/plugins/receivers/base"
)

type PluginService struct {
	base.UnimplementedService
}

// NewPluginService returns file receiver service struct. This service implement [receiver.Resolver] and [notification.Notifier] interface.
func NewPluginService() *PluginService {
	return &PluginService{}
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
	if err := s.validateFilePath(notificationConfig.URL); err != nil {
		return false, err
	}

	fileInstance, err := os.OpenFile(notificationConfig.URL, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return false, err
	}

	byteNewLine := []byte("\n")
	bodyBytes = append(bodyBytes, byteNewLine...)
	_, err = fileInstance.Write(bodyBytes)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (s *PluginService) validateFilePath(path string) error {
	dirs := strings.Split(path, "/")
	filename := dirs[len(dirs)-1]
	format := strings.Split(filename, ".")
	if len(format) != 2 {
		return fmt.Errorf("invalid filename for \"%s\"", path)
	}
	return nil
}
