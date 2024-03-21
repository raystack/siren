package template

import (
	"gopkg.in/yaml.v2"
)

func YamlStringToFile(str string) (*File, error) {
	var t File
	err := yaml.Unmarshal([]byte(str), &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ParseFile(fl *File) (*Template, error) {
	bodyStr, err := yaml.Marshal(fl.Body)
	if err != nil {
		return nil, err
	}
	return &Template{
		Name:      fl.Name,
		Body:      string(bodyStr),
		Tags:      fl.Tags,
		Variables: fl.Variables,
	}, nil
}

func RulesBody(tpl *Template) ([]Rule, error) {
	var rules []Rule

	if err := yaml.Unmarshal([]byte(tpl.Body), &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func MessagesFromBody(tpl *Template) ([]Message, error) {
	var messages []Message

	if err := yaml.Unmarshal([]byte(tpl.Body), &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
