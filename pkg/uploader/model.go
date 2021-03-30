package uploader

import (
	"github.com/odpf/siren/domain"
)

type variables struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type rule struct {
	Template  string      `json:"template"`
	Status    string      `json:"enabled"`
	Variables []variables `json:"variables"`
}

type ruleYaml struct {
	ApiVersion string          `json:"apiVersion"`
	Entity     string          `json:"entity"`
	Type       string          `json:"type"`
	Namespace  string          `json:"namespace"`
	Rules      map[string]rule `json:"rules"`
}

type templatedRule struct {
	Alert       string            `json:"alert"`
	Expr        string            `json:"expr"`
	For         string            `json:"for"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type template struct {
	Name       string            `json:"name"`
	ApiVersion string            `json:"apiVersion"`
	Type       string            `json:"type"`
	Body       []templatedRule   `json:"body"`
	Tags       []string          `json:"tags"`
	Variables  []domain.Variable `json:"variables"`
}

type yamlObject struct {
	Type string `json:"type"`
}
