package slack

import _ "embed"

var (
	//go:embed config/default_alert_template_body.goyaml
	defaultAlertTemplateBody string
)
