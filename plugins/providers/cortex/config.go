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
	Address string `yaml:"address" mapstructure:"address" default:"http://localhost:8080"`
}