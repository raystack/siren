package cortex

type Config struct {
	Address                              string `mapstructure:"address" default:"http://localhost:8080"`
	PrometheusAlertManagerConfigYaml     string `mapstructure:"address" default:""`
	PrometheusAlertManagerHelperTemplate string `mapstructure:"address" default:""`
}
