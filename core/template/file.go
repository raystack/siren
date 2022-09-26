package template

type TemplateFile struct {
	Name       string              `yaml:"name"`
	ApiVersion string              `yaml:"apiVersion"`
	Type       string              `yaml:"type"`
	Body       []templatedRuleFile `yaml:"body"`
	Tags       []string            `yaml:"tags"`
	Variables  []Variable          `yaml:"variables"`
}

type templatedRuleFile struct {
	Record      string            `yaml:"record,omitempty"`
	Alert       string            `yaml:"alert,omitempty"`
	Expr        string            `yaml:"expr"`
	For         string            `yaml:"for,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}
