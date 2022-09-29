package cortex

import _ "embed"

var (
	//go:embed config/helper.tmpl
	HelperTemplateString string
	//go:embed config/config.goyaml
	ConfigYamlString string
)

// Config is a cortex provider config
type Config struct {
	Address string `mapstructure:"address" default:"http://localhost:9009"`
}
