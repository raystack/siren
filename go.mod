module github.com/odpf/siren

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/antihax/optional v1.0.0
	github.com/go-openapi/loads v0.20.1 // indirect
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grafana/cortex-tools v0.7.2
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jeremywohl/flatten v1.0.1
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.3.0
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mcuadros/go-defaults v1.2.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/newrelic/go-agent/v3 v3.11.0
	github.com/newrelic/go-agent/v3/integrations/nrgorilla v1.1.0
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/alertmanager v0.21.1-0.20200911160112-1fdff6b3f939
	github.com/prometheus/prometheus v1.8.2-0.20201014093524-73e2ce1bd643
	github.com/purini-to/zapmw v1.1.0
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.14.1
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
	golang.org/x/sys v0.0.0-20201126233918-771906719818 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2

replace github.com/grafana/cortex-tools v0.7.2 => github.com/kevinbheda/cortex-tools v0.8.0
