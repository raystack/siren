package cortex

type Config struct {
	Address string `mapstructure:"address" default:"http://localhost:8080"`
}
