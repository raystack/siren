package pagerduty

import _ "embed"

var (
	//go:embed config/default_alert_template_body_v1.goyaml
	defaultAlertTemplateBodyV1 string
)
