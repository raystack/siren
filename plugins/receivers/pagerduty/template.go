package pagerduty

import _ "embed"

var (
	//go:embed config/default_cortex_alert_template_body_v1.goyaml
	DefaultCortexAlertTemplateBodyV1 string
)
