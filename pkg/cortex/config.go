package cortex

type Config struct {
	Address                              string `mapstructure:"address" default:"http://localhost:8080"`
	PrometheusAlertManagerConfigYaml     string `mapstructure:"configyaml" default:""`
	PrometheusAlertManagerHelperTemplate string `mapstructure:"helpertemplate" default:""`
}
