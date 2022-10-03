package slack

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/pkg/errors"
	goslack "github.com/slack-go/slack"
)

// MessageGoSlack is the contract of goslack message
// Deprecated: we are going with doing http call directly
// for better handling
type MessageGoSlack struct {
	ReceiverName string         `json:"receiver_name" validate:"required"`
	ReceiverType string         `json:"receiver_type" validate:"required,oneof=user channel"`
	Message      string         `json:"message"`
	Blocks       goslack.Blocks `json:"blocks"`
}

// Validate checks whether the message is valid or not
func (sm *MessageGoSlack) Validate() error {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if sm.Message == "" && len(sm.Blocks.BlockSet) == 0 {
		return errors.New("non empty message or non zero length block is required")
	}

	return sm.checkError(v.Struct(sm))
}

func (sm *MessageGoSlack) checkError(err error) error {
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		errStrs := []string{}
		for _, e := range errs {
			if e.Tag() == "oneof" {
				errStrValue := fmt.Sprintf("error value %q", e.Value())
				if e.Field() != "" {
					errStrValue = errStrValue + fmt.Sprintf(" for key %q", e.Field())
				}
				errStrValue = errStrValue + fmt.Sprintf(" not recognized, only support %q", e.Param())
				errStrs = append(errStrs, errStrValue)
				continue
			}

			if e.Tag() == "required" {
				errStrs = append(errStrs, fmt.Sprintf("field %q is required", e.Field()))
				continue
			}

			errStrs = append(errStrs, e.Field())
		}

		return errors.New(strings.Join(errStrs, " and "))
	}
	return nil
}

func (sm *MessageGoSlack) FromNotificationMessage(nm notification.Message) error {
	if nm.Configs["channel_type"] == "" {
		sm.ReceiverType = DefaultChannelType
	}
	sm.ReceiverName = fmt.Sprintf("%v", nm.Configs["channel_name"])

	sm.Message = fmt.Sprintf("%v", nm.Detail["message"])

	blocks := goslack.Blocks{}
	if err := mapstructure.Decode(nm.Detail["blocks"], &blocks); err != nil {
		return err
	}
	sm.Blocks = blocks

	return nil
}
