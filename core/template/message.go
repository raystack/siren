package template

import (
	"github.com/goto/siren/pkg/errors"
)

func MessageContentByReceiverType(messagesTemplate []Message, receiverType string) (string, error) {
	var messageTemplateMap = make(map[string]string)

	for _, msgTemplate := range messagesTemplate {
		messageTemplateMap[msgTemplate.ReceiverType] = msgTemplate.Content
	}

	content, ok := messageTemplateMap[receiverType]
	if !ok {
		errors.ErrInvalid.WithCausef("can't found template of receiver type %s", receiverType)
	}

	if content == "" {
		return "", errors.ErrInvalid.WithCausef("%s template is empty", receiverType)
	}

	return content, nil
}
