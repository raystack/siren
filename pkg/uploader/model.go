package uploader

import (
	"github.com/odpf/siren/domain"
)

type variables struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type rule struct {
	Template  string      `yaml:"template"`
	Status    string      `yaml:"status"`
	Variables []variables `yaml:"variables"`
}

type ruleYaml struct {
	ApiVersion string          `yaml:"apiVersion"`
	Entity     string          `yaml:"entity"`
	Type       string          `yaml:"type"`
	Namespace  string          `yaml:"namespace"`
	Rules      map[string]rule `yaml:"rules"`
}

type templatedRule struct {
	Record      string            `yaml:"record,omitempty"`
	Alert       string            `yaml:"alert,omitempty"`
	Expr        string            `yaml:"expr"`
	For         string            `yaml:"for,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type template struct {
	Name       string            `yaml:"name"`
	ApiVersion string            `yaml:"apiVersion"`
	Type       string            `yaml:"type"`
	Body       []templatedRule   `yaml:"body"`
	Tags       []string          `yaml:"tags"`
	Variables  []domain.Variable `yaml:"variables"`
}

type yamlObject struct {
	Type string `yaml:"type"`
}
