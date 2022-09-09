package cortex

type Config struct {
	Address                              string `yaml:"address" mapstructure:"address" default:"http://localhost:8080"`
	PrometheusAlertManagerConfigYaml     string `yaml:"-" mapstructure:"configyaml" default:""`
	PrometheusAlertManagerHelperTemplate string `yaml:"-" mapstructure:"helpertemplate" default:""`
}
