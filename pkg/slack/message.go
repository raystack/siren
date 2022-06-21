package slack

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	goslack "github.com/slack-go/slack"
)

type Message struct {
	ReceiverName string         `json:"receiver_name" validate:"required"`
	ReceiverType string         `json:"receiver_type" validate:"required,oneof=user channel"`
	Message      string         `json:"message"`
	Blocks       goslack.Blocks `json:"blocks"`
}

func (sm *Message) Validate() error {
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

func (sm *Message) checkError(err error) error {
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
