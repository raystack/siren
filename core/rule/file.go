package rule

type variables struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type RuleFile struct {
	ApiVersion        string `yaml:"apiVersion"`
	Entity            string `yaml:"entity"`
	Type              string `yaml:"type"`
	Namespace         string `yaml:"namespace"`
	Provider          string `yaml:"provider"`
	ProviderNamespace string `yaml:"providerNamespace"`
	Rules             map[string]struct {
		Template  string      `yaml:"template"`
		Enabled   bool        `yaml:"enabled"`
		Variables []variables `yaml:"variables"`
	} `yaml:"rules"`
}
