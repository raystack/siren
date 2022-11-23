package telemetry

type Config struct {
	// Debug sets the bind address for pprof & zpages server.
	Debug string `mapstructure:"debug_addr" yaml:"debug_addr" default:"localhost:8081"`

	// OpenCensus trace & metrics configurations.
	EnableCPU        bool    `mapstructure:"enable_cpu" yaml:"enable_cpu" default:"true"`
	EnableMemory     bool    `mapstructure:"enable_memory" yaml:"enable_memory" default:"true"`
	SamplingFraction float64 `mapstructure:"sampling_fraction" yaml:"sampling_fraction" default:"1"`

	// OpenCensus exporter configurations.
	ServiceName string `mapstructure:"service_name" yaml:"service_name" default:"siren"`

	// NewRelic exporter.
	EnableNewrelic  bool   `mapstructure:"enable_newrelic" yaml:"enable_newrelic" default:"false"`
	NewRelicAppName string `mapstructure:"newrelic_app_name" yaml:"newrelic_app_name"`
	NewRelicAPIKey  string `mapstructure:"newrelic_api_key" yaml:"newrelic_api_key" default:"____LICENSE_STRING_OF_40_CHARACTERS_____"`

	// OpenTelemetry Agent exporter.
	EnableOtelAgent  bool   `mapstructure:"enable_otel_agent" yaml:"enable_otel_agent" default:"false"`
	OpenTelAgentAddr string `mapstructure:"otel_agent_addr" yaml:"otel_agent_addr" default:"localhost:8088"`
}
