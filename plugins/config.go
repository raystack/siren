package plugins

type Config struct {
	PluginPath string                  `mapstructure:"plugin_path" yaml:"plugin_path" json:"plugin_path" default:"./plugin"`
	Plugins    map[string]PluginConfig `mapstructure:"plugins" yaml:"plugins" json:"plugins"`
}

type PluginConfig struct {
	Handshake     HandshakeConfig        `mapstructure:"handshake" yaml:"handshake" json:"handshake"`
	ServiceConfig map[string]interface{} `mapstructure:"service_config" yaml:"plugin_config" json:"plugin_config"`
}

type HandshakeConfig struct {
	// ProtocolVersion is the version that clients must match on to
	// agree they can communicate. This should match the ProtocolVersion
	// set on ClientConfig when using a plugin.
	// This field is not required if VersionedPlugins are being used in the
	// Client or Server configurations.
	ProtocolVersion uint `mapstructure:"protocol_version" yaml:"protocol_version" json:"protocol_version"`

	// MagicCookieKey and value are used as a very basic verification
	// that a plugin is intended to be launched. This is not a security
	// measure, just a UX feature. If the magic cookie doesn't match,
	// we show human-friendly output.
	MagicCookieKey   string `mapstructure:"magic_cookie_key" yaml:"magic_cookie_key" json:"magic_cookie_key"`
	MagicCookieValue string `mapstructure:"magic_cookie_value" yaml:"magic_cookie_value" json:"magic_cookie_value"`
}
