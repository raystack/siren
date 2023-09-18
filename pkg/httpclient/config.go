package httpclient

type Config struct {
	TimeoutMS           int `mapstructure:"timeout_ms" json:"timeout_ms" yaml:"timeout_ms"`
	MaxConnsPerHost     int `mapstructure:"max_conns_per_host" json:"max_conns_per_host" yaml:"max_conns_per_host"`
	MaxIdleConns        int `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxIdleConnsPerHost int `mapstructure:"max_idle_conns_per_host" json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host"`
	IdleConnTimeoutMS   int `mapstructure:"idle_conn_timeout_ms" json:"idle_conn_timeout_ms" yaml:"idle_conn_timeout_ms"`
}
