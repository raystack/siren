package testdatatemplate_test

import (
	_ "embed"
)

var (
	//go:embed template-rule-sample-1.yaml
	SampleRuleTemplate string
	//go:embed template-message-sample-1.yaml
	SampleMessageTemplate string
)
