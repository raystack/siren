package slack

import _ "embed"

var (
	//go:embed config/default_cortex_alert_template_body.goyaml
	defaultCortexAlertTemplateBody string
)
